package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	pb "go-grpc-example/proto/hello" // 引入proto包

	"go-grpc-example/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "hello", "service name")
	reg  = flag.String("reg", "http://192.168.1.5:2379", "register etcd address")
)

func main() {
	flag.Parse()
	r := etcd.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewHelloClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
	}
}
