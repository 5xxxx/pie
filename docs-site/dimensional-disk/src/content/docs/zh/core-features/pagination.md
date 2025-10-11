---
title: 分页查询
description: 学习 Pie 的高效分页实现，支持简单和完整分页模式
---

# 分页查询

Pie 提供了高效的分页实现，支持两种分页模式：完整分页（包含总数统计）和简单分页（不统计总数，性能更好）。

## 基本用法

### 完整分页

完整分页会统计总记录数，提供完整的分页信息：

```go
// 分页参数
params := pie.PaginateParams{
    Page:     1,  // 页码，从 1 开始
    PageSize: 10, // 每页大小
}

// 执行分页查询
result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    Paginate(ctx, params)

if err != nil {
    log.Fatal("Failed to paginate:", err)
}

// 使用结果
fmt.Printf("总记录数: %d\n", result.Total)
fmt.Printf("总页数: %d\n", result.TotalPages)
fmt.Printf("当前页: %d\n", result.Page)
fmt.Printf("每页大小: %d\n", result.PageSize)
fmt.Printf("是否有下一页: %t\n", result.HasNext)
fmt.Printf("是否有上一页: %t\n", result.HasPrev)
fmt.Printf("数据: %+v\n", result.Data)
```

### 简单分页

简单分页不统计总记录数，性能更好，适合大数据量场景：

```go
// 简单分页参数
params := pie.PaginateParams{
    Page:     1,
    PageSize: 10,
}

// 执行简单分页查询
result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    PaginateSimple(ctx, params)

if err != nil {
    log.Fatal("Failed to paginate:", err)
}

// 简单分页结果不包含 Total 和 TotalPages
fmt.Printf("当前页: %d\n", result.Page)
fmt.Printf("每页大小: %d\n", result.PageSize)
fmt.Printf("是否有下一页: %t\n", result.HasNext)
fmt.Printf("是否有上一页: %t\n", result.HasPrev)
fmt.Printf("数据: %+v\n", result.Data)
```

## 分页结果结构

### 完整分页结果

```go
type PaginateResult[T any] struct {
    Data       []T   `json:"data"`        // 当前页数据
    Total      int64 `json:"total"`       // 总记录数
    TotalPages int   `json:"total_pages"` // 总页数
    Page       int   `json:"page"`        // 当前页码
    PageSize   int   `json:"page_size"`   // 每页大小
    HasNext    bool  `json:"has_next"`    // 是否有下一页
    HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}
```

### 简单分页结果

```go
type SimplePaginateResult[T any] struct {
    Data     []T  `json:"data"`      // 当前页数据
    Page     int  `json:"page"`      // 当前页码
    PageSize int  `json:"page_size"` // 每页大小
    HasNext  bool `json:"has_next"`  // 是否有下一页
    HasPrev  bool `json:"has_prev"`  // 是否有上一页
}
```

## 实际应用示例

### HTTP API 分页

```go
// 分页请求结构
type PaginationRequest struct {
    Page     int    `form:"page" json:"page" binding:"min=1"`
    PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
    SortBy   string `form:"sort_by" json:"sort_by"`
    SortOrder string `form:"sort_order" json:"sort_order" binding:"oneof=asc desc"`
}

// 分页响应结构
type PaginationResponse[T any] struct {
    Data       []T   `json:"data"`
    Total      int64 `json:"total,omitempty"`
    TotalPages int   `json:"total_pages,omitempty"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
}

