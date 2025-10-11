---
title: Bulk Operations
description: Pie provides efficient bulk write operations
---

# Bulk Operations

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

## Bulk Write Types

### Insert Operations

```go
// Insert single document
bulkWrite.InsertOne(&User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   30,
})

// Insert multiple documents
users := []User{
    {Name: "Alice", Email: "alice@example.com", Age: 25},
    {Name: "Bob", Email: "bob@example.com", Age: 35},
    {Name: "Charlie", Email: "charlie@example.com", Age: 28},
}

for _, user := range users {
    bulkWrite.InsertOne(&user)
}
```

### Update Operations

```go
// Update one document
bulkWrite.UpdateOne(
    bson.D{{"email", "john@example.com"}},
    bson.D{{"$set", bson.D{{"age", 31}}}},
)

// Update many documents
bulkWrite.UpdateMany(
    bson.D{{"status", "pending"}},
    bson.D{{"$set", bson.D{{"status", "active"}}}},
)

// Upsert operation
bulkWrite.Upsert(
    bson.D{{"email", "newuser@example.com"}},
    bson.D{{"$set", bson.D{{"name", "New User"}}}},
)
```

### Delete Operations

```go
// Delete one document
bulkWrite.DeleteOne(bson.D{{"email", "old@example.com"}})

// Delete many documents
bulkWrite.DeleteMany(bson.D{{"status", "inactive"}})
```

## Execution Methods

### Ordered Execution

```go
// Execute operations in order
result, err := bulkWrite.ExecuteOrdered(ctx)

if err != nil {
    return err
}

fmt.Printf("Inserted: %d\n", result.InsertedCount)
fmt.Printf("Updated: %d\n", result.ModifiedCount)
fmt.Printf("Deleted: %d\n", result.DeletedCount)
fmt.Printf("Upserted: %d\n", result.UpsertedCount)
```

### Unordered Execution

```go
// Execute operations in parallel (faster)
result, err := bulkWrite.ExecuteUnordered(ctx)

if err != nil {
    return err
}

fmt.Printf("Inserted: %d\n", result.InsertedCount)
fmt.Printf("Updated: %d\n", result.ModifiedCount)
fmt.Printf("Deleted: %d\n", result.DeletedCount)
fmt.Printf("Upserted: %d\n", result.UpsertedCount)
```

## Real-world Examples

### User Import

```go
func ImportUsers(ctx context.Context, users []User) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    for _, user := range users {
        // Validate user
        if err := validateUser(&user); err != nil {
            log.Printf("Invalid user: %v", err)
            continue
        }
        
        // Insert user
        bulkWrite.InsertOne(&user)
    }
    
    // Execute bulk operation
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        return err
    }
    
    log.Printf("Imported %d users", result.InsertedCount)
    return nil
}

func validateUser(user *User) error {
    if user.Name == "" {
        return errors.New("name is required")
    }
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !isValidEmail(user.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
```

### Data Migration

```go
func MigrateUserData(ctx context.Context) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // Get all users
    cursor, err := session.FindCursor(ctx)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        // Update user data
        bulkWrite.UpdateOne(
            bson.D{{"_id", user.ID}},
            bson.D{{"$set", bson.D{
                {"email", strings.ToLower(user.Email)},
                {"name", strings.Title(user.Name)},
                {"migrated", true},
                {"migrated_at", time.Now()},
            }}},
        )
    }
    
    // Execute migration
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        return err
    }
    
    log.Printf("Migrated %d users", result.ModifiedCount)
    return nil
}
```

### Batch Status Update

```go
func UpdateUserStatus(ctx context.Context, userIDs []bson.ObjectID, status string) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    for _, userID := range userIDs {
        bulkWrite.UpdateOne(
            bson.D{{"_id", userID}},
            bson.D{{"$set", bson.D{
                {"status", status},
                {"updated_at", time.Now()},
            }}},
        )
    }
    
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        return err
    }
    
    log.Printf("Updated %d users to status: %s", result.ModifiedCount, status)
    return nil
}
```

### Data Cleanup

```go
func CleanupInactiveUsers(ctx context.Context) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // Delete users inactive for more than 1 year
    cutoffDate := time.Now().AddDate(-1, 0, 0)
    
    bulkWrite.DeleteMany(bson.D{
        {"status", "inactive"},
        {"last_login", bson.D{{"$lt", cutoffDate}}},
    })
    
    // Delete users with invalid email
    bulkWrite.DeleteMany(bson.D{
        {"email", bson.D{{"$regex", "^[^@]+$"}}},
    })
    
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        return err
    }
    
    log.Printf("Cleaned up %d users", result.DeletedCount)
    return nil
}
```

## Error Handling

### Handle Partial Failures

