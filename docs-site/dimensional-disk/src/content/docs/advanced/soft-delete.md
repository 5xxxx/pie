---
title: Soft Delete
description: Learn Pie's soft delete functionality to protect important data
---

# Soft Delete

Pie provides built-in soft delete functionality, allowing you to "delete" data without actually removing it from the database, protecting important data and providing data recovery capabilities.

## Basic Usage

### Define Soft Delete Fields

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

type Product struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Price     float64       `bson:"price"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
}
```

### Soft Delete Operations

```go
// Soft delete user
err := session.Where("email", "test@example.com").SoftDelete(ctx)

// Soft delete multiple users
err := session.Where("status", "inactive").SoftDeleteMany(ctx)

// Soft delete specific user
err := session.Where("_id", userID).SoftDelete(ctx)
```

### Restore Soft Deleted

```go
// Restore soft deleted user
err := session.Where("email", "test@example.com").Restore(ctx)

// Restore multiple users
err := session.Where("status", "inactive").RestoreMany(ctx)

// Restore specific user
err := session.Where("_id", userID).Restore(ctx)
```

### Force Delete

```go
// Force delete (physical delete)
err := session.Where("email", "test@example.com").ForceDelete(ctx)

// Force delete multiple records
err := session.Where("status", "inactive").ForceDeleteMany(ctx)
```

### Query Behavior

```go
// Query automatically excludes soft deleted records
var users []User
err := session.Find(ctx, &users) // Automatically filters deleted records

// Include soft deleted records in query
var allUsers []User
err := session.IncludeDeleted().Find(ctx, &allUsers)

// Query only soft deleted records
var deletedUsers []User
err := session.OnlyDeleted().Find(ctx, &deletedUsers)
```

## Advanced Usage

### Custom Soft Delete Field

```go
type User struct {
    ID           bson.ObjectID `bson:"_id,omitempty"`
    Name         string        `bson:"name"`
    Email        string        `bson:"email"`
    RemovedAt    *time.Time    `bson:"removed_at,omitempty" pie:"soft_delete"`
}

// Configure custom soft delete field
engine.WithSoftDeleteField("removed_at")
```

### Soft Delete with Hooks

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

// Before soft delete hook
func (u *User) BeforeSoftDelete(ctx context.Context) error {
    log.Printf("About to soft delete user: %s", u.Name)
    
    // Check if user has active orders
    if hasActiveOrders(u.ID) {
        return fmt.Errorf("cannot delete user with active orders")
    }
    
    return nil
}

// After soft delete hook
func (u *User) AfterSoftDelete(ctx context.Context) error {
    log.Printf("User soft deleted: %s", u.Name)
    
    // Send notification
    go sendDeletionNotification(u.Email)
    
    return nil
}

// Before restore hook
func (u *User) BeforeRestore(ctx context.Context) error {
    log.Printf("About to restore user: %s", u.Name)
    return nil
}

// After restore hook
func (u *User) AfterRestore(ctx context.Context) error {
    log.Printf("User restored: %s", u.Name)
    
    // Send restoration notification
    go sendRestorationNotification(u.Email)
    
    return nil
}
```

### Batch Operations

```go
// Batch soft delete
func softDeleteInactiveUsers() error {
    session := pie.Table[User](engine)
    
    // Soft delete users inactive for more than 30 days
    cutoffDate := time.Now().AddDate(0, 0, -30)
    
    result, err := session.
        Where("status", "inactive").
        Where("last_login", pie.Lt(cutoffDate)).
        SoftDeleteMany(ctx)
    
    if err != nil {
        return err
    }
    
    log.Printf("Soft deleted %d inactive users", result.DeletedCount)
    return nil
}

// Batch restore
func restoreUsersByStatus(status string) error {
    session := pie.Table[User](engine)
    
    result, err := session.
        OnlyDeleted().
        Where("status", status).
        RestoreMany(ctx)
    
    if err != nil {
        return err
    }
    
    log.Printf("Restored %d users with status %s", result.ModifiedCount, status)
    return nil
}
```

## Real-World Examples

### User Management System

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Username  string        `bson:"username"`
    Email     string        `bson:"email"`
    Status    string        `bson:"status"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

// Deactivate user (soft delete)
func deactivateUser(userID bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // Soft delete user
    err := session.Where("_id", userID).SoftDelete(ctx)
    if err != nil {
        return fmt.Errorf("failed to deactivate user: %w", err)
    }
    
    // Log the action
    log.Printf("User %s deactivated", userID.Hex())
    
    return nil
}

