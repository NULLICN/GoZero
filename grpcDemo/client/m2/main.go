// 对比示例：使用 go-zero 风格 greeterclient 包装器的客户端
// 与 m/main.go (原始 protobuf 客户端) 做对比
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"client/greeterclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
)

type NacosConfig struct {
	Ip                  string
	Port                uint64
	Namespace           string
	NotLoadCacheAtStart bool
	LogLevel            string
}

type Config struct {
	zrpc.RpcClientConf
	Nacos NacosConfig
}

var configFile = flag.String("f", "../etc/client.yaml", "the config file")

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	conn := zrpc.MustNewClient(c.RpcClientConf)
	defer conn.Conn().Close()

	// ====== 与 m/main.go 的区别在这里 ======
	// 原始方式: greeter.NewGreeterClient(conn.Conn())
	// 包装方式: greeterclient.NewGreeter(conn)  — 传 conn 而非 conn.Conn()
	client := greeterclient.NewGreeter(conn)

	// 原始方式: &greeter.HelloReq{Name: "nullicn"}
	// 包装方式: &greeterclient.HelloReq{...}  — 类型别名，不用多引一个 greeter 包
	res, err := client.SayHello(context.Background(), &greeterclient.HelloReq{
		Name: "nullicn",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Message)
}
