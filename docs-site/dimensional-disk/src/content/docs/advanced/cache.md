---
title: Cache Support
description: Learn Pie's multi-level caching features to improve query performance
---

# Cache Support

Pie provides powerful multi-level caching functionality, supporting memory cache, Redis cache, and two-level cache, which can significantly improve query performance.

## Basic Usage

### Memory Cache

```go
// Enable memory cache
engine.UseCache(pie.NewMemoryCache(), &pie.CacheConfig{
    TTL: 5 * time.Minute,
})

// Use cache in session
session := pie.Table[User](engine).WithCache(5 * time.Minute)

// Cache query results
var users []User
err := session.
    Where("status", "active").
    Cache("active_users").
    Find(ctx, &users)
```

### Redis Cache

```go
// Create Redis cache
redisCache := pie.NewRedisCache("localhost:6379", "", 0)

// Enable Redis cache
engine.UseCache(redisCache, &pie.CacheConfig{
    TTL: 10 * time.Minute,
})

// Use cache
session := pie.Table[User](engine).WithCache(10 * time.Minute)
```

### Two-Level Cache

```go
// Create two-level cache
memoryCache := pie.NewMemoryCache()
redisCache := pie.NewRedisCache("localhost:6379", "", 0)

// Enable two-level cache
engine.UseTwoLevelCache(
    memoryCache,  // L1 cache
    redisCache,   // L2 cache
    &pie.TwoLevelCacheConfig{
        L1TTL: 1 * time.Minute,  // L1 cache time
        L2TTL: 10 * time.Minute, // L2 cache time
    },
)
```

## Cache Configuration

### Basic Configuration

```go
// Cache configuration
config := &pie.CacheConfig{
    TTL:            5 * time.Minute,    // Cache expiration time
    MaxSize:        1000,               // Maximum cache entries
    CleanupInterval: 10 * time.Minute,  // Cleanup interval
    EnableMetrics:  true,               // Enable metrics
}

engine.UseCache(pie.NewMemoryCache(), config)
```

### Advanced Configuration

```go
// Advanced cache configuration
config := &pie.CacheConfig{
    TTL:            5 * time.Minute,
    MaxSize:        10000,
    CleanupInterval: 5 * time.Minute,
    EnableMetrics:  true,
    KeyPrefix:      "pie:cache:",           // Key prefix
    Serializer:     &pie.JSONSerializer{},  // Serializer
    Compressor:     &pie.GzipCompressor{},  // Compressor
}

engine.UseCache(pie.NewMemoryCache(), config)
```

## Cache Strategies

### Query Caching

```go
// Basic query caching
var users []User
err := session.
    Where("status", "active").
    Cache("active_users").
    Find(ctx, &users)

// Query caching with parameters
var users []User
err := session.
    Where("status", "active").
    Where("age", pie.Gte("age", 18)).
    Cache("active_adult_users").
    Find(ctx, &users)

// Dynamic cache keys
cacheKey := fmt.Sprintf("users_by_status_%s", status)
var users []User
err := session.
    Where("status", status).
    Cache(cacheKey).
    Find(ctx, &users)
```

### Pagination Caching

```go
// Pagination result caching
result, err := session.
    Where("status", "active").
    Cache("active_users_page_1").
    Paginate(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
```

### Aggregation Caching

```go
// Aggregation query caching
result, err := aggregate.
    MatchStage().Where("status", "active").
    GroupStage().
        By("role", "$role").
        Count("total").
        Done().
    Cache("user_stats_by_role").
    Exec(ctx)
```

## Real-World Applications

### User Information Caching

```go
func getUserWithCache(userID bson.ObjectID) (*User, error) {
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
    cacheKey := fmt.Sprintf("user_%s", userID.Hex())
    engine.Cache().Delete(cacheKey)
    
    // Clear list cache
    engine.Cache().DeleteByPattern("users_*")
    
    return nil
}
```

### Statistics Caching

```go
func getCachedUserStats() (*UserStats, error) {
    session := pie.Table[User](engine).WithCache(30 * time.Minute)
    
    var stats UserStats
    
    // Cache total user count
    totalCount, err := session.
        Cache("user_total_count").
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalCount = totalCount
    
    // Cache active user count
    activeCount, err := session.
        Where("status", "active").
        Cache("user_active_count").
        Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.ActiveCount = activeCount
    
    // Cache role-based statistics
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

### Configuration Caching

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
    
    // Clear cache
    engine.Cache().Delete(fmt.Sprintf("config_%s", key))
    
    return nil
}
```

## Advanced Usage

### Custom Cache Keys

```go
// Custom cache key generator
func customCacheKey(query *pie.Query) string {
    // Generate cache key based on query conditions
    conditions := query.GetConditions()
    hash := sha256.Sum256([]byte(fmt.Sprintf("%v", conditions)))
    return fmt.Sprintf("custom_%x", hash[:8])
}

// Use custom cache key
var users []User
err := session.
    Where("status", "active").
    Where("age", pie.Gte("age", 18)).
    CacheWithKey(customCacheKey).
    Find(ctx, &users)
```

