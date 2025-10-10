package pie

import (
	"context"
	"time"
)

// TwoLevelCacheConfig two-level cache configuration
type TwoLevelCacheConfig struct {
	L1TTL       time.Duration // L1 (memory) TTL
	L2TTL       time.Duration // L2 (Redis) TTL
	SyncOnWrite bool          // Whether to sync to both levels on write
}

// DefaultTwoLevelCacheConfig default two-level cache configuration
func DefaultTwoLevelCacheConfig() *TwoLevelCacheConfig {
	return &TwoLevelCacheConfig{
		L1TTL:       1 * time.Minute,
		L2TTL:       10 * time.Minute,
		SyncOnWrite: true,
	}
}

// TwoLevelCache two-level cache (L1: memory + L2: Redis)
type TwoLevelCache struct {
	l1     Cache // L1: memory cache
	l2     Cache // L2: Redis cache
	config *TwoLevelCacheConfig
	stats  *TwoLevelCacheStats
}

// TwoLevelCacheStats two-level cache statistics
type TwoLevelCacheStats struct {
	L1Hits       int64   // L1 hit count
	L2Hits       int64   // L2 hit count
	Misses       int64   // Miss count
	Total        int64   // Total request count
	L1HitRate    float64 // L1 hit rate
	L2HitRate    float64 // L2 hit rate
	TotalHitRate float64 // Total hit rate
}

// NewTwoLevelCache create two-level cache
func NewTwoLevelCache(l1, l2 Cache, config *TwoLevelCacheConfig) *TwoLevelCache {
	if config == nil {
		config = DefaultTwoLevelCacheConfig()
	}

	return &TwoLevelCache{
		l1:     l1,
		l2:     l2,
		config: config,
		stats:  &TwoLevelCacheStats{},
	}
}

// Get get cache (check L1 first, if miss then check L2, if L2 hit then write back to L1)
func (tlc *TwoLevelCache) Get(ctx context.Context, key string) ([]byte, error) {
	tlc.stats.Total++

	// Check L1 first
	val, err := tlc.l1.Get(ctx, key)
	if err == nil {
		tlc.stats.L1Hits++
		tlc.updateStats()
		return val, nil
	}

	// L1 miss, check L2
	val, err = tlc.l2.Get(ctx, key)
	if err == nil {
		tlc.stats.L2Hits++
		tlc.updateStats()

		// Write back to L1
		_ = tlc.l1.Set(ctx, key, val, tlc.config.L1TTL)

		return val, nil
	}

	// L2 also miss
	tlc.stats.Misses++
	tlc.updateStats()

	return nil, ErrCacheNotFound
}

// Set set cache (write to both L1 and L2, or only L2)
func (tlc *TwoLevelCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// Write to L2
	if err := tlc.l2.Set(ctx, key, value, tlc.config.L2TTL); err != nil {
		return err
	}

	// If sync on write is enabled, also write to L1
	if tlc.config.SyncOnWrite {
		_ = tlc.l1.Set(ctx, key, value, tlc.config.L1TTL)
	}

	return nil
}

// Delete delete cache (delete from both L1 and L2)
func (tlc *TwoLevelCache) Delete(ctx context.Context, key string) error {
	_ = tlc.l1.Delete(ctx, key)
	return tlc.l2.Delete(ctx, key)
}

// DeleteByPattern delete cache by pattern (delete from both L1 and L2)
func (tlc *TwoLevelCache) DeleteByPattern(ctx context.Context, pattern string) error {
	_ = tlc.l1.DeleteByPattern(ctx, pattern)
	return tlc.l2.DeleteByPattern(ctx, pattern)
}

// DeleteByTags delete cache by tags (delete from both L1 and L2)
func (tlc *TwoLevelCache) DeleteByTags(ctx context.Context, tags ...string) error {
	_ = tlc.l1.DeleteByTags(ctx, tags...)
	return tlc.l2.DeleteByTags(ctx, tags...)
}

// Exists check if cache exists (check L1 first, then L2)
func (tlc *TwoLevelCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, _ := tlc.l1.Exists(ctx, key)
	if exists {
		return true, nil
	}

	return tlc.l2.Exists(ctx, key)
}

// Clear clear all cache (clear L1 and L2)
func (tlc *TwoLevelCache) Clear(ctx context.Context) error {
	_ = tlc.l1.Clear(ctx)
	return tlc.l2.Clear(ctx)
}

// Stats get cache statistics
func (tlc *TwoLevelCache) Stats() *CacheStats {
	// Return merged statistics
	return &CacheStats{
		Hits:    tlc.stats.L1Hits + tlc.stats.L2Hits,
		Misses:  tlc.stats.Misses,
		Total:   tlc.stats.Total,
		HitRate: tlc.stats.TotalHitRate,
	}
}

// StatsDetailed get detailed two-level cache statistics
func (tlc *TwoLevelCache) StatsDetailed() *TwoLevelCacheStats {
	stats := *tlc.stats
	return &stats
}

// GetL1 get from L1 only
func (tlc *TwoLevelCache) GetL1(ctx context.Context, key string) ([]byte, error) {
	return tlc.l1.Get(ctx, key)
}

// GetL2 get from L2 only
func (tlc *TwoLevelCache) GetL2(ctx context.Context, key string) ([]byte, error) {
	return tlc.l2.Get(ctx, key)
}

// SetL1Only write to L1 only
func (tlc *TwoLevelCache) SetL1Only(ctx context.Context, key string, value []byte) error {
	return tlc.l1.Set(ctx, key, value, tlc.config.L1TTL)
}

// SetL2Only write to L2 only
func (tlc *TwoLevelCache) SetL2Only(ctx context.Context, key string, value []byte) error {
	return tlc.l2.Set(ctx, key, value, tlc.config.L2TTL)
}

// updateStats update statistics
func (tlc *TwoLevelCache) updateStats() {
	if tlc.stats.Total > 0 {
		tlc.stats.L1HitRate = float64(tlc.stats.L1Hits) / float64(tlc.stats.Total) * 100
		tlc.stats.L2HitRate = float64(tlc.stats.L2Hits) / float64(tlc.stats.Total) * 100
		totalHits := tlc.stats.L1Hits + tlc.stats.L2Hits
		tlc.stats.TotalHitRate = float64(totalHits) / float64(tlc.stats.Total) * 100
	}
}