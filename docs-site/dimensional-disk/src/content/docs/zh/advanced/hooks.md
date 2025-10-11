---
title: 钩子系统
description: 学习 Pie 的完整生命周期钩子支持
---

# 钩子系统

Pie 提供了完整的生命周期钩子支持，允许您在数据操作的各个阶段执行自定义逻辑。

## 基本用法

### 定义钩子方法

```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Password  string             `bson:"password"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// 创建前钩子
func (u *User) BeforeCreate(ctx context.Context) error {
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    
    // 加密密码
    hashedPassword, err := hashPassword(u.Password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }
    u.Password = hashedPassword
    
    return nil
}

// 创建后钩子
func (u *User) AfterCreate(ctx context.Context) error {
    log.Printf("User %s created with ID %s", u.Name, u.ID.Hex())
    
    // 发送欢迎邮件
    go sendWelcomeEmail(u.Email, u.Name)
    
    return nil
}

// 更新前钩子
func (u *User) BeforeUpdate(ctx context.Context) error {
    u.UpdatedAt = time.Now()
    
    // 如果密码被更新，重新加密
    if u.Password != "" {
        hashedPassword, err := hashPassword(u.Password)
        if err != nil {
            return fmt.Errorf("failed to hash password: %w", err)
        }
        u.Password = hashedPassword
    }
    
    return nil
}

// 更新后钩子
func (u *User) AfterUpdate(ctx context.Context) error {
    log.Printf("User %s updated", u.Name)
    
    // 清除相关缓存
    clearUserCache(u.ID)
    
    return nil
}

// 删除前钩子
func (u *User) BeforeDelete(ctx context.Context) error {
    log.Printf("About to delete user %s", u.Name)
    
    // 检查是否可以删除
    if u.Email == "admin@example.com" {
        return errors.New("cannot delete admin user")
    }
    
    return nil
}

// 删除后钩子
func (u *User) AfterDelete(ctx context.Context) error {
    log.Printf("User %s deleted", u.Name)
    
    // 清理相关数据
    go cleanupUserData(u.ID)
    
    return nil
}

// 查找后钩子
func (u *User) AfterFind(ctx context.Context) error {
    // 隐藏敏感信息
    u.Password = ""
    
    // 添加计算字段
    u.Age = calculateAge(u.CreatedAt)
    
    return nil
}
```

## 钩子类型

### 创建钩子

```go
// BeforeCreate - 创建前
func (u *User) BeforeCreate(ctx context.Context) error {
    // 设置默认值
    if u.Status == "" {
        u.Status = "active"
    }
    
    // 生成唯一标识
    u.UUID = generateUUID()
    
    // 验证数据
    if err := u.validate(); err != nil {
        return err
    }
    
    return nil
}

// AfterCreate - 创建后
func (u *User) AfterCreate(ctx context.Context) error {
    // 记录审计日志
    logAuditEvent("user_created", u.ID, u.Name)
    
    // 发送通知
    notifyUserCreated(u)
    
    return nil
}
```

### 更新钩子

```go
// BeforeUpdate - 更新前
func (u *User) BeforeUpdate(ctx context.Context) error {
    // 记录原始数据
    originalUser := *u
    
    // 更新修改时间
    u.UpdatedAt = time.Now()
    
    // 检查权限
    if !canUpdateUser(originalUser, u) {
        return errors.New("insufficient permissions")
    }
    
    return nil
}

// AfterUpdate - 更新后
func (u *User) AfterUpdate(ctx context.Context) error {
    // 记录变更
    logUserChange(u.ID, "updated")
    
    // 更新搜索索引
    updateSearchIndex(u)
    
    return nil
}
```

### 删除钩子

```go
// BeforeDelete - 删除前
func (u *User) BeforeDelete(ctx context.Context) error {
    // 检查依赖关系
    if hasActiveOrders(u.ID) {
        return errors.New("cannot delete user with active orders")
    }
    
    // 备份数据
    backupUserData(u)
    
    return nil
}

