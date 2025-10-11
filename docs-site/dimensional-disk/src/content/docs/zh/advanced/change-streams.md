---
title: 变更流
description: 学习 Pie 的变更流功能，实现实时数据监听
---

# 变更流

Pie 提供了变更流功能，允许您实时监听数据库中的数据变更，支持事件驱动架构。

## 基本用法

### 创建变更流监听器

```go
// 创建集合变更流监听器
watcher := pie.NewWatcher[User](engine)

// 创建数据库变更流监听器
dbWatcher := pie.NewDatabaseWatcher[User](engine)
```

### 监听集合变更

```go
// 监听用户集合的所有变更
err := watcher.
    WatchCollection().
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("变更类型: %s", change.OperationType)
        log.Printf("文档ID: %s", change.DocumentKey)
        
        switch change.OperationType {
        case "insert":
            log.Printf("新用户插入: %s", change.FullDocument.Name)
        case "update":
            log.Printf("用户更新: %s", change.FullDocument.Name)
        case "delete":
            log.Printf("用户删除: %s", change.DocumentKey)
        }
        
        return nil
    })

if err != nil {
    log.Fatal("Failed to start watcher:", err)
}
```

### 监听数据库变更

```go
// 监听整个数据库的变更
err := dbWatcher.
    WatchDatabase().
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("数据库变更: %s", change.OperationType)
        log.Printf("集合: %s", change.Collection)
        log.Printf("文档ID: %s", change.DocumentKey)
        
        return nil
    })
```

## 高级用法

### 过滤变更

```go
// 只监听插入操作
err := watcher.
    WatchCollection().
    Filter(bson.D{{"operationType", "insert"}}).
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("新用户: %s", change.FullDocument.Name)
        return nil
    })

// 只监听特定字段的更新
err = watcher.
    WatchCollection().
    Filter(bson.D{
        {"operationType", "update"},
        {"updateDescription.updatedFields.status", bson.D{{"$exists", true}}},
    }).
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("用户状态更新: %s", change.FullDocument.Name)
        return nil
    })
```

### 条件监听

```go
// 监听特定条件的变更
err := watcher.
    WatchCollection().
    Filter(bson.D{
        {"$or", []bson.D{
            {{"operationType", "insert"}},
            {{"operationType", "update"}},
        }},
        {"fullDocument.status", "active"},
    }).
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("活跃用户变更: %s", change.FullDocument.Name)
        return nil
    })
```

### 批量处理变更

```go
func processChangesInBatches() error {
    watcher := pie.NewWatcher[User](engine)
    
    var changes []*pie.ChangeEvent[User]
    batchSize := 100
    ticker := time.NewTicker(5 * time.Second)
    
    go func() {
        for {
            select {
            case <-ticker.C:
                if len(changes) > 0 {
                    processBatch(changes)
                    changes = changes[:0] // 清空切片
                }
            }
        }
    }()
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            changes = append(changes, change)
            
            // 达到批次大小时立即处理
            if len(changes) >= batchSize {
                processBatch(changes)
                changes = changes[:0]
            }
            
            return nil
        })
}

func processBatch(changes []*pie.ChangeEvent[User]) {
    log.Printf("处理 %d 个变更", len(changes))
    
    for _, change := range changes {
        // 处理单个变更
        processChange(change)
    }
}
```

## 实际应用场景

### 用户注册通知

```go
func watchUserRegistrations() error {
    watcher := pie.NewWatcher[User](engine)
    
    return watcher.
        WatchCollection().
        Filter(bson.D{{"operationType", "insert"}}).
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            user := change.FullDocument
            
            // 发送欢迎邮件
            go func() {
                if err := sendWelcomeEmail(user.Email, user.Name); err != nil {
                    log.Printf("Failed to send welcome email: %v", err)
                }
            }()
            
            // 创建用户目录
            go func() {
                if err := createUserDirectory(user.ID); err != nil {
                    log.Printf("Failed to create user directory: %v", err)
                }
            }()
            
            // 更新统计信息
            go updateUserStatistics()
            
            return nil
        })
}
```

### 订单状态变更

