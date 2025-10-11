---
title: 性能优化
description: 学习 Pie 的性能优化技巧和最佳实践
---

# 性能优化

Pie 提供了多种性能优化功能，帮助您构建高性能的数据库应用程序。

## 连接池优化

### 连接池配置

```go
// 优化连接池配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMaxPoolSize(100),        // 最大连接数
    pie.WithMinPoolSize(10),         // 最小连接数
    pie.WithMaxIdleTime(30*time.Minute), // 最大空闲时间
    pie.WithConnectTimeout(10*time.Second), // 连接超时
    pie.WithSocketTimeout(30*time.Second),  // 套接字超时
)
```

### 连接池监控

```go
func monitorConnectionPool(engine *pie.Engine) {
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                stats := engine.GetConnectionStats()
                log.Printf("Connection Pool Stats:")
                log.Printf("  Active: %d", stats.ActiveConnections)
                log.Printf("  Idle: %d", stats.IdleConnections)
                log.Printf("  Total: %d", stats.TotalConnections)
            }
        }
    }()
}
```

## 查询优化

### 索引优化

```go
// 为常用查询创建索引
func createOptimizedIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 单字段索引
    err := indexes.CreateIndex(ctx, "users", bson.D{{"email", 1}})
    if err != nil {
        return err
    }
    
    // 复合索引
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    })
    if err != nil {
        return err
    }
    
    // 部分索引
    err = indexes.CreateIndex(ctx, "users", bson.D{{"email", 1}}, &options.IndexOptions{
        PartialFilterExpression: bson.D{{"status", "active"}},
    })
    
    return err
}
```

### 查询投影

```go
// 只选择需要的字段
func getUsersOptimized() ([]User, error) {
    session := pie.Table[User](engine)
    
    var users []User
    err := session.
        Select("name", "email", "status"). // 只选择需要的字段
        Where("status", "active").
        Find(ctx, &users)
    
    return users, err
}
```

### 查询限制

```go
// 使用 LIMIT 避免返回过多数据
func getRecentUsers(limit int) ([]User, error) {
    session := pie.Table[User](engine)
    
    var users []User
    err := session.
        Where("status", "active").
        OrderByDesc("created_at").
        Limit(limit).
        Find(ctx, &users)
    
    return users, err
}
```

## 缓存优化

### 多级缓存

```go
// 配置多级缓存
func setupMultiLevelCache() error {
    // L1 缓存（内存）
    memoryCache := pie.NewMemoryCache()
    
    // L2 缓存（Redis）
    redisCache := pie.NewRedisCache("localhost:6379", "", 0)
    
    // 配置二级缓存
    engine.UseTwoLevelCache(
        memoryCache,
        redisCache,
        &pie.TwoLevelCacheConfig{
            L1TTL: 1 * time.Minute,  // L1 缓存 1 分钟
            L2TTL: 10 * time.Minute, // L2 缓存 10 分钟
        },
    )
    
    return nil
}
```

### 缓存预热

```go
func warmupCache() error {
    session := pie.Table[User](engine).WithCache(1 * time.Hour)
    
    // 预热常用查询
    queries := []struct {
        name string
        query func() *pie.Query
    }{
        {
            name: "active_users",
            query: func() *pie.Query {
                return session.Where("status", "active")
            },
        },
        {
            name: "admin_users",
            query: func() *pie.Query {
                return session.Where("role", "admin")
            },
        },
    }
    
    for _, q := range queries {
        var users []User
        err := q.query().Cache(q.name).Find(ctx, &users)
        if err != nil {
            log.Printf("Failed to warmup cache for %s: %v", q.name, err)
        } else {
            log.Printf("Warmed up cache for %s: %d users", q.name, len(users))
        }
    }
    
    return nil
}
```

### 缓存策略

```go
// 根据数据更新频率设置不同的缓存时间
const (
    UserCacheTTL      = 10 * time.Minute  // 用户信息，更新频率中等
    ConfigCacheTTL    = 1 * time.Hour     // 配置信息，更新频率低
    StatsCacheTTL     = 30 * time.Minute  // 统计数据，更新频率中等
    SessionCacheTTL   = 5 * time.Minute   // 会话信息，更新频率高
)

func getCachedData(dataType string) (interface{}, error) {
    session := pie.Table[User](engine)
    
    var ttl time.Duration
    switch dataType {
    case "user":
        ttl = UserCacheTTL
    case "config":
        ttl = ConfigCacheTTL
    case "stats":
        ttl = StatsCacheTTL
    case "session":
        ttl = SessionCacheTTL
    default:
        ttl = 5 * time.Minute
    }
    
    var data interface{}
    err := session.WithCache(ttl).Cache(dataType).Find(ctx, &data)
    return data, err
}
```