// 用户列表 API
func GetUsers(c *gin.Context) {
    var req PaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 设置默认值
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    if req.SortBy == "" {
        req.SortBy = "created_at"
    }
    if req.SortOrder == "" {
        req.SortOrder = "desc"
    }
    
    session := pie.Table[User](engine)
    query := session.Where("status", "active")
    
    // 添加排序
    if req.SortOrder == "desc" {
        query = query.OrderByDesc(req.SortBy)
    } else {
        query = query.OrderBy(req.SortBy)
    }
    
    // 执行分页查询
    params := pie.PaginateParams{
        Page:     req.Page,
        PageSize: req.PageSize,
    }
    
    result, err := query.Paginate(context.Background(), params)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 返回响应
    response := PaginationResponse[User]{
        Data:       result.Data,
        Total:      result.Total,
        TotalPages: result.TotalPages,
        Page:       result.Page,
        PageSize:   result.PageSize,
        HasNext:    result.HasNext,
        HasPrev:    result.HasPrev,
    }
    
    c.JSON(200, response)
}
```

### 条件分页查询

```go
// 用户搜索分页
func SearchUsers(c *gin.Context) {
    var req struct {
        PaginationRequest
        Name     string   `form:"name" json:"name"`
        Email    string   `form:"email" json:"email"`
        MinAge   int      `form:"min_age" json:"min_age"`
        MaxAge   int      `form:"max_age" json:"max_age"`
        Status   []string `form:"status" json:"status"`
        Roles    []string `form:"roles" json:"roles"`
    }
    
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 设置默认值
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    
    session := pie.Table[User](engine)
    query := session
    
    // 添加查询条件
    if req.Name != "" {
        query = query.Where("name", pie.Like("name", "%"+req.Name+"%"))
    }
    if req.Email != "" {
        query = query.Where("email", pie.Like("email", "%"+req.Email+"%"))
    }
    if req.MinAge > 0 {
        query = query.Where("age", pie.Gte("age", req.MinAge))
    }
    if req.MaxAge > 0 {
        query = query.Where("age", pie.Lte("age", req.MaxAge))
    }
    if len(req.Status) > 0 {
        query = query.WhereIn("status", req.Status)
    }
    if len(req.Roles) > 0 {
        query = query.WhereIn("role", req.Roles)
    }
    
    // 添加排序
    if req.SortBy != "" {
        if req.SortOrder == "desc" {
            query = query.OrderByDesc(req.SortBy)
        } else {
            query = query.OrderBy(req.SortBy)
        }
    } else {
        query = query.OrderByDesc("created_at")
    }
    
    // 执行分页查询
    params := pie.PaginateParams{
        Page:     req.Page,
        PageSize: req.PageSize,
    }
    
    result, err := query.Paginate(context.Background(), params)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, PaginationResponse[User]{
        Data:       result.Data,
        Total:      result.Total,
        TotalPages: result.TotalPages,
        Page:       result.Page,
        PageSize:   result.PageSize,
        HasNext:    result.HasNext,
        HasPrev:    result.HasPrev,
    })
}
```

## 高级用法

### 游标分页

对于大数据量场景，可以使用游标分页：

```go
// 游标分页参数
type CursorPaginationRequest struct {
    Cursor string `form:"cursor" json:"cursor"`
    Limit  int    `form:"limit" json:"limit" binding:"min=1,max=100"`
}

// 游标分页响应
type CursorPaginationResponse[T any] struct {
    Data   []T   `json:"data"`
    Cursor string `json:"cursor,omitempty"`
    HasMore bool `json:"has_more"`
}

