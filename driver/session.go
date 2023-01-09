package driver

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error)

	FilterBy(object interface{}) Session

	Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error)

	ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	FindOneAndReplace(doc interface{}, ctx ...context.Context) error

	FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error)

	FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error)

	FindPagination(page, count int64, doc interface{}, ctx ...context.Context) error

	FindAndDelete(doc interface{}, ctx ...context.Context) error

	// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
	FindOne(doc interface{}, ctx ...context.Context) error

	// FindAll Find executes a find command and returns a Cursor over the matching documents in the collectionByName.
	FindAll(rowsSlicePtr interface{}, ctx ...context.Context) error

	// InsertOne executes an insert command to insert a single document into the collectionByName.
	InsertOne(doc interface{}, ctx ...context.Context) (primitive.ObjectID, error)

	// InsertMany executes an insert command to insert multiple documents into the collectionByName.
	InsertMany(docs interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error)

	// DeleteOne executes a delete command to delete at most one document from the collectionByName.
	DeleteOne(doc interface{}, ctx ...context.Context) (*mongo.DeleteResult, error)
	SoftDeleteOne(doc interface{}, ctx ...context.Context) error

	// DeleteMany executes a delete command to delete documents from the collectionByName.
	DeleteMany(doc interface{}, ctx ...context.Context) (*mongo.DeleteResult, error)

	SoftDeleteMany(doc interface{}, ctx ...context.Context) error

	Clone() Session
	Limit(i int64) Session

	Skip(i int64) Session

	Count(i interface{}, ctx ...context.Context) (int64, error)

	UpdateOne(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error)

	RegexFilter(key, pattern string) Session

	ID(id interface{}) Session

	Asc(colNames ...string) Session

	Desc(colNames ...string) Session

	Sort(colNames ...string) Session
	Soft(f bool) Session
	Filter(key string, value interface{}) Session
	FilterBson(d bson.D) Session
	// Eq Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) Session

	// Gt {field: {$gt: value} } >
	Gt(key string, gt interface{}) Session

	// Gte { qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) Session

	// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) Session

	// Lt {field: {$lt: value} } <
	Lt(key string, lt interface{}) Session

	// Lte { field: { $lte: value} } <=
	Lte(key string, lte interface{}) Session

	// Ne {field: {$ne: value} } !=
	Ne(key string, ne interface{}) Session

	// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) Session

	// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(c Condition) Session

	// Not { field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) Session

	// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(c Condition) Session

	// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(c Condition) Session

	Exists(key string, exists bool, filter ...Condition) Session

	// SetArrayFilters sets the value for the ArrayFilters field.
	SetArrayFilters(filters options.ArrayFilters) Session

	// SetOrdered sets the value for the Ordered field.
	SetOrdered(ordered bool) Session

	// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
	SetBypassDocumentValidation(b bool) Session

	// SetReturnDocument sets the value for the ReturnDocument field.
	SetReturnDocument(rd options.ReturnDocument) Session

	// SetUpsert sets the value for the Upsert field.
	SetUpsert(b bool) Session

	// SetCollation sets the value for the Collation field.
	SetCollation(collation *options.Collation) Session

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) Session

	// SetProjection sets the value for the Projection field.
	SetProjection(projection interface{}) Session

	// SetSort sets the value for the Sort field.
	SetSort(sort interface{}) Session

	// SetHint sets the value for the Hint field.
	SetHint(hint interface{}) Session

	// Type { field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) Session

	// Expr Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(c Condition) Session

	// Regex todo 简单实现，后续增加支持
	Regex(key string, value string) Session

	SetDatabase(db string) Session

	SetCollRegistry(r *bsoncodec.Registry) Session

	SetCollReadPreference(rp *readpref.ReadPref) Session

	SetCollWriteConcern(wc *writeconcern.WriteConcern) Session

	SetReadConcern(rc *readconcern.ReadConcern) Session
}
