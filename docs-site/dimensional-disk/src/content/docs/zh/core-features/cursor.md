---
title: 游标操作
description: 学习 Pie 的灵活游标操作，支持多种数据遍历方式
---

# 游标操作

Pie 提供了灵活的游标操作，支持多种数据遍历方式，特别适合处理大量数据或需要流式处理的场景。

## 基本用法

### 获取游标

```go
// 创建游标
cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").
    FindCursor(ctx)

if err != nil {
    log.Fatal("Failed to create cursor:", err)
}
defer cursor.Close(ctx)
```

### 方式1: 使用 Next() 和 Decode()

```go
// 逐条处理数据
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        log.Printf("Failed to decode user: %v", err)
        continue
    }
    
    // 处理用户数据
    fmt.Printf("User: %s, Email: %s\n", user.Name, user.Email)
    
    // 可以在这里添加业务逻辑
    processUser(&user)
}
```

### 方式2: 使用 All() 一次性获取

```go
// 一次性获取所有数据
var users []User
err := cursor.All(ctx, &users)
if err != nil {
    log.Fatal("Failed to get all users:", err)
}

// 处理所有用户
for _, user := range users {
    fmt.Printf("User: %s, Email: %s\n", user.Name, user.Email)
    processUser(&user)
}
```

### 方式3: 使用 Iterate() 迭代处理

```go
// 使用迭代函数处理
err := cursor.Iterate(ctx, func(user *User) error {
    // 处理单个用户
    fmt.Printf("User: %s, Email: %s\n", user.Name, user.Email)
    
    // 可以返回错误来中断迭代
    if user.Name == "特殊用户" {
        return errors.New("遇到特殊用户，停止处理")
    }
    
    return nil
})

if err != nil {
    log.Printf("Iteration error: %v", err)
}
```

### 方式4: 使用 Take() 获取前N个

```go
// 获取前5个用户
topUsers, err := cursor.Take(ctx, 5)
if err != nil {
    log.Fatal("Failed to take users:", err)
}

fmt.Printf("获取到 %d 个用户\n", len(topUsers))
for _, user := range topUsers {
    fmt.Printf("User: %s\n", user.Name)
}
```

### 方式5: 使用 First() 获取第一个

```go
// 获取第一个用户
firstUser, err := cursor.First(ctx)
if err != nil {
    if pie.IsNotFoundError(err) {
        fmt.Println("没有找到用户")
    } else {
        log.Fatal("Failed to get first user:", err)
    }
} else {
    fmt.Printf("第一个用户: %s\n", firstUser.Name)
}
```

## 高级用法

### 批量处理

```go
const batchSize = 100

cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").
    FindCursor(ctx)

if err != nil {
    log.Fatal("Failed to create cursor:", err)
}
defer cursor.Close(ctx)

var batch []User
for cursor.Next(ctx) {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    
    batch = append(batch, user)
    
    // 当批次满了或游标结束时处理
    if len(batch) >= batchSize {
        processBatch(batch)
        batch = batch[:0] // 清空批次
    }
}

// 处理剩余的记录
if len(batch) > 0 {
    processBatch(batch)
}

func processBatch(users []User) {
    fmt.Printf("处理批次，包含 %d 个用户\n", len(users))
    
    // 批量处理逻辑
    for _, user := range users {
        // 处理用户
        processUser(&user)
    }
}
```

### 条件中断

```go
cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").
    FindCursor(ctx)

if err != nil {
    log.Fatal("Failed to create cursor:", err)
}
defer cursor.Close(ctx)

processedCount := 0
maxProcess := 1000

for cursor.Next(ctx) && processedCount < maxProcess {
    var user User
    if err := cursor.Decode(&user); err != nil {
        continue
    }
    
    // 处理用户
    processUser(&user)
    processedCount++
    
    // 每处理100个用户输出进度
    if processedCount%100 == 0 {
        fmt.Printf("已处理 %d 个用户\n", processedCount)
    }
}

fmt.Printf("总共处理了 %d 个用户\n", processedCount)
```

