---
title: 配置选项
description: 了解 Pie 的所有配置选项和参数
---

# 配置选项

Pie 提供了丰富的配置选项来满足不同部署环境和性能需求。

## 基本配置

### 创建引擎

```go
// 基础配置
engine, err := pie.NewEngine(
    context.Background(),
    "mydb",
    pie.WithURI("mongodb://localhost:27017"),
)
```

### 连接配置

```go
// 完整连接配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithURI("mongodb://localhost:27017"),
    pie.WithAuth("username", "password"),
    pie.WithSSL(true),
    pie.WithReplicaSet("rs0"),
    pie.WithReadPreference("secondary"),
    pie.WithWriteConcern("majority"),
    pie.WithReadConcern("majority"),
)
```

## 连接池配置

### 连接池大小

```go
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMaxPoolSize(100),        // 最大连接数
    pie.WithMinPoolSize(5),          // 最小连接数
    pie.WithMaxIdleTime(30*time.Minute), // 最大空闲时间
)
```

### 连接超时

```go
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithConnectTimeout(10*time.Second),      // 连接超时
    pie.WithSocketTimeout(30*time.Second),       // 套接字超时
    pie.WithServerSelectionTimeout(5*time.Second), // 服务器选择超时
)
```

## 读写配置

### 读偏好

```go
// 读偏好配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithReadPreference("primary"),           // 主节点
    pie.WithReadPreference("secondary"),         // 从节点
    pie.WithReadPreference("primaryPreferred"),  // 主节点优先
    pie.WithReadPreference("secondaryPreferred"), // 从节点优先
    pie.WithReadPreference("nearest"),           // 最近节点
)
```

### 写关注

```go
// 写关注配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithWriteConcern("majority"),            // 大多数节点确认
    pie.WithWriteConcern("1"),                   // 单个节点确认
    pie.WithWriteConcern("0"),                   // 无确认
)
```

### 读关注

```go
// 读关注配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithReadConcern("local"),                // 本地读关注
    pie.WithReadConcern("available"),            // 可用读关注
    pie.WithReadConcern("majority"),             // 大多数读关注
    pie.WithReadConcern("linearizable"),         // 线性化读关注
)
```

## 安全配置

### 认证配置

```go
// 用户名密码认证
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithAuth("username", "password"),
)

// 认证数据库
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithAuth("username", "password"),
    pie.WithAuthSource("admin"),
)
```

### SSL/TLS 配置

```go
// SSL 配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithSSL(true),
    pie.WithSSLInsecure(false),
    pie.WithSSLCertFile("/path/to/cert.pem"),
    pie.WithSSLKeyFile("/path/to/key.pem"),
    pie.WithSSLCAFile("/path/to/ca.pem"),
)
```

## 性能配置

### 查询配置

```go
// 查询配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithQueryLog(os.Stdout),                 // 查询日志
    pie.WithQueryTimeout(30*time.Second),        // 查询超时
)
```

### 缓存配置

```go
// 缓存配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithCache(pie.NewMemoryCache(), &pie.CacheConfig{
        TTL: 5 * time.Minute,
        MaxSize: 1000,
    }),
)
```

## 高级配置

### 命名映射

```go
// 命名映射配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SnakeMapper{}),          // 蛇形命名
    pie.WithMapper(&pie.CamelMapper{}),          // 驼峰命名
    pie.WithMapper(&pie.SameMapper{}),           // 相同命名
)

// 自定义映射器
type CustomMapper struct{}

func (m CustomMapper) TableName(structName string) string {
    return "t_" + strings.ToLower(structName)
}

func (m CustomMapper) FieldName(fieldName string) string {
    return strings.ToLower(fieldName)
}

engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(CustomMapper{}),
)
```

### 日志配置

```go
// 日志配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithLogger(log.New(os.Stdout, "PIE: ", log.LstdFlags)),
    pie.WithLogLevel(pie.LogLevelInfo),
    pie.WithQueryLog(os.Stdout),
)
```

### 监控配置

```go
// 监控配置
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMetrics(true),                       // 启用指标
    pie.WithHealthCheck(true),                   // 启用健康检查
    pie.WithPingInterval(30*time.Second),        // Ping 间隔
)
```

## 配置选项参考

### 连接选项

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `WithURI` | string | - | MongoDB 连接 URI |
| `WithAuth` | string, string | - | 用户名和密码 |
| `WithAuthSource` | string | "admin" | 认证数据库 |
| `WithSSL` | bool | false | 启用 SSL |
| `WithSSLInsecure` | bool | false | 跳过 SSL 验证 |
| `WithSSLCertFile` | string | - | SSL 证书文件 |
| `WithSSLKeyFile` | string | - | SSL 私钥文件 |
| `WithSSLCAFile` | string | - | SSL CA 文件 |

