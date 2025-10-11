---
title: 最佳实践
description: 学习 Pie 开发的最佳实践和设计模式
---

# 最佳实践

本指南提供了使用 Pie 开发 MongoDB 应用程序的最佳实践和设计模式。

## 项目结构

### 推荐的目录结构

```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── connection.go
│   ├── models/
│   │   ├── user.go
│   │   ├── order.go
│   │   └── product.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── order_repository.go
│   │   └── product_repository.go
│   ├── services/
│   │   ├── user_service.go
│   │   ├── order_service.go
│   │   └── product_service.go
│   └── handlers/
│       ├── user_handler.go
│       ├── order_handler.go
│       └── product_handler.go
├── pkg/
│   └── utils/
│       └── validation.go
├── migrations/
│   └── indexes.go
├── tests/
│   ├── integration/
│   └── unit/
└── go.mod
```

### 配置管理

```go
// internal/config/config.go
package config

import (
    "os"
    "time"
)

type Config struct {
    Database struct {
        URI      string `yaml:"uri"`
        Database string `yaml:"database"`
        Auth     struct {
            Username string `yaml:"username"`
            Password string `yaml:"password"`
        } `yaml:"auth"`
        Pool struct {
            MaxSize     int           `yaml:"max_size"`
            MinSize     int           `yaml:"min_size"`
            MaxIdleTime time.Duration `yaml:"max_idle_time"`
        } `yaml:"pool"`
    } `yaml:"database"`
    
    Cache struct {
        Type string        `yaml:"type"`
        TTL  time.Duration `yaml:"ttl"`
        Redis struct {
            Host     string `yaml:"host"`
            Port     int    `yaml:"port"`
            Password string `yaml:"password"`
            DB       int    `yaml:"db"`
        } `yaml:"redis"`
    } `yaml:"cache"`
}

func Load() (*Config, error) {
    config := &Config{}
    
    // 从环境变量加载配置
    config.Database.URI = os.Getenv("DATABASE_URI")
    config.Database.Database = os.Getenv("DATABASE_NAME")
    config.Database.Auth.Username = os.Getenv("DATABASE_USERNAME")
    config.Database.Auth.Password = os.Getenv("DATABASE_PASSWORD")
    
    // 设置默认值
    if config.Database.Pool.MaxSize == 0 {
        config.Database.Pool.MaxSize = 100
    }
    if config.Database.Pool.MinSize == 0 {
        config.Database.Pool.MinSize = 10
    }
    
    return config, nil
}
```

## 模型设计

### 基础模型

```go
// internal/models/base.go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type BaseModel struct {
    ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" pie:"soft_delete"`
}

func (m *BaseModel) BeforeCreate(ctx context.Context) error {
    now := time.Now()
    m.CreatedAt = now
    m.UpdatedAt = now
    return nil
}

func (m *BaseModel) BeforeUpdate(ctx context.Context) error {
    m.UpdatedAt = time.Now()
    return nil
}
```

### 用户模型

```go
// internal/models/user.go
package models

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
    BaseModel
    Name      string    `bson:"name" json:"name" pie:"index"`
    Email     string    `bson:"email" json:"email" pie:"unique"`
    Password  string    `bson:"password" json:"-"`
    Role      string    `bson:"role" json:"role" pie:"index"`
    Status    string    `bson:"status" json:"status" pie:"index"`
    LastLogin *time.Time `bson:"last_login,omitempty" json:"last_login,omitempty"`
    Profile   *Profile  `bson:"profile,omitempty" json:"profile,omitempty"`
}

type Profile struct {
    FirstName string `bson:"first_name" json:"first_name"`
    LastName  string `bson:"last_name" json:"last_name"`
    Avatar    string `bson:"avatar" json:"avatar"`
    Bio       string `bson:"bio" json:"bio"`
}

func (u *User) BeforeCreate(ctx context.Context) error {
    if err := u.BaseModel.BeforeCreate(ctx); err != nil {
        return err
    }
    
    // 设置默认值
    if u.Status == "" {
        u.Status = "active"
    }
    if u.Role == "" {
        u.Role = "user"
    }
    
    // 加密密码
    hashedPassword, err := hashPassword(u.Password)
    if err != nil {
        return err
    }
    u.Password = hashedPassword
    
    return nil
}

