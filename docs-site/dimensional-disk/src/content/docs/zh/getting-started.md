---
title: 快速开始
description: 学习如何连接数据库、定义模型和执行基本操作
---

# 快速开始

本指南将帮助您快速上手 Pie，学习如何连接数据库、定义模型和执行基本操作。

## 安装

```bash
go get github.com/5xxxx/pie
```

## 1. 连接数据库

首先，创建一个引擎实例来连接 MongoDB 数据库：

```go
package main

import (
    "context"
    "log"
    
    "github.com/5xxxx/pie"
)

func main() {
    // 创建引擎
    engine, err := pie.NewEngine(
        context.Background(),
        "mydb",
        pie.WithURI("mongodb://localhost:27017"),
        pie.WithMapper(&pie.SnakeMapper{}),
    )
    if err != nil {
        log.Fatal("Failed to create engine:", err)
    }
    defer engine.Disconnect(context.Background())
    
    // 使用引擎...
}
```

### 配置选项

Pie 提供了丰富的配置选项来满足不同需求：

```go
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithURI("mongodb://localhost:27017"),
    pie.WithAuth("username", "password"),
    pie.WithSSL(true),
    pie.WithReplicaSet("rs0"),
    pie.WithReadPreference("secondary"),
    pie.WithWriteConcern("majority"),
    pie.WithReadConcern("majority"),
    pie.WithMaxPoolSize(100),
    pie.WithMinPoolSize(5),
    pie.WithMaxIdleTime(30*time.Minute),
    pie.WithConnectTimeout(10*time.Second),
    pie.WithSocketTimeout(30*time.Second),
    pie.WithServerSelectionTimeout(5*time.Second),
)
```

## 2. 定义模型

定义您的数据结构，Pie 支持完整的生命周期钩子：

```go
import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Age       int                `bson:"age"`
    Status    string             `bson:"status"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// 钩子方法
func (u *User) BeforeCreate(ctx context.Context) error {
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return nil
}

func (u *User) AfterCreate(ctx context.Context) error {
    log.Printf("User %s created", u.Name)
    return nil
}

func (u *User) BeforeUpdate(ctx context.Context) error {
    u.UpdatedAt = time.Now()
    return nil
}

func (u *User) AfterUpdate(ctx context.Context) error {
    log.Printf("User %s updated", u.Name)
    return nil
}
```

### 模型标签

Pie 支持多种模型标签来定义索引和约束：

```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name" pie:"index"`
    Email     string             `bson:"email" pie:"unique"`
    Age       int                `bson:"age"`
    Status    string             `bson:"status" pie:"index"`
    CreatedAt time.Time          `bson:"created_at" pie:"index"`
    UpdatedAt time.Time          `bson:"updated_at"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" pie:"soft_delete"`
}
```

## 3. 基本操作

### 创建类型安全的会话

```go
// 创建类型安全的会话
session := pie.Table[User](engine)
```

### 插入文档

```go
// 插入单个文档
user := &User{
    Name:   "张三",
    Email:  "zhangsan@example.com",
    Age:    25,
    Status: "active",
}
result, err := session.Insert(context.Background(), user)
if err != nil {
    log.Fatal("Failed to insert user:", err)
}
log.Printf("Inserted user with ID: %s", result.InsertedID)

// 插入多个文档
users := []*User{
    {Name: "李四", Email: "lisi@example.com", Age: 30, Status: "active"},
    {Name: "王五", Email: "wangwu@example.com", Age: 28, Status: "pending"},
}
results, err := session.InsertMany(context.Background(), users)
```

### 查询文档

```go
// 查询所有文档
var users []User
err := session.Find(context.Background(), &users)

// 条件查询
var activeUsers []User
err = session.
    Where("status", "active").
    Where("age", pie.Gte("age", 18)).
    OrderBy("created_at").
    Limit(10).
    Find(context.Background(), &activeUsers)

// 查询单个文档
var user User
err = session.Where("email", "zhangsan@example.com").First(context.Background(), &user)

// 按 ID 查询
user, err := session.FindByID(context.Background(), userID)
```

### 更新文档

```go
// 更新单个文档
updateResult, err := session.
    Where("email", "zhangsan@example.com").
    Update(context.Background(), bson.D{{"$set", bson.D{{"age", 26}}}})

// 更新多个文档
updateResult, err = session.
    Where("status", "pending").
    UpdateMany(context.Background(), bson.D{{"$set", bson.D{{"status", "active"}}}})

// 使用 Upsert
upsertResult, err := session.
    Where("email", "newuser@example.com").
    Upsert(context.Background(), &User{
        Name:   "新用户",
        Email:  "newuser@example.com",
        Status: "active",
    })
```

### 删除文档

```go
// 删除单个文档
deleteResult, err := session.
    Where("email", "zhangsan@example.com").
    Delete(context.Background())

// 删除多个文档
deleteResult, err = session.
    Where("status", "inactive").
    DeleteMany(context.Background())

// 软删除（如果模型支持）
err = session.Where("email", "test@example.com").SoftDelete(context.Background())
```

## 4. 查询构建器

Pie 提供了丰富的查询方法：

### 基础查询

```go
// 等值查询
session.Where("status", "active")

// 比较查询
session.Where("age", pie.Gte("age", 18))
session.Where("age", pie.Lt("age", 65))
session.Where("age", pie.Between("age", 18, 65))

// 包含查询
session.WhereIn("role", []string{"admin", "user"})
session.WhereNotIn("status", []string{"deleted", "banned"})

// 空值查询
session.WhereNull("deleted_at")
session.WhereNotNull("email")
```

### 模糊查询

```go
// 模糊匹配
session.WhereLike("name", "%张%")
session.WhereStartsWith("email", "admin")
session.WhereEndsWith("domain", ".com")

// 正则表达式
session.WhereRegex("name", "^张")
```

### 日期查询

```go
// 日期范围查询
session.WhereRecentDays("created_at", 7)
session.WhereMonth("created_at", time.Now().Month())
session.WhereYear("created_at", 2024)
session.WhereBetween("created_at", startDate, endDate)
```

### 复杂条件

```go
// AND 条件
session.Where("status", "active").
    Where("age", pie.Gte("age", 18))

// OR 条件
session.Where("status", "active").
    OrWhere(func(q *pie.Query) {
        q.Where("role", "admin").WhereBetween("age", 30, 50)
    })

// 嵌套条件
session.Where(func(q *pie.Query) {
    q.Where("status", "active").
        OrWhere("role", "admin")
}).Where("age", pie.Gte("age", 18))
```

## 5. 便捷方法

Pie 提供了许多便捷方法来简化常见操作：

```go
// 检查存在性
exists, err := session.Where("email", "test@example.com").Exists(ctx)

// 快速计数
count, err := session.Where("status", "active").QuickCount(ctx)

// 查找或创建
user, isNew, err := session.
    Where("email", "test@example.com").
    FirstOrCreate(ctx, &User{Email: "test@example.com", Name: "测试用户"})

// 查找或失败
user, err := session.
    Where("email", "test@example.com").
    FirstOrFail(ctx)

// 获取最新记录
var latestUsers []User
err := session.Latest("created_at", 10).Find(ctx, &latestUsers)
```

## 下一步

现在您已经了解了 Pie 的基本用法，可以继续学习：

- [查询构建器](/core-features/query-builder/) - 深入了解查询功能
- [聚合查询](/core-features/aggregation/) - 学习复杂的数据聚合
- [事务管理](/core-features/transactions/) - 掌握事务操作
- [缓存支持](/advanced/cache/) - 提升查询性能
