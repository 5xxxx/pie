---
title: 事务管理
description: 学习 Pie 的事务管理功能，确保数据一致性
---

# 事务管理

Pie 提供了简单易用的事务管理功能，确保数据操作的原子性和一致性。

## 基本用法

### 使用引擎执行事务

```go
// 使用引擎执行事务
err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
    // 在事务中执行操作
    session := pie.Table[User](engine)
    
    // 插入用户
    user := &User{Name: "张三", Email: "zhangsan@example.com"}
    _, err := session.Insert(txCtx, user)
    if err != nil {
        return err
    }
    
    // 更新其他集合
    orderSession := pie.Table[Order](engine)
    order := &Order{UserID: user.ID, Amount: 100.0}
    _, err = orderSession.Insert(txCtx, order)
    return err
})

if err != nil {
    log.Fatal("Transaction failed:", err)
}
```

### 使用事务管理器

```go
// 创建事务管理器
tx := pie.MustTransaction(engine)

// 执行事务
err := tx.Transaction(ctx, func(txCtx context.Context) error {
    // 事务操作
    session := pie.Table[User](engine)
    
    // 插入用户
    _, err := session.Insert(txCtx, &User{
        Name:  "李四",
        Email: "lisi@example.com",
    })
    if err != nil {
        return err
    }
    
    // 更新用户状态
    _, err = session.
        Where("email", "lisi@example.com").
        Update(txCtx, bson.D{{"$set", bson.D{{"status", "active"}}}})
    
    return err
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

## 实际应用场景

### 用户注册流程

```go
func registerUser(ctx context.Context, userData *UserRegistrationData) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        userSession := pie.Table[User](engine)
        profileSession := pie.Table[UserProfile](engine)
        
        // 1. 检查邮箱是否已存在
        exists, err := userSession.Where("email", userData.Email).Exists(txCtx)
        if err != nil {
            return fmt.Errorf("failed to check email existence: %w", err)
        }
        if exists {
            return errors.New("email already exists")
        }
        
        // 2. 创建用户
        user := &User{
            Name:     userData.Name,
            Email:    userData.Email,
            Password: hashPassword(userData.Password),
            Status:   "pending",
        }
        
        result, err := userSession.Insert(txCtx, user)
        if err != nil {
            return fmt.Errorf("failed to create user: %w", err)
        }
        
        // 3. 创建用户资料
        profile := &UserProfile{
            UserID:   result.InsertedID.(bson.ObjectID),
            Bio:      userData.Bio,
            Location: userData.Location,
        }
        
        _, err = profileSession.Insert(txCtx, profile)
        if err != nil {
            return fmt.Errorf("failed to create user profile: %w", err)
        }
        
        // 4. 发送验证邮件（这里只是模拟）
        if err := sendVerificationEmail(user.Email); err != nil {
            // 如果邮件发送失败，整个事务会回滚
            return fmt.Errorf("failed to send verification email: %w", err)
        }
        
        return nil
    })
}
```

### 订单处理流程

```go
func processOrder(ctx context.Context, orderData *OrderData) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        orderSession := pie.Table[Order](engine)
        userSession := pie.Table[User](engine)
        productSession := pie.Table[Product](engine)
        
        // 1. 验证用户存在
        user, err := userSession.FindByID(txCtx, orderData.UserID)
        if err != nil {
            return fmt.Errorf("user not found: %w", err)
        }
        
        // 2. 验证产品库存
        for _, item := range orderData.Items {
            product, err := productSession.FindByID(txCtx, item.ProductID)
            if err != nil {
                return fmt.Errorf("product not found: %w", err)
            }
            
            if product.Stock < item.Quantity {
                return fmt.Errorf("insufficient stock for product %s", product.Name)
            }
        }
        
        // 3. 创建订单
        order := &Order{
            UserID:    orderData.UserID,
            Items:     orderData.Items,
            Total:     calculateTotal(orderData.Items),
            Status:    "pending",
            CreatedAt: time.Now(),
        }
        
        orderResult, err := orderSession.Insert(txCtx, order)
        if err != nil {
            return fmt.Errorf("failed to create order: %w", err)
        }
        
        // 4. 更新产品库存
        for _, item := range orderData.Items {
            _, err = productSession.
                Where("_id", item.ProductID).
                Update(txCtx, bson.D{{"$inc", bson.D{{"stock", -item.Quantity}}}})
            if err != nil {
                return fmt.Errorf("failed to update stock: %w", err)
            }
        }
        
        // 5. 扣除用户余额
        _, err = userSession.
            Where("_id", orderData.UserID).
            Update(txCtx, bson.D{{"$inc", bson.D{{"balance", -order.Total}}}})
        if err != nil {
            return fmt.Errorf("failed to deduct balance: %w", err)
        }
        
        // 6. 更新订单状态
        _, err = orderSession.
            Where("_id", orderResult.InsertedID).
            Update(txCtx, bson.D{{"$set", bson.D{{"status", "completed"}}}})
        if err != nil {
            return fmt.Errorf("failed to update order status: %w", err)
        }
        
        return nil
    })
}
```

### 积分转账

```go
func transferPoints(ctx context.Context, fromUserID, toUserID bson.ObjectID, points int) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        userSession := pie.Table[User](engine)
        
        // 1. 检查发送方积分是否足够
        fromUser, err := userSession.FindByID(txCtx, fromUserID)
        if err != nil {
            return fmt.Errorf("sender not found: %w", err)
        }
        
        if fromUser.Points < points {
            return errors.New("insufficient points")
        }
        
        // 2. 检查接收方是否存在
        _, err = userSession.FindByID(txCtx, toUserID)
        if err != nil {
            return fmt.Errorf("receiver not found: %w", err)
        }
        
        // 3. 扣除发送方积分
        _, err = userSession.
            Where("_id", fromUserID).
            Update(txCtx, bson.D{{"$inc", bson.D{{"points", -points}}}})
        if err != nil {
            return fmt.Errorf("failed to deduct points: %w", err)
        }
        
        // 4. 增加接收方积分
        _, err = userSession.
            Where("_id", toUserID).
            Update(txCtx, bson.D{{"$inc", bson.D{{"points", points}}}})
        if err != nil {
            return fmt.Errorf("failed to add points: %w", err)
        }
        
        // 5. 记录转账日志
        logSession := pie.Table[TransferLog](engine)
        _, err = logSession.Insert(txCtx, &TransferLog{
            FromUserID: fromUserID,
            ToUserID:   toUserID,
            Points:     points,
            Timestamp:  time.Now(),
        })
        if err != nil {
            return fmt.Errorf("failed to log transfer: %w", err)
        }
        
        return nil
    })
}
```

## 高级用法

### 嵌套事务

```go
func complexBusinessLogic(ctx context.Context) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        // 主事务逻辑
        if err := createUser(txCtx); err != nil {
            return err
        }
        
        // 嵌套事务（实际上会使用同一个事务上下文）
        if err := engine.WithTransaction(txCtx, func(nestedCtx context.Context) error {
            return createUserProfile(nestedCtx)
        }); err != nil {
            return err
        }
        
        return nil
    })
}
```

### 条件事务

```go
func conditionalTransaction(ctx context.Context, shouldUseTransaction bool) error {
    if shouldUseTransaction {
        return engine.WithTransaction(ctx, func(txCtx context.Context) error {
            return performOperations(txCtx)
        })
    } else {
        return performOperations(ctx)
    }
}

