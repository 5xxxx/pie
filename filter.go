/*
 *
 * filter.go
 * pie
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

	Filters() bson.D
	A() bson.A
}

type filter struct {
	d bson.D
}

func DefaultCondition() Condition {
	return &filter{d: bson.D{}}
}

func (f *filter) Filters() bson.D {
	return f.d
}

func (f *filter) RegexFilter(key, pattern string) *filter {
	v := primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

func (f *filter) ID(id interface{}) *filter {
	if id == nil {
		return f
	}
	switch id.(type) {
	case string:
		objectId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			panic("id type must be string or primitive.ObjectID")
		}
		if objectId == primitive.NilObjectID {
			panic("id type must be string or primitive.ObjectID")
		}
		f.d = append(f.d, bson.E{Key: "_id", Value: objectId})
	case primitive.ObjectID:
		if id == primitive.NilObjectID {
			panic("id type must be string or primitive.ObjectID")
		}
		f.d = append(f.d, bson.E{Key: "_id", Value: id})
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
	f.d = append(f.d, bson.E{Key: key, Value: value})
	return f
}

//{field: {$gt: value} } >
func (f *filter) Gt(key string, gt interface{}) *filter {
	v := bson.M{
		"$gt": gt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ qty: { $gte: 20 } } >=
func (f *filter) Gte(key string, gte interface{}) *filter {
	v := bson.M{
		"$gte": gte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (f *filter) In(key string, in interface{}) *filter {
	v := bson.M{
		"$in": in,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{field: {$lt: value} } <
func (f *filter) Lt(key string, lt interface{}) *filter {
	v := bson.M{
		"$lt": lt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ field: { $lte: value} } <=
func (f *filter) Lte(key string, lte interface{}) *filter {
	v := bson.M{
		"$lte": lte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{field: {$ne: value} } !=
func (f *filter) Ne(key string, ne interface{}) *filter {
	v := bson.M{
		"$ne": ne,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (f *filter) Nin(key string, nin interface{}) *filter {
	v := bson.M{
		"$nin": nin,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
//$and: [
//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
// ]
func (f *filter) And(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$and", Value: filter.A()})
	return f

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (f *filter) Not(key string, not interface{}) *filter {
	v := bson.M{
		"$not": not,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (f *filter) Nor(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$nor", Value: filter.A()})
	return f
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (f *filter) Or(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$or", Value: filter.A()})
	return f
}

func (f *filter) Exists(key string, exists bool, filter ...Condition) *filter {
	m := bson.M{
		"$exists": exists,
	}
	for _, v := range filter {
		for _, fv := range v.Filters() {
			m[fv.Key] = fv.Value
		}
	}

	f.d = append(f.d, bson.E{Key: key, Value: m})
	return f
}

//{ field: { $type: <BSON type> } }
// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
func (f *filter) Type(key string, t interface{}) *filter {
	v := bson.M{
		"$type": t,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (f *filter) Expr(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$expr", Value: filter.A()})
	return f
}

//todo 简单实现，后续增加支持
func (f *filter) Regex(key string, value interface{}) *filter {
	v := bson.M{
		"$regex":   value,
		"$options": "i",
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

func (f *filter) A() bson.A {
	var fs bson.A

	for _, value := range f.d {
		fs = append(fs, value)
	}
	return fs
}
