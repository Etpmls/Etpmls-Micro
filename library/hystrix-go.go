package em_library

import (
	"github.com/afex/hystrix-go/hystrix"
	"time"
)

func Init_HystrixGo(communicationTimeout time.Duration)  {
	hystrix.ConfigureCommand("default", hystrix.CommandConfig{
		Timeout:               int(communicationTimeout),
		ErrorPercentThreshold: 25,
	})
}

type hystrixGo struct {

}

func NewHystrixGo() *hystrixGo {
	return &hystrixGo{}
}

func (this *hystrixGo) Sync(name string, run func() error, fallBack func(error) error) error {
	return hystrix.Do(name, run, fallBack)
}

func (this *hystrixGo) Async(name string, run func() error, fallBack func(error) error) chan error {
	return hystrix.Go(name, run, fallBack)
}