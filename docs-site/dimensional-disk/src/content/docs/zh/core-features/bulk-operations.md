---
title: 批量操作
description: 学习 Pie 的高效批量写入操作，提升数据处理性能
---

# 批量操作

Pie 提供了高效的批量写入操作，特别适合需要处理大量数据的场景，能够显著提升数据处理性能。

## 基本用法

### 创建批量写入操作

```go
// 创建批量写入操作
bulkWrite := pie.NewBulkWrite[User](engine).
    CollectionForStruct(User{})

// 或者指定集合名称
bulkWrite := pie.NewBulkWrite[User](engine).
    Collection("users")
```

### 插入操作

```go
// 插入单个文档
bulkWrite.InsertOne(&User{
    Name:  "张三",
    Email: "zhangsan@example.com",
    Age:   25,
})

// 插入多个文档
users := []*User{
    {Name: "李四", Email: "lisi@example.com", Age: 30},
    {Name: "王五", Email: "wangwu@example.com", Age: 28},
    {Name: "赵六", Email: "zhaoliu@example.com", Age: 32},
}

for _, user := range users {
    bulkWrite.InsertOne(user)
}
```

### 更新操作

```go
// 更新单个文档
bulkWrite.UpdateOne(
    bson.D{{"email", "old@example.com"}},
    bson.D{{"$set", bson.D{{"email", "new@example.com"}}}},
)

// 更新多个文档
bulkWrite.UpdateMany(
    bson.D{{"status", "inactive"}},
    bson.D{{"$set", bson.D{{"status", "active"}}}},
)

// 使用 Upsert
bulkWrite.UpsertOne(
    bson.D{{"email", "unique@example.com"}},
    &User{
        Name:  "新用户",
        Email: "unique@example.com",
        Age:   25,
    },
)
```

### 删除操作

```go
// 删除单个文档
bulkWrite.DeleteOne(bson.D{{"email", "test@example.com"}})

// 删除多个文档
bulkWrite.DeleteMany(bson.D{{"age", bson.D{{"$lt", 18}}}})

// 删除所有文档（谨慎使用）
bulkWrite.DeleteMany(bson.D{})
```

### 替换操作

```go
// 替换文档
bulkWrite.ReplaceOne(
    bson.D{{"email", "old@example.com"}},
    &User{
        Name:  "新名称",
        Email: "new@example.com",
        Age:   30,
    },
)
```

## 执行批量操作

### 有序执行

```go
// 有序执行（按添加顺序执行，遇到错误会停止）
result, err := bulkWrite.ExecuteOrdered(ctx)
if err != nil {
    log.Fatal("Failed to execute ordered bulk write:", err)
}

fmt.Printf("插入数量: %d\n", result.InsertedCount)
fmt.Printf("更新数量: %d\n", result.ModifiedCount)
fmt.Printf("删除数量: %d\n", result.DeletedCount)
fmt.Printf("匹配数量: %d\n", result.MatchedCount)
fmt.Printf("Upsert数量: %d\n", result.UpsertedCount)
```

### 无序执行

```go
// 无序执行（并行执行，性能更好，但顺序不保证）
result, err := bulkWrite.ExecuteUnordered(ctx)
if err != nil {
    log.Fatal("Failed to execute unordered bulk write:", err)
}
```

## 实际应用场景

### 数据导入

