# Pie - MongoDB ORM Framework

Pie is a modern, type-safe MongoDB ORM framework for Go that provides a rich, high-performance database operation experience.

## Features

- **Type Safety**: Built on Go generics with compile-time type checking
- **Smart Query Builder**: Intuitive chainable API with support for complex query conditions
- **Struct Query**: Direct conversion of HTTP request parameters to query conditions
- **Cache Support**: Built-in multi-level caching with Redis and memory cache support
- **Hook System**: Complete lifecycle hook support
- **Transaction Management**: Simple and easy-to-use transaction operations
- **Index Management**: Automated index creation and management
- **Advanced Aggregation**: Comprehensive stage builders with 100+ expression functions
- **Change Streams**: Real-time data change monitoring
- **Soft Delete**: Built-in soft delete functionality
- **Pagination**: Efficient pagination implementation
- **Query Logging**: Detailed query logging and performance monitoring
- **Bulk Operations**: High-performance bulk write operations

## Installation

```bash
go get github.com/5xxxx/pie
```

## Quick Start

### 1. Connect to Database

```go
package main

import (
    "context"
    "log"
    
    "github.com/5xxxx/pie"
)

func main() {
    // Create engine
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
    
    // Use engine...
}
```

### 2. Define Models

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    Age       int           `bson:"age"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

// Hook methods
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

### 3. Basic Operations

```go
// Create type-safe session
session := pie.Table[User](engine)

// Insert document
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   25,
}
result, err := session.Insert(context.Background(), user)

// Query documents
users, err := session.
    Where("age", pie.Gte("age", 18)).
    OrderBy("name").
    Limit(10).
    Find(context.Background())

// Update document
updateResult, err := session.
    Where("email", "john@example.com").
    Update(context.Background(), bson.D{{"$set", bson.D{{"age", 26}}}})

// Delete document
deleteResult, err := session.
    Where("email", "john@example.com").
    Delete(context.Background())
```

## Core Features

### 1. Smart Query Builder

Pie provides rich query methods with chainable API:

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

### 2. Struct Query

Convert HTTP request parameters directly into query conditions:

```go
type UserQuery struct {
    Name     string   `pie:"name,like,omitempty" json:"name"`
    Email    string   `pie:"email,omitempty" json:"email"`
    MinAge   int      `pie:"age,gte,omitempty" json:"min_age"`
    MaxAge   int      `pie:"age,lte,omitempty" json:"max_age"`
    Status   []string `pie:"status,in,omitempty" json:"status"`
    Role     string   `pie:"role,omitempty" json:"role"`
}

// HTTP request: GET /users?name=John&min_age=20&max_age=40&status=active,pending
query := UserQuery{
    Name:   "John",
    MinAge: 20,
    MaxAge: 40,
    Status: []string{"active", "pending"},
}

var users []User
err := session.WhereStruct(query).Find(ctx, &users)
```

### 3. Convenience Methods

```go
// Find by ID
user, err := session.FindByID(ctx, userID)

// Check existence
exists, err := session.Where("email", "test@example.com").Exists(ctx)

// Quick count
count, err := session.Where("status", "active").QuickCount(ctx)

// Find or create
user, isNew, err := session.
    Where("email", "test@example.com").
    FirstOrCreate(ctx, &User{Email: "test@example.com", Name: "Test User"})

// Find or fail
user, err := session.
    Where("email", "test@example.com").
    FirstOrFail(ctx)
