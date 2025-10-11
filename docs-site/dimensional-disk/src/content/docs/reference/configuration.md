---
title: Configuration
description: Comprehensive configuration options for different application scenarios
---

# Configuration

Pie provides comprehensive configuration options for different application scenarios.

## Engine Configuration

```go
// Create engine with configuration
engine, err := pie.NewEngine(
    context.Background(),
    "mydb",
    pie.WithURI("mongodb://localhost:27017"),
    pie.WithMaxPoolSize(100),
    pie.WithMinPoolSize(10),
    pie.WithConnectTimeout(5*time.Second),
    pie.WithReadPreference(readpref.PrimaryPreferred()),
    pie.WithWriteConcern(writeconcern.New(writeconcern.WMajority())),
    pie.WithLogger(pie.NewDefaultLogger()),
    pie.WithMapper(&pie.SnakeMapper{}),
    pie.WithAutoIndex(true),
    pie.WithSoftDeleteField("deleted_at"),
    pie.WithCache(pie.NewMemoryCache(), &pie.CacheConfig{TTL: 5 * time.Minute}),
)
```

## Runtime Configuration

```go
// Update configuration at runtime
engine.SetMaxPoolSize(50)
engine.SetLogger(pie.NewLogger(log.Default()))
```

## Next Steps

- [Name Mappers](/reference/mappers/) - Learn about name mapping
- [Error Handling](/reference/error-handling/) - Learn about error handling
- [Performance](/reference/performance/) - Learn performance optimization