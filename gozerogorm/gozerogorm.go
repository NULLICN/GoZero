package main

import (
	"flag"
	"fmt"

	"gozerogorm/internal/config"
	"gozerogorm/internal/handler"
	"gozerogorm/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gozerogorm-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("启动服务于 %s:%d...\n", c.Host, c.Port)
	server.Start()
}
