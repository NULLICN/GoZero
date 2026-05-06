package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// 1. 用 rpc 链接服务器 --Dial()
	// conn, err := rpc.Dial("tcp", "47.109.80.234:8080")

	// json方式 需用net.Dial
	conn, err := net.Dial("tcp", "47.109.80.234:8080")

	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}
	defer conn.Close()

	// 2. 调用远程函数
	var reply string // 接受返回值 --- 传出参数

	// json方式
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// err = conn.Call("hello.SayHello", "张三", &reply)
	err = client.Call("hello.SayHello", "张三", &reply)
	if err != nil {
		fmt.Println("Call:", err)
		return
	}
	fmt.Println(reply)
}
