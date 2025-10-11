---
title: 缓存支持
description: 学习 Pie 的多级缓存功能，提升查询性能
---

# 缓存支持

Pie 提供了强大的多级缓存功能，支持内存缓存、Redis 缓存和二级缓存，能够显著提升查询性能。

## 基本用法

### 内存缓存

```go
// 启用内存缓存
engine.UseCache(pie.NewMemoryCache(), &pie.CacheConfig{
    TTL: 5 * time.Minute,
})

// 在会话中使用缓存
session := pie.Table[User](engine).WithCache(5 * time.Minute)

// 缓存查询结果
var users []User
err := session.
    Where("status", "active").
    Cache("active_users").
    Find(ctx, &users)
```

### Redis 缓存

```go
// 创建 Redis 缓存
redisCache := pie.NewRedisCache("localhost:6379", "", 0)

// 启用 Redis 缓存
engine.UseCache(redisCache, &pie.CacheConfig{
    TTL: 10 * time.Minute,
})

// 使用缓存
session := pie.Table[User](engine).WithCache(10 * time.Minute)
```

### 二级缓存

```go
// 创建二级缓存
memoryCache := pie.NewMemoryCache()
redisCache := pie.NewRedisCache("localhost:6379", "", 0)

// 启用二级缓存
engine.UseTwoLevelCache(
    memoryCache,  // L1 缓存
    redisCache,   // L2 缓存
    &pie.TwoLevelCacheConfig{
        L1TTL: 1 * time.Minute,  // L1 缓存时间
        L2TTL: 10 * time.Minute, // L2 缓存时间
    },
)
```

## 缓存配置

### 基础配置

```go
// 缓存配置
config := &pie.CacheConfig{
    TTL:           5 * time.Minute,    // 缓存过期时间
    MaxSize:       1000,               // 最大缓存条目数
    CleanupInterval: 10 * time.Minute, // 清理间隔
    EnableMetrics: true,               // 启用指标
}

engine.UseCache(pie.NewMemoryCache(), config)
```

### 高级配置

```go
// 高级缓存配置
config := &pie.CacheConfig{
    TTL:            5 * time.Minute,
    MaxSize:        10000,
    CleanupInterval: 5 * time.Minute,
    EnableMetrics:  true,
    KeyPrefix:      "pie:cache:",      // 键前缀
    Serializer:     &pie.JSONSerializer{}, // 序列化器
    Compressor:     &pie.GzipCompressor{},  // 压缩器
}

engine.UseCache(pie.NewMemoryCache(), config)
```

## 缓存策略

### 查询缓存

```go
// 基础查询缓存
var users []User
err := session.
    Where("status", "active").
    Cache("active_users").
    Find(ctx, &users)

// 带参数的查询缓存
var users []User
err := session.
    Where("status", "active").
    Where("age", pie.Gte("age", 18)).
    Cache("active_adult_users").
    Find(ctx, &users)

// 动态缓存键
cacheKey := fmt.Sprintf("users_by_status_%s", status)
var users []User
err := session.
    Where("status", status).
    Cache(cacheKey).
    Find(ctx, &users)
```

### 分页缓存

```go
// 分页结果缓存
result, err := session.
    Where("status", "active").
    Cache("active_users_page_1").
    Paginate(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
```

### 聚合缓存

```go
// 聚合查询缓存
result, err := aggregate.
    MatchStage().Where("status", "active").
    GroupStage().
        By("role", "$role").
        Count("total").
        Done().
    Cache("user_stats_by_role").
    Exec(ctx)
```

## 实际应用场景

### 用户信息缓存

```go
func getUserWithCache(userID primitive.ObjectID) (*User, error) {
    session := pie.Table[User](engine).WithCache(10 * time.Minute)
    
    var user User
    err := session.
        Where("_id", userID).
        Cache(fmt.Sprintf("user_%s", userID.Hex())).
        First(ctx, &user)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func updateUserWithCache(userID primitive.ObjectID, updates bson.D) error {
    session := pie.Table[User](engine)
    
    // 更新用户
    _, err := session.
        Where("_id", userID).
        Update(ctx, updates)
    
    if err != nil {
        return err
    }
    
    // 清除相关缓存
    cacheKey := fmt.Sprintf("user_%s", userID.Hex())
    engine.Cache().Delete(cacheKey)
    
    // 清除列表缓存
    engine.Cache().DeleteByPattern("users_*")
    
    return nil
}
```

