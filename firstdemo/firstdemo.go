// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

// 入口文件
package main

import (
	"flag"
	"fmt"

	"firstdemo/internal/config"
	"firstdemo/internal/handler"
	"firstdemo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/firstdemo-api.yaml", "the config file")

func main() {
	flag.Parse()
	// 初始化配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
