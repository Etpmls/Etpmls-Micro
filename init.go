package em

import (
	"context"
	"github.com/Etpmls/Etpmls-Micro/language"
	library "github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
)

var EA *Register

type Register struct {
	Config                    library.Configuration
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
}

func init()  {
	library.Init_Yaml()
	library.Init_Logrus()
	library.Init_Redis()
	library.Init_Consul()
	library.Init_Validator()
	library.Init_GoI18n()
	library.Init_HystrixGo()
	language.LoadLanguage()
}

func (this *Register) Run()  {
	// Set Custom Configuration
	library.Init_CustomYaml(this.CustomConfiguration.Path, this.CustomConfiguration.DebugPath, this.CustomConfiguration.StructAddr)

	if library.Config.App.Database {
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