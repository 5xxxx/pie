---
title: Query Scopes
description: Reusable query logic with scopes for better code organization
---

# Query Scopes

Pie provides query scopes for reusable query logic and better code organization.

## Basic Usage

```go
// Define scope
func ActiveScope(field string) pie.ScopeFunc {
    return func(q *pie.Query) {
        q.Where(field, "active")
    }
}

// Use scope
users, err := session.Scopes(ActiveScope("status")).Find(ctx)

// Multiple scopes
users, err = session.Scopes(
    ActiveScope("status"),
    RecentScope("created_at", 30),
).Find(ctx)
```

## Next Steps

- [Logging & Monitoring](/advanced/logging/) - Learn about logging
- [Advanced Aggregation](/advanced/aggregation-advanced/) - Learn advanced aggregation
- [Best Practices](/best-practices/) - Learn development best practices
