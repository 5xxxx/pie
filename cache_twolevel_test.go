package pie

import (
	"context"
	"testing"
	"time"
)

func TestDefaultTwoLevelCacheConfig(t *testing.T) {
	config := DefaultTwoLevelCacheConfig()

	if config.L1TTL != 1*time.Minute {
		t.Errorf("Expected L1TTL to be 1 minute, got %v", config.L1TTL)
	}
	if config.L2TTL != 10*time.Minute {
		t.Errorf("Expected L2TTL to be 10 minutes, got %v", config.L2TTL)
	}
	if !config.SyncOnWrite {
		t.Error("Expected SyncOnWrite to be true")
	}
}

func TestNewTwoLevelCache(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})

	// 测试使用默认配置
	tlc := NewTwoLevelCache(l1, l2, nil)
	if tlc == nil {
		t.Error("Expected TwoLevelCache to be created")
	}
	if tlc.l1 != l1 {
		t.Error("Expected L1 cache to be set correctly")
	}
	if tlc.l2 != l2 {
		t.Error("Expected L2 cache to be set correctly")
	}
	if tlc.config == nil {
		t.Error("Expected config to be set")
	}
	if tlc.stats == nil {
		t.Error("Expected stats to be initialized")
	}

	// 测试使用自定义配置
	customConfig := &TwoLevelCacheConfig{
		L1TTL:       2 * time.Minute,
		L2TTL:       20 * time.Minute,
		SyncOnWrite: false,
	}
	tlc2 := NewTwoLevelCache(l1, l2, customConfig)
	if tlc2.config.L1TTL != 2*time.Minute {
		t.Error("Expected custom L1TTL to be used")
	}
	if tlc2.config.L2TTL != 20*time.Minute {
		t.Error("Expected custom L2TTL to be used")
	}
	if tlc2.config.SyncOnWrite {
		t.Error("Expected custom SyncOnWrite to be used")
	}
}

