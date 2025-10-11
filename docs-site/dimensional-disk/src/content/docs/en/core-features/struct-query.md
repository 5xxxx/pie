---
title: Struct Query
description: Convert HTTP request parameters directly to query conditions
---

# Struct Query

Convert HTTP request parameters directly to query conditions:

```go
type UserQuery struct {
    Name     string   `pie:"name,like,omitempty" json:"name"`
    Email    string   `pie:"email,omitempty" json:"email"`
    MinAge   int      `pie:"age,gte,omitempty" json:"min_age"`
    MaxAge   int      `pie:"age,lte,omitempty" json:"max_age"`
    Status   []string `pie:"status,in,omitempty" json:"status"`
    Role     string   `pie:"role,omitempty" json:"role"`
}

// HTTP request: GET /users?name=John&min_age=20&max_age=40&status=active,pending
query := UserQuery{
    Name:   "John",
    MinAge: 20,
    MaxAge: 40,
    Status: []string{"active", "pending"},
}

var users []User
err := session.WhereStruct(query).Find(ctx, &users)
```

## Tag Syntax

### Basic Tags

```go
type Query struct {
    // Simple equality
    Name string `pie:"name,omitempty"`
    
    // Field mapping
    UserName string `pie:"username,omitempty"`
    
    // Multiple conditions on same field
    MinAge int `pie:"age,gte,omitempty"`
    MaxAge int `pie:"age,lte,omitempty"`
}
```

### Operators

```go
type Query struct {
    // Comparison operators
    Age     int    `pie:"age,eq,omitempty"`     // Equal
    Score   int    `pie:"score,gt,omitempty"`   // Greater than
    Rating  int    `pie:"rating,gte,omitempty"` // Greater than or equal
    Price   int    `pie:"price,lt,omitempty"`   // Less than
    Amount  int    `pie:"amount,lte,omitempty"` // Less than or equal
    Status  string `pie:"status,ne,omitempty"`  // Not equal
    
    // String operators
    Name    string `pie:"name,like,omitempty"`     // Like
    Email   string `pie:"email,regex,omitempty"`   // Regex
    Domain  string `pie:"domain,startsWith,omitempty"` // Starts with
    
    // Array operators
    Roles   []string `pie:"role,in,omitempty"`        // In array
    Tags    []string `pie:"tags,nin,omitempty"`       // Not in array
    Skills  []string `pie:"skills,all,omitempty"`     // All elements
    
    // Null operators
    Description string `pie:"description,null,omitempty"`     // Is null
    Bio         string `pie:"bio,notNull,omitempty"`          // Is not null
    Avatar      string `pie:"avatar,empty,omitempty"`         // Is empty
    Website     string `pie:"website,notEmpty,omitempty"`     // Is not empty
    
    // Date operators
    CreatedAt time.Time `pie:"created_at,date,omitempty"`     // Date
    UpdatedAt time.Time `pie:"updated_at,dateRange,omitempty"` // Date range
}
```

## Advanced Usage

### Nested Structs

```go
type UserQuery struct {
    Name     string      `pie:"name,like,omitempty"`
    Profile  ProfileQuery `pie:"profile,omitempty"`
    Settings SettingsQuery `pie:"settings,omitempty"`
}

type ProfileQuery struct {
    Age    int    `pie:"age,gte,omitempty"`
    Gender string `pie:"gender,omitempty"`
    City   string `pie:"city,like,omitempty"`
}

type SettingsQuery struct {
    Theme     string `pie:"theme,omitempty"`
    Language  string `pie:"language,omitempty"`
    Notifications bool `pie:"notifications,omitempty"`
}

// Usage
query := UserQuery{
    Name: "John",
    Profile: ProfileQuery{
        Age:    18,
        Gender: "male",
        City:   "New York",
    },
    Settings: SettingsQuery{
        Theme: "dark",
        Language: "en",
    },
}

var users []User
err := session.WhereStruct(query).Find(ctx, &users)
```

### Custom Mappers

```go
type CustomQuery struct {
    UserName string `pie:"user_name,like,omitempty"` // Maps to "username" field
    UserEmail string `pie:"user_email,omitempty"`    // Maps to "email" field
}

// Custom mapper
type CustomMapper struct {
    pie.SnakeMapper
}

func (m CustomMapper) FieldName(fieldName string) string {
    switch fieldName {
    case "user_name":
        return "username"
    case "user_email":
        return "email"
    default:
        return m.SnakeMapper.FieldName(fieldName)
    }
}

// Use custom mapper
engine, err := pie.NewEngine(ctx, "mydb", pie.WithMapper(CustomMapper{}))
```

### Validation Integration

```go
import "github.com/go-playground/validator/v10"

type UserQuery struct {
    Name     string   `pie:"name,like,omitempty" validate:"min=2,max=50"`
    Email    string   `pie:"email,omitempty" validate:"email"`
    MinAge   int      `pie:"age,gte,omitempty" validate:"min=0,max=120"`
    MaxAge   int      `pie:"age,lte,omitempty" validate:"min=0,max=120"`
    Status   []string `pie:"status,in,omitempty" validate:"dive,oneof=active inactive pending"`
}

func (q *UserQuery) Validate() error {
    validate := validator.New()
    return validate.Struct(q)
}

// Usage
query := UserQuery{...}
if err := query.Validate(); err != nil {
    return err
}

var users []User
err := session.WhereStruct(query).Find(ctx, &users)
```

