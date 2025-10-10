package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var defaultEngine *Engine

// SetDefaultEngine stores the engine instance that should be used when helper
// constructors are called without an explicit engine. It can be safely invoked
// at application startup to configure package-level defaults.
func SetDefaultEngine(engine *Engine) {
	defaultEngine = engine
}

// GetDefaultEngine returns the engine that has been registered through
// SetDefaultEngine. It may be nil when no default was configured.
func GetDefaultEngine() *Engine {
	return defaultEngine
}

// MustNewEngine constructs a new Engine and panics when initialization fails.
// It is useful in scenarios where startup failures should abort the
// application immediately.
func MustNewEngine(ctx context.Context, database string, opts ...EngineOption) *Engine {
	engine, err := NewEngine(ctx, database, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to create engine: %v", err))
	}
	return engine
}

// MustTable creates a strongly typed session bound to the provided engine and
// panics if the engine is nil. The helper simplifies setup code in examples
// where failing fast is preferred over manual error handling.
func MustTable[T any](engine *Engine) *Session[T] {
	if engine == nil {
		panic("engine is nil")
	}
	return Table[T](engine)
}

// MustTableWithDefault creates a strongly typed session using the default
// engine registered via SetDefaultEngine. A panic is triggered when no default
// engine has been configured, highlighting misconfiguration early.
func MustTableWithDefault[T any]() *Session[T] {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return Table[T](defaultEngine)
}

// TableWithDefault returns a strongly typed session backed by the default
// engine. An error is returned instead of panicking so callers can decide how
// to recover when the default engine is missing.
func TableWithDefault[T any]() (*Session[T], error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return Table[T](defaultEngine), nil
}

// MustAggregate constructs a new aggregation builder for the supplied engine
// and panics when the engine is nil. It mirrors MustTable to keep helpers
// consistent across different entry points.
func MustAggregate[T any](engine *Engine) *Aggregate[T] {
	if engine == nil {
		panic("engine is nil")
	}
	return NewAggregate[T](engine)
}

// MustAggregateWithDefault creates an aggregation builder using the default
// engine. A missing default engine is treated as a programmer error and causes
// a panic to surface the issue immediately.
func MustAggregateWithDefault[T any]() *Aggregate[T] {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return NewAggregate[T](defaultEngine)
}

// AggregateWithDefault creates an aggregation builder with the default engine
// and reports an error when no default engine is available.
func AggregateWithDefault[T any]() (*Aggregate[T], error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return NewAggregate[T](defaultEngine), nil
}

// MustIndexes creates an index manager for the given engine and panics if the
// engine is nil. The helper removes repetitive nil checks in applications that
// treat misconfiguration as fatal.
func MustIndexes(engine *Engine) *Indexes {
	if engine == nil {
		panic("engine is nil")
	}
	return NewIndexes(engine)
}

// MustIndexesWithDefault returns an index manager backed by the default engine
// or panics when the default engine has not been set.
func MustIndexesWithDefault() *Indexes {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return NewIndexes(defaultEngine)
}

// IndexesWithDefault is the error-returning counterpart of
// MustIndexesWithDefault. It enables graceful fallback when a default engine
// has not yet been registered.
func IndexesWithDefault() (*Indexes, error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return NewIndexes(defaultEngine), nil
}

// MustTransaction creates a transaction manager and panics when the provided
// engine is nil. Use it when your initialization logic should not proceed with
// an invalid engine reference.
func MustTransaction(engine *Engine) *Transaction {
	if engine == nil {
		panic("engine is nil")
	}
	return NewTransaction(engine)
}

// MustTransactionWithDefault constructs a transaction manager that uses the
// default engine. If the default engine has not been set, the function panics
// to surface the configuration error.
func MustTransactionWithDefault() *Transaction {
	if defaultEngine == nil {
		panic("default engine is not set")
	}
	return NewTransaction(defaultEngine)
}

// TransactionWithDefault creates a transaction manager from the default engine
// and returns an error when no default engine has been configured.
func TransactionWithDefault() (*Transaction, error) {
	if defaultEngine == nil {
		return nil, fmt.Errorf("default engine is not set")
	}
	return NewTransaction(defaultEngine), nil
}

// Convenience global helpers

// Connect instantiates a new Engine using the provided MongoDB connection URI
// and target database. Additional engine options can be supplied through opts.
func Connect(ctx context.Context, uri, database string, opts ...EngineOption) (*Engine, error) {
	engineOpts := append([]EngineOption{WithURI(uri)}, opts...)
	return NewEngine(ctx, database, engineOpts...)
}

// MustConnect wraps Connect and panics when establishing the connection fails.
// It is intended for simple programs where a failed connection should halt the
// process immediately.
func MustConnect(ctx context.Context, uri, database string, opts ...EngineOption) *Engine {
	engine, err := Connect(ctx, uri, database, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to connect: %v", err))
	}
	return engine
}

// ConnectWithTimeout opens a connection using a new context with the specified
// timeout. The helper ensures that connection attempts do not block
// indefinitely.
func ConnectWithTimeout(uri, database string, timeout time.Duration, opts ...EngineOption) (*Engine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	engineOpts := append([]EngineOption{WithURI(uri)}, opts...)
	return NewEngine(ctx, database, engineOpts...)
}

// MustConnectWithTimeout wraps ConnectWithTimeout and panics if the connection
// could not be established within the allotted time or when other errors
// occur.
func MustConnectWithTimeout(uri, database string, timeout time.Duration, opts ...EngineOption) *Engine {
	engine, err := ConnectWithTimeout(uri, database, timeout, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to connect with timeout: %v", err))
	}
	return engine
}

// init configures package-level defaults. Additional global initialization can
// be added here by applications embedding the library.
func init() {
	// Initialize default configuration. The block is intentionally left
	// minimal so downstream applications can modify package variables in
	// their own init functions if necessary.
}

// Version metadata constants
const (
	Version = "2.0.0"
	Author  = "pie-mongodb"
)

// GetVersion returns the semantic version identifier for the library.
func GetVersion() string {
	return Version
}

// GetAuthor returns the canonical author string associated with the project.
func GetAuthor() string {
	return Author
}

// NewWatcher constructs a collection-level change stream watcher. Callers can
// further configure the watcher before invoking Watch to start consuming
// change events.
func NewWatcher[T any](engine *Engine) *ChangeStreamWatcher[T] {
	return &ChangeStreamWatcher[T]{
		engine:    engine,
		options:   options.ChangeStream(),
		pipeline:  []bson.D{},
		watchType: WatchCollection,
	}
}

// NewDatabaseWatcher constructs a change stream watcher scoped to an entire
// database, enabling monitoring of multiple collections.
func NewDatabaseWatcher[T any](engine *Engine) *ChangeStreamWatcher[T] {
	return &ChangeStreamWatcher[T]{
		engine:    engine,
		database:  engine.database,
		options:   options.ChangeStream(),
		pipeline:  []bson.D{},
		watchType: WatchDatabase,
	}
}
