---
title: 错误处理
description: 学习 Pie 的错误处理机制和最佳实践
---

# 错误处理

Pie 提供了完善的错误处理机制，帮助您识别和处理各种数据库操作错误。

## 基本用法

### 错误类型检查

```go
// 检查特定错误类型
if pie.IsNotFoundError(err) {
    log.Println("Document not found")
}

if pie.IsDuplicateKeyError(err) {
    log.Println("Duplicate key error")
}

if pie.IsTimeoutError(err) {
    log.Println("Operation timeout")
}

if pie.IsConnectionError(err) {
    log.Println("Connection error")
}
```

### 错误详情获取

```go
// 获取 MongoDB 错误详情
if mongoErr, ok := err.(pie.MongoError); ok {
    log.Printf("MongoDB error code: %d", mongoErr.Code)
    log.Printf("MongoDB error message: %s", mongoErr.Message)
    log.Printf("MongoDB error details: %v", mongoErr.Details)
}
```

## 错误类型

### 基础错误类型

```go
// 文档未找到错误
if pie.IsNotFoundError(err) {
    return fmt.Errorf("user not found: %w", err)
}

// 重复键错误
if pie.IsDuplicateKeyError(err) {
    return fmt.Errorf("email already exists: %w", err)
}

// 超时错误
if pie.IsTimeoutError(err) {
    return fmt.Errorf("operation timeout: %w", err)
}

// 连接错误
if pie.IsConnectionError(err) {
    return fmt.Errorf("database connection failed: %w", err)
}

// 网络错误
if pie.IsNetworkError(err) {
    return fmt.Errorf("network error: %w", err)
}
```

### MongoDB 错误代码

```go
// 检查特定 MongoDB 错误代码
if mongoErr, ok := err.(pie.MongoError); ok {
    switch mongoErr.Code {
    case 11000: // DuplicateKey
        return fmt.Errorf("duplicate key error: %w", err)
    case 11001: // DuplicateKey
        return fmt.Errorf("duplicate key error: %w", err)
    case 112: // WriteConflict
        return fmt.Errorf("write conflict: %w", err)
    case 11600: // Interrupted
        return fmt.Errorf("operation interrupted: %w", err)
    case 11601: // InterruptedAtShutdown
        return fmt.Errorf("operation interrupted at shutdown: %w", err)
    case 11602: // InterruptedDueToReplStateChange
        return fmt.Errorf("operation interrupted due to replica set state change: %w", err)
    case 12500: // LockTimeout
        return fmt.Errorf("lock timeout: %w", err)
    case 12501: // LockTimeout
        return fmt.Errorf("lock timeout: %w", err)
    default:
        return fmt.Errorf("MongoDB error %d: %s", mongoErr.Code, mongoErr.Message)
    }
}
```

## 实际应用场景

### 用户操作错误处理

```go
func createUser(userData *UserData) error {
    session := pie.Table[User](engine)
    
    // 检查邮箱是否已存在
    exists, err := session.Where("email", userData.Email).Exists(ctx)
    if err != nil {
        return fmt.Errorf("failed to check email existence: %w", err)
    }
    if exists {
        return errors.New("email already exists")
    }
    
    // 创建用户
    user := &User{
        Name:  userData.Name,
        Email: userData.Email,
    }
    
    _, err = session.Insert(ctx, user)
    if err != nil {
        if pie.IsDuplicateKeyError(err) {
            return errors.New("email already exists")
        }
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func getUserByID(userID bson.ObjectID) (*User, error) {
    session := pie.Table[User](engine)
    
    var user User
    err := session.Where("_id", userID).First(ctx, &user)
    if err != nil {
        if pie.IsNotFoundError(err) {
            return nil, errors.New("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func updateUser(userID bson.ObjectID, updates bson.D) error {
    session := pie.Table[User](engine)
    
    result, err := session.Where("_id", userID).Update(ctx, updates)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    if result.ModifiedCount == 0 {
        return errors.New("user not found or no changes made")
    }
    
    return nil
}
```