```go
func importUsersFromCSV(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return fmt.Errorf("failed to read CSV: %w", err)
    }
    
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 跳过标题行
    for i, record := range records[1:] {
        if len(record) < 4 {
            log.Printf("跳过第 %d 行：字段不足", i+2)
            continue
        }
        
        age, err := strconv.Atoi(record[2])
        if err != nil {
            log.Printf("跳过第 %d 行：年龄格式错误", i+2)
            continue
        }
        
        user := &User{
            Name:  record[0],
            Email: record[1],
            Age:   age,
            Status: record[3],
        }
        
        bulkWrite.InsertOne(user)
        
        // 每1000条记录执行一次
        if (i+1)%1000 == 0 {
            if _, err := bulkWrite.ExecuteOrdered(context.Background()); err != nil {
                return fmt.Errorf("failed to execute batch at row %d: %w", i+1, err)
            }
            
            // 重新创建批量写入操作
            bulkWrite = pie.NewBulkWrite[User](engine).
                CollectionForStruct(User{})
            
            log.Printf("已处理 %d 条记录", i+1)
        }
    }
    
    // 执行剩余的记录
    if bulkWrite.OperationCount() > 0 {
        if _, err := bulkWrite.ExecuteOrdered(context.Background()); err != nil {
            return fmt.Errorf("failed to execute final batch: %w", err)
        }
    }
    
    log.Printf("导入完成，总共处理了 %d 条记录", len(records)-1)
    return nil
}
```

### 批量更新

```go
func updateUserStatuses() error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 批量更新状态
    bulkWrite.UpdateMany(
        bson.D{{"status", "pending"}},
        bson.D{{"$set", bson.D{{"status", "active"}}}},
    )
    
    // 批量更新年龄组
    bulkWrite.UpdateMany(
        bson.D{{"age", bson.D{{"$gte", 18}}}},
        bson.D{{"$set", bson.D{{"age_group", "adult"}}}},
    )
    
    bulkWrite.UpdateMany(
        bson.D{{"age", bson.D{{"$lt", 18}}}},
        bson.D{{"$set", bson.D{{"age_group", "minor"}}}},
    )
    
    // 批量更新最后登录时间
    bulkWrite.UpdateMany(
        bson.D{{"last_login", bson.D{{"$exists", false}}}},
        bson.D{{"$set", bson.D{{"last_login", time.Now()}}}},
    )
    
    result, err := bulkWrite.ExecuteOrdered(context.Background())
    if err != nil {
        return fmt.Errorf("failed to execute bulk update: %w", err)
    }
    
    log.Printf("批量更新完成：修改了 %d 个文档", result.ModifiedCount)
    return nil
}
```

### 数据清理

```go
func cleanupOldData() error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 删除过期的临时用户
    cutoffDate := time.Now().AddDate(0, -1, 0) // 1个月前
    bulkWrite.DeleteMany(bson.D{
        {"status", "temporary"},
        {"created_at", bson.D{{"$lt", cutoffDate}}},
    })
    
    // 软删除长期未活跃的用户
    inactiveDate := time.Now().AddDate(0, 0, -90) // 90天前
    bulkWrite.UpdateMany(
        bson.D{
            {"last_login", bson.D{{"$lt", inactiveDate}}},
            {"status", "active"},
        },
        bson.D{{"$set", bson.D{{"status", "inactive"}}}},
    )
    
    // 清理测试数据
    bulkWrite.DeleteMany(bson.D{
        {"email", bson.D{{"$regex", "test.*@example\\.com"}}},
    })
    
    result, err := bulkWrite.ExecuteOrdered(context.Background())
    if err != nil {
        return fmt.Errorf("failed to execute cleanup: %w", err)
    }
    
    log.Printf("数据清理完成：删除了 %d 个文档，修改了 %d 个文档", 
        result.DeletedCount, result.ModifiedCount)
    return nil
}
```

## 高级用法

### 条件批量操作

```go
func conditionalBulkOperations() error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 根据条件执行不同的操作
    users, err := session.
        Where("status", "pending").
        Find(context.Background())
    
    if err != nil {
        return err
    }
    
    for _, user := range users {
        if user.Age >= 18 {
            // 成年用户：激活账户
            bulkWrite.UpdateOne(
                bson.D{{"_id", user.ID}},
                bson.D{{"$set", bson.D{{"status", "active"}}}},
            )
        } else {
            // 未成年用户：标记为待审核
            bulkWrite.UpdateOne(
                bson.D{{"_id", user.ID}},
                bson.D{{"$set", bson.D{{"status", "pending_verification"}}}},
            )
        }
    }
    
    result, err := bulkWrite.ExecuteOrdered(context.Background())
    if err != nil {
        return fmt.Errorf("failed to execute conditional operations: %w", err)
    }
    
    log.Printf("条件批量操作完成：修改了 %d 个文档", result.ModifiedCount)
    return nil
}
```

