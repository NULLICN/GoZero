package config

import (
	"sync/atomic"

	"github.com/zeromicro/go-zero/zrpc"
)

// Config is the main business config, stored in Nacos config center.
type Config struct {
	zrpc.RpcServerConf
	AdvertisedIp string
	Nacos        NacosConfig
}

type NacosConfig struct {
	Ip                  string
	Port                uint64
	Namespace           string
	GroupName           string
	Username            string `json:",optional"`
	Password            string `json:",optional"`
	NotLoadCacheAtStart bool
	LogLevel            string
}

// Bootstrap is the minimal local config, only used to connect Nacos config center.
type Bootstrap struct {
	NacosIp          string
	NacosPort        uint64
	NacosNamespaceId string
	NacosUsername    string `json:",optional"`
	NacosPassword    string `json:",optional"`
	ConfigDataId     string
	ConfigGroup      string
}

// AtomicConfig holds a Config pointer that can be hot-reloaded safely.
type AtomicConfig struct {
	v atomic.Value
}

func NewAtomicConfig(c *Config) *AtomicConfig {
	ac := &AtomicConfig{}
	ac.v.Store(c)
	return ac
}

func (ac *AtomicConfig) Load() *Config {
	return ac.v.Load().(*Config)
}

func (ac *AtomicConfig) Store(c *Config) {
	ac.v.Store(c)
}
