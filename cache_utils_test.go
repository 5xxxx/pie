package pie

import (
        "context"
        "errors"
        "sync"
        "testing"
        "time"
)

type mockCache struct {
        mu              sync.Mutex
        values          map[string][]byte
        patternsDeleted []string
        tagsDeleted     [][]string
        cleared         bool
        ttl             map[string]time.Duration
}

func newMockCache() *mockCache {
        return &mockCache{
                values: make(map[string][]byte),
                ttl:    make(map[string]time.Duration),
        }
}

func (m *mockCache) Get(_ context.Context, key string) ([]byte, error) {
        m.mu.Lock()
        defer m.mu.Unlock()
        val, ok := m.values[key]
        if !ok {
                return nil, ErrCacheNotFound
        }
        return val, nil
}

func (m *mockCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.values[key] = append([]byte(nil), value...)
        m.ttl[key] = ttl
        return nil
}

func (m *mockCache) Delete(_ context.Context, key string) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        delete(m.values, key)
        delete(m.ttl, key)
        return nil
}

func (m *mockCache) DeleteByPattern(_ context.Context, pattern string) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.patternsDeleted = append(m.patternsDeleted, pattern)
        return nil
}

func (m *mockCache) DeleteByTags(_ context.Context, tags ...string) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        copied := append([]string(nil), tags...)
        m.tagsDeleted = append(m.tagsDeleted, copied)
        return nil
}

func (m *mockCache) Exists(_ context.Context, key string) (bool, error) {
        m.mu.Lock()
        defer m.mu.Unlock()
        _, ok := m.values[key]
        return ok, nil
}

func (m *mockCache) Clear(context.Context) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.cleared = true
        m.values = make(map[string][]byte)
        m.ttl = make(map[string]time.Duration)
        return nil
}

func (m *mockCache) Stats() *CacheStats {
        return &CacheStats{}
}

func TestRandomInt(t *testing.T) {
        if got := randomInt(5, 1); got != 5 {
                t.Fatalf("expected 5 when min >= max, got %d", got)
        }

        seen := make(map[int64]bool)
        for i := 0; i < 100; i++ {
                v := randomInt(0, 2)
                if v < 0 || v > 2 {
                        t.Fatalf("value %d out of range", v)
                }
                seen[v] = true
        }

        for i := int64(0); i <= 2; i++ {
                if !seen[i] {
                        t.Fatalf("value %d never produced", i)
                }
        }
}

func TestApplyJitter(t *testing.T) {
        base := 10 * time.Second
        if got := applyJitter(base, 0); got != base {
                t.Fatalf("expected ttl unchanged when jitter is zero")
        }

        min := base - 2*time.Second
        max := base + 2*time.Second
        for i := 0; i < 100; i++ {
                got := applyJitter(base, 2*time.Second)
                if got < min || got > max {
                        t.Fatalf("jittered ttl %v outside expected range [%v,%v]", got, min, max)
                }
        }

        if got := applyJitter(base, -base); got < base {
                t.Fatalf("negative jitter should not reduce ttl, got %v", got)
        }
}

func TestCacheManagerEnabled(t *testing.T) {
        mock := newMockCache()
        manager := NewCacheManager(mock, &CacheConfig{Enabled: true, KeyPrefix: "prefix:"})

        ctx := context.Background()
        if err := manager.Set(ctx, "key", []byte("value"), time.Minute); err != nil {
                t.Fatalf("Set returned error: %v", err)
        }

        if ttl := mock.ttl["prefix:key"]; ttl != time.Minute {
                t.Fatalf("unexpected ttl recorded: %v", ttl)
        }

        val, err := manager.Get(ctx, "key")
        if err != nil || string(val) != "value" {
                t.Fatalf("Get returned %v, %q", err, val)
        }

        exists, err := manager.Exists(ctx, "key")
        if err != nil || !exists {
                t.Fatalf("Exists returned %v, %v", err, exists)
        }

        if err := manager.Delete(ctx, "key"); err != nil {
                t.Fatalf("Delete error: %v", err)
        }

        if err := manager.DeleteByPattern(ctx, "pattern"); err != nil {
                t.Fatalf("DeleteByPattern error: %v", err)
        }
        if len(mock.patternsDeleted) != 1 || mock.patternsDeleted[0] != "prefix:pattern" {
                t.Fatalf("pattern deletion not recorded: %+v", mock.patternsDeleted)
        }

        if err := manager.DeleteByTags(ctx, "a", "b"); err != nil {
                t.Fatalf("DeleteByTags error: %v", err)
        }
        if len(mock.tagsDeleted) != 1 || len(mock.tagsDeleted[0]) != 2 {
                t.Fatalf("tags deletion not recorded: %+v", mock.tagsDeleted)
        }

        if err := manager.Clear(ctx); err != nil {
                t.Fatalf("Clear error: %v", err)
        }
        if !mock.cleared {
                t.Fatalf("expected mock cache to be cleared")
        }

        if err := manager.InvalidateTags(ctx, "x"); err != nil {
                t.Fatalf("InvalidateTags error: %v", err)
        }

        warmed := false
        if err := manager.Warm(ctx, func(context.Context) error {
                warmed = true
                return nil
        }); err != nil {
                t.Fatalf("Warm error: %v", err)
        }
        if !warmed {
                t.Fatalf("warm function not invoked")
        }

        if manager.Stats() == nil {
                t.Fatalf("Stats returned nil")
        }
}

func TestCacheManagerDisabled(t *testing.T) {
        mock := newMockCache()
        manager := NewCacheManager(mock, &CacheConfig{Enabled: false, KeyPrefix: "prefix:"})
        ctx := context.Background()

        if _, err := manager.Get(ctx, "key"); !errors.Is(err, ErrCacheDisabled) {
                t.Fatalf("expected ErrCacheDisabled, got %v", err)
        }

        if err := manager.Set(ctx, "key", []byte("v"), time.Second); err != nil {
                t.Fatalf("Set returned error: %v", err)
        }
        if len(mock.values) != 0 {
                t.Fatalf("cache should not receive writes when disabled")
        }

        if exists, err := manager.Exists(ctx, "key"); err != nil || exists {
                t.Fatalf("Exists should return false without touching cache, got %v, %v", exists, err)
        }

        if err := manager.Delete(ctx, "key"); err != nil {
                t.Fatalf("Delete error: %v", err)
        }
        if err := manager.DeleteByPattern(ctx, "p"); err != nil {
                t.Fatalf("DeleteByPattern error: %v", err)
        }
        if err := manager.DeleteByTags(ctx, "t"); err != nil {
                t.Fatalf("DeleteByTags error: %v", err)
        }
        if err := manager.Clear(ctx); err != nil {
                t.Fatalf("Clear error: %v", err)
        }
}
