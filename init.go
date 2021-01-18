package em

import (
	"context"
	"encoding/json"
	"github.com/Etpmls/Etpmls-Micro/define"
	_ "github.com/Etpmls/Etpmls-Micro/file"
	library "github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
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

var EA *Register


type Register struct {
	// APP
	// Version Infomation
	// main.exe -version
	Version                   map[string]string
	HandleExitFunc            func()
	EnabledFeatureName		[]string
	stoprun				  bool
	// Grpc
	GrpcServiceFunc           func(*grpc.Server)
	GrpcMiddleware            func() *grpc.Server
	GrpcEndpoint              func() *runtime.ServeMux
	ReturnRpcSuccessFunc      func(code string, message string, data interface{}) (*em_protobuf.Response, error)
	ReturnRpcErrorFunc    func(rcpStatusCode codes.Code, code string, message string, data interface{}, err error) (*em_protobuf.Response, error)
	// HTTP
	HttpServiceFunc           func(ctx context.Context, mux *runtime.ServeMux, grpcServerEndpoint *string, opts []grpc.DialOption) error
	RouteFunc                 func(mux *runtime.ServeMux)
	CorsOptions               func() cors.Options
	ReturnHttpSuccessFunc func(code string, message string, data interface{}) ([]byte, error)
	ReturnHttpErrorFunc   func(code string, message string, data interface{}, err error) ([]byte, error)
	// Database
	DatabaseMigrate           []interface{}
	InsertDatabaseInitialData []func()

	OverrideInitYaml		func() *library.Configuration
	OverrideInitKv	func(address string) *Interface_KV
	OverrideInitLog	func(level string) *Interface_Log
	OverrideInitCache	func(enableCache bool, address string, password string, db int) *Interface_Cache
	OverrideInitValidator func() *Interface_Validator
	// map [*.en.json](key) {"langKey":"langValue"}(value)
	OverrideInitI18n func(langMap map[string]string) *Interface_I18n
	OverrideInitCircuitBreaker func(time.Duration) *Interface_CircuitBreaker
	OverrideInitServiceDiscovery func(*library.ServiceConfig) *Interface_ServiceDiscovery
	OverrideInitJwtToken func() *Interface_Jwt
	OverrideInitCaptcha func() *Interface_Captcha
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
	this.checkFeatureEnable(this.EnabledFeatureName)

	Micro.Config = &library.Config
	EA = this


	// Generate key
	this.generateKey()
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
	if this.ReturnRpcSuccessFunc == nil {
		this.ReturnRpcSuccessFunc = defaultHandleRpcSuccessFunc
	}
	if this.ReturnRpcErrorFunc == nil {
		this.ReturnRpcErrorFunc = defaultHandleRpcErrorFunc
	}
	if this.ReturnHttpSuccessFunc == nil {
		this.ReturnHttpSuccessFunc = defaultHandleHttpSucessFunc
	}
	if this.ReturnHttpErrorFunc == nil {
		this.ReturnHttpErrorFunc = defaultHandleHttpErrorFunc
	}
	if this.GrpcMiddleware == nil {
		this.GrpcMiddleware = defaultGrpcMiddlewareFunc
	}
	if this.CorsOptions == nil {
		this.CorsOptions = defaultCorsOptions
	}
	if this.GrpcEndpoint == nil {
		this.GrpcEndpoint = defaultRegisterEndpoint
	}
	if this.HandleExitFunc == nil {
		this.HandleExitFunc = defaultHandleExit
	}


	go this.runGrpcServer()
	if this.HttpServiceFunc != nil {
		go this.runHttpServer()
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
			for _, v := range this.InsertDatabaseInitialData {
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

func (this *Register) panicIfMapValueEmpty(key string, m map[string]string) string {
	if len(m[key]) == 0 {
		library.InitLog.Println("[ERROR]", key, " is not configured!")
		panic(("[ERROR]"+ key+ " is not configured!"))
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
		}
		library.InitConsulKv(&conf)
	} else {
		Kv = *this.OverrideInitKv(library.Config.Kv.Address[idx])
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
	sDMap, err := Kv.List(define.KvServiceDiscovery)
	if err != nil || strings.ToLower(sDMap[define.KvServiceDiscoveryEnable]) != "true" {
		library.InitLog.Println("[WARNING]", define.KvServiceDiscoveryEnable, " is not configured or not enable!")
	} else {
		pkgConfig := api.DefaultConfig()
		pkgConfig.Address = this.panicIfMapValueEmpty(define.KvServiceDiscoveryAddress, sDMap)
		pkgConfig.WaitTime = timeout

		// Set tags
		// - Get global tag
		var glbRTag, glbHTag []string
		srvNameMap, err := Kv.List(define.MakeServiceConfField(library.Config.Service.RpcName, "/")) // service/rpcName/
		if err == nil {
			rtKey := define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceRpcTag)
			htKey := define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceHttpTag)
			_ = json.Unmarshal([]byte(this.logAndChangeJsonFormatIfMapValueEmpty(rtKey, srvNameMap)), &glbRTag) // service/rpcName/rpc-tag
			_ = json.Unmarshal([]byte(this.logAndChangeJsonFormatIfMapValueEmpty(htKey, srvNameMap)), &glbHTag) // service/rpcName/http-tag
		}

		// - Get current service node tag
		srvIdMap, err := Kv.List(define.MakeServiceConfField(library.Config.Service.RpcId, "/")) // service/rpcId/
		var rTag []string
		rtKey2 := define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceRpcTag)
		err = json.Unmarshal([]byte(this.logAndChangeJsonFormatIfMapValueEmpty(rtKey2, srvIdMap)), &rTag) // service/rpcID/rpc-tag
		if err != nil {
			library.InitLog.Println("[ERROR]", "RpcTag format error!", err)
			rTag = glbRTag
		}
		if len(rTag) == 0 {
			rTag = glbRTag
		}

		var hTag []string
		htKey2 := define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceHttpTag)
		err = json.Unmarshal([]byte(this.logAndChangeJsonFormatIfMapValueEmpty(htKey2, srvIdMap)), &hTag) // service/rpcID/http-tag
		if err != nil {
			library.InitLog.Println("[ERROR]", "HttpTag format error!", err)
			hTag = glbHTag
		}
		if len(hTag) == 0 {
			hTag = glbHTag
		}


		var config = library.ServiceConfig{
			Config: pkgConfig,
			RpcId:         library.Config.Service.RpcId,
			RpcName:       library.Config.Service.RpcName,
			RpcPort:       this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceRpcPort), srvIdMap),
			RpcTag:        rTag,
			HttpId:        this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceHttpId), srvIdMap),
			HttpName:      this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceHttpName), srvIdMap),
			HttpPort:      this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceHttpPort), srvIdMap),
			HttpTag:       hTag,
			Address:       this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceAddress), srvIdMap),
			CheckInterval: this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceCheckInterval), srvIdMap),
			CheckUrl:      this.panicIfMapValueEmpty(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceCheckUrl), srvIdMap),
		}

		if this.OverrideInitServiceDiscovery == nil {
			library.Init_Consul(&config)
		} else {
			this.OverrideInitServiceDiscovery(&config)
		}
	}
}

