---
title: 缓存插件架构
description: 学习 Pie 灵活的插件式缓存系统，支持 Ristretto 和 Redis
---

# 缓存插件架构

Pie 采用灵活的插件式缓存系统，支持多种缓存实现和链式组合。系统提供 Ristretto（默认）和 Redis 实现，同时允许用户实现自定义缓存插件。

## 概述

缓存插件架构允许您：
- 在链中使用多个缓存实例
- 实现自定义缓存后端
- 组合不同的缓存类型（内存 + Redis）
- 启用自动缓存回填
- 监控所有层级的缓存性能

## 基本用法

### 默认 Ristretto 缓存

```go
// 启用默认 Ristretto 内存缓存
engine.UseDefaultCache()

// 在会话中使用缓存
session := pie.Table[User](engine)
users, err := session.Cache(5 * time.Minute).Find(ctx)
```

### Redis 缓存

```go
// 启用 Redis 缓存
redisConfig := &pie.RedisCacheConfig{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    PoolSize: 10,
}
engine.UseRedis(redisConfig)

// 使用缓存
session := pie.Table[User](engine)
users, err := session.Cache(10 * time.Minute).Find(ctx)
```

### 多层缓存链

```go
// 创建多个缓存实例
ristrettoCache, _ := pie.NewRistrettoCache(nil)
redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{
    Addr: "localhost:6379",
})

// 使用链式缓存（L1: Ristretto, L2: Redis）
engine.UseCache(ristrettoCache, redisCache)

// 缓存操作将自动：
// 1. 首先检查 L1 缓存
// 2. 如果未命中，检查 L2 缓存
// 3. 如果 L2 命中，回填到 L1
// 4. Set 时写入所有缓存层
```

## 缓存配置

### Ristretto 配置

```go
// 自定义 Ristretto 配置
ristrettoConfig := &pie.RistrettoCacheConfig{
    NumCounters: 100000,              // ~10x 最大条目数
    MaxCost:     100 * 1024 * 1024,   // 100MB
    BufferItems: 64,                  // Get buffer 大小
}

ristrettoCache, err := pie.NewRistrettoCache(ristrettoConfig)
if err != nil {
    log.Fatal(err)
}

engine.UseCache(ristrettoCache)
```

### Redis 配置

```go
// Redis 配置
redisConfig := &pie.RedisCacheConfig{
    Addr:     "localhost:6379",
    Password: "your-password",
    DB:       0,
    PoolSize: 20,
}

redisCache, err := pie.NewRedisCache(redisConfig)
if err != nil {
    log.Fatal(err)
}

engine.UseCache(redisCache)
```

### 缓存管理器配置

```go
// 缓存管理器配置
config := &pie.CacheConfig{
    Enabled:       true,
    DefaultTTL:    5 * time.Minute,
    KeyPrefix:     "pie:",
    EnableJitter:  true,
    TTLJitter:     2 * time.Minute,
    EmptyCacheTTL: 30 * time.Second,
}

engine.UseCache(ristrettoCache, redisCache)
// 配置将应用到缓存管理器
```

## 自定义缓存实现

### 实现 Cache 接口

```go
// 自定义缓存实现
type MyCache struct {
    data  map[string][]byte
    stats *pie.CacheStats
    mu    sync.RWMutex
}

func NewMyCache() *MyCache {
    return &MyCache{
        data:  make(map[string][]byte),
        stats: &pie.CacheStats{},
    }
}

// 实现 Cache 接口
func (m *MyCache) Get(ctx context.Context, key string) ([]byte, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    m.stats.Total++
    if val, exists := m.data[key]; exists {
        m.stats.Hits++
        m.stats.HitRate = float64(m.stats.Hits) / float64(m.stats.Total) * 100
        return val, nil
    }
    
    m.stats.Misses++
    return nil, pie.ErrCacheNotFound
}

func (m *MyCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.data[key] = value
    m.stats.Keys++
    return nil
}

func (m *MyCache) Delete(ctx context.Context, key string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    delete(m.data, key)
    m.stats.Keys--
    return nil
}

func (m *MyCache) DeleteByPattern(ctx context.Context, pattern string) error {
    // 实现基于模式的删除
    return nil
}

func (m *MyCache) DeleteByTags(ctx context.Context, tags ...string) error {
    // 实现基于标签的删除
    return nil
}

func (m *MyCache) Exists(ctx context.Context, key string) (bool, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    _, exists := m.data[key]
    return exists, nil
}

func (m *MyCache) Clear(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.data = make(map[string][]byte)
    m.stats.Keys = 0
    return nil
}

func (m *MyCache) Stats() *pie.CacheStats {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    stats := *m.stats
    return &stats
}
```

### 使用自定义缓存

```go
// 使用自定义缓存
myCache := NewMyCache()
engine.UseCache(myCache)

// 或与其他缓存组合
engine.UseCache(myCache, ristrettoCache, redisCache)
```

## 会话缓存使用

