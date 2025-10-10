package pie

import (
	"context"
	"sync"
	"time"
)

// SingleFlight anti-breakdown (prevent cache breakdown)
type SingleFlight struct {
	mu sync.Mutex
	m  map[string]*call
}

// call single call
type call struct {
	wg  sync.WaitGroup
	val []byte
	err error
}

// NewSingleFlight create SingleFlight
func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		m: make(map[string]*call),
	}
}

// Do execute function, requests with the same key are executed only once
func (sf *SingleFlight) Do(key string, fn func() ([]byte, error)) ([]byte, error) {
	sf.mu.Lock()

	if c, ok := sf.m[key]; ok {
		// Same request is already executing, wait for result
		sf.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// Create new call
	c := &call{}
	c.wg.Add(1)
	sf.m[key] = c
	sf.mu.Unlock()

	// Execute function
	c.val, c.err = fn()
	c.wg.Done()

	// Clean up
	sf.mu.Lock()
	delete(sf.m, key)
	sf.mu.Unlock()

	return c.val, c.err
}

// CacheWithSingleFlight cache with SingleFlight protection
type CacheWithSingleFlight struct {
	cache        Cache
	singleFlight *SingleFlight
}

// NewCacheWithSingleFlight create cache with SingleFlight protection
func NewCacheWithSingleFlight(cache Cache) *CacheWithSingleFlight {
	return &CacheWithSingleFlight{
		cache:        cache,
		singleFlight: NewSingleFlight(),
	}
}

// Get get cache (with SingleFlight protection)
func (csf *CacheWithSingleFlight) Get(ctx context.Context, key string) ([]byte, error) {
	return csf.singleFlight.Do(key, func() ([]byte, error) {
		return csf.cache.Get(ctx, key)
	})
}

// Set set cache
func (csf *CacheWithSingleFlight) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return csf.cache.Set(ctx, key, value, ttl)
}

// Delete delete cache
func (csf *CacheWithSingleFlight) Delete(ctx context.Context, key string) error {
	return csf.cache.Delete(ctx, key)
}

// DeleteByPattern delete cache by pattern
func (csf *CacheWithSingleFlight) DeleteByPattern(ctx context.Context, pattern string) error {
	return csf.cache.DeleteByPattern(ctx, pattern)
}

// DeleteByTags delete cache by tags
func (csf *CacheWithSingleFlight) DeleteByTags(ctx context.Context, tags ...string) error {
	return csf.cache.DeleteByTags(ctx, tags...)
}

// Exists check if cache exists
func (csf *CacheWithSingleFlight) Exists(ctx context.Context, key string) (bool, error) {
	return csf.cache.Exists(ctx, key)
}

// Clear clear all cache
func (csf *CacheWithSingleFlight) Clear(ctx context.Context) error {
	return csf.cache.Clear(ctx)
}

// Stats get cache statistics
func (csf *CacheWithSingleFlight) Stats() *CacheStats {
	return csf.cache.Stats()
}