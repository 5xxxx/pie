package pie

import (
	"context"
	"testing"
	"time"
)

func newMockCache() *mockCache {
	return &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
}

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	if !config.Enabled {
		t.Error("Expected Enabled to be true")
	}
	if config.DefaultTTL != 5*time.Minute {
		t.Errorf("Expected DefaultTTL to be 5 minutes, got %v", config.DefaultTTL)
	}
	if config.KeyPrefix != "pie:" {
		t.Errorf("Expected KeyPrefix to be 'pie:', got %s", config.KeyPrefix)
	}
	if config.MaxSize != 10000 {
		t.Errorf("Expected MaxSize to be 10000, got %d", config.MaxSize)
	}
	if config.EnableJitter {
		t.Error("Expected EnableJitter to be false")
	}
	if config.TTLJitter != 0 {
		t.Errorf("Expected TTLJitter to be 0, got %v", config.TTLJitter)
	}
	if config.EmptyCacheTTL != 30*time.Second {
		t.Errorf("Expected EmptyCacheTTL to be 30 seconds, got %v", config.EmptyCacheTTL)
	}
}

func TestNewCacheManager(t *testing.T) {
	mockCache := newMockCache()

	// 测试使用默认配置
	manager := NewCacheManager(mockCache, nil)
	if manager == nil {
		t.Error("Expected manager to be created")
	}
	if manager.cache != mockCache {
		t.Error("Expected cache to be set correctly")
	}
	if manager.config == nil {
		t.Error("Expected config to be set")
	}

	// 测试使用自定义配置
	customConfig := &CacheConfig{
		Enabled:    false,
		DefaultTTL: 10 * time.Minute,
		KeyPrefix:  "test:",
	}
	manager2 := NewCacheManager(mockCache, customConfig)
	if manager2.config.Enabled {
		t.Error("Expected custom config to be used")
	}
	if manager2.config.DefaultTTL != 10*time.Minute {
		t.Error("Expected custom DefaultTTL to be used")
	}
	if manager2.config.KeyPrefix != "test:" {
		t.Error("Expected custom KeyPrefix to be used")
	}
}

