---
title: Getting Started
description: Learn how to get started with Pie quickly
---

# Getting Started

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

## Next Steps

- [Generics Guide](/core-features/generics/) - Learn about generics usage and benefits
- [Query Builder](/core-features/query-builder/) - Learn about query methods
- [Struct Query](/core-features/struct-query/) - Convert HTTP params to queries
- [Pagination](/core-features/pagination/) - Implement pagination
- [Transactions](/core-features/transactions/) - Use transactions
