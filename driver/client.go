package driver

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is an interface that provides various methods for interacting with a MongoDB database.
//
// FindPagination retrieves paginated documents from the specified collection based on the given page and count.
// FindOneAndReplace replaces a single document that matches the filter with the provided document.
// FindOneAndUpdate updates a single document that matches the filter with the provided update document.
// FindAndDelete deletes a single document that matches the filter.
// FindOne retrieves a single document that matches the filter.
// FindAll retrieves all documents that match the filter.
// RegexFilter applies a regular expression filter to the session.
// Distinct returns an array of distinct values for a specified field across a collection.
// FindOneAndUpdateBson updates a single document that matches the BSON filter with the provided update BSON document.
// InsertOne inserts a single document into the collection.
// InsertMany inserts multiple documents into the collection.
// BulkWrite performs multiple write operations in bulk on the collection.
// ReplaceOne replaces a single document that matches the filter with the provided replacement document.
// Update updates a single document that matches the filter with the provided update document.
// UpdateMany updates multiple documents that match the filter with the provided update document.
// UpdateOneBson updates a single document that matches the BSON filter with the provided update BSON document.
// UpdateManyBson updates multiple documents that match the BSON filter with the provided update BSON document.
// SoftDeleteOne performs a soft delete on a single document that matches the filter.
// SoftDeleteMany performs a soft delete on multiple documents that match the filter.
// DeleteOne deletes a single document that matches the filter.
// DeleteMany deletes multiple documents that match the filter.
// DataBase returns the underlying mongo.Database associated with the client.
// Collection returns the mongo.Collection with the specified name and options.
// Ping checks if the client is connected to the database and returns an error if not.
// Connect establishes a connection to the database.
// Disconnect closes the connection to the database.
// Soft sets the soft filter state of the session.
// FilterBy applies the specified object as a filter to the session.
// Filter adds a key-value filter to the session.
// Asc adds the specified column names as sort order in ascending order to the session.
// Eq adds an equal filter to the session.
// Ne adds a not equal filter to the session.
// Nin adds a not in filter to the session.
// Nor adds a nor condition to the session.
// Exists adds an exists condition to the session.
// Type adds a type condition to the session.
// Expr adds an expression filter to the session.
// Regex adds a regex filter to the session.
// ID adds an _id filter to the session.
// Gt adds a greater than filter to the session.
// Gte adds a greater than or equal to filter to the session.
// Lt adds a less than filter to the session.
// Lte adds a less than or equal to filter to the session.
// In adds an in filter to the session.
// And adds an and condition to the session.
// Not adds a not condition to the session.
// Or adds an or condition to the session.
// Limit sets the maximum number of documents the session should return.
// Skip sets the number of documents the session should skip.
// Count returns the number of documents that match the filter.
// Desc adds the specified column names as sort order in descending order to the session.
// FilterBson adds a BSON filter to the session.
// NewIndexes returns a new Indexes instance for managing collection indexes.
// DropAll drops all indexes from the specified document's collection.
// DropOne drops the specified index from the specified document's collection.
// AddIndex adds an index to the specified document's collection.
// NewSession returns a new session for performing multiple operations.
// Aggregate returns an aggregate operation for the specified document's collection.
// CollectionNameForStruct returns the collection name for the specified document struct.
// CollectionNameForSlice returns the collection name for the specified document slice.
// Transaction starts a transaction and executes the provided function within the transaction context.
// TransactionWithOptions starts a transaction with the provided options and executes the provided function within the transaction context.
type Client interface {
	FindPagination(needCount bool, doc interface{}, ctx ...context.Context) (int64, error)
	FindOneAndReplace(doc interface{}, ctx ...context.Context) error
	FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error)
	FindAndDelete(doc interface{}, ctx ...context.Context) error
	FindOne(doc interface{}, ctx ...context.Context) error
	FindAll(docs interface{}, ctx ...context.Context) error
	RegexFilter(key, pattern string) Session
	Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error)
	FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error)

	InsertOne(v interface{}, ctx ...context.Context) (primitive.ObjectID, error)
	InsertMany(v interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error)
	BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error)
	ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	Update(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	SoftDeleteOne(filter interface{}, ctx ...context.Context) error
	SoftDeleteMany(filter interface{}, ctx ...context.Context) error
	DeleteOne(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error)

	DataBase() *mongo.Database
	// Collection(name string, db ...string) *mongo.Collection
	Collection(name string, collOpts []*options.CollectionOptions, db ...string) *mongo.Collection
	Ping() error
	Connect(ctx ...context.Context) (err error)
	Disconnect(ctx ...context.Context) error

	// Soft filter
	Soft(s bool) Session
	FilterBy(object interface{}) Session
	Filter(key string, value interface{}) Session
	Asc(colNames ...string) Session
	Eq(key string, value interface{}) Session
	Ne(key string, ne interface{}) Session
	Nin(key string, nin interface{}) Session
	Nor(c Condition) Session
	Exists(key string, exists bool, filter ...Condition) Session
	Type(key string, t interface{}) Session
	Expr(filter Condition) Session
	Regex(key string, value string) Session
	ID(id interface{}) Session
	Gt(key string, value interface{}) Session
	Gte(key string, value interface{}) Session
	Lt(key string, value interface{}) Session
	Lte(key string, value interface{}) Session
	In(key string, value interface{}) Session
	And(filter Condition) Session
	Not(key string, value interface{}) Session
	Or(filter Condition) Session
	Limit(limit int64) Session
	Skip(skip int64) Session
	Count(i interface{}, ctx ...context.Context) (int64, error)
	Desc(s1 ...string) Session
	FilterBson(d bson.D) Session

	NewIndexes() Indexes
	DropAll(doc interface{}, ctx ...context.Context) error
	DropOne(doc interface{}, name string, ctx ...context.Context) error
	AddIndex(keys interface{}, opt ...*options.IndexOptions) Indexes

	NewSession() Session
	Aggregate() Aggregate

	CollectionNameForStruct(doc interface{}) (*schemas.Collection, error)
	CollectionNameForSlice(doc interface{}) (*schemas.Collection, error)
	Transaction(ctx context.Context, f schemas.TransFunc) error
	TransactionWithOptions(ctx context.Context, f schemas.TransFunc, opt ...*options.SessionOptions) error
}