### 批量操作与事务

```go
func bulkOperationsWithTransaction() error {
    return engine.WithTransaction(context.Background(), func(txCtx context.Context) error {
        // 在事务中执行批量操作
        bulkWrite := pie.NewBulkWrite[User](engine).
            CollectionForStruct(User{})
        
        // 添加批量操作
        bulkWrite.InsertOne(&User{Name: "用户1", Email: "user1@example.com"})
        bulkWrite.InsertOne(&User{Name: "用户2", Email: "user2@example.com"})
        
        // 执行批量操作
        result, err := bulkWrite.ExecuteOrdered(txCtx)
        if err != nil {
            return fmt.Errorf("failed to execute bulk write in transaction: %w", err)
        }
        
        log.Printf("事务中批量操作完成：插入了 %d 个文档", result.InsertedCount)
        return nil
    })
}
```

### 错误处理

```go
func bulkOperationsWithErrorHandling() error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 添加操作
    bulkWrite.InsertOne(&User{Name: "用户1", Email: "user1@example.com"})
    bulkWrite.InsertOne(&User{Name: "用户2", Email: "user2@example.com"})
    bulkWrite.UpdateOne(
        bson.D{{"email", "nonexistent@example.com"}},
        bson.D{{"$set", bson.D{{"status", "active"}}}},
    )
    
    result, err := bulkWrite.ExecuteOrdered(context.Background())
    if err != nil {
        // 检查是否是批量写入错误
        if bulkErr, ok := err.(pie.BulkWriteError); ok {
            log.Printf("批量操作部分失败：")
            log.Printf("成功插入: %d", result.InsertedCount)
            log.Printf("成功更新: %d", result.ModifiedCount)
            log.Printf("失败操作: %d", len(bulkErr.WriteErrors))
            
            for _, writeErr := range bulkErr.WriteErrors {
                log.Printf("操作 %d 失败: %s", writeErr.Index, writeErr.Message)
            }
            
            // 根据业务需求决定是否返回错误
            if result.InsertedCount > 0 {
                log.Println("部分操作成功，继续处理")
                return nil
            }
        }
        
        return fmt.Errorf("bulk operation failed: %w", err)
    }
    
    log.Printf("批量操作成功：插入了 %d 个文档", result.InsertedCount)
    return nil
}
```

## 性能优化

### 1. 批量大小优化

```go
const (
    SmallBatchSize  = 100   // 小批量，适合复杂操作
    MediumBatchSize = 1000  // 中等批量，平衡性能和内存
    LargeBatchSize  = 5000  // 大批量，适合简单操作
)

func processDataInBatches(data []User, batchSize int) error {
    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }
        
        batch := data[i:end]
        
        bulkWrite := pie.NewBulkWrite[User](engine).
            CollectionForStruct(User{})
        
        for _, user := range batch {
            bulkWrite.InsertOne(&user)
        }
        
        if _, err := bulkWrite.ExecuteOrdered(context.Background()); err != nil {
            return fmt.Errorf("failed to process batch %d-%d: %w", i, end-1, err)
        }
        
        log.Printf("已处理批次 %d-%d", i, end-1)
    }
    
    return nil
}
```

### 2. 使用无序执行提升性能

```go
// 对于不依赖顺序的操作，使用无序执行
result, err := bulkWrite.ExecuteUnordered(ctx)
```

### 3. 避免重复创建批量写入操作

```go
func efficientBulkOperations() error {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 重用同一个批量写入操作
    for i := 0; i < 10000; i++ {
        bulkWrite.InsertOne(&User{
            Name:  fmt.Sprintf("用户%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        })
        
        // 每1000个操作执行一次
        if (i+1)%1000 == 0 {
            if _, err := bulkWrite.ExecuteOrdered(context.Background()); err != nil {
                return err
            }
            
            // 重新创建批量写入操作
            bulkWrite = pie.NewBulkWrite[User](engine).
                CollectionForStruct(User{})
        }
    }
    
    // 执行剩余操作
    if bulkWrite.OperationCount() > 0 {
        _, err := bulkWrite.ExecuteOrdered(context.Background())
        return err
    }
    
    return nil
}
```

