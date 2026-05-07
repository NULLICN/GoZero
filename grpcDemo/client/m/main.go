package main

import (
	"client/greeter"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	// 注册客户端
	client := greeter.NewGreeterClient(conn)

	res, err := client.SayHello(context.Background(), &greeter.HelloReq{
		Name: "nullicn",
	})
	fmt.Printf("%#v\r\n", res)
	fmt.Println(res.Message)
	defer conn.Close()

}
