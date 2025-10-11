---
title: Hook System
description: Learn Pie's complete lifecycle hook support
---

# Hook System

Pie provides complete lifecycle hook support, allowing you to execute custom logic at various stages of data operations.

## Basic Usage

### Define Hook Methods

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name"`
    Email     string        `bson:"email"`
    Password  string        `bson:"password"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

// Before create hook
func (u *User) BeforeCreate(ctx context.Context) error {
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    
    // Hash password
    hashedPassword, err := hashPassword(u.Password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }
    u.Password = hashedPassword
    
    return nil
}

// After create hook
func (u *User) AfterCreate(ctx context.Context) error {
    log.Printf("User %s created with ID %s", u.Name, u.ID.Hex())
    
    // Send welcome email
    go sendWelcomeEmail(u.Email, u.Name)
    
    return nil
}

// Before update hook
func (u *User) BeforeUpdate(ctx context.Context) error {
    u.UpdatedAt = time.Now()
    
    // Validate email format
    if !isValidEmail(u.Email) {
        return fmt.Errorf("invalid email format: %s", u.Email)
    }
    
    return nil
}

// After update hook
func (u *User) AfterUpdate(ctx context.Context) error {
    log.Printf("User %s updated", u.Name)
    
    // Update user session
    go updateUserSession(u.ID)
    
    return nil
}

// Before delete hook
func (u *User) BeforeDelete(ctx context.Context) error {
    log.Printf("About to delete user %s", u.Name)
    
    // Check if user has active orders
    if hasActiveOrders(u.ID) {
        return fmt.Errorf("cannot delete user with active orders")
    }
    
    return nil
}

// After delete hook
func (u *User) AfterDelete(ctx context.Context) error {
    log.Printf("User %s deleted", u.Name)
    
    // Clean up related data
    go cleanupUserData(u.ID)
    
    return nil
}

// After find hook
func (u *User) AfterFind(ctx context.Context) error {
    // Remove sensitive data
    u.Password = ""
    
    // Format display name
    u.Name = strings.TrimSpace(u.Name)
    
    return nil
}
```

## Hook Types

### Create Hooks

```go
// BeforeCreate - Execute before document creation
func (u *User) BeforeCreate(ctx context.Context) error {
    // Set timestamps
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    
    // Generate ID if not set
    if u.ID.IsZero() {
        u.ID = bson.NewObjectID()
    }
    
    // Validate required fields
    if u.Name == "" {
        return fmt.Errorf("name is required")
    }
    
    return nil
}

// AfterCreate - Execute after document creation
func (u *User) AfterCreate(ctx context.Context) error {
    // Log creation
    log.Printf("Created user: %s", u.Name)
    
    // Send notifications
    go notifyUserCreated(u)
    
    return nil
}
```

### Update Hooks

```go
// BeforeUpdate - Execute before document update
func (u *User) BeforeUpdate(ctx context.Context) error {
    // Update timestamp
    u.UpdatedAt = time.Now()
    
    // Validate changes
    if u.Email != "" && !isValidEmail(u.Email) {
        return fmt.Errorf("invalid email: %s", u.Email)
    }
    
    return nil
}

// AfterUpdate - Execute after document update
func (u *User) AfterUpdate(ctx context.Context) error {
    // Log update
    log.Printf("Updated user: %s", u.Name)
    
    // Invalidate cache
    go invalidateUserCache(u.ID)
    
    return nil
}
```

### Delete Hooks

```go
// BeforeDelete - Execute before document deletion
func (u *User) BeforeDelete(ctx context.Context) error {
    // Check dependencies
    if hasDependentData(u.ID) {
        return fmt.Errorf("cannot delete user with dependent data")
    }
    
    // Backup user data
    go backupUserData(u)
    
    return nil
}

// AfterDelete - Execute after document deletion
func (u *User) AfterDelete(ctx context.Context) error {
    // Log deletion
    log.Printf("Deleted user: %s", u.Name)
    
    // Clean up related data
    go cleanupUserData(u.ID)
    
    return nil
}
```

### Find Hooks

```go
// AfterFind - Execute after document retrieval
func (u *User) AfterFind(ctx context.Context) error {
    // Remove sensitive data
    u.Password = ""
    
    // Format data
    u.Name = strings.TrimSpace(u.Name)
    u.Email = strings.ToLower(u.Email)
    
    // Add computed fields
    u.DisplayName = fmt.Sprintf("%s (%s)", u.Name, u.Email)
    
    return nil
}
```

## Advanced Usage

### Conditional Hooks

