// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gozeroapi/internal/config"
	"gozeroapi/internal/handler"
	"gozeroapi/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gozero-api.yaml", "the config file")

func main() {
	flag.Parse()

	// 设置时区为中国标准时间（UTC+8）
	_ = os.Setenv("TZ", "Asia/Shanghai")
	time.Local = time.FixedZone("CST", 8*3600)

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
