package em

import (
	"flag"
	"fmt"
)

const (
	Version_Framework = "1.2.2"
)

var (
	flag_Version = flag.Bool("version", false, "View current app version")
)

func init_version(reg *Register) bool {
	flag.Parse()
	if *flag_Version {
		VersionRegistration("Etpmls-Micro Version", Version_Framework)
		for k, v := range reg.Version_Service {
			VersionRegistration(k, v)
		}
		return true
	}
	return false
}

var versionMap = make(map[string]string)

// Version registration
// 版本注册
// Through this function, you can register the version number of your developed application in the global map
// 通过该函数，可以把你开发的应用版本号注册到全局map中
func VersionRegistration(key, value string)  {
	versionMap[key] = value
	return
}

// Print version information
// 打印版本信息
func VersionPrint()  {
	for k, v := range versionMap {
		fmt.Println(k, " : ", v)
	}
	return
}