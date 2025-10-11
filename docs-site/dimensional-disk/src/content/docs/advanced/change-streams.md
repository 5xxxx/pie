---
title: Change Streams
description: Real-time data change monitoring with event-driven architecture support
---

# Change Streams

Pie provides change stream functionality for real-time data monitoring.

## Basic Usage

```go
// Create change stream watcher
watcher := pie.NewWatcher[User](engine)

// Watch collection changes
err := watcher.
    WatchCollection().
    Start(ctx, func(change *pie.ChangeEvent[User]) error {
        log.Printf("Change type: %s", change.OperationType)
        log.Printf("Document ID: %s", change.DocumentKey)
        
        switch change.OperationType {
        case "insert":
            log.Printf("New user inserted: %s", change.FullDocument.Name)
        case "update":
            log.Printf("User updated: %s", change.FullDocument.Name)
        case "delete":
            log.Printf("User deleted: %s", change.DocumentKey)
        }
        
        return nil
    })
```

## Next Steps

- [Query Scopes](/advanced/scopes/) - Learn about query scopes
- [Logging & Monitoring](/advanced/logging/) - Learn about logging
- [Best Practices](/best-practices/) - Learn development best practices
