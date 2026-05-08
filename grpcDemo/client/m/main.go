package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"client/greeter"

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

var configFile = flag.String("f", "etc/client.yaml", "the config file")

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	conn := zrpc.MustNewClient(c.RpcClientConf)
	defer conn.Conn().Close()

	client := greeter.NewGreeterClient(conn.Conn())

	res, err := client.SayHello(context.Background(), &greeter.HelloReq{
		Name: "nullicn",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Message)
}
