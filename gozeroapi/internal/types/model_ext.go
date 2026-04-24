// Model extension methods
// 注意：这个文件用于存放自定义的结构体方法，不会被生成工具覆盖

package types

import (
	"database/sql"
	"time"

	"gozeroapi/model/mysql"
)

// User 模型扩展方法

// TableName 指定 User 结构体对应的数据库表名
// 这通常用于 ORM 框架（如 gorm）来确定表名
func (u *User) TableName() string {
	return "users" // 数据库表名
}

// ToDBModel 将 API User 类型转换为数据库 Users 类型
// 这是 types.User (单数) 和 mysql.Users (复数) 之间的适配器
// 巧妙的设计，在User结构体中增加转换为Users的方法
func (u *User) ToDBModel() *mysql.Users {
	dbUser := &mysql.Users{
		Id: int64(u.Id),
	}
	// 处理 Name -> Username 的字段映射
	if u.Username != "" {
		dbUser.Username = sql.NullString{
			String: u.Username,
			Valid:  true,
		}
	}
	// 处理 AddTime 字符串 -> sql.NullTime 的类型转换
	if u.AddTime != "" {
		// 解析字符串格式的时间
		parsedTime, err := time.Parse("2006-01-02 15:04:05", u.AddTime)
		if err == nil {
			dbUser.AddTime = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}
	}
	return dbUser
}

// FromDBModel 从数据库 Users 类型转换为 API User 类型
// 这是反向的适配器，用于查询结果的转换
// 巧妙的设计，在Users结构体中增加转换为User的方法
func UserFromDBModel(dbUser *mysql.Users) *User {
	return &User{
		Id:       int(dbUser.Id),
		Username: dbUser.Username.String,
		AddTime:  dbUser.AddTime.Time.Format("2006-01-02 15:04:05"),
	}
}

// Focus 模型扩展方法

// TableName 指定 Focus 结构体对应的数据库表名
func (f *Focus) TableName() string {
	return "focus" // 数据库表名
}
