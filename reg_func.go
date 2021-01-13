package em

import (
	"encoding/json"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

func defaultHandleExit() {
	e, err := Kv.ReadKey(KvServiceDiscoveryEnable)
	if err != nil {
		LogInfo.OutputSimplePath(err)
		return
	}

	if strings.ToLower(e) != "true" {
		if ServiceDiscovery != nil {
			err := ServiceDiscovery.CancelService()
			if err != nil {
				LogError.Output("Cancel service failed! " + MessageWithLineNum(err.Error()))
			}
		}
	}

	return
}

func defaultRegisterEndpoint() *runtime.ServeMux {
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(NewEndpoint().CustomErrorHandlerFunc),
		// https://grpc-ecosystem.github.io/grpc-gateway/docs/customizingyourgateway/
		runtime.WithIncomingHeaderMatcher(NewEndpoint().SetCustomMatcher),
	)
	return mux
}

func defaultCorsOptions() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type", "Language", "Token"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		// Debug: true,
	}
}

func defaultHandleHttpSucessFunc(code string, message string, data interface{}) ([]byte, error) {
	// First judge whether it is a string, if it is a string, then judge whether it belongs to json, if it already belongs to json, no longer convert json again
	// 首先判断是不是字符串，如果是字符串，则判断属不属于json，如果已经属于json，则不再二次转化json
	if v, ok := data.(string); ok {
		ok_json := json.Valid([]byte(v))
		if ok_json {
			b, err := json.Marshal(
				Response{
					Code:    code,
					Status:  "success",
					Message: message,
					Data:    v,
				})
			return b, err
		}
	}

	b, err := json.Marshal(
		Response{
			Code:    code,
			Status:  "success",
			Message: message,
			Data:    MustConvertJson(data),
		})
	return b, err
}

func defaultHandleHttpErrorFunc(code string, message string, data interface{}, err error) ([]byte, error) {
	// Data interface => Json string
	var tmp_data []byte
	if data != nil {
		// First judge whether it is a string, if it is a string, then judge whether it belongs to json, if it already belongs to json, no longer convert json again
		// 首先判断是不是字符串，如果是字符串，则判断属不属于json，如果已经属于json，则不再二次转化json
		v, ok := data.(string)
		if ok {
			ok_json := json.Valid([]byte(v))
			if ok_json {
				tmp_data = []byte(v)
			} else {
				tmp_data, _ = json.Marshal(data)
			}
		} else {
			tmp_data, _ = json.Marshal(data)
		}
	}

	b, err := json.Marshal(
		Response{
			Code: code,
			Status:	"error",
			Message: message + "Error: " + err.Error(),
			Data: string(tmp_data),
		})

	return b, err
}

// Return success information in json format
// 返回json格式的成功信息
func defaultHandleRpcSuccessFunc(code string, message string, data interface{}) (*em_protobuf.Response, error) {
	// First judge whether it is a string, if it is a string, then judge whether it belongs to json, if it already belongs to json, no longer convert json again
	// 首先判断是不是字符串，如果是字符串，则判断属不属于json，如果已经属于json，则不再二次转化json
	if v, ok := data.(string); ok {
		ok_json := json.Valid([]byte(v))
		if ok_json {
			return &em_protobuf.Response{
				Code:    code,
				Status:  SUCCESS_Status,
				Message: message,
				Data:    v,
			}, nil
		}
	}

	return &em_protobuf.Response{
		Code:    code,
		Status:  SUCCESS_Status,
		Message: message,
		Data:    MustConvertJson(data),
	}, nil
}

// Return error message in json format
// 返回json格式的错误信息
func defaultHandleRpcErrorFunc(rcpStatusCode codes.Code, code string, message string, data interface{}, err error) (*em_protobuf.Response, error) {
	// Data interface => Json string
	var tmp_data []byte
	if data != nil {
		// First judge whether it is a string, if it is a string, then judge whether it belongs to json, if it already belongs to json, no longer convert json again
		// 首先判断是不是字符串，如果是字符串，则判断属不属于json，如果已经属于json，则不再二次转化json
		v, ok := data.(string)
		if ok {
			ok_json := json.Valid([]byte(v))
			if ok_json {
				tmp_data = []byte(v)
			} else {
				tmp_data, _ = json.Marshal(data)
			}
		} else {
			tmp_data, _ = json.Marshal(data)
		}
	}

	// If enabled, use HTTP CODE instead of system default CODE
	// 如果开启使用HTTP CODE 代替系统的默认CODE
	e, err := Kv.ReadKey(KvAppUseHttpCode)
	if err != nil {
		LogInfo.OutputSimplePath(err)
	}
	if strings.ToLower(e) != "true" {
		code = strconv.Itoa(int(rcpStatusCode))
	}

	// If it is a Debug environment, return information with Error
	// 如果是Debug环境，返回带有Error的信息
	if IsDebug() && err != nil {
		err := status.Error(rcpStatusCode, Response{
			Code:    code,
			Status:  ERROR_Status,
			Message: message + "Error: " + err.Error(),
			Data:    string(tmp_data),
		}.String())
		return nil, err
	}

	err3 := status.Error(rcpStatusCode, Response{
		Code:    code,
		Status:  ERROR_Status,
		Message: message,
		Data:    string(tmp_data),
	}.String())
	return nil, err3
}

func defaultGrpcMiddlewareFunc() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// Panic recover
				grpc_recovery.UnaryServerInterceptor(),
				// I18n
				DefaultMiddleware().I18n(),
				// token Auth
				DefaultMiddleware().Auth(),
			),
		),
	)
	return s
}