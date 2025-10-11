# Pie - MongoDB ORM 框架

Pie 是一个现代化的 Go 语言 MongoDB ORM 框架，提供类型安全、高性能、功能丰富的数据库操作体验。

## 特性

- **类型安全**: 基于 Go 泛型，提供编译时类型检查
- **智能查询构建器**: 直观的链式调用，支持复杂查询条件
- **结构体查询**: HTTP 请求参数直接转换为查询条件
- **缓存支持**: 内置多级缓存，支持 Redis 和内存缓存
- **钩子系统**: 完整的生命周期钩子支持
- **事务管理**: 简单易用的事务操作
- **索引管理**: 自动化的索引创建和管理
- **高级聚合**: 全面的阶段构建器，包含 100+ 表达式函数
- **变更流**: 实时数据变更监听
- **软删除**: 内置软删除功能
- **分页查询**: 高效的分页实现
- **查询日志**: 详细的查询日志和性能监控
- **批量操作**: 高效的批量写入操作

## 安装

```bash
go get github.com/5xxxx/pie
```

## 快速开始

### 1. 连接数据库

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

### 2. 定义模型

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    Age       int           `bson:"age"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
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
```

### 3. 基本操作

```go
// 创建类型安全的会话
session := pie.Table[User](engine)

// 插入文档
user := &User{
    Name:  "张三",
    Email: "zhangsan@example.com",
    Age:   25,
}
result, err := session.Insert(context.Background(), user)

// 查询文档
users, err := session.
    Where("age", pie.Gte("age", 18)).
    OrderBy("name").
    Limit(10).
    Find(context.Background())

// 更新文档
updateResult, err := session.
    Where("email", "zhangsan@example.com").
    Update(context.Background(), bson.D{{"$set", bson.D{{"age", 26}}}})

// 删除文档
deleteResult, err := session.
    Where("email", "zhangsan@example.com").
    Delete(context.Background())
```

## 核心功能详解

### 1. 智能查询构建器

Pie 提供了丰富的查询方法，支持链式调用：

```go
// 基础查询
session.Where("status", "active")
session.Where("age", pie.Gte("age", 18))
session.WhereIn("role", []string{"admin", "user"})
session.WhereBetween("age", 18, 65)
session.WhereNull("deleted_at")
session.WhereNotNull("email")

// 模糊查询
session.WhereLike("name", "%张%")
session.WhereStartsWith("email", "admin")
session.WhereEndsWith("domain", ".com")

// 日期查询
session.WhereRecentDays("created_at", 7)
session.WhereMonth("created_at", time.Now().Month())
session.WhereYear("created_at", 2024)

// 复杂条件
session.Where("status", "active").
    OrWhere(func(q *pie.Query) {
        q.Where("role", "admin").WhereBetween("age", 30, 50)
    })
```

### 2. 结构体查询

可将 HTTP 请求参数直接转换为查询条件：

```go
type UserQuery struct {
    Name     string   `pie:"name,like,omitempty" json:"name"`
    Email    string   `pie:"email,omitempty" json:"email"`
    MinAge   int      `pie:"age,gte,omitempty" json:"min_age"`
    MaxAge   int      `pie:"age,lte,omitempty" json:"max_age"`
    Status   []string `pie:"status,in,omitempty" json:"status"`
    Role     string   `pie:"role,omitempty" json:"role"`
}

// HTTP 请求: GET /users?name=张&min_age=20&max_age=40&status=active,pending
query := UserQuery{
    Name:   "张",
    MinAge: 20,
    MaxAge: 40,
    Status: []string{"active", "pending"},
}

var users []User
err := session.WhereStruct(query).Find(ctx, &users)
```

### 3. 便捷方法

```go
// 按ID查找
user, err := session.FindByID(ctx, userID)

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
```

### 4. 分页查询

```go
// 完整分页（包含总数统计）
result, err := session.
    Where("status", "active").
    Paginate(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
// result.Total, result.TotalPages, result.HasNext, result.HasPrev

// 简单分页（不统计总数，更快）
simpleResult, err := session.
    PaginateSimple(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
```

### 5. 游标操作

