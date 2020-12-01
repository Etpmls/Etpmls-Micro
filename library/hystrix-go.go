package em_library

import (
	"github.com/afex/hystrix-go/hystrix"
)

func Init_HystrixGo()  {
	hystrix.ConfigureCommand("common", hystrix.CommandConfig{
		Timeout:               int(Config.App.CommunicationTimeout),
		ErrorPercentThreshold: 25,
	})
}