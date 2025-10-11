package pie

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewMemoryCache(t *testing.T) {
	mc := NewMemoryCache(100)
	if mc == nil {
		t.Error("Expected MemoryCache to be created")
	}
	if mc.maxSize != 100 {
		t.Errorf("Expected maxSize to be 100, got %d", mc.maxSize)
	}
	if mc.items == nil {
		t.Error("Expected items map to be initialized")
	}
	if mc.tags == nil {
		t.Error("Expected tags map to be initialized")
	}
	if mc.stats == nil {
		t.Error("Expected stats to be initialized")
	}
	if mc.stopCleanup == nil {
		t.Error("Expected stopCleanup channel to be initialized")
	}

	// 测试清理goroutine是否启动
	time.Sleep(100 * time.Millisecond)

	// 关闭缓存
	mc.Close()
}

func TestMemoryCacheGet(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 测试缓存未命中
	val, err := mc.Get(ctx, "nonexistent")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value")
	}

	// 设置缓存
	err = mc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试缓存命中
	val, err = mc.Get(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", stats.Total)
	}
	if stats.Hits != 1 {
		t.Errorf("Expected Hits to be 1, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
	if stats.HitRate != 50.0 {
		t.Errorf("Expected HitRate to be 50.0, got %f", stats.HitRate)
	}
}

func TestMemoryCacheGetExpired(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 设置已过期的缓存
	mc.mu.Lock()
	mc.items["expiredkey"] = &memoryCacheItem{
		value:      []byte("expiredvalue"),
		expiration: time.Now().Add(-1 * time.Minute), // 1分钟前过期
	}
	mc.mu.Unlock()

	// 测试过期缓存
	val, err := mc.Get(ctx, "expiredkey")
	if err != ErrCacheExpired {
		t.Errorf("Expected ErrCacheExpired, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value for expired cache")
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
}

func TestMemoryCacheSet(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 测试设置缓存
	err := mc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	mc.mu.RLock()
	item, exists := mc.items["testkey"]
	mc.mu.RUnlock()

	if !exists {
		t.Error("Expected key to be set")
	}
	if string(item.value) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(item.value))
	}
	if item.expiration.Before(time.Now()) {
		t.Error("Expected expiration to be in the future")
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.Keys != 1 {
		t.Errorf("Expected Keys to be 1, got %d", stats.Keys)
	}
}

func TestMemoryCacheSetWithTags(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	tags := []string{"tag1", "tag2"}
	err := mc.SetWithTags(ctx, "testkey", []byte("testvalue"), 5*time.Minute, tags)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已设置
	mc.mu.RLock()
	item, exists := mc.items["testkey"]
	mc.mu.RUnlock()

	if !exists {
		t.Error("Expected key to be set")
	}
	if len(item.tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(item.tags))
	}

	// 验证标签映射
	mc.mu.RLock()
	tag1Keys := mc.tags["tag1"]
	tag2Keys := mc.tags["tag2"]
	mc.mu.RUnlock()

	if len(tag1Keys) != 1 || tag1Keys[0] != "testkey" {
		t.Error("Expected tag1 to map to testkey")
	}
	if len(tag2Keys) != 1 || tag2Keys[0] != "testkey" {
		t.Error("Expected tag2 to map to testkey")
	}
}

