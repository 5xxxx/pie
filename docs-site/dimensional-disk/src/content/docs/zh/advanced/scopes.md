---
title: 查询作用域
description: 学习 Pie 的查询作用域功能，提高代码复用性
---

# 查询作用域

Pie 的查询作用域功能允许您定义可复用的查询逻辑，提高代码复用性和可维护性。

## 基本用法

### 定义作用域

```go
// 定义基础作用域
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

func AdultScope(field string) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where(field, pie.Gte(field, 18))
    }
}
```

### 使用作用域

```go
// 应用单个作用域
var users []User
err := session.Scopes(ActiveScope("status")).Find(ctx, &users)

// 应用多个作用域
err = session.Scopes(
    ActiveScope("status"),
    RecentScope("created_at", 30),
    AdultScope("age"),
).Find(ctx, &users)

// 作用域与普通条件组合
err = session.
    Scopes(ActiveScope("status")).
    Where("email", pie.Like("email", "%@company.com")).
    Find(ctx, &users)
```

## 高级用法

### 参数化作用域

```go
// 带参数的作用域
func StatusScope(status string) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("status", status)
    }
}

func AgeRangeScope(minAge, maxAge int) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("age", pie.Between("age", minAge, maxAge))
    }
}

func DateRangeScope(field string, startDate, endDate time.Time) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.WhereBetween(field, startDate, endDate)
    }
}

// 使用参数化作用域
err := session.Scopes(
    StatusScope("active"),
    AgeRangeScope(18, 65),
    DateRangeScope("created_at", startDate, endDate),
).Find(ctx, &users)
```

### 条件作用域

```go
// 根据条件应用作用域
func ConditionalScope(condition bool, scope pie.ScopeFunc) pie.ScopeFunc {
    if condition {
        return scope
    }
    return func(q *pie.Query) {
        // 空作用域，不添加任何条件
    }
}

// 使用条件作用域
includeInactive := false
err := session.Scopes(
    ConditionalScope(includeInactive, ActiveScope("status")),
    RecentScope("created_at", 30),
).Find(ctx, &users)
```

### 复合作用域

```go
// 复合作用域
func UserScope(userType string, includeInactive bool) pie.ScopeFunc {
    return func(q *pie.Query) {
        // 根据用户类型应用不同条件
        switch userType {
        case "admin":
            q.Where("role", "admin")
        case "user":
            q.Where("role", "user")
        case "premium":
            q.Where("role", "premium")
        }
        
        // 根据是否包含非活跃用户
        if !includeInactive {
            q.Where("status", "active")
        }
    }
}

// 使用复合作用域
err := session.Scopes(
    UserScope("premium", false),
    RecentScope("created_at", 90),
).Find(ctx, &users)
```

## 实际应用场景

### 用户管理作用域

```go
// 用户相关作用域
func ActiveUsersScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("status", "active")
    }
}

func VerifiedUsersScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("verified", true)
    }
}

func RecentUsersScope(days int) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.WhereRecentDays("created_at", days)
    }
}

func AdminUsersScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("role", "admin")
    }
}

func PremiumUsersScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("subscription_type", "premium")
    }
}

// 使用用户作用域
func getActiveVerifiedUsers() ([]User, error) {
    var users []User
    err := session.Scopes(
        ActiveUsersScope(),
        VerifiedUsersScope(),
    ).Find(ctx, &users)
    return users, err
}

func getRecentAdminUsers(days int) ([]User, error) {
    var users []User
    err := session.Scopes(
        AdminUsersScope(),
        RecentUsersScope(days),
    ).Find(ctx, &users)
    return users, err
}
```

### 订单管理作用域

```go
// 订单相关作用域
func CompletedOrdersScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("status", "completed")
    }
}

func PendingOrdersScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("status", "pending")
    }
}

func HighValueOrdersScope(minAmount float64) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("total", pie.Gte("total", minAmount))
    }
}

func RecentOrdersScope(days int) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.WhereRecentDays("created_at", days)
    }
}

func UserOrdersScope(userID bson.ObjectID) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("user_id", userID)
    }
}

// 使用订单作用域
func getCompletedHighValueOrders(minAmount float64) ([]Order, error) {
    var orders []Order
    err := session.Scopes(
        CompletedOrdersScope(),
        HighValueOrdersScope(minAmount),
    ).Find(ctx, &orders)
    return orders, err
}

func getUserRecentOrders(userID bson.ObjectID, days int) ([]Order, error) {
    var orders []Order
    err := session.Scopes(
        UserOrdersScope(userID),
        RecentOrdersScope(days),
    ).Find(ctx, &orders)
    return orders, err
}
```

### 产品管理作用域

```go
// 产品相关作用域
func AvailableProductsScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("status", "available")
        q.Where("stock", pie.Gt("stock", 0))
    }
}

func FeaturedProductsScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("featured", true)
    }
}

func CategoryScope(category string) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("category", category)
    }
}

func PriceRangeScope(minPrice, maxPrice float64) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.WhereBetween("price", minPrice, maxPrice)
    }
}

func SearchScope(keyword string) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("name", pie.Like("name", "%"+keyword+"%"))
    }
}

// 使用产品作用域
func getFeaturedElectronics() ([]Product, error) {
    var products []Product
    err := session.Scopes(
        AvailableProductsScope(),
        FeaturedProductsScope(),
        CategoryScope("electronics"),
    ).Find(ctx, &products)
    return products, err
}

func searchProducts(keyword string, minPrice, maxPrice float64) ([]Product, error) {
    var products []Product
    err := session.Scopes(
        AvailableProductsScope(),
        SearchScope(keyword),
        PriceRangeScope(minPrice, maxPrice),
    ).Find(ctx, &products)
    return products, err
}
```

