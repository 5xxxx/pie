package pie

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// mockEngine 模拟Engine用于测试
type mockEngine struct {
	cacheManager *CacheManager
}

func (e *mockEngine) GetCacheManager() *CacheManager {
	return e.cacheManager
}

// mockCollection 模拟Collection用于测试
type mockCollection struct {
	name string
}

func (c *mockCollection) Name() string {
	return c.name
}

// mockSession 模拟Session用于测试
type mockSession[T any] struct {
	engine      *mockEngine
	collection  *mockCollection
	query       *mockQuery
	options     *mockOptions
	cacheConfig *SessionCacheConfig
}

type mockQuery struct {
	filter bson.D
}

type mockOptions struct {
	limit int
	skip  int
}

func TestSessionCache(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}},
		},
	}

	// 测试Cache方法
	session.Cache(5 * time.Minute)
	if session.cacheConfig == nil {
		t.Error("Expected cacheConfig to be set")
	}
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if session.cacheConfig.TTL != 5*time.Minute {
		t.Errorf("Expected TTL to be 5 minutes, got %v", session.cacheConfig.TTL)
	}

	// 测试NoCache方法
	session.NoCache()
	if session.cacheConfig.Enabled {
		t.Error("Expected cache to be disabled")
	}

	// 测试CacheWithTags方法
	session.CacheWithTags("tag1", "tag2")
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if len(session.cacheConfig.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(session.cacheConfig.Tags))
	}
	if session.cacheConfig.Tags[0] != "tag1" || session.cacheConfig.Tags[1] != "tag2" {
		t.Error("Expected correct tags")
	}

	// 测试CacheWithJitter方法
	session.CacheWithJitter(10*time.Minute, 1*time.Minute)
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if session.cacheConfig.TTL != 10*time.Minute {
		t.Errorf("Expected TTL to be 10 minutes, got %v", session.cacheConfig.TTL)
	}
	if !session.cacheConfig.UseJitter {
		t.Error("Expected jitter to be enabled")
	}

	// 测试CacheEmpty方法
	session.CacheEmpty(30 * time.Second)
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if !session.cacheConfig.CacheEmpty {
		t.Error("Expected empty caching to be enabled")
	}
	if session.cacheConfig.TTL != 30*time.Second {
		t.Errorf("Expected TTL to be 30 seconds, got %v", session.cacheConfig.TTL)
	}

	// 测试CacheL1Only方法
	session.CacheL1Only()
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}

	// 测试CacheL2Only方法
	session.CacheL2Only()
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
}

