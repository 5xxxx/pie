---
title: Query Builder
description: Pie provides rich query methods with chainable calls
---

# Smart Query Builder

Pie provides rich query methods with chainable calls:

```go
// Basic queries
session.Where("status", "active")
session.Where("age", pie.Gte("age", 18))
session.WhereIn("role", []string{"admin", "user"})
session.WhereBetween("age", 18, 65)
session.WhereNull("deleted_at")
session.WhereNotNull("email")

// Fuzzy queries
session.WhereLike("name", "%John%")
session.WhereStartsWith("email", "admin")
session.WhereEndsWith("domain", ".com")

// Date queries
session.WhereRecentDays("created_at", 7)
session.WhereMonth("created_at", time.Now().Month())
session.WhereYear("created_at", 2024)

// Complex conditions
session.Where("status", "active").
    OrWhere(func(q *pie.Query) {
        q.Where("role", "admin").WhereBetween("age", 30, 50)
    })
```

## Query Methods

### Comparison Operators

```go
// Equal
session.Where("status", "active")

// Not equal
session.Where("status", pie.Ne("status", "inactive"))

// Greater than
session.Where("age", pie.Gt("age", 18))

// Greater than or equal
session.Where("age", pie.Gte("age", 18))

// Less than
session.Where("age", pie.Lt("age", 65))

// Less than or equal
session.Where("age", pie.Lte("age", 65))

// Between
session.WhereBetween("age", 18, 65)

// Not between
session.WhereNotBetween("age", 18, 65)
```

### Array Operations

```go
// In array
session.WhereIn("role", []string{"admin", "user"})

// Not in array
session.WhereNotIn("role", []string{"guest"})

// Array contains
session.WhereArrayContains("tags", "golang")

// Array size
session.WhereArraySize("tags", 3)
```

### String Operations

```go
// Like (regex)
session.WhereLike("name", "%John%")

// Starts with
session.WhereStartsWith("email", "admin")

// Ends with
session.WhereEndsWith("domain", ".com")

// Regex
session.WhereRegex("name", "^[A-Z]")
```

### Null Operations

```go
// Is null
session.WhereNull("deleted_at")

// Is not null
session.WhereNotNull("email")

// Is empty
session.WhereEmpty("description")

// Is not empty
session.WhereNotEmpty("description")
```

### Date Operations

```go
// Recent days
session.WhereRecentDays("created_at", 7)

// Recent hours
session.WhereRecentHours("updated_at", 24)

// Specific month
session.WhereMonth("created_at", time.January)

// Specific year
session.WhereYear("created_at", 2024)

// Date range
session.WhereDateBetween("created_at", startDate, endDate)
```

## Complex Queries

### Logical Operators

```go
// AND conditions
session.Where("status", "active").
    Where("role", "user").
    Where("age", pie.Gte("age", 18))

// OR conditions
session.Where("status", "active").
    OrWhere("role", "admin")

// Nested conditions
session.Where("status", "active").
    OrWhere(func(q *pie.Query) {
        q.Where("role", "admin").
            Where("verified", true)
    })
```

### Subqueries

```go
// Exists subquery
session.WhereExists(func(q *pie.Query) {
    q.Table("orders").
        Where("user_id", pie.Raw("users._id")).
        Where("status", "completed")
})

// Not exists subquery
session.WhereNotExists(func(q *pie.Query) {
    q.Table("orders").
        Where("user_id", pie.Raw("users._id"))
})
```

### Raw Queries

```go
// Raw MongoDB query
session.WhereRaw(bson.D{
    {"$or", []bson.D{
        {{"age", bson.D{{"$gte", 18}}}},
        {{"role", "admin"}},
    }},
})

// Raw aggregation
session.AggregateRaw(bson.A{
    bson.D{{"$match", bson.D{{"status", "active"}}}},
    bson.D{{"$group", bson.D{
        {"_id", "$role"},
        {"count", bson.D{{"$sum", 1}}},
    }}},
})
```

## Query Execution

### Find Operations

```go
// Find all
var users []User
err := session.Find(ctx, &users)

// Find first
var user User
err := session.First(ctx, &user)

// Find with limit
var users []User
err := session.Limit(10).Find(ctx, &users)

// Find with offset
var users []User
err := session.Offset(20).Find(ctx, &users)

// Find with order
var users []User
err := session.OrderBy("name").Find(ctx, &users)
```

### Count Operations

```go
// Count all
count, err := session.Count(ctx)

// Count with conditions
count, err := session.Where("status", "active").Count(ctx)

// Exists check
exists, err := session.Where("email", "test@example.com").Exists(ctx)
```

### Update Operations

```go
// Update one
result, err := session.
    Where("email", "test@example.com").
    Update(ctx, bson.D{{"$set", bson.D{{"status", "active"}}}})

// Update many
result, err := session.
    Where("role", "guest").
    UpdateMany(ctx, bson.D{{"$set", bson.D{{"role", "user"}}}})

// Upsert
result, err := session.
    Where("email", "test@example.com").
    Upsert(ctx, bson.D{{"$set", bson.D{{"name", "Test User"}}}})
```

### Delete Operations

```go
// Delete one
result, err := session.
    Where("email", "test@example.com").
    Delete(ctx)

// Delete many
result, err := session.
    Where("status", "inactive").
    DeleteMany(ctx)
```

## Performance Tips

### Index Usage

```go
// Use indexed fields for better performance
session.Where("email", "test@example.com") // email has index
session.Where("created_at", pie.Gte("created_at", time.Now().AddDate(0, -1, 0))) // created_at has index
```

### Query Optimization

```go
// Use projection to limit returned fields
var users []User
err := session.
    Select("name", "email"). // Only select needed fields
    Where("status", "active").
    Find(ctx, &users)

// Use limit to avoid large result sets
var users []User
err := session.
    Where("status", "active").
    Limit(100). // Limit result size
    Find(ctx, &users)
```

### Caching

```go
// Cache query results
var users []User
err := session.
    Where("status", "active").
    Cache("active_users", 5*time.Minute). // Cache for 5 minutes
    Find(ctx, &users)
```

## Next Steps

- [Struct Query](/core-features/struct-query/) - Convert HTTP params to queries
- [Pagination](/core-features/pagination/) - Implement pagination
- [Cursor Operations](/core-features/cursor/) - Use cursors for large datasets
- [Aggregation](/core-features/aggregation/) - Advanced aggregation queries