## 批量操作优化

### 批量插入

```go
func batchInsertUsers(users []*User) error {
    session := pie.Table[User](engine)
    
    // 分批插入，避免单次操作过大
    batchSize := 1000
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        _, err := session.InsertMany(ctx, batch)
        if err != nil {
            return fmt.Errorf("failed to insert batch %d-%d: %w", i, end-1, err)
        }
    }
    
    return nil
}
```

### 批量更新

```go
func batchUpdateUsers(updates []UserUpdate) error {
    bulkWrite := pie.NewBulkWrite[User](engine)
    
    for _, update := range updates {
        bulkWrite.UpdateOne(
            bson.D{{"_id", update.ID}},
            bson.D{{"$set", update.Data}},
        )
    }
    
    result, err := bulkWrite.ExecuteOrdered(ctx)
    if err != nil {
        return err
    }
    
    log.Printf("Updated %d users", result.ModifiedCount)
    return nil
}
```

## 聚合优化

### 聚合管道优化

```go
func optimizedAggregation() error {
    aggregate := pie.NewAggregate[User](engine)
    
    // 使用 $match 阶段减少数据量
    result, err := aggregate.
        MatchStage().
            Where("status", "active").
            Where("created_at", pie.Gte("created_at", time.Now().AddDate(0, -1, 0))).
        GroupStage().
            By("role", "$role").
            Count("total").
            Avg("avgAge", "$age").
            Done().
        SortStage().Desc("total").
        LimitStage(10). // 限制结果数量
        Exec(ctx)
    
    if err != nil {
        return err
    }
    
    log.Printf("Aggregation result: %+v", result.Data)
    return nil
}
```

### 聚合索引

```go
func createAggregationIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 为聚合查询创建复合索引
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
        {"role", 1},
    })
    if err != nil {
        return err
    }
    
    return nil
}
```

## 内存优化

### 游标使用

```go
func processLargeDataset() error {
    session := pie.Table[User](engine)
    
    cursor, err := session.
        Where("status", "active").
        OrderBy("created_at").
        FindCursor(ctx)
    
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    // 使用游标逐条处理，避免内存溢出
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        
        // 处理单个用户
        processUser(&user)
    }
    
    return nil
}
```

### 分页处理

```go
func processDataInPages(pageSize int) error {
    session := pie.Table[User](engine)
    
    page := 1
    for {
        result, err := session.
            Where("status", "active").
            OrderBy("created_at").
            Paginate(ctx, pie.PaginateParams{
                Page:     page,
                PageSize: pageSize,
            })
        
        if err != nil {
            return err
        }
        
        if len(result.Data) == 0 {
            break // 没有更多数据
        }
        
        // 处理当前页数据
        for _, user := range result.Data {
            processUser(&user)
        }
        
        page++
        
        // 强制垃圾回收
        runtime.GC()
    }
    
    return nil
}
```

## 监控和诊断

### 性能监控

```go
func monitorPerformance() {
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                stats := engine.GetPerformanceStats()
                log.Printf("Performance Stats:")
                log.Printf("  Queries/sec: %.2f", stats.QueriesPerSecond)
                log.Printf("  Avg Query Time: %v", stats.AvgQueryTime)
                log.Printf("  Slow Queries: %d", stats.SlowQueries)
                log.Printf("  Cache Hit Rate: %.2f%%", stats.CacheHitRate)
            }
        }
    }()
}
```

### 慢查询监控

```go
func setupSlowQueryMonitoring() {
    engine.SetSlowQueryThreshold(100 * time.Millisecond)
    
    engine.SetQueryLogFormatter(func(entry *pie.LogEntry) string {
        if entry.Duration > 100*time.Millisecond {
            return fmt.Sprintf("[SLOW QUERY] %s %s - %v", 
                entry.Operation, entry.Collection, entry.Duration)
        }
        return fmt.Sprintf("[QUERY] %s %s - %v", 
            entry.Operation, entry.Collection, entry.Duration)
    })
}
```

### 查询分析

