---
title: Aggregation
description: Pie provides powerful aggregation framework with stage builders and expression functions for complex data processing
---

# Aggregation Queries

Pie provides a powerful aggregation framework with stage builders and expression functions for complex data processing.

## Basic Aggregation

```go
// Create aggregation operation
aggregate := pie.NewAggregate[User](engine).
    CollectionForStruct(User{})

// Use stage builders to build aggregation pipeline
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

// Process results
for _, item := range result.Data {
    // item is bson.M type
}
```

## Advanced Aggregation with Expressions

```go
// Complex aggregation with expressions
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

## Join Operations

```go
// Join with orders collection
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

## Facet Analysis

```go
// Multi-dimensional analysis using facets
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

## Next Steps

- [Transactions](/core-features/transactions/) - Use transactions
- [Advanced Aggregation](/advanced/aggregation-advanced/) - Learn advanced aggregation features
- [Performance](/reference/performance/) - Learn performance optimization
