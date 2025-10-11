package pie

import (
	"context"
	"testing"
	"time"
)

func TestCacheManagerChain(t *testing.T) {
	// 创建两个缓存实例
	ristrettoCache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer ristrettoCache.Close()

	mockRedisCache := NewMockRedisCache()
	defer mockRedisCache.Close()

	// 创建链式缓存管理器
	manager := NewCacheManager([]Cache{ristrettoCache, mockRedisCache}, nil)
	ctx := context.Background()

	key := "test-key"
	value := []byte("test-value")

	// Test Set - 应该写入所有缓存
	err = manager.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Test Get - 应该从第一个缓存获取
	retrieved, err := manager.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}

	// 验证两个缓存都有数据 (使用带前缀的键)
	fullKey := "pie:" + key
	ristrettoValue, err := ristrettoCache.Get(ctx, fullKey)
	if err != nil {
		t.Errorf("Ristretto cache should have the value: %v", err)
	}
	if string(ristrettoValue) != string(value) {
		t.Errorf("Ristretto cache value mismatch")
	}

	mockRedisValue, err := mockRedisCache.Get(ctx, fullKey)
	if err != nil {
		t.Errorf("Mock Redis cache should have the value: %v", err)
	}
	if string(mockRedisValue) != string(value) {
		t.Errorf("Mock Redis cache value mismatch")
	}
}

func TestCacheManagerChainBackfill(t *testing.T) {
	// 创建两个缓存实例
	ristrettoCache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer ristrettoCache.Close()

	mockRedisCache := NewMockRedisCache()
	defer mockRedisCache.Close()

	// 创建链式缓存管理器
	manager := NewCacheManager([]Cache{ristrettoCache, mockRedisCache}, nil)
	ctx := context.Background()

	key := "test-key"
	value := []byte("test-value")

	// 只在第二个缓存中设置数据 (使用带前缀的键)
	fullKey := "pie:" + key
	err = mockRedisCache.Set(ctx, fullKey, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set to mock Redis failed: %v", err)
	}

	// 从链式缓存获取 - 应该从第二个缓存获取并回填到第一个
	retrieved, err := manager.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}

	// 验证第一个缓存现在也有数据（回填）
	ristrettoValue, err := ristrettoCache.Get(ctx, fullKey)
	if err != nil {
		t.Errorf("Ristretto cache should have been backfilled: %v", err)
	}
	if string(ristrettoValue) != string(value) {
		t.Errorf("Ristretto cache backfill value mismatch")
	}
}

func TestCacheManagerChainDelete(t *testing.T) {
	// 创建两个缓存实例
	ristrettoCache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer ristrettoCache.Close()

	mockRedisCache := NewMockRedisCache()
	defer mockRedisCache.Close()

	// 创建链式缓存管理器
	manager := NewCacheManager([]Cache{ristrettoCache, mockRedisCache}, nil)
	ctx := context.Background()

	key := "test-key"
	value := []byte("test-value")

	// 设置数据到两个缓存
	err = manager.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// 删除 - 应该从所有缓存删除
	err = manager.Delete(ctx, key)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// 验证两个缓存都没有数据
	_, err = ristrettoCache.Get(ctx, key)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound from Ristretto cache, got %v", err)
	}

	_, err = mockRedisCache.Get(ctx, key)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound from Mock Redis cache, got %v", err)
	}
}

func TestCacheManagerChainPattern(t *testing.T) {
	// 创建两个缓存实例
	ristrettoCache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer ristrettoCache.Close()

	mockRedisCache := NewMockRedisCache()
	defer mockRedisCache.Close()

	// 创建链式缓存管理器
	manager := NewCacheManager([]Cache{ristrettoCache, mockRedisCache}, nil)
	ctx := context.Background()

	// 设置多个键
	keys := []string{"user:1", "user:2", "product:1", "product:2"}
	for _, key := range keys {
		value := []byte("value-" + key)
		err = manager.Set(ctx, key, value, 5*time.Minute)
		if err != nil {
			t.Errorf("Set failed for key %s: %v", key, err)
		}
	}

	// 按模式删除
	err = manager.DeleteByPattern(ctx, "user:.*")
	if err != nil {
		t.Errorf("DeleteByPattern failed: %v", err)
	}

	// 验证 user 键被删除
	for _, key := range []string{"user:1", "user:2"} {
		_, err = manager.Get(ctx, key)
		if err != ErrCacheNotFound {
			t.Errorf("Expected ErrCacheNotFound for key %s, got %v", key, err)
		}
	}

	// 验证 product 键仍然存在
	for _, key := range []string{"product:1", "product:2"} {
		_, err = manager.Get(ctx, key)
		if err != nil {
			t.Errorf("Expected product key %s to exist, got error: %v", key, err)
		}
	}
}

