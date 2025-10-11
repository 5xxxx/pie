---
title: Generics Guide
description: Deep dive into generics usage, benefits, and best practices in Pie
---

# Generics Guide

Pie is built on Go generics, providing a type-safe MongoDB operation experience. This guide covers generics usage, benefits, and best practices in Pie.

## What are Generics?

Generics are a feature introduced in Go 1.18 that allows writing reusable code while maintaining type safety. In Pie, generics ensure compile-time type checking, preventing runtime type errors.

## Basic Usage

### Session[T] - The Most Common Generic Type

`Session[T]` is the core generic type in Pie, providing type-safe database operations:

```go
// Create type-safe session
session := pie.Table[User](engine)

// Query operations return []User, not []interface{}
users, err := session.Find(context.Background())

// Insert operations accept *User type
user := &User{Name: "John Doe", Email: "john@example.com"}
result, err := session.Insert(context.Background(), user)

// Type-safe update operations
updateResult, err := session.
    Where("email", "john@example.com").
    Update(context.Background(), bson.D{{"$set", bson.D{{"age", 25}}}})
```

### Aggregate[T] - Generic Support for Aggregation

```go
// Create type-safe aggregation builder
aggregate := pie.NewAggregate[User](engine)

// Aggregation results return []User
var results []User
err := aggregate.
    MatchStage().Where("status", "active").Done().
    GroupStage().By("department", "$department").Count("total").Done().
    Execute(context.Background(), &results)
```

### Cursor[T] - Generic Support for Cursors

```go
// Create cursor
cursor := pie.NewCursor[User](context.Background(), mongoCursor)

// Iteration provides type-safe User objects
for cursor.Next(context.Background()) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        // Handle error
    }
    // user is User type, not interface{}
}
```

### Generic Functions

Pie provides several generic functions for creating type-safe operations:

```go
// Basic functions
session := pie.Table[User](engine)                    // Create session
aggregate := pie.NewAggregate[User](engine)           // Create aggregation

// Must versions (panic on failure)
session := pie.MustTable[User](engine)                // Must create session
aggregate := pie.MustAggregate[User](engine)          // Must create aggregation

// Using default engine
session, err := pie.TableWithDefault[User]()          // Create session with default engine
session := pie.MustTableWithDefault[User]()           // Must create session with default engine
```

## Advanced Topics

### Benefits of Type Safety

Generics provide compile-time type checking, preventing runtime type errors:

```go
// ❌ Traditional approach - runtime errors possible
var users []interface{}
err := session.Find(context.Background(), &users)
// Manual type assertion required
for _, u := range users {
    user := u.(User) // May panic
}

// ✅ Generic approach - compile-time type safety
var users []User
err := session.Find(context.Background(), &users)
// Direct usage, no type assertion needed
for _, user := range users {
    fmt.Println(user.Name) // Type safe
}
```

### Type Inference and Explicit Type Declaration

Go's type inference makes generics usage more concise:

```go
// Type inference
session := pie.Table[User](engine)  // T is inferred as User

// Explicit type declaration (usually not needed)
var session *pie.Session[User] = pie.Table[User](engine)
```

### Generics with Interfaces

```go
// Define common interface
type Model interface {
    GetID() string
    GetCreatedAt() time.Time
}

// Use constraints
func ProcessModel[T Model](session *pie.Session[T], ctx context.Context) error {
    var models []T
    err := session.Find(ctx, &models)
    if err != nil {
        return err
    }
    
    for _, model := range models {
        fmt.Printf("ID: %s, Created: %v\n", model.GetID(), model.GetCreatedAt())
    }
    return nil
}
```

### Generic Return Types

Pie provides various generic return types:

```go
// Result[T] - Generic result type
func GetUser[T any](session *pie.Session[T], id string) *pie.Result[T] {
    var user T
    err := session.FindByID(context.Background(), id, &user)
    return &pie.Result[T]{
        Data:  user,
        Error: err,
    }
}

// PaginationResult[T] - Pagination result
func GetUsersWithPagination[T any](session *pie.Session[T], page, size int) *pie.PaginationResult[T] {
    var users []T
    total, err := session.Paginate(context.Background(), page, size, &users)
    return &pie.PaginationResult[T]{
        Data:  users,
        Total: total,
        Page:  int64(page),
        Size:  int64(size),
        Error: err,
    }
}

// AggregateResult[T] - Aggregation result
func GetUserStats[T any](aggregate *pie.Aggregate[T]) *pie.AggregateResult[T] {
    var results []T
    err := aggregate.Execute(context.Background(), &results)
    return &pie.AggregateResult[T]{
        Data:  results,
        Error: err,
    }
}
```

### Generics in Transactions

```go
// TransactionWithResult[T] - Transaction with result
func TransferMoney[T any](tx *pie.Transaction, fromID, toID string, amount float64) *pie.TransactionResult[T] {
    return pie.TransactionWithResult[T](tx, context.Background(), func(ctx context.Context) (T, error) {
        // Execute operations within transaction
        fromSession := pie.Table[Account](tx.engine)
        toSession := pie.Table[Account](tx.engine)
        
        // Deduct from sender's balance
        err := fromSession.Where("id", fromID).Update(ctx, bson.D{{"$inc", bson.D{{"balance", -amount}}}})
        if err != nil {
            return zero, err
        }
        
        // Add to receiver's balance
        err = toSession.Where("id", toID).Update(ctx, bson.D{{"$inc", bson.D{{"balance", amount}}}})
        if err != nil {
            return zero, err
        }
        
        return zero, nil
    })
}
```

