package em_library

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type grpc struct {

}

func NewGrpc() *grpc {
	return &grpc{}
}

func (this *grpc) GetServiceName(fullMethodName string) string {
	strs := strings.Split(fullMethodName, ".")
	return strs[len(strs)-1]
}

// https://github.com/johanbrandhorst/grpc-auth-example
// Get the specified header from the request
// 从请求中获取指定header
func (this *grpc) ExtractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}

func (this *grpc) SetValueToMetadata(ctx *context.Context, m map[string]string) {
	value := metadata.New(m)
	*ctx = metadata.NewOutgoingContext(*ctx, value)
	return
}

func (this *grpc) GetValueFromMetadata(ctx context.Context, key string) ([]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("Get metadata failed!")
	}
	if len(md.Get(key)) == 0 {
		return nil, errors.New("metadata no record")
	}
	return md.Get(key), nil
}

func (this *grpc) GetFirstValueFromMetadata(ctx context.Context, key string) (string, error) {
	sl, err := this.GetValueFromMetadata(ctx, key)
	if err != nil {
		return "", err
	}

	return sl[0], nil
}