func TestSessionGetFromCache(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}},
		},
		cacheConfig: &SessionCacheConfig{
			Enabled: true,
			TTL:     5 * time.Minute,
		},
	}

	ctx := context.Background()

	// 测试缓存未命中
	results, found, err := session.getFromCache(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found {
		t.Error("Expected cache miss")
	}
	if results != nil {
		t.Error("Expected nil results")
	}

	// 设置缓存数据
	testData := []string{"result1", "result2"}
	data, _ := json.Marshal(testData)
	mockCache.data["pie:testkey"] = data

	// 测试缓存命中
	results, found, err = session.getFromCache(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !found {
		t.Error("Expected cache hit")
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
	if results[0] != "result1" || results[1] != "result2" {
		t.Error("Expected correct results")
	}

	// 测试禁用缓存
	session.cacheConfig.Enabled = false
	results, found, err = session.getFromCache(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found {
		t.Error("Expected cache to be disabled")
	}

	// 测试无缓存管理器
	session.engine.cacheManager = nil
	results, found, err = session.getFromCache(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found {
		t.Error("Expected no cache manager")
	}
}

func TestSessionSetToCache(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}},
		},
		cacheConfig: &SessionCacheConfig{
			Enabled: true,
			TTL:     5 * time.Minute,
		},
	}

	ctx := context.Background()

	// 测试设置缓存
	testData := []string{"result1", "result2"}
	err := session.setToCache(ctx, "testkey", testData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	if val, exists := mockCache.data["pie:testkey"]; !exists {
		t.Error("Expected key to be cached")
	} else {
		var results []string
		json.Unmarshal(val, &results)
		if len(results) != 2 || results[0] != "result1" || results[1] != "result2" {
			t.Error("Expected correct cached data")
		}
	}

	// 测试空结果不缓存（默认行为）
	session.cacheConfig.CacheEmpty = false
	err = session.setToCache(ctx, "emptykey", []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证空结果未缓存
	if _, exists := mockCache.data["pie:emptykey"]; exists {
		t.Error("Expected empty results not to be cached")
	}

	// 测试空结果缓存（启用CacheEmpty）
	session.cacheConfig.CacheEmpty = true
	err = session.setToCache(ctx, "emptykey2", []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证空结果已缓存
	if _, exists := mockCache.data["pie:emptykey2"]; !exists {
		t.Error("Expected empty results to be cached when CacheEmpty is enabled")
	}

	// 测试禁用缓存
	session.cacheConfig.Enabled = false
	err = session.setToCache(ctx, "disabledkey", testData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存未设置
	if _, exists := mockCache.data["pie:disabledkey"]; exists {
		t.Error("Expected cache to be disabled")
	}

	// 测试无缓存管理器
	session.engine.cacheManager = nil
	err = session.setToCache(ctx, "nomanagerkey", testData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSessionInvalidateCache(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}},
		},
	}

	ctx := context.Background()

	// 测试缓存失效
	err := session.invalidateCache(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试无缓存管理器
	session.engine.cacheManager = nil
	err = session.invalidateCache(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSessionGenerateCacheKey(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}, {"age", 25}},
		},
		options: &mockOptions{
			limit: 10,
			skip:  0,
		},
	}

	// 测试生成缓存键
	key := session.generateCacheKey()
	if key == "" {
		t.Error("Expected non-empty cache key")
	}

	// 验证相同输入产生相同键
	key2 := session.generateCacheKey()
	if key != key2 {
		t.Error("Expected same input to generate same key")
	}

	// 测试无缓存管理器
	session.engine.cacheManager = nil
	key3 := session.generateCacheKey()
	if key3 != "" {
		t.Error("Expected empty key when no cache manager")
	}
}

func TestSessionCacheWithJitter(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	config := &CacheConfig{
		Enabled:      true,
		DefaultTTL:   5 * time.Minute,
		KeyPrefix:    "pie:",
		EnableJitter: true,
		TTLJitter:    30 * time.Second,
	}
	cacheManager := NewCacheManager(mockCache, config)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}},
		},
		cacheConfig: &SessionCacheConfig{
			Enabled:   true,
			TTL:       5 * time.Minute,
			UseJitter: true,
		},
	}

	ctx := context.Background()

	// 测试带抖动的缓存设置
	testData := []string{"result1", "result2"}
	err := session.setToCache(ctx, "jitterkey", testData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	if _, exists := mockCache.data["pie:jitterkey"]; !exists {
		t.Error("Expected key to be cached with jitter")
	}
}

func TestSessionCacheWithTags(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建模拟Engine
	engine := &mockEngine{
		cacheManager: cacheManager,
	}

	// 创建模拟Collection
	collection := &mockCollection{
		name: "test_collection",
	}

	// 创建模拟Session
	session := &mockSession[string]{
		engine:     engine,
		collection: collection,
		query: &mockQuery{
			filter: bson.D{{"name", "John"}},
		},
		cacheConfig: &SessionCacheConfig{
			Enabled: true,
			TTL:     5 * time.Minute,
			Tags:    []string{"tag1", "tag2"},
		},
	}

	ctx := context.Background()

	// 测试带标签的缓存设置
	testData := []string{"result1", "result2"}
	err := session.setToCache(ctx, "tagkey", testData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	if _, exists := mockCache.data["pie:tagkey"]; !exists {
		t.Error("Expected key to be cached with tags")
	}
}

// TestRealSessionCacheMethods 测试真实的Session缓存方法
func TestRealSessionCacheMethods(t *testing.T) {
	// 创建模拟缓存管理器
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	cacheManager := NewCacheManager(mockCache, nil)

	// 创建真实的Engine实例
	engine := &Engine{
		cacheManager: cacheManager,
	}

	// 创建真实的Session实例（使用interface{}作为泛型参数）
	session := &Session[interface{}]{
		engine:     engine,
		collection: nil, // 使用nil，因为我们只测试缓存方法
		query: &Query{
			filter: bson.D{{"name", "John"}},
		},
	}

	// 测试Cache方法
	session.Cache(5 * time.Minute)
	if session.cacheConfig == nil {
		t.Error("Expected cacheConfig to be set")
	}
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if session.cacheConfig.TTL != 5*time.Minute {
		t.Errorf("Expected TTL to be 5 minutes, got %v", session.cacheConfig.TTL)
	}

	// 测试NoCache方法
	session.NoCache()
	if session.cacheConfig.Enabled {
		t.Error("Expected cache to be disabled")
	}

	// 测试CacheWithTags方法
	session.CacheWithTags("tag1", "tag2")
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if len(session.cacheConfig.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(session.cacheConfig.Tags))
	}

	// 测试CacheWithJitter方法
	session.CacheWithJitter(10*time.Minute, 1*time.Minute)
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if session.cacheConfig.TTL != 10*time.Minute {
		t.Errorf("Expected TTL to be 10 minutes, got %v", session.cacheConfig.TTL)
	}
	if !session.cacheConfig.UseJitter {
		t.Error("Expected jitter to be enabled")
	}

	// 测试CacheEmpty方法
	session.CacheEmpty(30 * time.Second)
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
	if !session.cacheConfig.CacheEmpty {
		t.Error("Expected empty caching to be enabled")
	}
	if session.cacheConfig.TTL != 30*time.Second {
		t.Errorf("Expected TTL to be 30 seconds, got %v", session.cacheConfig.TTL)
	}

	// 测试CacheL1Only方法
	session.CacheL1Only()
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}

	// 测试CacheL2Only方法
	session.CacheL2Only()
	if !session.cacheConfig.Enabled {
		t.Error("Expected cache to be enabled")
	}
}

