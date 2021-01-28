package em

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"net/http"
)

type endpoint struct {

}

func NewEndpoint() *endpoint {
	return &endpoint{}
}

// https://grpc-ecosystem.github.io/grpc-gateway/docs/customizingyourgateway.html
// => https://mycodesmells.com/post/grpc-gateway-error-handler
func (this *endpoint) CustomErrorHandlerFunc (_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))
	w.Write([]byte(status.Convert(err).Message()))
}


func (this *endpoint) SetCustomMatcher(key string) (string, bool) {
	switch key {
	case "Token":
		return key, true
	case "Language":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func (this *endpoint) CustomRoutingErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	LogWarn.Output("httpStatus:", httpStatus , " request:", r)
	runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r ,httpStatus)
}