```go
// 获取游标
cursor, err := session.
    Where("status", "active").
    FindCursor(ctx)

// 方式1: 使用 Next() 和 Decode()
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    // 处理用户
}

// 方式2: 使用 All() 一次性获取
users, err := cursor.All(ctx)

// 方式3: 使用 Iterate() 迭代处理
cursor.Iterate(ctx, func(user *User) error {
    // 处理用户
    return nil
})

// 方式4: 使用 Take() 获取前N个
topUsers, err := cursor.Take(ctx, 5)

// 方式5: 使用 First() 获取第一个
firstUser, err := cursor.First(ctx)

cursor.Close(ctx)
```

### 6. 批量操作

```go
// 创建批量写入操作
bulkWrite := pie.NewBulkWrite[User](engine).
    CollectionForStruct(User{})

// 插入多个文档
bulkWrite.InsertOne(&User{Name: "用户1"})
bulkWrite.InsertOne(&User{Name: "用户2"})

// 更新操作
bulkWrite.UpdateOne(
    bson.D{{"email", "old@example.com"}},
    bson.D{{"$set", bson.D{{"email", "new@example.com"}}}},
)

// 批量更新
bulkWrite.UpdateMany(
    bson.D{{"status", "inactive"}},
    bson.D{{"$set", bson.D{{"status", "active"}}}},
)

// 删除操作
bulkWrite.DeleteMany(bson.D{{"age", bson.D{{"$lt", 18}}}})

// 执行批量操作
result, err := bulkWrite.ExecuteOrdered(ctx)
```

### 7. 聚合查询

Pie 提供了强大的聚合框架，包含阶段构建器和表达式函数，用于复杂的数据处理。

#### 基础聚合

```go
// 创建聚合操作
aggregate := pie.NewAggregate[User](engine).
    CollectionForStruct(User{})

// 使用阶段构建器构建聚合管道
result, err := aggregate.
    MatchStage().Where("status", "active").
    GroupStage().
        By("role", "$role").
        Count("total").
        Avg("avgAge", "$age").
        Max("maxAge", "$age").
        Min("minAge", "$age").
        Done().
    SortStage().Desc("total").
    Exec(ctx)

// 处理结果
for _, item := range result.Data {
    // item 是 bson.M 类型
}
```

#### 高级聚合与表达式

```go
// 使用表达式的复杂聚合
result, err := aggregate.
    MatchStage().
        Where("active", true).
        Between("age", 18, 65).
        In("status", "active", "pending").
    AddFieldsStage().
        Add("ageGroup", pie.Cond(
            pie.GteExpr("$age", 30),
            "adult",
            "young",
        )).
        Add("fullName", pie.Concat("$firstName", " ", "$lastName")).
        Add("scoreRounded", pie.Round("$score", 1)).
        Done().
    GroupStage().
        By("ageGroup", "$ageGroup").
        Count("total").
        Avg("avgScore", "$score").
        Push("names", "$fullName").
        Done().
    ProjectStage().
        Include("ageGroup", "total", "avgScore", "names").
        Field("nameCount", pie.SizeArray("$names")).
        Done().
    SortStage().Desc("total").
    LimitStage(10).
    Exec(ctx)
```

#### 关联查询操作

```go
// 与订单集合关联查询
result, err := aggregate.
    LookupStage("orders", "_id", "user_id", "user_orders").
        Pipeline(
            bson.M{"$match": bson.M{"status": "completed"}},
            bson.M{"$limit": 5},
        ).
        Done().
    AddFieldsStage().
        Add("orderCount", pie.SizeArray("$user_orders")).
        Add("totalSpent", pie.Sum("$user_orders.amount")).
        Done().
    MatchStage().Where("orderCount", pie.GtExpr(0, 0)).
    ProjectStage().
        Include("name", "email", "orderCount", "totalSpent").
        Done().
    Exec(ctx)
```

#### 分面分析