// AfterDelete - 删除后
func (u *User) AfterDelete(ctx context.Context) error {
    // 清理相关数据
    cleanupUserRelatedData(u.ID)
    
    // 记录删除日志
    logUserDeletion(u.ID, u.Name)
    
    return nil
}
```

### 查找钩子

```go
// AfterFind - 查找后
func (u *User) AfterFind(ctx context.Context) error {
    // 隐藏敏感字段
    u.Password = ""
    u.SecretKey = ""
    
    // 添加计算字段
    u.FullName = u.FirstName + " " + u.LastName
    u.Age = calculateAge(u.BirthDate)
    
    // 格式化数据
    u.CreatedAtFormatted = u.CreatedAt.Format("2006-01-02 15:04:05")
    
    return nil
}
```

## 高级用法

### 条件钩子

```go
// 根据条件执行不同的钩子逻辑
func (u *User) BeforeCreate(ctx context.Context) error {
    // 根据用户类型执行不同逻辑
    switch u.Type {
    case "admin":
        return u.beforeCreateAdmin(ctx)
    case "user":
        return u.beforeCreateUser(ctx)
    case "guest":
        return u.beforeCreateGuest(ctx)
    default:
        return errors.New("invalid user type")
    }
}

func (u *User) beforeCreateAdmin(ctx context.Context) error {
    // 管理员创建逻辑
    u.Permissions = []string{"read", "write", "delete", "admin"}
    u.Verified = true
    return nil
}

func (u *User) beforeCreateUser(ctx context.Context) error {
    // 普通用户创建逻辑
    u.Permissions = []string{"read", "write"}
    u.Verified = false
    return nil
}
```

### 钩子链

```go
// 使用钩子链处理复杂逻辑
func (u *User) BeforeCreate(ctx context.Context) error {
    // 钩子链：验证 -> 处理 -> 保存
    if err := u.validateUser(); err != nil {
        return err
    }
    
    if err := u.processUserData(); err != nil {
        return err
    }
    
    if err := u.saveUserMetadata(); err != nil {
        return err
    }
    
    return nil
}

func (u *User) validateUser() error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    if !isValidEmail(u.Email) {
        return errors.New("invalid email format")
    }
    return nil
}

func (u *User) processUserData() error {
    // 处理用户数据
    u.Email = strings.ToLower(u.Email)
    u.Name = strings.TrimSpace(u.Name)
    return nil
}

func (u *User) saveUserMetadata() error {
    // 保存用户元数据
    u.Metadata = map[string]interface{}{
        "created_by": getCurrentUserID(ctx),
        "created_ip": getClientIP(ctx),
        "user_agent": getUserAgent(ctx),
    }
    return nil
}
```

### 异步钩子

```go
// 异步执行钩子逻辑
func (u *User) AfterCreate(ctx context.Context) error {
    // 同步操作
    log.Printf("User %s created", u.Name)
    
    // 异步操作
    go func() {
        // 发送欢迎邮件
        if err := sendWelcomeEmail(u.Email, u.Name); err != nil {
            log.Printf("Failed to send welcome email: %v", err)
        }
        
        // 创建用户目录
        if err := createUserDirectory(u.ID); err != nil {
            log.Printf("Failed to create user directory: %v", err)
        }
        
        // 更新统计信息
        updateUserStatistics()
    }()
    
    return nil
}
```

### 错误处理钩子

```go
// 钩子错误处理
func (u *User) BeforeCreate(ctx context.Context) error {
    if err := u.validate(); err != nil {
        // 记录错误
        log.Printf("User validation failed: %v", err)
        
        // 返回用户友好的错误
        return fmt.Errorf("用户数据验证失败: %v", err)
    }
    
    return nil
}

