package em

import (
	"context"
	"encoding/json"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	_ "github.com/Etpmls/Etpmls-Micro/v2/file"
	library "github.com/Etpmls/Etpmls-Micro/v2/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strconv"
	"strings"
	"time"
)

const (
	EnableDatabase = "Database"
	EnableValidator = "Validator"
	EnableI18n = "I18n"
	EnableCircuitBreaker = "CircuitBreaker"
	EnableCache = "Cache"
	EnableServiceDiscovery = "ServiceDiscovery"
	EnableCaptcha = "Captcha"
)

var Reg *Register


type Register struct {
	// APP
	// Version Infomation
	// main.exe -version
	AppVersion            map[string]string
	AppHandleExitFunc     func()
	AppEnabledFeatureName []string
	// Rpc
	RpcServiceFunc       func(*grpc.Server)
	RpcMiddleware        func() *grpc.Server
	RpcEndpoint          func() *runtime.ServeMux
	RpcReturnSuccessFunc func(code string, message string, data interface{}) (*em_protobuf.Response, error)
	RpcReturnErrorFunc   func(rcpStatusCode codes.Code, code string, message string, data interface{}, err error) (*em_protobuf.Response, error)
	// HTTP
	HttpServiceFunc       func(ctx context.Context, mux *runtime.ServeMux, grpcServerEndpoint *string, opts []grpc.DialOption) error
	HttpRouteFunc         func(mux *runtime.ServeMux)
	HttpCorsOptions       func() cors.Options
	HttpReturnSuccessFunc func(code string, message string, data interface{}) ([]byte, error)
	HttpReturnErrorFunc   func(code string, message string, data interface{}, err error) ([]byte, error)
	// Custom Server
	CustomServerFunc                               []func()
	CustomServerServiceRegisterFunc                func() error
	CustomServerServiceExitFunc                    func() error
	// Database
	DatabaseMigrate           []interface{}
	DatabaseInsertInitialData []func()

	OverrideInitYaml		func() *library.Configuration
	OverrideInitKv	func() *Interface_KV
	OverrideInitLog	func(level string) *Interface_Log
	OverrideInitCache	func(enableCache bool, address string, password string, db int) *Interface_Cache
	OverrideInitValidator func() *Interface_Validator
	// map [*.en.json](key) {"langKey":"langValue"}(value)
	OverrideInitI18n func(langMap map[string]string) *Interface_I18n
	OverrideInitCircuitBreaker func(time.Duration) *Interface_CircuitBreaker
	OverrideInitServiceDiscovery func(*library.ServiceConfig) *Interface_ServiceDiscovery
	OverrideInitJwtToken func() *Interface_Jwt
	OverrideInitCaptcha func() *Interface_Captcha
	stoprun				  bool
	initFinished bool
}



