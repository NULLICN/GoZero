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
	"greeter/internal/nacosconfig"
	"greeter/internal/netutil"
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

var configFile = flag.String("f", "etc/greeter.yaml", "the bootstrap config file")

func main() {
	flag.Parse()

	// ================================================================
	// Step 1: 加载本地 bootstrap 配置
	// bootstrap 只包含连接 Nacos 配置中心所需的信息（IP/端口/Namespace/DataId/Group）
	// 真正的业务配置存放在 Nacos 配置中心，下一步从那里拉取
	// ================================================================
	var bs config.Bootstrap
	conf.MustLoad(*configFile, &bs)

	if bs.ConfigDataId == "" || bs.ConfigGroup == "" {
		panic("bootstrap: ConfigDataId and ConfigGroup are required")
	}

	// ================================================================
	// Step 2: 从 Nacos 配置中心拉取主配置
	// 通过 Nacos Config Client 连接配置中心，根据 DataId + Group 获取 YAML 配置内容
	// 获取到的 YAML 会被解析为 config.Config 结构体
	// ================================================================
	loader, err := nacosconfig.NewLoader(bs)
	if err != nil {
		panic(fmt.Sprintf("init nacos config loader: %v", err))
	}
	defer loader.Close()

	c, err := loader.Load()
	if err != nil {
		panic(fmt.Sprintf("load config from nacos: %v", err))
	}
	logx.Infof("[nacos] config loaded: dataId=%s group=%s", bs.ConfigDataId, bs.ConfigGroup)

	// 原子容器持有配置指针，后续热更新时通过 atomic.Value 无锁替换
	atomicCfg := config.NewAtomicConfig(c)

	// ================================================================
	// Step 3: 监听 Nacos 配置变更（热更新）
	// 注册 Watcher 回调，当 Nacos 控制台修改并发布配置后自动触发
	// 新配置会被解析并替换 atomicCfg 中的旧配置，服务无需重启
	// 注意：ListenOn（端口）、Name（服务名）等涉及 gRPC Server 生命周期的字段
	//       变更后仍需重启才能生效，热更新适合 Nacos、Log、AdvertisedIp 等运行时参数
	// ================================================================
	if err := loader.Listen(func(newCfg *config.Config) {
		oldCfg := atomicCfg.Load()
		atomicCfg.Store(newCfg)

		logx.Infof("[nacos] config hot-reloaded")
		newCfg.Print(newCfg.AdvertisedIp)

		if oldCfg.ListenOn != newCfg.ListenOn {
			logx.Errorf("[nacos] ListenOn changed (%s -> %s), restart required to take effect", oldCfg.ListenOn, newCfg.ListenOn)
		}
		if oldCfg.Name != newCfg.Name {
			logx.Errorf("[nacos] Name changed (%s -> %s), restart required to take effect", oldCfg.Name, newCfg.Name)
		}
	}); err != nil {
		logx.Errorf("[nacos] listen config failed: %v", err)
	}

	// ================================================================
	// Step 4: 自动探测公网 IP
	// 用于向 Nacos 注册的 advertise 地址（客户端通过此 IP 连接本服务）
	// 优先级：
	//   1. Nacos 配置中的 AdvertisedIp（手动指定）
	//   2. 环境变量 ADVERTISED_IP
	//   3. 外部 HTTP 服务探测（ifconfig.me / ipify.org 等）
	//   4. 本地网卡出口 IP（兜底）
	// ================================================================
	curCfg := atomicCfg.Load()
	advertiseIP := curCfg.AdvertisedIp
	if advertiseIP == "" {
		detected, err := netutil.DiscoverPublicIP()
		if err != nil {
			logx.Errorf("auto detect public IP failed: %v, falling back to listen address", err)
			_, port, _ := strings.Cut(curCfg.ListenOn, ":")
			advertiseIP = "0.0.0.0" + ":" + port
		} else {
			advertiseIP = detected
			logx.Infof("[auto-ip] detected public IP: %s", advertiseIP)
		}
	}

	// ================================================================
	// Step 5: Print final active configuration
	// ================================================================
	curCfg.Print(advertiseIP)

	// ================================================================
	// Step 6: 构建服务上下文并启动 gRPC Server
	// ServiceContext 持有配置，供 logic 层使用
	// gRPC Server 注册 Greeter 服务实现
	// dev/test 模式下开启 reflection，支持 grpcurl 调试
	// ================================================================
	svcCtx := svc.NewServiceContext(atomicCfg)

	// Nacos 服务发现的连接参数（用于后续注册实例）
	nacosSC := []constant.ServerConfig{
		*constant.NewServerConfig(curCfg.Nacos.Ip, curCfg.Nacos.Port),
	}

	nacosCC := &constant.ClientConfig{
		NamespaceId:         curCfg.Nacos.Namespace,
		Username:            curCfg.Nacos.Username,
		Password:            curCfg.Nacos.Password,
		TimeoutMs:           50000,
		NotLoadCacheAtStart: curCfg.Nacos.NotLoadCacheAtStart,
		LogDir:              "logs/nacos/naming",
		CacheDir:            "cache/nacos/naming",
		LogLevel:            curCfg.Nacos.LogLevel,
	}

	s := zrpc.MustNewServer(curCfg.RpcServerConf, func(grpcServer *grpc.Server) {
		greeter.RegisterGreeterServer(grpcServer, server.NewGreeterServer(svcCtx))
		if curCfg.Mode == service.DevMode || curCfg.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// ================================================================
	// Step 7: 解析注册地址并注册到 Nacos（持久化实例）
	// advertiseIP + ListenOn 端口拼接为最终注册地址
	// 通过 gRPC 协议直接注册持久化实例（绕过 SDK HTTP 代理）
	// 服务关闭时自动注销实例
	// ================================================================
	_, portStr, _ := strings.Cut(curCfg.ListenOn, ":")
	advertiseAddr := advertiseIP + ":" + portStr

	host, portStr, err := net.SplitHostPort(advertiseAddr)
	if err != nil {
		panic(fmt.Sprintf("failed parsing advertise address %s: %v", advertiseAddr, err))
	}
	port, _ := strconv.ParseUint(portStr, 10, 16)

	groupName := curCfg.Nacos.GroupName
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}

	grpcProxy := registerPersistent(curCfg.RpcServerConf.Name, groupName, host, port, nacosSC, *nacosCC)
	defer grpcProxy.CloseClient()

	proc.AddShutdownListener(func() {
		_, err := grpcProxy.DeregisterInstance(curCfg.RpcServerConf.Name, groupName, model.Instance{
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

	fmt.Printf("Starting rpc server at %s (registered as %s)...\n", curCfg.ListenOn, advertiseAddr)
	s.Start()
}

// registerPersistent 通过 SDK gRPC 代理注册持久化实例
// Nacos SDK 默认将 Ephemeral:false 路由到 HTTP API，
// 但部分 Nacos 部署不支持 HTTP naming API（返回 501），
// 此处绕过代理路由，直接使用 gRPC 代理完成持久化注册
func registerPersistent(
	serviceName, groupName, ip string, port uint64,
	nacosSC []constant.ServerConfig,
	nacosCC constant.ClientConfig,
) *naming_grpc.NamingGrpcProxy {

	// 创建服务信息缓存，用于本地缓存已发现的服务实例列表
	serviceInfoHolder := naming_cache.NewServiceInfoHolder(
		nacosCC.NamespaceId,
		nacosCC.CacheDir,
		nacosCC.UpdateCacheWhenEmpty,
		nacosCC.NotLoadCacheAtStart,
	)

	// 创建 Nacos Server 连接（底层 HTTP 客户端，用于获取集群信息等）
	nacosServer, err := nacos_server.NewNacosServer(
		context.Background(),
		nacosSC,
		nacosCC,
		&http_agent.HttpAgent{},
		nacosCC.TimeoutMs,
		"",  // endpoint - 不使用 address-server 模式
		nil, // endpoint 查询头
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating nacos server: %v", err))
	}

	// 创建 gRPC 代理，后续通过此代理完成持久化实例的注册/注销
	grpcProxy, err := naming_grpc.NewNamingGrpcProxy(
		context.Background(),
		nacosCC,
		nacosServer,
		serviceInfoHolder,
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating grpc proxy: %v", err))
	}

	// 注册持久化实例（Ephemeral: false）
	// 持久化实例会在 Nacos 服务端持久存储，即使服务端下线也不会立即移除
	_, err = grpcProxy.RegisterInstance(serviceName, groupName, model.Instance{
		Ip:        ip,
		Port:      port,
		Weight:    1.0,
		Enable:    true,
		Healthy:   true,
		Ephemeral: false,
		Metadata:  map[string]string{},
	})
	if err != nil {
		panic(fmt.Sprintf("failed registering persistent instance: %v", err))
	}

	logx.Infof("registered persistent instance %s:%d to nacos service %s group %s", ip, port, serviceName, groupName)
	return grpcProxy
}
