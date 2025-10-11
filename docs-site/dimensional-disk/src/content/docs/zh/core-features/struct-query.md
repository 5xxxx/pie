---
title: 结构体查询
description: 学习如何将 HTTP 请求参数直接转换为查询条件
---

# 结构体查询

Pie 的结构体查询功能允许您将 HTTP 请求参数直接转换为查询条件，大大简化了 API 开发流程。

## 基本用法

### 定义查询结构体

```go
type UserQuery struct {
    Name     string   `pie:"name,like,omitempty" json:"name"`
    Email    string   `pie:"email,omitempty" json:"email"`
    MinAge   int      `pie:"age,gte,omitempty" json:"min_age"`
    MaxAge   int      `pie:"age,lte,omitempty" json:"max_age"`
    Status   []string `pie:"status,in,omitempty" json:"status"`
    Role     string   `pie:"role,omitempty" json:"role"`
    IsActive *bool    `pie:"is_active,omitempty" json:"is_active"`
}
```

### 标签说明

Pie 结构体标签的格式为：`pie:"field,operator,omitempty"`

- `field`: 数据库字段名
- `operator`: 查询操作符（可选）
- `omitempty`: 如果字段为零值则忽略（可选）

### 支持的查询操作符

```go
type ProductQuery struct {
    // 等值查询
    Name        string `pie:"name,omitempty"`
    Category    string `pie:"category,omitempty"`
    
    // 比较查询
    MinPrice    float64 `pie:"price,gte,omitempty"`
    MaxPrice    float64 `pie:"price,lte,omitempty"`
    MinStock    int     `pie:"stock,gt,omitempty"`
    MaxStock    int     `pie:"stock,lt,omitempty"`
    
    // 包含查询
    Categories  []string `pie:"category,in,omitempty"`
    Tags        []string `pie:"tags,all,omitempty"`
    ExcludeTags []string `pie:"tags,nin,omitempty"`
    
    // 模糊查询
    Search      string `pie:"name,like,omitempty"`
    Description string `pie:"description,regex,omitempty"`
    
    // 空值查询
    HasImage    *bool `pie:"image,exists,omitempty"`
    NoImage     *bool `pie:"image,not_exists,omitempty"`
    
    // 数组查询
    TagCount    int `pie:"tags,size,omitempty"`
    
    // 日期查询
    CreatedAfter  time.Time `pie:"created_at,gte,omitempty"`
    CreatedBefore time.Time `pie:"created_at,lte,omitempty"`
}
```

## 使用示例

### HTTP 请求处理

```go
// HTTP 请求: GET /users?name=张&min_age=20&max_age=40&status=active,pending&role=user
func GetUsers(c *gin.Context) {
    var query UserQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    session := pie.Table[User](engine)
    var users []User
    err := session.WhereStruct(query).Find(context.Background(), &users)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"data": users})
}
```

### 手动构建查询

```go
query := UserQuery{
    Name:   "张",
    MinAge: 20,
    MaxAge: 40,
    Status: []string{"active", "pending"},
    Role:   "user",
}

session := pie.Table[User](engine)
var users []User
err := session.WhereStruct(query).Find(ctx, &users)
```

## 高级用法

### 嵌套结构体

```go
type SearchQuery struct {
    User UserQuery `pie:"user"`
    Product ProductQuery `pie:"product"`
    DateRange struct {
        Start time.Time `pie:"created_at,gte,omitempty"`
        End   time.Time `pie:"created_at,lte,omitempty"`
    } `pie:"date_range"`
}

// 使用嵌套查询
var searchQuery SearchQuery
// ... 填充数据
err := session.WhereStruct(searchQuery).Find(ctx, &results)
```

### 自定义字段映射

```go
type CustomQuery struct {
    UserName    string `pie:"name,like,omitempty" json:"user_name"`
    UserEmail   string `pie:"email,omitempty" json:"user_email"`
    UserAge     int    `pie:"age,gte,omitempty" json:"user_age"`
    UserStatus  string `pie:"status,omitempty" json:"user_status"`
}
```

### 条件查询

```go
type ConditionalQuery struct {
    // 只有指定了值才会应用条件
    Name     string `pie:"name,like,omitempty"`
    Email    string `pie:"email,omitempty"`
    
    // 使用指针类型支持 null 值
    IsActive *bool  `pie:"is_active,omitempty"`
    IsVip    *bool  `pie:"is_vip,omitempty"`
}

// 示例
query := ConditionalQuery{
    Name: "张",
    // IsActive 为 nil，不会添加到查询条件中
    // IsVip 为 nil，不会添加到查询条件中
}
```

## 复杂查询场景

### 多条件组合

```go
type AdvancedQuery struct {
    // 基础条件
    Name        string    `pie:"name,like,omitempty"`
    Email       string    `pie:"email,omitempty"`
    
    // 数值范围
    MinAge      int       `pie:"age,gte,omitempty"`
    MaxAge      int       `pie:"age,lte,omitempty"`
    MinScore    float64   `pie:"score,gte,omitempty"`
    MaxScore    float64   `pie:"score,lte,omitempty"`
    
    // 数组条件
    Roles       []string  `pie:"role,in,omitempty"`
    Tags        []string  `pie:"tags,all,omitempty"`
    ExcludeTags []string  `pie:"tags,nin,omitempty"`
    
    // 日期范围
    CreatedAfter  time.Time `pie:"created_at,gte,omitempty"`
    CreatedBefore time.Time `pie:"created_at,lte,omitempty"`
    
    // 布尔条件
    IsActive   *bool `pie:"is_active,omitempty"`
    IsVerified *bool `pie:"is_verified,omitempty"`
    
    // 空值条件
    HasPhone   *bool `pie:"phone,exists,omitempty"`
    NoAddress  *bool `pie:"address,not_exists,omitempty"`
}
```

