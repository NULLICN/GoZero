package main

import (
	"flag"
	"fmt"
	"strings"

	"greeter/greeter"
	"greeter/internal/config"
	"greeter/internal/server"
	"greeter/internal/svc"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/greeter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Nacos.Ip, c.Nacos.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         c.Nacos.Namespace,
		TimeoutMs:           50000,
		NotLoadCacheAtStart: c.Nacos.NotLoadCacheAtStart,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            c.Nacos.LogLevel,
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		greeter.RegisterGreeterServer(grpcServer, server.NewGreeterServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 用 AdvertisedIp 覆盖自动探测的私网 IP，注册到 Nacos
	advertiseAddr := c.ListenOn
	if c.AdvertisedIp != "" {
		_, port, _ := strings.Cut(c.ListenOn, ":")
		advertiseAddr = c.AdvertisedIp + ":" + port
	}
	opts := nacos.NewNacosConfig(c.RpcServerConf.Name, advertiseAddr, sc, cc)
	nacos.RegisterService(opts)

	fmt.Printf("Starting rpc server at %s (registered as %s)...\n", c.ListenOn, advertiseAddr)
	s.Start()
}
