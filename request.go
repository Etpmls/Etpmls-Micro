package em

import (
	"context"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
)

type request struct {

}

// Get token from ctx
// 从ctx获取令牌
func (this *request) GetValueFromCtx(ctx context.Context, value string) (string, error) {
	if ctx == nil {
		return "", LogError.OutputAndReturnError(MessageWithLineNum("Failed to obtain request!"))
	}

	v := ctx.Value(value);
	if v == nil {
		return "", LogInfo.OutputAndReturnError(MessageWithLineNum("Failed to obtain " + value +"!"))
	}

	return v.(string), nil
}

// Get value from header
// 从header获取值
func (this *request) Rpc_GetValueFromHeader(ctx context.Context, value string) (string, error) {
	if ctx == nil {
		return "", LogError.OutputAndReturnError(MessageWithLineNum("Failed to obtain " + value +"! Context is nil"))
	}

	// 1.Get header from grpc-gateway
	// 从grpc-gateway获取header
	g := em_library.NewGrpc()
	v, err := g.ExtractHeader(ctx, value)
	if err != nil {
		// 2.Get header from metadata
		// 从metadata中获取header
		v2, err2 := g.GetFirstValueFromMetadata(ctx, value)
		if err2 == nil {
			return v2, nil
		}

		return "", LogInfo.OutputAndReturnError(MessageWithLineNum(err.Error()))
	}
	return v, nil
}

// Set the value to the header
// 向header中设置值
func (this *request) Rpc_SetValueToHeader(ctx *context.Context, m map[string]string) {
	g := em_library.NewGrpc()
	g.SetValueToMetadata(ctx, m)
	return
}