// 钩子恢复机制
func (u *User) AfterCreate(ctx context.Context) error {
    // 尝试执行可能失败的操作
    if err := u.initializeUserData(); err != nil {
        // 记录错误但不中断流程
        log.Printf("Failed to initialize user data: %v", err)
        
        // 可以设置标志或发送通知
        u.InitializationFailed = true
    }
    
    return nil
}
```

## 实际应用场景

### 用户注册流程

```go
type UserRegistration struct {
    Name     string `bson:"name"`
    Email    string `bson:"email"`
    Password string `bson:"password"`
    Status   string `bson:"status"`
}

func (ur *UserRegistration) BeforeCreate(ctx context.Context) error {
    // 验证邮箱唯一性
    session := pie.Table[User](engine)
    exists, err := session.Where("email", ur.Email).Exists(ctx)
    if err != nil {
        return err
    }
    if exists {
        return errors.New("email already exists")
    }
    
    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ur.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    ur.Password = string(hashedPassword)
    
    // 设置默认状态
    ur.Status = "pending"
    
    return nil
}

func (ur *UserRegistration) AfterCreate(ctx context.Context) error {
    // 发送验证邮件
    go func() {
        if err := sendVerificationEmail(ur.Email); err != nil {
            log.Printf("Failed to send verification email: %v", err)
        }
    }()
    
    // 记录注册事件
    logUserEvent("user_registered", ur.Email)
    
    return nil
}
```

### 订单处理流程

```go
type Order struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    UserID      primitive.ObjectID `bson:"user_id"`
    Items       []OrderItem        `bson:"items"`
    Total       float64            `bson:"total"`
    Status      string             `bson:"status"`
    CreatedAt   time.Time          `bson:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at"`
}

func (o *Order) BeforeCreate(ctx context.Context) error {
    // 计算总价
    o.Total = o.calculateTotal()
    
    // 验证库存
    for _, item := range o.Items {
        if err := o.validateStock(item); err != nil {
            return err
        }
    }
    
    // 设置状态
    o.Status = "pending"
    o.CreatedAt = time.Now()
    o.UpdatedAt = time.Now()
    
    return nil
}

func (o *Order) AfterCreate(ctx context.Context) error {
    // 扣减库存
    for _, item := range o.Items {
        if err := o.reduceStock(item); err != nil {
            log.Printf("Failed to reduce stock for item %s: %v", item.ProductID, err)
        }
    }
    
    // 发送订单确认邮件
    go func() {
        if err := sendOrderConfirmation(o); err != nil {
            log.Printf("Failed to send order confirmation: %v", err)
        }
    }()
    
    return nil
}

func (o *Order) BeforeUpdate(ctx context.Context) error {
    // 记录状态变更
    if o.Status != "" {
        logOrderStatusChange(o.ID, o.Status)
    }
    
    o.UpdatedAt = time.Now()
    return nil
}
```

### 审计日志

```go
type AuditLog struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    EntityID  primitive.ObjectID `bson:"entity_id"`
    EntityType string            `bson:"entity_type"`
    Action    string             `bson:"action"`
    Changes   map[string]interface{} `bson:"changes"`
    UserID    primitive.ObjectID `bson:"user_id"`
    Timestamp time.Time          `bson:"timestamp"`
}

func (u *User) BeforeUpdate(ctx context.Context) error {
    // 记录原始数据
    originalUser := *u
    
    // 存储到上下文以供 AfterUpdate 使用
    ctx = context.WithValue(ctx, "original_user", originalUser)
    
    return nil
}