### 基础缓存

```go
// 基础缓存
session := pie.Table[User](engine)
users, err := session.Cache(5 * time.Minute).Find(ctx)
```

### 带标签的缓存

```go
// 带标签的缓存，便于失效
users, err := session.CacheWithTags("users", "active").Find(ctx)

// 按标签失效
_ = engine.Cache().DeleteByTags(ctx, "users")
```

### 带 TTL 抖动的缓存

```go
// 使用 TTL 抖动防止缓存雪崩
users, err := session.CacheWithJitter(10*time.Minute, 2*time.Minute).Find(ctx)
```

### 缓存空结果

```go
// 缓存空结果防止缓存穿透
users, err := session.CacheEmpty(30*time.Second).Find(ctx)
```

## 缓存链行为

### 读操作

从缓存链读取时：

1. **顺序查找**：按顺序检查缓存（L1 → L2 → L3...）
2. **回填**：如果在 L2+ 找到，自动回填到 L1
3. **返回**：返回第一个找到的值

```go
// 示例：L1 未命中，L2 命中，自动回填到 L1
ristrettoCache, _ := pie.NewRistrettoCache(nil)
redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{Addr: "localhost:6379"})

engine.UseCache(ristrettoCache, redisCache)

// 第一次调用：L1 未命中，L2 未命中，查询数据库
users, err := session.Cache(5*time.Minute).Find(ctx)

// 第二次调用：L1 命中（从 L2 回填）
users, err := session.Cache(5*time.Minute).Find(ctx)
```

### 写操作

当你在启用缓存的情况下执行查询（例如 `Find`、`First`、`Count`）时，Pie 会自动序列化结果并写入每一层缓存，大多数场景下无需手动写入。

如果需要缓存非查询的数据，可以直接操作缓存管理器：

```go
// 手动向所有缓存层写入数据
payload, _ := json.Marshal(users)
err := engine.Cache().Set(ctx, "users:active", payload, 5*time.Minute)
```

### 删除操作

从缓存链删除时：

1. **删除全部**：从所有缓存层删除
2. **模式匹配**：对所有层应用模式删除
3. **标签失效**：对所有层应用基于标签的删除

```go
// 从所有层删除
_ = engine.Cache().Delete(ctx, "users:active")

// 模式删除
_ = engine.Cache().DeleteByPattern(ctx, "users:*")

// 基于标签的删除
_ = engine.Cache().DeleteByTags(ctx, "users", "active")
```

## 性能监控

### 缓存统计

```go
// 获取所有缓存层的聚合统计
stats := engine.Cache().Stats()

fmt.Printf("缓存统计:")
fmt.Printf("  总请求数: %d", stats.Total)
fmt.Printf("  命中数: %d", stats.Hits)
fmt.Printf("  未命中数: %d", stats.Misses)
fmt.Printf("  命中率: %.2f%%", stats.HitRate)
fmt.Printf("  键数量: %d", stats.Keys)
fmt.Printf("  大小: %d 字节", stats.Size)
fmt.Printf("  被驱逐的键: %d", stats.EvictedKeys)
```

### 单个缓存统计

```go
// 获取单个缓存的统计
caches := engine.Cache().GetCaches()
for i, cache := range caches {
    stats := cache.Stats()
    fmt.Printf("缓存层 %d: 命中率 %.2f%%", i+1, stats.HitRate)
}
```

## 实际应用场景

### 用户信息缓存

```go
func getUserWithCache(userID bson.ObjectID) (*User, error) {
    session := pie.Table[User](engine)

    user, err := session.
        Where("_id", userID).
        Cache(10 * time.Minute).
        FindOne(ctx)

    if err != nil {
        return nil, err
    }

    return user, nil
}

func updateUserWithCache(userID bson.ObjectID, updates bson.D) error {
    session := pie.Table[User](engine)
    
    // 更新用户
    _, err := session.
        Where("_id", userID).
        Update(ctx, updates)
    
    if err != nil {
        return err
    }
    
    // 清除相关缓存
    _ = engine.Cache().DeleteByPattern(ctx, "user:*")
    
    return nil
}
```

### 统计数据缓存

```go
func getCachedUserStats() (*UserStats, error) {
    session := pie.Table[User](engine)
    
    var stats UserStats
    
    // 缓存用户总数
    totalCount, err := session.
        Cache(30 * time.Minute).
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalCount = totalCount
    
    // 缓存活跃用户数（带标签）
    activeCount, err := session.
        Where("status", "active").
        CacheWithTags("users", "stats").
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.ActiveCount = activeCount
    
    return &stats, nil
}
```

### 配置信息缓存

```go
func getCachedConfig(key string) (string, error) {
    session := pie.Table[Config](engine)

    config, err := session.
        Where("key", key).
        CacheWithTags("config").
        FindOne(ctx)

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
    
    // 按标签清除配置缓存
    _ = engine.Cache().DeleteByTags(ctx, "config")
    
    return nil
}
```

