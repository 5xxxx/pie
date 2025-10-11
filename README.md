# Pie - MongoDB ORM Framework

Pie is a modern, type-safe MongoDB ORM framework for Go that provides a rich, high-performance database operation experience.

## Features

- **Type Safety**: Built on Go generics with compile-time type checking
- **Smart Query Builder**: Intuitive chainable API with support for complex query conditions
- **Struct Query**: Direct conversion of HTTP request parameters to query conditions
- **Cache Support**: Built-in multi-level caching with Redis and memory cache support
- **Hook System**: Complete lifecycle hook support
- **Transaction Management**: Simple and easy-to-use transaction operations
- **Index Management**: Automated index creation and management
- **Advanced Aggregation**: Comprehensive stage builders with 100+ expression functions
- **Change Streams**: Real-time data change monitoring
- **Soft Delete**: Built-in soft delete functionality
- **Pagination**: Efficient pagination implementation
- **Query Logging**: Detailed query logging and performance monitoring
- **Bulk Operations**: High-performance bulk write operations

## Installation

```bash
go get github.com/5xxxx/pie
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "github.com/5xxxx/pie"
)

func main() {
    // Create engine
    engine, err := pie.NewEngine(
        context.Background(),
        "mydb",
        pie.WithURI("mongodb://localhost:27017"),
    )
    if err != nil {
        log.Fatal("Failed to create engine:", err)
    }
    defer engine.Disconnect(context.Background())
    
    // Create type-safe session
    session := pie.Table[User](engine)
    
    // Insert document
    user := &User{Name: "John Doe", Email: "john@example.com"}
    result, err := session.Insert(context.Background(), user)
    
    // Query documents
    var users []User
    err = session.Where("age", pie.Gte("age", 18)).Find(context.Background(), &users)
}
```

## Documentation

Complete documentation is available at: [https://5xxxx.github.io/pie](https://5xxxx.github.io/pie)

## License

MIT License

## Contributing

Issues and Pull Requests are welcome!