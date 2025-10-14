---
title: Pagination
description: Pie provides efficient pagination implementation
---

# Pagination

```go
// Full pagination (with total count)
result, err := session.
    Where("status", "active").
    Paginate(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
// result.Total, result.TotalPages, result.HasNext, result.HasPrev

// Simple pagination (no total count, faster)
simpleResult, err := session.
    PaginateSimple(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })
```

## Pagination Types

### Full Pagination

```go
// Full pagination with total count
result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    Paginate(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })

if err != nil {
    return err
}

// Access pagination metadata
fmt.Printf("Total items: %d\n", result.Total)
fmt.Printf("Total pages: %d\n", result.TotalPages)
fmt.Printf("Current page: %d\n", result.Page)
fmt.Printf("Page size: %d\n", result.PageSize)
fmt.Printf("Has next: %t\n", result.HasNext)
fmt.Printf("Has prev: %t\n", result.HasPrev)

// Access data
for _, user := range result.Data {
    fmt.Printf("User: %s\n", user.Name)
}
```

### Simple Pagination

```go
// Simple pagination without total count (faster)
result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    PaginateSimple(ctx, pie.PaginateParams{
        Page:     1,
        PageSize: 10,
    })

if err != nil {
    return err
}

// Access data
for _, user := range result.Data {
    fmt.Printf("User: %s\n", user.Name)
}
```

## Advanced Pagination

### Cursor-based Pagination

```go
// Cursor-based pagination for large datasets
result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    PaginateCursor(ctx, pie.CursorPaginateParams{
        Cursor:   lastCursor, // From previous page
        PageSize: 10,
    })

if err != nil {
    return err
}

// Get next cursor for next page
nextCursor := result.NextCursor
```

### Offset-style Pagination

You can express a traditional offset/limit request by mapping it to the existing
`Paginate` (or `PaginateSimple`) helpers. For example, to skip the first 20
records and return 10 results per page you would set `Page` to
`offset/pageSize + 1` and `PageSize` to the desired limit:

```go
offset := 20
pageSize := 10

result, err := session.
    Where("status", "active").
    OrderBy("created_at").
    Paginate(ctx, pie.PaginateParams{
        Page:     offset/pageSize + 1,
        PageSize: pageSize,
    })
```

If you do not need the total count you can call `PaginateSimple` with the same
parameters. For real cursor-based pagination support, prefer
`PaginateCursor` instead of relying on offsets.

## Real-world Examples

### API Pagination Handler

```go
type PaginationRequest struct {
    Page     int    `json:"page" validate:"min=1"`
    PageSize int    `json:"page_size" validate:"min=1,max=100"`
    SortBy   string `json:"sort_by"`
    SortDir  string `json:"sort_dir" validate:"oneof=asc desc"`
}

type PaginationResponse struct {
    Data       []User `json:"data"`
    Total      int64  `json:"total"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
    TotalPages int    `json:"total_pages"`
    HasNext    bool   `json:"has_next"`
    HasPrev    bool   `json:"has_prev"`
}

func GetUsers(ctx context.Context, req PaginationRequest) (*PaginationResponse, error) {
    session := pie.Table[User](engine)
    
    // Build query
    query := session.Where("status", "active")
    
    // Apply sorting
    if req.SortBy != "" {
        if req.SortDir == "desc" {
            query = query.OrderByDesc(req.SortBy)
        } else {
            query = query.OrderBy(req.SortBy)
        }
    }
    
    // Apply pagination
    result, err := query.Paginate(ctx, pie.PaginateParams{
        Page:     req.Page,
        PageSize: req.PageSize,
    })
    
    if err != nil {
        return nil, err
    }
    
    return &PaginationResponse{
        Data:       result.Data,
        Total:      result.Total,
        Page:       result.Page,
        PageSize:   result.PageSize,
        TotalPages: result.TotalPages,
        HasNext:    result.HasNext,
        HasPrev:    result.HasPrev,
    }, nil
}
```

### Search with Pagination

```go
type SearchRequest struct {
    Query    string `json:"query"`
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    SortBy   string `json:"sort_by"`
}

