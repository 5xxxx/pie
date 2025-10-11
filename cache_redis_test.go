package pie

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewRedisCache(t *testing.T) {
	config := &RedisCacheConfig{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
		PoolSize: 20,
	}

	rc := NewRedisCache(config)
	if rc == nil {
		t.Error("Expected RedisCache to be created")
	}
	if rc.config != config {
		t.Error("Expected config to be set correctly")
	}
	if rc.stats == nil {
		t.Error("Expected stats to be initialized")
	}

	// 测试默认PoolSize
	config2 := &RedisCacheConfig{
		Addr: "localhost:6379",
	}
	rc2 := NewRedisCache(config2)
	if rc2.config.PoolSize != 10 {
		t.Errorf("Expected default PoolSize to be 10, got %d", rc2.config.PoolSize)
	}
}

func TestRedisCacheIntegration(t *testing.T) {
	// 从环境变量获取Redis地址，默认为127.0.0.1:6379
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	config := &RedisCacheConfig{
		Addr: redisAddr,
		DB:   1, // 使用数据库1进行测试
	}

	rc := NewRedisCache(config)
	defer rc.Close()
	ctx := context.Background()

	// 测试连接
	err := rc.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
		return
	}

	// 清空测试数据库
	err = rc.Clear(ctx)
	if err != nil {
		t.Errorf("Failed to clear Redis: %v", err)
	}

	// 测试基本操作
	testKey := "test:integration:key"
	testValue := []byte("test value")

	// 测试设置
	err = rc.Set(ctx, testKey, testValue, 5*time.Minute)
	if err != nil {
		t.Errorf("Failed to set key: %v", err)
	}

	// 测试获取
	val, err := rc.Get(ctx, testKey)
	if err != nil {
		t.Errorf("Failed to get key: %v", err)
	}
	if string(val) != string(testValue) {
		t.Errorf("Expected '%s', got '%s'", string(testValue), string(val))
	}

	// 测试存在检查
	exists, err := rc.Exists(ctx, testKey)
	if err != nil {
		t.Errorf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// 测试删除
	err = rc.Delete(ctx, testKey)
	if err != nil {
		t.Errorf("Failed to delete key: %v", err)
	}

	// 验证删除
	val, err = rc.Get(ctx, testKey)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound after delete, got %v", err)
	}

	// 验证不存在
	exists, err = rc.Exists(ctx, testKey)
	if err != nil {
		t.Errorf("Failed to check existence after delete: %v", err)
	}
	if exists {
		t.Error("Expected key not to exist after delete")
	}
}

func TestRedisCacheSetWithTagsIntegration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	config := &RedisCacheConfig{
		Addr: redisAddr,
		DB:   1,
	}

	rc := NewRedisCache(config)
	defer rc.Close()
	ctx := context.Background()

	// 测试连接
	err := rc.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
		return
	}

	// 清空测试数据库
	err = rc.Clear(ctx)
	if err != nil {
		t.Errorf("Failed to clear Redis: %v", err)
	}

	// 测试带标签的设置
	testKey := "test:tagged:key"
	testValue := []byte("tagged value")
	tags := []string{"tag1", "tag2"}

	err = rc.SetWithTags(ctx, testKey, testValue, 5*time.Minute, tags)
	if err != nil {
		t.Errorf("Failed to set key with tags: %v", err)
	}

	// 验证设置成功
	val, err := rc.Get(ctx, testKey)
	if err != nil {
		t.Errorf("Failed to get tagged key: %v", err)
	}
	if string(val) != string(testValue) {
		t.Errorf("Expected '%s', got '%s'", string(testValue), string(val))
	}

	// 测试按标签删除
	err = rc.DeleteByTags(ctx, "tag1")
	if err != nil {
		t.Errorf("Failed to delete by tags: %v", err)
	}

	// 验证删除
	val, err = rc.Get(ctx, testKey)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound after tag delete, got %v", err)
	}
}

