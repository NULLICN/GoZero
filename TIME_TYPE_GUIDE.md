# Go Zero 时间类型处理指南

## 问题说明
在 Go Zero 的 API 定义文件 (`.api`) 中，**不能使用 `time.Time` 类型**，因为 API 定义语言只支持基本类型。

## 解决方案：使用 `string` 类型

### 1. API 定义 (users.api)

```api
type User {
    Id       int    `json:"id"`
    Name     string `json:"name"`
    AddTime  string `json:"add_time"`  // ✅ 使用 string
}
```

**为什么用 string？**
- API 通过 JSON 传输，JSON 中时间通常为字符串格式
- Go Zero 只支持基本类型
- 便于前后端时间格式的协商

---

### 2. 生成的类型文件 (types.go)

生成后会自动变成：
```go
type User struct {
    Id      int    `json:"id"`
    Name    string `json:"name"`
    AddTime string `json:"add_time"`
}
```

---

### 3. 数据库模型 (Model 层)

在数据库模型中使用 `sql.NullTime` 来处理 NULL 值：

```go
// usersmodel_gen.go
type Users struct {
    Id       int64          `db:"id"`
    Username sql.NullString `db:"username"`
    AddTime  sql.NullTime   `db:"add_time"`  // ✅ 使用 sql.NullTime
}
```

**为什么用 sql.NullTime？**
- 处理数据库中可能为 NULL 的时间戳字段
- 避免 time.Time 的零值问题
- 标准的数据库时间处理方式

---

### 4. Logic 层处理流程（✨ 使用适配器模式）

#### 新的推荐做法（使用 model_ext.go 中的转换函数）

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
    svcCtx *svc.ServiceContext
}

func (l *AddUserLogic) AddUser(req *types.UserAdd) (resp *types.CommonResponse, err error) {
    l.Infof("开始创建用户，Name: %s", req.Name)

    // 1. 获取当前时间（字符串格式）
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    // 2. 构建 API 层的 User 对象（单数）
    apiUser := &types.User{
        Id:      0,
        Name:    req.Name,
        AddTime: currentTime,
    }

    // 3. ✨ 使用转换函数将 API User 转换为数据库 Users 类型
    // 这是在 internal/types/model_ext.go 中定义的适配器方法
    dbUser := apiUser.ToDBModel()

    // 4. 调用数据库 Model 层的 Insert 方法
    result, err := l.svcCtx.UserModel.Insert(l.ctx, dbUser)
    if err != nil {
        l.Errorf("创建用户失败: %v", err)
        return &types.CommonResponse{
            Success: false,
            Code:    500,
            Message: "创建用户失败",
        }, err
    }

    // 5. 获取新插入的用户ID
    lastInsertId, _ := result.LastInsertId()
    apiUser.Id = int(lastInsertId)

    // 6. 返回 API 层的响应
    l.Infof("创建用户成功，用户ID: %d", apiUser.Id)
    return &types.CommonResponse{
        Success: true,
        Code:    200,
        Message: "创建用户成功",
        Data:    apiUser,
    }, nil
}
```

#### 查询时的反向转换

```go
package users

import (
    "context"
    "strconv"
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

    // 1. 转换 ID 字符串为 int64
    aUserId, parseErr := strconv.ParseInt(req.Id, 10, 64)
    if parseErr != nil {
        l.Errorf("ID格式错误: %s", req.Id)
        return &types.CommonResponse{
            Success: false,
            Code:    400,
            Message: "ID格式错误",
        }, nil
    }

    // 2. 从数据库查询 mysql.Users 类型
    dbUser, queryErr := l.svcCtx.UserModel.FindOne(l.ctx, aUserId)
    if queryErr != nil {
        l.Errorf("查询用户失败，ID: %s", req.Id)
        return &types.CommonResponse{
            Success: false,
            Code:    500,
            Message: "查询用户失败",
        }, nil
    }

    // 3. ✨ 使用反向转换函数将数据库 Users 转换为 API User
    // 这是在 internal/types/model_ext.go 中定义的反向适配器
    apiUser := types.UserFromDBModel(dbUser)

    // 4. 返回 API 层的响应
    l.Infof("查询用户成功，ID: %s", req.Id)
    return &types.CommonResponse{
        Success: true,
        Code:    200,
        Message: "查询用户成功",
        Data:    apiUser,
    }, nil
}
```

#### 适配器函数定义（internal/types/model_ext.go）

```go
package types

import (
    "database/sql"
    "time"
    "gozeroapi/model/mysql"
)

