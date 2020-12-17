package main

import (
	"context"
	"github.com/Etpmls/Etpmls-Micro"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
)

func main()  {
	var reg = em.Register{
		GrpcServiceFunc: RegisterRpcService,
		HttpServiceFunc: RegisterHttpService,
		RouteFunc:       RegisterRoute,
	}
	reg.Init()
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
	mux.HandlePath("GET", "/hello", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("world"))
	})
}


