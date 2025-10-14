---
title: Performance
description: Performance optimization strategies for MongoDB database operations
---

# Performance

Pie provides multiple strategies to help you optimize MongoDB database operation performance.

## 1. Proper Index Usage

Ensure appropriate indexes are created on fields frequently used in `Where`, `OrderBy`, and `Lookup` operations.

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Email     string        `bson:"email" pie:"index:,unique"` // Unique index
    Age       int           `bson:"age" pie:"index:age_idx"`   // Regular index
    CreatedAt time.Time     `bson:"created_at" pie:"index:created_at_desc,order=-1"` // Descending index
}

// Enable automatic index creation
engine.WithAutoIndex(true)
```

## 2. Bulk Operations

For large insert, update, or delete operations, use bulk writes to significantly reduce network round trips and improve performance.

```go
bulkWrite := pie.NewBulkWrite[User](engine).CollectionForStruct(User{})
for _, user := range users {
    bulkWrite.InsertOne(user)
}
_, err := bulkWrite.ExecuteOrdered(ctx)
```

## 3. Caching

Leverage Pie's built-in caching mechanism (memory cache or Redis cache) to cache infrequently changing but frequently read data.

```go
// Enable Redis cache
cfg := &pie.RedisCacheConfig{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    PoolSize: 10,
}
redisCache, err := pie.NewRedisCache(cfg)
if err != nil {
    log.Fatal(err)
}
engine.UseCache(redisCache)

// Use cache in queries
users, err := session.WithCache(5 * time.Minute).Cache("active_users").Find(ctx)

```

## 4. Projection

Only query the fields you need, avoiding unnecessary data transmission.

```go
type UserProjection struct {
    Name  string `bson:"name"`
    Email string `bson:"email"`
}

users, err := session.Project("name", "email").Find(ctx)
```

## 5. Limit and Skip

In pagination queries, avoid using large `Skip` values as they cause MongoDB to scan many documents. Consider using cursors or range-based pagination.

```go
// Avoid large Skip
// session.Skip(100000).Limit(10).Find(ctx)

// Use cursor
cursor, err := session.Where("status", "active").FindCursor(ctx)
// ... process cursor ...
```

## 6. Aggregation Pipeline Optimization

In aggregation queries, try to place `$match` stages early in the pipeline to filter out unnecessary documents as soon as possible.

```go
aggregate := pie.NewAggregate[User](engine).CollectionForStruct(User{})
result, err := aggregate.
    MatchStage().Where("status", "active").Done(). // Filter first
    GroupStage().By("role", "$role").Count("total").Done().
    Exec(ctx)
```

## 7. Connection Pool Configuration

Properly configure MongoDB client connection pool size to avoid too many or too few connections.

```go
engine.WithMaxPoolSize(100).WithMinPoolSize(10)
```

## 8. Query Logging and Monitoring

Enable query logging and regularly analyze query performance.

```go
engine.OnQuery(func(event *pie.QueryEvent) {
    log.Printf("Query: %s, Duration: %v", event.Query, event.Duration)
})
```

## Next Steps

- [Best Practices](/best-practices/) - Learn development best practices
- [Configuration](/reference/configuration/) - Learn about configuration options