func TestRedisCachePatternDeleteIntegration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	config := &RedisCacheConfig{
		Addr: redisAddr,
		DB:   1,
	}

	rc := NewRedisCache(config)
	defer rc.Close()
	ctx := context.Background()

	// 测试连接
	err := rc.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
		return
	}

	// 清空测试数据库
	err = rc.Clear(ctx)
	if err != nil {
		t.Errorf("Failed to clear Redis: %v", err)
	}

	// 设置多个测试键
	keys := []string{
		"test:pattern:key1",
		"test:pattern:key2",
		"test:other:key3",
	}
	values := [][]byte{
		[]byte("value1"),
		[]byte("value2"),
		[]byte("value3"),
	}

	for i, key := range keys {
		err = rc.Set(ctx, key, values[i], 5*time.Minute)
		if err != nil {
			t.Errorf("Failed to set key %s: %v", key, err)
		}
	}

	// 测试按模式删除
	err = rc.DeleteByPattern(ctx, "test:pattern:*")
	if err != nil {
		t.Errorf("Failed to delete by pattern: %v", err)
	}

	// 验证删除结果
	_, err1 := rc.Get(ctx, keys[0])
	_, err2 := rc.Get(ctx, keys[1])
	val3, err3 := rc.Get(ctx, keys[2])

	if err1 != ErrCacheNotFound {
		t.Error("Expected key1 to be deleted by pattern")
	}
	if err2 != ErrCacheNotFound {
		t.Error("Expected key2 to be deleted by pattern")
	}
	if err3 != nil {
		t.Error("Expected key3 to remain")
	}
	if string(val3) != string(values[2]) {
		t.Errorf("Expected '%s', got '%s'", string(values[2]), string(val3))
	}
}

func TestRedisCacheTTLIntegration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	config := &RedisCacheConfig{
		Addr: redisAddr,
		DB:   1,
	}

	rc := NewRedisCache(config)
	defer rc.Close()
	ctx := context.Background()

	// 测试连接
	err := rc.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
		return
	}

	// 清空测试数据库
	err = rc.Clear(ctx)
	if err != nil {
		t.Errorf("Failed to clear Redis: %v", err)
	}

	// 测试短期TTL
	testKey := "test:ttl:key"
	testValue := []byte("ttl value")
	shortTTL := 2 * time.Second

	err = rc.Set(ctx, testKey, testValue, shortTTL)
	if err != nil {
		t.Errorf("Failed to set key with TTL: %v", err)
	}

	// 立即获取应该成功
	val, err := rc.Get(ctx, testKey)
	if err != nil {
		t.Errorf("Failed to get key immediately: %v", err)
	}
	if string(val) != string(testValue) {
		t.Errorf("Expected '%s', got '%s'", string(testValue), string(val))
	}

	// 等待TTL过期
	time.Sleep(3 * time.Second)

	// 过期后获取应该失败
	val, err = rc.Get(ctx, testKey)
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound after TTL expiry, got %v", err)
	}
}

func TestRedisCacheStatsIntegration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	config := &RedisCacheConfig{
		Addr: redisAddr,
		DB:   1,
	}

	rc := NewRedisCache(config)
	defer rc.Close()
	ctx := context.Background()

	// 测试连接
	err := rc.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
		return
	}

	// 清空测试数据库
	err = rc.Clear(ctx)
	if err != nil {
		t.Errorf("Failed to clear Redis: %v", err)
	}

	// 执行一些操作
	rc.Get(ctx, "nonexistent1") // miss
	rc.Get(ctx, "nonexistent2") // miss
	rc.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	rc.Get(ctx, "testkey") // hit

	// 验证统计数据
	stats := rc.Stats()
	if stats.Total != 4 {
		t.Errorf("Expected Total to be 4, got %d", stats.Total)
	}
	if stats.Hits != 1 {
		t.Errorf("Expected Hits to be 1, got %d", stats.Hits)
	}
	if stats.Misses != 3 {
		t.Errorf("Expected Misses to be 3, got %d", stats.Misses)
	}
	if stats.HitRate != 25.0 {
		t.Errorf("Expected HitRate to be 25.0, got %f", stats.HitRate)
	}
}

func TestRedisCacheConcurrentIntegration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	config := &RedisCacheConfig{
		Addr: redisAddr,
		DB:   1,
	}

	rc := NewRedisCache(config)
	defer rc.Close()
	ctx := context.Background()

	// 测试连接
	err := rc.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available at %s, skipping integration test: %v", redisAddr, err)
		return
	}

	// 清空测试数据库
	err = rc.Clear(ctx)
	if err != nil {
		t.Errorf("Failed to clear Redis: %v", err)
	}

	// 并发测试
	numGoroutines := 10
	numOperations := 10

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("concurrent:key_%d_%d", goroutineID, j)
				value := fmt.Sprintf("value_%d_%d", goroutineID, j)

				// 设置
				err := rc.Set(ctx, key, []byte(value), 5*time.Minute)
				if err != nil {
					t.Errorf("Goroutine %d failed to set key %s: %v", goroutineID, key, err)
					return
				}

				// 获取
				val, err := rc.Get(ctx, key)
				if err != nil {
					t.Errorf("Goroutine %d failed to get key %s: %v", goroutineID, key, err)
					return
				}
				if string(val) != value {
					t.Errorf("Goroutine %d got wrong value for key %s: expected '%s', got '%s'", goroutineID, key, value, string(val))
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证统计数据
	stats := rc.Stats()
	if stats.Total == 0 {
		t.Error("Expected some operations to be recorded")
	}
}
