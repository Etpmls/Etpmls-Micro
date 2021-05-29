package em

import (
	"errors"
	"fmt"
	"github.com/Etpmls/Etpmls-Micro/v3/define"
	"runtime"
	"strconv"
	"strings"
)

type Level uint32
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// Parse log level
// 解析log等级
func ParseLogLevel(str string) (Level, error) {
	switch strings.ToLower(str) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("Not a valid log Level: %q", str)
}

const (
	LOG_MODE_ONLY = "1"
	CONSOLE_MODE_ONLY = "2"
	LOG_CONSOLE_MODE = "3"
)

var (
	LogPanic = log{Level: PanicLevel}
	LogFatal = log{Level: FatalLevel}
	LogError = log{Level: ErrorLevel}
	LogWarn = log{Level: WarnLevel}
	LogInfo = log{Level: InfoLevel}
	LogDebug = log{Level: DebugLevel}
	LogTrace = log{Level: TraceLevel}
)

type log struct {
	Level Level
}

// No matter whether it is in Debug mode, it will output an message
// 无论是否为Debug模式，都输出信息
func (o log) New(info ...interface{}) {
	m, err := Kv.List(em_define.KvLogOutputMethod)
	if err != nil {
		Log.Error(err)
		return
	}
	lvl, lerr := Kv.ReadKey(em_define.KvLogLevel)
	if lerr != nil {
		lvl = em_define.DefaultLogLevel
	}
	l, err := ParseLogLevel(lvl)
	if err != nil {
		Log.Panic(MessageWithLineNum(err.Error()))
		return
	}

	switch o.Level {
	case PanicLevel:
		switch m[em_define.KvLogOutputMethodPanic] {
		case LOG_MODE_ONLY:
			Log.Panic(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Panic(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Panic(info...)
		}

	case FatalLevel:
		switch m[em_define.KvLogOutputMethodFatal] {
		case LOG_MODE_ONLY:
			Log.Fatal(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Fatal(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Fatal(info...)
		}

	case ErrorLevel:
		switch m[em_define.KvLogOutputMethodError] {
		case LOG_MODE_ONLY:
			Log.Error(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Error(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Error(info...)
		}

	case WarnLevel:
		switch m[em_define.KvLogOutputMethodWarning] {
		case LOG_MODE_ONLY:
			Log.Warning(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Warning(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Warning(info...)
		}

	case InfoLevel:
		switch m[em_define.KvLogOutputMethodInfo] {
		case LOG_MODE_ONLY:
			Log.Info(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Info(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Info(info...)
		}

	case DebugLevel:
		switch m[em_define.KvLogOutputMethodDebug] {
		case LOG_MODE_ONLY:
			Log.Debug(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Debug(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Debug(info...)
		}

	case TraceLevel:
		switch m[em_define.KvLogOutputMethodTrace] {
		case LOG_MODE_ONLY:
			Log.Trace(info...)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(info...)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Trace(info...)
		default:
			if l >= o.Level {
				fmt.Println(info...)
			}
			Log.Trace(info...)
		}

	}
}

// No matter whether it is in Debug mode, it will output an message, and return Error
// 无论是否为Debug模式，都输出信息，并且返回错误
func (this log) Error(info ...interface{}) error {
	this.New(info...)
	return errors.New(fmt.Sprintf("%v", info...))
}

// New information with file line number
// 输出带文件行数的信息
func (o log) Path(info ...interface{}) {
	var p = []interface{}{MessageWithLineNumAdvanced("", 1, 1)}
	p = append(p, info...)
	o.New(p...)
	return
}

func (this log) PathWithError(info ...interface{}) error {
	this.Path(info...)
	return errors.New(fmt.Sprintf("%v", info...))
}

// New information with the number of file lines and include the caller path
// 输出带文件行数的信息，并且包含调用者路径
func (o log) FullPath(info ...interface{}) {
	var p = []interface{}{MessageWithLineNumAdvanced("", 1, 20)}
	p = append(p, info...)
	o.New(p...)
	return
}

func (this log) FullPathWithError(info ...interface{}) error {
	this.FullPath(info...)
	return errors.New(fmt.Sprintf("%v", info...))
}

// If it is currently in Debug mode, it will output an return message, if it is in production mode, it will output a custom message
// 若当前为Debug模式，则输出返回信息，若为生产模式，则输出自定义信息
func (o log) DebugOrProd(err error, msg interface{}) {
	mp, err := Kv.List(em_define.KvLogOutputMethod)
	if err != nil {
		Log.Error(err)
		return
	}
	lvl, lerr := Kv.ReadKey(em_define.KvLogLevel)
	if lerr != nil {
		lvl = em_define.DefaultLogLevel
	}
	l, err := ParseLogLevel(lvl)
	if err != nil {
		Log.Panic(MessageWithLineNum(err.Error()))
		return
	}

	var m interface{}
	if IsDebug() {
		m = err
	} else {
		m = msg
	}

	switch o.Level {
	case PanicLevel:
		switch mp[em_define.KvLogOutputMethodPanic] {
		case LOG_MODE_ONLY:
			Log.Panic(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
			fmt.Println(m)
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Panic(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Panic(m)
		}

	case FatalLevel:
		switch mp[em_define.KvLogOutputMethodFatal] {
		case LOG_MODE_ONLY:
			Log.Fatal(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Fatal(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Fatal(m)
		}

	case ErrorLevel:
		switch mp[em_define.KvLogOutputMethodError] {
		case LOG_MODE_ONLY:
			Log.Error(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Error(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Error(m)
		}

	case WarnLevel:
		switch mp[em_define.KvLogOutputMethodWarning] {
		case LOG_MODE_ONLY:
			Log.Warning(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Warning(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Warning(m)
		}

	case InfoLevel:
		switch mp[em_define.KvLogOutputMethodInfo] {
		case LOG_MODE_ONLY:
			Log.Info(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Info(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Info(m)
		}

	case DebugLevel:
		switch mp[em_define.KvLogOutputMethodDebug] {
		case LOG_MODE_ONLY:
			Log.Debug(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Debug(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Debug(m)
		}

	case TraceLevel:
		switch mp[em_define.KvLogOutputMethodTrace] {
		case LOG_MODE_ONLY:
			Log.Trace(m)
		case CONSOLE_MODE_ONLY:
			if l >= o.Level {
				fmt.Println(m)
			}
		case LOG_CONSOLE_MODE:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Trace(m)
		default:
			if l >= o.Level {
				fmt.Println(m)
			}
			Log.Trace(m)
		}

	}
}

// Automatically output Debug, if it is a debug environment, it will output custom information + Error, if it is not a Debug environment, it will output custom information
// 自动输出Debug，如果是debug环境，则输出自定义信息+Error，如果不是Debug环境，输出自定义信息
func (o log) DetailedIfDebug (msg interface{}, err error) {
	v, ok := msg.(string);
	if !ok {
		o.DebugOrProd(err, msg)
		return
	}

	o.DebugOrProd(GenerateErrorWithMessage(v + "Error: ", err), msg)
	return
}

// Message(or Error) with line number
// 消息(或错误)带行号
func MessageWithLineNum(msg string) string {
	var list []string
	for i := 1; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			list = append(list, file+":"+strconv.Itoa(line))
		} else {
			break
		}
	}
	return strings.Join(list, " => ") + " => Message: " + msg
}

/*func messageWithLineNum_Local(msg string) string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(file))
	sourceDir := strings.ReplaceAll(dir, "\\", "/")

	var list []string
	for i := 1; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && strings.HasPrefix(file, sourceDir) {
			list = append(list, file+":"+strconv.Itoa(line))
		} else {
			break
		}
	}
	return strings.Join(list, " => ") + " => Message: " + msg
}*/

// Message(or Error) with line number - Only one record
// 消息(或错误)带行号 - 仅一条记录
func MessageWithLineNumOneRecord(msg string) string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return file + ":" + strconv.Itoa(line) + " => Message: " + msg
	}
	return msg
}

// Message(or Error) with line number,Specify call level
// 消息(或错误)带行号，指定调用层级
func MessageWithLineNumAdvanced(msg string, level int, num int) string {
	var list []string
	for i := level + 1; i < level+1+num; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			list = append(list, file+":"+strconv.Itoa(line))
		} else {
			break
		}
	}
	return strings.Join(list, " => ") + " => Message: " + msg
}

