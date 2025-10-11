---
title: Advanced Aggregation
description: Comprehensive aggregation stage builders and 100+ expression functions
---

# Advanced Aggregation

Pie provides comprehensive aggregation stage builders and 100+ expression functions for complex data processing.

## Stage Builders

```go
aggregate := pie.NewAggregate[User](engine).CollectionForStruct(User{})

// $match stage
aggregate.MatchStage().Where("status", "active").Between("age", 18, 65).Done()

// $group stage
aggregate.GroupStage().
    By("status", "$status").
    Count("total").
    Avg("avgAge", "$age").
    Done()

// $project stage
aggregate.ProjectStage().Include("name", "email").Exclude("_id").Done()

// $lookup stage (join)
aggregate.LookupStage("orders", "_id", "user_id", "userOrders").Done()

// $facet stage (faceted search)
aggregate.FacetStage().
    Facet("activeUsers", bson.M{"$match": bson.M{"active": true}}).
    Facet("ageGroups", bson.M{"$bucket": bson.M{"groupBy": "$age", "boundaries": []int{0, 20, 40, 60, 100}}}).
    Done()
```

## Expression Functions

Pie provides 100+ aggregation expression functions:

```go
// Arithmetic expressions
pie.Add(1, 2, "$field")
pie.Multiply("$price", "$quantity")
pie.Divide("$total", "$count")

// String expressions
pie.Concat("$firstName", " ", "$lastName")
pie.SubstrCP("$text", 0, 5)
pie.ToUpper("$name")

// Date expressions
pie.Year("$createdAt")
pie.Month("$createdAt")
pie.DayOfMonth("$createdAt")

// Logical expressions
pie.And(pie.Eq("$status", "active"), pie.Gt("$age", 18))
pie.Or(pie.Eq("$role", "admin"), pie.Eq("$role", "editor"))
pie.Cond(pie.Gt("$score", 90), "Excellent", "Good")

// Array expressions
pie.SizeArray("$tags")
pie.ArrayElemAt("$items", 0)
pie.FilterArray("$items", "item", pie.GtExpr("$$item.price", 100))
```

## Next Steps

- [Best Practices](/best-practices/) - Learn development best practices
- [Performance](/reference/performance/) - Learn performance optimization
