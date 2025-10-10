package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var defaultEngine *Engine

// SetDefaultEngine 设置默认引擎
func SetDefaultEngine(engine *Engine) {
	defaultEngine = engine
}

// GetDefaultEngine 获取默认引擎
func GetDefaultEngine() *Engine {
	return defaultEngine
}

// MustNewEngine 创建新引擎，失败时panic
func MustNewEngine(ctx context.Context, database string, opts ...EngineOption) *Engine {
	engine, err := NewEngine(ctx, database, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to create engine: %v", err))
	}
	return engine
}

// MustTable 创建类型安全的会话，失败时panic
func MustTable[T any](engine *Engine) *Session[T] {
	if engine == nil {
		panic("engine is nil")
	}
	return Table[T](engine)
}

// MustTableWithDefault 使用默认引擎创建类型安全的会话，失败时panic
func MustTableWithDefault[T any]() *Session[T] {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return Table[T](defaultEngine)
}

// TableWithDefault 使用默认引擎创建类型安全的会话
func TableWithDefault[T any]() (*Session[T], error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return Table[T](defaultEngine), nil
}

// MustAggregate 创建聚合操作，失败时panic
func MustAggregate[T any](engine *Engine) *Aggregate[T] {
	if engine == nil {
		panic("engine is nil")
	}
	return NewAggregate[T](engine)
}

// MustAggregateWithDefault 使用默认引擎创建聚合操作，失败时panic
func MustAggregateWithDefault[T any]() *Aggregate[T] {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return NewAggregate[T](defaultEngine)
}

// AggregateWithDefault 使用默认引擎创建聚合操作
func AggregateWithDefault[T any]() (*Aggregate[T], error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return NewAggregate[T](defaultEngine), nil
}

// MustIndexes 创建索引管理器，失败时panic
func MustIndexes(engine *Engine) *Indexes {
	if engine == nil {
		panic("engine is nil")
	}
	return NewIndexes(engine)
}

// MustIndexesWithDefault 使用默认引擎创建索引管理器，失败时panic
func MustIndexesWithDefault() *Indexes {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return NewIndexes(defaultEngine)
}

// IndexesWithDefault 使用默认引擎创建索引管理器
func IndexesWithDefault() (*Indexes, error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return NewIndexes(defaultEngine), nil
}

// MustTransaction 创建事务管理器，失败时panic
func MustTransaction(engine *Engine) *Transaction {
	if engine == nil {
		panic("engine is nil")
	}
	return NewTransaction(engine)
}

// MustTransactionWithDefault 使用默认引擎创建事务管理器，失败时panic
func MustTransactionWithDefault() *Transaction {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return NewTransaction(defaultEngine)
}

// TransactionWithDefault 使用默认引擎创建事务管理器
func TransactionWithDefault() (*Transaction, error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return NewTransaction(defaultEngine), nil
}

// 便捷的全局方法

// Connect 连接到MongoDB
func Connect(ctx context.Context, uri, database string, opts ...EngineOption) (*Engine, error) {
	engineOpts := append([]EngineOption{WithURI(uri)}, opts...)
	return NewEngine(ctx, database, engineOpts...)
}

// MustConnect 连接到MongoDB，失败时panic
func MustConnect(ctx context.Context, uri, database string, opts ...EngineOption) *Engine {
	engine, err := Connect(ctx, uri, database, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to connect: %v", err))
	}
	return engine
}

// ConnectWithTimeout 带超时的连接
func ConnectWithTimeout(uri, database string, timeout time.Duration, opts ...EngineOption) (*Engine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	engineOpts := append([]EngineOption{WithURI(uri)}, opts...)
	return NewEngine(ctx, database, engineOpts...)
}

// MustConnectWithTimeout 带超时的连接，失败时panic
func MustConnectWithTimeout(uri, database string, timeout time.Duration, opts ...EngineOption) *Engine {
	engine, err := ConnectWithTimeout(uri, database, timeout, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to connect with timeout: %v", err))
	}
	return engine
}

// 初始化函数
func init() {
	// 设置默认配置
	// 这里可以添加一些全局初始化逻辑
}

// 版本信息
const (
	Version = "2.0.0"
	Author  = "pie-mongodb"
)

// GetVersion 获取版本信息
func GetVersion() string {
	return Version
}

// GetAuthor 获取作者信息
func GetAuthor() string {
	return Author
}

// NewWatcher 创建变更流监听器
func NewWatcher[T any](engine *Engine) *ChangeStreamWatcher[T] {
	return &ChangeStreamWatcher[T]{
		engine:    engine,
		options:   options.ChangeStream(),
		pipeline:  []bson.D{},
		watchType: WatchCollection,
	}
}

// NewDatabaseWatcher 创建数据库级别监听器
func NewDatabaseWatcher[T any](engine *Engine) *ChangeStreamWatcher[T] {
	return &ChangeStreamWatcher[T]{
		engine:    engine,
		database:  engine.database,
		options:   options.ChangeStream(),
		pipeline:  []bson.D{},
		watchType: WatchDatabase,
	}
}
