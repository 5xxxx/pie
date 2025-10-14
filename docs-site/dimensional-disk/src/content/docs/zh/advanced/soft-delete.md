---
title: 软删除
description: 学习 Pie 的软删除功能，保护重要数据
---

# 软删除

Pie 提供了内置的软删除功能，允许您"删除"数据而不实际从数据库中移除，保护重要数据并提供数据恢复能力。

## 基本用法

### 定义软删除字段

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

type Product struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Price     float64            `bson:"price"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" pie:"soft_delete"`
}
```

### 软删除操作

```go
// 软删除用户
err := session.Where("email", "test@example.com").SoftDelete(ctx)

// 软删除多个用户
err := session.Where("status", "inactive").SoftDeleteMany(ctx)

// 软删除特定用户
err := session.Where("_id", userID).SoftDelete(ctx)
```

### 恢复软删除

```go
// 恢复软删除的用户
err := session.Where("email", "test@example.com").Restore(ctx)

// 恢复多个用户
err := session.Where("status", "inactive").RestoreMany(ctx)

// 恢复特定用户
err := session.Where("_id", userID).Restore(ctx)
```

### 强制删除

```go
// 强制删除（物理删除）
err := session.Where("email", "test@example.com").ForceDelete(ctx)

// 强制删除多个
err := session.Where("status", "inactive").ForceDeleteMany(ctx)
```

## 查询行为

### 自动过滤软删除记录

```go
// 普通查询自动排除软删除的记录
users, err := session.Find(ctx) // 自动过滤 deleted_at 不为 null 的记录

// 条件查询也自动过滤
activeUsers, err := session.Where("status", "active").Find(ctx)

// 聚合查询也自动过滤
result, err := session.
    GroupStage().
        By("role", "$role").
        Count("total").
        Done().
    Exec(ctx)
```

### 包含软删除记录

```go
// 使用 WithTrashed 包含软删除的记录
allUsers, err := session.WithTrashed().Find(ctx)

// 只查询软删除的记录
deletedUsers, err := session.OnlyTrashed().Find(ctx)

// 在条件查询中包含软删除
users, err := session.
    WithTrashed().
    Where("email", "test@example.com").
    Find(ctx)
```

## 实际应用场景

### 用户管理

```go
func deleteUser(userID bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // 软删除用户
    err := session.Where("_id", userID).SoftDelete(ctx)
    if err != nil {
        return fmt.Errorf("failed to soft delete user: %w", err)
    }
    
    // 记录删除日志
    logUserDeletion(userID)
    
    return nil
}

func restoreUser(userID bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // 恢复用户
    err := session.Where("_id", userID).Restore(ctx)
    if err != nil {
        return fmt.Errorf("failed to restore user: %w", err)
    }
    
    // 记录恢复日志
    logUserRestoration(userID)
    
    return nil
}

func permanentlyDeleteUser(userID bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // 强制删除用户
    err := session.Where("_id", userID).ForceDelete(ctx)
    if err != nil {
        return fmt.Errorf("failed to permanently delete user: %w", err)
    }
    
    // 清理相关数据
    cleanupUserRelatedData(userID)
    
    return nil
}
```

### 订单管理

```go
func cancelOrder(orderID bson.ObjectID) error {
    session := pie.Table[Order](engine)
    
    // 软删除订单
    err := session.Where("_id", orderID).SoftDelete(ctx)
    if err != nil {
        return fmt.Errorf("failed to cancel order: %w", err)
    }
    
    // 恢复库存
    if err := restoreOrderStock(orderID); err != nil {
        log.Printf("Failed to restore stock for order %s: %v", orderID.Hex(), err)
    }
    
    // 发送取消通知
    go sendOrderCancellationNotification(orderID)
    
    return nil
}

