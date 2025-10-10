package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Database database operation wrapper
type Database struct {
	engine   *Engine
	database *mongo.Database
}

// NewDatabase create database operation wrapper
func NewDatabase(engine *Engine) *Database {
	return &Database{
		engine:   engine,
		database: engine.database,
	}
}

// Collection get collection
func (d *Database) Collection(name string) *mongo.Collection {
	return d.database.Collection(name)
}

// CollectionForStruct get collection by struct
func (d *Database) CollectionForStruct(v interface{}) (*mongo.Collection, error) {
	return d.engine.CollectionForStruct(v)
}

// ListCollections list all collections
func (d *Database) ListCollections(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	start := time.Now()

	opts := options.ListCollections()
	cursor, err := d.database.ListCollections(ctx, filter, opts)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "system.namespaces",
			Operation:  "listCollections",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return cursor, err
}

// ListCollectionNames list all collection names
func (d *Database) ListCollectionNames(ctx context.Context, filter interface{}) ([]string, error) {
	start := time.Now()

	names, err := d.database.ListCollectionNames(ctx, filter)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "system.namespaces",
			Operation:  "listCollectionNames",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return names, err
}

// CreateCollection create collection
func (d *Database) CreateCollection(ctx context.Context, name string, opts *options.CreateCollectionOptions) error {
	start := time.Now()

	// Convert to builder pattern for v2
	var builder *options.CreateCollectionOptionsBuilder
	if opts != nil {
		builder = options.CreateCollection()
		// Apply options from opts to builder
		// Note: This is a simplified implementation
	} else {
		builder = options.CreateCollection()
	}

	err := d.database.CreateCollection(ctx, name, builder)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: name,
			Operation:  "createCollection",
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return err
}

// DropCollection drop collection
func (d *Database) DropCollection(ctx context.Context, name string) error {
	start := time.Now()

	err := d.database.Collection(name).Drop(ctx)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: name,
			Operation:  "dropCollection",
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return err
}

// RunCommand execute database command
func (d *Database) RunCommand(ctx context.Context, command interface{}) (*mongo.SingleResult, error) {
	start := time.Now()

	result := d.database.RunCommand(ctx, command)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "$cmd",
			Operation:  "runCommand",
			Document:   command,
			Duration:   time.Since(start),
			Error:      result.Err(),
		})
	}

	return result, result.Err()
}

// RunCommandCursor execute database command returning cursor
func (d *Database) RunCommandCursor(ctx context.Context, command interface{}) (*mongo.Cursor, error) {
	start := time.Now()

	cursor, err := d.database.RunCommandCursor(ctx, command)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "$cmd",
			Operation:  "runCommandCursor",
			Document:   command,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return cursor, err
}

// Aggregate execute aggregation at database level
func (d *Database) Aggregate(ctx context.Context, pipeline interface{}) (*mongo.Cursor, error) {
	start := time.Now()

	cursor, err := d.database.Aggregate(ctx, pipeline)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "$cmd",
			Operation:  "aggregate",
			Pipeline:   pipeline,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return cursor, err
}

// Watch listen to database level changes
func (d *Database) Watch(ctx context.Context, pipeline interface{}) (*mongo.ChangeStream, error) {
	start := time.Now()

	stream, err := d.database.Watch(ctx, pipeline)

	// Record query log
	if d.engine.queryLogger.IsEnabled() {
		d.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "$cmd",
			Operation:  "watch",
			Pipeline:   pipeline,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return stream, err
}

// DatabaseStats get database stats
func (d *Database) DatabaseStats(ctx context.Context) (bson.M, error) {
	result, err := d.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}})
	if err != nil {
		return nil, err
	}

	var stats bson.M
	err = result.Decode(&stats)
	return stats, err
}

// ServerStatus get server status
func (d *Database) ServerStatus(ctx context.Context) (bson.M, error) {
	result, err := d.RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}})
	if err != nil {
		return nil, err
	}

	var status bson.M
	err = result.Decode(&status)
	return status, err
}

// BuildInfo get build info
func (d *Database) BuildInfo(ctx context.Context) (bson.M, error) {
	result, err := d.RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}})
	if err != nil {
		return nil, err
	}

	var info bson.M
	err = result.Decode(&info)
	return info, err
}

// Ping check database connection
func (d *Database) Ping(ctx context.Context) error {
	result, err := d.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}})
	if err != nil {
		return err
	}

	var response bson.M
	return result.Decode(&response)
}

// GetProfilingLevel get profiling level
func (d *Database) GetProfilingLevel(ctx context.Context) (int32, error) {
	result, err := d.RunCommand(ctx, bson.D{{Key: "profile", Value: -1}})
	if err != nil {
		return 0, err
	}

	var response bson.M
	if err := result.Decode(&response); err != nil {
		return 0, err
	}

	if level, ok := response["was"].(int32); ok {
		return level, nil
	}

	return 0, fmt.Errorf("unable to get profiling level")
}

// SetProfilingLevel set profiling level
// level: 0=off, 1=slow ops, 2=all ops
func (d *Database) SetProfilingLevel(ctx context.Context, level int32, slowMs int32) error {
	command := bson.D{
		{Key: "profile", Value: level},
	}
	if slowMs > 0 {
		command = append(command, bson.E{Key: "slowms", Value: slowMs})
	}

	_, err := d.RunCommand(ctx, command)
	return err
}

// GetProfilingData get profiling data
func (d *Database) GetProfilingData(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return d.Collection("system.profile").Find(ctx, filter)
}

// ClearProfilingData clear profiling data
func (d *Database) ClearProfilingData(ctx context.Context) error {
	return d.Collection("system.profile").Drop(ctx)
}

// CreateView create view
func (d *Database) CreateView(ctx context.Context, viewName, sourceCollection string, pipeline interface{}) error {
	command := bson.D{
		{Key: "create", Value: viewName},
		{Key: "viewOn", Value: sourceCollection},
		{Key: "pipeline", Value: pipeline},
	}

	_, err := d.RunCommand(ctx, command)
	return err
}

// DropView drop view
func (d *Database) DropView(ctx context.Context, viewName string) error {
	return d.DropCollection(ctx, viewName)
}

// ListViews list all views
func (d *Database) ListViews(ctx context.Context) (*mongo.Cursor, error) {
	filter := bson.D{
		{Key: "type", Value: "view"},
	}
	return d.ListCollections(ctx, filter)
}

// DatabaseInfo database info
type DatabaseInfo struct {
	Name        string  `bson:"name"`
	SizeOnDisk  int64   `bson:"sizeOnDisk"`
	Empty       bool    `bson:"empty"`
	Collections int32   `bson:"collections"`
	Views       int32   `bson:"views"`
	Objects     int64   `bson:"objects"`
	AvgObjSize  float64 `bson:"avgObjSize"`
	DataSize    int64   `bson:"dataSize"`
	StorageSize int64   `bson:"storageSize"`
	TotalSize   int64   `bson:"totalSize"`
	Indexes     int32   `bson:"indexes"`
	IndexSize   int64   `bson:"indexSize"`
	FileSize    int64   `bson:"fileSize"`
	Ok          int32   `bson:"ok"`
}

// GetDatabaseInfo get database info
func (d *Database) GetDatabaseInfo(ctx context.Context) (*DatabaseInfo, error) {
	result, err := d.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}})
	if err != nil {
		return nil, err
	}

	var info DatabaseInfo
	err = result.Decode(&info)
	return &info, err
}