```go
// Execute hooks based on conditions
func (u *User) BeforeUpdate(ctx context.Context) error {
    // Only validate email if it's being updated
    if u.Email != "" {
        if !isValidEmail(u.Email) {
            return fmt.Errorf("invalid email: %s", u.Email)
        }
    }
    
    return nil
}
```

### Error Handling in Hooks

```go
func (u *User) BeforeCreate(ctx context.Context) error {
    // Hash password with error handling
    if u.Password != "" {
        hashedPassword, err := hashPassword(u.Password)
        if err != nil {
            return fmt.Errorf("password hashing failed: %w", err)
        }
        u.Password = hashedPassword
    }
    
    return nil
}
```

### Context Usage

```go
func (u *User) AfterCreate(ctx context.Context) error {
    // Get user ID from context
    if userID, ok := ctx.Value("userID").(bson.ObjectID); ok {
        log.Printf("User %s created by %s", u.Name, userID.Hex())
    }
    
    // Check if operation is in transaction
    if isInTransaction(ctx) {
        log.Println("User created within transaction")
    }
    
    return nil
}
```

## Real-World Examples

### User Registration Flow

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Username  string        `bson:"username"`
    Email     string        `bson:"email"`
    Password  string        `bson:"password"`
    Status    string        `bson:"status"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

func (u *User) BeforeCreate(ctx context.Context) error {
    // Set timestamps
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    
    // Set default status
    u.Status = "pending"
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("password hashing failed: %w", err)
    }
    u.Password = string(hashedPassword)
    
    // Validate email
    if !isValidEmail(u.Email) {
        return fmt.Errorf("invalid email format")
    }
    
    return nil
}

func (u *User) AfterCreate(ctx context.Context) error {
    // Send verification email
    go sendVerificationEmail(u.Email, u.Username)
    
    // Log registration
    log.Printf("New user registered: %s (%s)", u.Username, u.Email)
    
    return nil
}
```

### Order Processing Flow

```go
type Order struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    UserID    bson.ObjectID `bson:"user_id"`
    Items     []OrderItem   `bson:"items"`
    Total     float64       `bson:"total"`
    Status    string        `bson:"status"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

func (o *Order) BeforeCreate(ctx context.Context) error {
    // Set timestamps
    o.CreatedAt = time.Now()
    o.UpdatedAt = time.Now()
    
    // Set initial status
    o.Status = "pending"
    
    // Calculate total
    o.Total = calculateOrderTotal(o.Items)
    
    // Validate items
    if len(o.Items) == 0 {
        return fmt.Errorf("order must have at least one item")
    }
    
    return nil
}

func (o *Order) AfterCreate(ctx context.Context) error {
    // Reserve inventory
    go reserveInventory(o.Items)
    
    // Send confirmation email
    go sendOrderConfirmation(o)
    
    return nil
}

func (o *Order) BeforeUpdate(ctx context.Context) error {
    // Update timestamp
    o.UpdatedAt = time.Now()
    
    // Validate status transition
    if !isValidStatusTransition(o.Status, getNewStatus()) {
        return fmt.Errorf("invalid status transition")
    }
    
    return nil
}
```

## Best Practices

### 1. Keep Hooks Lightweight

```go
// Good: Lightweight hook
func (u *User) BeforeCreate(ctx context.Context) error {
    u.CreatedAt = time.Now()
    return nil
}

// Bad: Heavy operations in hook
func (u *User) BeforeCreate(ctx context.Context) error {
    // Don't do heavy operations in hooks
    result := heavyComputation() // This blocks the operation
    u.SomeField = result
    return nil
}
```

### 2. Use Goroutines for Heavy Operations

```go
func (u *User) AfterCreate(ctx context.Context) error {
    // Use goroutines for non-critical operations
    go func() {
        sendWelcomeEmail(u.Email)
        updateAnalytics(u.ID)
        logUserActivity(u)
    }()
    
    return nil
}
```

### 3. Handle Errors Properly

```go
func (u *User) BeforeCreate(ctx context.Context) error {
    // Validate and return meaningful errors
    if u.Email == "" {
        return fmt.Errorf("email is required")
    }
    
    if !isValidEmail(u.Email) {
        return fmt.Errorf("invalid email format: %s", u.Email)
    }
    
    return nil
}
```

### 4. Use Context for Additional Data

```go
func (u *User) AfterCreate(ctx context.Context) error {
    // Get additional data from context
    if requestID, ok := ctx.Value("requestID").(string); ok {
        log.Printf("User created in request %s", requestID)
    }
    
    return nil
}
```

## Next Steps

- [Soft Delete](/advanced/soft-delete/) - Learn about soft delete
- [Query Scopes](/advanced/scopes/) - Learn about query scopes
- [Best Practices](/best-practices/) - Learn development best practices