func (u *User) AfterFind(ctx context.Context) error {
    // 隐藏敏感信息
    u.Password = ""
    return nil
}
```

### 订单模型

```go
// internal/models/order.go
package models

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type Order struct {
    BaseModel
    UserID      bson.ObjectID `bson:"user_id" json:"user_id" pie:"index"`
    Items       []OrderItem        `bson:"items" json:"items"`
    Total       float64            `bson:"total" json:"total"`
    Status      string             `bson:"status" json:"status" pie:"index"`
    PaymentID   string             `bson:"payment_id,omitempty" json:"payment_id,omitempty"`
    ShippingAddress *Address       `bson:"shipping_address,omitempty" json:"shipping_address,omitempty"`
}

type OrderItem struct {
    ProductID bson.ObjectID `bson:"product_id" json:"product_id"`
    Name      string             `bson:"name" json:"name"`
    Price     float64            `bson:"price" json:"price"`
    Quantity  int                `bson:"quantity" json:"quantity"`
}

type Address struct {
    Street    string `bson:"street" json:"street"`
    City      string `bson:"city" json:"city"`
    State     string `bson:"state" json:"state"`
    ZipCode   string `bson:"zip_code" json:"zip_code"`
    Country   string `bson:"country" json:"country"`
}

func (o *Order) BeforeCreate(ctx context.Context) error {
    if err := o.BaseModel.BeforeCreate(ctx); err != nil {
        return err
    }
    
    // 设置默认状态
    if o.Status == "" {
        o.Status = "pending"
    }
    
    // 计算总价
    o.Total = o.calculateTotal()
    
    return nil
}

func (o *Order) calculateTotal() float64 {
    total := 0.0
    for _, item := range o.Items {
        total += item.Price * float64(item.Quantity)
    }
    return total
}
```

## 仓储模式

### 基础仓储

```go
// internal/repositories/base_repository.go
package repositories

import (
    "context"
    "github.com/5xxxx/pie"
)

type BaseRepository[T any] struct {
    engine  *pie.Engine
    session *pie.Session[T]
}

func NewBaseRepository[T any](engine *pie.Engine) *BaseRepository[T] {
    return &BaseRepository[T]{
        engine:  engine,
        session: pie.Table[T](engine),
    }
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
    _, err := r.session.Insert(ctx, entity)
    return err
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id bson.ObjectID) (*T, error) {
    var entity T
    err := r.session.Where("_id", id).First(ctx, &entity)
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, id bson.ObjectID, updates bson.D) error {
    result, err := r.session.Where("_id", id).Update(ctx, updates)
    if err != nil {
        return err
    }
    if result.ModifiedCount == 0 {
        return errors.New("entity not found")
    }
    return nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id bson.ObjectID) error {
    result, err := r.session.Where("_id", id).Delete(ctx)
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return errors.New("entity not found")
    }
    return nil
}

func (r *BaseRepository[T]) SoftDelete(ctx context.Context, id bson.ObjectID) error {
    return r.session.Where("_id", id).SoftDelete(ctx)
}
```

### 用户仓储

```go
// internal/repositories/user_repository.go
package repositories

import (
    "context"
    "github.com/5xxxx/pie"
    "your-project/internal/models"
)

type UserRepository struct {
    *BaseRepository[models.User]
}