func restoreOrder(orderID bson.ObjectID) error {
    session := pie.Table[Order](engine)
    
    // 恢复订单
    err := session.Where("_id", orderID).Restore(ctx)
    if err != nil {
        return fmt.Errorf("failed to restore order: %w", err)
    }
    
    // 重新扣减库存
    if err := deductOrderStock(orderID); err != nil {
        log.Printf("Failed to deduct stock for order %s: %v", orderID.Hex(), err)
    }
    
    return nil
}
```

### 数据清理

```go
func cleanupOldSoftDeletedData() error {
    session := pie.Table[User](engine)
    
    // 查找30天前软删除的用户
    cutoffDate := time.Now().AddDate(0, 0, -30)
    
    deletedUsers, err := session.
        OnlyTrashed().
        Where("deleted_at", pie.Lt("deleted_at", cutoffDate)).
        Find(ctx)
    
    if err != nil {
        return err
    }
    
    // 永久删除这些用户
    for _, user := range deletedUsers {
        if err := session.Where("_id", user.ID).ForceDelete(ctx); err != nil {
            log.Printf("Failed to permanently delete user %s: %v", user.ID.Hex(), err)
        }
    }
    
    log.Printf("Permanently deleted %d old soft-deleted users", len(deletedUsers))
    return nil
}
```

## 高级用法

### 自定义软删除字段

```go
type CustomUser struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" pie:"soft_delete"`
    DeletedBy *bson.ObjectID `bson:"deleted_by,omitempty"`
    DeleteReason string          `bson:"delete_reason,omitempty"`
}

// 自定义软删除逻辑
func (u *CustomUser) BeforeSoftDelete(ctx context.Context) error {
    // 设置删除者
    if userID := getCurrentUserID(ctx); userID != nil {
        u.DeletedBy = userID
    }
    
    // 设置删除原因
    if reason := getDeleteReason(ctx); reason != "" {
        u.DeleteReason = reason
    }
    
    return nil
}
```

### 软删除钩子

```go
func (u *User) BeforeSoftDelete(ctx context.Context) error {
    // 检查是否可以删除
    if u.Email == "admin@example.com" {
        return errors.New("cannot delete admin user")
    }
    
    // 记录删除前的状态
    logUserSoftDeletion(u.ID, u.Name)
    
    return nil
}

func (u *User) AfterSoftDelete(ctx context.Context) error {
    // 清理相关数据
    go func() {
        cleanupUserSessions(u.ID)
        cleanupUserCache(u.ID)
    }()
    
    return nil
}

func (u *User) BeforeRestore(ctx context.Context) error {
    // 检查是否可以恢复
    if u.Email == "banned@example.com" {
        return errors.New("cannot restore banned user")
    }
    
    return nil
}

func (u *User) AfterRestore(ctx context.Context) error {
    // 重新初始化用户数据
    go func() {
        reinitializeUserData(u.ID)
    }()
    
    return nil
}
```

### 批量软删除

```go
func batchSoftDeleteUsers(userIDs []bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // 批量软删除
    err := session.WhereIn("_id", userIDs).SoftDeleteMany(ctx)
    if err != nil {
        return fmt.Errorf("failed to batch soft delete users: %w", err)
    }
    
    // 记录批量删除日志
    logBatchUserDeletion(userIDs)
    
    return nil
}

func batchRestoreUsers(userIDs []bson.ObjectID) error {
    session := pie.Table[User](engine)
    
    // 批量恢复
    err := session.WhereIn("_id", userIDs).RestoreMany(ctx)
    if err != nil {
        return fmt.Errorf("failed to batch restore users: %w", err)
    }
    
    // 记录批量恢复日志
    logBatchUserRestoration(userIDs)
    
    return nil
}
```

### 软删除统计

```go
func getSoftDeleteStatistics() (*SoftDeleteStats, error) {
    session := pie.Table[User](engine)
    
    stats := &SoftDeleteStats{}
    
    // 总用户数（包括软删除）
    totalCount, err := session.WithTrashed().Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalCount = totalCount
    
    // 活跃用户数
    activeCount, err := session.Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.ActiveCount = activeCount
    
    // 软删除用户数
    deletedCount, err := session.OnlyTrashed().Count(ctx)
    if err != nil {
        return nil, err
    }
    stats.DeletedCount = deletedCount
    
    // 按删除时间分组统计
    var deleteStats []bson.M
    err = session.
        OnlyTrashed().
        GroupStage().
            By("deleteDate", pie.DateToString("$deleted_at", "%Y-%m-%d", "UTC")).
            Count("count").
            Done().
        Exec(ctx, &deleteStats)
    if err != nil {
        return nil, err
    }
    
    stats.DeleteStatsByDate = deleteStats
    
    return stats, nil
}
```

## 性能优化

### 1. 索引优化

```go
// 为软删除字段创建索引
func createSoftDeleteIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // 为 deleted_at 字段创建索引
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"deleted_at", 1},
    })
    if err != nil {
        return err
    }
    
    // 为复合查询创建索引
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"deleted_at", 1},
    })
    if err != nil {
        return err
    }
    
    return nil
}
```

### 2. 查询优化

```go
// 使用投影减少数据传输
func getActiveUsers() ([]User, error) {
    session := pie.Table[User](engine)

    users, err := session.
        Select("name", "email", "status"). // 只选择需要的字段
        Find(ctx)

    return users, err
}