func (u *User) AfterUpdate(ctx context.Context) error {
    // 获取原始数据
    originalUser, ok := ctx.Value("original_user").(User)
    if !ok {
        return nil
    }
    
    // 计算变更
    changes := calculateChanges(originalUser, *u)
    
    // 记录审计日志
    auditLog := &AuditLog{
        EntityID:   u.ID,
        EntityType: "user",
        Action:     "update",
        Changes:    changes,
        UserID:     getCurrentUserID(ctx),
        Timestamp:  time.Now(),
    }
    
    session := pie.Table[AuditLog](engine)
    _, err := session.Insert(ctx, auditLog)
    return err
}
```

## 性能优化

### 1. 避免在钩子中执行重操作

```go
// 好的做法：异步执行重操作
func (u *User) AfterCreate(ctx context.Context) error {
    // 同步操作
    log.Printf("User %s created", u.Name)
    
    // 异步执行重操作
    go func() {
        sendWelcomeEmail(u.Email, u.Name)
        createUserDirectory(u.ID)
        updateUserStatistics()
    }()
    
    return nil
}

// 避免的做法：在钩子中执行重操作
func (u *User) AfterCreate(ctx context.Context) error {
    // 这会阻塞数据库操作
    sendWelcomeEmail(u.Email, u.Name) // 重操作
    createUserDirectory(u.ID)         // 重操作
    updateUserStatistics()            // 重操作
    
    return nil
}
```

### 2. 使用上下文传递数据

```go
func (u *User) BeforeCreate(ctx context.Context) error {
    // 在上下文中存储数据
    ctx = context.WithValue(ctx, "user_metadata", map[string]interface{}{
        "created_by": getCurrentUserID(ctx),
        "created_ip": getClientIP(ctx),
    })
    
    return nil
}

func (u *User) AfterCreate(ctx context.Context) error {
    // 从上下文获取数据
    if metadata, ok := ctx.Value("user_metadata").(map[string]interface{}); ok {
        u.Metadata = metadata
    }
    
    return nil
}
```

### 3. 批量操作优化

```go
// 批量插入时避免重复的钩子操作
func batchCreateUsers(users []*User) error {
    // 预处理所有用户
    for _, user := range users {
        if err := user.BeforeCreate(context.Background()); err != nil {
            return err
        }
    }
    
    // 批量插入
    session := pie.Table[User](engine)
    _, err := session.InsertMany(context.Background(), users)
    if err != nil {
        return err
    }
    
    // 批量后处理
    for _, user := range users {
        user.AfterCreate(context.Background())
    }
    
    return nil
}
```

## 最佳实践

### 1. 钩子方法命名规范

```go
// 使用标准的钩子方法名
func (u *User) BeforeCreate(ctx context.Context) error { return nil }
func (u *User) AfterCreate(ctx context.Context) error { return nil }
func (u *User) BeforeUpdate(ctx context.Context) error { return nil }
func (u *User) AfterUpdate(ctx context.Context) error { return nil }
func (u *User) BeforeDelete(ctx context.Context) error { return nil }
func (u *User) AfterDelete(ctx context.Context) error { return nil }
func (u *User) AfterFind(ctx context.Context) error { return nil }
```

### 2. 错误处理

```go
func (u *User) BeforeCreate(ctx context.Context) error {
    // 验证数据
    if err := u.validate(); err != nil {
        // 记录错误
        log.Printf("User validation failed: %v", err)
        
        // 返回有意义的错误
        return fmt.Errorf("用户数据验证失败: %v", err)
    }
    
    return nil
}
```

### 3. 钩子测试

```go
func TestUserHooks(t *testing.T) {
    // 测试 BeforeCreate 钩子
    user := &User{Name: "Test User", Email: "test@example.com"}
    
    err := user.BeforeCreate(context.Background())
    assert.NoError(t, err)
    assert.NotEmpty(t, user.CreatedAt)
    assert.NotEmpty(t, user.UpdatedAt)
    
    // 测试 AfterCreate 钩子
    err = user.AfterCreate(context.Background())
    assert.NoError(t, err)
}
```

## 下一步

- [软删除](/advanced/soft-delete/) - 学习软删除功能
- [索引管理](/advanced/indexes/) - 掌握索引优化
- [变更流](/advanced/change-streams/) - 学习实时监听