func NewUserRepository(engine *pie.Engine) *UserRepository {
    return &UserRepository{
        BaseRepository: NewBaseRepository[models.User](engine),
    }
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.session.Where("email", email).First(ctx, &user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetActiveUsers(ctx context.Context) ([]models.User, error) {
    var users []models.User
    err := r.session.
        Where("status", "active").
        OrderBy("created_at").
        Find(ctx, &users)
    return users, err
}

func (r *UserRepository) GetUsersByRole(ctx context.Context, role string) ([]models.User, error) {
    var users []models.User
    err := r.session.
        Where("role", role).
        Where("status", "active").
        Find(ctx, &users)
    return users, err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID bson.ObjectID) error {
    now := time.Now()
    return r.Update(ctx, userID, bson.D{{"$set", bson.D{{"last_login", now}}}})
}
```

## 服务层

### 用户服务

```go
// internal/services/user_service.go
package services

import (
    "context"
    "errors"
    "github.com/5xxxx/pie"
    "your-project/internal/models"
    "your-project/internal/repositories"
)

type UserService struct {
    userRepo *repositories.UserRepository
}

func NewUserService(engine *pie.Engine) *UserService {
    return &UserService{
        userRepo: repositories.NewUserRepository(engine),
    }
}

func (s *UserService) CreateUser(ctx context.Context, userData *CreateUserRequest) (*models.User, error) {
    // 验证邮箱唯一性
    existingUser, err := s.userRepo.GetByEmail(ctx, userData.Email)
    if err != nil && !pie.IsNotFoundError(err) {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("email already exists")
    }
    
    // 创建用户
    user := &models.User{
        Name:     userData.Name,
        Email:    userData.Email,
        Password: userData.Password,
        Role:     userData.Role,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) GetUser(ctx context.Context, userID bson.ObjectID) (*models.User, error) {
    return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) UpdateUser(ctx context.Context, userID bson.ObjectID, updates *UpdateUserRequest) error {
    updateDoc := bson.D{}
    
    if updates.Name != "" {
        updateDoc = append(updateDoc, bson.E{"$set", bson.D{{"name", updates.Name}}})
    }
    if updates.Email != "" {
        updateDoc = append(updateDoc, bson.E{"$set", bson.D{{"email", updates.Email}}})
    }
    
    return s.userRepo.Update(ctx, userID, updateDoc)
}

func (s *UserService) DeleteUser(ctx context.Context, userID bson.ObjectID) error {
    return s.userRepo.SoftDelete(ctx, userID)
}

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Role     string `json:"role" validate:"oneof=user admin"`
}

type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email" validate:"email"`
}
```

## 错误处理

### 自定义错误类型

```go
// pkg/errors/errors.go
package errors

import "fmt"

type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// 预定义错误
var (
    ErrUserNotFound     = &AppError{Code: "USER_NOT_FOUND", Message: "user not found"}
    ErrEmailExists      = &AppError{Code: "EMAIL_EXISTS", Message: "email already exists"}
    ErrInvalidPassword  = &AppError{Code: "INVALID_PASSWORD", Message: "invalid password"}
    ErrUnauthorized     = &AppError{Code: "UNAUTHORIZED", Message: "unauthorized"}
    ErrForbidden        = &AppError{Code: "FORBIDDEN", Message: "forbidden"}
)
```

### 错误处理中间件

```go
// internal/handlers/middleware.go
package handlers

import (
    "encoding/json"
    "net/http"
    "your-project/pkg/errors"
)

func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                handleError(w, err)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

func handleError(w http.ResponseWriter, err interface{}) {
    w.Header().Set("Content-Type", "application/json")
    
    switch e := err.(type) {
    case *errors.AppError:
        w.WriteHeader(getHTTPStatus(e.Code))
        json.NewEncoder(w).Encode(map[string]string{
            "error": e.Message,
            "code":  e.Code,
        })
    case error:
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "internal server error",
        })
    default:
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "unknown error",
        })
    }
}

func getHTTPStatus(code string) int {
    switch code {
    case "USER_NOT_FOUND":
        return http.StatusNotFound
    case "EMAIL_EXISTS":
        return http.StatusConflict
    case "UNAUTHORIZED":
        return http.StatusUnauthorized
    case "FORBIDDEN":
        return http.StatusForbidden
    default:
        return http.StatusInternalServerError
    }
}
```

## 测试策略

### 单元测试

```go
// tests/unit/user_service_test.go
package unit

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "your-project/internal/services"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := &services.UserService{UserRepo: mockRepo}
    
    // 测试用例：成功创建用户
    t.Run("Success", func(t *testing.T) {
        mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), errors.New("not found"))
        mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
        
        userData := &services.CreateUserRequest{
            Name:     "Test User",
            Email:    "test@example.com",
            Password: "password123",
            Role:     "user",
        }
        
        user, err := service.CreateUser(context.Background(), userData)
        
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "Test User", user.Name)
        assert.Equal(t, "test@example.com", user.Email)
        
        mockRepo.AssertExpectations(t)
    })
    
    // 测试用例：邮箱已存在
    t.Run("EmailExists", func(t *testing.T) {
        existingUser := &models.User{Email: "test@example.com"}
        mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)
        
        userData := &services.CreateUserRequest{
            Name:     "Test User",
            Email:    "test@example.com",
            Password: "password123",
            Role:     "user",
        }
        
        user, err := service.CreateUser(context.Background(), userData)
        
        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Contains(t, err.Error(), "email already exists")
        
        mockRepo.AssertExpectations(t)
    })
}
```

### 集成测试

```go
// tests/integration/user_integration_test.go
package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/5xxxx/pie"
    "your-project/internal/services"
)