func performOperations(ctx context.Context) error {
    session := pie.Table[User](engine)
    _, err := session.Insert(ctx, &User{Name: "测试用户"})
    return err
}
```

### 事务重试

```go
func transactionWithRetry(ctx context.Context, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
            return performBusinessLogic(txCtx)
        })
        
        if err == nil {
            return nil // 成功
        }
        
        lastErr = err
        
        // 检查是否是可重试的错误
        if !isRetryableError(err) {
            return err
        }
        
        // 等待一段时间后重试
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return fmt.Errorf("transaction failed after %d retries: %w", maxRetries, lastErr)
}

func isRetryableError(err error) bool {
    // 检查是否是网络错误、超时等可重试的错误
    if pie.IsTimeoutError(err) {
        return true
    }
    
    // 检查是否是 MongoDB 的临时错误
    if mongoErr, ok := err.(pie.MongoError); ok {
        return mongoErr.Code == 112 // WriteConflict
    }
    
    return false
}
```

## 错误处理

### 事务错误类型

```go
func handleTransactionError(err error) {
    if err == nil {
        return
    }
    
    // 检查特定错误类型
    if pie.IsNotFoundError(err) {
        log.Println("Transaction failed: resource not found")
        return
    }
    
    if pie.IsTimeoutError(err) {
        log.Println("Transaction failed: timeout")
        return
    }
    
    if pie.IsDuplicateKeyError(err) {
        log.Println("Transaction failed: duplicate key")
        return
    }
    
    // 检查 MongoDB 错误
    if mongoErr, ok := err.(pie.MongoError); ok {
        switch mongoErr.Code {
        case 112: // WriteConflict
            log.Println("Transaction failed: write conflict")
        case 11000: // DuplicateKey
            log.Println("Transaction failed: duplicate key")
        case 11001: // DuplicateKey
            log.Println("Transaction failed: duplicate key")
        default:
            log.Printf("Transaction failed: MongoDB error %d - %s", 
                mongoErr.Code, mongoErr.Message)
        }
        return
    }
    
    log.Printf("Transaction failed: %v", err)
}
```

### 事务回滚处理

```go
func transactionWithRollbackHandling(ctx context.Context) error {
    var rollbackFuncs []func() error
    
    err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
        // 执行操作并记录回滚函数
        if err := operation1(txCtx, &rollbackFuncs); err != nil {
            return err
        }
        
        if err := operation2(txCtx, &rollbackFuncs); err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        // 执行回滚操作
        for i := len(rollbackFuncs) - 1; i >= 0; i-- {
            if rollbackErr := rollbackFuncs[i](); rollbackErr != nil {
                log.Printf("Rollback operation %d failed: %v", i, rollbackErr)
            }
        }
    }
    
    return err
}

