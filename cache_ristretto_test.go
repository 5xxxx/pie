package pie

import (
	"context"
	"testing"
	"time"
)

func TestNewRistrettoCache(t *testing.T) {
	// 测试默认配置
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer cache.Close()

	if cache == nil {
		t.Error("Expected cache to be created")
	}
}

func TestRistrettoCacheBasicOperations(t *testing.T) {
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	ctx := context.Background()
	key := "test-key"
	value := []byte("test-value")

	// Test Set
	err = cache.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Test Get
	retrieved, err := cache.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}

	// Test Exists
	exists, err := cache.Exists(ctx, key)
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// Test Delete
	err = cache.Delete(ctx, key)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Test Get after delete
	_, err = cache.Get(ctx, key)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
}

func TestRistrettoCacheWithTags(t *testing.T) {
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	ctx := context.Background()
	key := "test-key"
	value := []byte("test-value")
	tags := []string{"tag1", "tag2"}

	// Test SetWithTags
	err = cache.SetWithTags(ctx, key, value, 5*time.Minute, tags)
	if err != nil {
		t.Errorf("SetWithTags failed: %v", err)
	}

	// Test Get
	retrieved, err := cache.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}

	// Test DeleteByTags
	err = cache.DeleteByTags(ctx, "tag1")
	if err != nil {
		t.Errorf("DeleteByTags failed: %v", err)
	}

	// Test Get after delete by tags
	_, err = cache.Get(ctx, key)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
}

func TestRistrettoCachePattern(t *testing.T) {
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Set multiple keys
	keys := []string{"user:1", "user:2", "product:1", "product:2"}
	for _, key := range keys {
		value := []byte("value-" + key)
		err = cache.Set(ctx, key, value, 5*time.Minute)
		if err != nil {
			t.Errorf("Set failed for key %s: %v", key, err)
		}
		// 等待一小段时间确保缓存设置完成
		time.Sleep(time.Millisecond * 10)
	}

	// Test DeleteByPattern
	err = cache.DeleteByPattern(ctx, "user:.*")
	if err != nil {
		t.Errorf("DeleteByPattern failed: %v", err)
	}

	// Check that user keys are deleted
	for _, key := range []string{"user:1", "user:2"} {
		_, err = cache.Get(ctx, key)
		if err != ErrCacheNotFound {
			t.Errorf("Expected ErrCacheNotFound for key %s, got %v", key, err)
		}
	}

	// Check that product keys still exist
	for _, key := range []string{"product:1", "product:2"} {
		_, err = cache.Get(ctx, key)
		if err != nil {
			t.Errorf("Expected product key %s to exist, got error: %v", key, err)
		}
	}
}

func TestRistrettoCacheClear(t *testing.T) {
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Set some data
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err = cache.Set(ctx, key, []byte("value"), 5*time.Minute)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}
	}

	// Clear cache
	err = cache.Clear(ctx)
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	// Check that all keys are gone
	for _, key := range keys {
		_, err = cache.Get(ctx, key)
		if err != ErrCacheNotFound {
			t.Errorf("Expected ErrCacheNotFound for key %s, got %v", key, err)
		}
	}
}

func TestRistrettoCacheStats(t *testing.T) {
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Initial stats
	stats := cache.Stats()
	if stats.Total != 0 {
		t.Errorf("Expected initial Total to be 0, got %d", stats.Total)
	}

	// Set and get some data
	key := "test-key"
	value := []byte("test-value")

	err = cache.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	_, err = cache.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	// Check stats
	stats = cache.Stats()
	if stats.Total != 1 {
		t.Errorf("Expected Total to be 1, got %d", stats.Total)
	}
	if stats.Hits != 1 {
		t.Errorf("Expected Hits to be 1, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected Misses to be 0, got %d", stats.Misses)
	}
	if stats.HitRate != 100.0 {
		t.Errorf("Expected HitRate to be 100.0, got %f", stats.HitRate)
	}
}

func TestRistrettoCacheTTL(t *testing.T) {
	cache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	ctx := context.Background()
	key := "test-key"
	value := []byte("test-value")

	// Set with short TTL
	err = cache.Set(ctx, key, value, 100*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Should exist immediately
	_, err = cache.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed immediately: %v", err)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, err = cache.Get(ctx, key)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound after TTL, got %v", err)
	}
}
