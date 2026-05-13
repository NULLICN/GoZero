package config

import "github.com/zeromicro/go-zero/core/logx"

// Print logs the active config key-value pairs to the log.
// Keys match the YAML config file fields; empty values print as empty strings.
func (c *Config) Print(advIP string) {
	logx.Info("========================================")
	logx.Info("[config] active configuration")
	logx.Info("========================================")

	logx.Infof("[config] Name                      : %s", c.Name)
	logx.Infof("[config] ListenOn                  : %s", c.ListenOn)
	logx.Infof("[config] AdvertisedIp              : %s", advIP)

	logx.Infof("[config] Nacos.Ip                  : %s", c.Nacos.Ip)
	logx.Infof("[config] Nacos.Port                : %d", c.Nacos.Port)
	logx.Infof("[config] Nacos.Namespace           : %s", c.Nacos.Namespace)
	logx.Infof("[config] Nacos.GroupName           : %s", c.Nacos.GroupName)
	logx.Infof("[config] Nacos.NotLoadCacheAtStart : %v", c.Nacos.NotLoadCacheAtStart)
	logx.Infof("[config] Nacos.LogLevel            : %s", c.Nacos.LogLevel)

	logx.Info("========================================")
}
