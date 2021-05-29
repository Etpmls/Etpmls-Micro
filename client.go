package em

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type client struct {
	Conn *grpc.ClientConn
	Context *context.Context
	Header map[string]string
}

func NewClient() *client {
	return &client{}
}

func (this *client) ConnectService(service_name string) (error) {
	if ServiceDiscovery == nil {
		LogError.Path("ServiceDiscovery is not enabled!")
		return errors.New("ServiceDiscovery is not enabled!")
	}
	addr, err := ServiceDiscovery.GetServiceAddr(service_name, nil)
	if err != nil {
		LogError.Path(err.Error())
		return err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		LogError.Path(err.Error())
		return err
	}
	this.Conn = conn
	return nil
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

func (this *client) Sync(run func() error, callback func(error) error) error {
	if CircuitBreaker == nil {
		LogError.Path("CircuitBreaker is not enabled!")
	}
	if this.Context == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		this.Context = &ctx
	}

	Micro.Request.SetValueToHeader(this.Context, this.Header)

	defer this.Conn.Close()
	return CircuitBreaker.Sync("default", run, callback)
}

func (this *client) Async(run func() error, callback func(error) error) chan error {
	if CircuitBreaker == nil {
		LogError.Path("CircuitBreaker is not enabled!")
	}
	if this.Context == nil {
		var cancel context.CancelFunc
		*this.Context, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	}

	Micro.Request.SetValueToHeader(this.Context, this.Header)

	defer this.Conn.Close()
	return CircuitBreaker.Async("default", run, callback)
}
