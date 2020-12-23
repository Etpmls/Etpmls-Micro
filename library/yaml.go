// https://github.com/go-yaml/yaml

package em_library

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type Configuration struct {
	App struct{
		RpcPort string	`yaml:"rpc-port"`
		HttpPort string	`yaml:"http-port"`
		Key string
		EnableDatabase bool	`yaml:"enable-database"`
		EnableCache bool	`yaml:"enable-cache"`
		EnableServiceDiscovery	bool	`yaml:"enable-service-discovery"`
		TokenExpirationTime time.Duration	`yaml:"token-expiration-time"`
		UseHttpCode bool	`yaml:"use-http-code"`
		CommunicationTimeout time.Duration	`yaml:"communication-timeout"`
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
	Database struct{
		Host string
		Port string
		Name string
		User string
		Password string
		Prefix string
		TimeZone string		`yaml:"time-zone"`
	}
	ServiceDiscovery struct{
		Address string
		Service struct{
			Rpc struct{
				Id string
				Name string
				Tag []string
			}
			Http struct{
				Id string
				Name string
				Tag []string
			}
			Prefix        string
			Address       string
			CheckInterval string	`yaml:"check-interval"`
			CheckUrl      string	`yaml:"check-url"`
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
		initLog.Println("[ERROR]", "Failed to read the Configuration file! Error:", err)
		return
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		initLog.Fatal("Failed to unmarshal the Configuration file! Error:", err)
		return
	}

	if len(Config.App.Key) < 50 {
		var y yamlV2
		Config.App.Key = y.utils_GenerateRandomString(50)

		out, err := yaml.Marshal(Config)
		if err != nil {
			initLog.Fatal("Failed to parse the Configuration file into yaml format!", err)
			return
		}

		err = ioutil.WriteFile(yamlPath, out, os.ModeAppend)
		if err != nil {
			initLog.Fatal("Failed to write yaml Configuration file!", err)
			return
		}
	}

	initLog.Println("[INFO]", "Successfully loaded configuration file!")
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
		log.Fatal("Failed to read the Configuration file! Error:", err)
		return
	}

	err = yaml.Unmarshal(b, structAddr)
	if err != nil {
		log.Fatal("Failed to unmarshal the Configuration file! Error:", err)
		return
	}
}

type yamlV2 struct {

}

// Generate random strings
// 生成随机字符串
func (this *yamlV2) utils_GenerateRandomString(l int) string {
	var code = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/~!@#$%^&*()_="

	data := make([]byte, l)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < l; i++ {
		idx := rand.Intn(len(code))
		data[i] = code[idx]
	}
	return string(data)
}







