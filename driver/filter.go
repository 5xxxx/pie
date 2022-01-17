package driver

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Condition interface {
	RegexFilter(key, pattern string) Condition
	ID(id interface{}) Condition
	// Eq Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) Condition

	// Gt {field: {$gt: value} } >
	Gt(key string, gt interface{}) Condition

	// Gte { qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) Condition

	// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) Condition

	// Lt {field: {$lt: value} } <
	Lt(key string, lt interface{}) Condition

	// Lte { field: { $lte: value} } <=
	Lte(key string, lte interface{}) Condition

	// Ne {field: {$ne: value} } !=
	Ne(key string, ne interface{}) Condition

	// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) Condition

	// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(filter Condition) Condition

	// Not { field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) Condition

	// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(filter Condition) Condition
	// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(filter Condition) Condition

	Exists(key string, exists bool, filter ...Condition) Condition

	// Type { field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) Condition

	// Expr Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(filter Condition) Condition

	// Regex todo 简单实现，后续增加支持
	Regex(key string, value string) Condition

	Filters() (bson.D, error)
	A() bson.A
	Err() error
	FilterBson(d interface{}) Condition
	FilterBy(object interface{}) Condition

	Clone() Condition
}
