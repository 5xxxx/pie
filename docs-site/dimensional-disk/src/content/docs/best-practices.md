---
title: Best Practices
description: Learn Pie development best practices and design patterns
---

# Best Practices

This guide provides best practices and design patterns for developing MongoDB applications with Pie.

## Project Structure

### Recommended Directory Structure

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

## Model Design

### Base Model

```go
// internal/models/base.go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
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

## Repository Pattern

### Base Repository

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

func (r *BaseRepository[T]) GetByID(ctx context.Context, id primitive.ObjectID) (*T, error) {
    var entity T
    err := r.session.Where("_id", id).First(ctx, &entity)
    if err != nil {
        return nil, err
    }
    return &entity, nil
}
```

## Error Handling

### Custom Error Types

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

// Predefined errors
var (
    ErrUserNotFound     = &AppError{Code: "USER_NOT_FOUND", Message: "user not found"}
    ErrEmailExists      = &AppError{Code: "EMAIL_EXISTS", Message: "email already exists"}
    ErrInvalidPassword  = &AppError{Code: "INVALID_PASSWORD", Message: "invalid password"}
    ErrUnauthorized     = &AppError{Code: "UNAUTHORIZED", Message: "unauthorized"}
    ErrForbidden        = &AppError{Code: "FORBIDDEN", Message: "forbidden"}
)
```

## Testing Strategy

### Unit Tests

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

func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := &services.UserService{UserRepo: mockRepo}
    
    // Test case: Success
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
}
```

## Performance Optimization

### Connection Pool Configuration

```go
// Optimize connection pool
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMaxPoolSize(100),        // Max connections
    pie.WithMinPoolSize(10),         // Min connections
    pie.WithMaxIdleTime(30*time.Minute), // Max idle time
)
```

### Caching Strategy

```go
// Multi-level caching
engine.UseTwoLevelCache(
    pie.NewMemoryCache(),  // L1 cache
    redisCache,            // L2 cache
    &pie.TwoLevelCacheConfig{
        L1TTL: 1 * time.Minute,  // L1 cache 1 minute
        L2TTL: 10 * time.Minute, // L2 cache 10 minutes
    },
)
```

## Security Best Practices

### Input Validation

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
```

## Summary

Following these best practices will help you build maintainable, scalable, and high-performance MongoDB applications:

1. **Project Structure**: Use clear layered architecture
2. **Model Design**: Use hooks and validation appropriately
3. **Repository Pattern**: Abstract data access layer
4. **Error Handling**: Unified error handling mechanism
5. **Testing Strategy**: Comprehensive unit and integration tests
6. **Performance**: Connection pooling and caching
7. **Security**: Input validation and access control

These practices will help you make the most of Pie's powerful features and build high-quality applications.