func SearchUsers(ctx context.Context, req SearchRequest) (*PaginationResponse, error) {
    session := pie.Table[User](engine)
    
    // Build search query
    query := session.Where("status", "active")
    
    if req.Query != "" {
        query = query.Where("name", pie.Like("name", "%"+req.Query+"%"))
    }
    
    // Apply sorting
    if req.SortBy != "" {
        query = query.OrderBy(req.SortBy)
    }
    
    // Apply pagination
    result, err := query.Paginate(ctx, pie.PaginateParams{
        Page:     req.Page,
        PageSize: req.PageSize,
    })
    
    if err != nil {
        return nil, err
    }
    
    return &PaginationResponse{
        Data:       result.Data,
        Total:      result.Total,
        Page:       result.Page,
        PageSize:   result.PageSize,
        TotalPages: result.TotalPages,
        HasNext:    result.HasNext,
        HasPrev:    result.HasPrev,
    }, nil
}
```

## Performance Optimization

### Use Indexes

```go
// Ensure pagination fields are indexed
func createPaginationIndexes() error {
    indexes := pie.MustIndexes(engine)
    
    // Index for sorting
    err := indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"created_at", -1},
    })
    if err != nil {
        return err
    }
    
    // Index for search
    err = indexes.CreateIndex(ctx, "users", bson.D{
        {"status", 1},
        {"name", "text"},
    })
    
    return err
}
```

### Caching

```go
// Cache pagination results
func GetUsersCached(ctx context.Context, page, pageSize int) (*PaginationResponse, error) {
    cacheKey := fmt.Sprintf("users_page_%d_size_%d", page, pageSize)
    
    // Try cache first
    var cached PaginationResponse
    if err := cache.Get(cacheKey, &cached); err == nil {
        return &cached, nil
    }
    
    // Query database
    result, err := getUsersFromDB(ctx, page, pageSize)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    cache.Set(cacheKey, result, 5*time.Minute)
    
    return result, nil
}
```

## Best Practices

### 1. Set Reasonable Page Sizes

```go
// Good: Reasonable page size
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
)

func validatePageSize(pageSize int) int {
    if pageSize <= 0 {
        return DefaultPageSize
    }
    if pageSize > MaxPageSize {
        return MaxPageSize
    }
    return pageSize
}
```

### 2. Use Cursor Pagination for Large Datasets

```go
// For large datasets, use cursor-based pagination
func GetLargeDataset(ctx context.Context, cursor string, pageSize int) (*CursorResult, error) {
    session := pie.Table[User](engine)
    
    result, err := session.
        Where("status", "active").
        OrderBy("created_at").
        PaginateCursor(ctx, pie.CursorPaginateParams{
            Cursor:   cursor,
            PageSize: pageSize,
        })
    
    return result, err
}
```

### 3. Handle Edge Cases

```go
func GetUsersSafe(ctx context.Context, page, pageSize int) (*PaginationResponse, error) {
    // Validate input
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    
    session := pie.Table[User](engine)
    
    result, err := session.
        Where("status", "active").
        Paginate(ctx, pie.PaginateParams{
            Page:     page,
            PageSize: pageSize,
        })
    
    if err != nil {
        return nil, err
    }
    
    // Handle empty results
    if len(result.Data) == 0 {
        return &PaginationResponse{
            Data:       []User{},
            Total:      0,
            Page:       page,
            PageSize:   pageSize,
            TotalPages: 0,
            HasNext:    false,
            HasPrev:    false,
        }, nil
    }
    
    return &PaginationResponse{
        Data:       result.Data,
        Total:      result.Total,
        Page:       result.Page,
        PageSize:   result.PageSize,
        TotalPages: result.TotalPages,
        HasNext:    result.HasNext,
        HasPrev:    result.HasPrev,
    }, nil
}
```

## Next Steps

- [Cursor Operations](/core-features/cursor/) - Use cursors for large datasets
- [Query Builder](/core-features/query-builder/) - Learn about query methods
- [Transactions](/core-features/transactions/) - Use transactions
- [Performance](/reference/performance/) - Learn performance optimization