func TestCacheManagerChainTags(t *testing.T) {
	// 创建两个缓存实例
	ristrettoCache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer ristrettoCache.Close()

	mockRedisCache := NewMockRedisCache()
	defer mockRedisCache.Close()

	// 创建链式缓存管理器
	manager := NewCacheManager([]Cache{ristrettoCache, mockRedisCache}, nil)
	ctx := context.Background()

	// 设置带标签的数据
	key1 := "key1"
	key2 := "key2"
	value1 := []byte("value1")
	value2 := []byte("value2")

	err = ristrettoCache.SetWithTags(ctx, key1, value1, 5*time.Minute, []string{"tag1", "tag2"})
	if err != nil {
		t.Errorf("SetWithTags failed: %v", err)
	}

	err = mockRedisCache.SetWithTags(ctx, key2, value2, 5*time.Minute, []string{"tag1", "tag3"})
	if err != nil {
		t.Errorf("SetWithTags failed: %v", err)
	}

	// 按标签删除
	err = manager.DeleteByTags(ctx, "tag1")
	if err != nil {
		t.Errorf("DeleteByTags failed: %v", err)
	}

	// 验证两个键都被删除
	_, err = manager.Get(ctx, key1)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound for key1, got %v", err)
	}

	_, err = manager.Get(ctx, key2)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound for key2, got %v", err)
	}
}

func TestCacheManagerChainStats(t *testing.T) {
	// 创建两个缓存实例
	ristrettoCache, err := NewRistrettoCache(nil)
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer ristrettoCache.Close()

	mockRedisCache := NewMockRedisCache()
	defer mockRedisCache.Close()

	// 创建链式缓存管理器
	manager := NewCacheManager([]Cache{ristrettoCache, mockRedisCache}, nil)
	ctx := context.Background()

	// 执行一些操作
	key := "test-key"
	value := []byte("test-value")

	// Set
	err = manager.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Get (hit)
	_, err = manager.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	// Get non-existent key (miss)
	_, err = manager.Get(ctx, "non-existent")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}

	// 检查统计信息 (链式缓存的统计是聚合的，所以数值会不同)
	stats := manager.Stats()
	if stats.Total < 2 {
		t.Errorf("Expected Total to be at least 2, got %d", stats.Total)
	}
	if stats.Hits < 1 {
		t.Errorf("Expected Hits to be at least 1, got %d", stats.Hits)
	}
	if stats.Misses < 1 {
		t.Errorf("Expected Misses to be at least 1, got %d", stats.Misses)
	}
	// HitRate 计算可能因为聚合而不同，这里只检查是否在合理范围内
	if stats.HitRate < 0 || stats.HitRate > 100 {
		t.Errorf("Expected HitRate to be between 0 and 100, got %f", stats.HitRate)
	}
}

func TestCacheManagerChainEmpty(t *testing.T) {
	// 测试空缓存链
	manager := NewCacheManager([]Cache{}, nil)
	ctx := context.Background()

	key := "test-key"
	value := []byte("test-value")

	// Set 应该不报错但也不做任何事
	err := manager.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set on empty cache chain should not error: %v", err)
	}

	// Get 应该返回 ErrCacheDisabled
	_, err = manager.Get(ctx, key)
	if err != ErrCacheDisabled {
		t.Errorf("Expected ErrCacheDisabled, got %v", err)
	}

	// Stats 应该返回空统计
	stats := manager.Stats()
	if stats.Total != 0 {
		t.Errorf("Expected empty stats, got %+v", stats)
	}
}
