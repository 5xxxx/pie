/*
 *
 * filter.go
 * tugrik
 *
 * Created by lintao on 2020/8/13 8:29 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Condition interface {
	RegexFilter(key, pattern string) *filter
	ID(id interface{}) *filter
	//Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) *filter

	//{field: {$gt: value} } >
	Gt(key string, gt interface{}) *filter

	//{ qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) *filter

	//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) *filter

	//{field: {$lt: value} } <
	Lt(key string, lt interface{}) *filter

	//{ field: { $lte: value} } <=
	Lte(key string, lte interface{}) *filter

	//{field: {$ne: value} } !=
	Ne(key string, ne interface{}) *filter

	//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) *filter

	//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(filter Condition) *filter

	//{ field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) *filter

	// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(filter Condition) *filter
	// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(filter Condition) *filter

	Exists(key string, exists bool, filter ...Condition) *filter

	//{ field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) *filter

	//Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(filter Condition) *filter

	//todo 简单实现，后续增加支持
	Regex(key string, value interface{}) *filter

	Filters() bson.M
	A() []bson.E
}

type filter struct {
	m bson.M
}

func DefaultCondition() Condition {
	return &filter{m: bson.M{}}
}

func (f *filter) Filters() bson.M {
	return f.m
}

func (f *filter) RegexFilter(key, pattern string) *filter {
	f.m[key] = primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	return f
}

func (f *filter) ID(id interface{}) *filter {
	if id == nil {
		return f
	}
	switch id.(type) {
	case string:
		objectId, _ := primitive.ObjectIDFromHex(id.(string))
		f.m["_id"] = objectId
	case primitive.ObjectID:
		f.m["_id"] = id
	default:
		panic("id type must be string or primitive.ObjectID")
	}
	return f
}

//Equals a Specified Value
//{ qty: 20 }
//Field in Embedded Document Equals a Value
//{"item.name": "ab" }
// Equals an Array Value
//{ tags: [ "A", "B" ] }
func (f *filter) Eq(key string, value interface{}) *filter {
	f.m[key] = value
	return f
}

//{field: {$gt: value} } >
func (f *filter) Gt(key string, gt interface{}) *filter {
	f.m[key] = bson.M{
		"$gt": gt,
	}
	return f
}

//{ qty: { $gte: 20 } } >=
func (f *filter) Gte(key string, gte interface{}) *filter {
	f.m[key] = bson.M{
		"$gte": gte,
	}
	return f
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (f *filter) In(key string, in interface{}) *filter {
	f.m[key] = bson.M{
		"$in": in,
	}
	return f
}

//{field: {$lt: value} } <
func (f *filter) Lt(key string, lt interface{}) *filter {
	f.m[key] = bson.M{
		"$lt": lt,
	}
	return f
}

//{ field: { $lte: value} } <=
func (f *filter) Lte(key string, lte interface{}) *filter {
	f.m[key] = bson.M{
		"$lte": lte,
	}
	return f
}

//{field: {$ne: value} } !=
func (f *filter) Ne(key string, ne interface{}) *filter {
	f.m[key] = bson.M{
		"$ne": ne,
	}
	return f
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (f *filter) Nin(key string, nin interface{}) *filter {
	f.m[key] = bson.M{
		"$nin": nin,
	}
	return f
}

//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
//$and: [
//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
// ]
func (f *filter) And(filter Condition) *filter {
	f.m["$and"] = filter.A()
	return f

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (f *filter) Not(key string, not interface{}) *filter {
	f.m[key] = bson.M{
		"$not": not,
	}
	return f
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (f *filter) Nor(filter Condition) *filter {
	f.m["$nor"] = filter.A()
	return f
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (f *filter) Or(filter Condition) *filter {
	f.m["$or"] = filter.A()
	return f
}

func (f *filter) Exists(key string, exists bool, filter ...Condition) *filter {
	m := bson.M{
		"$exists": exists,
	}
	for _, v := range filter {
		for fk, fv := range v.Filters() {
			m[fk] = fv
		}
	}
	f.m[key] = m
	return f
}

//{ field: { $type: <BSON type> } }
// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
func (f *filter) Type(key string, t interface{}) *filter {
	f.m[key] = bson.M{
		"$type": t,
	}
	return f
}

//Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (f *filter) Expr(filter Condition) *filter {
	f.m["$expr"] = filter.A()
	return f
}

//todo 简单实现，后续增加支持
func (f *filter) Regex(key string, value interface{}) *filter {
	f.m[key] = bson.M{
		"$regex":   value,
		"$options": "i",
	}

	return f
}

func (f *filter) A() []bson.E {
	var fs []bson.E

	for key, value := range f.m {
		fs = append(fs, bson.E{Key: key, Value: value})
	}
	return fs
}
