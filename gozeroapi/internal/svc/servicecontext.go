// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"gozeroapi/internal/config"
	"gozeroapi/model"
	"gozeroapi/model/mysql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	DataSource string
	UserModel  mysql.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	model.InitDB(c.Mysql.DataSource)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		DataSource: c.Mysql.DataSource, // 为业务访问挂载上
		UserModel:  mysql.NewUsersModel(conn),
	}
}