```go
func watchOrderStatusChanges() error {
    watcher := pie.NewWatcher[Order](engine)
    
    return watcher.
        WatchCollection().
        Filter(bson.D{
            {"operationType", "update"},
            {"updateDescription.updatedFields.status", bson.D{{"$exists", true}}},
        }).
        Start(ctx, func(change *pie.ChangeEvent[Order]) error {
            order := change.FullDocument
            
            // 根据状态执行不同操作
            switch order.Status {
            case "confirmed":
                go processOrderConfirmation(order)
            case "shipped":
                go processOrderShipment(order)
            case "delivered":
                go processOrderDelivery(order)
            case "cancelled":
                go processOrderCancellation(order)
            }
            
            return nil
        })
}

func processOrderConfirmation(order *Order) {
    // 发送确认邮件
    sendOrderConfirmationEmail(order)
    
    // 扣减库存
    deductInventory(order.Items)
    
    // 创建发货任务
    createShippingTask(order)
}
```

### 数据同步

```go
func syncDataToExternalSystem() error {
    watcher := pie.NewWatcher[User](engine)
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            // 同步到外部系统
            go func() {
                switch change.OperationType {
                case "insert":
                    syncUserToCRM(change.FullDocument)
                case "update":
                    updateUserInCRM(change.FullDocument)
                case "delete":
                    deleteUserFromCRM(change.DocumentKey)
                }
            }()
            
            return nil
        })
}

func syncUserToCRM(user *User) {
    // 同步用户到 CRM 系统
    crmClient := getCRMClient()
    
    crmUser := &CRMUser{
        ID:    user.ID.Hex(),
        Name:  user.Name,
        Email: user.Email,
        Status: user.Status,
    }
    
    if err := crmClient.CreateUser(crmUser); err != nil {
        log.Printf("Failed to sync user to CRM: %v", err)
    }
}
```

### 实时统计更新

```go
func updateRealTimeStatistics() error {
    watcher := pie.NewWatcher[User](engine)
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            // 更新实时统计
            go func() {
                switch change.OperationType {
                case "insert":
                    incrementUserCount()
                    updateUserRegistrationStats(change.FullDocument)
                case "update":
                    updateUserActivityStats(change.FullDocument)
                case "delete":
                    decrementUserCount()
                }
            }()
            
            return nil
        })
}

func incrementUserCount() {
    // 更新用户总数
    stats := getStatistics()
    stats.TotalUsers++
    saveStatistics(stats)
}

func updateUserRegistrationStats(user *User) {
    // 更新注册统计
    stats := getStatistics()
    
    // 按日期统计
    date := user.CreatedAt.Format("2006-01-02")
    stats.DailyRegistrations[date]++
    
    // 按地区统计
    if user.Region != "" {
        stats.RegionalRegistrations[user.Region]++
    }
    
    saveStatistics(stats)
}
```

## 错误处理

### 变更流错误处理

```go
func handleChangeStreamErrors() error {
    watcher := pie.NewWatcher[User](engine)
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            // 处理变更
            if err := processChange(change); err != nil {
                log.Printf("Failed to process change: %v", err)
                
                // 记录错误但不中断流
                logError(change, err)
                
                // 可以返回错误来中断流，或者返回 nil 继续
                return nil
            }
            
            return nil
        })
}

func processChange(change *pie.ChangeEvent[User]) error {
    // 处理变更逻辑
    switch change.OperationType {
    case "insert":
        return handleInsert(change.FullDocument)
    case "update":
        return handleUpdate(change.FullDocument)
    case "delete":
        return handleDelete(change.DocumentKey)
    default:
        return fmt.Errorf("unknown operation type: %s", change.OperationType)
    }
}
```

### 重连机制

```go
func watchWithReconnect() error {
    watcher := pie.NewWatcher[User](engine)
    
    for {
        err := watcher.
            WatchCollection().
            Start(ctx, func(change *pie.ChangeEvent[User]) error {
                return processChange(change)
            })
        
        if err != nil {
            log.Printf("Change stream error: %v", err)
            
            // 等待一段时间后重连
            time.Sleep(5 * time.Second)
            continue
        }
        
        // 正常退出
        break
    }
    
    return nil
}
```

### 优雅关闭

```go
func watchWithGracefulShutdown() error {
    watcher := pie.NewWatcher[User](engine)
    
    // 创建上下文用于取消
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 监听系统信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        log.Println("Received shutdown signal, stopping watcher...")
        cancel()
    }()
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            return processChange(change)
        })
}
```

