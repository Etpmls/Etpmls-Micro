package em_library

import (
	"log"
	"os"
	"path/filepath"
)

var (
	InitLog *log.Logger
)

type init_Library struct {
	logFile *os.File
}

func NewInit() *init_Library {
	return &init_Library{}
}

func (this *init_Library) Start()  {
	fileName := "../storage/log/init.log"
	_, err := os.Stat(filepath.Dir(fileName))
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	this.logFile, err  = os.OpenFile(fileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Fatalln(err)
	}
	InitLog = log.New(this.logFile,"",log.LstdFlags | log.Llongfile)
}

func (this *init_Library) Close()  {
	InitLog.Println("[INFO]", "Library initialization completed!")
	this.logFile.Close()
}