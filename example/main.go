package main

import (
	"context"
	"github.com/Etpmls/Etpmls-Micro"
	"github.com/Etpmls/Etpmls-Micro/library"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
)

func main()  {
	var reg = em.Register{
		Config:          em_library.Config,
		GrpcServiceFunc: RegisterRpcService,
		GrpcMiddleware:  RegisterGrpcMiddleware,
		HttpServiceFunc: RegisterHttpService,
		RouteFunc:       RegisterRoute,
	}

	reg.Run()
}

// Register Rpc Service
func RegisterRpcService(s *grpc.Server)  {
	// protobuf.RegisterUserServer(s, &service.ServiceUser{})
	return
}

// Register Http Service
func RegisterHttpService(ctx context.Context, mux *runtime.ServeMux, grpcServerEndpoint *string, opts []grpc.DialOption) error {
	/*err := protobuf.RegisterUserHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}*/
	return nil
}

// Register Route
func RegisterRoute(mux *runtime.ServeMux)  {
	mux.HandlePath("GET", em_library.Config.ServiceDiscovery.Service.CheckUrl, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("hello"))
	})
}


// Register GRPC middleware
// 注册GRPC中间件
func RegisterGrpcMiddleware() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// Panic recover
				grpc_recovery.UnaryServerInterceptor(),
				// token auth
				// grpc_auth.UnaryServerInterceptor(middleware.NewAuth().BasicVerify),
				// middleware.NewMiddleware().Auth(),
				// Captcha auth
				// middleware.NewMiddleware().Captcha(),
				// I18n
				em.NewMiddleware().I18n(),
			),
		),
	)
	return s
}