### 事务错误处理

```go
func transferPoints(fromUserID, toUserID bson.ObjectID, points int) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        session := pie.Table[User](engine)
        
        // 检查发送方用户
        var fromUser User
        err := session.Where("_id", fromUserID).First(txCtx, &fromUser)
        if err != nil {
            if pie.IsNotFoundError(err) {
                return errors.New("sender user not found")
            }
            return fmt.Errorf("failed to get sender user: %w", err)
        }
        
        // 检查积分是否足够
        if fromUser.Points < points {
            return errors.New("insufficient points")
        }
        
        // 检查接收方用户
        var toUser User
        err = session.Where("_id", toUserID).First(txCtx, &toUser)
        if err != nil {
            if pie.IsNotFoundError(err) {
                return errors.New("receiver user not found")
            }
            return fmt.Errorf("failed to get receiver user: %w", err)
        }
        
        // 扣除发送方积分
        _, err = session.Where("_id", fromUserID).Update(txCtx, 
            bson.D{{"$inc", bson.D{{"points", -points}}}})
        if err != nil {
            return fmt.Errorf("failed to deduct points: %w", err)
        }
        
        // 增加接收方积分
        _, err = session.Where("_id", toUserID).Update(txCtx, 
            bson.D{{"$inc", bson.D{{"points", points}}}})
        if err != nil {
            return fmt.Errorf("failed to add points: %w", err)
        }
        
        return nil
    })
}
```

### 批量操作错误处理

```go
func batchCreateUsers(users []*User) error {
    session := pie.Table[User](engine)
    
    result, err := session.InsertMany(ctx, users)
    if err != nil {
        if bulkErr, ok := err.(pie.BulkWriteError); ok {
            log.Printf("Bulk operation partially failed:")
            log.Printf("Inserted: %d", result.InsertedCount)
            log.Printf("Failed: %d", len(bulkErr.WriteErrors))
            
            for _, writeErr := range bulkErr.WriteErrors {
                log.Printf("Error at index %d: %s", writeErr.Index, writeErr.Message)
            }
            
            // 根据业务需求决定是否返回错误
            if result.InsertedCount > 0 {
                log.Println("Partial success, continuing...")
                return nil
            }
        }
        
        return fmt.Errorf("failed to create users: %w", err)
    }
    
    log.Printf("Successfully created %d users", result.InsertedCount)
    return nil
}
```

## 高级错误处理

### 错误包装

```go
// 自定义错误类型
type UserError struct {
    Code    string
    Message string
    Err     error
}

func (e *UserError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *UserError) Unwrap() error {
    return e.Err
}

// 错误包装函数
func wrapUserError(code, message string, err error) *UserError {
    return &UserError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// 使用错误包装
func createUserWithErrorWrapping(userData *UserData) error {
    session := pie.Table[User](engine)
    
    _, err := session.Insert(ctx, &User{
        Name:  userData.Name,
        Email: userData.Email,
    })
    
    if err != nil {
        if pie.IsDuplicateKeyError(err) {
            return wrapUserError("DUPLICATE_EMAIL", "email already exists", err)
        }
        return wrapUserError("CREATE_FAILED", "failed to create user", err)
    }
    
    return nil
}
```

### 错误恢复

```go
func resilientOperation() error {
    maxRetries := 3
    retryDelay := 1 * time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := performOperation()
        if err == nil {
            return nil
        }
        
        // 检查是否可重试
        if !isRetryableError(err) {
            return err
        }
        
        log.Printf("Operation failed (attempt %d/%d): %v", i+1, maxRetries, err)
        
        if i < maxRetries-1 {
            time.Sleep(retryDelay)
            retryDelay *= 2 // 指数退避
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}

func isRetryableError(err error) bool {
    // 网络错误可重试
    if pie.IsNetworkError(err) {
        return true
    }
    
    // 超时错误可重试
    if pie.IsTimeoutError(err) {
        return true
    }
    
    // 连接错误可重试
    if pie.IsConnectionError(err) {
        return true
    }
    
    // 写冲突可重试
    if mongoErr, ok := err.(pie.MongoError); ok {
        return mongoErr.Code == 112 // WriteConflict
    }
    
    return false
}
```

