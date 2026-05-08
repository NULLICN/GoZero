package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"gozerogorm/internal/config"
	"gozerogorm/internal/handler"
	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"

	"github.com/zeromicro/go-zero/core/conf"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/x/errors"
)

//go:embed etc/gozerogorm-api.yaml
var embeddedConfig []byte

var configFile = flag.String("f", "", "config file path (optional, uses embedded config by default)")

func main() {
	flag.Parse()

	var c config.Config
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		if err := conf.LoadFromYamlBytes(embeddedConfig, &c); err != nil {
			log.Fatalf("load embedded config error: %v", err)
		}
	}

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(
		func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.CommonResponse{
				Success: false,
				Code:    401,
				Message: "JWT认证失败：" + err.Error(),
			})
		},
	))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 统一的错误处理 https://go-zero.dev/zh-cn/guides/http/server/error/
	httpx.SetErrorHandler(func(err error) (int, any) {
		switch e := err.(type) {
		case *errors.CodeMsg:
			return http.StatusOK, types.CommonResponse{
				Code:    e.Code,
				Success: false,
				Message: "错误原因：" + e.Msg,
			}
		default:
			return http.StatusInternalServerError, types.CommonResponse{
				Code:    500,
				Success: false,
				Message: "错误原因：" + e.Error(),
			}
		}
	})

	fmt.Printf("启动服务于 %s:%d...\n", c.Host, c.Port)
	server.Start()
}
