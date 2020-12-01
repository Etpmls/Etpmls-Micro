package em_library

import (
	Package_Consul "github.com/hashicorp/consul/api"
	"math/rand"
	"strconv"
	"time"
)

var (
	Instance_Consul *Package_Consul.Client
)

func Init_Consul()  {
	//	If the service discovery is not turned on, skip consul initialization
	//	如果没有开启服务发现则跳过consul初始化
	if !Config.App.ServiceDiscovery {
		return
	}

	config := Package_Consul.DefaultConfig()
	config.Address = Config.ServiceDiscovery.Address

	// Establish connection
	// 建立连接
	var err error
	Instance_Consul, err = Package_Consul.NewClient(config)
	if err != nil {
		Instance_Logrus.Warning("Consul initialization failed.")
	} else {
		Instance_Logrus.Info("Consul initialized successfully.")
	}

	// Registration Service
	// 注册服务
	c := NewConsul()
	_ = c.RegistrationService()
}

type consul struct {

}

func NewConsul() *consul {
	return &consul{}
}

// Registration Service
// 注册服务
func (this *consul) RegistrationService() error {
	// string -> int
	port, err := strconv.Atoi(Config.App.RpcPort)
	if err != nil {
		Instance_Logrus.Error("port is not a int type! Error:", err.Error())
		return err
	}

	r := Package_Consul.AgentServiceRegistration{
		ID:      Config.ServiceDiscovery.Service.Id,
		Name:    Config.ServiceDiscovery.Service.Name,
		Tags:    Config.ServiceDiscovery.Service.Tag,
		Port:    port,
		Address: Config.ServiceDiscovery.Service.Address,
	}

	c := Package_Consul.AgentServiceCheck{
		Interval: Config.ServiceDiscovery.Service.CheckInterval,
		HTTP:     "http://" + Config.ServiceDiscovery.Service.Address + ":" + Config.App.HttpPort + Config.ServiceDiscovery.Service.CheckUrl,
	}

	r.Check = &c
	err = Instance_Consul.Agent().ServiceRegister(&r)
	if err != nil {
		Instance_Logrus.Error("Consul Service registration failed! Error:", err.Error())
		return err
	}

	return nil
}

// Cancel Service
// 取消服务
func (this *consul) CancelService() error {
	err := Instance_Consul.Agent().ServiceDeregister(Config.ServiceDiscovery.Service.Id)
	if err != nil {
		Instance_Logrus.Error("Cancel Consul service failed! Error:", err.Error())
		return err
	}
	return nil
}

// Get a random service address
// 随机获取一个服务的地址
func (this *consul) GetServiceAddress_Random(service_name string, options map[string]interface{}) (string, error) {
	conf := Package_Consul.DefaultConfig()
	conf.Address = Config.ServiceDiscovery.Address

	// Get a new client
	client, err := Package_Consul.NewClient(conf)
	if err != nil {
		return "", err
	}
	list, _, err := client.Health().Service(service_name, "",true, &Package_Consul.QueryOptions{})
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		return "", err
	}

	// Get random number
	// 获取随机数
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(list))

	return list[index].Service.TaggedAddresses["lan_ipv4"].Address + ":" + strconv.Itoa(list[index].Service.TaggedAddresses["lan_ipv4"].Port), nil
}
