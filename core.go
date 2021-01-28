package em

import (
	"encoding/json"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
	"google.golang.org/grpc/codes"
	"os"
	"strings"
)

// Error Code
// 错误码
const (
	ERROR_Code = "400000"
	ERROR_Status = "error"
)


// Success Code
// 成功码
const (
	SUCCESS_Code = "200000"
	SUCCESS_Status = "success"
)

// Whether it is a DEBUG environment
// 是否为DEBUG环境
func IsDebug() bool {
	if strings.ToUpper(os.Getenv("DEBUG")) == "TRUE" {
		return true
	}
	return false
}


type Response em_protobuf.Response

// Used to quickly transform itself into json, and need to be used when returning a response
// 用于快速把自身转化json，返回响应时需要使用
func (this Response) String() (string) {
	b, err := json.Marshal(this)
	if err != nil {
		return ""
	}
	return string(b)
}

// Return error message in json format
// 返回json格式的错误信息
func ErrorRpc(rcpStatusCode codes.Code, code string, message string, data interface{}, err error) (*em_protobuf.Response, error) {
	return Reg.RpcReturnErrorFunc(rcpStatusCode, code, message, data, err)
}

func ErrorHttp(code string, message string, data interface{}, err error) ([]byte, error) {
	return Reg.HttpReturnErrorFunc(code, message, data, err)
}

// Return success information in json format
// 返回json格式的成功信息
func SuccessRpc(code string, message string, data interface{}) (*em_protobuf.Response, error) {
	return Reg.RpcReturnSuccessFunc(code, message, data)
}

func SuccessHttp(code string, message string, data interface{}) ([]byte, error) {
	return Reg.HttpReturnSuccessFunc(code, message, data)
}
