package em

import (
	"context"
	"errors"
	library "github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	utils "github.com/Etpmls/Etpmls-Micro/utils"
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"time"
)

type client struct {
}

func NewClient() *client {
	return &client{}
}

func (this *client) ConnectService(service_name string) (*grpc.ClientConn, error) {
	// Get service address
	// 获取服务地址
	host, err := library.ServiceDiscovery.GetServiceAddress_Random(service_name, nil)
	if err != nil {
		LogError.Output(utils.MessageWithLineNum(err.Error()))
		return nil, err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		LogError.Output(utils.MessageWithLineNum(err.Error()))
		return nil, err
	}

	return conn, nil
}

func (this *client) Go(name string, run func() error, fallBack func(error) error) chan error {
	return hystrix.Go(name, run, fallBack)
}

func (this *client) Do(name string, run func() error, fallBack func(error) error) error {
	return hystrix.Do(name, run, fallBack)
}

func (this *client) AuthCheck(authServiceName string, currentServiceName string, userId uint) (bool) {
	cl := NewClient()
	err := cl.Do("common", func() error {

		// Connect Service
		conn, err := cl.ConnectService(authServiceName)
		if err != nil {
			return err
		}
		defer conn.Close()
		c := em_protobuf.NewAuthClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.Check(ctx, &em_protobuf.AuthCheck{
			Service:       currentServiceName,
			UserId:        uint32(userId),
		})
		if err != nil {
			LogError.Output(utils.MessageWithLineNum(err.Error()))
			return err
		}

		if r.GetSuccess() == true {
			return nil
		} else {
			LogInfo.Output(utils.MessageWithLineNum("Check failed!"))
			return errors.New("Check failed!")
		}

	}, nil)
	if err != nil {
		return false
	}

	return true
}

func (this *client) SetValueToClientHeader(ctx *context.Context, m map[string]string) {
	g := library.NewGrpc()
	g.SetValueToMetadata(ctx, m)
	return
}