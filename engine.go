package pie

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Engine MongoDB engine, responsible for connection management and global configuration
type Engine struct {
	client       *mongo.Client
	database     *mongo.Database
	config       *Config
	nameMapper   NameMapper
	tagParser    *TagParser
	hooks        *HookManager
	queryLogger  *QueryLogger
	cacheManager *CacheManager
}

// NewEngine creates a new MongoDB engine
func NewEngine(ctx context.Context, database string, opts ...EngineOption) (*Engine, error) {
	// Create default configuration
	engineConfig := &EngineConfig{
		config:      DefaultConfig(),
		nameMapper:  &SnakeMapper{},
		tagParser:   NewTagParser(&SnakeMapper{}),
		queryLogger: NewQueryLogger(os.Stdout),
		hooks:       NewHookManager(),
	}

	// Apply options
	for _, opt := range opts {
		opt(engineConfig)
	}

	// Set database name
	engineConfig.config.Database = database

	// Build client options
	clientOpts, err := engineConfig.buildClientOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to build client options: %w", err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get database
	db := client.Database(database)

	// Update tag parser with name mapper
	engineConfig.tagParser = NewTagParser(engineConfig.nameMapper)

	engine := &Engine{
		client:      client,
		database:    db,
		config:      engineConfig.config,
		nameMapper:  engineConfig.nameMapper,
		tagParser:   engineConfig.tagParser,
		hooks:       engineConfig.hooks,
		queryLogger: engineConfig.queryLogger,
	}

	return engine, nil
}

// Database gets the database instance
func (e *Engine) Database() *mongo.Database {
	return e.database
}

// Client gets the MongoDB client
func (e *Engine) Client() *mongo.Client {
	return e.client
}

// Collection gets a collection
func (e *Engine) Collection(name string) *mongo.Collection {
	return e.database.Collection(name)
}

// CollectionForStruct gets a collection based on struct
func (e *Engine) CollectionForStruct(v any) (*mongo.Collection, error) {
	info, err := e.tagParser.ParseStruct(v)
	if err != nil {
		return nil, fmt.Errorf("failed to parse struct: %w", err)
	}

	return e.database.Collection(info.Name), nil
}

// Table creates a type-safe session
func Table[T any](engine *Engine) *Session[T] {
	var zero T
	collection, err := engine.CollectionForStruct(zero)
	return &Session[T]{
		engine:     engine,
		collection: collection,
		query:      NewQuery(),
		options:    NewSessionOptions(),
		initErr:    err, // New field to save initialization error
	}
}

// Ping test connection
func (e *Engine) Ping(ctx context.Context) error {
	return e.client.Ping(ctx, nil)
}

// Disconnect disconnect
func (e *Engine) Disconnect(ctx context.Context) error {
	return e.client.Disconnect(ctx)
}

// Close close engine (alias method)
func (e *Engine) Close(ctx context.Context) error {
	return e.Disconnect(ctx)
}

// IsConnected check if connected
func (e *Engine) IsConnected() bool {
	return e.client != nil
}

// GetConfig get configuration
func (e *Engine) GetConfig() *Config {
	return e.config
}

// GetNameMapper get name mapper
func (e *Engine) GetNameMapper() NameMapper {
	return e.nameMapper
}

// GetTagParser get tag parser
func (e *Engine) GetTagParser() *TagParser {
	return e.tagParser
}

// Hooks get hook manager
func (e *Engine) Hooks() *HookManager {
	return e.hooks
}

// QueryLogger get query logger
func (e *Engine) QueryLogger() *QueryLogger {
	return e.queryLogger
}

// EnableQueryLog enable query log
func (e *Engine) EnableQueryLog() {
	e.queryLogger.Enable()
}

// DisableQueryLog disable query log
func (e *Engine) DisableQueryLog() {
	e.queryLogger.Disable()
}

// SetQueryLogOutput set query log output destination
func (e *Engine) SetQueryLogOutput(writer io.Writer) {
	e.queryLogger.writer = writer
}

// SetQueryLogFormatter set log formatter
func (e *Engine) SetQueryLogFormatter(formatter LogFormatter) {
	e.queryLogger.SetFormatter(formatter)
}

// WithTransaction execute transaction
func (e *Engine) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	session, err := e.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// WithTransactionOptions execute transaction with options
func (e *Engine) WithTransactionOptions(ctx context.Context, fn func(context.Context) error, opts ...options.Lister[options.SessionOptions]) error {
	session, err := e.client.StartSession(opts...)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// ========== Cache Management Methods ==========

// UseCache enables caching with multiple cache instances
func (e *Engine) UseCache(caches ...Cache) *Engine {
	if len(caches) > 0 {
		e.cacheManager = NewCacheManager(caches, nil)
	}
	return e
}

// UseRistretto enables Ristretto memory cache
func (e *Engine) UseRistretto(config *RistrettoCacheConfig) *Engine {
	ristrettoCache, err := NewRistrettoCache(config)
	if err == nil {
		e.cacheManager = NewCacheManager([]Cache{ristrettoCache}, nil)
	}
	return e
}

// UseRedis enables Redis cache
func (e *Engine) UseRedis(config *RedisCacheConfig) *Engine {
	redisCache, err := NewRedisCache(config)
	if err == nil {
		e.cacheManager = NewCacheManager([]Cache{redisCache}, nil)
	}
	return e
}

// UseDefaultCache enables default Ristretto cache
func (e *Engine) UseDefaultCache() *Engine {
	return e.UseRistretto(nil)
}

// Cache gets the cache manager
func (e *Engine) Cache() *CacheManager {
	return e.cacheManager
}

// DisableCache disables caching
func (e *Engine) DisableCache() *Engine {
	if e.cacheManager != nil {
		e.cacheManager.config.Enabled = false
	}
	return e
}

// EnableCache enables caching
func (e *Engine) EnableCache() *Engine {
	if e.cacheManager != nil {
		e.cacheManager.config.Enabled = true
	}
	return e
}
