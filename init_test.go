package em

import (
	"context"
	em_define "github.com/Etpmls/Etpmls-Micro/v3/define"
	helloworld "github.com/Etpmls/Etpmls-Micro/v3/test/proto_test/pb_test"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestInit(t *testing.T)  {
	var reg = Register{
		OverrideInterface: OverrideInterface{},
		OverrideFunction:  OverrideFunction{},
		Version:           map[string]string{"Test Version": "v0.0.1"},
		EnabledFeature:     []string{
			em_define.EnableCaptcha,
			em_define.EnableCircuitBreaker,
			em_define.EnableTranslate,
			em_define.EnableServiceDiscovery,
			em_define.EnableValidator,
		},
		RegisterService: func(s *grpc.Server) {
			helloworld.RegisterGreeterServer(s, &server{})
		},
		RegisterMiddleware: nil,
	}
	reg.Init()
	reg.Run()
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello again " + in.GetName()}, nil
}