func TestTwoLevelCacheGetL1Hit(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L1中设置缓存
	err := l1.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试L1命中
	val, err := tlc.Get(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 验证统计数据
	stats := tlc.StatsDetailed()
	if stats.L1Hits != 1 {
		t.Errorf("Expected L1Hits to be 1, got %d", stats.L1Hits)
	}
	if stats.L2Hits != 0 {
		t.Errorf("Expected L2Hits to be 0, got %d", stats.L2Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected Misses to be 0, got %d", stats.Misses)
	}
	if stats.Total != 1 {
		t.Errorf("Expected Total to be 1, got %d", stats.Total)
	}
}

func TestTwoLevelCacheGetL2Hit(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L2中设置缓存（模拟）
	// 注意：RedisCache是模拟版本，实际不会存储数据
	// 这里我们直接测试L2 miss的情况

	// 测试L1和L2都miss
	val, err := tlc.Get(ctx, "testkey")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value")
	}

	// 验证统计数据
	stats := tlc.StatsDetailed()
	if stats.L1Hits != 0 {
		t.Errorf("Expected L1Hits to be 0, got %d", stats.L1Hits)
	}
	if stats.L2Hits != 0 {
		t.Errorf("Expected L2Hits to be 0, got %d", stats.L2Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
	if stats.Total != 1 {
		t.Errorf("Expected Total to be 1, got %d", stats.Total)
	}
}

func TestTwoLevelCacheSet(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 测试设置缓存
	err := tlc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1中也有缓存（因为SyncOnWrite默认为true）
	val, err := l1.Get(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error from L1, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue' from L1, got %s", string(val))
	}
}

func TestTwoLevelCacheSetWithoutSync(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	config := &TwoLevelCacheConfig{
		L1TTL:       1 * time.Minute,
		L2TTL:       10 * time.Minute,
		SyncOnWrite: false,
	}
	tlc := NewTwoLevelCache(l1, l2, config)
	ctx := context.Background()

	// 测试设置缓存
	err := tlc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1中没有缓存（因为SyncOnWrite为false）
	val, err := l1.Get(ctx, "testkey")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound from L1, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value from L1")
	}
}

func TestTwoLevelCacheDelete(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L1中设置缓存
	err := l1.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试删除
	err = tlc.Delete(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1中已删除
	val, err := l1.Get(ctx, "testkey")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound from L1, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value from L1")
	}
}

func TestTwoLevelCacheDeleteByPattern(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L1中设置测试数据
	err := l1.Set(ctx, "testkey1", []byte("value1"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = l1.Set(ctx, "testkey2", []byte("value2"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = l1.Set(ctx, "otherkey", []byte("value3"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试按模式删除
	err = tlc.DeleteByPattern(ctx, "test.*")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1中删除结果
	_, err1 := l1.Get(ctx, "testkey1")
	_, err2 := l1.Get(ctx, "testkey2")
	val3, err3 := l1.Get(ctx, "otherkey")

	if err1 != ErrCacheNotFound {
		t.Error("Expected testkey1 to be deleted from L1")
	}
	if err2 != ErrCacheNotFound {
		t.Error("Expected testkey2 to be deleted from L1")
	}
	if err3 != nil {
		t.Error("Expected otherkey to remain in L1")
	}
	if string(val3) != "value3" {
		t.Error("Expected otherkey value to remain")
	}
}

func TestTwoLevelCacheDeleteByTags(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 测试按标签删除
	err := tlc.DeleteByTags(ctx, "tag1", "tag2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTwoLevelCacheExists(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L1中设置缓存
	err := l1.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试存在
	exists, err := tlc.Exists(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// 测试不存在
	exists, err = tlc.Exists(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected key not to exist")
	}
}

func TestTwoLevelCacheClear(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L1中设置测试数据
	err := l1.Set(ctx, "testkey1", []byte("value1"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = l1.Set(ctx, "testkey2", []byte("value2"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试清空
	err = tlc.Clear(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1已清空
	val1, err1 := l1.Get(ctx, "testkey1")
	val2, err2 := l1.Get(ctx, "testkey2")

	if err1 != ErrCacheNotFound {
		t.Error("Expected testkey1 to be cleared from L1")
	}
	if err2 != ErrCacheNotFound {
		t.Error("Expected testkey2 to be cleared from L1")
	}
	if val1 != nil || val2 != nil {
		t.Error("Expected nil values from L1")
	}
}

func TestTwoLevelCacheStats(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 执行一些操作
	l1.Set(ctx, "key1", []byte("value1"), 5*time.Minute)
	tlc.Get(ctx, "key1") // L1 hit
	tlc.Get(ctx, "key2") // miss

	stats := tlc.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected Hits to be 1, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
	if stats.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", stats.Total)
	}
	if stats.HitRate != 50.0 {
		t.Errorf("Expected HitRate to be 50.0, got %f", stats.HitRate)
	}
}

func TestTwoLevelCacheStatsDetailed(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 执行一些操作
	l1.Set(ctx, "key1", []byte("value1"), 5*time.Minute)
	tlc.Get(ctx, "key1") // L1 hit
	tlc.Get(ctx, "key2") // miss

	stats := tlc.StatsDetailed()
	if stats.L1Hits != 1 {
		t.Errorf("Expected L1Hits to be 1, got %d", stats.L1Hits)
	}
	if stats.L2Hits != 0 {
		t.Errorf("Expected L2Hits to be 0, got %d", stats.L2Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
	if stats.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", stats.Total)
	}
	if stats.L1HitRate != 50.0 {
		t.Errorf("Expected L1HitRate to be 50.0, got %f", stats.L1HitRate)
	}
	if stats.L2HitRate != 0.0 {
		t.Errorf("Expected L2HitRate to be 0.0, got %f", stats.L2HitRate)
	}
	if stats.TotalHitRate != 50.0 {
		t.Errorf("Expected TotalHitRate to be 50.0, got %f", stats.TotalHitRate)
	}
}

func TestTwoLevelCacheGetL1(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 在L1中设置缓存
	err := l1.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 测试直接从L1获取
	val, err := tlc.GetL1(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", string(val))
	}

	// 测试L1中不存在的key
	val, err = tlc.GetL1(ctx, "nonexistent")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value")
	}
}

func TestTwoLevelCacheGetL2(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 测试直接从L2获取（模拟版本总是返回未找到）
	val, err := tlc.GetL2(ctx, "testkey")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value")
	}
}

func TestTwoLevelCacheSetL1Only(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 测试只设置L1
	err := tlc.SetL1Only(ctx, "testkey", []byte("testvalue"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1中有缓存
	val, err := l1.Get(ctx, "testkey")
	if err != nil {
		t.Errorf("Expected no error from L1, got %v", err)
	}
	if string(val) != "testvalue" {
		t.Errorf("Expected 'testvalue' from L1, got %s", string(val))
	}
}

func TestTwoLevelCacheSetL2Only(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)
	ctx := context.Background()

	// 测试只设置L2
	err := tlc.SetL2Only(ctx, "testkey", []byte("testvalue"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证L1中没有缓存
	val, err := l1.Get(ctx, "testkey")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound from L1, got %v", err)
	}
	if val != nil {
		t.Error("Expected nil value from L1")
	}
}

func TestTwoLevelCacheUpdateStats(t *testing.T) {
	l1 := NewMemoryCache(100)
	defer l1.Close()
	l2 := NewRedisCache(&RedisCacheConfig{Addr: "localhost:6379"})
	tlc := NewTwoLevelCache(l1, l2, nil)

	// 手动设置统计数据
	tlc.stats.L1Hits = 5
	tlc.stats.L2Hits = 3
	tlc.stats.Misses = 2
	tlc.stats.Total = 10

	// 调用updateStats
	tlc.updateStats()

	// 验证统计数据更新
	if tlc.stats.L1HitRate != 50.0 {
		t.Errorf("Expected L1HitRate to be 50.0, got %f", tlc.stats.L1HitRate)
	}
	if tlc.stats.L2HitRate != 30.0 {
		t.Errorf("Expected L2HitRate to be 30.0, got %f", tlc.stats.L2HitRate)
	}
	if tlc.stats.TotalHitRate != 80.0 {
		t.Errorf("Expected TotalHitRate to be 80.0, got %f", tlc.stats.TotalHitRate)
	}

	// 测试Total为0的情况
	tlc.stats.Total = 0
	tlc.updateStats()

	if tlc.stats.L1HitRate != 0.0 {
		t.Errorf("Expected L1HitRate to be 0.0 when Total is 0, got %f", tlc.stats.L1HitRate)
	}
	if tlc.stats.L2HitRate != 0.0 {
		t.Errorf("Expected L2HitRate to be 0.0 when Total is 0, got %f", tlc.stats.L2HitRate)
	}
	if tlc.stats.TotalHitRate != 0.0 {
		t.Errorf("Expected TotalHitRate to be 0.0 when Total is 0, got %f", tlc.stats.TotalHitRate)
	}
}
