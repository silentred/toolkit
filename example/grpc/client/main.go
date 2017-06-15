package main

import (
	"flag"
	"log"
	"time"

	"github.com/silentred/toolkit/example/grpc/proto"
	"github.com/silentred/toolkit/service/discovery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	host     string
	svcName  string
	etcd     bool
	etcdHost string
	loop     bool
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1:28080", "gRPC server address")
	flag.StringVar(&svcName, "svc", "hello", "gRPC service name")
	flag.StringVar(&etcdHost, "etcdHost", "http://127.0.0.1:2379", "etcd host")
	flag.BoolVar(&etcd, "etcd", false, "use etcd as registry")
	flag.BoolVar(&loop, "loop", false, "call rpc every 0.5s forever")
}

func main() {
	flag.Parse()

	var err error
	var conn *grpc.ClientConn

	if etcd {
		resolver := discovery.NewResolver(svcName, discovery.DefaultPrefix)
		opt := grpc.WithBalancer(grpc.RoundRobin(resolver))
		conn, err = grpc.Dial(etcdHost, grpc.WithInsecure(), opt)
	} else {
		conn, err = grpc.Dial(host, grpc.WithInsecure())
	}

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	callFunc := func() {
		resp, err := c.SayHello(context.Background(), &proto.HelloReq{Name: "world", Times: 100})
		if err != nil {
			log.Printf("err:%v \n", err)
		}
		log.Printf("resp:%s \n", resp)
		time.Sleep(time.Second / 2)
	}

	callFunc()
	for loop {
		callFunc()
	}
}
