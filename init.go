package em

import (
	"encoding/json"
	"github.com/Etpmls/Etpmls-Micro/v3/define"
	_ "github.com/Etpmls/Etpmls-Micro/v3/file"
	em_library "github.com/Etpmls/Etpmls-Micro/v3/library"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Reg *Register
	//Judge whether the initialization is completed, if the initialization is not completed, use InitLog to record the log, otherwise use LogError to record the log
	// 判断是否初始化完成，如果初始化未完成，使用InitLog记录日志，否则使用LogError记录日志
	initFinished bool
)

type Register struct {
	// Version Infomation
	// main.exe -version
	Version map[string]string

	// Set the enabled function to prevent the program from using an empty interface to report a nil error.
	// If during the initialization process, it is checked that the relevant parameters are not configured in the KV table, it will panic the program and prompt the user to configure the information
	// 设置启用的功能，防止程序使用空的接口出现nil报错。
	// 如果在初始化过程中，检查到KV表中未配置相关参数，会panic程序并提示要求用户配置的信息
	EnabledFeature []string

	// Register RPC service and middleware
	// 注册RPC服务和中间件
	RegisterService    func(*grpc.Server)
	RegisterMiddleware func() *grpc.Server

	// Override Interface & Function
	// 重写接口和函数
	OverrideInterface OverrideInterface
	OverrideFunction OverrideFunction
}

type OverrideInterface struct {
	Yaml		func() *em_library.Configuration
	Kv	func() *Interface_KV
	Log	func(level string) *Interface_Log
	Cache	func(enableCache bool, address string, password string, db int) *Interface_Cache
	Validator func() *Interface_Validator
	Translate func(langMap map[string]string) *Interface_Translate 	// map [*.en.json](key) {"langKey":"langValue"}(value)
	CircuitBreaker func(time.Duration) *Interface_CircuitBreaker
	ServiceDiscovery func(*em_library.ServiceConfig) *Interface_ServiceDiscovery
	JwtToken func() *Interface_Jwt
	Captcha func() *Interface_Captcha
}

type OverrideFunction struct {
	HandleExitFunc func()
}

func (this *Register) Init() {
	// Print Version info
	if init_version(this) {
		VersionPrint()
		os.Exit(1)
	}

	// The configuration is registered in global variables for easy program reference
	// 配置注册到全局变量中，方便程序引用
	Reg = this

	// Write log to File
	// 写日志到文件中
	i := em_library.NewInit()
	i.Start()
	defer i.Close()

	// Yaml config - Get KV address
	this.initYaml()

	// Connect KV
	this.initKv()

	// Init Log
	this.initLog()

	// Init Validator
	this.initValidator()

	// Init I18n
	this.initTranslate()

	t, err := Kv.ReadKey(em_define.KvAppCommunicationTimeout)
	if err != nil {
		em_library.InitLog.Println("[INFO][DEFAULT: " +em_define.DefaultAppCommunicationTimeout+ "]", em_define.KvAppCommunicationTimeout, " is not configured!")
		t = em_define.DefaultAppCommunicationTimeout
	}
	timeout, err := time.ParseDuration(t)
	if err != nil {
		em_library.InitLog.Println("[INFO]","The format of " , em_define.KvAppCommunicationTimeout, " is incorrect!")
		return
	}

	// Init CircuitBreaker
	this.initCircuitBreaker(timeout)

	// Init Cache
	this.initCache()

	// Init service discovery
	this.initServiceDiscovery(timeout)

	/*
		If the user registers a custom method to implement the interface, then the interface will use the user-registered custom method to override the default method
		如果用户注册了自定义方法实现了接口，那么接口将使用用户注册的自定义方法，覆盖默认方法
	*/
	// Init jwt token
	if this.OverrideInterface.JwtToken != nil {
		JwtToken = *this.OverrideInterface.JwtToken()
	}
	// Init captcha
	if this.OverrideInterface.Captcha != nil {
		Captcha = *this.OverrideInterface.Captcha()
	}

	// Check whether the function is turned on to prevent interface nil during operation
	// 检查功能是否开启，防止运行期间出现接口nil的情况
	this.checkFeatureEnable(this.EnabledFeature)

	Micro.Config = &em_library.Config

	// Generate key
	this.generateKey()

	// Complete initialization
	// 完成初始化
	initFinished = true
}

func (this *Register) Run()  {
	// Default Func
	if this.RegisterMiddleware == nil {
		this.RegisterMiddleware = defaultGrpcMiddlewareFunc
	}
	if this.OverrideFunction.HandleExitFunc == nil {
		this.OverrideFunction.HandleExitFunc = defaultHandleExit
	}

	// Rpc Server
	go this.runGrpcServer()

	this.monitorExit()
}

