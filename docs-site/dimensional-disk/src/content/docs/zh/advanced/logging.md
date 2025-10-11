---
title: 日志记录与监控
description: 内置查询日志和性能监控
---

# 日志记录与监控

Pie 提供内置的查询日志和性能监控功能。

## 基础用法

```go
// 启用默认日志
engine.UseLogger(pie.NewDefaultLogger())

// 自定义日志
customLogger := log.New(os.Stdout, "[PIE] ", log.LstdFlags)
engine.UseLogger(pie.NewLogger(customLogger))

// 查询事件监控
engine.OnQuery(func(event *pie.QueryEvent) {
    log.Printf("Query: %s, Duration: %v, Error: %v", 
        event.Query, event.Duration, event.Error)
})
```

## 下一步

- [高级聚合](/advanced/aggregation-advanced/) - 学习高级聚合
- [性能](/reference/performance/) - 学习性能优化
- [最佳实践](/best-practices/) - 学习开发最佳实践
