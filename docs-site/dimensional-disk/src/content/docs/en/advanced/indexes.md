---
title: Index Management
description: Automated index creation and management with struct tags
---

# Index Management

Pie provides automated index creation and management functionality.

## Basic Usage

```go
// Create index manager
indexes := pie.MustIndexes(engine)

// Create indexes for struct
err := indexes.CreateIndexes(ctx, User{})

// Manual index creation
err := indexes.CreateIndex(ctx, "users", bson.D{
    {"email", 1},
}, &options.IndexOptions{
    Unique: pie.Bool(true),
})
```

## Struct Tags

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    Name      string        `bson:"name" pie:"index"`
    Email     string        `bson:"email" pie:"unique"`
    Age       int           `bson:"age" pie:"index,sparse"`
    CreatedAt time.Time     `bson:"created_at" pie:"index"`
}
```

## Next Steps

- [Change Streams](/advanced/change-streams/) - Learn about change streams
- [Query Scopes](/advanced/scopes/) - Learn about query scopes
- [Performance](/reference/performance/) - Learn performance optimization