## 作用域组合

### 作用域链

```go
// 创建作用域链
func createUserQueryScopes() []pie.ScopeFunc {
    return []pie.ScopeFunc{
        ActiveUsersScope(),
        VerifiedUsersScope(),
        RecentUsersScope(30),
    }
}

// 使用作用域链
func getActiveVerifiedRecentUsers() ([]User, error) {
    var users []User
    err := session.Scopes(createUserQueryScopes()...).Find(ctx, &users)
    return users, err
}
```

### 动态作用域

```go
// 动态构建作用域
func buildDynamicScopes(filters map[string]any) []pie.ScopeFunc {
    var scopes []pie.ScopeFunc
    
    if status, ok := filters["status"]; ok {
        scopes = append(scopes, StatusScope(status.(string)))
    }
    
    if minAge, ok := filters["min_age"]; ok {
        if maxAge, ok := filters["max_age"]; ok {
            scopes = append(scopes, AgeRangeScope(minAge.(int), maxAge.(int)))
        }
    }
    
    if startDate, ok := filters["start_date"]; ok {
        if endDate, ok := filters["end_date"]; ok {
            scopes = append(scopes, DateRangeScope("created_at", 
                startDate.(time.Time), endDate.(time.Time)))
        }
    }
    
    return scopes
}

// 使用动态作用域
func searchUsersWithFilters(filters map[string]any) ([]User, error) {
    var users []User
    scopes := buildDynamicScopes(filters)
    err := session.Scopes(scopes...).Find(ctx, &users)
    return users, err
}
```

### 作用域继承

```go
// 基础作用域
func BaseUserScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where("deleted_at", pie.Null("deleted_at"))
    }
}

// 继承基础作用域
func ActiveUserScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        BaseUserScope()(q) // 应用基础作用域
        q.Where("status", "active")
    }
}

func AdminUserScope() pie.ScopeFunc {
    return func(q *pie.Query) {
        BaseUserScope()(q) // 应用基础作用域
        q.Where("role", "admin")
    }
}
```

## 性能优化

### 1. 作用域缓存

```go
// 缓存常用作用域
var (
    activeUsersScope = ActiveUsersScope()
    verifiedUsersScope = VerifiedUsersScope()
    recentUsersScope = RecentUsersScope(30)
)

func getCachedScopedUsers() ([]User, error) {
    var users []User
    err := session.Scopes(
        activeUsersScope,
        verifiedUsersScope,
        recentUsersScope,
    ).Find(ctx, &users)
    return users, err
}
```

### 2. 作用域预编译

```go
// 预编译作用域查询
func precompileScopedQuery() *pie.Query {
    return session.Scopes(
        ActiveUsersScope(),
        VerifiedUsersScope(),
        RecentUsersScope(30),
    )
}

func usePrecompiledQuery() ([]User, error) {
    var users []User
    err := precompileScopedQuery().Find(ctx, &users)
    return users, err
}
```

### 3. 作用域索引优化

```go
// 确保作用域使用的字段有索引
func createScopeIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 为常用作用域字段创建索引
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"verified", 1},
        {"created_at", -1},
    })
    if err != nil {
        return err
    }
    
    return nil
}
```

## 最佳实践

### 1. 作用域命名规范

```go
// 使用清晰的作用域命名
func ActiveUsersScope() pie.ScopeFunc { ... }
func VerifiedUsersScope() pie.ScopeFunc { ... }
func RecentUsersScope(days int) pie.ScopeFunc { ... }
func AdminUsersScope() pie.ScopeFunc { ... }
func PremiumUsersScope() pie.ScopeFunc { ... }

// 避免模糊的命名
func Scope1() pie.ScopeFunc { ... } // 不好
func UserScope() pie.ScopeFunc { ... } // 不够具体
```

### 2. 作用域文档化

```go
// UserScopes 包含所有用户相关的作用域
var UserScopes = struct {
    // ActiveUsersScope 返回状态为 active 的用户
    ActiveUsers func() pie.ScopeFunc
    
    // VerifiedUsersScope 返回已验证的用户
    VerifiedUsers func() pie.ScopeFunc
    
    // RecentUsersScope 返回最近 N 天创建的用户
    RecentUsers func(days int) pie.ScopeFunc
}{
    ActiveUsers: ActiveUsersScope,
    VerifiedUsers: VerifiedUsersScope,
    RecentUsers: RecentUsersScope,
}
```

### 3. 作用域测试

```go
func TestUserScopes(t *testing.T) {
    // 测试 ActiveUsersScope
    query := session.Scopes(ActiveUsersScope())
    conditions := query.GetConditions()
    
    assert.Contains(t, conditions, bson.D{{"status", "active"}})
    
    // 测试 RecentUsersScope
    query = session.Scopes(RecentUsersScope(7))
    conditions = query.GetConditions()
    
    // 验证日期条件
    assert.NotEmpty(t, conditions)
}
```

### 4. 作用域组合测试

```go
func TestScopeCombination(t *testing.T) {
    // 测试多个作用域组合
    query := session.Scopes(
        ActiveUsersScope(),
        VerifiedUsersScope(),
        RecentUsersScope(30),
    )
    
    conditions := query.GetConditions()
    
    // 验证所有条件都存在
    assert.Contains(t, conditions, bson.D{{"status", "active"}})
    assert.Contains(t, conditions, bson.D{{"verified", true}})
    // 验证日期条件
    assert.NotEmpty(t, conditions)
}
```

## 下一步

- [日志监控](/advanced/logging/) - 学习查询日志
- [高级聚合](/advanced/aggregation-advanced/) - 深入了解聚合功能
- [最佳实践](/best-practices/) - 掌握开发最佳实践
