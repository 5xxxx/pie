---
title: Cache Plugin Architecture
description: Learn Pie's flexible plugin-based caching system with Ristretto and Redis support
---

# Cache Plugin Architecture

Pie features a flexible plugin-based caching system that supports multiple cache implementations and chaining. The system provides Ristretto (default) and Redis implementations, while allowing users to implement custom cache plugins.

## Overview

The cache plugin architecture allows you to:
- Use multiple cache instances in a chain
- Implement custom cache backends
- Combine different cache types (memory + Redis)
- Enable automatic cache backfilling
- Monitor cache performance across all layers

## Basic Usage

### Default Ristretto Cache

```go
// Enable default Ristretto memory cache
engine.UseDefaultCache()

// Use cache in session
session := pie.Table[User](engine)
users, err := session.Cache(5 * time.Minute).Find(ctx)
```

### Redis Cache

```go
// Enable Redis cache
redisConfig := &pie.RedisCacheConfig{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    PoolSize: 10,
}
engine.UseRedis(redisConfig)

// Use cache
session := pie.Table[User](engine)
users, err := session.Cache(10 * time.Minute).Find(ctx)
```

### Multi-Level Cache Chain

```go
// Create multiple cache instances
ristrettoCache, _ := pie.NewRistrettoCache(nil)
redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{
    Addr: "localhost:6379",
})

// Use chained caching (L1: Ristretto, L2: Redis)
engine.UseCache(ristrettoCache, redisCache)

// Cache operations will automatically:
// 1. Check L1 cache first
// 2. If miss, check L2 cache
// 3. If L2 hit, backfill to L1
// 4. Write to all cache layers on Set
```

## Cache Configuration

### Ristretto Configuration

```go
// Custom Ristretto configuration
ristrettoConfig := &pie.RistrettoCacheConfig{
    NumCounters: 100000,              // ~10x max entries
    MaxCost:     100 * 1024 * 1024,  // 100MB
    BufferItems: 64,                  // Get buffer size
}

ristrettoCache, err := pie.NewRistrettoCache(ristrettoConfig)
if err != nil {
    log.Fatal(err)
}

engine.UseCache(ristrettoCache)
```

### Redis Configuration

```go
// Redis configuration
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

### Cache Manager Configuration

```go
// Cache manager configuration
config := &pie.CacheConfig{
    Enabled:       true,
    DefaultTTL:    5 * time.Minute,
    KeyPrefix:     "pie:",
    EnableJitter:  true,
    TTLJitter:     2 * time.Minute,
    EmptyCacheTTL: 30 * time.Second,
}

engine.UseCache(ristrettoCache, redisCache)
// Configuration is applied to the cache manager
```

## Custom Cache Implementation

### Implementing the Cache Interface

```go
// Custom cache implementation
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

// Implement Cache interface
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
    // Implement pattern-based deletion
    return nil
}

