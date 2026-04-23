// Model extension methods
// 注意：这个文件用于存放自定义的结构体方法，不会被生成工具覆盖

package types

// User 模型扩展方法

// TableName 指定 User 结构体对应的数据库表名
// 这通常用于 ORM 框架（如 gorm）来确定表名
func (u *User) TableName() string {
	return "users" // 数据库表名
}

// Focus 模型扩展方法

// TableName 指定 Focus 结构体对应的数据库表名
func (f *Focus) TableName() string {
	return "focus" // 数据库表名
}
