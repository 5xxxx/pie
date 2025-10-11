---
title: 高级聚合
description: 全面的聚合阶段构建器和 100+ 表达式函数
---

# 高级聚合

Pie 提供全面的聚合阶段构建器和 100+ 表达式函数，用于复杂的数据处理。

## 阶段构建器

```go
aggregate := pie.NewAggregate[User](engine).CollectionForStruct(User{})

// $match 阶段
aggregate.MatchStage().Where("status", "active").Between("age", 18, 65).Done()

// $group 阶段
aggregate.GroupStage().
    By("status", "$status").
    Count("total").
    Avg("avgAge", "$age").
    Done()

// $project 阶段
aggregate.ProjectStage().Include("name", "email").Exclude("_id").Done()

// $lookup 阶段 (连接)
aggregate.LookupStage("orders", "_id", "user_id", "userOrders").Done()

// $facet 阶段 (分面搜索)
aggregate.FacetStage().
    Facet("activeUsers", bson.M{"$match": bson.M{"active": true}}).
    Facet("ageGroups", bson.M{"$bucket": bson.M{"groupBy": "$age", "boundaries": []int{0, 20, 40, 60, 100}}}).
    Done()
```

## 表达式函数

Pie 提供 100+ 聚合表达式函数：

```go
// 算术表达式
pie.Add(1, 2, "$field")
pie.Multiply("$price", "$quantity")
pie.Divide("$total", "$count")

// 字符串表达式
pie.Concat("$firstName", " ", "$lastName")
pie.SubstrCP("$text", 0, 5)
pie.ToUpper("$name")

// 日期表达式
pie.Year("$createdAt")
pie.Month("$createdAt")
pie.DayOfMonth("$createdAt")

// 逻辑表达式
pie.And(pie.Eq("$status", "active"), pie.Gt("$age", 18))
pie.Or(pie.Eq("$role", "admin"), pie.Eq("$role", "editor"))
pie.Cond(pie.Gt("$score", 90), "Excellent", "Good")

// 数组表达式
pie.SizeArray("$tags")
pie.ArrayElemAt("$items", 0)
pie.FilterArray("$items", "item", pie.GtExpr("$$item.price", 100))
```

## 下一步

- [最佳实践](/best-practices/) - 学习开发最佳实践
- [性能](/reference/performance/) - 学习性能优化
