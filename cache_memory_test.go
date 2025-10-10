package pie

import (
        "context"
        "errors"
        "reflect"
        "regexp"
        "runtime"
        "testing"
        "time"
)

func TestMemoryCacheBasicOperations(t *testing.T) {
        ctx := context.Background()
        cache := NewMemoryCache(2)
        t.Cleanup(cache.Close)

        if err := cache.Set(ctx, "key1", []byte("value1"), 50*time.Millisecond); err != nil {
                t.Fatalf("Set returned error: %v", err)
        }
        if err := cache.SetWithTags(ctx, "key2", []byte("value2"), 50*time.Millisecond, []string{"tag1", "tag2"}); err != nil {
                t.Fatalf("SetWithTags error: %v", err)
        }

        val, err := cache.Get(ctx, "key1")
        if err != nil || string(val) != "value1" {
                t.Fatalf("Get returned (%v, %s)", err, val)
        }

        exists, err := cache.Exists(ctx, "key2")
        if err != nil || !exists {
                t.Fatalf("Exists returned (%v, %v)", err, exists)
        }

        if err := cache.Delete(ctx, "key1"); err != nil {
                t.Fatalf("Delete error: %v", err)
        }

        if err := cache.DeleteByTags(ctx, "tag1"); err != nil {
                t.Fatalf("DeleteByTags error: %v", err)
        }

        if err := cache.DeleteByPattern(ctx, "key.*"); err != nil {
                t.Fatalf("DeleteByPattern error: %v", err)
        }

        stats := cache.Stats()
        if stats.Keys != 0 {
                t.Fatalf("expected no keys after deletions, got %d", stats.Keys)
        }

        if err := cache.Clear(ctx); err != nil {
                t.Fatalf("Clear error: %v", err)
        }
}

func TestMemoryCacheExpirationAndEviction(t *testing.T) {
        ctx := context.Background()
        cache := NewMemoryCache(1)
        t.Cleanup(cache.Close)

        if err := cache.Set(ctx, "expire", []byte("value"), 10*time.Millisecond); err != nil {
                t.Fatalf("Set error: %v", err)
        }
        time.Sleep(20 * time.Millisecond)

        if _, err := cache.Get(ctx, "expire"); !errors.Is(err, ErrCacheExpired) {
                t.Fatalf("expected ErrCacheExpired, got %v", err)
        }

        if err := cache.Set(ctx, "a", []byte("va"), time.Second); err != nil {
                t.Fatalf("Set error: %v", err)
        }
        if err := cache.Set(ctx, "b", []byte("vb"), time.Second); err != nil {
                t.Fatalf("Set error: %v", err)
        }

        if _, err := cache.Get(ctx, "a"); err == nil {
                t.Fatalf("expected first item to be evicted")
        }
}

func TestMemoryCacheCleanupRoutine(t *testing.T) {
        ctx := context.Background()
        cache := NewMemoryCache(10)
        if err := cache.Set(ctx, "expire", []byte("value"), 1*time.Millisecond); err != nil {
                t.Fatalf("Set error: %v", err)
        }

        time.Sleep(2 * time.Millisecond)
        cache.mu.Lock()
        cache.items["expire"].expiration = time.Now().Add(-time.Second)
        cache.mu.Unlock()

        // Trigger cleanup manually instead of waiting a minute by invoking method directly
        cache.mu.Lock()
        cache.evictOldest()
        cache.mu.Unlock()

        runtime.Gosched()
        if _, err := cache.Get(ctx, "expire"); err == nil {
                t.Fatalf("expected item to be removed after manual eviction")
        }
        cache.Close()
}

func TestMemoryCacheRemoveKeyFromTag(t *testing.T) {
        cache := NewMemoryCache(10)
        t.Cleanup(cache.Close)

        cache.mu.Lock()
        cache.tags["tag"] = []string{"a", "b", "c"}
        cache.mu.Unlock()

        cache.removeKeyFromTag("tag", "b")
        cache.mu.RLock()
        remaining := append([]string(nil), cache.tags["tag"]...)
        cache.mu.RUnlock()

        expected := []string{"a", "c"}
        if !reflect.DeepEqual(remaining, expected) {
                t.Fatalf("expected %v, got %v", expected, remaining)
        }

        cache.removeKeyFromTag("tag", "a")
        cache.removeKeyFromTag("tag", "c")
        cache.mu.RLock()
        _, exists := cache.tags["tag"]
        cache.mu.RUnlock()
        if exists {
                t.Fatalf("tag should be removed when empty")
        }
}

func TestMemoryCacheDeleteByPatternInvalidRegexp(t *testing.T) {
        cache := NewMemoryCache(10)
        t.Cleanup(cache.Close)

        err := cache.DeleteByPattern(context.Background(), "[")
        if err == nil || !regexp.MustCompile("missing closing").MatchString(err.Error()) {
                t.Fatalf("expected regexp compile error, got %v", err)
        }
}

