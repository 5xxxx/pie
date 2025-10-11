---
title: Cursor Operations
description: Pie provides flexible cursor operations
---

# Cursor Operations

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

## Cursor Methods

### Basic Iteration

```go
// Get cursor
cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").
    FindCursor(ctx)

if err != nil {
    return err
}
defer cursor.Close(ctx)

// Iterate through results
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        log.Printf("Failed to decode user: %v", err)
        continue
    }
    
    // Process user
    processUser(&user)
}
```

### Batch Processing

```go
// Process in batches
cursor, err := session.
    Where("status", "active").
    FindCursor(ctx)

if err != nil {
    return err
}
defer cursor.Close(ctx)

batchSize := 100
batch := make([]User, 0, batchSize)

for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    
    batch = append(batch, user)
    
    if len(batch) >= batchSize {
        // Process batch
        if err := processBatch(batch); err != nil {
            return err
        }
        batch = batch[:0] // Reset batch
    }
}

// Process remaining items
if len(batch) > 0 {
    if err := processBatch(batch); err != nil {
        return err
    }
}
```

### Error Handling

```go
cursor, err := session.
    Where("status", "active").
    FindCursor(ctx)

if err != nil {
    return err
}
defer cursor.Close(ctx)

for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        // Log error but continue
        log.Printf("Failed to decode user: %v", err)
        continue
    }
    
    // Process user
    if err := processUser(&user); err != nil {
        // Handle processing error
        log.Printf("Failed to process user %s: %v", user.ID, err)
        continue
    }
}

// Check for cursor errors
if err := cursor.Err(); err != nil {
    return fmt.Errorf("cursor error: %w", err)
}
```

## Advanced Cursor Operations

### Cursor with Projection

```go
// Get cursor with specific fields
cursor, err := session.
    Select("name", "email", "status").
    Where("status", "active").
    FindCursor(ctx)

if err != nil {
    return err
}
defer cursor.Close(ctx)

for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    
    // Only name, email, and status are populated
    fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
}
```

### Cursor with Sorting

```go
// Get cursor with sorting
cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").
    FindCursor(ctx)

if err != nil {
    return err
}
defer cursor.Close(ctx)

// Process in chronological order
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    
    processUser(&user)
}
```

### Cursor with Limit

```go
// Get cursor with limit
cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").
    Limit(1000).
    FindCursor(ctx)

if err != nil {
    return err
}
defer cursor.Close(ctx)

// Process up to 1000 users
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    
    processUser(&user)
}
```

## Real-world Examples

### Data Migration

```go
func MigrateUsers(ctx context.Context) error {
    cursor, err := session.
        Where("status", "active").
        FindCursor(ctx)
    
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    batchSize := 100
    batch := make([]User, 0, batchSize)
    
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        // Transform user data
        user.Email = strings.ToLower(user.Email)
        user.Name = strings.Title(user.Name)
        
        batch = append(batch, user)
        
        if len(batch) >= batchSize {
            if err := updateBatch(ctx, batch); err != nil {
                return err
            }
            batch = batch[:0]
        }
    }
    
    // Process remaining batch
    if len(batch) > 0 {
        if err := updateBatch(ctx, batch); err != nil {
            return err
        }
    }
    
    return nil
}

func updateBatch(ctx context.Context, users []User) error {
    for _, user := range users {
        _, err := session.
            Where("_id", user.ID).
            Update(ctx, bson.D{{"$set", bson.D{
                {"email", user.Email},
                {"name", user.Name},
            }}})
        if err != nil {
            return err
        }
    }
    return nil
}
```

### Data Export

```go
func ExportUsers(ctx context.Context, writer io.Writer) error {
    cursor, err := session.
        Where("status", "active").
        OrderBy("created_at").
        FindCursor(ctx)
    
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    // Write CSV header
    writer.Write([]byte("ID,Name,Email,Status,CreatedAt\n"))
    
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        // Write CSV row
        row := fmt.Sprintf("%s,%s,%s,%s,%s\n",
            user.ID.Hex(),
            user.Name,
            user.Email,
            user.Status,
            user.CreatedAt.Format(time.RFC3339),
        )
        writer.Write([]byte(row))
    }
    
    return nil
}
```

### Data Processing

```go
func ProcessUsers(ctx context.Context) error {
    cursor, err := session.
        Where("status", "active").
        FindCursor(ctx)
    
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        // Process user asynchronously
        go func(u User) {
            if err := processUserAsync(u); err != nil {
                log.Printf("Failed to process user %s: %v", u.ID, err)
            }
        }(user)
    }
    
    return nil
}

func processUserAsync(user User) error {
    // Simulate processing
    time.Sleep(100 * time.Millisecond)
    
    // Update user status
    _, err := session.
        Where("_id", user.ID).
        Update(context.Background(), bson.D{{"$set", bson.D{
            {"processed", true},
            {"processed_at", time.Now()},
        }}})
    
    return err
}
```

## Performance Tips

### Use Projection

```go
// Only select needed fields
cursor, err := session.
    Select("name", "email"). // Only select needed fields
    Where("status", "active").
    FindCursor(ctx)
```

### Use Indexes

```go
// Ensure cursor fields are indexed
func createCursorIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // Index for cursor queries
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    })
    
    return err
}
```

### Batch Processing

```go
// Process in batches to avoid memory issues
func processLargeDataset(ctx context.Context) error {
    cursor, err := session.
        Where("status", "active").
        FindCursor(ctx)
    
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    batchSize := 1000
    batch := make([]User, 0, batchSize)
    
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        batch = append(batch, user)
        
        if len(batch) >= batchSize {
            if err := processBatch(batch); err != nil {
                return err
            }
            batch = batch[:0]
            
            // Force garbage collection
            runtime.GC()
        }
    }
    
    return nil
}
```

## Best Practices

### 1. Always Close Cursors

```go
// Good: Always close cursors
cursor, err := session.FindCursor(ctx)
if err != nil {
    return err
}
defer cursor.Close(ctx)

// Bad: Forgetting to close cursors
cursor, err := session.FindCursor(ctx)
// Missing defer cursor.Close(ctx)
```

### 2. Handle Errors Properly

```go
// Good: Handle cursor errors
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        log.Printf("Failed to decode user: %v", err)
        continue
    }
    
    processUser(&user)
}

if err := cursor.Err(); err != nil {
    return fmt.Errorf("cursor error: %w", err)
}

// Bad: Ignoring errors
for cursor.Next(ctx) {
    var user User
    cursor.Decode(&user) // Ignoring error
    processUser(&user)
}
```

### 3. Use Appropriate Batch Sizes

```go
// Good: Use appropriate batch size
const BatchSize = 1000

func processUsers(ctx context.Context) error {
    cursor, err := session.FindCursor(ctx)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    batch := make([]User, 0, BatchSize)
    
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        batch = append(batch, user)
        
        if len(batch) >= BatchSize {
            if err := processBatch(batch); err != nil {
                return err
            }
            batch = batch[:0]
        }
    }
    
    return nil
}
```

## Next Steps

- [Bulk Operations](/core-features/bulk-operations/) - Learn about bulk operations
- [Query Builder](/core-features/query-builder/) - Learn about query methods
- [Transactions](/core-features/transactions/) - Use transactions
- [Performance](/reference/performance/) - Learn performance optimization