func operation1(ctx context.Context, rollbackFuncs *[]func() error) error {
    // 执行操作
    session := pie.Table[User](engine)
    result, err := session.Insert(ctx, &User{Name: "用户1"})
    if err != nil {
        return err
    }
    
    // 记录回滚函数
    *rollbackFuncs = append(*rollbackFuncs, func() error {
        _, err := session.Where("_id", result.InsertedID).Delete(context.Background())
        return err
    })
    
    return nil
}
```

## 性能优化

### 1. 合理使用事务

```go
// 好的做法：将相关操作放在一个事务中
func goodTransactionUsage(ctx context.Context) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        // 相关的操作放在一起
        if err := createUser(txCtx); err != nil {
            return err
        }
        if err := createUserProfile(txCtx); err != nil {
            return err
        }
        return nil
    })
}

// 避免的做法：将不相关的操作放在一个事务中
func badTransactionUsage(ctx context.Context) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        // 不相关的操作会增加事务时间
        if err := createUser(txCtx); err != nil {
            return err
        }
        if err := updateSystemConfig(txCtx); err != nil { // 不相关
            return err
        }
        return nil
    })
}
```

### 2. 避免长时间事务

```go
func avoidLongTransactions(ctx context.Context) error {
    // 将长时间操作移到事务外
    if err := prepareData(ctx); err != nil {
        return err
    }
    
    // 只将必要的数据库操作放在事务中
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        return performDatabaseOperations(txCtx)
    })
}
```

### 3. 使用适当的隔离级别

```go
// 对于读多写少的场景，可以使用读已提交
func readCommittedTransaction(ctx context.Context) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        // 设置读已提交隔离级别
        // 注意：MongoDB 的事务隔离级别是有限的
        return performOperations(txCtx)
    })
}
```

## 最佳实践

### 1. 事务边界设计

```go
// 好的事务边界：业务逻辑的原子性
func createUserWithProfile(ctx context.Context, userData *UserData) error {
    return engine.WithTransaction(ctx, func(txCtx context.Context) error {
        // 用户和用户资料应该作为一个整体创建
        if err := createUser(txCtx, userData); err != nil {
            return err
        }
        if err := createUserProfile(txCtx, userData); err != nil {
            return err
        }
        return nil
    })
}
```

### 2. 错误处理策略

```go
func robustTransaction(ctx context.Context) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
            return performBusinessLogic(txCtx)
        })
        
        if err == nil {
            return nil
        }
        
        // 检查是否是可重试的错误
        if !isRetryableError(err) {
            return err
        }
        
        // 指数退避
        time.Sleep(time.Duration(1<<uint(i)) * time.Second)
    }
    
    return errors.New("max retries exceeded")
}
```

### 3. 事务监控

```go
func monitoredTransaction(ctx context.Context) error {
    start := time.Now()
    
    err := engine.WithTransaction(ctx, func(txCtx context.Context) error {
        return performOperations(txCtx)
    })
    
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("Transaction failed after %v: %v", duration, err)
    } else {
        log.Printf("Transaction completed in %v", duration)
    }
    
    return err
}
```

## 下一步

- [缓存支持](/advanced/cache/) - 学习缓存功能
- [钩子系统](/advanced/hooks/) - 掌握生命周期钩子
- [软删除](/advanced/soft-delete/) - 学习软删除功能
