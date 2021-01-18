package em_library

import (
	"fmt"
	"github.com/Etpmls/Etpmls-Micro/define"
	"github.com/hashicorp/consul/api"
	Package_Logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

var Instance_Logrus *Package_Logrus.Logger

func Init_Logrus(logLevel string) {
	Instance_Logrus = Package_Logrus.New()
	// Instance_Logrus as JSON instead of the default ASCII formatter.
	Instance_Logrus.Formatter = new(Package_Logrus.JSONFormatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Instance_Logrus.Out = &lumberjack.Logger{
		Filename:   "../storage/log/app.log",
		MaxSize:    500, // megabytes
		MaxAge:     30, //days
		Compress:   true, // disabled by default
	}

	// Only log the warning severity or above.
	level, err := Package_Logrus.ParseLevel(logLevel)
	if err != nil {
		level = Package_Logrus.WarnLevel
		InitLog.Fatalln("[FATAL]", "Set Instance_Logrus Level Failed!", " Error:", err)
	} else {
		InitLog.Println("[INFO]", "logrus initialized successfully.")
	}
	Instance_Logrus.Level = level
}

type logrus struct {

}

func NewLogrus() *logrus {
	return &logrus{}
}

func (this *logrus) Panic(args ...interface{}) {
	Instance_Logrus.Panic("[PANIC]", args)
	go this.logToKv("[PANIC]", Config.Service.RpcId + ": ", args)
	return
}

func (this *logrus) Fatal(args ...interface{}) {
	Instance_Logrus.Fatal("[FATAL]", args)
	if Instance_Logrus.Level >= Package_Logrus.FatalLevel {
		go this.logToKv("[FATAL]", Config.Service.RpcId + ": ", args)
	}
	return
}

func (this *logrus) Error(args ...interface{}) {
	Instance_Logrus.Error("[ERROR]", args)
	if Instance_Logrus.Level >= Package_Logrus.ErrorLevel {
		go this.logToKv("[ERROR]", Config.Service.RpcId + ": ", args)
	}
	return
}

func (this *logrus) Warning(args ...interface{}) {
	Instance_Logrus.Warning("[WARNING]", args)
	if Instance_Logrus.Level >= Package_Logrus.WarnLevel {
		go this.logToKv("[WARNING]", Config.Service.RpcId + ": ", args)
	}

	return
}

func (this *logrus) Info(args ...interface{}) {
	Instance_Logrus.Info("[INFO]", args)
	if Instance_Logrus.Level >= Package_Logrus.InfoLevel {
		go this.logToKv("[INFO]", Config.Service.RpcId + ": ", args)
	}

	return
}

func (this *logrus) Debug(args ...interface{}) {
	Instance_Logrus.Debug("[DEBUG]", args)
	if Instance_Logrus.Level >= Package_Logrus.DebugLevel {
		go this.logToKv("[DEBUG]", Config.Service.RpcId + ": ", args)
	}

	return
}

func (this *logrus) Trace(args ...interface{}) {
	Instance_Logrus.Trace("[TRACE]", args)
	if Instance_Logrus.Level >= Package_Logrus.TraceLevel {
		go this.logToKv("[TRACE]", Config.Service.RpcId + ": ", args)
	}

	return
}

func (this *logrus) logToKv(args ...interface{}) {
	k := define.KvLogLog + time.Now().Format("2006-01")
	pair, _, err := kv.Get(k, nil)
	if err != nil || pair == nil {
		p := &api.KVPair{Key: k, Value: []byte(fmt.Sprint(args))}
		_, err := kv.Put(p, nil)
		if err != nil {
			Instance_Logrus.Error(err)
			return
		}
		return
	}

	p := &api.KVPair{Key: k, Value: []byte(string(pair.Value) + "\n" + fmt.Sprint(args))}
	_, err = kv.Put(p, nil)
	if err != nil {
		Instance_Logrus.Error(err)
	}
	return
}

