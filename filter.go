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

func (s *filter) Filters() bson.M {
	return s.m
}

func (s *filter) RegexFilter(key, pattern string) *filter {
	s.m[key] = primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	return s
}

func (s *filter) ID(id interface{}) *filter {
	if id == nil {
		return s
	}
	switch id.(type) {
	case string:
		objectId, _ := primitive.ObjectIDFromHex(id.(string))
		s.m["_id"] = objectId
	case primitive.ObjectID:
		s.m["_id"] = id
	default:
		panic("id type must be string or primitive.ObjectID")
	}
	return s
}

//Equals a Specified Value
//{ qty: 20 }
//Field in Embedded Document Equals a Value
//{"item.name": "ab" }
// Equals an Array Value
//{ tags: [ "A", "B" ] }
func (s *filter) Eq(key string, value interface{}) *filter {
	s.m[key] = value
	return s
}

//{field: {$gt: value} } >
func (s *filter) Gt(key string, gt interface{}) *filter {
	s.m[key] = bson.M{
		"$gt": gt,
	}
	return s
}

//{ qty: { $gte: 20 } } >=
func (s *filter) Gte(key string, gte interface{}) *filter {
	s.m[key] = bson.M{
		"$gte": gte,
	}
	return s
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (s *filter) In(key string, in interface{}) *filter {
	s.m[key] = bson.M{
		"$in": in,
	}
	return s
}

//{field: {$lt: value} } <
func (s *filter) Lt(key string, lt interface{}) *filter {
	s.m[key] = bson.M{
		"$lt": lt,
	}
	return s
}

//{ field: { $lte: value} } <=
func (s *filter) Lte(key string, lte interface{}) *filter {
	s.m[key] = bson.M{
		"$lte": lte,
	}
	return s
}

//{field: {$ne: value} } !=
func (s *filter) Ne(key string, ne interface{}) *filter {
	s.m[key] = bson.M{
		"$ne": ne,
	}
	return s
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (s *filter) Nin(key string, nin interface{}) *filter {
	s.m[key] = bson.M{
		"$nin": nin,
	}
	return s
}

//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
//$and: [
//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
// ]
func (s *filter) And(filter Condition) *filter {
	s.m["$and"] = filter.A()
	return s

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (s *filter) Not(key string, not interface{}) *filter {
	s.m[key] = bson.M{
		"$not": not,
	}
	return s
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (s *filter) Nor(filter Condition) *filter {
	s.m["$nor"] = filter.A()
	return s
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (s *filter) Or(filter Condition) *filter {
	s.m["$or"] = filter.A()
	return s
}

func (s *filter) Exists(key string, exists bool, filter ...Condition) *filter {
	m := bson.M{
		"$exists": exists,
	}
	for _, v := range filter {
		for fk, fv := range v.Filters() {
			m[fk] = fv
		}
	}
	s.m[key] = m
	return s
}

//{ field: { $type: <BSON type> } }
// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
func (s *filter) Type(key string, t interface{}) *filter {
	s.m[key] = bson.M{
		"$type": t,
	}
	return s
}

//Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (s *filter) Expr(filter Condition) *filter {
	s.m["$expr"] = filter.A()
	return s
}

//todo 简单实现，后续增加支持
func (s *filter) Regex(key string, value interface{}) *filter {
	s.m[key] = bson.M{
		"$regex":   value,
		"$options": "i",
	}

	return s
}

func (s *filter) A() []bson.E {
	var fs []bson.E

	for key, value := range s.m {
		fs = append(fs, bson.E{Key: key, Value: value})
	}
	return fs
}
