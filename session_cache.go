package pie

import (
	"context"
	"encoding/json"
	"time"
)

// SessionCacheConfig Session cache configuration
type SessionCacheConfig struct {
	Enabled    bool          // Whether to enable cache
	TTL        time.Duration // Cache TTL
	Tags       []string      // Cache tags
	UseJitter  bool          // Whether to use TTL jitter
	CacheEmpty bool          // Whether to cache empty results
}

// Cache enable cache (use default TTL)
func (s *Session[T]) Cache(ttl ...time.Duration) *Session[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}

	s.cacheConfig.Enabled = true

	if len(ttl) > 0 {
		s.cacheConfig.TTL = ttl[0]
	} else if s.engine.cacheManager != nil {
		s.cacheConfig.TTL = s.engine.cacheManager.config.DefaultTTL
	} else {
		s.cacheConfig.TTL = 5 * time.Minute
	}

	return s
}

// NoCache disable cache
func (s *Session[T]) NoCache() *Session[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}
	s.cacheConfig.Enabled = false
	return s
}

// CacheWithTags cache with tags
func (s *Session[T]) CacheWithTags(tags ...string) *Session[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{
			Enabled: true,
		}
	}

	if s.cacheConfig.TTL == 0 && s.engine.cacheManager != nil {
		s.cacheConfig.TTL = s.engine.cacheManager.config.DefaultTTL
	}

	s.cacheConfig.Tags = tags
	s.cacheConfig.Enabled = true

	return s
}

// CacheWithJitter use TTL jitter
func (s *Session[T]) CacheWithJitter(ttl, jitter time.Duration) *Session[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}

	s.cacheConfig.Enabled = true
	s.cacheConfig.TTL = ttl
	s.cacheConfig.UseJitter = true

	return s
}

// CacheEmpty cache empty results (anti-penetration)
func (s *Session[T]) CacheEmpty(ttl time.Duration) *Session[T] {
	if s.cacheConfig == nil {
		s.cacheConfig = &SessionCacheConfig{}
	}

	s.cacheConfig.Enabled = true
	s.cacheConfig.CacheEmpty = true

	if ttl > 0 {
		s.cacheConfig.TTL = ttl
	} else if s.engine.cacheManager != nil {
		s.cacheConfig.TTL = s.engine.cacheManager.config.EmptyCacheTTL
	} else {
		s.cacheConfig.TTL = 30 * time.Second
	}

	return s
}

// getFromCache get result from cache
func (s *Session[T]) getFromCache(ctx context.Context, key string) ([]T, bool, error) {
	if s.engine.cacheManager == nil || s.cacheConfig == nil || !s.cacheConfig.Enabled {
		return nil, false, nil
	}

	data, err := s.engine.cacheManager.Get(ctx, key)
	if err != nil {
		if err == ErrCacheNotFound || err == ErrCacheExpired {
			return nil, false, nil
		}
		return nil, false, err
	}

	var results []T
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, false, err
	}

	return results, true, nil
}

// setToCache set cache
func (s *Session[T]) setToCache(ctx context.Context, key string, results []T) error {
	if s.engine.cacheManager == nil || s.cacheConfig == nil || !s.cacheConfig.Enabled {
		return nil
	}

	// If result is empty and not caching empty results, skip
	if len(results) == 0 && !s.cacheConfig.CacheEmpty {
		return nil
	}

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	ttl := s.cacheConfig.TTL
	if s.cacheConfig.UseJitter && s.engine.cacheManager.config.EnableJitter {
		ttl = applyJitter(ttl, s.engine.cacheManager.config.TTLJitter)
	}

	// If there are tags, use the method with tags
	if len(s.cacheConfig.Tags) > 0 {
		// 对于链式缓存，我们需要遍历所有缓存实例
		for _, cache := range s.engine.cacheManager.GetCaches() {
			if ristrettoCache, ok := cache.(*RistrettoCache); ok {
				ristrettoCache.SetWithTags(ctx, key, data, ttl, s.cacheConfig.Tags)
			} else if redisCache, ok := cache.(*RedisCache); ok {
				redisCache.SetWithTags(ctx, key, data, ttl, s.cacheConfig.Tags)
			} else {
				// 对于其他缓存类型，使用普通的 Set 方法
				cache.Set(ctx, key, data, ttl)
			}
		}
		return nil
	}

	return s.engine.cacheManager.Set(ctx, key, data, ttl)
}

// invalidateCache invalidate cache
func (s *Session[T]) invalidateCache(ctx context.Context) error {
	if s.engine.cacheManager == nil {
		return nil
	}

	// Invalidate cache of entire collection
	collectionName := s.collection.Name()
	keyGen := NewCacheKeyGenerator(s.engine.cacheManager.config.KeyPrefix)
	pattern := keyGen.GenerateCollectionPattern(collectionName)

	return s.engine.cacheManager.DeleteByPattern(ctx, pattern)
}

// generateCacheKey generate cache key
func (s *Session[T]) generateCacheKey() string {
	if s.engine.cacheManager == nil {
		return ""
	}

	collectionName := s.collection.Name()
	keyGen := NewCacheKeyGenerator(s.engine.cacheManager.config.KeyPrefix)

	return keyGen.GenerateQueryKey(collectionName, s.query.filter, s.options)
}
