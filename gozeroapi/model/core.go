package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
// dsn: 数据库连接字符串 (Data Source Name)
// 示例: "root:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
func InitDB(dsn string) error {
	// 配置日志记录器
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢查询阈值，超过这个时间会标记为 SLOW SQL
			LogLevel:                  logger.Info,            // 设置为 Info 级别，打印所有 SQL
			IgnoreRecordNotFoundError: true,                   // 忽略 ErrRecordNotFound 错误
			Colorful:                  true,                   // 彩色输出
		},
	)

	fmt.Println("数据源:", dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, // 应用日志配置
	})
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return err
	}

	fmt.Println("数据库连接成功!")
	return nil
}