// 使用分页避免大量数据
func getActiveUsersPaginated(page, pageSize int) (*PaginatedResult[User], error) {
    session := pie.Table[User](engine)
    
    result, err := session.
        Paginate(ctx, pie.PaginateParams{
            Page:     page,
            PageSize: pageSize,
        })
    
    return result, err
}
```

### 3. 缓存策略

```go
func getCachedActiveUsers() ([]User, error) {
    session := pie.Table[User](engine).WithCache(5 * time.Minute)

    users, err := session.
        Cache("active_users").
        Find(ctx)

    return users, err
}

func invalidateUserCache() error {
    cache := engine.Cache()
    
    // 清除相关缓存
    cache.Delete("active_users")
    cache.DeleteByPattern("users_*")
    
    return nil
}
```

## 最佳实践

### 1. 合理使用软删除

```go
// 好的做法：对重要数据使用软删除
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" pie:"soft_delete"`
}

// 避免的做法：对临时数据使用软删除
type TempData struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Data      string             `bson:"data"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" pie:"soft_delete"` // 不需要
}
```

### 2. 定期清理软删除数据

```go
func scheduleSoftDeleteCleanup() {
    // 每天凌晨2点清理30天前的软删除数据
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if time.Now().Hour() == 2 { // 凌晨2点
                    if err := cleanupOldSoftDeletedData(); err != nil {
                        log.Printf("Failed to cleanup soft deleted data: %v", err)
                    }
                }
            }
        }
    }()
}
```

### 3. 软删除权限控制

```go
func canSoftDeleteUser(userID bson.ObjectID, currentUserID bson.ObjectID) bool {
    // 检查权限
    if !hasPermission(currentUserID, "delete_user") {
        return false
    }
    
    // 不能删除自己
    if userID == currentUserID {
        return false
    }
    
    // 不能删除管理员
    user, err := getUserByID(userID)
    if err != nil {
        return false
    }
    
    if user.Role == "admin" {
        return false
    }
    
    return true
}
```

### 4. 软删除审计

```go
type SoftDeleteAudit struct {
    ID          bson.ObjectID `bson:"_id,omitempty"`
    EntityID    bson.ObjectID `bson:"entity_id"`
    EntityType  string             `bson:"entity_type"`
    Action      string             `bson:"action"` // "soft_delete" or "restore"
    DeletedBy   bson.ObjectID `bson:"deleted_by"`
    DeletedAt   time.Time          `bson:"deleted_at"`
    Reason      string             `bson:"reason,omitempty"`
}

func logSoftDeleteAudit(entityID bson.ObjectID, entityType, action string, reason string) {
    audit := &SoftDeleteAudit{
        EntityID:   entityID,
        EntityType: entityType,
        Action:     action,
        DeletedBy:  getCurrentUserID(ctx),
        DeletedAt:  time.Now(),
        Reason:     reason,
    }
    
    session := pie.Table[SoftDeleteAudit](engine)
    session.Insert(ctx, audit)
}
```

## 下一步

- [索引管理](/advanced/indexes/) - 学习索引优化
- [变更流](/advanced/change-streams/) - 掌握实时监听
- [查询作用域](/advanced/scopes/) - 学习作用域定义