func (m *MyCache) DeleteByTags(ctx context.Context, tags ...string) error {
    // Implement tag-based deletion
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

### Using Custom Cache

```go
// Use custom cache
myCache := NewMyCache()
engine.UseCache(myCache)

// Or combine with other caches
engine.UseCache(myCache, ristrettoCache, redisCache)
```

## Session Cache Usage

### Basic Caching

```go
// Basic caching
session := pie.Table[User](engine)
users, err := session.Cache(5 * time.Minute).Find(ctx)
```

### Cache with Tags

```go
// Cache with tags for easy invalidation
users, err := session.CacheWithTags("users", "active").Find(ctx)

// Invalidate by tags
engine.Cache().DeleteByTags("users")
```

### Cache with TTL Jitter

```go
// Use TTL jitter to prevent cache stampede
users, err := session.CacheWithJitter(10*time.Minute, 2*time.Minute).Find(ctx)
```

### Cache Empty Results

```go
// Cache empty results to prevent cache penetration
users, err := session.CacheEmpty(30*time.Second).Find(ctx)
```

## Cache Chain Behavior

### Read Operations

When reading from a cache chain:

1. **Sequential Lookup**: Check caches in order (L1 → L2 → L3...)
2. **Backfill**: If found in L2+, automatically backfill to L1
3. **Return**: Return the first found value

```go
// Example: L1 miss, L2 hit, automatic backfill to L1
ristrettoCache, _ := pie.NewRistrettoCache(nil)
redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{Addr: "localhost:6379"})

engine.UseCache(ristrettoCache, redisCache)

// First call: L1 miss, L2 miss, query database
users, err := session.Cache(5*time.Minute).Find(ctx)

// Second call: L1 hit (backfilled from L2)
users, err := session.Cache(5*time.Minute).Find(ctx)
```

### Write Operations

When writing to a cache chain:

1. **Write All**: Write to all cache layers simultaneously
2. **Error Handling**: Continue on individual cache errors
3. **Consistency**: All layers get the same data

```go
// Write to all cache layers
err := session.Cache(5*time.Minute).Set(ctx, "key", data)
// Writes to both Ristretto and Redis
```

### Delete Operations

When deleting from a cache chain:

1. **Delete All**: Delete from all cache layers
2. **Pattern Matching**: Apply pattern deletion to all layers
3. **Tag Invalidation**: Apply tag-based deletion to all layers

```go
// Delete from all layers
engine.Cache().Delete("key")

// Pattern deletion
engine.Cache().DeleteByPattern("users:*")

// Tag-based deletion
engine.Cache().DeleteByTags("users", "active")
```

## Performance Monitoring

### Cache Statistics

```go
// Get aggregated statistics from all cache layers
stats := engine.Cache().Stats()

fmt.Printf("Cache Statistics:")
fmt.Printf("  Total Requests: %d", stats.Total)
fmt.Printf("  Hits: %d", stats.Hits)
fmt.Printf("  Misses: %d", stats.Misses)
fmt.Printf("  Hit Rate: %.2f%%", stats.HitRate)
fmt.Printf("  Keys: %d", stats.Keys)
fmt.Printf("  Size: %d bytes", stats.Size)
fmt.Printf("  Evicted Keys: %d", stats.EvictedKeys)
```

### Individual Cache Statistics

```go
// Get statistics from individual caches
caches := engine.Cache().GetCaches()
for i, cache := range caches {
    stats := cache.Stats()
    fmt.Printf("Cache Layer %d: Hit Rate %.2f%%", i+1, stats.HitRate)
}
```

## Real-World Applications

### User Information Caching

```go
func getUserWithCache(userID bson.ObjectID) (*User, error) {
    session := pie.Table[User](engine)
    
    var user User
    err := session.
        Where("_id", userID).
        Cache(10 * time.Minute).
        First(ctx, &user)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func updateUserWithCache(userID bson.ObjectID, updates bson.D) error {
    session := pie.Table[User](engine)
    
    // Update user
    _, err := session.
        Where("_id", userID).
        Update(ctx, updates)
    
    if err != nil {
        return err
    }
    
    // Clear related cache
    engine.Cache().DeleteByPattern("user:*")
    
    return nil
}
```

### Statistics Caching

```go
func getCachedUserStats() (*UserStats, error) {
    session := pie.Table[User](engine)
    
    var stats UserStats
    
    // Cache total user count
    totalCount, err := session.
        Cache(30 * time.Minute).
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalCount = totalCount
    
    // Cache active user count with tags
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

### Configuration Caching

```go
func getCachedConfig(key string) (string, error) {
    session := pie.Table[Config](engine)
    
    var config Config
    err := session.
        Where("key", key).
        CacheWithTags("config").
        First(ctx, &config)
    
    if err != nil {
        return "", err
    }
    
    return config.Value, nil
}

func setConfigWithCache(key, value string) error {
    session := pie.Table[Config](engine)
    
    // Update or insert configuration
    _, err := session.
        Where("key", key).
        Upsert(ctx, &Config{
            Key:   key,
            Value: value,
        })
    
    if err != nil {
        return err
    }
    
    // Clear config cache by tags
    engine.Cache().DeleteByTags("config")
    
    return nil
}
```

## Advanced Usage

### Cache Warming

```go
func warmupCache() error {
    session := pie.Table[User](engine)
    
    // Warm up common queries
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
        var users []User
        err := q.query().Find(ctx, &users)
        if err != nil {
            log.Printf("Failed to warmup cache for %s: %v", q.name, err)
        } else {
            log.Printf("Warmed up cache for %s: %d users", q.name, len(users))
        }
    }
    
    return nil
}
```

### Conditional Caching

```go
func getUsersWithConditionalCache(useCache bool) ([]User, error) {
    session := pie.Table[User](engine)
    
    query := session.Where("status", "active")
    
    if useCache {
        query = query.Cache(5 * time.Minute)
    }
    
    var users []User
    err := query.Find(ctx, &users)
    return users, err
}
```

### Cache Invalidation Strategies

```go
func invalidateUserCache(userID bson.ObjectID) error {
    cache := engine.Cache()
    
    // Clear specific user cache
    cache.Delete(fmt.Sprintf("user:%s", userID.Hex()))
    
    // Clear related list cache
    cache.DeleteByPattern("users:*")
    cache.DeleteByPattern("active_users:*")
    
    return nil
}