```go
func BulkInsertWithErrorHandling(ctx context.Context, users []User) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    for _, user := range users {
        bulkWrite.InsertOne(&user)
    }
    
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        // Check if it's a bulk write error
        if bulkErr, ok := err.(pie.BulkWriteError); ok {
            log.Printf("Bulk operation partially failed:")
            log.Printf("Inserted: %d", result.InsertedCount)
            log.Printf("Failed: %d", len(bulkErr.WriteErrors))
            
            for _, writeErr := range bulkErr.WriteErrors {
                log.Printf("Error at index %d: %s", writeErr.Index, writeErr.Message)
            }
            
            // Return error if too many failures
            if len(bulkErr.WriteErrors) > len(users)/2 {
                return err
            }
        } else {
            return err
        }
    }
    
    return nil
}
```

### Retry Logic

```go
func BulkInsertWithRetry(ctx context.Context, users []User) error {
    maxRetries := 3
    retryDelay := 1 * time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := bulkInsert(ctx, users)
        if err == nil {
            return nil
        }
        
        log.Printf("Bulk insert failed (attempt %d/%d): %v", i+1, maxRetries, err)
        
        if i < maxRetries-1 {
            time.Sleep(retryDelay)
            retryDelay *= 2 // Exponential backoff
        }
    }
    
    return fmt.Errorf("bulk insert failed after %d attempts", maxRetries)
}

func bulkInsert(ctx context.Context, users []User) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    for _, user := range users {
        bulkWrite.InsertOne(&user)
    }
    
    _, err := bulkWrite.ExecuteOrdered(ctx)
    return err
}
```

## Performance Optimization

### Batch Size Optimization

```go
func OptimizedBulkInsert(ctx context.Context, users []User) error {
    batchSize := 1000
    totalUsers := len(users)
    
    for i := 0; i < totalUsers; i += batchSize {
        end := i + batchSize
        if end > totalUsers {
            end = totalUsers
        }
        
        batch := users[i:end]
        if err := insertBatch(ctx, batch); err != nil {
            return err
        }
        
        log.Printf("Processed %d/%d users", end, totalUsers)
    }
    
    return nil
}

func insertBatch(ctx context.Context, users []User) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    for _, user := range users {
        bulkWrite.InsertOne(&user)
    }
    
    _, err := bulkWrite.ExecuteOrdered(ctx)
    return err
}
```

### Use Unordered Execution

```go
// Use unordered execution for better performance when order doesn't matter
func FastBulkUpdate(ctx context.Context, updates []UpdateOperation) error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    for _, update := range updates {
        bulkWrite.UpdateOne(update.Filter, update.Update)
    }
    
    // Use unordered execution for better performance
    _, err := bulkWrite.ExecuteUnordered(ctx)
    return err
}
```

## Best Practices

### 1. Use Appropriate Batch Sizes

```go
// Good: Use appropriate batch size
const OptimalBatchSize = 1000

func BulkInsert(ctx context.Context, users []User) error {
    for i := 0; i < len(users); i += OptimalBatchSize {
        end := i + OptimalBatchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        if err := insertBatch(ctx, batch); err != nil {
            return err
        }
    }
    
    return nil
}

// Bad: Processing all at once
func BadBulkInsert(ctx context.Context, users []User) error {
    bulkWrite := pie.NewBulkWrite[User](engine)
    
    for _, user := range users {
        bulkWrite.InsertOne(&user)
    }
    
    _, err := bulkWrite.ExecuteOrdered(ctx)
    return err
}
```

### 2. Handle Errors Appropriately

```go
// Good: Handle errors properly
func BulkInsertWithErrorHandling(ctx context.Context, users []User) error {
    bulkWrite := pie.NewBulkWrite[User](engine)
    
    for _, user := range users {
        bulkWrite.InsertOne(&user)
    }
    
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        if bulkErr, ok := err.(pie.BulkWriteError); ok {
            log.Printf("Partial failure: %d inserted, %d failed", 
                result.InsertedCount, len(bulkErr.WriteErrors))
        }
        return err
    }
    
    return nil
}
```

### 3. Use Transactions for Critical Operations

```go
// Good: Use transactions for critical operations
func CriticalBulkUpdate(ctx context.Context, updates []UpdateOperation) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        bulkWrite := pie.NewBulkWrite[User](engine)
        
        for _, update := range updates {
            bulkWrite.UpdateOne(update.Filter, update.Update)
        }
        
        _, err := bulkWrite.ExecuteOrdered(txCtx)
        return err
    })
}
```

## Next Steps

- [Aggregation](/core-features/aggregation/) - Learn about aggregation queries
- [Transactions](/core-features/transactions/) - Use transactions
- [Performance](/reference/performance/) - Learn performance optimization
- [Best Practices](/best-practices/) - Learn development best practices
