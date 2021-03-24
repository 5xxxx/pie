/*
 *
 * filter.go
 * pie
 *
 * Created by lintao on 2020/8/13 8:29 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package driver

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Condition interface {
	RegexFilter(key, pattern string) Condition
	ID(id interface{}) Condition
	//Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) Condition

	//{field: {$gt: value} } >
	Gt(key string, gt interface{}) Condition

	//{ qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) Condition

	//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) Condition

	//{field: {$lt: value} } <
	Lt(key string, lt interface{}) Condition

	//{ field: { $lte: value} } <=
	Lte(key string, lte interface{}) Condition

	//{field: {$ne: value} } !=
	Ne(key string, ne interface{}) Condition

	//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) Condition

	//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(filter Condition) Condition

	//{ field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) Condition

	// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(filter Condition) Condition
	// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(filter Condition) Condition

	Exists(key string, exists bool, filter ...Condition) Condition

	//{ field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) Condition

	//Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(filter Condition) Condition

	//todo 简单实现，后续增加支持
	Regex(key string, value string) Condition

	Filters() (bson.D, error)
	A() bson.A
	Err() error
	FilterBson(d bson.D) Condition
	FilterBy(object interface{}) Condition
}