### 错误监控

```go
type ErrorMonitor struct {
    errorCounts map[string]int
    mutex       sync.RWMutex
}

func NewErrorMonitor() *ErrorMonitor {
    return &ErrorMonitor{
        errorCounts: make(map[string]int),
    }
}

func (m *ErrorMonitor) RecordError(err error) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    errorType := getErrorType(err)
    m.errorCounts[errorType]++
}

func (m *ErrorMonitor) GetErrorStats() map[string]int {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    stats := make(map[string]int)
    for k, v := range m.errorCounts {
        stats[k] = v
    }
    return stats
}

func getErrorType(err error) string {
    if pie.IsNotFoundError(err) {
        return "not_found"
    }
    if pie.IsDuplicateKeyError(err) {
        return "duplicate_key"
    }
    if pie.IsTimeoutError(err) {
        return "timeout"
    }
    if pie.IsConnectionError(err) {
        return "connection"
    }
    return "unknown"
}

// 使用错误监控
func monitoredOperation() error {
    monitor := NewErrorMonitor()
    
    err := performOperation()
    if err != nil {
        monitor.RecordError(err)
    }
    
    return err
}
```

## 最佳实践

### 1. 错误处理策略

```go
// 好的错误处理
func goodErrorHandling() error {
    session := pie.Table[User](engine)
    
    var user User
    err := session.Where("email", "test@example.com").First(ctx, &user)
    if err != nil {
        if pie.IsNotFoundError(err) {
            return errors.New("user not found")
        }
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    return nil
}

// 避免的错误处理
func badErrorHandling() error {
    session := pie.Table[User](engine)
    
    var user User
    err := session.Where("email", "test@example.com").First(ctx, &user)
    if err != nil {
        return err // 没有提供上下文信息
    }
    
    return nil
}
```

### 2. 错误日志记录

```go
func logErrorWithContext(operation string, err error, context map[string]interface{}) {
    log.Printf("Operation %s failed: %v", operation, err)
    
    for key, value := range context {
        log.Printf("  %s: %v", key, value)
    }
    
    // 记录堆栈跟踪
    if pie.IsConnectionError(err) {
        log.Printf("Stack trace: %+v", err)
    }
}

// 使用错误日志记录
func createUserWithLogging(userData *UserData) error {
    session := pie.Table[User](engine)
    
    _, err := session.Insert(ctx, &User{
        Name:  userData.Name,
        Email: userData.Email,
    })
    
    if err != nil {
        logErrorWithContext("create_user", err, map[string]interface{}{
            "user_name":  userData.Name,
            "user_email": userData.Email,
        })
        return err
    }
    
    return nil
}
```

### 3. 错误测试

```go
func TestErrorHandling(t *testing.T) {
    // 测试 NotFound 错误
    t.Run("NotFoundError", func(t *testing.T) {
        session := pie.Table[User](engine)
        
        var user User
        err := session.Where("_id", primitive.NewObjectID()).First(ctx, &user)
        
        assert.Error(t, err)
        assert.True(t, pie.IsNotFoundError(err))
    })
    
    // 测试 DuplicateKey 错误
    t.Run("DuplicateKeyError", func(t *testing.T) {
        session := pie.Table[User](engine)
        
        user := &User{
            Name:  "Test User",
            Email: "test@example.com",
        }
        
        // 第一次插入
        _, err := session.Insert(ctx, user)
        require.NoError(t, err)
        
        // 第二次插入应该失败
        _, err = session.Insert(ctx, user)
        assert.Error(t, err)
        assert.True(t, pie.IsDuplicateKeyError(err))
    })
}
```

## 下一步

- [性能优化](/reference/performance/) - 了解性能优化
- [最佳实践](/best-practices/) - 掌握开发最佳实践
- [快速开始](/getting-started/) - 开始使用 Pie