// Reactivate user (restore)
func reactivateUser(userID bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // Restore user
    err := session.Where("_id", userID).Restore(ctx)
    if err != nil {
        return fmt.Errorf("failed to reactivate user: %w", err)
    }
    
    // Log the action
    log.Printf("User %s reactivated", userID.Hex())
    
    return nil
}

// Get active users (excludes soft deleted)
func getActiveUsers() ([]User, error) {
    session := pie.Table[User](engine)
    
    var users []User
    err := session.
        Where("status", "active").
        Find(ctx, &users)
    
    return users, err
}

// Get all users including soft deleted
func getAllUsers() ([]User, error) {
    session := pie.Table[User](engine)
    
    var users []User
    err := session.IncludeDeleted().Find(ctx, &users)
    
    return users, err
}
```

### Product Catalog Management

```go
type Product struct {
    ID          bson.ObjectID `bson:"_id,omitempty"`
    Name        string        `bson:"name"`
    SKU         string        `bson:"sku"`
    Price       float64       `bson:"price"`
    Category    string        `bson:"category"`
    Status      string        `bson:"status"`
    DeletedAt   *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete"`
    CreatedAt   time.Time     `bson:"created_at"`
    UpdatedAt   time.Time     `bson:"updated_at"`
}

// Discontinue product (soft delete)
func discontinueProduct(productID bson.ObjectID) error {
    session := pie.Table[Product](engine)
    
    // Check if product has active orders
    if hasActiveOrders(productID) {
        return fmt.Errorf("cannot discontinue product with active orders")
    }
    
    // Soft delete product
    err := session.Where("_id", productID).SoftDelete(ctx)
    if err != nil {
        return fmt.Errorf("failed to discontinue product: %w", err)
    }
    
    // Update inventory status
    go updateInventoryStatus(productID, "discontinued")
    
    return nil
}

// Reintroduce product (restore)
func reintroduceProduct(productID bson.ObjectID) error {
    session := pie.Table[Product](engine)
    
    // Restore product
    err := session.Where("_id", productID).Restore(ctx)
    if err != nil {
        return fmt.Errorf("failed to reintroduce product: %w", err)
    }
    
    // Update inventory status
    go updateInventoryStatus(productID, "available")
    
    return nil
}

// Get available products (excludes discontinued)
func getAvailableProducts() ([]Product, error) {
    session := pie.Table[Product](engine)
    
    var products []Product
    err := session.
        Where("status", "available").
        Find(ctx, &products)
    
    return products, err
}
```

## Configuration

### Global Configuration

```go
// Configure soft delete field globally
engine, err := pie.NewEngine(
    context.Background(),
    "mydb",
    pie.WithURI("mongodb://localhost:27017"),
    pie.WithSoftDeleteField("deleted_at"), // Default field name
)
```

### Per-Session Configuration

```go
// Configure soft delete for specific session
session := pie.Table[User](engine).WithSoftDeleteField("removed_at")
```

## Best Practices

### 1. Use Appropriate Field Names

```go
// Good: Clear and descriptive
type User struct {
    DeletedAt *time.Time `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

// Good: Context-specific
type Order struct {
    CancelledAt *time.Time `bson:"cancelled_at,omitempty" pie:"soft_delete"`
}

// Bad: Unclear
type User struct {
    Flag *time.Time `bson:"flag,omitempty" pie:"soft_delete"`
}
```

### 2. Handle Dependencies

```go
func (u *User) BeforeSoftDelete(ctx context.Context) error {
    // Check for dependent data
    if hasActiveOrders(u.ID) {
        return fmt.Errorf("cannot delete user with active orders")
    }
    
    if hasUnpaidInvoices(u.ID) {
        return fmt.Errorf("cannot delete user with unpaid invoices")
    }
    
    return nil
}
```

### 3. Use Indexes for Performance

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    DeletedAt *time.Time    `bson:"deleted_at,omitempty" pie:"soft_delete,index"`
    CreatedAt time.Time     `bson:"created_at"`
}
```

### 4. Regular Cleanup

```go
// Clean up old soft deleted records
func cleanupOldSoftDeleted() error {
    session := pie.Table[User](engine)
    
    // Delete records soft deleted more than 1 year ago
    cutoffDate := time.Now().AddDate(-1, 0, 0)
    
    result, err := session.
        OnlyDeleted().
        Where("deleted_at", pie.Lt(cutoffDate)).
        ForceDeleteMany(ctx)
    
    if err != nil {
        return err
    }
    
    log.Printf("Permanently deleted %d old records", result.DeletedCount)
    return nil
}
```

## Next Steps

- [Query Scopes](/advanced/scopes/) - Learn about query scopes
- [Index Management](/advanced/indexes/) - Learn about index management
- [Best Practices](/best-practices/) - Learn development best practices
