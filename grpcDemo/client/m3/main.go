// 对比示例：使用 nacos-sdk-go/v2 直接做服务发现，再直连 gRPC 的客户端
// 与 m/main.go (go-zero nacos:// resolver) 和 m2/main.go (go-zero wrapper) 做对比

// m3.exe -name nullicn -count 1
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"client/greeter"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type NacosConfig struct {
	Ip                  string
	Port                uint64
	Namespace           string
	GroupName           string 
	NotLoadCacheAtStart bool
	LogLevel            string
}

type Config struct {
	Nacos NacosConfig
}

var (
	configFile = flag.String("f", "../etc/client.yaml", "the config file")
	name       = flag.String("name", "nullicn", "name to greet")
	count      = flag.Int("count", 1, "number of concurrent calls (0 for infinite loop)")
	interval   = flag.Duration("interval", time.Second, "interval between calls in loop mode")
	timeout    = flag.Duration("timeout", 3*time.Second, "call timeout")
)

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(logx.LogConf{})

	groupName := c.Nacos.GroupName
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}

	// 1. 创建 Nacos 服务发现客户端
	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(c.Nacos.Ip, c.Nacos.Port),
		},
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         c.Nacos.Namespace,
			TimeoutMs:           50000,
			NotLoadCacheAtStart: c.Nacos.NotLoadCacheAtStart,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            c.Nacos.LogLevel,
		},
	})
	if err != nil {
		logx.Errorf("创建 Nacos 客户端失败: %v", err)
		os.Exit(1)
	}

	// 2. 从 Nacos 获取一个健康的服务实例
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "greeter.rpc",
		GroupName:   groupName,
	})
	if err != nil {
		logx.Errorf("服务发现失败: %v", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
	logx.Infof("通过 Nacos 发现服务实例: %s", addr)

	// 3. 直连发现的实例 (不使用 go-zero 的 nacos:// resolver)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logx.Errorf("连接失败: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 触发懒连接从 Idle → Connecting → Ready 的状态迁移
	conn.Connect()
	if !waitReady(conn, 10*time.Second) {
		logx.Error("connection not ready after timeout")
		os.Exit(1)
	}

	// 4. 使用原始 proto 客户端 (与 m2 的 greeterclient 包装器做对比)
	client := greeter.NewGreeterClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 信号处理：优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logx.Info("received shutdown signal")
		cancel()
	}()

	if *count == 0 {
		runLoop(ctx, client, *interval)
	} else {
		runBatch(ctx, client, *count)
	}
}

func waitReady(conn *grpc.ClientConn, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if conn.GetState() == connectivity.Ready {
			logx.Info("connection ready")
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func runBatch(ctx context.Context, client greeter.GreeterClient, n int) {
	sem := make(chan struct{}, 10) // 最多 10 个并发
	errCh := make(chan error, n)

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			logx.Infof("batch cancelled after %d/%d calls", i, n)
			return
		case sem <- struct{}{}:
		}

		go func(idx int) {
			defer func() { <-sem }()
			if err := callWithRetry(ctx, client, fmt.Sprintf("%s-%d", *name, idx)); err != nil {
				errCh <- err
			}
		}(i)
	}

	// 等待所有完成
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	close(errCh)

	failures := 0
	for range errCh {
		failures++
	}
	logx.Infof("batch complete: %d calls, %d failures", n, failures)
}

func runLoop(ctx context.Context, client greeter.GreeterClient, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logx.Info("loop stopped")
			return
		case <-ticker.C:
			callWithRetry(ctx, client, *name)
		}
	}
}

func callWithRetry(ctx context.Context, client greeter.GreeterClient, name string) error {
	const maxRetries = 3

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * 100 * time.Millisecond
			logx.Infof("retry %d/%d for %s after %v", attempt, maxRetries-1, name, backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		callCtx, cancel := context.WithTimeout(ctx, *timeout)
		res, err := client.SayHello(callCtx, &greeter.HelloReq{Name: name})
		cancel()

		if err != nil {
			lastErr = err
			logx.Errorf("call failed (attempt %d/%d): %v", attempt+1, maxRetries, err)
			continue
		}

		logx.Infof("response: %s", res.Message)
		return nil
	}

	logx.Errorf("all %d attempts exhausted for %s: %v", maxRetries, name, lastErr)
	return lastErr
}
