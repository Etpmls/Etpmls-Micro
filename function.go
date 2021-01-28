package em

import (
	"errors"
	"github.com/Etpmls/Etpmls-Micro/define"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
)

func PanicIfMapValueEmpty(key string, m map[string]string) string {
	if len(m[key]) == 0 {
		LogError.OutputFullPath("[ERROR]", key, " is not configured!")
		panic(("[ERROR]"+ key+ " is not configured!"))
	}
	return m[key]
}

func MustGetServiceKvKey(key string) string {
	nameKey := define.MakeServiceConfField(em_library.Config.Service.RpcName, key)
	idKey := define.MakeServiceConfField(em_library.Config.Service.RpcId, key)

	// Get key
	idValue, err := Kv.ReadKey(idKey)
	if err == nil && len(idValue) != 0 {
		return idValue
	}

	nameValue, err := Kv.ReadKey(nameKey)
	if err == nil && len(nameValue) != 0 {
		return nameValue
	}

	if !Reg.initFinished {
		em_library.InitLog.Panicln("[PANIC]", "Please configure the ", idKey, " or ", nameKey)
	} else {
		LogError.Output("Please configure the ", idKey, " or ", nameKey)
	}
	panic("Please configure the "+ idKey+ " or "+ nameKey)
}

// Id Key
func MustGetServiceIdKvKey(key string) string {
	idKey := define.MakeServiceConfField(em_library.Config.Service.RpcId, key)

	// Get key
	idValue, err := Kv.ReadKey(idKey)
	if err == nil && len(idValue) != 0 {
		return idValue
	}

	if !Reg.initFinished {
		em_library.InitLog.Panicln("[PANIC]", "Please configure the ", idKey)
	} else {
		LogError.Output("Please configure the ", idKey)
	}
	panic("Please configure the "+ idKey)
}

func MustListServiceIdKvKey(key string) map[string]string {
	idKey := define.MakeServiceConfField(em_library.Config.Service.RpcId, key)

	// Get key
	idValue, err := Kv.List(idKey)
	if err == nil && len(idValue) != 0 {
		return idValue
	}

	if !Reg.initFinished {
		em_library.InitLog.Panicln("[PANIC]", idKey, " has no value")
	} else {
		LogError.Output(idKey, " has no value")
	}
	panic(idKey+ " has no value")
}

// Name Key
func MustGetServiceNameKvKey(key string) string {
	nameKey := define.MakeServiceConfField(em_library.Config.Service.RpcName, key)

	// Get key
	nameValue, err := Kv.ReadKey(nameKey)
	if err == nil && len(nameValue) != 0 {
		return nameValue
	}

	if !Reg.initFinished {
		em_library.InitLog.Panicln("[PANIC]", "Please configure the ", nameKey)
	} else {
		LogError.Output("Please configure the ", nameKey)
	}
	panic("Please configure the "+ nameKey)
}

func MustListServiceNameKvKey(key string) map[string]string {
	nameKey := define.MakeServiceConfField(em_library.Config.Service.RpcName, key)

	// Get key
	nameValue, err := Kv.List(nameKey)
	if err == nil && len(nameValue) != 0 {
		return nameValue
	}

	if !Reg.initFinished {
		em_library.InitLog.Panicln("[PANIC]", nameKey, " has no value")
	} else {
		LogError.Output(nameKey, " has no value")
	}
	panic(nameKey+ " has no value")
}

func MustGetKvKey(key string) string {
	s, err := Kv.ReadKey(key)
	if err != nil || len(s) == 0 {
		if !Reg.initFinished {
			em_library.InitLog.Panicln("[PANIC]", "Please configure the ", key)
		} else {
			LogError.Output("Please configure the ", key)
		}
		panic("Please configure the "+ key)
	}
	return s
}

func GetServiceKvKey(key string) (string, error) {
	nameKey := define.MakeServiceConfField(em_library.Config.Service.RpcName, key)
	idKey := define.MakeServiceConfField(em_library.Config.Service.RpcId, key)

	// Get key
	idValue, err := Kv.ReadKey(idKey)
	if err == nil && len(idValue) != 0 {
		return idValue, nil
	}

	nameValue, err := Kv.ReadKey(nameKey)
	if err == nil && len(nameValue) != 0 {
		return nameValue, nil
	}

	return "", errors.New(idKey+ " or "+ nameKey+ " not found.")
}

func GetServiceIdKvKey(key string)  (string, error) {
	idKey := define.MakeServiceConfField(em_library.Config.Service.RpcId, key)

	// Get key
	idValue, err := Kv.ReadKey(idKey)
	if err == nil && len(idValue) != 0 {
		return idValue, nil
	}

	return "", errors.New(idKey+ " not found.")
}

func ListServiceIdKvKey(key string)  (map[string]string, error) {
	idKey := define.MakeServiceConfField(em_library.Config.Service.RpcId, key)

	// Get key
	idValue, err := Kv.List(idKey)
	if err == nil && len(idValue) != 0 {
		return idValue, nil
	}

	return nil, errors.New(idKey+ " not found.")
}

func GetServiceNameKvKey(key string)  (string, error) {
	nameKey := define.MakeServiceConfField(em_library.Config.Service.RpcName, key)

	// Get key
	nameValue, err := Kv.ReadKey(nameKey)
	if err == nil && len(nameValue) != 0 {
		return nameValue, nil
	}

	return "", errors.New(nameKey+ " not found.")
}

func ListServiceNameKvKey(key string)  (map[string]string, error) {
	nameKey := define.MakeServiceConfField(em_library.Config.Service.RpcName, key)

	// Get key
	nameValue, err := Kv.List(nameKey)
	if err == nil && len(nameValue) != 0 {
		return nameValue, nil
	}

	return nil, errors.New(nameKey+ " not found.")
}