### 统计数据缓存

```go
func getCachedUserStats() (*UserStats, error) {
    session := pie.Table[User](engine).WithCache(30 * time.Minute)
    
    var stats UserStats
    
    // 缓存用户总数
    totalCount, err := session.
        Cache("user_total_count").
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalCount = totalCount
    
    // 缓存活跃用户数
    activeCount, err := session.
        Where("status", "active").
        Cache("user_active_count").
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.ActiveCount = activeCount
    
    // 缓存按角色分组的统计
    var roleStats []bson.M
    err = session.
        GroupStage().
            By("role", "$role").
            Count("count").
            Done().
        Cache("user_stats_by_role").
        Exec(ctx, &roleStats)
    if err != nil {
        return nil, err
    }
    
    stats.RoleStats = roleStats
    return &stats, nil
}
```

### 配置信息缓存

```go
func getCachedConfig(key string) (string, error) {
    session := pie.Table[Config](engine).WithCache(1 * time.Hour)
    
    var config Config
    err := session.
        Where("key", key).
        Cache(fmt.Sprintf("config_%s", key)).
        First(ctx, &config)
    
    if err != nil {
        return "", err
    }
    
    return config.Value, nil
}

func setConfigWithCache(key, value string) error {
    session := pie.Table[Config](engine)
    
    // 更新或插入配置
    _, err := session.
        Where("key", key).
        Upsert(ctx, &Config{
            Key:   key,
            Value: value,
        })
    
    if err != nil {
        return err
    }
    
    // 清除缓存
    engine.Cache().Delete(fmt.Sprintf("config_%s", key))
    
    return nil
}
```

## 高级用法

### 自定义缓存键

```go
// 自定义缓存键生成器
func customCacheKey(query *pie.Query) string {
    // 基于查询条件生成缓存键
    conditions := query.GetConditions()
    hash := sha256.Sum256([]byte(fmt.Sprintf("%v", conditions)))
    return fmt.Sprintf("custom_%x", hash[:8])
}

// 使用自定义缓存键
var users []User
err := session.
    Where("status", "active").
    Where("age", pie.Gte("age", 18)).
    CacheWithKey(customCacheKey).
    Find(ctx, &users)
```

### 条件缓存

```go
// 根据条件决定是否使用缓存
func getUsersWithConditionalCache(useCache bool) ([]User, error) {
    session := pie.Table[User](engine)
    
    query := session.Where("status", "active")
    
    if useCache {
        query = query.Cache("active_users")
    }
    
    var users []User
    err := query.Find(ctx, &users)
    return users, err
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

### 缓存失效

```go
func invalidateUserCache(userID primitive.ObjectID) error {
    cache := engine.Cache()
    
    // 清除特定用户缓存
    cache.Delete(fmt.Sprintf("user_%s", userID.Hex()))
    
    // 清除相关列表缓存
    cache.DeleteByPattern("users_*")
    cache.DeleteByPattern("active_users_*")
    
    return nil
}

func invalidateAllUserCache() error {
    cache := engine.Cache()
    
    // 清除所有用户相关缓存
    cache.DeleteByPattern("user_*")
    cache.DeleteByPattern("users_*")
    
    return nil
}
```

## 性能优化

### 1. 合理设置缓存时间

```go
// 根据数据更新频率设置不同的缓存时间
const (
    UserCacheTTL      = 10 * time.Minute  // 用户信息，更新频率中等
    ConfigCacheTTL    = 1 * time.Hour     // 配置信息，更新频率低
    StatsCacheTTL     = 30 * time.Minute  // 统计数据，更新频率中等
    SessionCacheTTL   = 5 * time.Minute   // 会话信息，更新频率高
)
```

### 2. 使用缓存预热

```go
func startCacheWarmup() {
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if err := warmupCache(); err != nil {
                    log.Printf("Cache warmup failed: %v", err)
                }
            }
        }
    }()
}
```

### 3. 监控缓存性能

```go
func monitorCachePerformance() {
    cache := engine.Cache()
    
    // 获取缓存统计信息
    stats := cache.GetStats()
    
    log.Printf("Cache Stats:")
    log.Printf("  Hits: %d", stats.Hits)
    log.Printf("  Misses: %d", stats.Misses)
    log.Printf("  Hit Rate: %.2f%%", stats.HitRate())
    log.Printf("  Size: %d", stats.Size)
    log.Printf("  Memory Usage: %d bytes", stats.MemoryUsage)
}
```

### 4. 使用压缩减少内存使用

```go
// 启用压缩
config := &pie.CacheConfig{
    TTL:        5 * time.Minute,
    Compressor: &pie.GzipCompressor{},
}

