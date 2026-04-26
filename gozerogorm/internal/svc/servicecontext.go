// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"gozerogorm/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Error(err)
		return nil
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