func TestUserIntegration(t *testing.T) {
    // 设置测试数据库
    engine, err := createTestEngine()
    require.NoError(t, err)
    defer engine.Disconnect(context.Background())
    
    service := services.NewUserService(engine)
    
    t.Run("CreateAndGetUser", func(t *testing.T) {
        // 创建用户
        userData := &services.CreateUserRequest{
            Name:     "Test User",
            Email:    "test@example.com",
            Password: "password123",
            Role:     "user",
        }
        
        user, err := service.CreateUser(context.Background(), userData)
        require.NoError(t, err)
        assert.NotNil(t, user.ID)
        
        // 获取用户
        retrievedUser, err := service.GetUser(context.Background(), user.ID)
        require.NoError(t, err)
        assert.Equal(t, user.ID, retrievedUser.ID)
        assert.Equal(t, user.Name, retrievedUser.Name)
        assert.Equal(t, user.Email, retrievedUser.Email)
    })
}

func createTestEngine() (*pie.Engine, error) {
    return pie.NewEngine(
        context.Background(),
        "test_db",
        pie.WithURI("mongodb://localhost:27017"),
        pie.WithMapper(&pie.SnakeMapper{}),
    )
}
```

## 部署和运维

### 健康检查

```go
// internal/handlers/health.go
package handlers

import (
    "context"
    "net/http"
    "time"
    "github.com/5xxxx/pie"
)

type HealthHandler struct {
    engine *pie.Engine
}

func NewHealthHandler(engine *pie.Engine) *HealthHandler {
    return &HealthHandler{engine: engine}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 检查数据库连接
    if err := h.engine.Ping(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Database connection failed"))
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
    // 检查应用是否准备好接收请求
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Ready"))
}
```

### 指标监控

```go
// internal/monitoring/metrics.go
package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // 数据库操作指标
    dbOperationsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_operations_total",
            Help: "Total number of database operations",
        },
        []string{"operation", "collection", "status"},
    )
    
    dbOperationDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_operation_duration_seconds",
            Help:    "Duration of database operations",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation", "collection"},
    )
    
    // 缓存指标
    cacheHitsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache_type"},
    )
    
    cacheMissesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_misses_total",
            Help: "Total number of cache misses",
        },
        []string{"cache_type"},
    )
)

func RecordDBOperation(operation, collection, status string, duration float64) {
    dbOperationsTotal.WithLabelValues(operation, collection, status).Inc()
    dbOperationDuration.WithLabelValues(operation, collection).Observe(duration)
}

func RecordCacheHit(cacheType string) {
    cacheHitsTotal.WithLabelValues(cacheType).Inc()
}

func RecordCacheMiss(cacheType string) {
    cacheMissesTotal.WithLabelValues(cacheType).Inc()
}
```

## 安全最佳实践

### 输入验证

```go
// pkg/validation/validation.go
package validation

import (
    "regexp"
    "strings"
    "unicode"
)

func ValidateEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}

func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters long")
    }
    
    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    if !hasUpper {
        return errors.New("password must contain at least one uppercase letter")
    }
    if !hasLower {
        return errors.New("password must contain at least one lowercase letter")
    }
    if !hasNumber {
        return errors.New("password must contain at least one number")
    }
    if !hasSpecial {
        return errors.New("password must contain at least one special character")
    }
    
    return nil
}

func SanitizeString(input string) string {
    // 移除危险字符
    input = strings.TrimSpace(input)
    input = strings.ReplaceAll(input, "<", "&lt;")
    input = strings.ReplaceAll(input, ">", "&gt;")
    input = strings.ReplaceAll(input, "\"", "&quot;")
    input = strings.ReplaceAll(input, "'", "&#x27;")
    return input
}
```

### 权限控制

```go
// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "your-project/pkg/errors"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        user, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := r.Context().Value("user").(*models.User)
            if user.Role != role {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func extractToken(r *http.Request) string {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return ""
    }
    
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return ""
    }
    
    return parts[1]
}
```

## 总结

遵循这些最佳实践可以帮助您构建可维护、可扩展、高性能的 MongoDB 应用程序：

1. **项目结构**: 使用清晰的分层架构
2. **模型设计**: 合理使用钩子和验证
3. **仓储模式**: 抽象数据访问层
4. **错误处理**: 统一的错误处理机制
5. **测试策略**: 全面的单元测试和集成测试
6. **监控运维**: 健康检查和指标监控
7. **安全实践**: 输入验证和权限控制

这些实践将帮助您充分利用 Pie 的强大功能，构建高质量的应用程序。
