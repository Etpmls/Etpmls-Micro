package em

import (
	"context"
	library "github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
	_ "github.com/Etpmls/Etpmls-Micro/file"
	consul "github.com/hashicorp/consul/api"
)

var EA *Register

type Register struct {
	Version_Service map[string]string
	CustomConfiguration       struct{
		Path string
		DebugPath string
		StructAddr interface{}
	}
	GrpcServiceFunc           func(*grpc.Server)
	HttpServiceFunc           func(ctx context.Context, mux *runtime.ServeMux, grpcServerEndpoint *string, opts []grpc.DialOption) error
	RouteFunc                 func(mux *runtime.ServeMux)
	AuthServiceName       string
	DatabaseMigrate           []interface{}
	InsertDatabaseInitialData []func()
	HandleExitFunc        func()
	GrpcMiddleware        func() *grpc.Server
	GrpcEndpoint          func() *runtime.ServeMux
	CorsOptions           func() cors.Options
	ReturnRpcSuccessFunc  func(code string, message string, data interface{}) (*em_protobuf.Response, error)
	ReturnRpcErrorFunc    func(rcpStatusCode codes.Code, code string, message string, data interface{}, err error) (*em_protobuf.Response, error)
	ReturnHttpSuccessFunc func(code string, message string, data interface{}) ([]byte, error)
	ReturnHttpErrorFunc   func(code string, message string, data interface{}, err error) ([]byte, error)
	stoprun				  bool
}


func (this *Register) Init() {
	switch  {
	case init_version(this):
		this.stoprun = true
		VersionPrint()
		return
	default:
	}

	i := library.NewInit()
	i.Start()
	defer i.Close()

	library.Init_Yaml()
	library.Init_Logrus(library.Config.Log.Level)
	library.Init_Redis(library.Config.App.EnableCache, library.Config.Cache.Address, library.Config.Cache.Password, library.Config.Cache.DB)

	// Consul Config
	pkgConfig := consul.DefaultConfig()
	pkgConfig.Address = library.Config.ServiceDiscovery.Address
	pkgConfig.WaitTime = library.Config.App.CommunicationTimeout
	var config = library.ConsulConfig{
		Config: pkgConfig,
		Enable: library.Config.App.EnableServiceDiscovery,
		ConsulAddress:  library.Config.ServiceDiscovery.Address,
		RpcId:          library.Config.ServiceDiscovery.Service.Rpc.Id,
		RpcName:        library.Config.ServiceDiscovery.Service.Rpc.Name,
		RpcPort:        library.Config.App.RpcPort,
		RpcTag:         library.Config.ServiceDiscovery.Service.Rpc.Tag,
		HttpId:         library.Config.ServiceDiscovery.Service.Http.Id,
		HttpName:       library.Config.ServiceDiscovery.Service.Http.Name,
		HttpPort:       library.Config.App.HttpPort,
		HttpTag:        library.Config.ServiceDiscovery.Service.Http.Tag,
		Prefix:         library.Config.ServiceDiscovery.Service.Prefix,
		ServiceAddress: library.Config.ServiceDiscovery.Service.Address,
		CheckInterval:  library.Config.ServiceDiscovery.Service.CheckInterval,
		CheckUrl:       library.Config.ServiceDiscovery.Service.CheckUrl,
	}
	library.Init_Consul(&config)

	library.Init_Validator()
	library.Init_GoI18n()
	library.Init_HystrixGo(library.Config.App.CommunicationTimeout)
}

func (this *Register) Run()  {
	if this.stoprun {
		return
	}
	// Set Custom Configuration
	library.Init_CustomYaml(this.CustomConfiguration.Path, this.CustomConfiguration.DebugPath, this.CustomConfiguration.StructAddr)

	if library.Config.App.EnableDatabase {
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
	if len(this.AuthServiceName) == 0 {
		this.AuthServiceName = defaultAuthServiceName
	}

	go this.runGrpcServer()
	go this.runHttpServer()

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