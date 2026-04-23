// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	// Mysql为自定义的配置属性名，取值为对应的yaml中
	Mysql struct {
		DataSource string
	}
}
