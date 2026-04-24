# 如何在 Logic 中使用配置文件

## 📖 配置文件位置和结构

### 1. 配置文件定义：`internal/config/config.go`

```go
package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
    rest.RestConf
    // 下面添加你的自定义配置字段
    DbUserName string `yaml:"db_username"`
    DbPassword string `yaml:"db_password"`
    Timeout    int64  `yaml:"timeout"`  // 毫秒
    MaxRetry   int    `yaml:"max_retry"`
}
```

### 2. YAML 配置文件：`etc/gozero-api.yaml`

```yaml
Name: gozero-api
Host: 0.0.0.0
Port: 8888

# 自定义配置字段
db_username: root
db_password: password123
timeout: 5000
max_retry: 3
```

### 3. 配置文件读取流程

```
gozero.go (main函数)
    ↓
读取 etc/gozero-api.yaml
    ↓
解析为 config.Config 结构体
    ↓
传入 ServiceContext
    ↓
在 Logic 中通过 l.svcCtx.Config 访问
```

## 🔧 在 Logic 中使用配置

### 示例 1：在 AddUserLogic 中使用配置

```go
package users

import (
    "context"
    "time"

    "gozeroapi/internal/svc"
    "gozeroapi/internal/types"

    "github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext  // ✨ 包含配置
}

func (l *AddUserLogic) AddUser(req *types.UserAdd) (resp *types.CommonResponse, err error) {
    l.Infof("开始创建用户，Name: %s", req.Name)

    // ✨ 获取配置中的超时时间
    timeout := time.Duration(l.svcCtx.Config.Timeout) * time.Millisecond
    
    // ✨ 创建带超时的 Context
    ctx, cancel := context.WithTimeout(l.ctx, timeout)
    defer cancel()

    // 获取当前时间
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    // 创建用户
    apiUser := &types.User{
        Id:      0,
        Name:    req.Name,
        AddTime: currentTime,
    }

    dbUser := apiUser.ToDBModel()

    // ✨ 使用带超时的 Context 执行数据库操作
    result, err := l.svcCtx.UserModel.Insert(ctx, dbUser)
    if err != nil {
        l.Errorf("创建用户失败: %v", err)
        
        // ✨ 根据配置的重试次数重试
        for i := 0; i < l.svcCtx.Config.MaxRetry; i++ {
            l.Infof("重试第 %d 次", i+1)
            result, err = l.svcCtx.UserModel.Insert(ctx, dbUser)
            if err == nil {
                break
            }
        }
        
        if err != nil {
            resp = &types.CommonResponse{
                Success: false,
                Code:    500,
                Message: "创建用户失败",
            }
            return
        }
    }

    lastInsertId, _ := result.LastInsertId()
    apiUser.Id = int(lastInsertId)

    resp = &types.CommonResponse{
        Success: true,
        Code:    200,
        Message: "创建用户成功",
        Data:    apiUser,
    }

    return
}
```

### 示例 2：在 GetUsersByIdLogic 中使用配置

```go
package users

import (
    "context"
    "strconv"
    "time"

    "gozeroapi/internal/svc"
    "gozeroapi/internal/types"

    "github.com/zeromicro/go-zero/core/logx"
)

type GetUsersByIdLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (l *GetUsersByIdLogic) GetUsersById(req *types.UserQuestById) (resp *types.CommonResponse, err error) {
    l.Infof("开始查询用户，ID: %s", req.Id)

    // ✨ 使用配置的超时时间
    timeout := time.Duration(l.svcCtx.Config.Timeout) * time.Millisecond
    ctx, cancel := context.WithTimeout(l.ctx, timeout)
    defer cancel()

    // 转换 ID
    aUserId, parseErr := strconv.ParseInt(req.Id, 10, 64)
    if parseErr != nil {
        l.Errorf("ID格式错误: %s", req.Id)
        resp = &types.CommonResponse{
            Success: false,
            Code:    400,
            Message: "ID格式错误",
        }
        return
    }

    // 查询用户
    dbUser, queryErr := l.svcCtx.UserModel.FindOne(ctx, aUserId)
    if queryErr != nil {
        l.Errorf("查询用户失败，ID: %s, 错误: %v", req.Id, queryErr)
        resp = &types.CommonResponse{
            Success: false,
            Code:    500,
            Message: "查询用户失败",
        }
        return
    }

    // 转换为 API 类型
    apiUser := types.UserFromDBModel(dbUser)

    resp = &types.CommonResponse{
        Success: true,
        Code:    200,
        Message: "查询用户成功",
        Data:    apiUser,
    }
    return
}
```

