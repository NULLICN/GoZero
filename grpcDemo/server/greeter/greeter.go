package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	"greeter/greeter"
	"greeter/internal/config"
	"greeter/internal/server"
	"greeter/internal/svc"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client/naming_cache"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client/naming_grpc"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/v2/common/nacos_server"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
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

	// 用 AdvertisedIp 覆盖自动探测的私网 IP，注册到 Nacos（持久化实例）
	advertiseAddr := c.ListenOn
	if c.AdvertisedIp != "" {
		_, port, _ := strings.Cut(c.ListenOn, ":")
		advertiseAddr = c.AdvertisedIp + ":" + port
	}

	host, portStr, err := net.SplitHostPort(advertiseAddr)
	if err != nil {
		panic(fmt.Sprintf("failed parsing advertise address %s: %v", advertiseAddr, err))
	}
	port, _ := strconv.ParseUint(portStr, 10, 16)

	groupName := c.Nacos.GroupName
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}

	// 注册持久化实例 (Ephemeral: false) — 通过 gRPC 协议直连，绕过 SDK HTTP 代理
	grpcProxy := registerPersistent(c.RpcServerConf.Name, groupName, host, port, sc, *cc)
	defer grpcProxy.CloseClient()

	proc.AddShutdownListener(func() {
		_, err := grpcProxy.DeregisterInstance(c.RpcServerConf.Name, groupName, model.Instance{
			Ip:        host,
			Port:      port,
			Ephemeral: false,
		})
		if err != nil {
			logx.Infof("deregister service error: %v", err)
		} else {
			logx.Infof("deregistered persistent instance from nacos")
		}
	})

	fmt.Printf("Starting rpc server at %s (registered as %s)...\n", c.ListenOn, advertiseAddr)
	s.Start()
}

// registerPersistent 通过 SDK gRPC 代理注册持久化实例。
// Nacos SDK 默认将 Ephemeral:false 路由到 HTTP API，
// 但部分 Nacos 部署不支持 HTTP naming API (返回 501)。
// 此处绕过代理路由，直接使用 gRPC 代理完成持久化注册。
func registerPersistent(
	serviceName, groupName, ip string, port uint64,
	serverConfigs []constant.ServerConfig,
	clientCfg constant.ClientConfig,
) *naming_grpc.NamingGrpcProxy {

	serviceInfoHolder := naming_cache.NewServiceInfoHolder(
		clientCfg.NamespaceId,
		clientCfg.CacheDir,
		clientCfg.UpdateCacheWhenEmpty,
		clientCfg.NotLoadCacheAtStart,
	)

	nacosServer, err := nacos_server.NewNacosServer(
		context.Background(),
		serverConfigs,
		clientCfg,
		&http_agent.HttpAgent{},
		clientCfg.TimeoutMs,
		"",   // endpoint — 不使用 address-server 模式
		nil,  // endpoint query headers
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating nacos server: %v", err))
	}

	grpcProxy, err := naming_grpc.NewNamingGrpcProxy(
		context.Background(),
		clientCfg,
		nacosServer,
		serviceInfoHolder,
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating grpc proxy: %v", err))
	}

	_, err = grpcProxy.RegisterInstance(serviceName, groupName, model.Instance{
		Ip:        ip,
		Port:      port,
		Weight:    1.0,
		Enable:    true,
		Healthy:   true,
		Ephemeral: false, // 持久化实例
		Metadata:  map[string]string{},
	})
	if err != nil {
		panic(fmt.Sprintf("failed registering persistent instance: %v", err))
	}

	logx.Infof("registered persistent instance %s:%d to nacos service %s group %s (gRPC)", ip, port, serviceName, groupName)
	return grpcProxy
}
