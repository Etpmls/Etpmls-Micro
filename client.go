package em

import (
	"context"
	"errors"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
	"google.golang.org/grpc"
	"time"
)

type client struct {
	Conn *grpc.ClientConn
	Context *context.Context
	Header map[string]string
}

func (this *client) NewClient() *client {
	return &client{}
}

func (this *client) ConnectService(service_name string) (error) {
	if ServiceDiscovery == nil {
		LogError.OutputSimplePath("ServiceDiscovery is not enabled!")
	}
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

func (this *client) Sync(run func() error, callback func(error) error) error {
	if CircuitBreaker == nil {
		LogError.OutputSimplePath("CircuitBreaker is not enabled!")
	}
	if this.Context == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		this.Context = &ctx
	}

	Micro.Request.Rpc_SetValueToHeader(this.Context, this.Header)

	defer this.Conn.Close()
	return CircuitBreaker.Sync("default", run, callback)
}

func (this *client) Async(run func() error, callback func(error) error) chan error {
	if CircuitBreaker == nil {
		LogError.OutputSimplePath("CircuitBreaker is not enabled!")
	}
	if this.Context == nil {
		var cancel context.CancelFunc
		*this.Context, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	}

	Micro.Request.Rpc_SetValueToHeader(this.Context, this.Header)

	defer this.Conn.Close()
	return CircuitBreaker.Async("default", run, callback)
}

func (this *client) ConnectServiceWithToken(service_name string, ctx *context.Context) ( error) {
	if ctx == nil {
		return errors.New("*context.Context is nil")
	}

	// 1.Connect Service
	err := this.ConnectService(service_name)
	if err != nil {
		return err
	}

	// 2. Set Header
	// Get token By Request
	this.Context = ctx
	token, err := Micro.Auth.GetTokenFromCtx(*ctx)
	if err != nil {
		return err
	}
	this.Header = map[string]string{"token": token}

	return nil
}

func (this *client) IsSuccess(r *em_protobuf.Response) error {
	if r.GetStatus() == SUCCESS_Status {
		return nil
	} else {
		LogWarn.Output(MessageWithLineNum_Advanced("Request failed!", 1, 20))
		return errors.New("Request failed!")
	}
}

// Deprecated: Use Sync_SimpleV2
func (this *client) Sync_Simple(run func() (*em_protobuf.Response, error), callback func(error) error) error {
	return this.Sync(func() error {
		r, err := run()
		if err != nil {
			LogWarn.OutputSimplePath("Run failed!", err)
			return err
		}
		return Micro.Client.IsSuccess(r)
	}, callback)
}

func (this *client) Sync_SimpleV2(run func() (*em_protobuf.Response, error), callback func(error) error) ([]byte, error) {
	var b []byte

	err := this.Sync(func() error {
		r, err := run()
		if err != nil {
			LogWarn.OutputSimplePath("Run failed!", err)
			return err
		}

		b = []byte(r.GetData())

		return Micro.Client.IsSuccess(r)
	}, callback)

	return b, err
}