### 分页查询

```go
type PaginatedQuery struct {
    // 查询条件
    UserQuery
    
    // 分页参数
    Page     int `json:"page" form:"page"`
    PageSize int `json:"page_size" form:"page_size"`
    
    // 排序
    SortBy    string `json:"sort_by" form:"sort_by"`
    SortOrder string `json:"sort_order" form:"sort_order"`
}

func GetUsersWithPagination(c *gin.Context) {
    var query PaginatedQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 设置默认值
    if query.Page <= 0 {
        query.Page = 1
    }
    if query.PageSize <= 0 {
        query.PageSize = 10
    }
    if query.SortBy == "" {
        query.SortBy = "created_at"
    }
    if query.SortOrder == "" {
        query.SortOrder = "desc"
    }
    
    session := pie.Table[User](engine)
    
    // 构建查询
    q := session.WhereStruct(query.UserQuery)
    
    // 添加排序
    if query.SortOrder == "desc" {
        q = q.OrderByDesc(query.SortBy)
    } else {
        q = q.OrderBy(query.SortBy)
    }
    
    // 执行分页查询
    result, err := q.Paginate(context.Background(), pie.PaginateParams{
        Page:     query.Page,
        PageSize: query.PageSize,
    })
    
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "data":       result.Data,
        "total":      result.Total,
        "page":       result.Page,
        "page_size":  result.PageSize,
        "total_pages": result.TotalPages,
        "has_next":   result.HasNext,
        "has_prev":   result.HasPrev,
    })
}
```

## 验证和错误处理

### 结构体验证

```go
import "github.com/go-playground/validator/v10"

type ValidatedQuery struct {
    Name     string `pie:"name,like,omitempty" json:"name" validate:"max=50"`
    Email    string `pie:"email,omitempty" json:"email" validate:"email"`
    MinAge   int    `pie:"age,gte,omitempty" json:"min_age" validate:"min=0,max=150"`
    MaxAge   int    `pie:"age,lte,omitempty" json:"max_age" validate:"min=0,max=150"`
    Status   []string `pie:"status,in,omitempty" json:"status" validate:"dive,oneof=active pending inactive"`
}

func ValidateQuery(query *ValidatedQuery) error {
    validate := validator.New()
    return validate.Struct(query)
}
```

### 自定义验证

```go
type CustomValidatedQuery struct {
    Name     string `pie:"name,like,omitempty" json:"name"`
    MinAge   int    `pie:"age,gte,omitempty" json:"min_age"`
    MaxAge   int    `pie:"age,lte,omitempty" json:"max_age"`
}

func (q *CustomValidatedQuery) Validate() error {
    if q.MinAge > 0 && q.MaxAge > 0 && q.MinAge > q.MaxAge {
        return errors.New("min_age cannot be greater than max_age")
    }
    return nil
}
```

## 性能优化

### 索引优化

```go
// 确保查询字段有适当的索引
type OptimizedQuery struct {
    Email    string `pie:"email,omitempty"`     // 应该有唯一索引
    Status   string `pie:"status,omitempty"`    // 应该有索引
    CreatedAt time.Time `pie:"created_at,gte,omitempty"` // 应该有索引
}

// 复合索引建议：(status, created_at)
```

### 字段投影

```go
// 只选择需要的字段
type ProjectedQuery struct {
    Name     string `pie:"name,like,omitempty"`
    Email    string `pie:"email,omitempty"`
    Status   string `pie:"status,omitempty"`
}

func GetUsersOptimized(query ProjectedQuery) ([]User, error) {
    session := pie.Table[User](engine)
    var users []User
    
    return users, session.
        WhereStruct(query).
        Select("name", "email", "status"). // 只选择需要的字段
        Find(context.Background(), &users)
}
```

## 最佳实践

### 1. 使用指针类型处理可选布尔值

```go
type GoodQuery struct {
    Name     string `pie:"name,like,omitempty"`
    IsActive *bool  `pie:"is_active,omitempty"` // 使用指针，支持 null
}

type BadQuery struct {
    Name     string `pie:"name,like,omitempty"`
    IsActive bool   `pie:"is_active,omitempty"` // 布尔值无法区分 false 和未设置
}
```

### 2. 合理使用 omitempty 标签

```go
type GoodQuery struct {
    Name     string `pie:"name,like,omitempty"` // 字符串零值 "" 应该被忽略
    Age      int    `pie:"age,gte,omitempty"`   // 整数零值 0 应该被忽略
    IsActive *bool  `pie:"is_active,omitempty"` // 指针零值 nil 应该被忽略
}
```

### 3. 避免过于复杂的查询结构体

```go
// 好的做法：拆分为多个简单的查询结构体
type UserBasicQuery struct {
    Name  string `pie:"name,like,omitempty"`
    Email string `pie:"email,omitempty"`
}

type UserAdvancedQuery struct {
    UserBasicQuery
    MinAge int      `pie:"age,gte,omitempty"`
    MaxAge int      `pie:"age,lte,omitempty"`
    Roles  []string `pie:"role,in,omitempty"`
}
```

## 下一步

- [分页查询](/core-features/pagination/) - 学习分页实现
- [聚合查询](/core-features/aggregation/) - 掌握复杂数据聚合
- [缓存支持](/advanced/cache/) - 提升查询性能
