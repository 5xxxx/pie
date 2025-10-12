package pie

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewSingleFlight(t *testing.T) {
	sf := NewSingleFlight()
	if sf == nil {
		t.Error("Expected SingleFlight to be created")
	}
	if sf.m == nil {
		t.Error("Expected map to be initialized")
	}
	if len(sf.m) != 0 {
		t.Error("Expected map to be empty")
	}
}

func TestSingleFlightDo(t *testing.T) {
	sf := NewSingleFlight()

	// 测试基本功能
	val, err := sf.Do("key1", func() ([]byte, error) {
		return []byte("value1"), nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "value1" {
		t.Errorf("Expected 'value1', got %s", string(val))
	}

	// 测试错误处理
	val, err = sf.Do("key2", func() ([]byte, error) {
		return nil, errors.New("test error")
	})
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %s", err.Error())
	}
	if val != nil {
		t.Error("Expected nil value on error")
	}
}

func TestSingleFlightConcurrentSameKey(t *testing.T) {
	sf := NewSingleFlight()

	var wg sync.WaitGroup
	numGoroutines := 10
	results := make([]string, numGoroutines)
	errors := make([]error, numGoroutines)

	// 并发执行相同key的操作
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			val, err := sf.Do("samekey", func() ([]byte, error) {
				// 模拟耗时操作
				time.Sleep(100 * time.Millisecond)
				return []byte("samevalue"), nil
			})
			results[index] = string(val)
			errors[index] = err
		}(i)
	}

	wg.Wait()

	// 验证所有结果都相同
	for i := 0; i < numGoroutines; i++ {
		if errors[i] != nil {
			t.Errorf("Goroutine %d got error: %v", i, errors[i])
		}
		if results[i] != "samevalue" {
			t.Errorf("Goroutine %d got wrong value: %s", i, results[i])
		}
	}

	// 验证map被清理
	if len(sf.m) != 0 {
		t.Error("Expected map to be cleaned up")
	}
}

func TestSingleFlightConcurrentDifferentKeys(t *testing.T) {
	sf := NewSingleFlight()

	var wg sync.WaitGroup
	numGoroutines := 5
	results := make([]string, numGoroutines)
	errors := make([]error, numGoroutines)

	// 并发执行不同key的操作
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", index)
			expectedValue := fmt.Sprintf("value%d", index)
			val, err := sf.Do(key, func() ([]byte, error) {
				// 模拟耗时操作
				time.Sleep(50 * time.Millisecond)
				return []byte(expectedValue), nil
			})
			results[index] = string(val)
			errors[index] = err
		}(i)
	}

	wg.Wait()

	// 验证每个结果都不同且正确
	for i := 0; i < numGoroutines; i++ {
		if errors[i] != nil {
			t.Errorf("Goroutine %d got error: %v", i, errors[i])
		}
		expectedValue := fmt.Sprintf("value%d", i)
		if results[i] != expectedValue {
			t.Errorf("Goroutine %d got wrong value: %s, expected: %s", i, results[i], expectedValue)
		}
	}

	// 验证map被清理
	if len(sf.m) != 0 {
		t.Error("Expected map to be cleaned up")
	}
}

func TestSingleFlightPanicRecovery(t *testing.T) {
	sf := NewSingleFlight()

	// 测试panic恢复
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected panic to be handled, got: %v", r)
		}
	}()

	val, err := sf.Do("panickey", func() ([]byte, error) {
		panic("test panic")
	})

	if err == nil {
		t.Error("Expected error from panic")
	}
	if val != nil {
		t.Error("Expected nil value from panic")
	}
}

func TestNewCacheWithSingleFlight(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}

	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	if csf == nil {
		t.Error("Expected CacheWithSingleFlight to be created")
	}
	if csf.cache != singleFlightMockCache {
		t.Error("Expected cache to be set correctly")
	}
	if csf.singleFlight == nil {
		t.Error("Expected singleFlight to be initialized")
	}
}

func TestCacheWithSingleFlightGet(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 设置缓存
	singleFlightMockCache.data["testkey"] = []byte("testvalue")

	// 测试获取缓存
	val, err := csf.Get(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 测试缓存未命中
	val, err = csf.Get(ctx, "nonexistent")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value")
	}
}

func TestCacheWithSingleFlightSet(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 测试设置缓存
	err := csf.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	if val, exists := singleFlightMockCache.data["testkey"]; !exists {
		t.Error("Expected key to be set")
	} else if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}
}