### Conditional Caching

```go
// Decide whether to use cache based on conditions
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

### Cache Warming

```go
func warmupCache() error {
    session := pie.Table[User](engine).WithCache(1 * time.Hour)
    
    // Warm up common queries
    queries := []struct {
        name  string
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

### Cache Invalidation

```go
func invalidateUserCache(userID bson.ObjectID) error {
    cache := engine.Cache()
    
    // Clear specific user cache
    cache.Delete(fmt.Sprintf("user_%s", userID.Hex()))
    
    // Clear related list cache
    cache.DeleteByPattern("users_*")
    cache.DeleteByPattern("active_users_*")
    
    return nil
}

func invalidateAllUserCache() error {
    cache := engine.Cache()
    
    // Clear all user-related cache
    cache.DeleteByPattern("user_*")
    cache.DeleteByPattern("users_*")
    
    return nil
}
```

## Performance Optimization

### 1. Reasonable Cache Time Settings

```go
// Set different cache times based on data update frequency
const (
    UserCacheTTL      = 10 * time.Minute  // User info, medium update frequency
    ConfigCacheTTL    = 1 * time.Hour     // Config info, low update frequency
    StatsCacheTTL     = 30 * time.Minute  // Statistics, medium update frequency
    SessionCacheTTL   = 5 * time.Minute   // Session info, high update frequency
)
```

### 2. Use Cache Warming

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

### 3. Monitor Cache Performance

```go
func monitorCachePerformance() {
    cache := engine.Cache()
    
    // Get cache statistics
    stats := cache.GetStats()
    
    log.Printf("Cache Stats:")
    log.Printf("  Hits: %d", stats.Hits)
    log.Printf("  Misses: %d", stats.Misses)
    log.Printf("  Hit Rate: %.2f%%", stats.HitRate())
    log.Printf("  Size: %d", stats.Size)
    log.Printf("  Memory Usage: %d bytes", stats.MemoryUsage)
}
```

### 4. Use Compression to Reduce Memory Usage

```go
// Enable compression
config := &pie.CacheConfig{
    TTL:        5 * time.Minute,
    Compressor: &pie.GzipCompressor{},
}

engine.UseCache(pie.NewMemoryCache(), config)
```

## Error Handling

### Cache Error Handling

```go
func handleCacheError(err error) {
    if err == nil {
        return
    }
    
    // Check if it's a cache error
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

### Fallback Handling

```go
func getUsersWithFallback() ([]User, error) {
    session := pie.Table[User](engine).WithCache(5 * time.Minute)
    
    var users []User
    
    // Try to get from cache
    err := session.
        Where("status", "active").
        Cache("active_users").
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
        
        // Asynchronously update cache
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

## Best Practices

### 1. Choose Appropriate Cache Strategy

```go
// Read-heavy, write-light data: long cache time
func getStaticData() ([]StaticData, error) {
    session := pie.Table[StaticData](engine).WithCache(1 * time.Hour)
    // ...
}

// Read-light, write-heavy data: short cache time or no cache
func getFrequentlyUpdatedData() ([]DynamicData, error) {
    session := pie.Table[DynamicData](engine).WithCache(1 * time.Minute)
    // ...
}
```

### 2. Cache Key Naming Conventions

```go
// Use meaningful cache keys
const (
    UserCachePrefix      = "user:"
    UserListCachePrefix  = "users:list:"
    UserStatsCachePrefix = "users:stats:"
)

func getUserCacheKey(userID bson.ObjectID) string {
    return fmt.Sprintf("%s%s", UserCachePrefix, userID.Hex())
}

func getUserListCacheKey(status string, page int) string {
    return fmt.Sprintf("%s%s:page:%d", UserListCachePrefix, status, page)
}
```

### 3. Cache Update Strategy

```go
// Update cache when updating data
func updateUserWithCache(userID bson.ObjectID, updates bson.D) error {
    session := pie.Table[User](engine)
    
    // Update database
    _, err := session.
        Where("_id", userID).
        Update(ctx, updates)
    
    if err != nil {
        return err
    }
    
    // Update cache
    var user User
    err = session.
        Where("_id", userID).
        First(ctx, &user)
    
    if err == nil {
        cacheKey := getUserCacheKey(userID)
        engine.Cache().Set(cacheKey, user, 10*time.Minute)
    }
    
    // Clear related list cache
    engine.Cache().DeleteByPattern("users:list:*")
    
    return nil
}
```

## Next Steps

- [Hook System](/advanced/hooks/) - Learn about lifecycle hooks
- [Soft Delete](/advanced/soft-delete/) - Master soft delete functionality
- [Index Management](/advanced/indexes/) - Learn index optimization
