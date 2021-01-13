package em

import (
	"context"
	"encoding/json"
	library "github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strconv"
	"strings"
	_ "github.com/Etpmls/Etpmls-Micro/file"
	"github.com/hashicorp/consul/api"
	"time"
)

var EA *Register


type Register struct {
	// APP
	// Version Infomation
	// main.exe -version
	Version                   map[string]string
	HandleExitFunc            func()
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
	if this.OverrideInitYaml == nil {
		library.Init_Yaml()
	} else {
		library.Config = *this.OverrideInitYaml()
	}

	// Connect KV
	if this.OverrideInitKv == nil {
		var conf = api.Config{
			Address:    library.Config.Kv.Address,
		}
		library.InitConsulKv(&conf)
	} else {
		Kv = *this.OverrideInitKv(library.Config.Kv.Address)
	}

	// Init Log
	level, err := Kv.ReadKey(KvLogLevel)
	if err != nil {
		library.InitLog.Println("[ERROR]", KvLogLevel, " is not configured!")
		return
	}
	if this.OverrideInitLog == nil {
		library.Init_Logrus(level)
	} else {
		Log = *this.OverrideInitLog(level)
	}

	// Init Validator
	validatorEnable, err := Kv.ReadKey(KvValidatorEnable)
	if err != nil || strings.ToLower(validatorEnable) != "true" {
		library.InitLog.Println("[WARNING]", KvValidatorEnable, " is not configured or not enable!")
	} else {
		if this.OverrideInitValidator == nil {
			library.Init_Validator()
		} else {
			Validator = *this.OverrideInitValidator()
		}
	}


	// Init I18n
	// - Check if i18n is enable
	i18nEnable, err := Kv.ReadKey(KvI18nEnable)
	if err != nil || strings.ToLower(i18nEnable) != "true" {
		library.InitLog.Println("[WARNING]", KvI18nEnable, " is not configured or not enable!")
	} else {
		i18nMap, err := Kv.List(KvI18nLanguage)
		if err != nil {
			library.InitLog.Println("[ERROR]", KvI18nLanguage, " is not configured!")
		}
		if len(i18nMap) > 0 {
			if this.OverrideInitI18n == nil {
				library.Init_GoI18n(i18nMap)
			} else {
				I18n = *this.OverrideInitI18n(i18nMap)
			}
		}
	}


	// Init CircuitBreaker
	// - Check if i18n is enable
	t, err := Kv.ReadKey(KvAppCommunicationTimeout)
	if err != nil {
		library.InitLog.Println("[WARNING]", KvAppCommunicationTimeout, " is not configured!")
		library.InitLog.Println("[INFO]", KvAppCommunicationTimeout, " is set to 5s!")
		t = "5s"
	}
	timeout, err := time.ParseDuration(t)
	if err != nil {
		library.InitLog.Println("[WARNING]","The format of " ,KvAppCommunicationTimeout, " is incorrect!")
		return
	}
	cbEnable, err := Kv.ReadKey(KvCircuitBreakerEnable)
	if err != nil || strings.ToLower(cbEnable) != "true" {
		library.InitLog.Println("[WARNING]", KvCircuitBreakerEnable, " is not configured or not enable!")
	} else {
		if this.OverrideInitCircuitBreaker == nil {
			library.Init_HystrixGo(timeout)
		} else {
			CircuitBreaker = *this.OverrideInitCircuitBreaker(timeout)
		}
	}


	// Init Cache
	// - if enable
	m, err := Kv.List(KvCache)
	if err != nil {
		library.InitLog.Println("[WARNING]", KvCache, " is not configured!")
		library.InitLog.Println("[WARNING]", "Cache is not running")
	}
	if strings.ToLower(m[KvCacheEnable]) == "true" {
		cacheDbInt, err := strconv.Atoi(m[KvCacheDb])
		if err != nil {
			library.InitLog.Println("[ERROR]", KvCacheDb, " is not number!")
			return
		}
		if len(m[KvCacheAddress]) == 0 {
			library.InitLog.Println("[ERROR]", KvCacheAddress, " is not configured!")
			return
		}

		if this.OverrideInitCache == nil {
			library.Init_Redis(true, m[KvCacheAddress], m[KvCachePassword], cacheDbInt)
		} else {
			Cache = *this.OverrideInitCache(true, m[KvCacheAddress], m[KvCachePassword], cacheDbInt)
		}
	}


	// Init service discovery
	// - if enable
	sDMap, err := Kv.List(KvServiceDiscovery)
	if err != nil || strings.ToLower(sDMap[KvServiceDiscoveryEnable]) != "true" {
		library.InitLog.Println("[WARNING]", KvServiceDiscoveryEnable, " is not configured or not enable!")
	} else {
		pkgConfig := api.DefaultConfig()
		pkgConfig.Address = sDMap[KvServiceDiscoveryServerAddress]
		pkgConfig.WaitTime = timeout

		// Set tags
		// - Get global tag
		var glbRTag, glbHTag []string
		srvIdMap, err := Kv.List(MakeServiceConfField(library.Config.Service.RpcName, "")) // service/rpcName/
		if err == nil {
			_ = json.Unmarshal([]byte(srvIdMap[MakeServiceConfField(library.Config.Service.RpcName, KvServiceRpcTag)]), &glbRTag)  // service/rpcName/rpc-tag
			_ = json.Unmarshal([]byte(srvIdMap[MakeServiceConfField(library.Config.Service.RpcName, KvServiceHttpTag)]), &glbHTag)  // service/rpcName/http-tag
		}

		var rTag []string
		err = json.Unmarshal([]byte(m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceRpcTag)]), &rTag) // service/rpcID/rpc-tag
		if err != nil {
			library.InitLog.Println("[ERROR]", "RpcTag format error!", err)
			rTag = glbRTag
		}
		if len(rTag) == 0 {
			rTag = glbRTag
		}

		var hTag []string
		err = json.Unmarshal([]byte(m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceHttpTag)]), &hTag) // service/rpcID/http-tag
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
			RpcPort:       m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceRpcPort)],
			RpcTag:        rTag,
			HttpId:        m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceHttpId)],
			HttpName:      m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceHttpName)],
			HttpPort:      m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceHttpPort)],
			HttpTag:       hTag,
			Address:       m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceAddress)],
			CheckInterval: m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceCheckInterval)],
			CheckUrl:      m[MakeServiceConfField(library.Config.Service.RpcId, KvServiceCheckUrl)],
		}

		if this.OverrideInitServiceDiscovery == nil {
			library.Init_Consul(&config)
		} else {
			this.OverrideInitServiceDiscovery(&config)
		}
	}

	// Init jwt token
	if this.OverrideInitJwtToken != nil {
		JwtToken = *this.OverrideInitJwtToken()
	}
	// Init captcha
	if this.OverrideInitCaptcha != nil {
		Captcha = *this.OverrideInitCaptcha()
	}

	Micro.Config = &library.Config
}

func (this *Register) Run()  {
	// Only print version info
	if this.stoprun {
		return
	}

	dbEnable, err := Kv.ReadKey(KvDatabaseEnable)
	if err != nil || strings.ToLower(dbEnable) != "true" {
		library.InitLog.Println("[WARNING]", KvDatabaseEnable, " is not configured or not enable!!")
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

	EA = this

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
