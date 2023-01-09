package driver

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	// FindPagination find
	FindPagination(page, count int64, doc interface{}, ctx ...context.Context) error
	FindOneAndReplace(doc interface{}, ctx ...context.Context) error
	FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error)
	FindAndDelete(doc interface{}, ctx ...context.Context) error
	FindOne(doc interface{}, ctx ...context.Context) error
	FindAll(docs interface{}, ctx ...context.Context) error
	RegexFilter(key, pattern string) Session
	Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error)
	FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error)

	// InsertOne insert
	InsertOne(v interface{}, ctx ...context.Context) (primitive.ObjectID, error)
	InsertMany(v interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error)
	BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error)
	ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	// Update update
	Update(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)
	// UpdateMany The following operation updates all of the documents with quantity value less than 50.
	UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	// SoftDeleteOne delete
	SoftDeleteOne(filter interface{}, ctx ...context.Context) error
	SoftDeleteMany(filter interface{}, ctx ...context.Context) error
	DeleteOne(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error)

	// DataBase db operation
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
	// NewIndexes indexes
	NewIndexes() Indexes
	DropAll(doc interface{}, ctx ...context.Context) error
	DropOne(doc interface{}, name string, ctx ...context.Context) error
	AddIndex(keys interface{}, opt ...*options.IndexOptions) Indexes

	// NewSession session
	NewSession() Session
	Aggregate() Aggregate
	// CollectionNameForStruct SetDatabase(string string) Client
	CollectionNameForStruct(doc interface{}) (*schemas.Collection, error)
	CollectionNameForSlice(doc interface{}) (*schemas.Collection, error)
	Transaction(ctx context.Context, f schemas.TransFunc) error
	TransactionWithOptions(ctx context.Context, f schemas.TransFunc, opt ...*options.SessionOptions) error
}
