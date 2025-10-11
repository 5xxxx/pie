# Pie - MongoDB ORM 框架

Pie 是一个现代化的 Go 语言 MongoDB ORM 框架，提供类型安全、高性能、功能丰富的数据库操作体验。

## 特性

- **类型安全**: 基于 Go 泛型，提供编译时类型检查
- **智能查询构建器**: 直观的链式调用，支持复杂查询条件
- **结构体查询**: HTTP 请求参数直接转换为查询条件
- **缓存支持**: 基于插件的缓存架构，支持 Ristretto 和 Redis
- **钩子系统**: 完整的生命周期钩子支持
- **事务管理**: 简单易用的事务操作
- **索引管理**: 自动化的索引创建和管理
- **高级聚合**: 全面的阶段构建器，包含 100+ 表达式函数
- **变更流**: 实时数据变更监听
- **软删除**: 内置软删除功能
- **分页查询**: 高效的分页实现
- **查询日志**: 详细的查询日志和性能监控
- **批量操作**: 高效的批量写入操作

## 安装

```bash
go get github.com/5xxxx/pie
```

## 快速开始

```go
package main

import (
    "context"
    "log"
    "github.com/5xxxx/pie"
)

func main() {
    // 创建引擎
    engine, err := pie.NewEngine(
        context.Background(),
        "mydb",
        pie.WithURI("mongodb://localhost:27017"),
    )
    if err != nil {
        log.Fatal("Failed to create engine:", err)
    }
    defer engine.Disconnect(context.Background())
    
// 创建类型安全的会话
session := pie.Table[User](engine)

// 插入文档
    user := &User{Name: "张三", Email: "zhangsan@example.com"}
result, err := session.Insert(context.Background(), user)

// 查询文档
var users []User
    err = session.Where("age", pie.Gte("age", 18)).Find(context.Background(), &users)
}
```

## 缓存插件架构

Pie 采用灵活的插件化缓存系统，支持多种缓存实现和链式组合。

### 默认 Ristretto 缓存

```go
// 启用默认 Ristretto 内存缓存
engine.UseDefaultCache()

// 或使用自定义配置
ristrettoConfig := &pie.RistrettoCacheConfig{
    NumCounters: 100000,
    MaxCost:     100 * 1024 * 1024, // 100MB
    BufferItems: 64,
}
engine.UseRistretto(ristrettoConfig)
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
```

### 多层缓存链

```go
// 创建多个缓存实例
ristrettoCache, _ := pie.NewRistrettoCache(nil)
redisCache, _ := pie.NewRedisCache(&pie.RedisCacheConfig{
    Addr: "localhost:6379",
})

// 使用链式缓存 (L1: Ristretto, L2: Redis)
engine.UseCache(ristrettoCache, redisCache)
```

### 自定义缓存实现

```go
type MyCache struct {
    data map[string][]byte
}

func (m *MyCache) Get(ctx context.Context, key string) ([]byte, error) {
    // 实现你的缓存逻辑
}

func (m *MyCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // 实现你的缓存逻辑
}

// ... 实现其他 Cache 接口方法

// 使用自定义缓存
myCache := &MyCache{}
engine.UseCache(myCache)
```

### Session 缓存使用

```go
// 基础缓存
products, err := session.Cache(5*time.Minute).Find(ctx)

// 带标签缓存
products, err := session.CacheWithTags(10*time.Minute, "electronics").Find(ctx)

// 带 TTL 抖动缓存
products, err := session.CacheWithJitter(10*time.Minute, 2*time.Minute).Find(ctx)

// 缓存空结果（防穿透）
products, err := session.CacheEmpty(30*time.Second).Find(ctx)
```

## 文档

完整的文档请访问：[https://5xxxx.github.io/pie](https://5xxxx.github.io/pie)

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！