### 连接池选项

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `WithMaxPoolSize` | uint64 | 100 | 最大连接数 |
| `WithMinPoolSize` | uint64 | 0 | 最小连接数 |
| `WithMaxIdleTime` | time.Duration | 0 | 最大空闲时间 |
| `WithConnectTimeout` | time.Duration | 30s | 连接超时 |
| `WithSocketTimeout` | time.Duration | 0 | 套接字超时 |
| `WithServerSelectionTimeout` | time.Duration | 30s | 服务器选择超时 |

### 读写选项

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `WithReadPreference` | string | "primary" | 读偏好 |
| `WithWriteConcern` | string | "majority" | 写关注 |
| `WithReadConcern` | string | "local" | 读关注 |

### 性能选项

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `WithQueryLog` | io.Writer | - | 查询日志输出 |
| `WithQueryTimeout` | time.Duration | - | 查询超时 |
| `WithCache` | Cache, *CacheConfig | - | 缓存配置 |

## 环境配置

### 开发环境

```go
func createDevEngine() (*pie.Engine, error) {
    return pie.NewEngine(ctx, "dev_db",
        pie.WithURI("mongodb://localhost:27017"),
        pie.WithQueryLog(os.Stdout),
        pie.WithMapper(&pie.SnakeMapper{}),
    )
}
```

### 测试环境

```go
func createTestEngine() (*pie.Engine, error) {
    return pie.NewEngine(ctx, "test_db",
        pie.WithURI("mongodb://localhost:27017"),
        pie.WithMaxPoolSize(10),
        pie.WithMinPoolSize(1),
        pie.WithMapper(&pie.SameMapper{}),
    )
}
```

### 生产环境

```go
func createProdEngine() (*pie.Engine, error) {
    return pie.NewEngine(ctx, "prod_db",
        pie.WithURI("mongodb://prod-cluster:27017"),
        pie.WithAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
        pie.WithSSL(true),
        pie.WithReplicaSet("rs0"),
        pie.WithReadPreference("secondaryPreferred"),
        pie.WithWriteConcern("majority"),
        pie.WithMaxPoolSize(100),
        pie.WithMinPoolSize(10),
        pie.WithCache(pie.NewRedisCache("redis:6379", "", 0), &pie.CacheConfig{
            TTL: 10 * time.Minute,
        }),
    )
}
```

## 配置验证

### 配置检查

```go
func validateConfig(engine *pie.Engine) error {
    // 检查连接
    if err := engine.Ping(ctx); err != nil {
        return fmt.Errorf("database connection failed: %w", err)
    }
    
    // 检查权限
    if err := engine.CheckPermissions(ctx); err != nil {
        return fmt.Errorf("insufficient permissions: %w", err)
    }
    
    return nil
}
```

### 配置测试

```go
func testConfiguration() error {
    engine, err := createProdEngine()
    if err != nil {
        return err
    }
    defer engine.Disconnect(ctx)
    
    // 测试基本操作
    session := pie.Table[User](engine)
    _, err = session.Count(ctx)
    if err != nil {
        return fmt.Errorf("basic query failed: %w", err)
    }
    
    // 测试事务
    err = engine.WithTransaction(ctx, func(txCtx context.Context) error {
        return nil
    })
    if err != nil {
        return fmt.Errorf("transaction test failed: %w", err)
    }
    
    return nil
}
```

## 最佳实践

### 1. 环境特定配置

```go
func createEngineForEnvironment() (*pie.Engine, error) {
    env := os.Getenv("ENVIRONMENT")
    
    switch env {
    case "development":
        return createDevEngine()
    case "testing":
        return createTestEngine()
    case "production":
        return createProdEngine()
    default:
        return nil, fmt.Errorf("unknown environment: %s", env)
    }
}
```

### 2. 配置管理

```go
type Config struct {
    Database struct {
        URI      string `yaml:"uri"`
        Database string `yaml:"database"`
        Auth     struct {
            Username string `yaml:"username"`
            Password string `yaml:"password"`
        } `yaml:"auth"`
    } `yaml:"database"`
    
    Cache struct {
        Type string `yaml:"type"`
        TTL  int    `yaml:"ttl"`
    } `yaml:"cache"`
}

func loadConfig() (*Config, error) {
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = yaml.Unmarshal(data, &config)
    return &config, err
}
```

### 3. 配置验证

```go
func validateEngineConfig(config *Config) error {
    if config.Database.URI == "" {
        return errors.New("database URI is required")
    }
    
    if config.Database.Database == "" {
        return errors.New("database name is required")
    }
    
    return nil
}
```

## 下一步

- [命名映射](/reference/mappers/) - 学习自定义映射器
- [错误处理](/reference/error-handling/) - 掌握错误处理
- [性能优化](/reference/performance/) - 了解性能优化
