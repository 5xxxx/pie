---
title: 聚合查询
description: 学习 Pie 的强大聚合框架，包含阶段构建器和表达式函数
---

# 聚合查询

Pie 提供了强大的聚合框架，包含阶段构建器和 100+ 表达式函数，用于复杂的数据处理和分析。

## 基本用法

### 创建聚合操作

```go
// 创建聚合操作
aggregate := pie.NewAggregate[User](engine).
    CollectionForStruct(User{})

// 或者指定集合名称
aggregate := pie.NewAggregate[User](engine).
    Collection("users")
```

### 基础聚合管道

```go
// 基础聚合管道
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

if err != nil {
    log.Fatal("Failed to execute aggregation:", err)
}

// 处理结果
for _, item := range result.Data {
    // item 是 bson.M 类型
    fmt.Printf("角色: %s, 总数: %d, 平均年龄: %.2f\n", 
        item["role"], item["total"], item["avgAge"])
}
```

## 阶段构建器

### Match 阶段

```go
// 基础匹配
result, err := aggregate.
    MatchStage().
        Where("status", "active").
        Where("age", pie.Gte("age", 18)).
        WhereIn("role", []string{"user", "admin"}).
        WhereLike("name", "%张%").
        WhereNotNull("email").
    Exec(ctx)

// 复杂逻辑条件
result, err := aggregate.
    MatchStage().
        And(
            bson.D{{"age", bson.D{{"$gte", 18}}}},
            bson.D{{"status", "active"}},
        ).
        Or(
            bson.D{{"role", "admin"}},
            bson.D{{"role", "moderator"}},
        ).
    Exec(ctx)

// 日期范围匹配
result, err := aggregate.
    MatchStage().
        WhereBetween("created_at", startDate, endDate).
        WhereRecentDays("updated_at", 30).
    Exec(ctx)
```

### Group 阶段

```go
// 基础分组
result, err := aggregate.
    GroupStage().
        By("category", "$category").
        By("year", pie.Year("$created_at")).
        Count("total").
        Sum("totalAmount", "$amount").
        Avg("avgAmount", "$amount").
        Max("maxAmount", "$amount").
        Min("minAmount", "$amount").
        Done().
    Exec(ctx)

// 复杂分组
result, err := aggregate.
    GroupStage().
        By("status", "$status").
        By("ageGroup", pie.Cond(
            pie.GteExpr("$age", 30),
            "adult",
            "young",
        )).
        Count("total").
        Sum("totalScore", "$score").
        Avg("avgScore", "$score").
        Push("names", "$name").
        AddToSet("uniqueEmails", "$email").
        First("firstUser", "$name").
        Last("lastUser", "$name").
        Done().
    Exec(ctx)
```

### Project 阶段

```go
// 字段投影
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
        Field("nameLength", pie.StrLenCP("$name")).
        Done().
    Exec(ctx)

// 计算字段
result, err := aggregate.
    ProjectStage().
        Field("total", pie.Add("$price", "$tax", "$shipping")).
        Field("discount", pie.Multiply("$price", 0.1)).
        Field("finalPrice", pie.Subtract("$price", pie.Multiply("$price", 0.1))).
        Field("rounded", pie.Round("$price", 2)).
        Done().
    Exec(ctx)
```

### Lookup 阶段

```go
// 基础关联查询
result, err := aggregate.
    LookupStage("orders", "_id", "user_id", "user_orders").
    Exec(ctx)

// 高级关联查询
result, err := aggregate.
    LookupStage("orders", "_id", "user_id", "user_orders").
        Pipeline(
            bson.M{"$match": bson.M{"status": "completed"}},
            bson.M{"$sort": bson.M{"created_at": -1}},
            bson.M{"$limit": 5},
        ).
        Done().
    Exec(ctx)

// 使用 Let 变量
result, err := aggregate.
    LookupStage("orders", "_id", "user_id", "user_orders").
        Let(pie.M{"userId": "$_id"}).
        Pipeline(
            bson.M{"$match": bson.M{
                "$expr": bson.M{"$eq": []string{"$user_id", "$$userId"}},
                "status": "completed",
            }},
        ).
        Done().
    Exec(ctx)
```

### Unwind 阶段