## 性能优化

### 1. 批量处理

```go
func optimizedBatchProcessing() error {
    watcher := pie.NewWatcher[User](engine)
    
    var changes []*pie.ChangeEvent[User]
    batchSize := 1000
    flushInterval := 10 * time.Second
    
    ticker := time.NewTicker(flushInterval)
    defer ticker.Stop()
    
    go func() {
        for {
            select {
            case <-ticker.C:
                if len(changes) > 0 {
                    processBatch(changes)
                    changes = changes[:0]
                }
            }
        }
    }()
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            changes = append(changes, change)
            
            if len(changes) >= batchSize {
                processBatch(changes)
                changes = changes[:0]
            }
            
            return nil
        })
}
```

### 2. 异步处理

```go
func asyncProcessing() error {
    watcher := pie.NewWatcher[User](engine)
    
    // 创建处理队列
    changeQueue := make(chan *pie.ChangeEvent[User], 1000)
    
    // 启动工作协程
    for i := 0; i < 5; i++ {
        go func() {
            for change := range changeQueue {
                processChange(change)
            }
        }()
    }
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            select {
            case changeQueue <- change:
                // 成功发送到队列
            default:
                // 队列满了，记录错误
                log.Printf("Change queue is full, dropping change")
            }
            
            return nil
        })
}
```

### 3. 过滤优化

```go
func optimizedFiltering() error {
    watcher := pie.NewWatcher[User](engine)
    
    // 使用精确的过滤条件减少不必要的变更
    return watcher.
        WatchCollection().
        Filter(bson.D{
            {"$or", []bson.D{
                {{"operationType", "insert"}},
                {{"operationType", "update"}},
            }},
            {"fullDocument.status", bson.D{{"$in", []string{"active", "pending"}}}},
        }).
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            return processChange(change)
        })
}
```

## 最佳实践

### 1. 错误处理策略

```go
func robustChangeStream() error {
    watcher := pie.NewWatcher[User](engine)
    
    maxRetries := 5
    retryDelay := 1 * time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := watcher.
            WatchCollection().
            Start(ctx, func(change *pie.ChangeEvent[User]) error {
                return processChange(change)
            })
        
        if err == nil {
            return nil // 成功
        }
        
        log.Printf("Change stream failed (attempt %d/%d): %v", i+1, maxRetries, err)
        
        if i < maxRetries-1 {
            time.Sleep(retryDelay)
            retryDelay *= 2 // 指数退避
        }
    }
    
    return fmt.Errorf("failed to start change stream after %d attempts", maxRetries)
}
```

### 2. 监控和日志

```go
func monitoredChangeStream() error {
    watcher := pie.NewWatcher[User](engine)
    
    var processedCount int64
    var errorCount int64
    
    return watcher.
        WatchCollection().
        Start(ctx, func(change *pie.ChangeEvent[User]) error {
            start := time.Now()
            
            err := processChange(change)
            
            duration := time.Since(start)
            atomic.AddInt64(&processedCount, 1)
            
            if err != nil {
                atomic.AddInt64(&errorCount, 1)
                log.Printf("Change processing failed: %v (duration: %v)", err, duration)
            } else {
                log.Printf("Change processed successfully (duration: %v)", duration)
            }
            
            // 定期输出统计信息
            if processedCount%1000 == 0 {
                log.Printf("Processed: %d, Errors: %d", processedCount, errorCount)
            }
            
            return err
        })
}
```

### 3. 资源管理

```go
func managedChangeStream() error {
    watcher := pie.NewWatcher[User](engine)
    
    // 创建带超时的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
    defer cancel()
    
    // 创建信号通道
    done := make(chan error, 1)
    
    go func() {
        done <- watcher.
            WatchCollection().
            Start(ctx, func(change *pie.ChangeEvent[User]) error {
                return processChange(change)
            })
    }()
    
    // 等待完成或超时
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

## 下一步

- [查询作用域](/advanced/scopes/) - 学习作用域定义
- [日志监控](/advanced/logging/) - 掌握查询日志
- [高级聚合](/advanced/aggregation-advanced/) - 深入了解聚合功能