// ToDBModel - 将 API User 类型转换为数据库 Users 类型
// types.User (单数) → mysql.Users (复数)
func (u *User) ToDBModel() *mysql.Users {
    dbUser := &mysql.Users{
        Id: int64(u.Id),
    }
    
    // 字段映射：Name → Username
    if u.Name != "" {
        dbUser.Username = sql.NullString{
            String: u.Name,
            Valid:  true,
        }
    }
    
    // 类型转换：string → sql.NullTime
    if u.AddTime != "" {
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

// UserFromDBModel - 将数据库 Users 类型转换为 API User 类型
// mysql.Users (复数) → types.User (单数)
func UserFromDBModel(dbUser *mysql.Users) *User {
    timeStr := ""
    if dbUser.AddTime.Valid {
        timeStr = dbUser.AddTime.Time.Format("2006-01-02 15:04:05")
    }
    
    return &User{
        Id:      int(dbUser.Id),
        Name:    dbUser.Username.String,
        AddTime: timeStr,
    }
}
```

---

## 时间类型转换对应表

| 层级 | 类型 | 格式示例 | 说明 |
|------|------|---------|------|
| **API 定义** | `string` | 如上所示 | API 文件中的定义 |
| **API 生成** (types.go) | `string` | `"2006-01-02 15:04:05"` | JSON 序列化/反序列化 |
| **Logic 层** | 两者混用 | 见代码示例 | 接收 string，转换为 time.Time |
| **Model 层** | `sql.NullTime` | `{Time: time.Now(), Valid: true}` | 处理数据库 NULL 值 |
| **数据库字段** | `timestamp` | `2024-04-24 10:30:45` | MySQL 字段类型 |
| **前端接收** | JSON 字符串 | `"2024-04-24 10:30:45"` | JSON 格式 |

---

## 关键要点

✅ **DO：** 在 API 文件中使用 `string` 表示时间
✅ **DO：** 在 Model 层使用 `sql.NullTime` 处理数据库时间
✅ **DO：** 在 Logic 层进行时间类型转换
✅ **DO：** 使用 `model_ext.go` 中的适配器函数进行类型转换
✅ **DO：** 按照适配器模式进行分层转换

❌ **DON'T：** 在 API 文件中使用 `time.Time`
❌ **DON'T：** 在 Model 层使用 `time.Time`（会导致 NULL 处理问题）
❌ **DON'T：** 在 API 响应中直接使用 `time.Time` 对象
❌ **DON'T：** 混乱使用 `types.User` 和 `mysql.Users`，应该使用转换函数

---

## 常用时间格式

```go
time.Now().Format("2006-01-02 15:04:05")        // 标准格式：2024-04-24 10:30:45
time.Now().Format("2006-01-02")                  // 仅日期：2024-04-24
time.Now().Format(time.RFC3339)                  // RFC3339：2024-04-24T10:30:45Z
time.Now().Format("20060102150405")              // 紧凑格式：20240424103045
```

---

## 总结

你的 `users.api` 文件现在是 **正确的**：
```api
AddTime  string `json:"add_time"`  ✅
```

这样设置确保了：
1. ✅ API 定义符合 Go Zero 规范
2. ✅ JSON 序列化/反序列化正常工作
3. ✅ 与数据库的时间戳字段完全兼容
4. ✅ 时间类型在各层的正确转换
5. ✅ 通过适配器模式清晰地分离 `types.User` 和 `mysql.Users`

---

## 完整的分层架构总结

```
┌─────────────────────────────────────────┐
│     HTTP JSON 请求/响应                  │
│  "add_time": "2006-01-02 15:04:05"      │
└────────────┬────────────────────────────┘
             │ JSON 反序列化
┌────────────▼────────────────────────────┐
│    types.User.AddTime (string)          │
│    "2006-01-02 15:04:05"               │
└────────────┬────────────────────────────┘
             │ 转换函数：ToDBModel()
┌────────────▼────────────────────────────┐
│   mysql.Users.AddTime (sql.NullTime)    │
│   {Time: 2024-04-24 10:30:45, Valid:t} │
└────────────┬────────────────────────────┘
             │ SQL Insert/Query
┌────────────▼────────────────────────────┐
│  数据库字段 (TIMESTAMP/DATETIME)        │
│  2024-04-24 10:30:45                   │
└─────────────────────────────────────────┘
```

这种分层结构的优势：
- **清晰的职责分工**：每层只处理自己的事情
- **易于维护**：改变数据库类型时，只需修改转换函数
- **不会被覆盖**：转换函数在 `model_ext.go` 中，不会被 goctl 重新生成覆盖
- **类型安全**：每层使用最合适的类型
- **易于测试**：可以独立测试转换逻辑