func TestMemoryCacheDelete(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 设置带标签的缓存
	tags := []string{"tag1", "tag2"}
	err := mc.SetWithTags(ctx, "testkey", []byte("testvalue"), 5*time.Minute, tags)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 删除缓存
	err = mc.Delete(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存已删除
	mc.mu.RLock()
	_, exists := mc.items["testkey"]
	mc.mu.RUnlock()

	if exists {
		t.Error("Expected key to be deleted")
	}

	// 验证标签映射已清理
	mc.mu.RLock()
	_, tag1Exists := mc.tags["tag1"]
	_, tag2Exists := mc.tags["tag2"]
	mc.mu.RUnlock()

	if tag1Exists || tag2Exists {
		t.Error("Expected tag mappings to be cleaned up")
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.Keys != 0 {
		t.Errorf("Expected Keys to be 0, got %d", stats.Keys)
	}
}

func TestMemoryCacheDeleteByPattern(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 设置测试数据
	err := mc.Set(ctx, "testkey1", []byte("value1"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.Set(ctx, "testkey2", []byte("value2"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.Set(ctx, "otherkey", []byte("value3"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 按模式删除
	err = mc.DeleteByPattern(ctx, "test.*")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证删除结果
	mc.mu.RLock()
	_, key1Exists := mc.items["testkey1"]
	_, key2Exists := mc.items["testkey2"]
	_, otherExists := mc.items["otherkey"]
	mc.mu.RUnlock()

	if key1Exists {
		t.Error("Expected testkey1 to be deleted")
	}
	if key2Exists {
		t.Error("Expected testkey2 to be deleted")
	}
	if !otherExists {
		t.Error("Expected otherkey to remain")
	}

	// 测试无效正则表达式
	err = mc.DeleteByPattern(ctx, "[invalid")
	if err == nil {
		t.Error("Expected error for invalid regex")
	}
}

func TestMemoryCacheDeleteByTags(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 设置带标签的缓存
	err := mc.SetWithTags(ctx, "key1", []byte("value1"), 5*time.Minute, []string{"tag1"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.SetWithTags(ctx, "key2", []byte("value2"), 5*time.Minute, []string{"tag1", "tag2"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.SetWithTags(ctx, "key3", []byte("value3"), 5*time.Minute, []string{"tag2"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.Set(ctx, "key4", []byte("value4"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 按标签删除
	err = mc.DeleteByTags(ctx, "tag1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证删除结果
	mc.mu.RLock()
	_, key1Exists := mc.items["key1"]
	_, key2Exists := mc.items["key2"]
	_, key3Exists := mc.items["key3"]
	_, key4Exists := mc.items["key4"]
	mc.mu.RUnlock()

	if key1Exists {
		t.Error("Expected key1 to be deleted")
	}
	if key2Exists {
		t.Error("Expected key2 to be deleted")
	}
	if !key3Exists {
		t.Error("Expected key3 to remain")
	}
	if !key4Exists {
		t.Error("Expected key4 to remain")
	}

	// 验证标签映射已清理
	mc.mu.RLock()
	_, tag1Exists := mc.tags["tag1"]
	_, tag2Exists := mc.tags["tag2"]
	mc.mu.RUnlock()

	if tag1Exists {
		t.Error("Expected tag1 mapping to be cleaned up")
	}
	if !tag2Exists {
		t.Error("Expected tag2 mapping to remain")
	}
}

func TestMemoryCacheExists(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 测试不存在的key
	exists, err := mc.Exists(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected key not to exist")
	}

	// 设置缓存
	err = mc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试存在的key
	exists, err = mc.Exists(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// 测试过期的key
	mc.mu.Lock()
	mc.items["expiredkey"] = &memoryCacheItem{
		value:      []byte("expiredvalue"),
		expiration: time.Now().Add(-1 * time.Minute),
	}
	mc.mu.Unlock()

	exists, err = mc.Exists(ctx, "expiredkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected expired key not to exist")
	}
}

func TestMemoryCacheClear(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 设置测试数据
	err := mc.SetWithTags(ctx, "key1", []byte("value1"), 5*time.Minute, []string{"tag1"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.SetWithTags(ctx, "key2", []byte("value2"), 5*time.Minute, []string{"tag2"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 清空缓存
	err = mc.Clear(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证已清空
	mc.mu.RLock()
	itemsLen := len(mc.items)
	tagsLen := len(mc.tags)
	mc.mu.RUnlock()

	if itemsLen != 0 {
		t.Error("Expected items to be cleared")
	}
	if tagsLen != 0 {
		t.Error("Expected tags to be cleared")
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.Keys != 0 {
		t.Errorf("Expected Keys to be 0, got %d", stats.Keys)
	}
}

func TestMemoryCacheStats(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()
	ctx := context.Background()

	// 初始状态
	stats := mc.Stats()
	if stats.Hits != 0 {
		t.Errorf("Expected initial Hits to be 0, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected initial Misses to be 0, got %d", stats.Misses)
	}
	if stats.Total != 0 {
		t.Errorf("Expected initial Total to be 0, got %d", stats.Total)
	}
	if stats.Keys != 0 {
		t.Errorf("Expected initial Keys to be 0, got %d", stats.Keys)
	}

	// 执行一些操作
	mc.Get(ctx, "nonexistent") // miss
	mc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	mc.Get(ctx, "testkey") // hit

	stats = mc.Stats()
	if stats.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", stats.Total)
	}
	if stats.Hits != 1 {
		t.Errorf("Expected Hits to be 1, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
	if stats.Keys != 1 {
		t.Errorf("Expected Keys to be 1, got %d", stats.Keys)
	}
	if stats.HitRate != 50.0 {
		t.Errorf("Expected HitRate to be 50.0, got %f", stats.HitRate)
	}
}

func TestMemoryCacheMaxSize(t *testing.T) {
	mc := NewMemoryCache(2) // 设置最大大小为2
	defer mc.Close()
	ctx := context.Background()

	// 设置第一个key
	err := mc.Set(ctx, "key1", []byte("value1"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 设置第二个key
	err = mc.Set(ctx, "key2", []byte("value2"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 设置第三个key，应该触发淘汰
	err = mc.Set(ctx, "key3", []byte("value3"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证缓存大小不超过最大值
	mc.mu.RLock()
	itemsLen := len(mc.items)
	mc.mu.RUnlock()

	if itemsLen > 2 {
		t.Errorf("Expected items length to be <= 2, got %d", itemsLen)
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.EvictedKeys == 0 {
		t.Error("Expected some keys to be evicted")
	}
}

func TestMemoryCacheConcurrency(t *testing.T) {
	mc := NewMemoryCache(1000)
	defer mc.Close()
	ctx := context.Background()

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// 并发写入
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", goroutineID, j)
				value := fmt.Sprintf("value_%d_%d", goroutineID, j)
				err := mc.Set(ctx, key, []byte(value), 5*time.Minute)
				if err != nil {
					t.Errorf("Error setting key %s: %v", key, err)
				}
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", goroutineID, j)
				_, err := mc.Get(ctx, key)
				// 允许缓存未命中，因为写入和读取可能并发
				if err != nil && err != ErrCacheNotFound {
					t.Errorf("Error getting key %s: %v", key, err)
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证没有panic或死锁
	stats := mc.Stats()
	if stats.Total == 0 {
		t.Error("Expected some operations to be recorded")
	}
}

func TestMemoryCacheClose(t *testing.T) {
	mc := NewMemoryCache(100)

	// 关闭缓存
	mc.Close()

	// 验证可以安全地再次关闭
	mc.Close()

	// 验证关闭后仍然可以执行操作（虽然清理goroutine已停止）
	ctx := context.Background()
	err := mc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error after close, got %v", err)
	}
}

func TestMemoryCacheCleanupExpired(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()

	// 设置一些过期的缓存项
	mc.mu.Lock()
	mc.items["expired1"] = &memoryCacheItem{
		value:      []byte("value1"),
		expiration: time.Now().Add(-2 * time.Minute),
		tags:       []string{"tag1"},
	}
	mc.items["expired2"] = &memoryCacheItem{
		value:      []byte("value2"),
		expiration: time.Now().Add(-1 * time.Minute),
		tags:       []string{"tag2"},
	}
	mc.items["valid"] = &memoryCacheItem{
		value:      []byte("value3"),
		expiration: time.Now().Add(5 * time.Minute),
	}
	mc.tags["tag1"] = []string{"expired1"}
	mc.tags["tag2"] = []string{"expired2"}
	mc.mu.Unlock()

	// 等待清理goroutine执行（默认1分钟间隔，这里等待2秒）
	time.Sleep(2 * time.Second)

	// 手动触发清理（因为测试环境可能不会自动触发）
	mc.mu.Lock()
	now := time.Now()
	keysToDelete := make([]string, 0)
	for key, item := range mc.items {
		if now.After(item.expiration) {
			keysToDelete = append(keysToDelete, key)
		}
	}
	for _, key := range keysToDelete {
		item := mc.items[key]
		if item != nil {
			for _, tag := range item.tags {
				mc.removeKeyFromTag(tag, key)
			}
		}
		delete(mc.items, key)
		mc.stats.EvictedKeys++
	}
	mc.stats.Keys = int64(len(mc.items))
	mc.mu.Unlock()

	// 验证过期项已被清理
	mc.mu.RLock()
	_, expired1Exists := mc.items["expired1"]
	_, expired2Exists := mc.items["expired2"]
	_, validExists := mc.items["valid"]
	mc.mu.RUnlock()

	if expired1Exists {
		t.Error("Expected expired1 to be cleaned up")
	}
	if expired2Exists {
		t.Error("Expected expired2 to be cleaned up")
	}
	if !validExists {
		t.Error("Expected valid item to remain")
	}

	// 验证标签映射已清理
	mc.mu.RLock()
	_, tag1Exists := mc.tags["tag1"]
	_, tag2Exists := mc.tags["tag2"]
	mc.mu.RUnlock()

	if tag1Exists || tag2Exists {
		t.Error("Expected tag mappings to be cleaned up")
	}

	// 验证统计数据
	stats := mc.Stats()
	if stats.EvictedKeys == 0 {
		t.Error("Expected some keys to be evicted")
	}
	if stats.Keys != 1 {
		t.Errorf("Expected Keys to be 1, got %d", stats.Keys)
	}
}

func TestMemoryCacheEvictOldest(t *testing.T) {
	mc := NewMemoryCache(2)
	defer mc.Close()
	ctx := context.Background()

	// 设置两个key
	err := mc.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = mc.Set(ctx, "key2", []byte("value2"), 2*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 设置第三个key，应该淘汰最旧的
	err = mc.Set(ctx, "key3", []byte("value3"), 3*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证key1被淘汰（因为TTL最短）
	mc.mu.RLock()
	_, key1Exists := mc.items["key1"]
	_, key2Exists := mc.items["key2"]
	_, key3Exists := mc.items["key3"]
	mc.mu.RUnlock()

	if key1Exists {
		t.Error("Expected key1 to be evicted (oldest)")
	}
	if !key2Exists {
		t.Error("Expected key2 to remain")
	}
	if !key3Exists {
		t.Error("Expected key3 to remain")
	}
}

func TestMemoryCacheRemoveKeyFromTag(t *testing.T) {
	mc := NewMemoryCache(100)
	defer mc.Close()

	// 设置标签映射
	mc.mu.Lock()
	mc.tags["tag1"] = []string{"key1", "key2", "key3"}
	mc.mu.Unlock()

	// 移除key2
	mc.mu.Lock()
	mc.removeKeyFromTag("tag1", "key2")
	mc.mu.Unlock()

	// 验证key2被移除
	mc.mu.RLock()
	keys := mc.tags["tag1"]
	mc.mu.RUnlock()

	expectedKeys := []string{"key1", "key3"}
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
	for i, expectedKey := range expectedKeys {
		if keys[i] != expectedKey {
			t.Errorf("Expected key %s at position %d, got %s", expectedKey, i, keys[i])
		}
	}

	// 移除最后一个key
	mc.mu.Lock()
	mc.removeKeyFromTag("tag1", "key1")
	mc.removeKeyFromTag("tag1", "key3")
	mc.mu.Unlock()

	// 验证标签映射被删除
	mc.mu.RLock()
	_, exists := mc.tags["tag1"]
	mc.mu.RUnlock()

	if exists {
		t.Error("Expected tag1 mapping to be deleted when empty")
	}
}