### 错误处理和重试

```go
func processUsersWithRetry(session *pie.Session[User], maxRetries int) error {
    cursor, err := session.
        Where("status", "active").
        OrderBy("created_at").
        FindCursor(context.Background())
    
    if err != nil {
        return fmt.Errorf("failed to create cursor: %w", err)
    }
    defer cursor.Close(context.Background())
    
    for cursor.Next(context.Background()) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            log.Printf("Failed to decode user: %v", err)
            continue
        }
        
        // 带重试的处理
        if err := processUserWithRetry(&user, maxRetries); err != nil {
            log.Printf("Failed to process user %s after %d retries: %v", 
                user.Name, maxRetries, err)
            continue
        }
    }
    
    return nil
}

func processUserWithRetry(user *User, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := processUser(user); err != nil {
            if i == maxRetries-1 {
                return err
            }
            
            // 等待一段时间后重试
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        return nil
    }
    return errors.New("max retries exceeded")
}
```

## 实际应用场景

### 数据迁移

```go
func migrateUsers() error {
    // 源数据库游标
    sourceCursor, err := sourceSession.
        Where("status", "active").
        OrderBy("created_at").
        FindCursor(context.Background())
    
    if err != nil {
        return fmt.Errorf("failed to create source cursor: %w", err)
    }
    defer sourceCursor.Close(context.Background())
    
    // 目标数据库会话
    targetSession := pie.Table[User](targetEngine)
    
    batchSize := 100
    var batch []User
    processedCount := 0
    
    for sourceCursor.Next(context.Background()) {
        var user User
        if err := sourceCursor.Decode(&user); err != nil {
            log.Printf("Failed to decode user: %v", err)
            continue
        }
        
        // 转换数据
        migratedUser := migrateUserData(user)
        batch = append(batch, migratedUser)
        
        // 批量插入
        if len(batch) >= batchSize {
            if err := insertBatch(targetSession, batch); err != nil {
                return fmt.Errorf("failed to insert batch: %w", err)
            }
            processedCount += len(batch)
            batch = batch[:0]
            log.Printf("已迁移 %d 个用户", processedCount)
        }
    }
    
    // 处理剩余数据
    if len(batch) > 0 {
        if err := insertBatch(targetSession, batch); err != nil {
            return fmt.Errorf("failed to insert final batch: %w", err)
        }
        processedCount += len(batch)
    }
    
    log.Printf("迁移完成，总共迁移了 %d 个用户", processedCount)
    return nil
}

func insertBatch(session *pie.Session[User], users []User) error {
    bulkWrite := pie.NewBulkWrite[User](session.Engine())
    for _, user := range users {
        bulkWrite.InsertOne(&user)
    }
    _, err := bulkWrite.ExecuteOrdered(context.Background())
    return err
}
```

### 数据导出

```go
func exportUsersToCSV(w io.Writer) error {
    cursor, err := session.
        Where("status", "active").
        OrderBy("created_at").
        FindCursor(context.Background())
    
    if err != nil {
        return fmt.Errorf("failed to create cursor: %w", err)
    }
    defer cursor.Close(context.Background())
    
    // 写入CSV头部
    csvWriter := csv.NewWriter(w)
    defer csvWriter.Flush()
    
    if err := csvWriter.Write([]string{"ID", "Name", "Email", "Age", "CreatedAt"}); err != nil {
        return err
    }
    
    // 逐行写入数据
    for cursor.Next(context.Background()) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            log.Printf("Failed to decode user: %v", err)
            continue
        }
        
        record := []string{
            user.ID.Hex(),
            user.Name,
            user.Email,
            strconv.Itoa(user.Age),
            user.CreatedAt.Format(time.RFC3339),
        }
        
        if err := csvWriter.Write(record); err != nil {
            return err
        }
    }
    
    return nil
}
```

### 实时数据处理