func (this *Register) logAndChangeJsonFormatIfMapValueEmpty(key string, m map[string]string) string {
	if len(m[key]) == 0 {
		em_library.InitLog.Println("[WARNING]", key, " is not configured!")
	}
	return IfEmptyChangeJsonFormat(m[key])
}

func (this *Register) logIfMapValueEmpty(key string, m map[string]string) string {
	if len(m[key]) == 0 {
		em_library.InitLog.Println("[WARNING]", key, " is not configured!")
	}
	return m[key]
}

func (this *Register) initYaml() {
	if this.OverrideInterface.Yaml == nil {
		em_library.Init_Yaml()
	} else {
		em_library.Config = *this.OverrideInterface.Yaml()
	}
}

func (this *Register) initKv() {
	if len(em_library.Config.Kv.Address) == 0 {
		em_library.InitLog.Println("[ERROR]", "Kv address is not configured!")
		panic("Kv address is not configured!")
	}
	idx := GetRandomNumberByLength(len(em_library.Config.Kv.Address))
	if this.OverrideInterface.Kv == nil {
		var conf = api.Config{
			Address:    em_library.Config.Kv.Address[idx],
			Token:      em_library.Config.Kv.Token,
		}
		em_library.InitConsulKv(&conf)
	} else {
		Kv = *this.OverrideInterface.Kv()
	}
}

func (this *Register) initLog() {
	level, err := Kv.ReadKey(em_define.KvLogLevel)
	if err != nil {
		em_library.InitLog.Println("[INFO][DEFAULT: " +em_define.DefaultLogLevel+ "]", em_define.KvLogLevel, " is not configured!")
		level = "info"
	}
	if this.OverrideInterface.Log == nil {
		em_library.Init_Logrus(level)
	} else {
		Log = *this.OverrideInterface.Log(level)
	}
}

func (this *Register) initValidator() {
	validatorEnable, err := Kv.ReadKey(em_define.KvValidatorEnable)
	if err != nil || strings.ToLower(validatorEnable) != "true" {
		em_library.InitLog.Println("[WARNING]", em_define.KvValidatorEnable, " is not configured or not enable!")
	} else {
		if this.OverrideInterface.Validator == nil {
			em_library.Init_Validator()
		} else {
			Validator = *this.OverrideInterface.Validator()
		}
	}
}

func (this *Register) initTranslate() {
	// - Check if i18n is enable
	i18nEnable, err := Kv.ReadKey(em_define.KvTranslateEnable)
	if err != nil || strings.ToLower(i18nEnable) != "true" {
		em_library.InitLog.Println("[WARNING]", em_define.KvTranslateEnable, " is not configured or not enable!")
	} else {
		i18nMap, err := Kv.List(em_define.KvTranslateLanguage)
		if err != nil || len(i18nMap) == 0 {
			em_library.InitLog.Println("[ERROR]", em_define.KvTranslateLanguage, " is not configured!")
			panic("[ERROR]"+ em_define.KvTranslateLanguage + " is not configured!")
		}
		if this.OverrideInterface.Translate == nil {
			em_library.Init_GoI18n(i18nMap)
		} else {
			Translate = *this.OverrideInterface.Translate(i18nMap)
		}
	}
}

func (this *Register) initCircuitBreaker(timeout time.Duration) {
	cbEnable, err := Kv.ReadKey(em_define.KvCircuitBreakerEnable)
	if err != nil || strings.ToLower(cbEnable) != "true" {
		em_library.InitLog.Println("[WARNING]", em_define.KvCircuitBreakerEnable, " is not configured or not enable!")
	} else {
		if this.OverrideInterface.CircuitBreaker == nil {
			em_library.Init_HystrixGo(timeout)
		} else {
			CircuitBreaker = *this.OverrideInterface.CircuitBreaker(timeout)
		}
	}
}