```go
// 展开数组
result, err := aggregate.
    UnwindStage("$tags").
    GroupStage().
        By("tag", "$tags").
        Count("count").
        Done().
    SortStage().Desc("count").
    Exec(ctx)

// 展开数组并保留空值
result, err := aggregate.
    UnwindStage("$tags").
        PreserveNullAndEmptyArrays(true).
        IncludeArrayIndex("tagIndex").
        Done().
    Exec(ctx)
```

### Facet 阶段

```go
// 多维度分析
result, err := aggregate.
    FacetStage().
        Facet("activeUsers",
            bson.M{"$match": bson.M{"status": "active"}},
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

## 表达式函数

### 日期表达式

```go
// 日期操作
result, err := aggregate.
    AddFieldsStage().
        Add("year", pie.Year("$created_at")).
        Add("month", pie.Month("$created_at")).
        Add("dayOfWeek", pie.DayOfWeek("$created_at")).
        Add("formattedDate", pie.DateToString("$created_at", "%Y-%m-%d", "UTC")).
        Add("daysSince", pie.DateDiff("$created_at", pie.Now(), "day")).
        Add("nextWeek", pie.DateAdd("$created_at", 7, "day")).
        Add("isWeekend", pie.In(pie.DayOfWeek("$created_at"), []int{0, 6})).
        Done().
    Exec(ctx)
```

### 算术表达式

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
        Add("abs", pie.Abs("$value")).
        Add("ceil", pie.Ceil("$value")).
        Add("floor", pie.Floor("$value")).
        Done().
    Exec(ctx)
```

### 字符串表达式

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
        Add("trimmed", pie.Trim("$name")).
        Add("replaced", pie.ReplaceAll("$name", " ", "_")).
        Done().
    Exec(ctx)
```

### 数组表达式

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
        Add("sorted", pie.SortArray("$items", 1)).
        Add("unique", pie.SetUnion("$items")).
        Done().
    Exec(ctx)
```

### 条件表达式

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
        Add("isValid", pie.And(
            pie.NeExpr("$email", ""),
            pie.GteExpr("$age", 18),
        )).
        Done().
    Exec(ctx)
```

## 实际应用场景

### 用户统计分析

```go
func getUserStatistics() error {
    result, err := aggregate.
        MatchStage().
            Where("status", "active").
            Where("created_at", pie.Gte("created_at", time.Now().AddDate(0, -1, 0))).
        GroupStage().
            By("role", "$role").
            Count("total").
            Avg("avgAge", "$age").
            Max("maxAge", "$age").
            Min("minAge", "$age").
            Push("names", "$name").
            Done().
        SortStage().Desc("total").
        Exec(context.Background())
    
    if err != nil {
        return err
    }
    
    fmt.Println("用户统计:")
    for _, item := range result.Data {
        fmt.Printf("角色: %s, 总数: %d, 平均年龄: %.2f\n", 
            item["role"], item["total"], item["avgAge"])
    }
    
    return nil
}
```

### 销售数据分析

```go
func getSalesAnalysis() error {
    result, err := aggregate.
        MatchStage().
            Where("status", "completed").
            WhereBetween("created_at", startDate, endDate).
        AddFieldsStage().
            Add("year", pie.Year("$created_at")).
            Add("month", pie.Month("$created_at")).
            Add("totalAmount", pie.Add("$price", "$tax", "$shipping")).
            Done().
        GroupStage().
            By("year", "$year").
            By("month", "$month").
            Count("orderCount").
            Sum("totalRevenue", "$totalAmount").
            Avg("avgOrderValue", "$totalAmount").
            Done().
        SortStage().Asc("year").Asc("month").
        Exec(context.Background())
    
    if err != nil {
        return err
    }
    
    fmt.Println("销售分析:")
    for _, item := range result.Data {
        fmt.Printf("%d年%d月: 订单数 %d, 总收入 %.2f, 平均订单价值 %.2f\n",
            item["year"], item["month"], item["orderCount"], 
            item["totalRevenue"], item["avgOrderValue"])
    }
    
    return nil
}
```

### 产品推荐

```go
func getProductRecommendations(userID primitive.ObjectID) error {
    result, err := aggregate.
        MatchStage().Where("_id", userID).
        LookupStage("orders", "_id", "user_id", "user_orders").
            Pipeline(
                bson.M{"$match": bson.M{"status": "completed"}},
                bson.M{"$unwind": "$items"},
                bson.M{"$group": bson.M{
                    "_id": "$items.category",
                    "count": bson.M{"$sum": 1},
                }},
            ).
            Done().
        UnwindStage("$user_orders").
        GroupStage().
            By("category", "$user_orders._id").
            Sum("purchaseCount", "$user_orders.count").
            Done().
        LookupStage("products", "category", "category", "recommended_products").
            Pipeline(
                bson.M{"$match": bson.M{"status": "active"}},
                bson.M{"$sample": bson.M{"size": 5}},
            ).
            Done().
        Exec(context.Background())
    
    if err != nil {
        return err
    }
    
    fmt.Println("推荐产品:")
    for _, item := range result.Data {
        fmt.Printf("分类: %s, 购买次数: %d\n", 
            item["category"], item["purchaseCount"])
    }
    
    return nil
}
```

## 性能优化

### 1. 使用索引优化

```go
// 确保匹配字段有索引
result, err := aggregate.
    MatchStage().
        Where("status", "active").  // status 字段应该有索引
        Where("created_at", pie.Gte("created_at", startDate)). // created_at 应该有索引
    Exec(ctx)
