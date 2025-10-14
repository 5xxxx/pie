---
title: 索引管理
description: 学习 Pie 的索引管理功能，优化查询性能
---

# 索引管理

Pie 提供了自动化的索引创建和管理功能，帮助您优化查询性能。

## 基本用法

### 创建索引管理器

```go
// 创建索引管理器
indexes := pie.MustIndexes(engine)

// 或者从引擎获取
indexes := engine.Indexes()
```

### 为结构体创建索引

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name" pie:"index"`
    Email     string             `bson:"email" pie:"unique"`
    Age       int                `bson:"age"`
    Status    string             `bson:"status" pie:"index"`
    CreatedAt time.Time          `bson:"created_at" pie:"index"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// 为结构体创建所有索引
err := indexes.CreateIndexes(ctx, User{})
```

### 手动创建索引

```go
// 创建单字段索引
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"email", 1},
}, &options.IndexOptions{
    Unique: pie.Bool(true),
})

// 创建复合索引
err = indexes.CreateIndex(ctx, "users", bson.D{
    {"status", 1},
    {"created_at", -1},
}, &options.IndexOptions{
    Name: pie.String("status_created_at_idx"),
})

// 创建文本索引
err = indexes.CreateIndex(ctx, "users", bson.D{
    {"name", "text"},
    {"email", "text"},
}, &options.IndexOptions{
    Name: pie.String("text_search_idx"),
})
```

## 索引标签

### 基础索引标签

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name" pie:"index"`           // 普通索引
    Email     string             `bson:"email" pie:"unique"`         // 唯一索引
    Age       int                `bson:"age" pie:"index,sparse"`     // 稀疏索引
    Status    string             `bson:"status" pie:"index,partial"` // 部分索引
    CreatedAt time.Time          `bson:"created_at" pie:"index"`     // 时间索引
}
```

### 复合索引标签

```go
type Order struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    UserID    bson.ObjectID `bson:"user_id" pie:"index:user_status"`
    Status    string             `bson:"status" pie:"index:user_status"`
    CreatedAt time.Time          `bson:"created_at" pie:"index:user_created"`
    UserID2   bson.ObjectID `bson:"user_id2" pie:"index:user_created"`
}
```

### 高级索引标签

```go
type Product struct {
    ID          bson.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name" pie:"index,text"`
    Description string             `bson:"description" pie:"index,text"`
    Price       float64            `bson:"price" pie:"index,sparse"`
    Category    string             `bson:"category" pie:"index,partial"`
    Tags        []string           `bson:"tags" pie:"index,multikey"`
    Location    GeoJSON            `bson:"location" pie:"index,2dsphere"`
    CreatedAt   time.Time          `bson:"created_at" pie:"index,expire,3600"` // TTL 索引
}
```

## 实际应用场景

### 用户管理索引

```go
func setupUserIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 用户基础索引
    err := indexes.CreateIndexes(ctx, User{})
    if err != nil {
        return err
    }
    
    // 复合索引：状态 + 创建时间
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    }, &options.IndexOptions{
        Name: pie.String("status_created_at_idx"),
    })
    if err != nil {
        return err
    }
    
    // 复合索引：角色 + 状态
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"role", 1},
        {"status", 1},
    }, &options.IndexOptions{
        Name: pie.String("role_status_idx"),
    })
    if err != nil {
        return err
    }
    
    // 文本搜索索引
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"name", "text"},
        {"email", "text"},
    }, &options.IndexOptions{
        Name: pie.String("user_text_search_idx"),
    })
    
    return err
}
```

### 订单管理索引

```go
func setupOrderIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 订单基础索引
    err := indexes.CreateIndexes(ctx, Order{})
    if err != nil {
        return err
    }
    
    // 用户订单查询索引
    err = indexes.CreateIndex(ctx, "orders", bson.D{
        {"user_id", 1},
        {"status", 1},
        {"created_at", -1},
    }, &options.IndexOptions{
        Name: pie.String("user_orders_idx"),
    })
    if err != nil {
        return err
    }
    
    // 订单状态统计索引
    err = indexes.CreateIndex(ctx, "orders", bson.D{
        {"status", 1},
        {"created_at", 1},
    }, &options.IndexOptions{
        Name: pie.String("order_status_time_idx"),
    })
    if err != nil {
        return err
    }
    
    // 金额范围查询索引
    err = indexes.CreateIndex(ctx, "orders", bson.D{
        {"total", 1},
    }, &options.IndexOptions{
        Name: pie.String("order_total_idx"),
    })
    
    return err
}
```

### 地理位置索引

```go
func setupLocationIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 2dsphere 索引用于地理位置查询
    err := indexes.CreateIndex(ctx, "locations", bson.D{
        {"location", "2dsphere"},
    }, &options.IndexOptions{
        Name: pie.String("location_2dsphere_idx"),
    })
    if err != nil {
        return err
    }
    
    // 复合索引：位置 + 类型
    err = indexes.CreateIndex(ctx, "locations", bson.D{
        {"location", "2dsphere"},
        {"type", 1},
    }, &options.IndexOptions{
        Name: pie.String("location_type_idx"),
    })
    
    return err
}
```

## 高级用法

### 条件索引

```go
// 部分索引：只为特定条件的文档创建索引
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"email", 1},
}, &options.IndexOptions{
    Name: pie.String("active_user_email_idx"),
    PartialFilterExpression: bson.D{
        {"status", "active"},
    },
})