// Check whether the function is turned on to prevent interface nil during operation
// 检查功能是否开启，防止运行期间出现接口nil的情况
func (this *Register) checkFeatureEnable(feature []string) {
	for _, v := range feature {
		if v == EnableDatabase {
			e, err := Kv.ReadKey(define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceDatabaseEnable))
			if err != nil || strings.ToLower(e) != "true" {
				panic("Please enable " + define.MakeServiceConfField(library.Config.Service.RpcName, define.KvServiceDatabaseEnable))
			}
		}
		if v == EnableValidator {
			e, err := Kv.ReadKey(define.KvValidatorEnable)
			if err != nil || strings.ToLower(e) != "true" {
				panic("Please enable validator/enable")
			}
		}
		if v == EnableI18n {
			e, err := Kv.ReadKey(define.KvI18nEnable)
			if err != nil || strings.ToLower(e) != "true" {
				panic("Please enable i18n/enable")
			}
		}
		if v == EnableCircuitBreaker {
			e, err := Kv.ReadKey(define.KvCircuitBreakerEnable)
			if err != nil || strings.ToLower(e) != "true" {
				panic("Please enable circuit-breaker/enable")
			}
		}
		if v == EnableCache {
			e, err := Kv.ReadKey(define.KvCacheEnable)
			if err != nil || strings.ToLower(e) != "true" {
				panic("Please enable cache/enable")
			}
		}
		if v == EnableServiceDiscovery {
			e, err := Kv.ReadKey(define.KvServiceDiscoveryEnable)
			if err != nil || strings.ToLower(e) != "true" {
				panic("Please enable service-discovery/enable")
			}
		}
		if v == EnableCaptcha {
			e, err := Kv.ReadKey(define.KvCaptchaEnable)
			if err != nil || strings.ToLower(e) != "true" {
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