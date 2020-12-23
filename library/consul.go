package em_library

import (
	"errors"
	Package_Consul "github.com/hashicorp/consul/api"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type ConsulConfig struct {
	Config *Package_Consul.Config
	Enable bool
	ConsulAddress string
	RpcId string
	RpcName string
	RpcPort string
	RpcTag []string
	HttpId string
	HttpName string
	HttpPort string
	HttpTag []string
	Prefix        string
	ServiceAddress       string
	CheckInterval string
	CheckUrl      string
}

var config *ConsulConfig

func Init_Consul(conf *ConsulConfig)  {
	//	If the service discovery is not turned on, skip consul initialization
	//	如果没有开启服务发现则跳过consul初始化
	if !conf.Enable {
		return
	}

	// Establish connection
	// 建立连接
	var err error
	Instance_Consul, err = Package_Consul.NewClient(conf.Config)
	if err != nil {
		initLog.Fatalln("[WARNING]", "Consul initialization failed.", " Error:", err)
	}

	// Registration Service
	// 注册服务
	config = conf
	var c = NewConsul()
	err = c.RegistrationService()
	if err != nil {
		initLog.Println("[WARNING]", "Registration Consul Service failed.", " Error:", err)
		go c.automaticRetry()
	} else {
		initLog.Println("[INFO]", "Registration Consul Service successfully.")
	}
}


var (
	Instance_Consul *Package_Consul.Client
)

type consul struct {}

func NewConsul() *consul {
	return &consul{}
}

// Registration Service
// 注册服务
func (this *consul) RegistrationService() error {
	// string -> int
	rpcPort, err := strconv.Atoi(config.RpcPort)
	if err != nil {
		return errors.New("rpcPort is not a int type! Error:" + err.Error())
	}
	httpPort, err := strconv.Atoi(config.HttpPort)
	if err != nil {
		return errors.New("httpPort is not a int type! Error:" + err.Error())
	}

	// Service Check
	Check := Package_Consul.AgentServiceCheck{
		Interval: config.CheckInterval,
		HTTP:     "http://" + config.ServiceAddress + ":" + config.HttpPort + config.CheckUrl,
	}

	// Configuration service
	conf := Package_Consul.AgentServiceRegistration{
		Address:   config.ServiceAddress,
		Check: &Check,
	}

	rpcConf := conf
	rpcConf.ID = config.RpcId
	rpcConf.Name = config.Prefix + config.RpcName
	rpcConf.Tags = config.RpcTag
	rpcConf.Port = rpcPort

	httpConf := conf
	httpConf.ID = config.HttpId
	httpConf.Name = config.Prefix + config.HttpName
	httpConf.Tags = config.HttpTag
	httpConf.Port = httpPort

	err = Instance_Consul.Agent().ServiceRegister(&rpcConf)
	if err != nil {
		return errors.New("Consul RPC Service registration failed! Error:" + err.Error())
	}

	err = Instance_Consul.Agent().ServiceRegister(&httpConf)
	if err != nil {
		return errors.New("Consul HTTP Service registration failed! Error:" + err.Error())
	}

	return nil
}

// Cancel Service
// 取消服务
func (this *consul) CancelService() error {
	err := Instance_Consul.Agent().ServiceDeregister(config.RpcId)
	if err != nil {
		initLog.Println("[ERROR]", "Cancel Consul RPC service failed!", " Error:", err)
		return err
	}
	err = Instance_Consul.Agent().ServiceDeregister(config.HttpId)
	if err != nil {
		initLog.Println("[ERROR]", "Cancel Consul HTTP service failed!", " Error:", err)
		return err
	}
	return nil
}

// When initial registration fails, automatically retry registration
// 当初始化注册失败时，自动重试注册
func (this *consul) automaticRetry() {
	for {
		time.Sleep(time.Second * 5)
		err := this.RegistrationService()
		if err == nil {
			initLog.Println("[INFO]", "Service registered successfully!")
			break
		}
	}
}

// Get service address
// 获取服务地址
func (this *consul) GetServiceAddr(service_name string, options map[string]interface{}) (string, error) {
	list, _, err := Instance_Consul.Health().Service(service_name, "",true, &Package_Consul.QueryOptions{})
	if err != nil {
		return "", err
	}

	// No record
	if len(list) == 0 {
		Instance_Logrus.Warning(service_name + " service not found.", " Error:", err)
		return "", errors.New(service_name + " service not found.")
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(list))

	addr := net.JoinHostPort(list[index].Service.Address, strconv.Itoa(list[index].Service.Port))
	return addr, nil
}