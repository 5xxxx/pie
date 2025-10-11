---
title: 泛型使用指南
description: 深入了解 Pie 中泛型的使用方式、优势和最佳实践
---

# 泛型使用指南

Pie 基于 Go 泛型构建，提供类型安全的 MongoDB 操作体验。本指南将详细介绍泛型在 Pie 中的使用方式、优势和最佳实践。

## 什么是泛型？

泛型是 Go 1.18 引入的特性，允许编写可重用的代码，同时保持类型安全。在 Pie 中，泛型确保了编译时的类型检查，避免了运行时的类型错误。

## 基础用法

### Session[T] - 最常用的泛型类型

`Session[T]` 是 Pie 中最核心的泛型类型，提供了类型安全的数据库操作：

```go
// 创建类型安全的会话
session := pie.Table[User](engine)

// 查询操作返回 []User，而不是 []interface{}
users, err := session.Find(context.Background())

// 插入操作接受 *User 类型
user := &User{Name: "张三", Email: "zhangsan@example.com"}
result, err := session.Insert(context.Background(), user)

// 更新操作的类型安全
updateResult, err := session.
    Where("email", "zhangsan@example.com").
    Update(context.Background(), bson.D{{"$set", bson.D{{"age", 25}}}})
```

### Aggregate[T] - 聚合查询的泛型支持

```go
// 创建类型安全的聚合构建器
aggregate := pie.NewAggregate[User](engine)

// 聚合结果返回 []User
var results []User
err := aggregate.
    MatchStage().Where("status", "active").Done().
    GroupStage().By("department", "$department").Count("total").Done().
    Execute(context.Background(), &results)
```

### Cursor[T] - 游标的泛型支持

```go
// 创建游标
cursor := pie.NewCursor[User](context.Background(), mongoCursor)

// 迭代时获得类型安全的 User 对象
for cursor.Next(context.Background()) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        // 处理错误
    }
    // user 是 User 类型，不是 interface{}
}
```

### 泛型函数

Pie 提供了多个泛型函数来创建类型安全的操作：

```go
// 基础函数
session := pie.Table[User](engine)                    // 创建会话
aggregate := pie.NewAggregate[User](engine)           // 创建聚合

// Must 版本（失败时 panic）
session := pie.MustTable[User](engine)                // 必须成功创建会话
aggregate := pie.MustAggregate[User](engine)          // 必须成功创建聚合

// 使用默认引擎
session, err := pie.TableWithDefault[User]()          // 使用默认引擎创建会话
session := pie.MustTableWithDefault[User]()           // 必须成功创建会话
```

## 进阶内容

### 类型安全的好处

泛型提供了编译时类型检查，避免了运行时的类型错误：

```go
// ❌ 传统方式 - 运行时可能出错
var users []interface{}
err := session.Find(context.Background(), &users)
// 需要手动类型断言
for _, u := range users {
    user := u.(User) // 可能 panic
}

// ✅ 泛型方式 - 编译时类型安全
var users []User
err := session.Find(context.Background(), &users)
// 直接使用，无需类型断言
for _, user := range users {
    fmt.Println(user.Name) // 类型安全
}
```

### 类型推断和显式类型声明

Go 的类型推断让泛型使用更加简洁：

```go
// 类型推断
session := pie.Table[User](engine)  // T 被推断为 User

// 显式类型声明（通常不需要）
var session *pie.Session[User] = pie.Table[User](engine)
```

### 泛型与接口的配合使用

```go
// 定义通用接口
type Model interface {
    GetID() string
    GetCreatedAt() time.Time
}

// 使用约束
func ProcessModel[T Model](session *pie.Session[T], ctx context.Context) error {
    var models []T
    err := session.Find(ctx, &models)
    if err != nil {
        return err
    }
    
    for _, model := range models {
        fmt.Printf("ID: %s, Created: %v\n", model.GetID(), model.GetCreatedAt())
    }
    return nil
}
```

### 泛型返回类型

Pie 提供了多种泛型返回类型：

```go
// Result[T] - 通用结果类型
func GetUser[T any](session *pie.Session[T], id string) *pie.Result[T] {
    var user T
    err := session.FindByID(context.Background(), id, &user)
    return &pie.Result[T]{
        Data:  user,
        Error: err,
    }
}

// PaginationResult[T] - 分页结果
func GetUsersWithPagination[T any](session *pie.Session[T], page, size int) *pie.PaginationResult[T] {
    var users []T
    total, err := session.Paginate(context.Background(), page, size, &users)
    return &pie.PaginationResult[T]{
        Data:  users,
        Total: total,
        Page:  int64(page),
        Size:  int64(size),
        Error: err,
    }
}

// AggregateResult[T] - 聚合结果
func GetUserStats[T any](aggregate *pie.Aggregate[T]) *pie.AggregateResult[T] {
    var results []T
    err := aggregate.Execute(context.Background(), &results)
    return &pie.AggregateResult[T]{
        Data:  results,
        Error: err,
    }
}
```

### 泛型在事务中的使用

```go
// TransactionWithResult[T] - 带结果的事务
func TransferMoney[T any](tx *pie.Transaction, fromID, toID string, amount float64) *pie.TransactionResult[T] {
    return pie.TransactionWithResult[T](tx, context.Background(), func(ctx context.Context) (T, error) {
        // 在事务中执行操作
        fromSession := pie.Table[Account](tx.engine)
        toSession := pie.Table[Account](tx.engine)
        
        // 扣除发送方余额
        err := fromSession.Where("id", fromID).Update(ctx, bson.D{{"$inc", bson.D{{"balance", -amount}}}})
        if err != nil {
            return zero, err
        }
        
        // 增加接收方余额
        err = toSession.Where("id", toID).Update(ctx, bson.D{{"$inc", bson.D{{"balance", amount}}}})
        if err != nil {
            return zero, err
        }
        
        return zero, nil
    })
}
```

