package em_library

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type ServiceConfig struct {
	Config   *api.Config
	RpcId    string
	RpcName  string
	RpcPort  string
	RpcTag   []string
	HttpId   string
	HttpName string
	HttpPort string
	HttpTag  []string
	Address  string
	/*
		Health check
	*/
	CheckInterval string
	CheckUrl      string
}

var config *ServiceConfig

func Init_Consul(conf *ServiceConfig)  {
	if conf.Config == nil {
		InitLog.Println("[ERROR]", "Consul is not configured!")
		return
	}

	// Establish connection
	// 建立连接
	var err error
	Instance_Consul, err = api.NewClient(conf.Config)
	if err != nil {
		InitLog.Fatalln("[WARNING]", "Consul initialization failed.", " Error:", err)
	}

	// Registration Service
	// 注册服务
	config = conf
	var c = NewConsul()
	err = c.RegistrationService()
	if err != nil {
		InitLog.Println("[ERROR]", "Registration Consul Service failed.", " Error:", err)
		go c.automaticRetry()
	} else {
		InitLog.Println("[INFO]", "Registration Consul Service successfully.")
	}
}


func InitConsulKv(c *api.Config)  {
	// Establish connection
	// 建立连接
	cl, err := api.NewClient(c)
	if err != nil {
		InitLog.Fatalln("[WARNING]", "Consul initialization failed.", " Error:", err)
	}

	// Registration KV
	kv = cl.KV()
	return
}

var (
	Instance_Consul *api.Client
	kv *api.KV
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
	Check := api.AgentServiceCheck{
		Interval: config.CheckInterval,
		HTTP:     "http://" + config.Address + ":" + config.HttpPort + config.CheckUrl,
	}

	// Configuration service
	conf := api.AgentServiceRegistration{
		Address:   config.Address,
		Check: &Check,
	}

	rpcConf := conf
	rpcConf.ID = config.RpcId
	rpcConf.Name = config.RpcName
	rpcConf.Tags = config.RpcTag
	rpcConf.Port = rpcPort

	httpConf := conf
	httpConf.ID = config.HttpId
	httpConf.Name = config.HttpName
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
		InitLog.Println("[ERROR]", "Cancel Consul RPC service failed!", " Error:", err)
		return err
	}
	err = Instance_Consul.Agent().ServiceDeregister(config.HttpId)
	if err != nil {
		InitLog.Println("[ERROR]", "Cancel Consul HTTP service failed!", " Error:", err)
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
			InitLog.Println("[INFO]", "Service registered successfully!")
			break
		}
	}
}

// Get service address
// 获取服务地址
func (this *consul) GetServiceAddr(service_name string, options map[string]interface{}) (string, error) {
	list, _, err := Instance_Consul.Health().Service(service_name, "",true, &api.QueryOptions{})
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

// Read Key
// 读取Key
func (this *consul) ReadKey(key string) (string, error) {
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", err
	}
	if pair == nil {
		return "", errors.New("Could not find the key of "+ key)
	}
	return string(pair.Value), nil
}

// Update if there is a key, otherwise create
// 如果存在key则更新，否则创建
func (this *consul) CrateOrUpdateKey(key, value string) error {
	p := &api.KVPair{Key: key, Value:[]byte(value)}
	_, err := kv.Put(p, nil)
	if err != nil {
		return err
	}
	return nil
}

// Delete key
// 删除Key
func (this *consul) DeleteKey(key string) error {
	_, err := kv.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (this *consul) List(prefix string) (map[string]string, error) {
	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, v := range pairs {
		m[v.Key] = string(v.Value)
	}
	return m, nil
}