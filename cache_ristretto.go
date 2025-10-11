package pie

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

// RistrettoCacheConfig Ristretto cache configuration
type RistrettoCacheConfig struct {
	NumCounters int64 // 约为最大条目数的10倍
	MaxCost     int64 // 最大缓存大小(字节)
	BufferItems int64 // Get buffer大小
}

// DefaultRistrettoCacheConfig creates default Ristretto cache configuration
func DefaultRistrettoCacheConfig() *RistrettoCacheConfig {
	return &RistrettoCacheConfig{
		NumCounters: 100000,            // 10万条目
		MaxCost:     100 * 1024 * 1024, // 100MB
		BufferItems: 64,
	}
}

// RistrettoCache Ristretto cache implementation
type RistrettoCache struct {
	cache  *ristretto.Cache[string, []byte]
	config *RistrettoCacheConfig
	stats  *CacheStats
	mu     sync.RWMutex
	tags   map[string][]string // tag -> keys mapping
}

// NewRistrettoCache create Ristretto cache
func NewRistrettoCache(config *RistrettoCacheConfig) (*RistrettoCache, error) {
	if config == nil {
		config = DefaultRistrettoCacheConfig()
	}

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: config.NumCounters,
		MaxCost:     config.MaxCost,
		BufferItems: config.BufferItems,
		OnEvict: func(item *ristretto.Item[[]byte]) {
			// 可以在这里处理缓存淘汰事件
		},
	})
	if err != nil {
		return nil, err
	}

	return &RistrettoCache{
		cache:  cache,
		config: config,
		stats:  &CacheStats{},
		tags:   make(map[string][]string),
	}, nil
}

// Get get cache
func (rc *RistrettoCache) Get(ctx context.Context, key string) ([]byte, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	rc.stats.Total++

	value, found := rc.cache.Get(key)
	if !found {
		rc.stats.Misses++
		return nil, ErrCacheNotFound
	}

	rc.stats.Hits++
	rc.stats.HitRate = float64(rc.stats.Hits) / float64(rc.stats.Total) * 100

	return value, nil
}

// Set set cache
func (rc *RistrettoCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// 计算成本 (字节数)
	cost := int64(len(value))

	// 设置缓存
	success := rc.cache.SetWithTTL(key, value, cost, ttl)
	if !success {
		return fmt.Errorf("failed to set cache key: %s", key)
	}

	// 等待缓存设置完成
	rc.cache.Wait()

	rc.stats.Size += cost
	rc.stats.Keys++

	return nil
}

// SetWithTags set cache with tags
func (rc *RistrettoCache) SetWithTags(ctx context.Context, key string, value []byte, ttl time.Duration, tags []string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// 计算成本
	cost := int64(len(value))

	// 设置缓存
	success := rc.cache.SetWithTTL(key, value, cost, ttl)
	if !success {
		return fmt.Errorf("failed to set cache key with tags: %s", key)
	}

	// 等待缓存设置完成
	rc.cache.Wait()

	// 更新标签映射
	for _, tag := range tags {
		rc.tags[tag] = append(rc.tags[tag], key)
	}

	rc.stats.Size += cost
	rc.stats.Keys++

	return nil
}

// Delete delete cache
func (rc *RistrettoCache) Delete(ctx context.Context, key string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache.Del(key)
	rc.stats.Keys--

	return nil
}

// DeleteByPattern delete cache by pattern
func (rc *RistrettoCache) DeleteByPattern(ctx context.Context, pattern string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// 由于 Ristretto 不提供 Keys() 方法，这里简化实现
	// 在实际应用中，可能需要维护一个键列表或使用其他方式
	// 这里简化处理，实际应用中需要根据具体需求实现
	// 可以考虑维护一个键列表或使用其他缓存实现

	// 暂时不实现模式删除，因为 Ristretto 不提供 Keys() 方法
	// 如果需要此功能，建议使用 Redis 缓存或维护键列表

	return nil
}

// DeleteByTags delete cache by tags
func (rc *RistrettoCache) DeleteByTags(ctx context.Context, tags ...string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	keysToDelete := make(map[string]bool)

	for _, tag := range tags {
		if keys, exists := rc.tags[tag]; exists {
			for _, key := range keys {
				keysToDelete[key] = true
			}
			delete(rc.tags, tag)
		}
	}

	for key := range keysToDelete {
		rc.cache.Del(key)
		rc.stats.Keys--
	}

	return nil
}

// Exists check if cache exists
func (rc *RistrettoCache) Exists(ctx context.Context, key string) (bool, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	_, found := rc.cache.Get(key)
	return found, nil
}

// Clear clear all cache
func (rc *RistrettoCache) Clear(ctx context.Context) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache.Clear()
	rc.tags = make(map[string][]string)
	rc.stats.Keys = 0
	rc.stats.Size = 0

	return nil
}

// Stats get cache statistics
func (rc *RistrettoCache) Stats() *CacheStats {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	stats := *rc.stats
	// Ristretto 不提供 Keys() 方法，使用统计信息中的 Keys
	return &stats
}

// Close close cache
func (rc *RistrettoCache) Close() {
	rc.cache.Close()
}
