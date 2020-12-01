package em_utils

import (
	"errors"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

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
func messageWithLineNum_Local(msg string) string {
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
}


// Generate errors with both custom messages and error messages
// 生成同时带有自定义信息和错误信息的错误
func GenerateErrorWithMessage(msg string, err error) error {
	return errors.New(msg + err.Error())
}

// Debug errors and custom errors are used as parameters at the same time, and appropriate errors are output according to environment variables.
// 把Debug错误和自定义错误同时作为参数，根据环境变量输出适合的错误。
func GetErrorByIfDebug(debug bool, err error, msg string) error {
	if debug {
		return err
	}
	return errors.New(msg)
}