engine.UseCache(pie.NewMemoryCache(), config)
```

## 错误处理

### 缓存错误处理

```go
func handleCacheError(err error) {
    if err == nil {
        return
    }
    
    // 检查是否是缓存错误
    if cacheErr, ok := err.(pie.CacheError); ok {
        switch cacheErr.Type {
        case pie.CacheErrorTypeNotFound:
            log.Println("Cache miss")
        case pie.CacheErrorTypeExpired:
            log.Println("Cache expired")
        case pie.CacheErrorTypeSerialization:
            log.Println("Cache serialization error")
        case pie.CacheErrorTypeDeserialization:
            log.Println("Cache deserialization error")
        default:
            log.Printf("Cache error: %v", err)
        }
        return
    }
    
    log.Printf("Unexpected error: %v", err)
}
```

### 降级处理

```go
func getUsersWithFallback() ([]User, error) {
    session := pie.Table[User](engine).WithCache(5 * time.Minute)
    
    var users []User
    
    // 尝试从缓存获取
    err := session.
        Where("status", "active").
        Cache("active_users").
        Find(ctx, &users)
    
    if err != nil {
        // 缓存失败，从数据库获取
        log.Printf("Cache miss, falling back to database: %v", err)
        
        err = session.
            Where("status", "active").
            Find(ctx, &users)
        
        if err != nil {
            return nil, err
        }
        
        // 异步更新缓存
        go func() {
            session.
                Where("status", "active").
                Cache("active_users").
                Find(context.Background(), &users)
        }()
    }
    
    return users, nil
}
```

## 最佳实践

### 1. 合理选择缓存策略

```go
// 读多写少的数据：长时间缓存
func getStaticData() ([]StaticData, error) {
    session := pie.Table[StaticData](engine).WithCache(1 * time.Hour)
    // ...
}

// 读少写多的数据：短时间缓存或不缓存
func getFrequentlyUpdatedData() ([]DynamicData, error) {
    session := pie.Table[DynamicData](engine).WithCache(1 * time.Minute)
    // ...
}
```

### 2. 缓存键命名规范

```go
// 使用有意义的缓存键
const (
    UserCachePrefix     = "user:"
    UserListCachePrefix = "users:list:"
    UserStatsCachePrefix = "users:stats:"
)

func getUserCacheKey(userID primitive.ObjectID) string {
    return fmt.Sprintf("%s%s", UserCachePrefix, userID.Hex())
}

func getUserListCacheKey(status string, page int) string {
    return fmt.Sprintf("%s%s:page:%d", UserListCachePrefix, status, page)
}
```

### 3. 缓存更新策略

```go
// 更新数据时同时更新缓存
func updateUserWithCache(userID primitive.ObjectID, updates bson.D) error {
    session := pie.Table[User](engine)
    
    // 更新数据库
    _, err := session.
        Where("_id", userID).
        Update(ctx, updates)
    
    if err != nil {
        return err
    }
    
    // 更新缓存
    var user User
    err = session.
        Where("_id", userID).
        First(ctx, &user)
    
    if err == nil {
        cacheKey := getUserCacheKey(userID)
        engine.Cache().Set(cacheKey, user, 10*time.Minute)
    }
    
    // 清除相关列表缓存
    engine.Cache().DeleteByPattern("users:list:*")
    
    return nil
}
```

## 下一步

- [钩子系统](/advanced/hooks/) - 学习生命周期钩子
- [软删除](/advanced/soft-delete/) - 掌握软删除功能
- [索引管理](/advanced/indexes/) - 学习索引优化
