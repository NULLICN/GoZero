package main

import (
	"context"
	"fmt"
	"greeter/greeter"
	"net"

	"google.golang.org/grpc"
)

type Hello struct {
	greeter.UnimplementedGreeterServer
}

func (this Hello) SayHello(c context.Context, req *greeter.HelloReq) (*greeter.HelloRes, error) {
	fmt.Println(req)
	return &greeter.HelloRes{
		Message: "你好" + req.Name,
	}, nil
}

func main() {
	// 1.初始化一个grpc对象
	grpcServer := grpc.NewServer()
	// 2.注册服务
	greeter.RegisterGreeterServer(grpcServer, new(Hello))
	// 3.监听服务
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	grpcServer.Serve(listener)
}
