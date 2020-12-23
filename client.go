package em

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

type client struct {
}

func (this *client) NewClient() *cli {
	return &cli{}
}

type cli struct {
	Conn *grpc.ClientConn
	Context *context.Context
	Header map[string]string
}

func (this *cli) ConnectService(service_name string) (error) {
	addr, err := ServiceDiscovery.GetServiceAddr(service_name, nil)
	if err != nil {
		LogError.Output(MessageWithLineNum_OneRecord(err.Error()))
		return err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		LogError.Output(MessageWithLineNum(err.Error()))
		return err
	}
	this.Conn = conn
	return nil
}

func (this *cli) Sync(run func() error, callback func(error) error) error {
	if this.Context == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		this.Context = &ctx
	}

	Micro.Request.Rpc_SetValueToHeader(this.Context, this.Header)

	defer this.Conn.Close()
	return CircuitBreaker.Sync("default", run, callback)
}

func (this *cli) Async(run func() error, callback func(error) error) chan error {
	if this.Context == nil {
		var cancel context.CancelFunc
		*this.Context, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	}

	Micro.Request.Rpc_SetValueToHeader(this.Context, this.Header)

	defer this.Conn.Close()
	return CircuitBreaker.Async("default", run, callback)
}




/*
func (this *cli) AuthCheck(authServiceName string, currentServiceName string, userId uint) (bool) {
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
*/