package em

import (
	"context"
	"errors"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/url"
)

type middleware struct {

}

func NewMiddleware() *middleware {
	return &middleware{}
}

func (this *middleware) I18n() grpc.UnaryServerInterceptor {
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

func (this *middleware) SetCors(mux *runtime.ServeMux, options cors.Options) http.Handler {
	// CORS
	// https://github.com/rs/cors
	c := cors.New(options)
	// Insert the middleware
	return c.Handler(mux)
}

func (this *middleware) Auth() grpc.UnaryServerInterceptor {
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

type MiddlewareFunc func(http.ResponseWriter, *http.Request, map[string]string) error
func WithMiddleware(f runtime.HandlerFunc, middlware ...MiddlewareFunc) runtime.HandlerFunc {
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

func HttpVerifyToken(w http.ResponseWriter, r *http.Request, pathParams map[string]string) error  {
	token := r.Header.Get("token")
	if len(token) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission Denied"))
		return errors.New("Get token error")
	}


	err := NewAuth().VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission Denied"))
		return err
	}

	return nil
}

func HttpVerifyPermissions(w http.ResponseWriter, r *http.Request, pathParams map[string]string) error  {
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
}