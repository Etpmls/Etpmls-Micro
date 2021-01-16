package em

import (
	"context"
	"flag"
	"github.com/Etpmls/Etpmls-Micro/define"
	library "github.com/Etpmls/Etpmls-Micro/library"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


/*
	[GRPC]
*/
// https://github.com/grpc/grpc-go/blob/15a78f19307d5faf10cfdd9d4e664c65a387cbd1/examples/helloworld/greeter_server/main.go#L46
func (this *Register) runGrpcServer()  {
	k, err := Kv.ReadKey(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceRpcPort))
	if err != nil {
		LogInfo.OutputSimplePath(err)
		panic(err)
	}

	lis, err := net.Listen("tcp", ":" + k)
	if err != nil {
		LogFatal.Output("failed to listen: " + err.Error())
	}

	s := this.GrpcMiddleware()

	// Register Service
	// 注册服务
	this.GrpcServiceFunc(s)

	if err := s.Serve(lis); err != nil {
		LogFatal.Output("failed to serve: " + err.Error())
	}
}

/*
	[HTTP]
*/
func (this *Register) runHttpServer()  {
	rp, err := Kv.ReadKey(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceRpcPort))
	if err != nil {
		LogInfo.OutputSimplePath(err)
		panic(err)
	}
	hp, err := Kv.ReadKey(define.MakeServiceConfField(library.Config.Service.RpcId, define.KvServiceHttpPort))
	if err != nil {
		LogInfo.OutputSimplePath(err)
		panic(err)
	}

	// https://github.com/grpc-ecosystem/grpc-gateway#usage
	var (
		// command-line options:
		// gRPC server endpoint
		grpcServerEndpoint = flag.String("grpc-server-endpoint",  "localhost:" + rp, "gRPC server endpoint")
	)

	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := this.GrpcEndpoint()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	// err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	err = this.HttpServiceFunc(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		glog.Fatal(err)
	}

	// Custom Route
	this.RouteFunc(mux)

	// Set CORS
	options := this.CorsOptions()

	handler := DefaultMiddleware().SetCors(mux, options)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	err = http.ListenAndServe(":" + hp, handler)
	if err != nil {
		glog.Fatal(err)
	}
}

func (this *Register) monitorExit()  {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-c
	this.HandleExitFunc()
}


