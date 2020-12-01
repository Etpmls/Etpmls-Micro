// https://github.com/go-yaml/yaml

package em_library

import (
	utils "github.com/Etpmls/Etpmls-Micro/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type Configuration struct {
	App struct{
		RpcPort string	`yaml:"rpc-port"`
		HttpPort string	`yaml:"http-port"`
		Captcha bool
		Register bool
		Key string
		Database bool
		Cache bool
		ServiceDiscovery	bool	`yaml:"service-discovery"`
		TokenExpirationTime time.Duration	`yaml:"token-expiration-time"`
		UseHttpCode bool	`yaml:"use-http-code"`
		TimeZone string		`yaml:"time-zone"`
		CommunicationTimeout time.Duration	`yaml:"communication-timeout"`
	}
	Database struct{
		Host string
		Port string
		Name string
		User string
		Password string
		Prefix string
	}
	ServiceDiscovery struct{
		Address string
		Service struct{
			Id string
			Name string
			Address string
			Tag []string
			CheckInterval string	`yaml:"check-interval"`
			CheckUrl string	`yaml:"check-url"`
		}
	}	`yaml:"service-discovery"`
	Cache struct{
		Address string
		Password string
		DB int
	}
	Captcha struct{
		Host   string
		Secret string
		Timeout time.Duration
	}
	Log struct {
		Level string
		Panic	int
		Fatal	int
		Error	int
		Warning	int
		Info	int
		Debug	int
		Trace	int
	}
	Field struct{
		Pagination struct {
			Number string
			Size string
			Count string
		}
	}
}

var Config = Configuration{}

func Init_Yaml() {
	var yamlPath string

	if os.Getenv("DEBUG") == "TRUE" {
		yamlPath = "storage/config/app_debug.yaml"
	} else{
		yamlPath = "storage/config/app.yaml"
	}

	b, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		Instance_Logrus.Fatal("Failed to read the Configuration file! Error:", err)
		return
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		Instance_Logrus.Fatal("Failed to unmarshal the Configuration file! Error:", err)
		return
	}

	if len(Config.App.Key) < 50 {
		Config.App.Key = utils.GenerateRandomString(50)

		out, err := yaml.Marshal(Config)
		if err != nil {
			Instance_Logrus.Fatal("Failed to parse the Configuration file into yaml format!", err)
			return
		}

		err = ioutil.WriteFile(yamlPath, out, os.ModeAppend)
		if err != nil {
			Instance_Logrus.Fatal("Failed to write yaml Configuration file!", err)
			return
		}
	}

	return
}

func Init_CustomYaml(path, debug_path string, structAddr interface{})  {
	// If it is empty, skip initialization
	// 如果为空，则跳过初始化
	if path == "" || debug_path == "" || structAddr == nil {
		return
	}

	var yamlPath string

	if os.Getenv("DEBUG") == "TRUE" {
		yamlPath = debug_path
	} else{
		yamlPath = path
	}

	b, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		Instance_Logrus.Fatal("Failed to read the Configuration file! Error:", err)
		return
	}

	err = yaml.Unmarshal(b, structAddr)
	if err != nil {
		Instance_Logrus.Fatal("Failed to unmarshal the Configuration file! Error:", err)
		return
	}
}









