package interceptor

import (
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/silentred/echorus"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewLogInterceptor(logger *echorus.Echorus) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		stop := time.Now()

		json := log.JSON{
			"time_unix":   strconv.FormatInt(time.Now().Unix(), 10),
			"method":      info.FullMethod,
			"latency":     strconv.FormatInt(int64(stop.Sub(start)), 10),
			"latency_str": stop.Sub(start).String(),
			"req":         marshal(req),
			"resp":        marshal(resp),
			"err":         err,
		}

		logger.Infoj(json)

		return resp, err
	}
}
