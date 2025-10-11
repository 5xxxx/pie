# Pie - MongoDB ORM Framework

Pie is a modern, type-safe MongoDB ORM framework for Go that provides a rich, high-performance database operation experience.

## Features

- **Type Safety**: Built on Go generics with compile-time type checking
- **Smart Query Builder**: Intuitive chainable API with support for complex query conditions
- **Struct Query**: Direct conversion of HTTP request parameters to query conditions
- **Cache Support**: Plugin-based caching architecture with Ristretto and Redis support
- **Hook System**: Complete lifecycle hook support
- **Transaction Management**: Simple and easy-to-use transaction operations
- **Index Management**: Automated index creation and management
- **Advanced Aggregation**: Comprehensive stage builders with 100+ expression functions
- **Change Streams**: Real-time data change monitoring
- **Soft Delete**: Built-in soft delete functionality
- **Pagination**: Efficient pagination implementation
- **Query Logging**: Detailed query logging and performance monitoring
- **Bulk Operations**: High-performance bulk write operations

## Installation

```bash
go get github.com/5xxxx/pie
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "github.com/5xxxx/pie"
)

func main() {
    // Create engine
    engine, err := pie.NewEngine(
        context.Background(),
        "mydb",
        pie.WithURI("mongodb://localhost:27017"),
    )
    if err != nil {
        log.Fatal("Failed to create engine:", err)
    }
    defer engine.Disconnect(context.Background())
    
    // Create type-safe session
    session := pie.Table[User](engine)
    
    // Insert document
    user := &User{Name: "John Doe", Email: "john@example.com"}
    result, err := session.Insert(context.Background(), user)
    
    // Query documents
    var users []User
    err = session.Where("age", pie.Gte("age", 18)).Find(context.Background(), &users)
}
```

## Cache Plugin Architecture

Pie features a flexible plugin-based caching system that supports multiple cache implementations and chaining.

### Default Ristretto Cache

```go
// Enable default Ristretto memory cache
engine.UseDefaultCache()

// Or with custom configuration
ristrettoConfig := &pie.RistrettoCacheConfig{
    NumCounters: 100000,
    MaxCost:     100 * 1024 * 1024, // 100MB
    BufferItems: 64,
}
engine.UseRistretto(ristrettoConfig)
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
```

### Custom Cache Implementation

```go
type MyCache struct {
    data map[string][]byte
}

func (m *MyCache) Get(ctx context.Context, key string) ([]byte, error) {
    // Implement your cache logic
}

func (m *MyCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // Implement your cache logic
}

// ... implement other Cache interface methods

// Use custom cache
myCache := &MyCache{}
engine.UseCache(myCache)
```

### Session Cache Usage

```go
// Basic caching
products, err := session.Cache(5*time.Minute).Find(ctx)

// Cache with tags
products, err := session.CacheWithTags(10*time.Minute, "electronics").Find(ctx)

// Cache with TTL jitter
products, err := session.CacheWithJitter(10*time.Minute, 2*time.Minute).Find(ctx)

// Cache empty results (anti-penetration)
products, err := session.CacheEmpty(30*time.Second).Find(ctx)
```

## Documentation

Complete documentation is available at: [https://5xxxx.github.io/pie](https://5xxxx.github.io/pie)

## License

MIT License

## Contributing

Issues and Pull Requests are welcome!