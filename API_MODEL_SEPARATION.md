# Go Zero API 和 Model 分层设计指南

## 📋 问题回顾

在使用 Go Zero + sqlx 时遇到的命名冲突问题：

```
问题：表名 users (复数) vs 业务实体 User (单数)

层级结构如下：
┌──────────────────────────────────────────────────┐
│           HTTP API 请求/响应                      │
│    types.User (单数) - 用于 JSON 序列化           │
└──────────────────┬───────────────────────────────┘
                   │ 需要转换
┌──────────────────▼───────────────────────────────┐
│      Business Logic (AddUserLogic等)              │
│      转换层：types.User ←→ mysql.Users           │
└──────────────────┬───────────────────────────────┘
                   │ 需要转换
┌──────────────────▼───────────────────────────────┐
│      Database Model Layer (sqlx 生成)             │
│  mysql.Users (复数) - 对应数据库表 users          │
└──────────────────────────────────────────────────┘
```

## ✅ 解决方案：适配器模式（Adapter Pattern）

### 原理

利用 Go Zero 的特性：**自定义的 `model_ext.go` 文件不会被 goctl 覆盖**

在这个文件中添加转换函数作为 API 类型和 Database 类型之间的**适配器**。

### 关键文件

#### 1. `internal/types/model_ext.go` - 转换层

```go
// 正向转换：API → Database
func (u *User) ToDBModel() *mysql.Users {
    // types.User → mysql.Users 的转换逻辑
    // 处理字段映射：Name → Username
    // 处理类型转换：string → sql.NullTime
}

// 反向转换：Database → API
func UserFromDBModel(dbUser *mysql.Users) *User {
    // mysql.Users → types.User 的转换逻辑
    // 回复字段映射和类型转换
}
```

#### 2. `internal/logic/users/adduserlogic.go` - 创建流程

```
请求报文 (JSON)
    ↓
types.UserAdd (API 请求类型)
    ↓ [✨ 转换]
types.User (API 业务类型)
    ↓ [✨ 调用 ToDBModel()]
mysql.Users (数据库类型)
    ↓ [数据库操作]
数据库
```

**代码示例：**
```go
func (l *AddUserLogic) AddUser(req *types.UserAdd) (resp *types.CommonResponse, err error) {
    // 1. 获取当前时间
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    
    // 2. 创建 API 类型对象
    apiUser := &types.User{
        Id:      0,
        Name:    req.Name,
        AddTime: currentTime,
    }
    
    // 3. 转换为数据库类型（✨ 使用适配器）
    dbUser := apiUser.ToDBModel()
    
    // 4. 调用数据库操作
    result, err := l.svcCtx.UserModel.Insert(l.ctx, dbUser)
    
    // 返回 API 类型的响应
    return
}
```

#### 3. `internal/logic/users/getusersbyidlogic.go` - 查询流程

```
HTTP 请求  ?id=123
    ↓
types.UserQuestById (API 请求类型)
    ↓ [数据库查询]
mysql.Users (数据库类型)
    ↓ [✨ 调用 UserFromDBModel()]
types.User (API 业务类型)
    ↓
JSON 响应报文
```

**代码示例：**
```go
func (l *GetUsersByIdLogic) GetUsersById(req *types.UserQuestById) (resp *types.CommonResponse, err error) {
    // 1. 从数据库查询
    dbUser, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
    
    // 2. 转换为 API 类型（✨ 使用适配器）
    apiUser := types.UserFromDBModel(dbUser)
    
    // 3. 返回 API 类型的响应
    return
}
```

## 📊 字段映射表

| API Layer (types.User) | ↔ | Database Layer (mysql.Users) | 说明 |
|------------------------|---|------------------------------|------|
| `Id: int` | ↔ | `Id: int64` | ID字段转换 |
| `Name: string` | ↔ | `Username: sql.NullString` | 字段名和类型转换 |
| `AddTime: string` | ↔ | `AddTime: sql.NullTime` | 字符串时间 ↔ sql时间类型 |

## 🔑 核心要点

### ✅ 为什么要分层？

1. **职责单一**：每层只负责自己的事情
   - API层：请求/响应格式、HTTP序列化
   - Logic层：业务逻辑、数据转换
   - Model层：数据库操作、SQL执行

2. **代码重用**：Model层可被多个不同的API使用

3. **易于测试**：每层可独立测试

4. **不会被覆盖**：自定义的转换函数在 `model_ext.go` 中，不会被 goctl 重新生成时覆盖

### ✅ 为什么不直接使用 mysql.Users？

❌ **不建议**：在 API 响应中直接使用 `mysql.Users`，因为：

```go
// ❌ 错误做法
type User struct {
    Username sql.NullString  // SQL 类型暴露给 API
    AddTime  sql.NullTime    // 需要特殊处理
}
// 导致前端接收到的 JSON 格式奇怪

// ✅ 正确做法
type User struct {
    Name    string  // 简洁的字段名
    AddTime string  // 标准的 JSON 时间格式
}
```

## 🚀 最佳实践总结

| 操作 | 使用类型 | 说明 |
|------|---------|------|
| API 请求参数 | `types.UserAdd` | Go Zero 自动生成 |
| 业务逻辑中创建 | `types.User` | API 类型 |
| 转换为数据库 | `apiUser.ToDBModel()` | 使用适配器转换 |
| 数据库操作 | `mysql.Users` | 由 goctl 生成 |
| 数据库查询结果 | `UserFromDBModel(dbUser)` | 使用反向适配器 |
| API 响应类型 | `types.User` | API 类型 |

## 🔧 时间类型处理流程

```
前端请求：2006-01-02 15:04:05
    ↓
types.User.AddTime (string)
    ↓ [转换]
time.Parse("2006-01-02 15:04:05", str)
    ↓
sql.NullTime {Time: time.Time, Valid: true}
    ↓ [数据库]
数据库字段：TIMESTAMP

查询返回：
数据库 TIMESTAMP
    ↓
sql.NullTime {Time: time.Time, Valid: true}
    ↓ [转换]
parsedTime.Format("2006-01-02 15:04:05")
    ↓
types.User.AddTime (string)
    ↓
JSON 响应
```

## 📚 在应用中使用配置文件

在 `svc/servicecontext.go` 中，你可以访问配置：

```go
type ServiceContext struct {
    Config    config.Config     // ✨ 配置对象
    UserModel mysql.UsersModel  // Model
}

// 在 logic 中
func (l *AddUserLogic) AddUser(req *types.UserAdd) {
    // 访问配置
    timeout := l.svcCtx.Config.SomeTimeout  // 从 yaml 读取
    maxRetry := l.svcCtx.Config.MaxRetry
    
    // 根据配置执行逻辑
}
```

配置来自 `etc/gozero-api.yaml` 并在 `internal/config/config.go` 中定义。

## 📖 相关文件参考

- `internal/types/model_ext.go` - **转换函数定义**
- `internal/logic/users/adduserlogic.go` - **创建流程示例**
- `internal/logic/users/getusersbyidlogic.go` - **查询流程示例**
- `model/mysql/usersmodel_gen.go` - goctl 生成的数据库模型（勿编辑）
- `model/mysql/usersmodel.go` - 自定义的模型扩展（可编辑）

