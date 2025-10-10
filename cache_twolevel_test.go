package pie

import (
        "context"
        "errors"
        "testing"
        "time"
)

type recordingCache struct {
        *mockCache
        name string
}

func newRecordingCache(name string) *recordingCache {
        return &recordingCache{mockCache: newMockCache(), name: name}
}

func (r *recordingCache) Get(ctx context.Context, key string) ([]byte, error) {
        val, err := r.mockCache.Get(ctx, r.name+":"+key)
        if err != nil {
                return nil, err
        }
        return val, nil
}

func (r *recordingCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
        return r.mockCache.Set(ctx, r.name+":"+key, value, ttl)
}

func (r *recordingCache) Delete(ctx context.Context, key string) error {
        return r.mockCache.Delete(ctx, r.name+":"+key)
}

func (r *recordingCache) DeleteByPattern(ctx context.Context, pattern string) error {
        return r.mockCache.DeleteByPattern(ctx, r.name+":"+pattern)
}

func (r *recordingCache) DeleteByTags(ctx context.Context, tags ...string) error {
        return r.mockCache.DeleteByTags(ctx, append([]string{r.name}, tags...)...)
}

func (r *recordingCache) Exists(ctx context.Context, key string) (bool, error) {
        return r.mockCache.Exists(ctx, r.name+":"+key)
}

func (r *recordingCache) Clear(ctx context.Context) error {
        return r.mockCache.Clear(ctx)
}

func (r *recordingCache) Stats() *CacheStats { return r.mockCache.Stats() }

func TestTwoLevelCache(t *testing.T) {
        ctx := context.Background()
        l1 := newRecordingCache("l1")
        l2 := newRecordingCache("l2")
        cache := NewTwoLevelCache(l1, l2, &TwoLevelCacheConfig{L1TTL: time.Second, L2TTL: 2 * time.Second, SyncOnWrite: true})

        if err := cache.Set(ctx, "key", []byte("value"), time.Minute); err != nil {
                t.Fatalf("Set error: %v", err)
        }

        l1Key := "l1:key"
        if _, ok := l1.values[l1Key]; !ok {
                t.Fatalf("expected value written to L1")
        }
        if _, ok := l2.values["l2:key"]; !ok {
                t.Fatalf("expected value written to L2")
        }

        l1.Delete(ctx, "key")
        val, err := cache.Get(ctx, "key")
        if err != nil || string(val) != "value" {
                t.Fatalf("Get returned (%v, %s)", err, val)
        }
        if _, ok := l1.values[l1Key]; !ok {
                t.Fatalf("expected value to be written back to L1 after L2 hit")
        }

        l1.Delete(ctx, "key")
        l2.Delete(ctx, "key")
        if _, err := cache.Get(ctx, "key"); !errors.Is(err, ErrCacheNotFound) {
                t.Fatalf("expected ErrCacheNotFound, got %v", err)
        }

        if err := cache.Delete(ctx, "x"); err != nil {
                t.Fatalf("Delete error: %v", err)
        }
        if err := cache.DeleteByPattern(ctx, "p"); err != nil {
                t.Fatalf("DeleteByPattern error: %v", err)
        }
        if err := cache.DeleteByTags(ctx, "t"); err != nil {
                t.Fatalf("DeleteByTags error: %v", err)
        }

        if exists, err := cache.Exists(ctx, "key"); err != nil || exists {
                t.Fatalf("Exists returned (%v, %v)", err, exists)
        }

        if err := cache.Clear(ctx); err != nil {
                t.Fatalf("Clear error: %v", err)
        }

        stats := cache.Stats()
        if stats.Total == 0 {
                t.Fatalf("expected stats to track operations")
        }

        detailed := cache.StatsDetailed()
        if detailed.Total == 0 {
                t.Fatalf("expected detailed stats to be populated")
        }

        if err := cache.SetL1Only(ctx, "only1", []byte("1")); err != nil {
                t.Fatalf("SetL1Only error: %v", err)
        }
        if err := cache.SetL2Only(ctx, "only2", []byte("2")); err != nil {
                t.Fatalf("SetL2Only error: %v", err)
        }

        if _, err := cache.GetL1(ctx, "only1"); err != nil {
                t.Fatalf("GetL1 error: %v", err)
        }
        if _, err := cache.GetL2(ctx, "only2"); err != nil {
                t.Fatalf("GetL2 error: %v", err)
        }
}

