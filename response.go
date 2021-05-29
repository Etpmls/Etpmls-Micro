package em

import (
	"context"
	"encoding/json"
	"github.com/Etpmls/Etpmls-Micro/v3/proto/empb"
	"google.golang.org/grpc/codes"
)

// Return error message in json format
// 返回json格式的错误信息
func Error(code codes.Code, msg string, data interface{}, err error) (*empb.Response, error) {
	b, err := json.Marshal(data)
	LogError.FullPath(err)

	return &empb.Response{
		Code:    uint32(code),
		Status:  false,
		Message: msg,
		Data:    string(b),
	}, err
}

// Return error message in json format and translate
// 返回json格式的错误信息并翻译
func ErrorTranslate(ctx context.Context, code codes.Code, msg string, data interface{}, err error) (*empb.Response, error) {
	return Error(code, Translate.TranslateFromRequest(ctx, msg), data, err)
}

// Return success information in json format
// 返回json格式的成功信息
func Success(code codes.Code, msg string, data interface{}) (*empb.Response, error) {
	b, err := json.Marshal(data)
	LogError.FullPath(err)

	return &empb.Response{
		Code:    uint32(code),
		Status:  true,
		Message: msg,
		Data:    string(b),
	}, err
}

// Return success information in json format and translate
// 返回json格式的成功信息并翻译
func SuccessTranslate(ctx context.Context, code codes.Code, msg string, data interface{}) (*empb.Response, error) {
	return Success(code, Translate.TranslateFromRequest(ctx, msg), data)
}