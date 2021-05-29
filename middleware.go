package em

import (
	"context"
	em_library "github.com/Etpmls/Etpmls-Micro/v3/library"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type middleware struct {

}

func Middleware() *middleware {
	return &middleware{}
}

func (this *middleware) Translate() grpc.UnaryServerInterceptor {
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

// Only Verify Token
// 仅验证token
// Ensure the security of the intranet (without going through the API gateway)
// 保证内网安全（不经过API网关）
func (this *middleware) Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Get token from header
		token, err:= Micro.Auth.GetTokenFromHeader(ctx)
		if err != nil || token == "" {
			return nil, status.Error(codes.PermissionDenied, "Permission Denied")
		}

		b, _ := Micro.Auth.VerifyToken(token)
		if !b {
			return nil, status.Error(codes.PermissionDenied, "Permission Denied")
		}

		// Pass the token to the method
		// 把token传递到方法中
		ctx = context.WithValue(ctx,"token", token)

		LogDebug.New("Auth middleware runs successfully!") // Debug
		return handler(ctx, req)
	}
}

