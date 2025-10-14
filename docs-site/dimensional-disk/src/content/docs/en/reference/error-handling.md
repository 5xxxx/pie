---
title: Error Handling
description: Comprehensive error handling mechanism for database operations
---

# Error Handling

Pie provides a comprehensive error handling mechanism to help developers identify and resolve database operation issues.

## Common Error Types

Pie defines some common error types for easy error checking:

- `pie.ErrNotFound`: Document not found
- `pie.ErrDuplicateKey`: Unique index conflict
- `pie.ErrValidation`: Data validation failed
- `pie.ErrInvalidID`: Invalid ID format
- `pie.ErrTransactionAborted`: Transaction aborted

```go
_, err := session.Where("email", "nonexistent@example.com").FindOne(ctx)
if err != nil {
    if errors.Is(err, pie.ErrNotFound) {
        log.Println("User not found")
    } else {
        log.Printf("Error finding user: %v", err)
    }
}
```

## Error Wrapping and Unwrapping

Pie wraps underlying MongoDB driver errors, and you can use `errors.Is` or `errors.As` to check original errors.

```go
_, err := session.Insert(ctx, &User{Email: "existing@example.com"})
if err != nil {
    var mongoErr mongo.WriteException
    if errors.As(err, &mongoErr) {
        for _, e := range mongoErr.WriteErrors {
            if e.Code == 11000 { // Duplicate key error code
                log.Println("Duplicate key error:", e.Message)
            }
        }
    } else {
        log.Printf("Other error: %v", err)
    }
}
```

## Transaction Error Handling

In transactions, if the callback function returns an error, the transaction will automatically rollback.

```go
err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
    // ... transaction operations ...
    return errors.New("something went wrong in transaction")
})

if err != nil {
    log.Printf("Transaction failed and rolled back: %v", err)
}
```

## Next Steps

- [Performance](/reference/performance/) - Learn performance optimization
- [Best Practices](/best-practices/) - Learn development best practices