func (this *Register) initCache() {
	// - if enable
	m, err := Kv.List(em_define.KvCache)
	if err != nil {
		em_library.InitLog.Println("[WARNING]", em_define.KvCache, " is not configured!")
		em_library.InitLog.Println("[WARNING]", "Cache is not running")
	}
	if strings.ToLower(m[em_define.KvCacheEnable]) == "true" {
		cacheDbInt, err := strconv.Atoi(m[em_define.KvCacheDb])
		if err != nil {
			em_library.InitLog.Println("[ERROR]", em_define.KvCacheDb, " is not number!")
			panic("[ERROR]"+ em_define.KvCacheDb + " is not number!")
		}
		if len(m[em_define.KvCacheAddress]) == 0 {
			em_library.InitLog.Println("[ERROR]", em_define.KvCacheAddress, " is not configured!")
			panic("[ERROR]"+ em_define.KvCacheAddress + " is not configured!")
		}

		if this.OverrideInterface.Cache == nil {
			em_library.Init_Redis(true, m[em_define.KvCacheAddress], m[em_define.KvCachePassword], cacheDbInt)
		} else {
			Cache = *this.OverrideInterface.Cache(true, m[em_define.KvCacheAddress], m[em_define.KvCachePassword], cacheDbInt)
		}
	}
}

func (this *Register) initServiceDiscovery(timeout time.Duration)  {
	// - if enable
	s, err := Kv.ReadKey(em_define.KvServiceDiscoveryEnable)
	if err != nil || strings.ToLower(s) != "true" {
		em_library.InitLog.Println("[WARNING]", em_define.KvServiceDiscoveryEnable, " is not configured or not enable!")
	} else {
		pkgConfig := api.DefaultConfig()
		pkgConfig.Address = MustGetKvKey(em_define.KvServiceDiscoveryAddress)
		pkgConfig.WaitTime = timeout
		pkgConfig.Token = em_library.Config.Kv.Token

		// Set tags
		rt, err := GetServiceKvKey(em_define.KvServiceTag)
		if err != nil {
			rt = "[]"
		}

		var rTag []string
		_ = json.Unmarshal([]byte(rt), &rTag)

		var config = em_library.ServiceConfig{
			Config:  pkgConfig,
			Id:      em_library.Config.Service.RpcId,
			Name:    em_library.Config.Service.RpcName,
			Port:    MustGetServiceIdKvKey(em_define.KvServicePort),
			Tag:     rTag,
			Address: MustGetServiceIdKvKey(em_define.KvServiceAddress),
		}

		if this.OverrideInterface.ServiceDiscovery == nil {
			em_library.Init_Consul(&config)
		} else {
			ServiceDiscovery = *this.OverrideInterface.ServiceDiscovery(&config)
		}

	}
}

// Check whether the function is turned on to prevent interface nil during operation
// 检查功能是否开启，防止运行期间出现接口nil的情况
func (this *Register) checkFeatureEnable(feature []string) {
	for _, v := range feature {
		if v == em_define.EnableValidator {
			e, err := Kv.ReadKey(em_define.KvValidatorEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable validator/enable")
			}
		}
		if v == em_define.EnableTranslate {
			e, err := Kv.ReadKey(em_define.KvTranslateEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable i18n/enable")
			}
		}
		if v == em_define.EnableCircuitBreaker {
			e, err := Kv.ReadKey(em_define.KvCircuitBreakerEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable circuit-breaker/enable")
			}
		}
		if v == em_define.EnableCache {
			e, err := Kv.ReadKey(em_define.KvCacheEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable cache/enable")
			}
		}
		if v == em_define.EnableServiceDiscovery {
			e, err := Kv.ReadKey(em_define.KvServiceDiscoveryEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable service-discovery/enable")
			}
		}
		if v == em_define.EnableCaptcha {
			e, err := Kv.ReadKey(em_define.KvCaptchaEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable captcha/enable")
			}
		}
	}
}

func (this *Register) generateKey() {
	k, err := Kv.ReadKey(em_define.KvAppKey)
	if err != nil || len(k) == 0 {
		err := Kv.CrateOrUpdateKey(em_define.KvAppKey, GenerateRandomString(50))
		if err != nil {
			panic(err)
		}
	}
}

/*
	[GRPC]
*/
// https://github.com/grpc/grpc-go/blob/15a78f19307d5faf10cfdd9d4e664c65a387cbd1/examples/helloworld/greeter_server/main.go#L46
func (this *Register) runGrpcServer()  {
	k, err := Kv.ReadKey(em_define.GetPathByFieldName(em_library.Config.Service.RpcId, em_define.KvServicePort))
	if err != nil {
		LogInfo.Path(err)
		panic(err)
	}

	lis, err := net.Listen("tcp", ":" + k)
	if err != nil {
		LogFatal.New("failed to listen: " + err.Error())
	}

	s := this.RegisterMiddleware()

	// Register Service
	// 注册服务
	this.RegisterService(s)

	if err := s.Serve(lis); err != nil {
		LogFatal.New("failed to serve: " + err.Error())
	}
}


func (this *Register) monitorExit()  {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-c
	this.OverrideFunction.HandleExitFunc()
}