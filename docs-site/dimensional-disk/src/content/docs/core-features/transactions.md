---
title: Transactions
description: Pie provides simple and easy-to-use transaction operations
---

# Transaction Management

```go
// Use engine to execute transaction
err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
    // Execute operations in transaction
    session := pie.Table[User](engine)

    // Insert user
    _, err := session.Insert(txCtx, &User{Name: "Transaction User"})
    if err != nil {
        return err
    }

    // Update other collection
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

## Basic Transaction Usage

```go
func TransferPoints(fromUserID, toUserID primitive.ObjectID, points int) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        session := pie.Table[User](engine)
        
        // Check sender user
        var fromUser User
        err := session.Where("_id", fromUserID).First(txCtx, &fromUser)
        if err != nil {
            return err
        }
        
        // Check if user has enough points
        if fromUser.Points < points {
            return errors.New("insufficient points")
        }
        
        // Check receiver user
        var toUser User
        err = session.Where("_id", toUserID).First(txCtx, &toUser)
        if err != nil {
            return err
        }
        
        // Deduct points from sender
        _, err = session.Where("_id", fromUserID).Update(txCtx, 
            bson.D{{"$inc", bson.D{{"points", -points}}}})
        if err != nil {
            return err
        }
        
        // Add points to receiver
        _, err = session.Where("_id", toUserID).Update(txCtx, 
            bson.D{{"$inc", bson.D{{"points", points}}}})
        if err != nil {
            return err
        }
        
        return nil
    })
}
```

## Transaction Manager

```go
func UseTransactionManager() error {
    tx := pie.MustTransaction(engine)
    
    return tx.Execute(ctx, func(txCtx context.Context) error {
        session := pie.Table[User](engine)
        
        // Insert user
        user := &User{Name: "John Doe", Email: "john@example.com"}
        result, err := session.Insert(txCtx, user)
        if err != nil {
            return err
        }
        
        // Insert user profile
        profile := &Profile{UserID: result.InsertedID, Bio: "Software Developer"}
        _, err = pie.Table[Profile](engine).Insert(txCtx, profile)
        if err != nil {
            return err
        }
        
        return nil
    })
}
```

## Error Handling

```go
func TransactionWithErrorHandling() error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        session := pie.Table[User](engine)
        
        // Try to insert user
        user := &User{Name: "John Doe", Email: "john@example.com"}
        _, err := session.Insert(txCtx, user)
        if err != nil {
            // Transaction will be rolled back automatically
            return fmt.Errorf("failed to insert user: %w", err)
        }
        
        // Try to insert profile
        profile := &Profile{UserID: user.ID, Bio: "Software Developer"}
        _, err = pie.Table[Profile](engine).Insert(txCtx, profile)
        if err != nil {
            // Transaction will be rolled back automatically
            return fmt.Errorf("failed to insert profile: %w", err)
        }
        
        return nil
    })
}
```

## Next Steps

- [Query Builder](/core-features/query-builder/) - Learn about query methods
- [Bulk Operations](/core-features/bulk-operations/) - Learn about bulk operations
- [Best Practices](/best-practices/) - Learn development best practices