func (this *Register) Init() {
	// Version info
	switch  {
	case init_version(this):
		this.stoprun = true
		VersionPrint()
		return
	default:
	}

	Reg = this

	// Write log to File
	i := library.NewInit()
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
	this.initI18n()

	t, err := Kv.ReadKey(define.KvAppCommunicationTimeout)
	if err != nil {
		library.InitLog.Println("[INFO][DEFAULT: " + define.DefaultAppCommunicationTimeout + "]", define.KvAppCommunicationTimeout, " is not configured!")
		t = define.DefaultAppCommunicationTimeout
	}
	timeout, err := time.ParseDuration(t)
	if err != nil {
		library.InitLog.Println("[INFO]","The format of " , define.KvAppCommunicationTimeout, " is incorrect!")
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
	if this.OverrideInitJwtToken != nil {
		JwtToken = *this.OverrideInitJwtToken()
	}
	// Init captcha
	if this.OverrideInitCaptcha != nil {
		Captcha = *this.OverrideInitCaptcha()
	}

	// Check whether the function is turned on to prevent interface nil during operation
	// 检查功能是否开启，防止运行期间出现接口nil的情况
	this.checkFeatureEnable(this.AppEnabledFeatureName)

	Micro.Config = &library.Config



	// Generate key
	this.generateKey()

	Reg.initFinished = true
}

func (this *Register) Run()  {
	// Only print version info
	if this.stoprun {
		return
	}

	dbEnable, err := Kv.ReadKey(define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceDatabaseEnable))
	if err != nil || strings.ToLower(dbEnable) != "true" {
		library.InitLog.Println("[WARNING]", define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceDatabaseEnable), " is not configured or not enable!!")
	} else {
		// Init Database
		this.RunDatabase()
		// Insert database initial data
		this.InsertDataToDatabase()
	}


	// Default Func
	if this.RpcReturnSuccessFunc == nil {
		this.RpcReturnSuccessFunc = defaultHandleRpcSuccessFunc
	}
	if this.RpcReturnErrorFunc == nil {
		this.RpcReturnErrorFunc = defaultHandleRpcErrorFunc
	}
	if this.HttpReturnSuccessFunc == nil {
		this.HttpReturnSuccessFunc = defaultHandleHttpSucessFunc
	}
	if this.HttpReturnErrorFunc == nil {
		this.HttpReturnErrorFunc = defaultHandleHttpErrorFunc
	}
	if this.RpcMiddleware == nil {
		this.RpcMiddleware = defaultGrpcMiddlewareFunc
	}
	if this.HttpCorsOptions == nil {
		this.HttpCorsOptions = defaultCorsOptions
	}
	if this.RpcEndpoint == nil {
		this.RpcEndpoint = defaultRegisterEndpoint
	}
	if this.AppHandleExitFunc == nil {
		this.AppHandleExitFunc = defaultHandleExit
	}

	// Rpc Server
	go this.runGrpcServer()
	// Http Server
	if this.HttpServiceFunc != nil {
		go this.runHttpServer()
	}
	// Other Service
	if this.CustomServerFunc != nil {
		for _, v := range this.CustomServerFunc {
			go v()
		}
	}

	this.monitorExit()
}

// Insert database initial data
// 插入数据库初始数据
func (this *Register) InsertDataToDatabase()  {
	env, err := godotenv.Read("./.env")
	if err != nil {
		library.Instance_Logrus.Error(err)
		return
	}

	if _, ok := env["INIT_DATABASE"]; ok {
		if strings.ToUpper(env["INIT_DATABASE"]) == "TRUE" {
			for _, v := range this.DatabaseInsertInitialData {
				v()
			}
			env["INIT_DATABASE"] = "FALSE"
		}
	}

	err = godotenv.Write(env, "./.env")
	if err != nil {
		library.Instance_Logrus.Error(err)
		return
	}
}

func (this *Register) logAndChangeJsonFormatIfMapValueEmpty(key string, m map[string]string) string {
	if len(m[key]) == 0 {
		library.InitLog.Println("[WARNING]", key, " is not configured!")
	}
	return IfEmptyChangeJsonFormat(m[key])
}

func (this *Register) logIfMapValueEmpty(key string, m map[string]string) string {
	if len(m[key]) == 0 {
		library.InitLog.Println("[WARNING]", key, " is not configured!")
	}
	return m[key]
}

func (this *Register) initYaml() {
	if this.OverrideInitYaml == nil {
		library.Init_Yaml()
	} else {
		library.Config = *this.OverrideInitYaml()
	}
}

func (this *Register) initKv() {
	if len(library.Config.Kv.Address) == 0 {
		library.InitLog.Println("[ERROR]", "Kv address is not configured!")
		panic("Kv address is not configured!")
	}
	idx := GetRandomNumberByLength(len(library.Config.Kv.Address))
	if this.OverrideInitKv == nil {
		var conf = api.Config{
			Address:    library.Config.Kv.Address[idx],
			Token:      library.Config.Kv.Token,
		}
		library.InitConsulKv(&conf)
	} else {
		Kv = *this.OverrideInitKv()
	}
}

