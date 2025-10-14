---
title: 查询构建器
description: 深入了解 Pie 的智能查询构建器，支持复杂查询条件构建
---

# 查询构建器

Pie 提供了丰富的查询方法，支持链式调用，让您能够构建复杂的查询条件。

## 基础查询

### 等值查询

```go
// 精确匹配
session.Where("status", "active")
session.Where("age", 25)
session.Where("email", "user@example.com")
```

### 比较查询

```go
// 大于等于
session.Where("age", pie.Gte("age", 18))

// 大于
session.Where("age", pie.Gt("age", 18))

// 小于等于
session.Where("age", pie.Lte("age", 65))

// 小于
session.Where("age", pie.Lt("age", 65))

// 不等于
session.Where("status", pie.Ne("status", "deleted"))

// 范围查询
session.Where("age", pie.Between("age", 18, 65))
session.Where("age", pie.NotBetween("age", 18, 65))
```

### 包含查询

```go
// 包含在数组中
session.WhereIn("role", []string{"admin", "user", "moderator"})
session.WhereIn("status", []int{1, 2, 3})

// 不包含在数组中
session.WhereNotIn("status", []string{"deleted", "banned"})

// 数组包含元素
session.WhereArrayContains("tags", "golang")
session.WhereArrayNotContains("tags", "deprecated")
```

### 空值查询

```go
// 字段为空
session.WhereNull("deleted_at")
session.WhereNull("updated_at")

// 字段不为空
session.WhereNotNull("email")
session.WhereNotNull("phone")
```

## 模糊查询

### 字符串匹配

```go
// LIKE 查询
session.WhereLike("name", "%张%")
session.WhereLike("email", "%@gmail.com")

// 前缀匹配
session.WhereStartsWith("email", "admin")
session.WhereStartsWith("domain", "www.")

// 后缀匹配
session.WhereEndsWith("email", "@company.com")
session.WhereEndsWith("file", ".pdf")

// 正则表达式
session.WhereRegex("name", "^张")
session.WhereRegex("email", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
```

### 文本搜索

```go
// 全文搜索
session.WhereText("golang mongodb tutorial")

// 指定语言
session.WhereText("golang mongodb tutorial", "en")
```

## 日期查询

### 日期范围

```go
// 最近 N 天
session.WhereRecentDays("created_at", 7)
session.WhereRecentDays("updated_at", 30)

// 指定月份
session.WhereMonth("created_at", time.Now().Month())
session.WhereMonth("created_at", time.January)

// 指定年份
session.WhereYear("created_at", 2024)
session.WhereYear("created_at", 2023)

// 日期范围
startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
session.WhereBetween("created_at", startDate, endDate)

// 今天
session.WhereToday("created_at")

// 昨天
session.WhereYesterday("created_at")

// 本周
session.WhereThisWeek("created_at")

// 上周
session.WhereLastWeek("created_at")

// 本月
session.WhereThisMonth("created_at")

// 上月
session.WhereLastMonth("created_at")
```

### 日期比较

```go
// 早于指定日期
session.WhereBefore("created_at", time.Now().AddDate(0, -1, 0))

// 晚于指定日期
session.WhereAfter("created_at", time.Now().AddDate(0, -6, 0))

// 在指定日期之后
session.WhereSince("created_at", time.Now().AddDate(0, -1, 0))
```

## 数组查询

### 数组操作

```go
// 数组大小
session.WhereArraySize("tags", 3)
session.WhereArraySize("tags", pie.Gte("size", 1))

// 数组包含所有元素
session.WhereArrayContainsAll("tags", []string{"golang", "mongodb"})

// 数组包含任意元素
session.WhereArrayContainsAny("tags", []string{"golang", "python"})

// 数组为空
session.WhereArrayEmpty("tags")

// 数组不为空
session.WhereArrayNotEmpty("tags")
```

### 嵌套数组查询

```go
// 嵌套字段查询
session.Where("address.city", "北京")
session.Where("profile.settings.theme", "dark")

// 嵌套数组查询
session.WhereArrayContains("users.roles", "admin")
```

## 复杂条件

### 逻辑操作符

