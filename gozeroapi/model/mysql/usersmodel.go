package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel
		FindAll(ctx context.Context) ([]*Users, error) // 此处定义新的数据模型方法
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

// FindAll 实现自定义模型方法：查询所有用户
func (m *customUsersModel) FindAll(ctx context.Context) ([]*Users, error) {
	var users []*Users
	query := "select * from `users`"
	err := m.conn.QueryRowsCtx(ctx, &users, query)
	return users, err
}

// Insert 自定义Insert方法：仅插入username，由数据库自动设置id和add_time
func (m *customUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
	// 由于数据库表的 add_time 字段设置了默认值 CURRENT_TIMESTAMP，
	// 所以这里只需要插入 username 字段
	query := fmt.Sprintf("insert into %s (`username`) values (?)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Username)
	return ret, err
}