func (this *Register) initLog() {
	level, err := Kv.ReadKey(define.KvLogLevel)
	if err != nil {
		library.InitLog.Println("[INFO][DEFAULT: " + define.DefaultLogLevel + "]", define.KvLogLevel, " is not configured!")
		level = "info"
	}
	if this.OverrideInitLog == nil {
		library.Init_Logrus(level)
	} else {
		Log = *this.OverrideInitLog(level)
	}
}

func (this *Register) initValidator() {
	validatorEnable, err := Kv.ReadKey(define.KvValidatorEnable)
	if err != nil || strings.ToLower(validatorEnable) != "true" {
		library.InitLog.Println("[WARNING]", define.KvValidatorEnable, " is not configured or not enable!")
	} else {
		if this.OverrideInitValidator == nil {
			library.Init_Validator()
		} else {
			Validator = *this.OverrideInitValidator()
		}
	}
}

func (this *Register) initI18n() {
	// - Check if i18n is enable
	i18nEnable, err := Kv.ReadKey(define.KvI18nEnable)
	if err != nil || strings.ToLower(i18nEnable) != "true" {
		library.InitLog.Println("[WARNING]", define.KvI18nEnable, " is not configured or not enable!")
	} else {
		i18nMap, err := Kv.List(define.KvI18nLanguage)
		if err != nil || len(i18nMap) == 0 {
			library.InitLog.Println("[ERROR]", define.KvI18nLanguage, " is not configured!")
			panic("[ERROR]"+ define.KvI18nLanguage+ " is not configured!")
		}
		if this.OverrideInitI18n == nil {
			library.Init_GoI18n(i18nMap)
		} else {
			I18n = *this.OverrideInitI18n(i18nMap)
		}
	}
}

func (this *Register) initCircuitBreaker(timeout time.Duration) {
	cbEnable, err := Kv.ReadKey(define.KvCircuitBreakerEnable)
	if err != nil || strings.ToLower(cbEnable) != "true" {
		library.InitLog.Println("[WARNING]", define.KvCircuitBreakerEnable, " is not configured or not enable!")
	} else {
		if this.OverrideInitCircuitBreaker == nil {
			library.Init_HystrixGo(timeout)
		} else {
			CircuitBreaker = *this.OverrideInitCircuitBreaker(timeout)
		}
	}
}

func (this *Register) initCache() {
	// - if enable
	m, err := Kv.List(define.KvCache)
	if err != nil {
		library.InitLog.Println("[WARNING]", define.KvCache, " is not configured!")
		library.InitLog.Println("[WARNING]", "Cache is not running")
	}
	if strings.ToLower(m[define.KvCacheEnable]) == "true" {
		cacheDbInt, err := strconv.Atoi(m[define.KvCacheDb])
		if err != nil {
			library.InitLog.Println("[ERROR]", define.KvCacheDb, " is not number!")
			panic("[ERROR]"+ define.KvCacheDb+ " is not number!")
		}
		if len(m[define.KvCacheAddress]) == 0 {
			library.InitLog.Println("[ERROR]", define.KvCacheAddress, " is not configured!")
			panic("[ERROR]"+ define.KvCacheAddress+ " is not configured!")
		}

		if this.OverrideInitCache == nil {
			library.Init_Redis(true, m[define.KvCacheAddress], m[define.KvCachePassword], cacheDbInt)
		} else {
			Cache = *this.OverrideInitCache(true, m[define.KvCacheAddress], m[define.KvCachePassword], cacheDbInt)
		}
	}
}

