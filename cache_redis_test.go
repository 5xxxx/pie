package pie

import (
        "context"
        "testing"
        "time"
)

func TestRedisCacheSimulatedBehavior(t *testing.T) {
        cache := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
        ctx := context.Background()

        if _, err := cache.Get(ctx, "key"); err == nil {
                t.Fatalf("simulated cache should report miss")
        }

        if err := cache.Set(ctx, "key", []byte("value"), time.Second); err != nil {
                t.Fatalf("Set error: %v", err)
        }
        if err := cache.SetWithTags(ctx, "key", []byte("value"), time.Second, []string{"tag"}); err != nil {
                t.Fatalf("SetWithTags error: %v", err)
        }

        if err := cache.Delete(ctx, "key"); err != nil {
                t.Fatalf("Delete error: %v", err)
        }
        if err := cache.DeleteByPattern(ctx, "pattern"); err != nil {
                t.Fatalf("DeleteByPattern error: %v", err)
        }
        if err := cache.DeleteByTags(ctx, "tag"); err != nil {
                t.Fatalf("DeleteByTags error: %v", err)
        }

        if exists, err := cache.Exists(ctx, "key"); err != nil || exists {
                t.Fatalf("Exists returned (%v, %v)", err, exists)
        }
        if err := cache.Clear(ctx); err != nil {
                t.Fatalf("Clear error: %v", err)
        }

        if cache.Stats() == nil {
                t.Fatalf("Stats returned nil")
        }

        if err := cache.Close(); err != nil {
                t.Fatalf("Close error: %v", err)
        }

        if err := cache.Ping(ctx); err == nil {
                t.Fatalf("Ping should report unavailable simulation")
        }
}

