package em

import (
	"github.com/Etpmls/Etpmls-Micro/v3/define"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"strings"
)

func defaultHandleExit() {
	e, err := Kv.ReadKey(em_define.KvServiceDiscoveryEnable)
	if err != nil {
		LogInfo.Path(err)
		return
	}

	if strings.ToLower(e) == "true" {
		if ServiceDiscovery != nil {
			err := ServiceDiscovery.CancelService()
			if err != nil {
				LogError.New("Cancel service failed! " + MessageWithLineNum(err.Error()))
			}
		}
	}

	return
}

func defaultGrpcMiddlewareFunc() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// Panic recover
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)
	return s
}



