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
	MaxSize       int           // Maximum cache count (memory cache)
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
		MaxSize:       10000,
		EnableJitter:  false,
		TTLJitter:     0,
		EmptyCacheTTL: 30 * time.Second,
	}
}

// CacheManager cache manager
type CacheManager struct {
	cache  Cache
	config *CacheConfig
}

// NewCacheManager create cache manager
func NewCacheManager(cache Cache, config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}
	return &CacheManager{
		cache:  cache,
		config: config,
	}
}

// Get get cache
func (cm *CacheManager) Get(ctx context.Context, key string) ([]byte, error) {
	if !cm.config.Enabled {
		return nil, ErrCacheDisabled
	}
	fullKey := cm.config.KeyPrefix + key
	return cm.cache.Get(ctx, fullKey)
}

// Set set cache
func (cm *CacheManager) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if !cm.config.Enabled {
		return nil
	}

	fullKey := cm.config.KeyPrefix + key

	// Apply TTL jitter
	if cm.config.EnableJitter && cm.config.TTLJitter > 0 {
		ttl = applyJitter(ttl, cm.config.TTLJitter)
	}

	return cm.cache.Set(ctx, fullKey, value, ttl)
}

// Delete delete cache
func (cm *CacheManager) Delete(ctx context.Context, key string) error {
	if !cm.config.Enabled {
		return nil
	}
	fullKey := cm.config.KeyPrefix + key
	return cm.cache.Delete(ctx, fullKey)
}

// DeleteByPattern delete cache by pattern
func (cm *CacheManager) DeleteByPattern(ctx context.Context, pattern string) error {
	if !cm.config.Enabled {
		return nil
	}
	fullPattern := cm.config.KeyPrefix + pattern
	return cm.cache.DeleteByPattern(ctx, fullPattern)
}

// DeleteByTags delete cache by tags
func (cm *CacheManager) DeleteByTags(ctx context.Context, tags ...string) error {
	if !cm.config.Enabled {
		return nil
	}
	return cm.cache.DeleteByTags(ctx, tags...)
}

// Exists check if cache exists
func (cm *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	if !cm.config.Enabled {
		return false, nil
	}
	fullKey := cm.config.KeyPrefix + key
	return cm.cache.Exists(ctx, fullKey)
}

// Clear clear all cache
func (cm *CacheManager) Clear(ctx context.Context) error {
	if !cm.config.Enabled {
		return nil
	}
	return cm.cache.Clear(ctx)
}

// Stats get cache statistics
func (cm *CacheManager) Stats() *CacheStats {
	return cm.cache.Stats()
}

// InvalidateTags invalidate tags
func (cm *CacheManager) InvalidateTags(ctx context.Context, tags ...string) error {
	return cm.DeleteByTags(ctx, tags...)
}

// Warm cache warm-up
func (cm *CacheManager) Warm(ctx context.Context, warmFunc func(ctx context.Context) error) error {
	return warmFunc(ctx)
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