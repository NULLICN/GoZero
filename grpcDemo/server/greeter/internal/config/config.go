package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Nacos NacosConfig
}

type NacosConfig struct {
    Ip                  string
    Port                uint64
    Namespace           string
    NotLoadCacheAtStart bool
    LogLevel            string // 通常为 "info" 或 "debug"
}
