package pie

import (
	"context"
	"time"
)

// Cache cache interface
type Cache interface {
	// Get gets cache value
	Get(ctx context.Context, key string) ([]byte, error)

	// Set sets cache value
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete deletes cache by key
	Delete(ctx context.Context, key string) error

	// DeleteByPattern deletes cache by pattern
	DeleteByPattern(ctx context.Context, pattern string) error

	// DeleteByTags deletes cache by tags
	DeleteByTags(ctx context.Context, tags ...string) error

	// Exists checks if cache exists
	Exists(ctx context.Context, key string) (bool, error)

	// Clear clears all cache
	Clear(ctx context.Context) error

	// Stats gets cache statistics
	Stats() *CacheStats
}

// CacheStats cache statistics
type CacheStats struct {
	Hits        int64   // Hit count
	Misses      int64   // Miss count
	Total       int64   // Total request count
	HitRate     float64 // Hit rate
	Size        int64   // Current cache size
	Keys        int64   // Current key count
	EvictedKeys int64   // Evicted key count
}

// CacheConfig cache configuration
type CacheConfig struct {
	Enabled       bool          // Whether caching is enabled
	DefaultTTL    time.Duration // Default TTL
	KeyPrefix     string        // Key prefix
	EnableJitter  bool          // Whether to enable TTL jitter
	TTLJitter     time.Duration // TTL jitter range (±)
	EmptyCacheTTL time.Duration // Empty result cache TTL (anti-penetration)
}

// DefaultCacheConfig creates default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled:       true,
		DefaultTTL:    5 * time.Minute,
		KeyPrefix:     "pie:",
		EnableJitter:  false,
		TTLJitter:     0,
		EmptyCacheTTL: 30 * time.Second,
	}
}

// CacheChainConfig chain cache configuration
type CacheChainConfig struct {
	Caches         []Cache // 缓存链(按优先级排序)
	EnableBackfill bool    // 是否启用回填(L2命中时写回L1)
	DefaultTTL     time.Duration
	KeyPrefix      string
}

// CacheManager cache manager
type CacheManager struct {
	caches []Cache
	config *CacheConfig
}

// NewCacheManager create cache manager
func NewCacheManager(caches []Cache, config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}
	return &CacheManager{
		caches: caches,
		config: config,
	}
}

// NewSingleCacheManager create single cache manager (backward compatibility)
func NewSingleCacheManager(cache Cache, config *CacheConfig) *CacheManager {
	return NewCacheManager([]Cache{cache}, config)
}

// Get get cache
func (cm *CacheManager) Get(ctx context.Context, key string) ([]byte, error) {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return nil, ErrCacheDisabled
	}

	fullKey := cm.config.KeyPrefix + key

	// 按顺序查找缓存
	for i, cache := range cm.caches {
		value, err := cache.Get(ctx, fullKey)
		if err == nil {
			// 找到缓存，回填到前面的层级
			if i > 0 {
				cm.backfillToPreviousCaches(ctx, fullKey, value, i)
			}
			return value, nil
		}
	}

	return nil, ErrCacheNotFound
}

// backfillToPreviousCaches 回填到前面的缓存层级
func (cm *CacheManager) backfillToPreviousCaches(ctx context.Context, key string, value []byte, foundIndex int) {
	// 回填到前面的所有缓存
	for i := 0; i < foundIndex; i++ {
		cm.caches[i].Set(ctx, key, value, cm.config.DefaultTTL)
	}
}

// Set set cache
func (cm *CacheManager) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return nil
	}

	fullKey := cm.config.KeyPrefix + key

	// Apply TTL jitter
	if cm.config.EnableJitter && cm.config.TTLJitter > 0 {
		ttl = applyJitter(ttl, cm.config.TTLJitter)
	}

	// 写入所有缓存层级
	var lastErr error
	for _, cache := range cm.caches {
		if err := cache.Set(ctx, fullKey, value, ttl); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// Delete delete cache
func (cm *CacheManager) Delete(ctx context.Context, key string) error {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return nil
	}

	fullKey := cm.config.KeyPrefix + key

	// 从所有缓存层级删除
	var lastErr error
	for _, cache := range cm.caches {
		if err := cache.Delete(ctx, fullKey); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// DeleteByPattern delete cache by pattern
func (cm *CacheManager) DeleteByPattern(ctx context.Context, pattern string) error {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return nil
	}

	fullPattern := cm.config.KeyPrefix + pattern

	// 从所有缓存层级按模式删除
	var lastErr error
	for _, cache := range cm.caches {
		if err := cache.DeleteByPattern(ctx, fullPattern); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// DeleteByTags delete cache by tags
func (cm *CacheManager) DeleteByTags(ctx context.Context, tags ...string) error {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return nil
	}

	// 从所有缓存层级按标签删除
	var lastErr error
	for _, cache := range cm.caches {
		if err := cache.DeleteByTags(ctx, tags...); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// Exists check if cache exists
func (cm *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return false, nil
	}

	fullKey := cm.config.KeyPrefix + key

	// 检查任意一个缓存是否存在
	for _, cache := range cm.caches {
		exists, err := cache.Exists(ctx, fullKey)
		if err == nil && exists {
			return true, nil
		}
	}

	return false, nil
}

// Clear clear all cache
func (cm *CacheManager) Clear(ctx context.Context) error {
	if !cm.config.Enabled || len(cm.caches) == 0 {
		return nil
	}

	// 清除所有缓存层级
	var lastErr error
	for _, cache := range cm.caches {
		if err := cache.Clear(ctx); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// Stats get cache statistics
func (cm *CacheManager) Stats() *CacheStats {
	if len(cm.caches) == 0 {
		return &CacheStats{}
	}

	// 聚合所有缓存的统计信息
	stats := &CacheStats{}
	for _, cache := range cm.caches {
		cacheStats := cache.Stats()
		stats.Hits += cacheStats.Hits
		stats.Misses += cacheStats.Misses
		stats.Total += cacheStats.Total
		stats.Size += cacheStats.Size
		stats.Keys += cacheStats.Keys
		stats.EvictedKeys += cacheStats.EvictedKeys
	}

	if stats.Total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(stats.Total) * 100
	}

	return stats
}

// InvalidateTags invalidate tags
func (cm *CacheManager) InvalidateTags(ctx context.Context, tags ...string) error {
	return cm.DeleteByTags(ctx, tags...)
}

// Warm cache warm-up
func (cm *CacheManager) Warm(ctx context.Context, warmFunc func(ctx context.Context) error) error {
	return warmFunc(ctx)
}

// AddCache add cache to chain
func (cm *CacheManager) AddCache(cache Cache) {
	cm.caches = append(cm.caches, cache)
}

// GetCaches get all caches
func (cm *CacheManager) GetCaches() []Cache {
	return cm.caches
}

// SetCaches set caches
func (cm *CacheManager) SetCaches(caches []Cache) {
	cm.caches = caches
}

// applyJitter apply TTL jitter
func applyJitter(ttl, jitter time.Duration) time.Duration {
	if jitter == 0 {
		return ttl
	}
	// Generate random offset in range [-jitter, +jitter]
	offset := time.Duration(randomInt(-int64(jitter), int64(jitter)))
	newTTL := ttl + offset
	if newTTL < 0 {
		return ttl
	}
	return newTTL
}
