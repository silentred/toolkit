package interceptor

import (
	"runtime"

	"github.com/silentred/toolkit/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	// MaxStackSize when recovering panic
	MaxStackSize = 4096
)

// NewRecovery return a recover interceptor for gRPC
func NewRecovery(logger util.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// recovery func
		defer func() {
			if r := recover(); r != nil {
				// log stack
				stack := make([]byte, MaxStackSize)
				stack = stack[:runtime.Stack(stack, false)]
				logger.Errorf("panic grpc invoke: %s, err=%v, stack:\n%s", info.FullMethod, r, string(stack))

				// if panic, set custom error to 'err', in order that client and sense it.
				err = grpc.Errorf(codes.Internal, "panic error: %v", r)
			}
		}()

		return handler(ctx, req)
	}
}