func TestCacheWithSingleFlightDelete(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 设置测试数据
	singleFlightMockCache.data["testkey"] = []byte("testvalue")

	// 测试删除
	err := csf.Delete(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证已删除
	if _, exists := singleFlightMockCache.data["testkey"]; exists {
		t.Error("Expected key to be deleted")
	}
}

func TestCacheWithSingleFlightDeleteByPattern(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 设置测试数据
	singleFlightMockCache.data["testkey1"] = []byte("value1")
	singleFlightMockCache.data["testkey2"] = []byte("value2")
	singleFlightMockCache.data["otherkey"] = []byte("value3")

	// 测试按模式删除
	err := csf.DeleteByPattern(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证删除结果
	if _, exists := singleFlightMockCache.data["testkey1"]; exists {
		t.Error("Expected testkey1 to be deleted")
	}
	if _, exists := singleFlightMockCache.data["testkey2"]; exists {
		t.Error("Expected testkey2 to be deleted")
	}
	if _, exists := singleFlightMockCache.data["otherkey"]; !exists {
		t.Error("Expected otherkey to remain")
	}
}

func TestCacheWithSingleFlightDeleteByTags(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 测试按标签删除
	err := csf.DeleteByTags(ctx, "tag1", "tag2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCacheWithSingleFlightExists(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 设置测试数据
	singleFlightMockCache.data["testkey"] = []byte("testvalue")

	// 测试存在
	exists, err := csf.Exists(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// 测试不存在
	exists, err = csf.Exists(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected key not to exist")
	}
}

func TestCacheWithSingleFlightClear(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 设置测试数据
	singleFlightMockCache.data["testkey1"] = []byte("value1")
	singleFlightMockCache.data["testkey2"] = []byte("value2")

	// 测试清空
	err := csf.Clear(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证已清空
	if len(singleFlightMockCache.data) != 0 {
		t.Error("Expected cache to be cleared")
	}
}

func TestCacheWithSingleFlightStats(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)

	// 设置一些统计数据
	singleFlightMockCache.stats.Hits = 10
	singleFlightMockCache.stats.Misses = 5
	singleFlightMockCache.stats.Total = 15
	singleFlightMockCache.stats.HitRate = 66.67

	stats := csf.Stats()
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

func TestCacheWithSingleFlightConcurrentGet(t *testing.T) {
	singleFlightMockCache := &singleFlightMockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	csf := NewCacheWithSingleFlight(singleFlightMockCache)
	ctx := context.Background()

	// 设置缓存
	singleFlightMockCache.data["testkey"] = []byte("testvalue")

	var wg sync.WaitGroup
	numGoroutines := 10
	results := make([]string, numGoroutines)
	errors := make([]error, numGoroutines)

	// 并发获取相同key
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			val, err := csf.Get(ctx, "testkey")
			results[index] = string(val)
			errors[index] = err
		}(i)
	}

	wg.Wait()

	// 验证所有结果都相同
	for i := 0; i < numGoroutines; i++ {
		if errors[i] != nil {
			t.Errorf("Goroutine %d got error: %v", i, errors[i])
		}
		if results[i] != "testvalue" {
			t.Errorf("Goroutine %d got wrong value: %s", i, results[i])
		}
	}
}

func TestSingleFlightCleanup(t *testing.T) {
	sf := NewSingleFlight()

	// 测试map清理
	sf.mu.Lock()
	sf.m["testkey"] = &call{}
	sf.mu.Unlock()

	// 执行操作，应该清理map
	val, err := sf.Do("testkey", func() ([]byte, error) {
		return []byte("testvalue"), nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 验证map被清理
	sf.mu.Lock()
	if len(sf.m) != 0 {
		t.Error("Expected map to be cleaned up")
	}
	sf.mu.Unlock()
}

// singleFlightMockCache 用于测试CacheWithSingleFlight
type singleFlightMockCache struct {
	data  map[string][]byte
	stats *CacheStats
}

func (m *singleFlightMockCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.stats.Total++
	if val, exists := m.data[key]; exists {
		m.stats.Hits++
		m.stats.HitRate = float64(m.stats.Hits) / float64(m.stats.Total) * 100
		return val, nil
	}
	m.stats.Misses++
	return nil, ErrCacheNotFound
}

func (m *singleFlightMockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.data[key] = value
	m.stats.Keys = int64(len(m.data))
	return nil
}

func (m *singleFlightMockCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	m.stats.Keys = int64(len(m.data))
	return nil
}

func (m *singleFlightMockCache) DeleteByPattern(ctx context.Context, pattern string) error {
	// 简单实现，删除包含pattern的key
	for key := range m.data {
		if strings.Contains(key, pattern) {
			delete(m.data, key)
		}
	}
	m.stats.Keys = int64(len(m.data))
	return nil
}

func (m *singleFlightMockCache) DeleteByTags(ctx context.Context, tags ...string) error {
	m.data = make(map[string][]byte)
	m.stats.Keys = 0
	return nil
}

func (m *singleFlightMockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.data[key]
	return exists, nil
}

func (m *singleFlightMockCache) Clear(ctx context.Context) error {
	m.data = make(map[string][]byte)
	m.stats.Keys = 0
	return nil
}

func (m *singleFlightMockCache) Stats() *CacheStats {
	stats := *m.stats
	return &stats
}