```

### 4. Pagination

```go
// Full pagination (with total count)
result, err := session.
    Where("status", "active").
    Paginate(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
// result.Total, result.TotalPages, result.HasNext, result.HasPrev

// Simple pagination (without total count, faster)
simpleResult, err := session.
    PaginateSimple(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
```

### 5. Cursor Operations

```go
// Get cursor
cursor, err := session.
    Where("status", "active").
    FindCursor(ctx)

// Method 1: Use Next() and Decode()
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    // Process user
}

// Method 2: Use All() to get all at once
users, err := cursor.All(ctx)

// Method 3: Use Iterate() for iteration
cursor.Iterate(ctx, func(user *User) error {
    // Process user
    return nil
})

// Method 4: Use Take() to get first N
topUsers, err := cursor.Take(ctx, 5)

// Method 5: Use First() to get first one
firstUser, err := cursor.First(ctx)

cursor.Close(ctx)
```

### 6. Bulk Operations

```go
// Create bulk write operation
bulkWrite := pie.NewBulkWrite[User](engine).
    CollectionForStruct(User{})

// Insert multiple documents
bulkWrite.InsertOne(&User{Name: "User 1"})
bulkWrite.InsertOne(&User{Name: "User 2"})

// Update operations
bulkWrite.UpdateOne(
    bson.D{{"email", "old@example.com"}},
    bson.D{{"$set", bson.D{{"email", "new@example.com"}}}},
)

// Bulk update
bulkWrite.UpdateMany(
    bson.D{{"status", "inactive"}},
    bson.D{{"$set", bson.D{{"status", "active"}}}},
)

// Delete operations
bulkWrite.DeleteMany(bson.D{{"age", bson.D{{"$lt", 18}}}})

// Execute bulk operation
result, err := bulkWrite.ExecuteOrdered(ctx)
```

### 7. Aggregation Queries

Pie provides a powerful aggregation framework with stage builders and expression functions for complex data processing.

#### Basic Aggregation

```go
// Create aggregation operation
aggregate := pie.NewAggregate[User](engine).
    CollectionForStruct(User{})

// Build aggregation pipeline using stage builders
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

#### Advanced Aggregation with Expressions

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

#### Lookup and Join Operations

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

#### Facet Analysis

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

### 8. Transaction Management

```go
// Execute transaction with engine
err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
    // Execute operations in transaction
    session := pie.Table[User](engine)
    
    // Insert user
    _, err := session.Insert(txCtx, &User{Name: "Transaction User"})
    if err != nil {
        return err
    }
    
    // Update other collections
    orderSession := pie.Table[Order](engine)
    _, err = orderSession.Insert(txCtx, &Order{UserID: userID})
    return err
})

// Use transaction manager
tx := pie.MustTransaction(engine)
err = tx.Execute(ctx, func(txCtx context.Context) error {
    // Transaction operations
    return nil
})
```

### 9. Cache Support

```go
// Enable memory cache
engine.UseCache(pie.NewMemoryCache(), &pie.CacheConfig{
    TTL: 5 * time.Minute,
})

// Enable Redis cache
redisCache := pie.NewRedisCache("localhost:6379", "", 0)
engine.UseCache(redisCache, &pie.CacheConfig{
    TTL: 10 * time.Minute,
})

// Enable two-level cache
engine.UseTwoLevelCache(
    pie.NewMemoryCache(),  // L1 cache
    redisCache,            // L2 cache
    &pie.TwoLevelCacheConfig{
        L1TTL: 1 * time.Minute,
        L2TTL: 10 * time.Minute,
    },
)

// Use cache in session
session := pie.Table[User](engine).
    WithCache(5 * time.Minute)

// Cache query results
var users []User
err := session.
    Where("status", "active").
    Cache("active_users").
    Find(ctx, &users)
```

### 10. Hook System

```go
type User struct {
    // Field definitions...
}

// Before create hook
func (u *User) BeforeCreate(ctx context.Context) error {
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return nil
}

// After create hook
func (u *User) AfterCreate(ctx context.Context) error {
    log.Printf("User %s created", u.Name)
    return nil
}

// Before update hook
func (u *User) BeforeUpdate(ctx context.Context) error {
    u.UpdatedAt = time.Now()
    return nil
}

// After update hook
func (u *User) AfterUpdate(ctx context.Context) error {
    log.Printf("User %s updated", u.Name)
    return nil
}

// Before delete hook
func (u *User) BeforeDelete(ctx context.Context) error {
    log.Printf("About to delete user %s", u.Name)
    return nil
}

// After delete hook
func (u *User) AfterDelete(ctx context.Context) error {
    log.Printf("User %s deleted", u.Name)
    return nil
}

// After find hook
func (u *User) AfterFind(ctx context.Context) error {
    // Handle post-find logic
    return nil
}
```

### 11. Soft Delete

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

// Soft delete
err := session.Where("email", "test@example.com").SoftDelete(ctx)

// Restore soft deleted
err := session.Where("email", "test@example.com").Restore(ctx)

// Force delete (physical delete)
err := session.Where("email", "test@example.com").ForceDelete(ctx)

// Queries automatically exclude soft deleted records
var users []User
err := session.Find(ctx, &users) // Automatically filters out records where deleted_at is not null
```

### 12. Index Management

```go
// Create index manager
indexes := pie.MustIndexes(engine)

// Create indexes for struct
err := indexes.CreateIndexes(ctx, User{})

// Manually create index
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"email", 1},
}, &options.IndexOptions{
    Unique: pie.Bool(true),
})

// Create compound index
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"status", 1},
    {"created_at", -1},
})

// Drop index
err := indexes.DropIndex(ctx, "users", "email_1")
```

### 13. Change Stream Monitoring

```go
// Create change stream watcher
watcher := pie.NewWatcher[User](engine)

// Watch collection changes
err := watcher.
    WatchCollection().
    Filter(bson.D{{"operationType", "insert"}}).
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("New user inserted: %s", change.FullDocument.Name)
        return nil
    })

// Watch database changes
dbWatcher := pie.NewDatabaseWatcher[User](engine)
err = dbWatcher.
    WatchDatabase().
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("Database change: %s", change.OperationType)
        return nil
    })
```

### 14. Query Scopes

```go
// Define scopes
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

// Use scopes
var users []User
err := session.
    Scopes(
        ActiveScope("status"),
        RecentScope("created_at", 30),
    ).
    Latest("created_at", 10).
    Find(ctx, &users)