```go
// 使用分面进行多维度分析
result, err := aggregate.
    FacetStage().
        Facet("activeUsers",
            bson.M{"$match": bson.M{"active": true}},
            bson.M{"$count": "count"},
        ).
        Facet("scoreStats",
            bson.M{"$group": bson.M{
                "_id":      nil,
                "avgScore": bson.M{"$avg": "$score"},
                "maxScore": bson.M{"$max": "$score"},
                "minScore": bson.M{"$min": "$score"},
            }},
        ).
        Facet("ageGroups",
            bson.M{"$bucket": bson.M{
                "groupBy":    "$age",
                "boundaries": []int{0, 25, 30, 35, 100},
                "default":    "other",
                "output":     bson.M{"count": bson.M{"$sum": 1}},
            }},
        ).
        Done().
    Exec(ctx)
```

### 8. 事务管理

```go
// 使用引擎执行事务
err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
    // 在事务中执行操作
    session := pie.Table[User](engine)
    
    // 插入用户
    _, err := session.Insert(txCtx, &User{Name: "事务用户"})
    if err != nil {
        return err
    }
    
    // 更新其他集合
    orderSession := pie.Table[Order](engine)
    _, err = orderSession.Insert(txCtx, &Order{UserID: userID})
    return err
})

// 使用事务管理器
tx := pie.MustTransaction(engine)
err = tx.Execute(ctx, func(txCtx context.Context) error {
    // 事务操作
    return nil
})
```

### 9. 缓存支持

```go
// 启用内存缓存
engine.UseCache(pie.NewMemoryCache(), &pie.CacheConfig{
    TTL: 5 * time.Minute,
})

// 启用 Redis 缓存
redisCache := pie.NewRedisCache("localhost:6379", "", 0)
engine.UseCache(redisCache, &pie.CacheConfig{
    TTL: 10 * time.Minute,
})

// 启用二级缓存
engine.UseTwoLevelCache(
    pie.NewMemoryCache(),  // L1 缓存
    redisCache,            // L2 缓存
    &pie.TwoLevelCacheConfig{
        L1TTL: 1 * time.Minute,
        L2TTL: 10 * time.Minute,
    },
)

// 在会话中使用缓存
session := pie.Table[User](engine).
    WithCache(5 * time.Minute)

// 缓存查询结果
var users []User
err := session.
    Where("status", "active").
    Cache("active_users").
    Find(ctx, &users)
```

### 10. 钩子系统

```go
type User struct {
    // 字段定义...
}

// 创建前钩子
func (u *User) BeforeCreate(ctx context.Context) error {
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return nil
}

// 创建后钩子
func (u *User) AfterCreate(ctx context.Context) error {
    log.Printf("User %s created", u.Name)
    return nil
}

// 更新前钩子
func (u *User) BeforeUpdate(ctx context.Context) error {
    u.UpdatedAt = time.Now()
    return nil
}

// 更新后钩子
func (u *User) AfterUpdate(ctx context.Context) error {
    log.Printf("User %s updated", u.Name)
    return nil
}

// 删除前钩子
func (u *User) BeforeDelete(ctx context.Context) error {
    log.Printf("About to delete user %s", u.Name)
    return nil
}

// 删除后钩子
func (u *User) AfterDelete(ctx context.Context) error {
    log.Printf("User %s deleted", u.Name)
    return nil
}

// 查找后钩子
func (u *User) AfterFind(ctx context.Context) error {
    // 处理查找后的逻辑
    return nil
}
```

### 11. 软删除

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

// 软删除
err := session.Where("email", "test@example.com").SoftDelete(ctx)

// 恢复软删除
err := session.Where("email", "test@example.com").Restore(ctx)

// 强制删除（物理删除）
err := session.Where("email", "test@example.com").ForceDelete(ctx)

// 查询时自动排除软删除的记录
var users []User
err := session.Find(ctx, &users) // 自动过滤 deleted_at 不为 null 的记录
```

### 12. 索引管理

```go
// 创建索引管理器
indexes := pie.MustIndexes(engine)

// 为结构体创建索引
err := indexes.CreateIndexes(ctx, User{})

// 手动创建索引
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"email", 1},
}, &options.IndexOptions{
    Unique: pie.Bool(true),
})

// 创建复合索引
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"status", 1},
    {"created_at", -1},
})