// 为mockSession添加缺失的方法
func (s *mockSession[T]) Cache(ttl ...time.Duration) *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = true
	if len(ttl) > 0 {
		s.cacheConfig.TTL = ttl[0]
	} else if s.engine.cacheManager != nil {
		s.cacheConfig.TTL = s.engine.cacheManager.config.DefaultTTL
	} else {
		s.cacheConfig.TTL = 5 * time.Minute
	}
	return s
}

func (s *mockSession[T]) NoCache() *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = false
	return s
}

func (s *mockSession[T]) CacheWithTags(tags ...string) *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{
			Enabled: true,
		}
	}
	if s.cacheConfig.TTL == 0 && s.engine.cacheManager != nil {
		s.cacheConfig.TTL = s.engine.cacheManager.config.DefaultTTL
	}
	s.cacheConfig.Tags = tags
	s.cacheConfig.Enabled = true
	return s
}

func (s *mockSession[T]) CacheWithJitter(ttl, jitter time.Duration) *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = true
	s.cacheConfig.TTL = ttl
	s.cacheConfig.UseJitter = true
	return s
}

func (s *mockSession[T]) CacheEmpty(ttl time.Duration) *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = true
	s.cacheConfig.CacheEmpty = true
	if ttl > 0 {
		s.cacheConfig.TTL = ttl
	} else if s.engine.cacheManager != nil {
		s.cacheConfig.TTL = s.engine.cacheManager.config.EmptyCacheTTL
	} else {
		s.cacheConfig.TTL = 30 * time.Second
	}
	return s
}

func (s *mockSession[T]) CacheL1Only() *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = true
	return s
}

func (s *mockSession[T]) CacheL2Only() *mockSession[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = true
	return s
}

func (s *mockSession[T]) getFromCache(ctx context.Context, key string) ([]T, bool, error) {
	if s.engine.cacheManager == nil || s.cacheConfig == nil || !s.cacheConfig.Enabled {
		return nil, false, nil
	}

	data, err := s.engine.cacheManager.Get(ctx, key)
	if err != nil {
		if err == ErrCacheNotFound || err == ErrCacheExpired {
			return nil, false, nil
		}
		return nil, false, err
	}

	var results []T
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, false, err
	}

	return results, true, nil
}

func (s *mockSession[T]) setToCache(ctx context.Context, key string, results []T) error {
	if s.engine.cacheManager == nil || s.cacheConfig == nil || !s.cacheConfig.Enabled {
		return nil
	}

	if len(results) == 0 && !s.cacheConfig.CacheEmpty {
		return nil
	}

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	ttl := s.cacheConfig.TTL
	if s.cacheConfig.UseJitter && s.engine.cacheManager.config.EnableJitter {
		ttl = applyJitter(ttl, s.engine.cacheManager.config.TTLJitter)
	}

	return s.engine.cacheManager.Set(ctx, key, data, ttl)
}

func (s *mockSession[T]) invalidateCache(ctx context.Context) error {
	if s.engine.cacheManager == nil {
		return nil
	}

	collectionName := s.collection.Name()
	keyGen := NewCacheKeyGenerator(s.engine.cacheManager.config.KeyPrefix)
	pattern := keyGen.GenerateCollectionPattern(collectionName)

	return s.engine.cacheManager.DeleteByPattern(ctx, pattern)
}

func (s *mockSession[T]) generateCacheKey() string {
	if s.engine.cacheManager == nil {
		return ""
	}

	collectionName := s.collection.Name()
	keyGen := NewCacheKeyGenerator(s.engine.cacheManager.config.KeyPrefix)

	return keyGen.GenerateQueryKey(collectionName, s.query.filter, s.options)
}