```go
func analyzeQueryPerformance() error {
    session := pie.Table[User](engine)
    
    // 启用查询分析
    session.EnableQueryAnalysis()
    
    start := time.Now()
    
    var users []User
    err := session.
        Where("status", "active").
        Where("age", pie.Gte("age", 18)).
        OrderBy("created_at").
        Find(ctx, &users)
    
    duration := time.Since(start)
    
    if err != nil {
        return err
    }
    
    log.Printf("Query executed in %v", duration)
    log.Printf("Returned %d users", len(users))
    
    // 分析查询计划
    plan := session.GetQueryPlan()
    log.Printf("Query plan: %+v", plan)
    
    return nil
}
```

## 最佳实践

### 1. 查询优化原则

```go
// 好的查询实践
func goodQueryPractices() error {
    session := pie.Table[User](engine)
    
    // 使用索引字段进行查询
    var users []User
    err := session.
        Where("email", "test@example.com"). // email 有索引
        Find(ctx, &users)
    
    if err != nil {
        return err
    }
    
    // 使用投影减少数据传输
    var userNames []string
    err = session.
        Select("name"). // 只选择需要的字段
        Where("status", "active").
        Find(ctx, &userNames)
    
    return err
}

// 避免的查询实践
func badQueryPractices() error {
    session := pie.Table[User](engine)
    
    // 避免全表扫描
    var users []User
    err := session.Find(ctx, &users) // 没有 WHERE 条件
    
    if err != nil {
        return err
    }
    
    // 避免选择所有字段
    var allUsers []User
    err = session.Find(ctx, &allUsers) // 没有 Select 限制
    
    return err
}
```

### 2. 缓存策略

```go
// 合理的缓存策略
func reasonableCachingStrategy() error {
    session := pie.Table[User](engine)
    
    // 静态数据长时间缓存
    var configs []Config
    err := session.
        WithCache(1 * time.Hour).
        Cache("configs").
        Find(ctx, &configs)
    
    if err != nil {
        return err
    }
    
    // 动态数据短时间缓存
    var recentUsers []User
    err = session.
        WithCache(5 * time.Minute).
        Cache("recent_users").
        WhereRecentDays("created_at", 1).
        Find(ctx, &recentUsers)
    
    return err
}
```

### 3. 错误处理

```go
// 性能相关的错误处理
func handlePerformanceErrors(err error) {
    if pie.IsTimeoutError(err) {
        log.Println("Query timeout, consider optimizing query or increasing timeout")
    }
    
    if pie.IsConnectionError(err) {
        log.Println("Connection error, check connection pool settings")
    }
    
    if mongoErr, ok := err.(pie.MongoError); ok {
        switch mongoErr.Code {
        case 12500: // LockTimeout
            log.Println("Lock timeout, consider reducing concurrent operations")
        case 11600: // Interrupted
            log.Println("Operation interrupted, check system resources")
        }
    }
}
```

## 性能测试

### 基准测试

```go
func BenchmarkUserQuery(b *testing.B) {
    engine, err := createTestEngine()
    require.NoError(b, err)
    defer engine.Disconnect(context.Background())
    
    session := pie.Table[User](engine)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var users []User
        err := session.Where("status", "active").Find(context.Background(), &users)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkUserQueryWithCache(b *testing.B) {
    engine, err := createTestEngine()
    require.NoError(b, err)
    defer engine.Disconnect(context.Background())
    
    session := pie.Table[User](engine).WithCache(5 * time.Minute)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var users []User
        err := session.Cache("active_users").Where("status", "active").Find(context.Background(), &users)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 压力测试

```go
func TestConcurrentQueries(t *testing.T) {
    engine, err := createTestEngine()
    require.NoError(t, err)
    defer engine.Disconnect(context.Background())
    
    session := pie.Table[User](engine)
    
    // 并发查询测试
    concurrency := 100
    queries := 1000
    
    var wg sync.WaitGroup
    errors := make(chan error, concurrency)
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            for j := 0; j < queries/concurrency; j++ {
                var users []User
                err := session.Where("status", "active").Find(context.Background(), &users)
                if err != nil {
                    errors <- err
                    return
                }
            }
        }()
    }
    
    wg.Wait()
    close(errors)
    
    // 检查错误
    for err := range errors {
        t.Errorf("Query error: %v", err)
    }
}
```

## 下一步

- [最佳实践](/best-practices/) - 掌握开发最佳实践
- [快速开始](/getting-started/) - 开始使用 Pie
- [查询构建器](/core-features/query-builder/) - 学习查询功能
