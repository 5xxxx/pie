package pie

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCacheConfig Redis cache configuration
type RedisCacheConfig struct {
	Addr     string // Redis address
	Password string // Redis password
	DB       int    // Redis database number
	PoolSize int    // Connection pool size
}

// DefaultRedisCacheConfig creates default Redis cache configuration
func DefaultRedisCacheConfig() *RedisCacheConfig {
	return &RedisCacheConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	}
}

// RedisCache Redis cache implementation
type RedisCache struct {
	client *redis.Client
	config *RedisCacheConfig
	stats  *CacheStats
	mu     sync.RWMutex
}

// NewRedisCache create Redis cache
func NewRedisCache(config *RedisCacheConfig) (*RedisCache, error) {
	if config == nil {
		config = DefaultRedisCacheConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		config: config,
		stats:  &CacheStats{},
	}, nil
}

// Get get cache
func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	rc.mu.Lock()
	rc.stats.Total++
	rc.mu.Unlock()

	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			rc.mu.Lock()
			rc.stats.Misses++
			rc.mu.Unlock()
			return nil, ErrCacheNotFound
		}
		return nil, err
	}

	rc.mu.Lock()
	rc.stats.Hits++
	rc.stats.HitRate = float64(rc.stats.Hits) / float64(rc.stats.Total) * 100
	rc.mu.Unlock()

	return []byte(val), nil
}

// Set set cache
func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return rc.client.Set(ctx, key, value, ttl).Err()
}

// SetWithTags set cache with tags
func (rc *RedisCache) SetWithTags(ctx context.Context, key string, value []byte, ttl time.Duration, tags []string) error {
	// 设置主键
	err := rc.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	// 为每个标签创建映射
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s:%s", tag, key)
		err = rc.client.Set(ctx, tagKey, key, ttl).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete delete cache
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// DeleteByPattern delete cache by pattern
func (rc *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rc.client.Del(ctx, keys...).Err()
	}

	return nil
}

// DeleteByTags delete cache by tags
func (rc *RedisCache) DeleteByTags(ctx context.Context, tags ...string) error {
	var keysToDelete []string

	for _, tag := range tags {
		pattern := fmt.Sprintf("tag:%s:*", tag)
		tagKeys, err := rc.client.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}

		for _, tagKey := range tagKeys {
			// 获取实际缓存键
			actualKey, err := rc.client.Get(ctx, tagKey).Result()
			if err != nil {
				continue
			}
			keysToDelete = append(keysToDelete, actualKey)
		}
	}

	if len(keysToDelete) > 0 {
		return rc.client.Del(ctx, keysToDelete...).Err()
	}

	return nil
}

// Exists check if cache exists
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Clear clear all cache
func (rc *RedisCache) Clear(ctx context.Context) error {
	return rc.client.FlushDB(ctx).Err()
}

// Stats get cache statistics
func (rc *RedisCache) Stats() *CacheStats {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	stats := *rc.stats
	return &stats
}

// Close close Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// Ping check Redis connection
func (rc *RedisCache) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}

// MockRedisCache mock Redis cache for testing
type MockRedisCache struct {
	data  map[string]string
	stats *CacheStats
	mu    sync.RWMutex
}

// NewMockRedisCache create mock Redis cache
func NewMockRedisCache() *MockRedisCache {
	return &MockRedisCache{
		data:  make(map[string]string),
		stats: &CacheStats{},
	}
}

// Get get cache
func (m *MockRedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.stats.Total++

	if val, exists := m.data[key]; exists {
		m.stats.Hits++
		m.stats.HitRate = float64(m.stats.Hits) / float64(m.stats.Total) * 100
		return []byte(val), nil
	}

	m.stats.Misses++
	return nil, ErrCacheNotFound
}

// Set set cache
func (m *MockRedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = string(value)
	m.stats.Keys++
	return nil
}

// SetWithTags set cache with tags
func (m *MockRedisCache) SetWithTags(ctx context.Context, key string, value []byte, ttl time.Duration, tags []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = string(value)

	// 为每个标签创建映射
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s:%s", tag, key)
		m.data[tagKey] = key
	}

	m.stats.Keys++
	return nil
}

// Delete delete cache
func (m *MockRedisCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	m.stats.Keys--
	return nil
}

// DeleteByPattern delete cache by pattern
func (m *MockRedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	re, err := regexp.Compile(strings.ReplaceAll(pattern, "*", ".*"))
	if err != nil {
		return err
	}

	keysToDelete := make([]string, 0)
	for key := range m.data {
		if re.MatchString(key) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		delete(m.data, key)
		m.stats.Keys--
	}

	return nil
}

// DeleteByTags delete cache by tags
func (m *MockRedisCache) DeleteByTags(ctx context.Context, tags ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	keysToDelete := make(map[string]bool)

	for _, tag := range tags {
		pattern := fmt.Sprintf("tag:%s:.*", tag)
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		for key := range m.data {
			if re.MatchString(key) {
				// 获取实际缓存键
				if actualKey, exists := m.data[key]; exists {
					keysToDelete[actualKey] = true
				}
			}
		}
	}

	for key := range keysToDelete {
		delete(m.data, key)
		m.stats.Keys--
	}

	return nil
}

// Exists check if cache exists
func (m *MockRedisCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists, nil
}

// Clear clear all cache
func (m *MockRedisCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]string)
	m.stats.Keys = 0
	return nil
}

// Stats get cache statistics
func (m *MockRedisCache) Stats() *CacheStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := *m.stats
	return &stats
}

// Close close mock cache
func (m *MockRedisCache) Close() error {
	return nil
}