func invalidateAllUserCache() error {
    cache := engine.Cache()
    
    // Clear all user-related cache using tags
    cache.DeleteByTags("users", "user_stats")
    
    return nil
}
```

## Best Practices

### 1. Choose Appropriate Cache Strategy

```go
// Read-heavy, write-light data: long cache time
func getStaticData() ([]StaticData, error) {
    session := pie.Table[StaticData](engine)
    return session.Cache(1 * time.Hour).Find(ctx)
}

// Read-light, write-heavy data: short cache time
func getFrequentlyUpdatedData() ([]DynamicData, error) {
    session := pie.Table[DynamicData](engine)
    return session.Cache(1 * time.Minute).Find(ctx)
}
```

### 2. Use Multi-Level Caching

```go
// L1: Fast memory cache for hot data
// L2: Persistent Redis cache for warm data
ristrettoCache, _ := pie.NewRistrettoCache(&pie.RistrettoCacheConfig{
    NumCounters: 100000,
    MaxCost:     50 * 1024 * 1024, // 50MB
})

redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{
    Addr: "localhost:6379",
})

engine.UseCache(ristrettoCache, redisCache)
```

### 3. Monitor Cache Performance

```go
func monitorCachePerformance() {
    cache := engine.Cache()
    
    // Get cache statistics
    stats := cache.Stats()
    
    log.Printf("Cache Performance:")
    log.Printf("  Hit Rate: %.2f%%", stats.HitRate)
    log.Printf("  Total Requests: %d", stats.Total)
    log.Printf("  Memory Usage: %d bytes", stats.Size)
    
    // Alert if hit rate is too low
    if stats.HitRate < 50.0 {
        log.Printf("WARNING: Low cache hit rate: %.2f%%", stats.HitRate)
    }
}
```

### 4. Use TTL Jitter

```go
// Prevent cache stampede with TTL jitter
session := pie.Table[User](engine)
users, err := session.CacheWithJitter(10*time.Minute, 2*time.Minute).Find(ctx)
```

## Error Handling

### Cache Error Handling

```go
func handleCacheError(err error) {
    if err == nil {
        return
    }
    
    switch err {
    case pie.ErrCacheNotFound:
        log.Println("Cache miss - this is normal")
    case pie.ErrCacheExpired:
        log.Println("Cache expired - refreshing")
    case pie.ErrCacheDisabled:
        log.Println("Cache is disabled")
    default:
        log.Printf("Cache error: %v", err)
    }
}
```

### Fallback Handling

```go
func getUsersWithFallback() ([]User, error) {
    session := pie.Table[User](engine)
    
    var users []User
    
    // Try to get from cache
    err := session.
        Where("status", "active").
        Cache(5 * time.Minute).
        Find(ctx, &users)
    
    if err != nil {
        // Cache miss, fallback to database
        log.Printf("Cache miss, falling back to database: %v", err)
        
        err = session.
            Where("status", "active").
            Find(ctx, &users)
        
        if err != nil {
            return nil, err
        }
    }
    
    return users, nil
}
```

## Next Steps

- [Hook System](/advanced/hooks/) - Learn about lifecycle hooks
- [Soft Delete](/advanced/soft-delete/) - Master soft delete functionality
- [Index Management](/advanced/indexes/) - Learn index optimization
