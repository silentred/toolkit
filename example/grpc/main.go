package main

import (
	"github.com/silentred/toolkit/example/grpc/proto"
	"github.com/silentred/toolkit/interceptor"
	"github.com/silentred/toolkit/service"
	"github.com/silentred/toolkit/service/discovery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type helloSvc struct{}

// Hello implements helloworld.GreeterServer
func (s *helloSvc) SayHello(ctx context.Context, in *proto.HelloReq) (*proto.HelloResp, error) {
	// log ctx values
	return &proto.HelloResp{Message: "Hello " + in.Name}, nil
}

func main() {
	app := service.NewGrpcApp(nil)
	app.RegisterHook(service.ServiceHook, initService)
	app.Initialize()

	// create grpc server
	chain := interceptor.UnaryInterceptorChain(interceptor.NewRecovery(app.DefaultLogger()),
		interceptor.NewLogInterceptor(app.DefaultLogger()))
	opt := grpc.UnaryInterceptor(chain)
	s := grpc.NewServer(opt)

	// register
	hello := &helloSvc{}
	proto.RegisterGreeterServer(s, hello)

	// set service
	app.SetServer(s)

	app.ListenAndServe()
}

func initService(app service.Application) error {
	s := discovery.NewService("hello", "127.0.0.1", app.GetConfig().Port)
	p := discovery.NewEtcdPublisher([]string{"http://localhost:2379"}, 10)
	app.Inject(p)
	p.Register(s)
	go p.Heartbeat(s)

	return nil
}
