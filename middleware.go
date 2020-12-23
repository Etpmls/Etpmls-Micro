package em

import (
	"context"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type defaultMiddleware struct {

}

func DefaultMiddleware() *defaultMiddleware {
	return &defaultMiddleware{}
}

func (this *defaultMiddleware) I18n() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Get language
		// 获取语言
		g := em_library.NewGrpc()
		lang, err := g.ExtractHeader(ctx, "language")
		if err != nil {
			return handler(ctx, req)
		}

		// Pass the language to the method
		// 把language传递到方法中
		ctx = context.WithValue(ctx,"language", lang)

		return handler(ctx, req)
	}
}

func (this *defaultMiddleware) SetCors(mux *runtime.ServeMux, options cors.Options) http.Handler {
	// CORS
	// https://github.com/rs/cors
	c := cors.New(options)
	// Insert the defaultMiddleware
	return c.Handler(mux)
}

// Only Verify Token
// 仅验证token
// Ensure the security of the intranet (without going through the API gateway)
// 保证内网安全（不经过API网关）
func (this *defaultMiddleware) Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Get token from header
		token, err:= Micro.Auth.Rpc_GetTokenFromHeader(ctx)
		if err != nil || token == "" {
			return nil, status.Error(codes.PermissionDenied, I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
		}

		b, _ := Micro.Auth.VerifyToken(token)
		if !b {
			return nil, status.Error(codes.PermissionDenied, I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
		}

		// Pass the token to the method
		// 把token传递到方法中
		ctx = context.WithValue(ctx,"token", token)

		LogDebug.Output("Auth defaultMiddleware runs successfully!")	// Debug
		return handler(ctx, req)
	}
}

type middleware struct {

}

// Implement http middleware
// 实现http中间件
type MiddlewareFunc func(http.ResponseWriter, *http.Request, map[string]string) error
func (this *middleware) WithMiddleware(f runtime.HandlerFunc, middlware ...MiddlewareFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		for _, v := range middlware {
			err := v(w, r, pathParams)
			if err != nil {
				return
			}
		}

		f(w, r, pathParams)
	}
}


/*func (this *defaultMiddleware) Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// fullMethodName: /protobuf.User/GetCurrent
		service := em_library.NewGrpc().GetServiceName(info.FullMethod)

		token, err:= NewAuth().GetTokenFromHeader(ctx)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
		}

		err = NewAuth().VerifyPermissions(token, service)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
		}

		// Pass the token to the method
		// 把token传递到方法中
		ctx = context.WithValue(ctx,"token", token)

		return handler(ctx, req)
	}
}



func HttpVerifyToken(w http.ResponseWriter, r *http.request, pathParams map[string]string) error  {
	token := r.Header.Get("token")
	if len(token) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission Denied"))
		return errors.New("Get token error")
	}


	b, _ := NewAuth().VerifyToken(token)
	if b {
		return nil
	}

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Permission Denied"))
	return errors.New("Permission Denied")
}

func HttpVerifyPermissions(w http.ResponseWriter, r *http.request, pathParams map[string]string) error  {
	token := r.Header.Get("token")
	if len(token) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission Denied"))
		return errors.New("Get token error")
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission Denied"))
		return err
	}

	err = NewAuth().VerifyPermissions(token, u.Path)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission Denied"))
		return err
	}

	return nil
}*/