// 删除索引
err := indexes.DropIndex(ctx, "users", "email_1")
```

### 13. 变更流监听

```go
// 创建变更流监听器
watcher := pie.NewWatcher[User](engine)

// 监听集合变更
err := watcher.
    WatchCollection().
    Filter(bson.D{{"operationType", "insert"}}).
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("New user inserted: %s", change.FullDocument.Name)
        return nil
    })

// 监听数据库变更
dbWatcher := pie.NewDatabaseWatcher[User](engine)
err = dbWatcher.
    WatchDatabase().
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("Database change: %s", change.OperationType)
        return nil
    })
```

### 14. 查询作用域

```go
// 定义作用域
func ActiveScope(field string) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where(field, "active")
    }
}

func RecentScope(field string, days int) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.WhereRecentDays(field, days)
    }
}

// 使用作用域
var users []User
err := session.
    Scopes(
        ActiveScope("status"),
        RecentScope("created_at", 30),
    ).
    Latest("created_at", 10).
    Find(ctx, &users)
```

### 15. 查询日志和监控

```go
// 启用查询日志
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithQueryLog(os.Stdout),
    pie.WithSlowQueryThreshold(50*time.Millisecond),
)

// 自定义日志格式
engine.SetQueryLogFormatter(func(entry *pie.LogEntry) string {
    return fmt.Sprintf("[%s] %s %s - %v", 
        entry.Timestamp.Format("15:04:05"),
        entry.Operation,
        entry.Collection,
        entry.Duration,
    )
})

// 设置慢查询阈值
engine.SetSlowQueryThreshold(100 * time.Millisecond)
```

## 高级功能

### 1. 聚合阶段构建器

Pie 提供了全面的 MongoDB 聚合管道阶段构建器，使复杂的数据处理变得直观且类型安全。

#### 匹配阶段构建器

```go
// 基础匹配，支持链式条件
result, err := aggregate.
    MatchStage().
        Where("status", "active").
        Between("age", 18, 65).
        In("role", "admin", "user").
        Regex("name", "^张").
        Exists("email", true).
        Text("搜索词").
    Exec(ctx)

// 复杂逻辑条件
result, err := aggregate.
    MatchStage().
        And(
            bson.D{{"age", bson.D{{"$gte", 18}}}},
            bson.D{{"active", true}},
        ).
        Or(
            bson.D{{"status", "active"}},
            bson.D{{"status", "pending"}},
        ).
    Exec(ctx)
```

#### 分组阶段构建器

```go
// 全面的分组聚合，支持多种累加器
result, err := aggregate.
    GroupStage().
        By("category", "$category").
        By("year", pie.Year("$created_at")).
        Count("total").
        Sum("totalAmount", "$amount").
        Avg("avgAmount", "$amount").
        Max("maxAmount", "$amount").
        Min("minAmount", "$amount").
        StdDevPop("stdDev", "$amount").
        Push("items", "$item").
        AddToSet("uniqueItems", "$item").
        First("firstItem", "$item").
        Last("lastItem", "$item").
        Done().
    Exec(ctx)
```

#### 投影阶段构建器

```go
// 字段投影和计算字段
result, err := aggregate.
    ProjectStage().
        Include("name", "email", "status").
        Exclude("password", "secret").
        Field("fullName", pie.Concat("$firstName", " ", "$lastName")).
        Field("ageGroup", pie.Cond(
            pie.GteExpr("$age", 30),
            "adult",
            "young",
        )).
        Slice("recentTags", "$tags", 3).
        Done().
    Exec(ctx)
```

#### 关联查询阶段构建器

```go
// 高级关联查询，支持管道
result, err := aggregate.
    LookupStage("orders", "_id", "user_id", "user_orders").
        Let(pie.M{"userId": "$_id"}).
        Pipeline(
            bson.M{"$match": bson.M{
                "$expr": bson.M{"$eq": []string{"$user_id", "$$userId"}},
                "status": "completed",
            }},
            bson.M{"$sort": bson.M{"created_at": -1}},
            bson.M{"$limit": 5},
        ).
        Done().
    Exec(ctx)
