package language

import (
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"path/filepath"
	"runtime"
)

func LoadLanguage()  {
	_, fileStr, _, _ := runtime.Caller(0)

	list, err := filepath.Glob(filepath.Dir(fileStr) + "/*.toml")
	if err != nil || len(list) < 1 {
		em_library.Instance_Logrus.Fatal("Failed to load language pack!")
		return
	}

	for _, v:=range list {
		em_library.Instance_GoI18n.MustLoadMessageFile(v)
	}
	return
}
