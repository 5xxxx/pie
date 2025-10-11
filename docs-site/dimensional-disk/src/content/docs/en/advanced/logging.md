---
title: Logging & Monitoring
description: Built-in query logging and performance monitoring
---

# Logging & Monitoring

Pie provides built-in query logging and performance monitoring capabilities.

## Basic Usage

```go
// Enable default logging
engine.UseLogger(pie.NewDefaultLogger())

// Custom logger
customLogger := log.New(os.Stdout, "[PIE] ", log.LstdFlags)
engine.UseLogger(pie.NewLogger(customLogger))

// Query event monitoring
engine.OnQuery(func(event *pie.QueryEvent) {
    log.Printf("Query: %s, Duration: %v, Error: %v", 
        event.Query, event.Duration, event.Error)
})
```

## Next Steps

- [Advanced Aggregation](/advanced/aggregation-advanced/) - Learn advanced aggregation
- [Performance](/reference/performance/) - Learn performance optimization
- [Best Practices](/best-practices/) - Learn development best practices