func TestCacheManagerGet(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 测试缓存命中
	mockCache.data["pie:testkey"] = []byte("testvalue")
	val, err := manager.Get(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 测试缓存未命中
	val, err = manager.Get(ctx, "nonexistent")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value")
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	val, err = manager.Get(ctx, "testkey")
	if err != ErrCacheDisabled {
		t.Errorf("Expected ErrCacheDisabled, got %v", err)
	}
}

func TestCacheManagerSet(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 测试设置缓存
	err := manager.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	if val, exists := mockCache.data["pie:testkey"]; !exists {
		t.Error("Expected key to be set")
	} else if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	err = manager.Set(ctx, "testkey2", []byte("testvalue2"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error when cache disabled, got %v", err)
	}

	// 验证缓存未设置
	if _, exists := mockCache.data["pie:testkey2"]; exists {
		t.Error("Expected key not to be set when cache disabled")
	}
}

func TestCacheManagerDelete(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 设置测试数据
	mockCache.data["pie:testkey"] = []byte("testvalue")

	// 测试删除
	err := manager.Delete(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证已删除
	if _, exists := mockCache.data["pie:testkey"]; exists {
		t.Error("Expected key to be deleted")
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	err = manager.Delete(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error when cache disabled, got %v", err)
	}
}

func TestCacheManagerDeleteByPattern(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 设置测试数据
	mockCache.data["pie:testkey1"] = []byte("value1")
	mockCache.data["pie:testkey2"] = []byte("value2")
	mockCache.data["pie:otherkey"] = []byte("value3")

	// 测试按模式删除
	err := manager.DeleteByPattern(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证删除结果
	if _, exists := mockCache.data["pie:testkey1"]; exists {
		t.Error("Expected testkey1 to be deleted")
	}
	if _, exists := mockCache.data["pie:testkey2"]; exists {
		t.Error("Expected testkey2 to be deleted")
	}
	if _, exists := mockCache.data["pie:otherkey"]; !exists {
		t.Error("Expected otherkey to remain")
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	err = manager.DeleteByPattern(ctx, "pattern")
	if err != nil {
		t.Errorf("Expected no error when cache disabled, got %v", err)
	}
}

func TestCacheManagerDeleteByTags(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 测试按标签删除
	err := manager.DeleteByTags(ctx, "tag1", "tag2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	err = manager.DeleteByTags(ctx, "tag1")
	if err != nil {
		t.Errorf("Expected no error when cache disabled, got %v", err)
	}
}

func TestCacheManagerExists(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 设置测试数据
	mockCache.data["pie:testkey"] = []byte("testvalue")

	// 测试存在
	exists, err := manager.Exists(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// 测试不存在
	exists, err = manager.Exists(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected key not to exist")
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	exists, err = manager.Exists(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected false when cache disabled")
	}
}

func TestCacheManagerClear(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 设置测试数据
	mockCache.data["pie:testkey1"] = []byte("value1")
	mockCache.data["pie:testkey2"] = []byte("value2")

	// 测试清空
	err := manager.Clear(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证已清空
	if len(mockCache.data) != 0 {
		t.Error("Expected cache to be cleared")
	}

	// 测试禁用缓存
	manager.config.Enabled = false
	err = manager.Clear(ctx)
	if err != nil {
		t.Errorf("Expected no error when cache disabled, got %v", err)
	}
}

func TestCacheManagerStats(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)

	// 设置一些统计数据
	mockCache.stats.Hits = 10
	mockCache.stats.Misses = 5
	mockCache.stats.Total = 15
	mockCache.stats.HitRate = 66.67

	stats := manager.Stats()
	if stats.Hits != 10 {
		t.Errorf("Expected Hits to be 10, got %d", stats.Hits)
	}
	if stats.Misses != 5 {
		t.Errorf("Expected Misses to be 5, got %d", stats.Misses)
	}
	if stats.Total != 15 {
		t.Errorf("Expected Total to be 15, got %d", stats.Total)
	}
	if stats.HitRate != 66.67 {
		t.Errorf("Expected HitRate to be 66.67, got %f", stats.HitRate)
	}
}

func TestCacheManagerInvalidateTags(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 测试标签失效
	err := manager.InvalidateTags(ctx, "tag1", "tag2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCacheManagerWarm(t *testing.T) {
	mockCache := newMockCache()
	manager := NewCacheManager(mockCache, nil)
	ctx := context.Background()

	// 测试预热
	warmFunc := func(ctx context.Context) error {
		return nil
	}

	err := manager.Warm(ctx, warmFunc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试预热函数返回错误
	errorWarmFunc := func(ctx context.Context) error {
		return ErrCacheOperation
	}

	err = manager.Warm(ctx, errorWarmFunc)
	if err != ErrCacheOperation {
		t.Errorf("Expected ErrCacheOperation, got %v", err)
	}
}

func TestCacheManagerTTLJitter(t *testing.T) {
	mockCache := newMockCache()
	config := &CacheConfig{
		Enabled:      true,
		DefaultTTL:   5 * time.Minute,
		KeyPrefix:    "pie:",
		EnableJitter: true,
		TTLJitter:    30 * time.Second,
	}
	manager := NewCacheManager(mockCache, config)
	ctx := context.Background()

	// 测试TTL抖动功能
	err := manager.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置（具体TTL值由applyJitter函数处理）
	if _, exists := mockCache.data["pie:testkey"]; !exists {
		t.Error("Expected key to be set")
	}
}

func TestCacheManagerKeyPrefix(t *testing.T) {
	mockCache := newMockCache()
	config := &CacheConfig{
		Enabled:    true,
		DefaultTTL: 5 * time.Minute,
		KeyPrefix:  "custom:",
	}
	manager := NewCacheManager(mockCache, config)
	ctx := context.Background()

	// 测试自定义前缀
	err := manager.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证使用了自定义前缀
	if _, exists := mockCache.data["custom:testkey"]; !exists {
		t.Error("Expected key with custom prefix to be set")
	}
	if _, exists := mockCache.data["pie:testkey"]; exists {
		t.Error("Expected default prefix not to be used")
	}
}