## 高级用法

### 缓存预热

```go
func warmupCache() error {
    session := pie.Table[User](engine)
    
    // 预热常用查询
    queries := []struct {
        name  string
        query func() *pie.Session[User]
    }{
        {
            name: "active_users",
            query: func() *pie.Session[User] {
                return session.Where("status", "active").Cache(1*time.Hour)
            },
        },
        {
            name: "admin_users",
            query: func() *pie.Session[User] {
                return session.Where("role", "admin").Cache(1*time.Hour)
            },
        },
    }
    
    for _, q := range queries {
        users, err := q.query().Find(ctx)
        if err != nil {
            log.Printf("预热缓存失败 %s: %v", q.name, err)
        } else {
            log.Printf("预热缓存成功 %s: %d 用户", q.name, len(users))
        }
    }
    
    return nil
}
```

### 条件缓存

```go
// 根据条件决定是否使用缓存
func getUsersWithConditionalCache(useCache bool) ([]User, error) {
    session := pie.Table[User](engine)
    
    query := session.Where("status", "active")
    
    if useCache {
        query = query.Cache(5 * time.Minute)
    }
    
    users, err := query.Find(ctx)
    return users, err
}
```

### 缓存失效策略

```go
func invalidateUserCache(userID bson.ObjectID) error {
    cache := engine.Cache()
    
    // 清除特定用户缓存
    _ = cache.Delete(ctx, fmt.Sprintf("user:%s", userID.Hex()))

    // 清除相关列表缓存
    _ = cache.DeleteByPattern(ctx, "users:*")
    _ = cache.DeleteByPattern(ctx, "active_users:*")
    
    return nil
}

func invalidateAllUserCache() error {
    cache := engine.Cache()
    
    // 使用标签清除所有用户相关缓存
    _ = cache.DeleteByTags(ctx, "users", "user_stats")
    
    return nil
}
```

## 最佳实践

### 1. 选择合适的缓存策略

```go
// 读多写少的数据：长时间缓存
func getStaticData() ([]StaticData, error) {
    session := pie.Table[StaticData](engine)
    return session.Cache(1 * time.Hour).Find(ctx)
}

// 读少写多的数据：短时间缓存
func getFrequentlyUpdatedData() ([]DynamicData, error) {
    session := pie.Table[DynamicData](engine)
    return session.Cache(1 * time.Minute).Find(ctx)
}
```

### 2. 使用多层缓存

```go
// L1：快速内存缓存用于热数据
// L2：持久化 Redis 缓存用于温数据
ristrettoCache, _ := pie.NewRistrettoCache(&pie.RistrettoCacheConfig{
    NumCounters: 100000,
    MaxCost:     50 * 1024 * 1024, // 50MB
})

redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{
    Addr: "localhost:6379",
})

engine.UseCache(ristrettoCache, redisCache)
```

### 3. 监控缓存性能

```go
func monitorCachePerformance() {
    cache := engine.Cache()
    
    // 获取缓存统计
    stats := cache.Stats()
    
    log.Printf("缓存性能:")
    log.Printf("  命中率: %.2f%%", stats.HitRate)
    log.Printf("  总请求数: %d", stats.Total)
    log.Printf("  内存使用: %d 字节", stats.Size)
    
    // 命中率过低时告警
    if stats.HitRate < 50.0 {
        log.Printf("警告: 缓存命中率过低: %.2f%%", stats.HitRate)
    }
}
```

### 4. 使用 TTL 抖动

```go
// 使用 TTL 抖动防止缓存雪崩
session := pie.Table[User](engine)
users, err := session.CacheWithJitter(10*time.Minute, 2*time.Minute).Find(ctx)
```

## 错误处理

### 缓存错误处理

```go
func handleCacheError(err error) {
    if err == nil {
        return
    }
    
    switch err {
    case pie.ErrCacheNotFound:
        log.Println("缓存未命中 - 这是正常的")
    case pie.ErrCacheExpired:
        log.Println("缓存已过期 - 正在刷新")
    case pie.ErrCacheDisabled:
        log.Println("缓存已禁用")
    default:
        log.Printf("缓存错误: %v", err)
    }
}
```

### 降级处理

```go
func getUsersWithFallback() ([]User, error) {
    session := pie.Table[User](engine)

    // 尝试从缓存获取
    users, err := session.
        Where("status", "active").
        Cache(5 * time.Minute).
        Find(ctx)

    if err != nil {
        // 缓存未命中，降级到数据库
        log.Printf("缓存未命中，降级到数据库: %v", err)

        users, err = session.
            Where("status", "active").
            Find(ctx)

        if err != nil {
            return nil, err
        }
    }
    
    return users, nil
}
```

## 下一步

- [钩子系统](/advanced/hooks/) - 学习生命周期钩子
- [软删除](/advanced/soft-delete/) - 掌握软删除功能
- [索引管理](/advanced/indexes/) - 学习索引优化