```

### 2. 合理使用投影

```go
// 只选择需要的字段
result, err := aggregate.
    ProjectStage().
        Include("name", "email", "status"). // 只选择需要的字段
        Done().
    Exec(ctx)
```

### 3. 限制结果数量

```go
// 使用 Limit 阶段限制结果
result, err := aggregate.
    MatchStage().Where("status", "active").
    SortStage().Desc("created_at").
    LimitStage(100). // 限制结果数量
    Exec(ctx)
```

### 4. 使用 $facet 进行并行处理

```go
// 使用 $facet 并行处理多个聚合
result, err := aggregate.
    FacetStage().
        Facet("stats", 
            bson.M{"$group": bson.M{
                "_id": nil,
                "total": bson.M{"$sum": 1},
                "avg": bson.M{"$avg": "$score"},
            }},
        ).
        Facet("topUsers",
            bson.M{"$sort": bson.M{"score": -1}},
            bson.M{"$limit": 10},
        ).
        Done().
    Exec(ctx)
```

## 错误处理

```go
func handleAggregationError(err error) {
    if pie.IsNotFoundError(err) {
        log.Println("No data found for aggregation")
        return
    }
    
    if pie.IsTimeoutError(err) {
        log.Println("Aggregation timeout")
        return
    }
    
    log.Printf("Aggregation error: %v", err)
}
```

## 最佳实践

### 1. 分阶段构建复杂聚合

```go
func buildComplexAggregation() *pie.Aggregate[User] {
    return aggregate.
        MatchStage().
            Where("status", "active").
            Where("age", pie.Gte("age", 18)).
        AddFieldsStage().
            Add("ageGroup", pie.Cond(
                pie.GteExpr("$age", 30),
                "adult",
                "young",
            )).
            Done().
        GroupStage().
            By("ageGroup", "$ageGroup").
            Count("total").
            Avg("avgScore", "$score").
            Done().
        SortStage().Desc("total")
}
```

### 2. 使用类型安全的聚合

```go
type UserStats struct {
    Role     string  `bson:"role"`
    Total    int     `bson:"total"`
    AvgAge   float64 `bson:"avgAge"`
    MaxAge   int     `bson:"maxAge"`
    MinAge   int     `bson:"minAge"`
}

func getTypedAggregation() ([]UserStats, error) {
    result, err := aggregate.
        GroupStage().
            By("role", "$role").
            Count("total").
            Avg("avgAge", "$age").
            Max("maxAge", "$age").
            Min("minAge", "$age").
            Done().
        Exec(context.Background())
    
    if err != nil {
        return nil, err
    }
    
    var stats []UserStats
    for _, item := range result.Data {
        var stat UserStats
        if err := bson.Unmarshal(bson.M(item), &stat); err != nil {
            continue
        }
        stats = append(stats, stat)
    }
    
    return stats, nil
}
```

## 下一步

- [事务管理](/core-features/transactions/) - 学习事务操作
- [高级聚合](/advanced/aggregation-advanced/) - 深入了解阶段构建器和表达式函数
- [缓存支持](/advanced/cache/) - 提升查询性能
