package main

import (
	"fmt"
	"proto_demo/userService"

	"google.golang.org/protobuf/proto"
)

func main() {
	// 1. 创建 Protobuf 对象并赋值
	u := &userService.Userinfo{
		Username: "张三",
		Age:      20,
		Hobby:    []string{"吃饭", "睡觉", "写代码"},
	}

	// 打印原始字段
	fmt.Println(u.GetUsername())
	fmt.Println(u.GetHobby())

	// 2. Protobuf 的序列化 (数据 -> 二进制)
	data, _ := proto.Marshal(u)
	fmt.Println(data) // 打印二进制字节切片

	// 3. Protobuf 的反序列化 (二进制 -> 数据)
	user := userService.Userinfo{}
	proto.Unmarshal(data, &user)

	// 打印反序列化后的结构体
	fmt.Printf("%#v\n", &user) // 使用 %#v 可以打印更详细的结构体信息
	fmt.Println(user.GetHobby())
}