package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/5xxxx/pie"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Product product model
type Product struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Price     float64       `bson:"price"`
	Category  string        `bson:"category"`
	CreatedAt time.Time     `bson:"created_at"`
}

// CustomCache 自定义缓存实现示例
type CustomCache struct {
	data  map[string][]byte
	stats *pie.CacheStats
	mu    sync.RWMutex
}

func NewCustomCache() *CustomCache {
	return &CustomCache{
		data:  make(map[string][]byte),
		stats: &pie.CacheStats{},
	}
}

func (c *CustomCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.stats.Total++
	if val, exists := c.data[key]; exists {
		c.stats.Hits++
		c.stats.HitRate = float64(c.stats.Hits) / float64(c.stats.Total) * 100
		return val, nil
	}

	c.stats.Misses++
	return nil, pie.ErrCacheNotFound
}

func (c *CustomCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.stats.Keys++
	return nil
}

func (c *CustomCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	c.stats.Keys--
	return nil
}

func (c *CustomCache) DeleteByPattern(ctx context.Context, pattern string) error {
	// 简化实现，只支持精确匹配
	return c.Delete(ctx, pattern)
}

func (c *CustomCache) DeleteByTags(ctx context.Context, tags ...string) error {
	// 简化实现，不做任何事
	return nil
}

func (c *CustomCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.data[key]
	return exists, nil
}

func (c *CustomCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string][]byte)
	c.stats.Keys = 0
	return nil
}

func (c *CustomCache) Stats() *pie.CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := *c.stats
	return &stats
}

func main() {
	ctx := context.Background()

	// 创建 MongoDB 引擎
	engine, err := pie.NewEngine(ctx, "testdb")
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Disconnect(ctx)

	fmt.Println("========== 缓存插件架构示例 ==========")

	// 1. 默认 Ristretto 缓存
	fmt.Println("\n1. 默认 Ristretto 缓存")
	demonstrateRistrettoCache(ctx, engine)

	// 2. Redis 缓存
	fmt.Println("\n2. Redis 缓存")
	demonstrateRedisCache(ctx, engine)

	// 3. 多层缓存链
	fmt.Println("\n3. 多层缓存链 (Ristretto + Redis)")
	demonstrateCacheChain(ctx, engine)

	// 4. 自定义缓存
	fmt.Println("\n4. 自定义缓存实现")
	demonstrateCustomCache(ctx, engine)

	// 5. Session 缓存集成
	fmt.Println("\n5. Session 缓存集成")
	demonstrateSessionCache(ctx, engine)
}

// 1. 默认 Ristretto 缓存
func demonstrateRistrettoCache(ctx context.Context, engine *pie.Engine) {
	// 启用默认 Ristretto 缓存
	engine.UseDefaultCache()

	// 创建产品
	product := &Product{
		Name:     "笔记本电脑",
		Price:    5999.0,
		Category: "电子产品",
	}

	// 使用缓存插入
	session := pie.Table[Product](engine)
	_, err := session.Cache(5*time.Minute).Insert(ctx, product)
	if err != nil {
		log.Printf("Insert failed: %v", err)
		return
	}

	// 查询并缓存
	products, err := session.Cache(5 * time.Minute).Find(ctx)
	if err != nil {
		log.Printf("Find failed: %v", err)
		return
	}

	fmt.Printf("找到 %d 个产品\n", len(products))
	for _, p := range products {
		fmt.Printf("- %s: ¥%.2f (%s)\n", p.Name, p.Price, p.Category)
	}

	// 显示缓存统计
	if cache := engine.Cache(); cache != nil {
		stats := cache.Stats()
		fmt.Printf("缓存统计: 命中率 %.2f%%, 总请求 %d\n", stats.HitRate, stats.Total)
	}
}

// 2. Redis 缓存
func demonstrateRedisCache(ctx context.Context, engine *pie.Engine) {
	// 尝试使用 Redis 缓存
	redisConfig := &pie.RedisCacheConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	}

	engine.UseRedis(redisConfig)

	// 创建产品
	product := &Product{
		Name:     "智能手机",
		Price:    2999.0,
		Category: "电子产品",
	}

	// 使用缓存插入
	session := pie.Table[Product](engine)
	_, err := session.Cache(10*time.Minute).Insert(ctx, product)
	if err != nil {
		log.Printf("Insert failed: %v", err)
		return
	}

	// 查询并缓存
	products, err := session.Cache(10 * time.Minute).Find(ctx)
	if err != nil {
		log.Printf("Find failed: %v", err)
		return
	}

	fmt.Printf("找到 %d 个产品\n", len(products))
	for _, p := range products {
		fmt.Printf("- %s: ¥%.2f (%s)\n", p.Name, p.Price, p.Category)
	}

	// 显示缓存统计
	if cache := engine.Cache(); cache != nil {
		stats := cache.Stats()
		fmt.Printf("Redis 缓存统计: 命中率 %.2f%%, 总请求 %d\n", stats.HitRate, stats.Total)
	}
}