## 📝 如何添加新的配置字段

### 第 1 步：修改 `internal/config/config.go`

```go
type Config struct {
    rest.RestConf
    
    // 原有字段
    DbUserName string `yaml:"db_username"`
    DbPassword string `yaml:"db_password"`
    Timeout    int64  `yaml:"timeout"`
    MaxRetry   int    `yaml:"max_retry"`
    
    // 新加字段
    LogLevel   string `yaml:"log_level"`        // 日志级别
    CacheTTL   int64  `yaml:"cache_ttl"`        // 缓存过期时间（秒）
    EnableMetrics bool `yaml:"enable_metrics"`  // 是否启用指标收集
}
```

### 第 2 步：修改 `etc/gozero-api.yaml`

```yaml
Name: gozero-api
Host: 0.0.0.0
Port: 8888

# 数据库配置
db_username: root
db_password: password123

# 超时和重试配置
timeout: 5000
max_retry: 3

# 新加配置
log_level: info
cache_ttl: 3600
enable_metrics: true
```

### 第 3 步：在 Logic 中使用

```go
func (l *AddUserLogic) AddUser(req *types.UserAdd) (resp *types.CommonResponse, err error) {
    // 访问新的配置字段
    l.Infof("日志级别: %s", l.svcCtx.Config.LogLevel)
    l.Infof("缓存过期时间: %d 秒", l.svcCtx.Config.CacheTTL)
    
    if l.svcCtx.Config.EnableMetrics {
        // 记录指标
    }
    
    // ... 业务逻辑
}
```

## 🎯 配置最佳实践

### 1. 配置分类

```go
// 根据功能分类组织配置
type Config struct {
    rest.RestConf
    
    // 数据库配置
    Database struct {
        Host     string
        Port     int
        Username string
        Password string
    }
    
    // 缓存配置
    Cache struct {
        TTL    int64
        Enable bool
    }
    
    // 超时和重试
    TimeOut struct {
        Request int64
        Database int64
    }
    Retry int
}
```

对应的 YAML：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: password123

cache:
  ttl: 3600
  enable: true

timeout:
  request: 5000
  database: 3000
retry: 3
```

### 2. 配置验证

```go
// 在初始化时验证配置
func (c *Config) Validate() error {
    if c.Timeout <= 0 {
        return fmt.Errorf("timeout must be positive")
    }
    if c.MaxRetry < 0 {
        return fmt.Errorf("max_retry cannot be negative")
    }
    return nil
}
```

### 3. 日志记录配置

```go
// 在 Logic 中记录配置信息（便于调试）
func (l *AddUserLogic) AddUser(req *types.UserAdd) (resp *types.CommonResponse, err error) {
    l.Debugf("配置信息 - 超时: %dms, 重试: %d次", 
        l.svcCtx.Config.Timeout, 
        l.svcCtx.Config.MaxRetry)
    
    // ... 业务逻辑
}
```

## 🔍 查看完整的配置文件

### 当前项目的配置文件

- **定义**：`internal/config/config.go`
- **YAML**：`etc/gozero-api.yaml`
- **读取**：`gozero.go` 中的 `main()` 函数
- **传递**：通过 `svc.ServiceContext` 传递到 Logic

### 查看当前配置

```go
// 在 main.go 或任意 logic 中
fmt.Printf("当前配置: %+v\n", l.svcCtx.Config)
```

## 📚 参考链接

- [Go Zero 官方文档 - 配置管理](https://go-zero.dev/)
- [YAML 语法指南](https://yaml.org/)