## Real-world Examples

### User Search API

```go
type UserSearchRequest struct {
    Query    string   `pie:"name,like,omitempty" json:"q"`
    Role     string   `pie:"role,omitempty" json:"role"`
    Status   []string `pie:"status,in,omitempty" json:"status"`
    MinAge   int      `pie:"age,gte,omitempty" json:"min_age"`
    MaxAge   int      `pie:"age,lte,omitempty" json:"max_age"`
    City     string   `pie:"profile.city,like,omitempty" json:"city"`
    Country  string   `pie:"profile.country,omitempty" json:"country"`
    SortBy   string   `json:"sort_by"`
    SortDir  string   `json:"sort_dir"`
    Page     int      `json:"page"`
    PageSize int      `json:"page_size"`
}

func SearchUsers(ctx context.Context, req UserSearchRequest) ([]User, error) {
    session := pie.Table[User](engine)
    
    // Apply struct query
    query := session.WhereStruct(req)
    
    // Apply sorting
    if req.SortBy != "" {
        if req.SortDir == "desc" {
            query = query.OrderByDesc(req.SortBy)
        } else {
            query = query.OrderBy(req.SortBy)
        }
    }
    
    // Apply pagination
    if req.Page > 0 && req.PageSize > 0 {
        query = query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
    }
    
    var users []User
    err := query.Find(ctx, &users)
    return users, err
}
```

### Product Filter API

```go
type ProductFilter struct {
    Name        string    `pie:"name,like,omitempty" json:"name"`
    Category    string    `pie:"category,omitempty" json:"category"`
    MinPrice    float64   `pie:"price,gte,omitempty" json:"min_price"`
    MaxPrice    float64   `pie:"price,lte,omitempty" json:"max_price"`
    Brand       []string  `pie:"brand,in,omitempty" json:"brand"`
    InStock     bool      `pie:"stock,gt,omitempty" json:"in_stock"`
    Tags        []string  `pie:"tags,all,omitempty" json:"tags"`
    CreatedFrom time.Time `pie:"created_at,gte,omitempty" json:"created_from"`
    CreatedTo   time.Time `pie:"created_at,lte,omitempty" json:"created_to"`
}

func FilterProducts(ctx context.Context, filter ProductFilter) ([]Product, error) {
    session := pie.Table[Product](engine)
    
    var products []Product
    err := session.WhereStruct(filter).Find(ctx, &products)
    return products, err
}
```

### Order Query API

```go
type OrderQuery struct {
    UserID      string    `pie:"user_id,omitempty" json:"user_id"`
    Status      []string  `pie:"status,in,omitempty" json:"status"`
    MinAmount   float64   `pie:"total,gte,omitempty" json:"min_amount"`
    MaxAmount   float64   `pie:"total,lte,omitempty" json:"max_amount"`
    DateFrom    time.Time `pie:"created_at,gte,omitempty" json:"date_from"`
    DateTo      time.Time `pie:"created_at,lte,omitempty" json:"date_to"`
    PaymentMethod string  `pie:"payment_method,omitempty" json:"payment_method"`
}

func QueryOrders(ctx context.Context, query OrderQuery) ([]Order, error) {
    session := pie.Table[Order](engine)
    
    var orders []Order
    err := session.WhereStruct(query).Find(ctx, &orders)
    return orders, err
}
```

## Best Practices

### 1. Use Appropriate Operators

```go
// Good: Use specific operators
type Query struct {
    Name string `pie:"name,like,omitempty"`     // For search
    Age  int    `pie:"age,gte,omitempty"`       // For range
    Role string `pie:"role,omitempty"`          // For exact match
}

// Avoid: Using wrong operators
type BadQuery struct {
    Name string `pie:"name,omitempty"`          // Should be 'like' for search
    Age  int    `pie:"age,omitempty"`           // Should be 'gte' for range
}
```

### 2. Handle Empty Values

```go
// Good: Use omitempty to skip empty values
type Query struct {
    Name string `pie:"name,like,omitempty"`
    Age  int    `pie:"age,gte,omitempty"`
}

// Bad: Without omitempty, zero values will be included
type BadQuery struct {
    Name string `pie:"name,like"`
    Age  int    `pie:"age,gte"`
}
```

### 3. Validate Input

```go
type Query struct {
    Name   string `pie:"name,like,omitempty" validate:"min=2,max=50"`
    Age    int    `pie:"age,gte,omitempty" validate:"min=0,max=120"`
    Email  string `pie:"email,omitempty" validate:"email"`
}

func (q *Query) Validate() error {
    validate := validator.New()
    return validate.Struct(q)
}
```

### 4. Use Meaningful Field Names

```go
// Good: Clear and descriptive
type Query struct {
    UserName    string `pie:"user_name,like,omitempty"`
    UserEmail   string `pie:"user_email,omitempty"`
    MinAge      int    `pie:"age,gte,omitempty"`
    MaxAge      int    `pie:"age,lte,omitempty"`
}

// Bad: Unclear field names
type BadQuery struct {
    N string `pie:"n,like,omitempty"`
    E string `pie:"e,omitempty"`
    A int    `pie:"a,gte,omitempty"`
}
```

## Next Steps

- [Pagination](/core-features/pagination/) - Implement pagination
- [Query Builder](/core-features/query-builder/) - Learn about query methods
- [Transactions](/core-features/transactions/) - Use transactions
- [Best Practices](/best-practices/) - Learn development best practices
