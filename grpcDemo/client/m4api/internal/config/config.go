package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	RpcClient zrpc.RpcClientConf
	Nacos     NacosConfig
}

type NacosConfig struct {
	Ip                  string
	Port                uint64
	Namespace           string
	NotLoadCacheAtStart bool
	LogLevel            string
}