// 3. 多层缓存链
func demonstrateCacheChain(ctx context.Context, engine *pie.Engine) {
	// 创建 Ristretto 缓存
	ristrettoCache, err := pie.NewRistrettoCache(nil)
	if err != nil {
		log.Printf("Failed to create Ristretto cache: %v", err)
		return
	}
	defer ristrettoCache.Close()

	// 创建 Mock Redis 缓存（用于演示）
	mockRedisCache := pie.NewMockRedisCache()
	defer mockRedisCache.Close()

	// 使用多层缓存
	engine.UseCache(ristrettoCache, mockRedisCache)

	// 创建产品
	product := &Product{
		Name:     "平板电脑",
		Price:    1999.0,
		Category: "电子产品",
	}

	// 使用缓存插入
	session := pie.Table[Product](engine)
	_, err = session.Cache(15*time.Minute).Insert(ctx, product)
	if err != nil {
		log.Printf("Insert failed: %v", err)
		return
	}

	// 查询并缓存
	products, err := session.Cache(15 * time.Minute).Find(ctx)
	if err != nil {
		log.Printf("Find failed: %v", err)
		return
	}

	fmt.Printf("找到 %d 个产品\n", len(products))
	for _, p := range products {
		fmt.Printf("- %s: ¥%.2f (%s)\n", p.Name, p.Price, p.Category)
	}

	// 显示缓存统计
	if cache := engine.Cache(); cache != nil {
		stats := cache.Stats()
		fmt.Printf("多层缓存统计: 命中率 %.2f%%, 总请求 %d\n", stats.HitRate, stats.Total)
	}
}

// 4. 自定义缓存
func demonstrateCustomCache(ctx context.Context, engine *pie.Engine) {
	// 创建自定义缓存
	customCache := NewCustomCache()

	// 使用自定义缓存
	engine.UseCache(customCache)

	// 创建产品
	product := &Product{
		Name:     "智能手表",
		Price:    1299.0,
		Category: "电子产品",
	}

	// 使用缓存插入
	session := pie.Table[Product](engine)
	_, err := session.Cache(20*time.Minute).Insert(ctx, product)
	if err != nil {
		log.Printf("Insert failed: %v", err)
		return
	}

	// 查询并缓存
	products, err := session.Cache(20 * time.Minute).Find(ctx)
	if err != nil {
		log.Printf("Find failed: %v", err)
		return
	}

	fmt.Printf("找到 %d 个产品\n", len(products))
	for _, p := range products {
		fmt.Printf("- %s: ¥%.2f (%s)\n", p.Name, p.Price, p.Category)
	}

	// 显示缓存统计
	if cache := engine.Cache(); cache != nil {
		stats := cache.Stats()
		fmt.Printf("自定义缓存统计: 命中率 %.2f%%, 总请求 %d\n", stats.HitRate, stats.Total)
	}
}

// 5. Session 缓存集成
func demonstrateSessionCache(ctx context.Context, engine *pie.Engine) {
	// 启用默认缓存
	engine.UseDefaultCache()

	// 创建产品
	products := []*Product{
		{Name: "游戏机", Price: 3999.0, Category: "电子产品"},
		{Name: "耳机", Price: 299.0, Category: "电子产品"},
		{Name: "充电器", Price: 99.0, Category: "配件"},
	}

	session := pie.Table[Product](engine)

	// 批量插入
	for _, product := range products {
		_, err := session.Insert(ctx, product)
		if err != nil {
			log.Printf("Insert failed: %v", err)
			return
		}
	}

	// 使用不同的缓存策略
	fmt.Println("\n--- 使用标签缓存 ---")
	results, err := session.CacheWithTags("electronics").Where("category", pie.Eq("category", "电子产品")).Find(ctx)
	if err != nil {
		log.Printf("Find with tags failed: %v", err)
		return
	}
	fmt.Printf("电子产品: %d 个\n", len(results))

	fmt.Println("\n--- 使用 TTL 抖动 ---")
	results, err = session.CacheWithJitter(10*time.Minute, 2*time.Minute).Where("price", pie.Gt("price", 1000)).Find(ctx)
	if err != nil {
		log.Printf("Find with jitter failed: %v", err)
		return
	}
	fmt.Printf("高价产品: %d 个\n", len(results))

	fmt.Println("\n--- 缓存空结果 ---")
	results, err = session.CacheEmpty(30*time.Second).Where("name", pie.Eq("name", "不存在的产品")).Find(ctx)
	if err != nil {
		log.Printf("Find empty cache failed: %v", err)
		return
	}
	fmt.Printf("空结果缓存: %d 个\n", len(results))

	// 显示最终缓存统计
	if cache := engine.Cache(); cache != nil {
		stats := cache.Stats()
		fmt.Printf("\n最终缓存统计: 命中率 %.2f%%, 总请求 %d, 键数 %d\n",
			stats.HitRate, stats.Total, stats.Keys)
	}
}