func GetUsersWithCursor(c *gin.Context) {
    var req CursorPaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if req.Limit <= 0 {
        req.Limit = 20
    }
    
    session := pie.Table[User](engine)
    query := session.Where("status", "active").OrderByDesc("created_at")
    
    // 如果有游标，从游标位置开始
    if req.Cursor != "" {
        // 解析游标（通常是最后一条记录的 ID 或时间戳）
        cursorTime, err := time.Parse(time.RFC3339, req.Cursor)
        if err == nil {
            query = query.Where("created_at", pie.Lt("created_at", cursorTime))
        }
    }
    
    // 获取比限制多一条记录来判断是否还有更多
    var users []User
    err := query.Limit(req.Limit + 1).Find(context.Background(), &users)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    hasMore := len(users) > req.Limit
    if hasMore {
        users = users[:req.Limit] // 移除多余的一条记录
    }
    
    var nextCursor string
    if hasMore && len(users) > 0 {
        nextCursor = users[len(users)-1].CreatedAt.Format(time.RFC3339)
    }
    
    c.JSON(200, CursorPaginationResponse[User]{
        Data:    users,
        Cursor:  nextCursor,
        HasMore: hasMore,
    })
}
```

### 分页缓存

对于频繁访问的分页数据，可以使用缓存：

```go
func GetCachedUsers(c *gin.Context) {
    var req PaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 生成缓存键
    cacheKey := fmt.Sprintf("users:page:%d:size:%d:sort:%s:%s", 
        req.Page, req.PageSize, req.SortBy, req.SortOrder)
    
    // 尝试从缓存获取
    var cachedResult PaginationResponse[User]
    if err := cache.Get(cacheKey, &cachedResult); err == nil {
        c.JSON(200, cachedResult)
        return
    }
    
    // 缓存未命中，从数据库查询
    session := pie.Table[User](engine).WithCache(5 * time.Minute)
    query := session.Where("status", "active")
    
    if req.SortBy != "" {
        if req.SortOrder == "desc" {
            query = query.OrderByDesc(req.SortBy)
        } else {
            query = query.OrderBy(req.SortBy)
        }
    }
    
    params := pie.PaginateParams{
        Page:     req.Page,
        PageSize: req.PageSize,
    }
    
    result, err := query.Paginate(context.Background(), params)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    response := PaginationResponse[User]{
        Data:       result.Data,
        Total:      result.Total,
        TotalPages: result.TotalPages,
        Page:       result.Page,
        PageSize:   result.PageSize,
        HasNext:    result.HasNext,
        HasPrev:    result.HasPrev,
    }
    
    // 缓存结果
    cache.Set(cacheKey, response, 5*time.Minute)
    
    c.JSON(200, response)
}
```

## 性能优化

### 1. 使用简单分页

对于大数据量场景，使用简单分页避免统计总数：

```go
// 大数据量场景，使用简单分页
result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    PaginateSimple(ctx, pie.PaginateParams{
        Page:     page,
        PageSize: pageSize,
    })
```

### 2. 合理设置页面大小

```go
// 限制最大页面大小
if pageSize > 100 {
    pageSize = 100
}

// 根据场景设置合适的页面大小
const (
    DefaultPageSize = 10
    MaxPageSize     = 100
    LargePageSize   = 50  // 用于管理后台
)
```

### 3. 使用索引优化

```go
// 确保排序字段有索引
session.Where("status", "active").OrderBy("created_at") // created_at 应该有索引

// 复合索引优化
// 对于 WHERE status = ? ORDER BY created_at 的查询
// 建议创建复合索引: {status: 1, created_at: -1}
```

### 4. 避免深分页

```go
// 避免使用过大的页码
if page > 1000 {
    // 使用游标分页或其他策略
    return getUsersWithCursor(cursor, limit)
}
```

## 最佳实践

### 1. 统一分页参数处理

```go
func parsePaginationParams(c *gin.Context) (pie.PaginateParams, error) {
    var req PaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        return pie.PaginateParams{}, err
    }
    
    // 设置默认值和限制
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    if req.PageSize > 100 {
        req.PageSize = 100
    }
    
    return pie.PaginateParams{
        Page:     req.Page,
        PageSize: req.PageSize,
    }, nil
}
```

### 2. 分页元数据

```go
type PaginationMeta struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    Total      int64 `json:"total,omitempty"`
    TotalPages int   `json:"total_pages,omitempty"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
}

type PaginatedResponse[T any] struct {
    Data []T            `json:"data"`
    Meta PaginationMeta `json:"meta"`
}
```

### 3. 错误处理

```go
func handlePaginationError(err error) gin.H {
    if pie.IsNotFoundError(err) {
        return gin.H{"error": "No data found"}
    }
    if pie.IsTimeoutError(err) {
        return gin.H{"error": "Query timeout"}
    }
    return gin.H{"error": "Internal server error"}
}
```

## 下一步

- [游标操作](/core-features/cursor/) - 学习高效的数据遍历
- [聚合查询](/core-features/aggregation/) - 掌握复杂数据聚合
- [缓存支持](/advanced/cache/) - 提升查询性能