## 监控和调试

### 操作计数

```go
func monitorBulkOperations() {
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 添加操作
    bulkWrite.InsertOne(&User{Name: "用户1"})
    bulkWrite.UpdateOne(bson.D{{"name", "用户2"}}, bson.D{{"$set", bson.D{{"status", "active"}}}})
    bulkWrite.DeleteOne(bson.D{{"name", "用户3"}})
    
    // 检查操作数量
    fmt.Printf("总操作数: %d\n", bulkWrite.OperationCount())
    
    // 执行操作
    result, err := bulkWrite.ExecuteOrdered(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    // 检查结果
    fmt.Printf("插入: %d, 更新: %d, 删除: %d\n", 
        result.InsertedCount, result.ModifiedCount, result.DeletedCount)
}
```

### 性能测量

```go
func measureBulkPerformance() error {
    start := time.Now()
    
    bulkWrite := pie.NewBulkWrite[User](engine).
        CollectionForStruct(User{})
    
    // 添加大量操作
    for i := 0; i < 10000; i++ {
        bulkWrite.InsertOne(&User{
            Name:  fmt.Sprintf("用户%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        })
    }
    
    result, err := bulkWrite.ExecuteOrdered(context.Background())
    if err != nil {
        return err
    }
    
    duration := time.Since(start)
    opsPerSecond := float64(result.InsertedCount) / duration.Seconds()
    
    log.Printf("批量操作完成：")
    log.Printf("操作数量: %d", result.InsertedCount)
    log.Printf("耗时: %v", duration)
    log.Printf("性能: %.2f 操作/秒", opsPerSecond)
    
    return nil
}
```

## 最佳实践

### 1. 合理选择批量大小

```go
// 根据操作复杂度选择批量大小
func chooseBatchSize(operationType string) int {
    switch operationType {
    case "insert":
        return 5000  // 插入操作可以大批量
    case "update":
        return 1000  // 更新操作中等批量
    case "delete":
        return 2000  // 删除操作可以较大批量
    case "upsert":
        return 500   // Upsert操作建议小批量
    default:
        return 1000
    }
}
```

### 2. 错误处理策略

```go
func robustBulkOperations() error {
    maxRetries := 3
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        bulkWrite := pie.NewBulkWrite[User](engine).
            CollectionForStruct(User{})
        
        // 添加操作...
        
        result, err := bulkWrite.ExecuteOrdered(context.Background())
        if err != nil {
            if attempt < maxRetries {
                log.Printf("批量操作失败，第 %d 次重试", attempt)
                time.Sleep(time.Duration(attempt) * time.Second)
                continue
            }
            return fmt.Errorf("批量操作失败，已重试 %d 次: %w", maxRetries, err)
        }
        
        log.Printf("批量操作成功：插入了 %d 个文档", result.InsertedCount)
        return nil
    }
    
    return errors.New("达到最大重试次数")
}
```

### 3. 内存管理

```go
func memoryEfficientBulkOperations() error {
    const batchSize = 1000
    
    // 分批处理，避免内存溢出
    for offset := 0; ; offset += batchSize {
        users, err := session.
            Skip(offset).
            Limit(batchSize).
            Find(context.Background())
        
        if err != nil {
            return err
        }
        
        if len(users) == 0 {
            break // 没有更多数据
        }
        
        // 处理当前批次
        if err := processBatch(users); err != nil {
            return fmt.Errorf("处理批次 %d 失败: %w", offset/batchSize+1, err)
        }
        
        // 强制垃圾回收
        runtime.GC()
    }
    
    return nil
}
```

## 下一步

- [聚合查询](/core-features/aggregation/) - 学习复杂数据聚合
- [事务管理](/core-features/transactions/) - 掌握事务操作
- [缓存支持](/advanced/cache/) - 提升查询性能
