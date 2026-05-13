// m4api: 标准 go-zero API 网关 → RPC 调用链测试
//
// 分层架构 (由 goctl 生成 + 手工定制):
//
//	m4api.go                    入口: 加载配置、创建 RPC 客户端、启动 HTTP 服务
//	internal/config/config.go   配置定义 (HTTP + RPC + Nacos)
//	internal/svc/servicecontext.go  服务上下文 (持有 greeterclient.Greeter)
//	internal/handler/routes.go  路由注册 (goctl 生成, DO NOT EDIT)
//	internal/handler/*.go       HTTP handler — 解析请求 → 调用 logic
//	internal/logic/*.go         业务逻辑 — 调用 greeterclient → RPC
//	internal/types/types.go     请求/响应 DTO
//
// 调用链:
//
//	HTTP client → [m4api :8081] → greeterclient → Nacos → greeter.rpc
//
// 用法: m4api.exe [-f etc/m4api.yaml]
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"client/greeterclient"
	"client/m4api/internal/config"
	"client/m4api/internal/handler"
	"client/m4api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
)

var configFile = flag.String("f", "etc/m4api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(logx.LogConf{})
	logx.Infof("starting m4api API gateway on %s:%d", c.Host, c.Port)

	// 0. Nacos SDK 环境变量 (zero-contrib 插件原生支持, URL 参数因 bool 字段缺 string 标签无法传递)
	os.Setenv("NACOS_NOT_LOAD_CACHE_AT_START", "true")

	// 1. 初始化 RPC 客户端 (通过 Nacos 自动发现 greeter.rpc)
	// NonBlock:true 下连接懒加载，首次 RPC 调用时自动建连
	rpcConn := zrpc.MustNewClient(c.RpcClient)
	defer rpcConn.Conn().Close()

	// 2. 装配服务上下文 (注入 RPC 客户端)
	ctx := svc.NewServiceContext(c, greeterclient.NewGreeter(rpcConn))

	// 3. 创建 HTTP 引擎 + 注册路由
	server := rest.MustNewServer(c.RestConf,
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpx.WriteJson(w, http.StatusNotFound, map[string]string{
				"error":   "not found",
				"message": fmt.Sprintf("no route for %s %s", r.Method, r.URL.Path),
			})
		})),
	)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/health",
		Handler: healthHandler,
	})

	// 4. 启动 (go-zero 内置 SIGINT/SIGTERM 优雅关闭)
	fmt.Printf("m4api API gateway listening on http://%s:%d\n", c.Host, c.Port)
	fmt.Printf("  GET/POST http://localhost:%d/sayhello?name=xxx\n", c.Port)
	fmt.Printf("  GET  http://localhost:%d/health\n", c.Port)
	server.Start()
}

// ---------- health check ----------

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status":"ok"}`)
}