### Generics in Change Streams

```go
// ChangeStreamWatcher[T] - Change stream watcher
func WatchUserChanges[T any](engine *pie.Engine) {
    watcher := pie.NewWatcher[T](engine)
    
    watcher.OnChange(func(event *pie.ChangeEvent[T]) {
        switch event.OperationType {
        case "insert":
            fmt.Printf("New user created: %+v\n", event.FullDocument)
        case "update":
            fmt.Printf("User updated: %+v\n", event.FullDocument)
        case "delete":
            fmt.Printf("User deleted: %+v\n", event.DocumentKey)
        }
    })
    
    watcher.Start(context.Background())
}
```

## Best Practices

### When to Use Generics vs interface{}

```go
// ✅ Recommended: Use generics for type safety
func ProcessUsers[T any](session *pie.Session[T], ctx context.Context) ([]T, error) {
    var users []T
    err := session.Find(ctx, &users)
    return users, err
}

// ❌ Not recommended: Use interface{}, loses type safety
func ProcessUsersGeneric(session *pie.Session[interface{}], ctx context.Context) ([]interface{}, error) {
    var users []interface{}
    err := session.Find(ctx, &users)
    return users, err
}
```

### Best Practices for Model Definition

```go
// ✅ Recommended: Define clear models
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Age       int                `bson:"age"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// ✅ Recommended: Implement interface methods
func (u *User) GetID() string {
    return u.ID.Hex()
}

func (u *User) GetCreatedAt() time.Time {
    return u.CreatedAt
}
```

### Techniques to Avoid Type Conversion

```go
// ✅ Recommended: Use generics to avoid type conversion
func GetUserByEmail[T any](session *pie.Session[T], email string) (*T, error) {
    var user T
    err := session.Where("email", email).First(context.Background(), &user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ❌ Not recommended: Requires type conversion
func GetUserByEmailGeneric(session *pie.Session[interface{}], email string) (interface{}, error) {
    var user interface{}
    err := session.Where("email", email).First(context.Background(), &user)
    if err != nil {
        return nil, err
    }
    // Type assertion required
    return user.(User), nil
}
```

### Performance Considerations

Generics in Go provide zero-cost abstraction, with compiled code performing the same as hand-written code:

```go
// Generic code
func FindByID[T any](session *pie.Session[T], id string) (*T, error) {
    var result T
    err := session.FindByID(context.Background(), id, &result)
    return &result, err
}

// Compiles to equivalent non-generic hand-written code
func FindUserByID(session *pie.Session[User], id string) (*User, error) {
    var result User
    err := session.FindByID(context.Background(), id, &result)
    return &result, err
}
```

## Common Issues

### Handling Generic Type Mismatch Errors

```go
// Issue: Type mismatch
var users []User
err := session.Find(context.Background(), &users) // ✅ Correct

// ❌ Error: Type mismatch
var users []string
err := session.Find(context.Background(), &users) // Compile error

// Solution: Ensure type consistency
var users []User
err := session.Find(context.Background(), &users)
```

### Converting Between Different Generic Types

```go
// If conversion is needed, use intermediate variables
func ConvertUsersToDTOs(users []User) []UserDTO {
    var dtos []UserDTO
    for _, user := range users {
        dtos = append(dtos, UserDTO{
            ID:    user.ID.Hex(),
            Name:  user.Name,
            Email: user.Email,
        })
    }
    return dtos
}

// Or use generic conversion function
func ConvertSlice[T, U any](slice []T, converter func(T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = converter(item)
    }
    return result
}
```

### Generics with BSON Serialization

```go
// BSON tags work with generics
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Age       int                `bson:"age"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// Generics ensure type safety, BSON tags handle serialization
session := pie.Table[User](engine)
user := &User{Name: "John Doe", Email: "john@example.com"}
result, err := session.Insert(context.Background(), user)
```

### Common Compilation Errors and Solutions

```go
// Error 1: Missing type parameter
// ❌ Compilation error
session := pie.Table(engine)

// ✅ Correct
session := pie.Table[User](engine)

// Error 2: Type mismatch
// ❌ Compilation error
var users []string
err := session.Find(context.Background(), &users)

// ✅ Correct
var users []User
err := session.Find(context.Background(), &users)

// Error 3: Generic constraint not satisfied
// ❌ If User doesn't implement Model interface
func ProcessModel[T Model](session *pie.Session[T], ctx context.Context) error {
    // ...
}

// ✅ Ensure User implements Model interface
type User struct {
    // ...
}

func (u *User) GetID() string { return u.ID.Hex() }
func (u *User) GetCreatedAt() time.Time { return u.CreatedAt }
```

## Summary

Generics are a core feature of Pie, providing:

1. **Type Safety**: Compile-time checking, preventing runtime errors
2. **Code Reuse**: One set of code handles multiple types
3. **Performance**: Zero-cost abstraction, no runtime overhead
4. **Developer Experience**: Better IDE support and code completion

By using generics properly, you can write safer, more efficient, and more maintainable MongoDB operation code.

## Next Steps

- [Query Builder](/core-features/query-builder/) - Learn about query functionality
- [Aggregation](/core-features/aggregation/) - Learn about data aggregation
- [Transactions](/core-features/transactions/) - Master transaction operations
- [Cache Support](/advanced/cache/) - Improve query performance