```

### 15. Query Logging and Monitoring

```go
// Enable query logging
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithQueryLog(os.Stdout),
    pie.WithSlowQueryThreshold(50*time.Millisecond),
)

// Custom log format
engine.SetQueryLogFormatter(func(entry *pie.LogEntry) string {
    return fmt.Sprintf("[%s] %s %s - %v", 
        entry.Timestamp.Format("15:04:05"),
        entry.Operation,
        entry.Collection,
        entry.Duration,
    )
})

// Set slow query threshold
engine.SetSlowQueryThreshold(100 * time.Millisecond)
```

## Advanced Features

### 1. Aggregation Stage Builders

Pie provides comprehensive stage builders for MongoDB aggregation pipelines, making complex data processing intuitive and type-safe.

#### Match Stage Builder

```go
// Basic matching with chainable conditions
result, err := aggregate.
    MatchStage().
        Where("status", "active").
        Between("age", 18, 65).
        In("role", "admin", "user").
        Regex("name", "^A").
        Exists("email", true).
        Text("search term").
    Exec(ctx)

// Complex logical conditions
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

#### Group Stage Builder

```go
// Comprehensive grouping with multiple accumulators
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

#### Project Stage Builder

```go
// Field projection and computed fields
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

#### Lookup Stage Builder

```go
// Advanced lookup with pipeline
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

#### Unwind Stage Builder

```go
// Array unwinding with options
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

#### Facet Stage Builder

```go
// Multi-dimensional analysis
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

### 2. Aggregation Expression Functions

Pie provides 100+ expression functions for complex data transformations and calculations.

#### Date Expressions

```go
// Date manipulation and formatting
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

#### Arithmetic Expressions

```go
// Mathematical operations
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

#### String Expressions

```go
// String manipulation
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

#### Array Expressions

```go
// Array operations
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

#### Conditional Expressions

```go
// Conditional logic
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

### 3. Custom Name Mapping

```go
// Use snake case naming
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SnakeMapper{}),
)

// Use camel case naming
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.CamelMapper{}),
)

// Use same naming
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SameMapper{}),
)

// Custom mapper
type CustomMapper struct{}

func (m CustomMapper) TableName(structName string) string {
    return "t_" + strings.ToLower(structName)
}

func (m CustomMapper) FieldName(fieldName string) string {
    return strings.ToLower(fieldName)
}
```

### 2. Configuration Options

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

### 3. Error Handling

```go
// Check specific error types
if pie.IsDuplicateKeyError(err) {
    log.Println("Duplicate key error")
}

if pie.IsNotFoundError(err) {
    log.Println("Document not found")
}

if pie.IsTimeoutError(err) {
    log.Println("Operation timeout")
}

// Get error details
if mongoErr, ok := err.(pie.MongoError); ok {
    log.Printf("MongoDB error code: %d", mongoErr.Code)
    log.Printf("MongoDB error message: %s", mongoErr.Message)
}
```

## Performance Optimization

### 1. Connection Pool Configuration

```go
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMaxPoolSize(100),        // Maximum connections
    pie.WithMinPoolSize(5),          // Minimum connections
    pie.WithMaxIdleTime(30*time.Minute), // Maximum idle time
)
```

### 2. Query Optimization

```go
// Use projection to reduce data transfer
var users []User
err := session.
    Select("name", "email").  // Only select needed fields
    Find(ctx, &users)

// Use indexes to optimize queries
err := session.
    Where("email", "test@example.com").  // email field has index
    Find(ctx, &users)

// Use cursor for large datasets
cursor, err := session.FindCursor(ctx)
defer cursor.Close(ctx)

for cursor.Next(ctx) {
    var user User
    cursor.Decode(&user)
    // Process single user
}
```

### 3. Bulk Operation Optimization

```go
// Use bulk writes for better performance
bulkWrite := pie.NewBulkWrite[User](engine)
for _, user := range users {
    bulkWrite.InsertOne(user)
}
result, err := bulkWrite.ExecuteOrdered(ctx)
```

## Best Practices

### 1. Model Design

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

### 2. Error Handling

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    session := pie.Table[User](s.engine)
    
    // Check if email already exists
    exists, err := session.Where("email", user.Email).Exists(ctx)
    if err != nil {
        return fmt.Errorf("failed to check email existence: %w", err)
    }
    if exists {
        return errors.New("email already exists")
    }
    
    // Create user
    _, err = session.Insert(ctx, user)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

### 3. Transaction Usage

```go
func (s *UserService) TransferPoints(ctx context.Context, fromUserID, toUserID bson.ObjectID, points int) error {
    return s.engine.WithTransaction(ctx, func(txCtx context.Context) error {
        userSession := pie.Table[User](s.engine)
        
        // Deduct points from sender
        _, err := userSession.
            Where("_id", pie.ID(fromUserID)).
            Update(txCtx, bson.D{{"$inc", bson.D{{"points", -points}}}})
        if err != nil {
            return fmt.Errorf("failed to deduct points: %w", err)
        }
        
        // Add points to receiver
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

## License

MIT License