### 泛型在 Change Streams 中的使用

```go
// ChangeStreamWatcher[T] - 变更流监听器
func WatchUserChanges[T any](engine *pie.Engine) {
    watcher := pie.NewWatcher[T](engine)
    
    watcher.OnChange(func(event *pie.ChangeEvent[T]) {
        switch event.OperationType {
        case "insert":
            fmt.Printf("新用户创建: %+v\n", event.FullDocument)
        case "update":
            fmt.Printf("用户更新: %+v\n", event.FullDocument)
        case "delete":
            fmt.Printf("用户删除: %+v\n", event.DocumentKey)
        }
    })
    
    watcher.Start(context.Background())
}
```

## 最佳实践

### 何时使用泛型 vs 何时使用 interface{}

```go
// ✅ 推荐：使用泛型，类型安全
func ProcessUsers[T any](session *pie.Session[T], ctx context.Context) ([]T, error) {
    var users []T
    err := session.Find(ctx, &users)
    return users, err
}

// ❌ 不推荐：使用 interface{}，失去类型安全
func ProcessUsersGeneric(session *pie.Session[interface{}], ctx context.Context) ([]interface{}, error) {
    var users []interface{}
    err := session.Find(ctx, &users)
    return users, err
}
```

### 模型定义的最佳实践

```go
// ✅ 推荐：定义清晰的模型
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Age       int                `bson:"age"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// ✅ 推荐：实现接口方法
func (u *User) GetID() string {
    return u.ID.Hex()
}

func (u *User) GetCreatedAt() time.Time {
    return u.CreatedAt
}
```

### 避免类型转换的技巧

```go
// ✅ 推荐：使用泛型避免类型转换
func GetUserByEmail[T any](session *pie.Session[T], email string) (*T, error) {
    var user T
    err := session.Where("email", email).First(context.Background(), &user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ❌ 不推荐：需要类型转换
func GetUserByEmailGeneric(session *pie.Session[interface{}], email string) (interface{}, error) {
    var user interface{}
    err := session.Where("email", email).First(context.Background(), &user)
    if err != nil {
        return nil, err
    }
    // 需要类型断言
    return user.(User), nil
}
```

### 性能考虑

泛型在 Go 中提供了零成本抽象，编译后的代码与手写代码性能相同：

```go
// 泛型代码
func FindByID[T any](session *pie.Session[T], id string) (*T, error) {
    var result T
    err := session.FindByID(context.Background(), id, &result)
    return &result, err
}

// 编译后等同于手写的非泛型代码
func FindUserByID(session *pie.Session[User], id string) (*User, error) {
    var result User
    err := session.FindByID(context.Background(), id, &result)
    return &result, err
}
```

## 常见问题

### 泛型类型不匹配的错误处理

```go
// 问题：类型不匹配
var users []User
err := session.Find(context.Background(), &users) // ✅ 正确

// ❌ 错误：类型不匹配
var users []string
err := session.Find(context.Background(), &users) // 编译错误

// 解决方案：确保类型一致
var users []User
err := session.Find(context.Background(), &users)
```

### 如何在不同泛型类型之间转换

```go
// 如果需要转换，使用中间变量
func ConvertUsersToDTOs(users []User) []UserDTO {
    var dtos []UserDTO
    for _, user := range users {
        dtos = append(dtos, UserDTO{
            ID:    user.ID.Hex(),
            Name:  user.Name,
            Email: user.Email,
        })
    }
    return dtos
}

// 或者使用泛型转换函数
func ConvertSlice[T, U any](slice []T, converter func(T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = converter(item)
    }
    return result
}
```

### 泛型与 BSON 序列化的配合

```go
// BSON 标签与泛型配合使用
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Age       int                `bson:"age"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// 泛型确保类型安全，BSON 标签处理序列化
session := pie.Table[User](engine)
user := &User{Name: "张三", Email: "zhangsan@example.com"}
result, err := session.Insert(context.Background(), user)
```

### 常见编译错误及解决方案

```go
// 错误 1：类型参数缺失
// ❌ 编译错误
session := pie.Table(engine)

// ✅ 正确
session := pie.Table[User](engine)

// 错误 2：类型不匹配
// ❌ 编译错误
var users []string
err := session.Find(context.Background(), &users)

// ✅ 正确
var users []User
err := session.Find(context.Background(), &users)

// 错误 3：泛型约束不满足
// ❌ 如果 User 没有实现 Model 接口
func ProcessModel[T Model](session *pie.Session[T], ctx context.Context) error {
    // ...
}

// ✅ 确保 User 实现了 Model 接口
type User struct {
    // ...
}

func (u *User) GetID() string { return u.ID.Hex() }
func (u *User) GetCreatedAt() time.Time { return u.CreatedAt }
```

## 总结

泛型是 Pie 的核心特性，提供了：

1. **类型安全**：编译时检查，避免运行时错误
2. **代码复用**：一套代码处理多种类型
3. **性能优化**：零成本抽象，无运行时开销
4. **开发体验**：更好的 IDE 支持和代码提示

通过合理使用泛型，您可以编写更安全、更高效、更易维护的 MongoDB 操作代码。

## 下一步

- [查询构建器](/core-features/query-builder/) - 学习查询功能
- [聚合查询](/core-features/aggregation/) - 学习数据聚合
- [事务管理](/core-features/transactions/) - 掌握事务操作
- [缓存支持](/advanced/cache/) - 提升查询性能