func (this *Register) initServiceDiscovery(timeout time.Duration)  {
	// - if enable
	s, err := Kv.ReadKey(define.KvServiceDiscoveryEnable)
	if err != nil || strings.ToLower(s) != "true" {
		library.InitLog.Println("[WARNING]", define.KvServiceDiscoveryEnable, " is not configured or not enable!")
	} else {
		pkgConfig := api.DefaultConfig()
		pkgConfig.Address = MustGetKvKey(define.KvServiceDiscoveryAddress)
		pkgConfig.WaitTime = timeout
		pkgConfig.Token = library.Config.Kv.Token

		// Set tags
		rt, err := GetServiceKvKey(define.KvServiceRpcTag)
		if err != nil {
			rt = "[]"
		}
		ht, err := GetServiceKvKey(define.KvServiceHttpTag)
		if err != nil {
			ht = "[]"
		}
		var rTag, hTag []string
		_ = json.Unmarshal([]byte(rt), &rTag)
		_ = json.Unmarshal([]byte(ht), &hTag)

		var config = library.ServiceConfig{
			Config:        pkgConfig,
			RpcId:         library.Config.Service.RpcId,
			RpcName:       library.Config.Service.RpcName,
			RpcPort:       MustGetServiceIdKvKey(define.KvServiceRpcPort),
			RpcTag:        rTag,
			HttpId:        MustGetServiceIdKvKey(define.KvServiceHttpId),
			HttpName:      MustGetServiceKvKey(define.KvServiceHttpName),
			HttpPort:      MustGetServiceIdKvKey(define.KvServiceHttpPort),
			HttpTag:       hTag,
			Address:       MustGetServiceIdKvKey(define.KvServiceAddress),
			CheckInterval: MustGetServiceKvKey(define.KvServiceCheckInterval),
			CheckUrl:      MustGetServiceIdKvKey(define.KvServiceCheckUrl),
		}

		if this.OverrideInitServiceDiscovery == nil {
			library.Init_Consul(&config)
		} else {
			ServiceDiscovery = *this.OverrideInitServiceDiscovery(&config)
		}

		if this.CustomServerServiceRegisterFunc != nil {
			err := this.CustomServerServiceRegisterFunc()
			if err != nil {
				go this.reRegisterCustomServerService()
			}
		}

	}
}

// Check whether the function is turned on to prevent interface nil during operation
// 检查功能是否开启，防止运行期间出现接口nil的情况
func (this *Register) checkFeatureEnable(feature []string) {
	for _, v := range feature {
		if v == EnableDatabase {
			e, err := Kv.ReadKey(define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceDatabaseEnable))
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable " + define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceDatabaseEnable))
			}
		}
		if v == EnableValidator {
			e, err := Kv.ReadKey(define.KvValidatorEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable validator/enable")
			}
		}
		if v == EnableI18n {
			e, err := Kv.ReadKey(define.KvI18nEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable i18n/enable")
			}
		}
		if v == EnableCircuitBreaker {
			e, err := Kv.ReadKey(define.KvCircuitBreakerEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable circuit-breaker/enable")
			}
		}
		if v == EnableCache {
			e, err := Kv.ReadKey(define.KvCacheEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable cache/enable")
			}
		}
		if v == EnableServiceDiscovery {
			e, err := Kv.ReadKey(define.KvServiceDiscoveryEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable service-discovery/enable")
			}
		}
		if v == EnableCaptcha {
			e, err := Kv.ReadKey(define.KvCaptchaEnable)
			if err != nil || (strings.ToLower(e) != "true" && strings.ToLower(e) != "false") {
				panic("Please enable captcha/enable")
			}
		}
	}
}

func (this *Register) generateKey() {
	k, err := Kv.ReadKey(define.KvAppKey)
	if err != nil || len(k) == 0 {
		err := Kv.CrateOrUpdateKey(define.KvAppKey, GenerateRandomString(50))
		if err != nil {
			panic(err)
		}
	}
}

// When initial registration fails, automatically retry registration
// 当初始化注册失败时，自动重试注册
func (this *Register) reRegisterCustomServerService() {
	for {
		time.Sleep(time.Second * 5)
		err := this.CustomServerServiceRegisterFunc()
		if err == nil {
			library.InitLog.Println("[INFO]", "Custom Server Service registered successfully!")
			break
		}
	}
}