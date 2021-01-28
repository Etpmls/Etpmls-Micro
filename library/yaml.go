// https://github.com/go-yaml/yaml

package em_library

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Kv struct{
		Address []string
		Token string
	}
	Service struct{
		RpcId	string `yaml:"rpc-id"`
		RpcName	string `yaml:"rpc-name"`
	}
}

var Config = Configuration{}

func Init_Yaml() {
	var yamlPath string

	if os.Getenv("DEBUG") == "TRUE" {
		yamlPath = "./storage/config/app_debug.yaml"
	} else{
		yamlPath = "./storage/config/app.yaml"
	}

	b, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		InitLog.Println("[ERROR]", "Failed to read the Configuration file! Error:", err)
		return
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		InitLog.Fatal("Failed to unmarshal the Configuration file! Error:", err)
		return
	}

	// Validate
	switch  {
	case len(Config.Kv.Address) == 0:
		InitLog.Fatal("[FATAL]", "You need to configure Config.Kv.Address!")
		return
	case Config.Service.RpcId == "":
		InitLog.Fatal("[FATAL]", "You need to configure Config.Service.RpcId!")
		return
	case Config.Service.RpcName == "":
		InitLog.Fatal("[FATAL]", "You need to configure Config.Service.RpcName!")
		return
	default:

	}

 	if len(Config.Kv.Token) == 0 {
		InitLog.Println("[WARNING]", "Config.Kv.Token is not configured.")
	}

	InitLog.Println("[INFO]", "Successfully loaded configuration file!")
	return
}







