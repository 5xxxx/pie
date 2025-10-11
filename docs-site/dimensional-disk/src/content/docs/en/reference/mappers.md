---
title: Name Mappers
description: Flexible field and table name mapping strategies
---

# Name Mappers

Pie supports multiple name mapping strategies to adapt to different Go struct field and MongoDB field naming conventions.

## Default Mapping

By default, Pie converts Go struct field names to camelCase as MongoDB field names.

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    FirstName string        `bson:"first_name"` // Explicit specification
    LastName  string        `bson:"last_name"`  // Explicit specification
    Email     string        // Default mapping to "email"
    CreatedAt time.Time     // Default mapping to "createdAt"
}
```

## Snake Case Mapping

If you want to convert Go struct field names to snake_case, use `SnakeMapper`.

```go
engine, err := pie.NewEngine(
    context.Background(),
    "mydb",
    pie.WithURI("mongodb://localhost:27017"),
    pie.WithMapper(&pie.SnakeMapper{}), // Use SnakeMapper
)

type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    FirstName string        // Maps to "first_name"
    LastName  string        // Maps to "last_name"
    Email     string        // Maps to "email"
    CreatedAt time.Time     // Maps to "created_at"
}
```

## Custom Mapping

You can implement the `pie.Mapper` interface to create custom name mapping logic.

```go
type CustomMapper struct{}

func (m *CustomMapper) FieldToBSON(fieldName string) string {
    // Custom logic, e.g., convert to uppercase
    return strings.ToUpper(fieldName)
}

func (m *CustomMapper) BSONToField(bsonName string) string {
    // Custom logic, e.g., convert to lowercase
    return strings.ToLower(bsonName)
}

engine, err := pie.NewEngine(
    context.Background(),
    "mydb",
    pie.WithURI("mongodb://localhost:27017"),
    pie.WithMapper(&CustomMapper{}), // Use custom Mapper
)
```

## Next Steps

- [Error Handling](/reference/error-handling/) - Learn about error handling
- [Performance](/reference/performance/) - Learn performance optimization
