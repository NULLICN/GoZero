package nacosconfig

import (
	"fmt"
	"strings"

	"greeter/internal/config"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type Loader struct {
	client    config_client.IConfigClient
	dataId    string
	group     string
	namespace string
}

func NewLoader(bs config.Bootstrap) (*Loader, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(bs.NacosIp, bs.NacosPort),
	}

	cc := &constant.ClientConfig{
		NamespaceId: bs.NacosNamespaceId,
		Username:    bs.NacosUsername,
		Password:    bs.NacosPassword,
		TimeoutMs:   10000,
		LogDir:      "logs/nacos/config",
		CacheDir:    "cache/nacos/config",
		LogLevel:    "warn",
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  cc,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("nacos config client: %w", err)
	}

	return &Loader{
		client:    client,
		dataId:    bs.ConfigDataId,
		group:     bs.ConfigGroup,
		namespace: bs.NacosNamespaceId,
	}, nil
}

func (l *Loader) Load() (*config.Config, error) {
	logx.Infof("[nacos] fetching config: namespace=%s group=%s dataId=%s", l.namespace, l.group, l.dataId)

	content, err := l.client.GetConfig(vo.ConfigParam{
		DataId: l.dataId,
		Group:  l.group,
	})
	if err != nil {
		return nil, fmt.Errorf("get config from nacos: %w", err)
	}

	content = strings.TrimSpace(content)
	if content == "" || content == "\"\"" {
		return nil, fmt.Errorf(
			"config not found or empty in nacos (namespace=%s, group=%s, dataId=%s). "+
				"please create this config in nacos console first",
			l.namespace, l.group, l.dataId,
		)
	}

	var c config.Config
	if err := conf.LoadFromYamlBytes([]byte(content), &c); err != nil {
		return nil, fmt.Errorf("parse config from nacos (namespace=%s, group=%s, dataId=%s): %w", l.namespace, l.group, l.dataId, err)
	}
	return &c, nil
}

func (l *Loader) Listen(onChange func(newConfig *config.Config)) error {
	return l.client.ListenConfig(vo.ConfigParam{
		DataId: l.dataId,
		Group:  l.group,
		OnChange: func(namespace, group, dataId, data string) {
			logx.Infof("[nacos] config changed: ns=%s group=%s dataId=%s", namespace, group, dataId)
			data = strings.TrimSpace(data)
			if data == "" || data == "\"\"" {
				logx.Errorf("[nacos] new config is empty, ignoring")
				return
			}
			var c config.Config
			if err := conf.LoadFromYamlBytes([]byte(data), &c); err != nil {
				logx.Errorf("[nacos] parse new config failed: %v", err)
				return
			}
			onChange(&c)
		},
	})
}

func (l *Loader) Close() {
	l.client.CloseClient()
}