```go
func processRealTimeData() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := processPendingData(); err != nil {
                log.Printf("Failed to process pending data: %v", err)
            }
        }
    }
}

func processPendingData() error {
    cursor, err := session.
        Where("status", "pending").
        Where("created_at", pie.Gte("created_at", time.Now().Add(-1*time.Hour))).
        OrderBy("created_at").
        FindCursor(context.Background())
    
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())
    
    processedCount := 0
    for cursor.Next(context.Background()) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        // 处理用户
        if err := processUser(&user); err != nil {
            log.Printf("Failed to process user %s: %v", user.Name, err)
            continue
        }
        
        // 更新状态
        if err := session.Where("_id", user.ID).Update(context.Background(), 
            bson.D{{"$set", bson.D{{"status", "processed"}}}}); err != nil {
            log.Printf("Failed to update user status: %v", err)
        }
        
        processedCount++
    }
    
    if processedCount > 0 {
        log.Printf("处理了 %d 个待处理用户", processedCount)
    }
    
    return nil
}
```

## 性能优化

### 1. 使用投影减少数据传输

```go
// 只选择需要的字段
cursor, err := session.
    Select("name", "email", "status").
    Where("status", "active").
    FindCursor(ctx)
```

### 2. 合理设置批次大小

```go
// 根据内存和性能需求调整批次大小
const (
    SmallBatchSize  = 50   // 内存受限环境
    MediumBatchSize = 100  // 一般场景
    LargeBatchSize  = 500  // 高性能场景
)
```

### 3. 使用索引优化

```go
// 确保排序字段有索引
cursor, err := session.
    Where("status", "active").
    OrderBy("created_at").  // created_at 应该有索引
    FindCursor(ctx)
```

### 4. 避免在游标中执行复杂操作

```go
// 好的做法：先获取数据，再处理
var users []User
err := cursor.All(ctx, &users)
if err != nil {
    return err
}

// 在内存中处理数据
for _, user := range users {
    processUser(&user)
}

// 避免的做法：在游标循环中执行复杂操作
for cursor.Next(ctx) {
    var user User
    cursor.Decode(&user)
    
    // 避免在游标中执行复杂查询
    relatedData, err := getRelatedData(user.ID) // 这会导致N+1查询问题
}
```

## 错误处理

### 游标错误处理

```go
func processUsersSafely() error {
    cursor, err := session.
        Where("status", "active").
        FindCursor(context.Background())
    
    if err != nil {
        return fmt.Errorf("failed to create cursor: %w", err)
    }
    defer func() {
        if closeErr := cursor.Close(context.Background()); closeErr != nil {
            log.Printf("Failed to close cursor: %v", closeErr)
        }
    }()
    
    for cursor.Next(context.Background()) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            log.Printf("Failed to decode user: %v", err)
            continue
        }
        
        if err := processUser(&user); err != nil {
            log.Printf("Failed to process user %s: %v", user.Name, err)
            // 根据业务需求决定是否继续
            continue
        }
    }
    
    // 检查游标错误
    if err := cursor.Err(); err != nil {
        return fmt.Errorf("cursor error: %w", err)
    }
    
    return nil
}
```

## 最佳实践

### 1. 总是关闭游标

```go
cursor, err := session.FindCursor(ctx)
if err != nil {
    return err
}
defer cursor.Close(ctx) // 确保游标被关闭
```

### 2. 处理游标错误

```go
for cursor.Next(ctx) {
    // 处理数据
}

// 检查游标是否有错误
if err := cursor.Err(); err != nil {
    log.Printf("Cursor error: %v", err)
}
```

### 3. 合理使用内存

```go
// 对于大数据集，使用流式处理而不是一次性加载
for cursor.Next(ctx) {
    var user User
    cursor.Decode(&user)
    processUser(&user) // 处理完立即释放内存
}

// 而不是
users, err := cursor.All(ctx, &users) // 可能消耗大量内存
```

## 下一步

- [批量操作](/core-features/bulk-operations/) - 学习高效的批量写入
- [聚合查询](/core-features/aggregation/) - 掌握复杂数据聚合
- [缓存支持](/advanced/cache/) - 提升查询性能