```

#### 展开阶段构建器

```go
// 数组展开，支持选项
result, err := aggregate.
    UnwindStage("$tags").
        PreserveNullAndEmptyArrays(true).
        IncludeArrayIndex("tagIndex").
        Done().
    GroupStage().
        By("tag", "$tags").
        Count("count").
        Done().
    SortStage().Desc("count").
    Exec(ctx)
```

#### 分面阶段构建器

```go
// 多维度分析
result, err := aggregate.
    FacetStage().
        Facet("activeUsers",
            bson.M{"$match": bson.M{"active": true}},
            bson.M{"$count": "count"},
        ).
        Facet("scoreDistribution",
            bson.M{"$bucket": bson.M{
                "groupBy":    "$score",
                "boundaries": []int{0, 60, 80, 90, 100},
                "default":    "other",
                "output":     bson.M{"count": bson.M{"$sum": 1}},
            }},
        ).
        Facet("topPerformers",
            bson.M{"$match": bson.M{"score": bson.M{"$gte": 90}}},
            bson.M{"$sort": bson.M{"score": -1}},
            bson.M{"$limit": 10},
        ).
        Done().
    Exec(ctx)
```

### 2. 聚合表达式函数

Pie 提供了 100+ 个表达式函数，用于复杂的数据转换和计算。

#### 日期表达式

```go
// 日期操作和格式化
result, err := aggregate.
    AddFieldsStage().
        Add("year", pie.Year("$created_at")).
        Add("month", pie.Month("$created_at")).
        Add("dayOfWeek", pie.DayOfWeek("$created_at")).
        Add("formattedDate", pie.DateToString("$created_at", "%Y-%m-%d", "UTC")).
        Add("daysSince", pie.DateDiff("$created_at", pie.Now(), "day")).
        Add("nextWeek", pie.DateAdd("$created_at", 7, "day")).
        Done().
    Exec(ctx)
```

#### 算术表达式

```go
// 数学运算
result, err := aggregate.
    AddFieldsStage().
        Add("total", pie.Add("$price", "$tax", "$shipping")).
        Add("discount", pie.Multiply("$price", 0.1)).
        Add("finalPrice", pie.Subtract("$price", pie.Multiply("$price", 0.1))).
        Add("rounded", pie.Round("$price", 2)).
        Add("power", pie.Pow("$base", 2)).
        Add("sqrt", pie.Sqrt("$value")).
        Done().
    Exec(ctx)
```

#### 字符串表达式

```go
// 字符串操作
result, err := aggregate.
    AddFieldsStage().
        Add("fullName", pie.Concat("$firstName", " ", "$lastName")).
        Add("upperName", pie.ToUpper("$name")).
        Add("lowerEmail", pie.ToLower("$email")).
        Add("initials", pie.Concat(
            pie.SubStr("$firstName", 0, 1),
            pie.SubStr("$lastName", 0, 1),
        )).
        Add("nameLength", pie.StrLenCP("$name")).
        Add("words", pie.Split("$description", " ")).
        Done().
    Exec(ctx)
```

#### 数组表达式

```go
// 数组操作
result, err := aggregate.
    AddFieldsStage().
        Add("firstItem", pie.First("$items")).
        Add("lastItem", pie.Last("$items")).
        Add("itemCount", pie.SizeArray("$items")).
        Add("firstThree", pie.Slice("$items", 3)).
        Add("filtered", pie.FilterArray("$items", pie.GtExpr("$$item", 0))).
        Add("mapped", pie.MapArray("$items", "item", pie.Multiply("$$item", 2))).
        Add("reversed", pie.ReverseArray("$items")).
        Done().
    Exec(ctx)
```

#### 条件表达式

```go
// 条件逻辑
result, err := aggregate.
    AddFieldsStage().
        Add("status", pie.Cond(
            pie.GteExpr("$score", 80),
            "excellent",
            pie.Cond(
                pie.GteExpr("$score", 60),
                "good",
                "needs_improvement",
            ),
        )).
        Add("displayName", pie.IfNull("$nickname", "$name")).
        Add("grade", pie.Switch([]pie.M{
            {"case": pie.GteExpr("$score", 90), "then": "A"},
            {"case": pie.GteExpr("$score", 80), "then": "B"},
            {"case": pie.GteExpr("$score", 70), "then": "C"},
        }, "F")).
        Done().
    Exec(ctx)
