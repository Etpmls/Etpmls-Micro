package em

import (
	"context"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	em_utils "github.com/Etpmls/Etpmls-Micro/utils"
)

type Request struct {

}

// Get value from header
// 从header获取值
func (this *Request) GetValueFromHeader(ctx context.Context, value string) (string, error) {
	if ctx == nil {
		return "", LogError.OutputAndReturnError((em_utils.MessageWithLineNum("Failed to obtain " + value +"! Context is nil")))
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

		return "", LogInfo.OutputAndReturnError((em_utils.MessageWithLineNum(err.Error())))
	}
	return v, nil
}

// Get token from ctx
// 从ctx获取令牌
func (this *Request) GetValueFromCtx(ctx context.Context, value string) (string, error) {
	if ctx == nil {
		return "", LogError.OutputAndReturnError((em_utils.MessageWithLineNum("Failed to obtain request!")))
	}

	v := ctx.Value(value);
	if v == nil {
		return "", LogInfo.OutputAndReturnError((em_utils.MessageWithLineNum("Failed to obtain " + value +"!")))
	}

	return v.(string), nil
}
