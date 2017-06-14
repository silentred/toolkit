package main

import (
	"flag"
	"log"

	"github.com/silentred/toolkit/example/grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	Host string
)

func init() {
	flag.StringVar(&Host, "h", "127.0.0.1:28080", "gRPC server address")
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(Host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	resp, err := c.SayHello(context.Background(), &proto.HelloReq{Name: "world", Times: 100})
	if err != nil {
		log.Printf("err:%v \n", err)
	}

	log.Printf("resp:%s \n", resp)
}