```

### 3. 自定义命名映射

```go
// 使用蛇形命名
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SnakeMapper{}),
)

// 使用驼峰命名
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.CamelMapper{}),
)

// 使用相同命名
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SameMapper{}),
)

// 自定义映射器
type CustomMapper struct{}

func (m CustomMapper) TableName(structName string) string {
    return "t_" + strings.ToLower(structName)
}

func (m CustomMapper) FieldName(fieldName string) string {
    return strings.ToLower(fieldName)
}
```

### 2. 配置选项

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

### 3. 错误处理

```go
// 检查特定错误类型
if pie.IsDuplicateKeyError(err) {
    log.Println("Duplicate key error")
}

if pie.IsNotFoundError(err) {
    log.Println("Document not found")
}

if pie.IsTimeoutError(err) {
    log.Println("Operation timeout")
}

// 获取错误详情
if mongoErr, ok := err.(pie.MongoError); ok {
    log.Printf("MongoDB error code: %d", mongoErr.Code)
    log.Printf("MongoDB error message: %s", mongoErr.Message)
}
```

## 性能优化

### 1. 连接池配置

```go
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMaxPoolSize(100),        // 最大连接数
    pie.WithMinPoolSize(5),          // 最小连接数
    pie.WithMaxIdleTime(30*time.Minute), // 最大空闲时间
)
```

### 2. 查询优化

```go
// 使用投影减少数据传输
var users []User
err := session.
    Select("name", "email").  // 只选择需要的字段
    Find(ctx, &users)

// 使用索引优化查询
err := session.
    Where("email", "test@example.com").  // email 字段有索引
    Find(ctx, &users)

// 使用游标处理大量数据
cursor, err := session.FindCursor(ctx)
defer cursor.Close(ctx)

for cursor.Next(ctx) {
    var user User
    cursor.Decode(&user)
    // 处理单个用户
}
```

### 3. 批量操作优化

```go
// 使用批量写入提高性能
bulkWrite := pie.NewBulkWrite[User](engine)
for _, user := range users {
    bulkWrite.InsertOne(user)
}
result, err := bulkWrite.ExecuteOrdered(ctx)
```

## 最佳实践

### 1. 模型设计

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name" pie:"index"`
    Email     string        `bson:"email" pie:"unique"`
    Age       int           `bson:"age"`
    Status    string        `bson:"status" pie:"index"`
    CreatedAt time.Time     `bson:"created_at" pie:"index"`
    UpdatedAt time.Time     `bson:"updated_at"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
}
```

### 2. 错误处理

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    session := pie.Table[User](s.engine)
    
    // 检查邮箱是否已存在
    exists, err := session.Where("email", user.Email).Exists(ctx)
    if err != nil {
        return fmt.Errorf("failed to check email existence: %w", err)
    }
    if exists {
        return errors.New("email already exists")
    }
    
    // 创建用户
    _, err = session.Insert(ctx, user)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

### 3. 事务使用

```go
func (s *UserService) TransferPoints(ctx context.Context, fromUserID, toUserID bson.ObjectID, points int) error {
    return s.engine.WithTransaction(ctx, func(txCtx context.Context) error {
        userSession := pie.Table[User](s.engine)
        
        // 减少发送方积分
        _, err := userSession.
            Where("_id", pie.ID(fromUserID)).
            Update(txCtx, bson.D{{"$inc", bson.D{{"points", -points}}}})
        if err != nil {
            return fmt.Errorf("failed to deduct points: %w", err)
        }
        
        // 增加接收方积分
        _, err = userSession.
            Where("_id", pie.ID(toUserID)).
            Update(txCtx, bson.D{{"$inc", bson.D{{"points", points}}}})
        if err != nil {
            return fmt.Errorf("failed to add points: %w", err)
        }
        
        return nil
    })
}
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

