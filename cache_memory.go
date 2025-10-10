package pie

import (
	"context"
	"regexp"
	"sync"
	"time"
)

// MemoryCache in-memory cache implementation
type MemoryCache struct {
	mu          sync.RWMutex
	items       map[string]*memoryCacheItem
	tags        map[string][]string // tag -> keys mapping
	maxSize     int
	stats       *CacheStats
	stopCleanup chan struct{}
}

// memoryCacheItem in-memory cache item
type memoryCacheItem struct {
	value      []byte
	expiration time.Time
	tags       []string
}

// NewMemoryCache create in-memory cache
func NewMemoryCache(maxSize int) *MemoryCache {
	mc := &MemoryCache{
		items:       make(map[string]*memoryCacheItem),
		tags:        make(map[string][]string),
		maxSize:     maxSize,
		stats:       &CacheStats{},
		stopCleanup: make(chan struct{}),
	}

	// Start goroutine to clean up expired items
	go mc.cleanupExpired()

	return mc
}

// Get get cache
func (mc *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	mc.stats.Total++

	item, exists := mc.items[key]
	if !exists {
		mc.stats.Misses++
		return nil, ErrCacheNotFound
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		mc.stats.Misses++
		return nil, ErrCacheExpired
	}

	mc.stats.Hits++
	mc.stats.HitRate = float64(mc.stats.Hits) / float64(mc.stats.Total) * 100

	return item.value, nil
}

// Set set cache
func (mc *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Check if maximum size is exceeded
	if len(mc.items) >= mc.maxSize && mc.items[key] == nil {
		// Evict oldest item
		mc.evictOldest()
	}

	mc.items[key] = &memoryCacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
		tags:       nil,
	}

	mc.stats.Size++
	mc.stats.Keys = int64(len(mc.items))

	return nil
}

// SetWithTags set cache with tags
func (mc *MemoryCache) SetWithTags(ctx context.Context, key string, value []byte, ttl time.Duration, tags []string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Check if maximum size is exceeded
	if len(mc.items) >= mc.maxSize && mc.items[key] == nil {
		mc.evictOldest()
	}

	mc.items[key] = &memoryCacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
		tags:       tags,
	}

	// Update tag mapping
	for _, tag := range tags {
		mc.tags[tag] = append(mc.tags[tag], key)
	}

	mc.stats.Size++
	mc.stats.Keys = int64(len(mc.items))

	return nil
}

// Delete delete cache
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item := mc.items[key]
	if item != nil {
		// Delete tag mapping
		for _, tag := range item.tags {
			mc.removeKeyFromTag(tag, key)
		}
	}

	delete(mc.items, key)
	mc.stats.Keys = int64(len(mc.items))

	return nil
}

// DeleteByPattern delete cache by pattern
func (mc *MemoryCache) DeleteByPattern(ctx context.Context, pattern string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	keysToDelete := make([]string, 0)
	for key := range mc.items {
		if re.MatchString(key) {
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
	}

	mc.stats.Keys = int64(len(mc.items))

	return nil
}

// DeleteByTags delete cache by tags
func (mc *MemoryCache) DeleteByTags(ctx context.Context, tags ...string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	keysToDelete := make(map[string]bool)

	for _, tag := range tags {
		if keys, exists := mc.tags[tag]; exists {
			for _, key := range keys {
				keysToDelete[key] = true
			}
			delete(mc.tags, tag)
		}
	}

	for key := range keysToDelete {
		delete(mc.items, key)
	}

	mc.stats.Keys = int64(len(mc.items))

	return nil
}

// Exists check if cache exists
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.items[key]
	if !exists {
		return false, nil
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		return false, nil
	}

	return true, nil
}

// Clear clear all cache
func (mc *MemoryCache) Clear(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.items = make(map[string]*memoryCacheItem)
	mc.tags = make(map[string][]string)
	mc.stats.Keys = 0
	mc.stats.Size = 0

	return nil
}

// Stats get cache statistics
func (mc *MemoryCache) Stats() *CacheStats {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	stats := *mc.stats
	return &stats
}

// Close close cache
func (mc *MemoryCache) Close() {
	close(mc.stopCleanup)
}

// cleanupExpired clean up expired items
func (mc *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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

		case <-mc.stopCleanup:
			return
		}
	}
}

// evictOldest evict oldest item (LRU)
func (mc *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range mc.items {
		if oldestKey == "" || item.expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.expiration
		}
	}

	if oldestKey != "" {
		item := mc.items[oldestKey]
		if item != nil {
			for _, tag := range item.tags {
				mc.removeKeyFromTag(tag, oldestKey)
			}
		}
		delete(mc.items, oldestKey)
		mc.stats.EvictedKeys++
	}
}

// removeKeyFromTag remove key from tag
func (mc *MemoryCache) removeKeyFromTag(tag, key string) {
	keys := mc.tags[tag]
	newKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}

	if len(newKeys) == 0 {
		delete(mc.tags, tag)
	} else {
		mc.tags[tag] = newKeys
	}
}