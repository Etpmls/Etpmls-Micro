// https://github.com/go-yaml/yaml

package em_library

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Kv struct{
		Address string
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
		yamlPath = "../storage/config/app_debug.yaml"
	} else{
		yamlPath = "../storage/config/app.yaml"
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

	InitLog.Println("[INFO]", "Successfully loaded configuration file!")
	return
}







