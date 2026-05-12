// m4: 生产级 API 网关 → RPC 调用链测试
// 监听 :8081，对外暴露 HTTP /sayhello，内部经过 Nacos 发现 greeter.rpc 并转发
// 调用链: HTTP client → [m4 API Gateway] → greeter.rpc
//
// 用法: m4.exe [-f etc/m4.yaml]
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"client/greeterclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
)

// ---------- config ----------

type NacosConfig struct {
	Ip                  string
	Port                uint64
	Namespace           string
	NotLoadCacheAtStart bool
	LogLevel            string
}

type Config struct {
	rest.RestConf
	RpcClient zrpc.RpcClientConf
	Nacos     NacosConfig
}

// ---------- request / response DTO ----------

type SayHelloReq struct {
	Name string `json:"name"`
}

type SayHelloRes struct {
	Message string `json:"message"`
}

// ---------- service context ----------

type ServiceContext struct {
	Config     Config
	GreeterRpc greeterclient.Greeter
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{Config: c}
}

// ---------- handler ----------

func makeSayHelloHandler(svcCtx *ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SayHelloReq
		if r.Method == http.MethodGet {
			req.Name = r.URL.Query().Get("name")
		} else {
			if err := httpx.ParseJsonBody(r, &req); err != nil {
				httpx.Error(w, err)
				return
			}
		}

		if req.Name == "" {
			req.Name = "anonymous"
		}

		logx.Infof("[gateway] incoming request name=%q", req.Name)

		resp, err := callRpcWithRetry(r.Context(), svcCtx.GreeterRpc, req.Name)
		if err != nil {
			logx.Errorf("[gateway] rpc call failed: %v", err)
			httpx.WriteJson(w, http.StatusBadGateway, SayHelloRes{
				Message: fmt.Sprintf("rpc error: %v", err),
			})
			return
		}

		logx.Infof("[gateway] rpc response: %s", resp.Message)
		httpx.WriteJson(w, http.StatusOK, SayHelloRes{Message: resp.Message})
	}
}

// ---------- RPC 调用 (带超时+重试) ----------

func callRpcWithRetry(ctx context.Context, client greeterclient.Greeter, name string) (*greeterclient.HelloRes, error) {
	const (
		maxRetries  = 3
		callTimeout = 3 * time.Second
	)

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * 100 * time.Millisecond
			logx.Infof("[gateway] retry %d/%d after %v", attempt, maxRetries-1, backoff)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		callCtx, cancel := context.WithTimeout(ctx, callTimeout)
		resp, err := client.SayHello(callCtx, &greeterclient.HelloReq{Name: name})
		cancel()

		if err != nil {
			lastErr = err
			logx.Errorf("[gateway] call failed (attempt %d/%d): %v", attempt+1, maxRetries, err)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("all %d attempts exhausted: %w", maxRetries, lastErr)
}

// ---------- health check ----------

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status":"ok"}`)
}

// ---------- main ----------

var configFile = flag.String("f", "etc/m4.yaml", "the config file")

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(logx.LogConf{})
	logx.Infof("starting m4 API gateway on %s:%d", c.Host, c.Port)

	// 1. 初始化 RPC 客户端 (通过 Nacos 发现 greeter.rpc)
	rpcConn := zrpc.MustNewClient(c.RpcClient)
	defer rpcConn.Conn().Close()

	if !waitRpcReady(rpcConn, 10*time.Second) {
		logx.Error("rpc connection not ready after timeout")
		os.Exit(1)
	}

	svcCtx := NewServiceContext(c)
	svcCtx.GreeterRpc = greeterclient.NewGreeter(rpcConn)

	// 2. 构建 HTTP 引擎
	engine := rest.MustNewServer(c.RestConf,
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpx.WriteJson(w, http.StatusNotFound, map[string]string{
				"error":   "not found",
				"message": fmt.Sprintf("no route for %s %s", r.Method, r.URL.Path),
			})
		})),
	)
	defer engine.Stop()

	// 3. 注册路由
	sayHello := makeSayHelloHandler(svcCtx)
	engine.AddRoute(rest.Route{Method: http.MethodGet, Path: "/sayhello", Handler: sayHello})
	engine.AddRoute(rest.Route{Method: http.MethodPost, Path: "/sayhello", Handler: sayHello})
	engine.AddRoute(rest.Route{Method: http.MethodGet, Path: "/health", Handler: healthHandler})

	// 4. 优雅关闭 (signal handler 与 go-zero 内置的 shutdown 叠加 — 双保险)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		logx.Infof("received signal %s, shutting down...", sig)
		engine.Stop()
	}()

	// 5. 启动
	fmt.Printf("m4 API gateway listening on http://%s:%d\n", c.Host, c.Port)
	fmt.Printf("  GET/POST http://localhost:%d/sayhello?name=xxx\n", c.Port)
	fmt.Printf("  GET  http://localhost:%d/health\n", c.Port)
	engine.Start()
}

func waitRpcReady(conn zrpc.Client, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if conn.Conn().GetState().String() == "READY" {
			logx.Info("rpc connection ready")
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}
