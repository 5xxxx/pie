/*
 *
 * interface.go
 * pie
 *
 * Created by lintao on 2020/8/9 11:02 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Interface interface {
	BulkWrite(ctx context.Context, docs interface{}) (*mongo.BulkWriteResult, error)

	FilterBy(object interface{}) *Session

	Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error)

	ReplaceOne(ctx context.Context, doc interface{}) (*mongo.UpdateResult, error)

	FindOneAndReplace(ctx context.Context, doc interface{}) error

	FindOneAndUpdate(ctx context.Context, doc interface{}) error

	FindAndDelete(ctx context.Context, doc interface{}) error

	// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
	FindOne(ctx context.Context, doc interface{}) error

	// Find executes a find command and returns a Cursor over the matching documents in the collectionByName.
	FindAll(ctx context.Context, rowsSlicePtr interface{}) error

	// InsertOne executes an insert command to insert a single document into the collectionByName.
	InsertOne(ctx context.Context, doc interface{}) (primitive.ObjectID, error)

	// InsertMany executes an insert command to insert multiple documents into the collectionByName.
	InsertMany(ctx context.Context, docs interface{}) (*mongo.InsertManyResult, error)

	// DeleteOne executes a delete command to delete at most one document from the collectionByName.
	DeleteOne(ctx context.Context, doc interface{}) (*mongo.DeleteResult, error)

	SoftDeleteOne(ctx context.Context, doc interface{}) error

	// DeleteMany executes a delete command to delete documents from the collectionByName.
	DeleteMany(ctx context.Context, doc interface{}) (*mongo.DeleteResult, error)

	SoftDeleteMany(ctx context.Context, doc interface{}) error

	Clone() *Session

	Limit(i int64) *Session

	Skip(i int64) *Session

	Count(i interface{}) (int64, error)

	Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error)

	RegexFilter(key, pattern string) *Session

	ID(id interface{}) *Session

	Asc(colNames ...string) *Session
	Desc(colNames ...string) *Session

	Filter(key string, value interface{}) *Session

	//Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) *Session

	//{field: {$gt: value} } >
	Gt(key string, gt interface{}) *Session

	//{ qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) *Session

	//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) *Session

	//{field: {$lt: value} } <
	Lt(key string, lt interface{}) *Session

	//{ field: { $lte: value} } <=
	Lte(key string, lte interface{}) *Session

	//{field: {$ne: value} } !=
	Ne(key string, ne interface{}) *Session

	//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) *Session

	//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(c Condition) *Session

	//{ field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) *Session

	// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(c Condition) *Session
	// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(c Condition) *Session

	Exists(key string, exists bool, filter ...Condition) *Session

	// SetArrayFilters sets the value for the ArrayFilters field.
	SetArrayFilters(filters options.ArrayFilters) *Session

	// SetOrdered sets the value for the Ordered field.
	SetOrdered(ordered bool) *Session

	// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
	SetBypassDocumentValidation(b bool) *Session

	// SetReturnDocument sets the value for the ReturnDocument field.
	SetReturnDocument(rd options.ReturnDocument) *Session

	// SetUpsert sets the value for the Upsert field.
	SetUpsert(b bool) *Session

	// SetCollation sets the value for the Collation field.
	SetCollation(collation *options.Collation) *Session

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) *Session

	// SetProjection sets the value for the Projection field.
	SetProjection(projection interface{}) *Session

	// SetSort sets the value for the Sort field.
	SetSort(sort interface{}) *Session

	// SetHint sets the value for the Hint field.
	SetHint(hint interface{}) *Session

	//{ field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) *Session

	//Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(c Condition) *Session

	//todo 简单实现，后续增加支持
	Regex(key string, value interface{}) *Session
}
