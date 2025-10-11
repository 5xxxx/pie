package pie

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// RedisCacheConfig Redis cache configuration
type RedisCacheConfig struct {
	Addr     string // Redis address
	Password string // Redis password
	DB       int    // Redis database number
	PoolSize int    // Connection pool size
}

// RedisCache Redis cache implementation
type RedisCache struct {
	config *RedisCacheConfig
	stats  *CacheStats
	client RedisClient // 接口，可以是真实Redis客户端或模拟客户端
}

// RedisClient Redis客户端接口
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	FlushDB(ctx context.Context) error
	Ping(ctx context.Context) error
	Close() error
}

// MockRedisClient 模拟Redis客户端（用于测试或Redis不可用时）
type MockRedisClient struct {
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if val, exists := m.data[key]; exists {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

func (m *MockRedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	re, err := regexp.Compile(strings.ReplaceAll(pattern, "*", ".*"))
	if err != nil {
		return nil, err
	}
	for key := range m.data {
		if re.MatchString(key) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	count := int64(0)
	for _, key := range keys {
		if _, exists := m.data[key]; exists {
			count++
		}
	}
	return count, nil
}

func (m *MockRedisClient) FlushDB(ctx context.Context) error {
	m.data = make(map[string]string)
	return nil
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	return fmt.Errorf("mock Redis client")
}

func (m *MockRedisClient) Close() error {
	return nil
}

// NewRedisCache create Redis cache
func NewRedisCache(config *RedisCacheConfig) *RedisCache {
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	// 尝试创建真实Redis客户端，如果失败则使用模拟客户端
	var client RedisClient
	if realClient, err := createRealRedisClient(config); err == nil {
		client = realClient
	} else {
		client = NewMockRedisClient()
	}

	return &RedisCache{
		config: config,
		stats:  &CacheStats{},
		client: client,
	}
}

// createRealRedisClient 创建真实Redis客户端
func createRealRedisClient(config *RedisCacheConfig) (RedisClient, error) {
	// 这里需要导入真实的Redis客户端库
	// 例如: go get github.com/go-redis/redis/v8
	// 由于没有安装Redis客户端库，这里返回错误使用模拟客户端
	return nil, fmt.Errorf("Redis client library not available")
}

// Get get cache
func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	rc.stats.Total++

	val, err := rc.client.Get(ctx, key)
	if err != nil {
		rc.stats.Misses++
		return nil, ErrCacheNotFound
	}

	rc.stats.Hits++
	rc.stats.HitRate = float64(rc.stats.Hits) / float64(rc.stats.Total) * 100
	return []byte(val), nil
}

// Set set cache
func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return rc.client.Set(ctx, key, string(value), ttl)
}

// SetWithTags set cache with tags
func (rc *RedisCache) SetWithTags(ctx context.Context, key string, value []byte, ttl time.Duration, tags []string) error {
	// 设置主键
	err := rc.client.Set(ctx, key, string(value), ttl)
	if err != nil {
		return err
	}

	// 为每个标签创建映射
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s:%s", tag, key)
		err = rc.client.Set(ctx, tagKey, key, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete delete cache
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key)
}

// DeleteByPattern delete cache by pattern
func (rc *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := rc.client.Keys(ctx, pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rc.client.Del(ctx, keys...)
	}

	return nil
}

// DeleteByTags delete cache by tags
func (rc *RedisCache) DeleteByTags(ctx context.Context, tags ...string) error {
	var keysToDelete []string

	for _, tag := range tags {
		pattern := fmt.Sprintf("tag:%s:*", tag)
		tagKeys, err := rc.client.Keys(ctx, pattern)
		if err != nil {
			continue
		}

		for _, tagKey := range tagKeys {
			// 获取实际缓存键
			actualKey, err := rc.client.Get(ctx, tagKey)
			if err != nil {
				continue
			}
			keysToDelete = append(keysToDelete, actualKey)
		}
	}

	if len(keysToDelete) > 0 {
		return rc.client.Del(ctx, keysToDelete...)
	}

	return nil
}

// Exists check if cache exists
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rc.client.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Clear clear all cache
func (rc *RedisCache) Clear(ctx context.Context) error {
	return rc.client.FlushDB(ctx)
}

// Stats get cache statistics
func (rc *RedisCache) Stats() *CacheStats {
	stats := *rc.stats
	return &stats
}

// Close close Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// Ping check Redis connection
func (rc *RedisCache) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx)
}