// 稀疏索引：只为存在该字段的文档创建索引
err = indexes.CreateIndex(ctx, "users", bson.D{
    {"phone", 1},
}, &options.IndexOptions{
    Name: pie.String("user_phone_idx"),
    Sparse: pie.Bool(true),
})
```

### TTL 索引

```go
// 创建 TTL 索引，自动删除过期文档
err := indexes.CreateIndex(ctx, "sessions", bson.D{
    {"expires_at", 1},
}, &options.IndexOptions{
    Name: pie.String("session_ttl_idx"),
    ExpireAfterSeconds: pie.Int32(0), // 立即过期
})

// 或者使用标签
type Session struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    UserID    bson.ObjectID `bson:"user_id"`
    ExpiresAt time.Time          `bson:"expires_at" pie:"index,expire,3600"` // 1小时后过期
}
```

### 文本索引

```go
// 创建文本搜索索引
err := indexes.CreateIndex(ctx, "articles", bson.D{
    {"title", "text"},
    {"content", "text"},
    {"tags", "text"},
}, &options.IndexOptions{
    Name: pie.String("article_text_idx"),
    DefaultLanguage: pie.String("english"),
    Weights: bson.D{
        {"title", 10},
        {"content", 5},
        {"tags", 3},
    },
})
```

### 复合索引优化

```go
// 为常见查询模式创建复合索引
func createOptimizedIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 查询模式：WHERE status = ? AND created_at > ? ORDER BY created_at DESC
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    }, &options.IndexOptions{
        Name: pie.String("status_created_desc_idx"),
    })
    if err != nil {
        return err
    }
    
    // 查询模式：WHERE user_id = ? AND status IN (?) ORDER BY created_at DESC
    err = indexes.CreateIndex(ctx, "orders", bson.D{
        {"user_id", 1},
        {"status", 1},
        {"created_at", -1},
    }, &options.IndexOptions{
        Name: pie.String("user_status_created_idx"),
    })
    if err != nil {
        return err
    }
    
    // 查询模式：WHERE category = ? AND price BETWEEN ? AND ? ORDER BY created_at DESC
    err = indexes.CreateIndex(ctx, "products", bson.D{
        {"category", 1},
        {"price", 1},
        {"created_at", -1},
    }, &options.IndexOptions{
        Name: pie.String("category_price_created_idx"),
    })
    
    return err
}
```

## 索引管理

### 列出所有索引

```go
func listIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 列出集合的所有索引
    indexList, err := indexes.ListIndexes(ctx, "users")
    if err != nil {
        return err
    }
    
    for _, index := range indexList {
        fmt.Printf("索引名称: %s\n", index.Name)
        fmt.Printf("索引键: %v\n", index.Key)
        fmt.Printf("是否唯一: %t\n", index.Unique)
        fmt.Printf("是否稀疏: %t\n", index.Sparse)
        fmt.Println("---")
    }
    
    return nil
}
```

### 删除索引

```go
func deleteIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 删除指定索引
    err := indexes.DropIndex(ctx, "users", "email_1")
    if err != nil {
        return err
    }
    
    // 删除所有索引（除了 _id 索引）
    err = indexes.DropAllIndexes(ctx, "users")
    if err != nil {
        return err
    }
    
    return nil
}
```

### 重建索引

```go
func rebuildIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 重建集合的所有索引
    err := indexes.RebuildIndexes(ctx, "users")
    if err != nil {
        return err
    }
    
    return nil
}
```

## 性能优化

### 1. 索引选择策略

```go
// 好的做法：为查询模式创建合适的索引
func createQueryOptimizedIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 单字段查询
    err := indexes.CreateIndex(ctx, "users", bson.D{{"email", 1}})
    if err != nil {
        return err
    }
    
    // 范围查询
    err = indexes.CreateIndex(ctx, "users", bson.D{{"age", 1}})
    if err != nil {
        return err
    }
    
    // 复合查询：WHERE status = ? AND created_at > ? ORDER BY created_at DESC
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    })
    
    return err
}
```

### 2. 索引监控

```go
func monitorIndexUsage() error {
    indexes := pie.MustIndexes(engine)
    
    // 获取索引使用统计
    stats, err := indexes.GetIndexStats(ctx, "users")
    if err != nil {
        return err
    }
    
    for _, stat := range stats {
        fmt.Printf("索引: %s\n", stat.Name)
        fmt.Printf("访问次数: %d\n", stat.Accesses)
        fmt.Printf("最后访问: %v\n", stat.LastAccess)
        fmt.Println("---")
    }
    
    return nil
}
```

### 3. 索引维护

```go
func maintainIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 定期重建索引
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if err := indexes.RebuildIndexes(ctx, "users"); err != nil {
                    log.Printf("Failed to rebuild indexes: %v", err)
                }
            }
        }
    }()
    
    return nil
}
```

## 最佳实践

### 1. 索引命名规范

```go
// 使用有意义的索引名称
func createNamedIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"email", 1},
    }, &options.IndexOptions{
        Name: pie.String("idx_users_email_unique"),
    })
    
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    }, &options.IndexOptions{
        Name: pie.String("idx_users_status_created_desc"),
    })
    
    return err
}
```

### 2. 索引数量控制

```go
// 避免创建过多索引
func createEssentialIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 只创建必要的索引
    essentialIndexes := []struct {
        name string
        keys bson.D
    }{
        {"idx_users_email", bson.D{{"email", 1}}},
        {"idx_users_status_created", bson.D{{"status", 1}, {"created_at", -1}}},
        {"idx_users_role_status", bson.D{{"role", 1}, {"status", 1}}},
    }
    
    for _, idx := range essentialIndexes {
        err := indexes.CreateIndex(ctx, "users", idx.keys, &options.IndexOptions{
            Name: pie.String(idx.name),
        })
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 3. 索引测试

```go
func testIndexPerformance() error {
    // 测试查询性能
    start := time.Now()
    
    users, err := session.
        Where("status", "active").
        Where("created_at", pie.Gte("created_at", time.Now().AddDate(0, -1, 0))).
        OrderByDesc("created_at").
        Find(ctx)
    
    duration := time.Since(start)
    
    if err != nil {
        return err
    }
    
    log.Printf("查询耗时: %v, 结果数量: %d", duration, len(users))
    
    // 如果查询时间过长，考虑创建索引
    if duration > 100*time.Millisecond {
        log.Println("查询较慢，建议检查索引")
    }
    
    return nil
}
```

## 下一步

- [变更流](/advanced/change-streams/) - 学习实时监听
- [查询作用域](/advanced/scopes/) - 掌握作用域定义
- [日志监控](/advanced/logging/) - 学习查询日志
