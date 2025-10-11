package pie

import (
	"io"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Config configuration structure
type Config struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
	SocketTimeout  time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
	MaxIdleTime    time.Duration
	ReadPreference string
	WriteConcern   string
	ReadConcern    string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		ConnectTimeout: 10 * time.Second,
		SocketTimeout:  30 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    0,
		MaxIdleTime:    30 * time.Minute,
		ReadPreference: "primary",
		WriteConcern:   "majority",
		ReadConcern:    "majority",
	}
}

// EngineOption engine option function type
type EngineOption func(*EngineConfig)

// EngineConfig engine configuration
type EngineConfig struct {
	config      *Config
	nameMapper  NameMapper
	tagParser   *TagParser
	queryLogger *QueryLogger
	hooks       *HookManager
}

// WithURI set MongoDB connection URI
func WithURI(uri string) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.URI = uri
	}
}

// WithDatabase set database name
func WithDatabase(db string) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.Database = db
	}
}

// WithMapper set name mapper
func WithMapper(mapper NameMapper) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.nameMapper = mapper
	}
}

// WithConnectTimeout set connection timeout
func WithConnectTimeout(timeout time.Duration) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.ConnectTimeout = timeout
	}
}

// WithSocketTimeout set Socket timeout
func WithSocketTimeout(timeout time.Duration) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.SocketTimeout = timeout
	}
}

// WithMaxPoolSize set maximum connection pool size
func WithMaxPoolSize(size uint64) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.MaxPoolSize = size
	}
}

// WithMinPoolSize set minimum connection pool size
func WithMinPoolSize(size uint64) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.MinPoolSize = size
	}
}

// WithMaxIdleTime set maximum idle time
func WithMaxIdleTime(duration time.Duration) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.MaxIdleTime = duration
	}
}

// WithReadPreference set read preference
func WithReadPreference(pref string) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.ReadPreference = pref
	}
}

// WithWriteConcern set write concern
func WithWriteConcern(concern string) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.WriteConcern = concern
	}
}

// WithReadConcern set read concern
func WithReadConcern(concern string) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.config.ReadConcern = concern
	}
}

// buildClientOptions build MongoDB client options based on configuration
func (cfg *EngineConfig) buildClientOptions() ([]*options.ClientOptions, error) {
	var opts []*options.ClientOptions

	if cfg.config.URI != "" {
		opts = append(opts, options.Client().ApplyURI(cfg.config.URI))
	}

	if cfg.config.ConnectTimeout > 0 {
		opts = append(opts, options.Client().SetConnectTimeout(cfg.config.ConnectTimeout))
	}

	// SocketTimeout is not available in v2, using ServerSelectionTimeout instead
	if cfg.config.SocketTimeout > 0 {
		opts = append(opts, options.Client().SetServerSelectionTimeout(cfg.config.SocketTimeout))
	}

	if cfg.config.MaxPoolSize > 0 {
		opts = append(opts, options.Client().SetMaxPoolSize(cfg.config.MaxPoolSize))
	}

	if cfg.config.MinPoolSize > 0 {
		opts = append(opts, options.Client().SetMinPoolSize(cfg.config.MinPoolSize))
	}

	// MaxIdleTime is not available in v2, using MaxConnIdleTime instead
	if cfg.config.MaxIdleTime > 0 {
		opts = append(opts, options.Client().SetMaxConnIdleTime(cfg.config.MaxIdleTime))
	}

	return opts, nil
}

// WithQueryLog enable query logging and set output destination
func WithQueryLog(writer io.Writer) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.queryLogger = NewQueryLogger(writer)
		cfg.queryLogger.Enable()
	}
}

// WithQueryLogFormatter set query log formatter
func WithQueryLogFormatter(formatter LogFormatter) EngineOption {
	return func(cfg *EngineConfig) {
		if cfg.queryLogger == nil {
			cfg.queryLogger = NewQueryLogger(nil)
		}
		cfg.queryLogger.SetFormatter(formatter)
	}
}

// WithSlowQueryThreshold set slow query threshold
func WithSlowQueryThreshold(duration time.Duration) EngineOption {
	return func(cfg *EngineConfig) {
		if cfg.queryLogger == nil {
			cfg.queryLogger = NewQueryLogger(nil)
		}
		cfg.queryLogger.SetSlowQueryThreshold(duration)
	}
}

// WithHooks set hook manager
func WithHooks(manager *HookManager) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.hooks = manager
	}
}

// WithCache set cache instances
func WithCache(caches ...Cache) EngineOption {
	return func(cfg *EngineConfig) {
		// 这个选项需要在 Engine 创建后设置，因为需要 Engine 实例
		// 这里只是占位，实际逻辑在 Engine 的 UseCache 方法中
	}
}

// WithDefaultCache enable default Ristretto cache
func WithDefaultCache() EngineOption {
	return func(cfg *EngineConfig) {
		// 这个选项需要在 Engine 创建后设置
		// 这里只是占位，实际逻辑在 Engine 的 UseDefaultCache 方法中
	}
}