```go
// AND 条件
session.Where("status", "active").
    Where("age", pie.Gte("age", 18)).
    Where("email", pie.Like("email", "%@company.com"))

// OR 条件
session.Where("status", "active").
    OrWhere("role", "admin")

// 复杂 OR 条件
session.Where("status", "active").
    OrWhere(func(q *pie.Query) {
        q.Where("role", "admin").WhereBetween("age", 30, 50)
    })

// NOT 条件
session.WhereNot("status", "deleted")
session.WhereNot("role", "guest")
```

### 嵌套条件

```go
// 嵌套 AND
session.Where(func(q *pie.Query) {
    q.Where("status", "active").
        Where("age", pie.Gte("age", 18))
}).Where("email", pie.Like("email", "%@company.com"))

// 嵌套 OR
session.Where(func(q *pie.Query) {
    q.Where("status", "active").
        OrWhere("role", "admin")
}).Where("age", pie.Between("age", 18, 65))

// 复杂嵌套
session.Where(func(q *pie.Query) {
    q.Where("status", "active").
        OrWhere(func(subQ *pie.Query) {
            subQ.Where("role", "admin").
                Where("verified", true)
        })
}).Where("age", pie.Gte("age", 18))
```

## 排序和限制

### 排序

```go
// 单字段排序
session.OrderBy("created_at")
session.OrderByDesc("created_at")

// 多字段排序
session.OrderBy("status").OrderByDesc("created_at")
session.OrderBy("name").OrderBy("age")

// 随机排序
session.OrderByRandom()
```

### 限制和跳过

```go
// 限制数量
session.Limit(10)

// 跳过记录
session.Skip(20)

// 分页
session.Skip(20).Limit(10)

// 获取最新记录
session.Latest("created_at", 10)

// 获取最旧记录
session.Oldest("created_at", 10)
```

## 字段选择

### 投影

```go
// 只选择指定字段
session.Select("name", "email", "status")

// 排除指定字段
session.Exclude("password", "secret_key")

// 选择嵌套字段
session.Select("name", "profile.avatar", "settings.theme")
```

## 聚合查询

### 计数

```go
// 总计数
count, err := session.Count(ctx)

// 条件计数
count, err := session.Where("status", "active").Count(ctx)

// 快速计数（不统计总数，更快）
count, err := session.Where("status", "active").QuickCount(ctx)
```

### 存在性检查

```go
// 检查是否存在
exists, err := session.Where("email", "test@example.com").Exists(ctx)

// 检查是否不存在
notExists, err := session.Where("email", "test@example.com").NotExists(ctx)
```

## 查询作用域

### 定义作用域

```go
// 定义可复用的查询作用域
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
users, err := session.Scopes(ActiveScope("status")).Find(ctx)

// 应用多个作用域
users, err = session.Scopes(
    ActiveScope("status"),
    RecentScope("created_at", 30),
    AdultScope("age"),
).Find(ctx)

// 作用域与普通条件组合
users, err = session.
    Scopes(ActiveScope("status")).
    Where("email", pie.Like("email", "%@company.com")).
    Find(ctx)
```

## 查询链式调用

Pie 的查询构建器支持链式调用，让您能够构建复杂的查询：

```go
users, err := session.
    Where("status", "active").
    Where("age", pie.Between("age", 18, 65)).
    WhereIn("role", []string{"user", "premium"}).
    WhereNotNull("email").
    WhereLike("name", "%张%").
    OrderBy("created_at").
    Limit(20).
    Find(ctx)
```

## 性能优化

### 使用索引

```go
// 确保查询字段有索引
session.Where("email", "test@example.com") // email 字段应该有索引
session.Where("status", "active").OrderBy("created_at") // 复合索引 (status, created_at)
```

### 字段投影

```go
// 只选择需要的字段，减少数据传输
users, err := session.
    Select("name", "email", "status").
    Where("status", "active").
    Find(ctx)
```

### 限制结果数量

```go
// 使用 Limit 避免返回过多数据
users, err := session.
    Where("status", "active").
    Limit(100).
    Find(ctx)
```

## 下一步

- [泛型使用指南](/core-features/generics/) - 深入了解泛型的使用方式和优势
- [结构体查询](/core-features/struct-query/) - 学习如何将 HTTP 参数转换为查询条件
- [分页查询](/core-features/pagination/) - 掌握分页实现
- [游标操作](/core-features/cursor/) - 学习高效的数据遍历
