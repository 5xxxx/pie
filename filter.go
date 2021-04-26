package pie

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Condition interface {
	RegexFilter(key, pattern string) *filter
	ID(id interface{}) *filter
	// Eq Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) *filter

	// Gt {field: {$gt: value} } >
	Gt(key string, gt interface{}) *filter

	// Gte { qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) *filter

	// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) *filter

	// Lt {field: {$lt: value} } <
	Lt(key string, lt interface{}) *filter

	// Lte { field: { $lte: value} } <=
	Lte(key string, lte interface{}) *filter

	// Ne {field: {$ne: value} } !=
	Ne(key string, ne interface{}) *filter

	// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) *filter

	// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(filter Condition) *filter

	// Not { field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) *filter

	// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(filter Condition) *filter
	// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(filter Condition) *filter

	Exists(key string, exists bool, filter ...Condition) *filter

	// Type { field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) *filter

	// Expr Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(filter Condition) *filter

	// Regex todo 简单实现，后续增加支持
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
	switch id.(type) {
	case string:
		objectId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			panic(fmt.Sprintf("id type must be string or primitive.ObjectID \n %s,%s", err.Error(), id))
		}

		if objectId == primitive.NilObjectID {
			panic(fmt.Sprintf("id type must be string or primitive.ObjectID \n %s", id))
		}

		f.d = append(f.d, bson.E{Key: "_id", Value: objectId})
	case primitive.ObjectID:
		if id == primitive.NilObjectID {
			panic(fmt.Sprintf("id type must be string or primitive.ObjectID \n %s", id))
		}
		f.d = append(f.d, bson.E{Key: "_id", Value: id})
	default:
		panic("id type must be string or primitive.ObjectID")
	}

	return f
}

// Eq Equals a Specified Value
//{ qty: 20 }
//Field in Embedded Document Equals a Value
//{"item.name": "ab" }
// Equals an Array Value
//{ tags: [ "A", "B" ] }
func (f *filter) Eq(key string, value interface{}) *filter {
	f.d = append(f.d, bson.E{Key: key, Value: value})
	return f
}

// Gt {field: {$gt: value} } >
func (f *filter) Gt(key string, gt interface{}) *filter {
	v := bson.M{
		"$gt": gt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Gte { qty: { $gte: 20 } } >=
func (f *filter) Gte(key string, gte interface{}) *filter {
	v := bson.M{
		"$gte": gte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (f *filter) In(key string, in interface{}) *filter {
	v := bson.M{
		"$in": in,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Lt {field: {$lt: value} } <
func (f *filter) Lt(key string, lt interface{}) *filter {
	v := bson.M{
		"$lt": lt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Lte { field: { $lte: value} } <=
func (f *filter) Lte(key string, lte interface{}) *filter {
	v := bson.M{
		"$lte": lte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Ne {field: {$ne: value} } !=
func (f *filter) Ne(key string, ne interface{}) *filter {
	v := bson.M{
		"$ne": ne,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (f *filter) Nin(key string, nin interface{}) *filter {
	v := bson.M{
		"$nin": nin,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
//$and: [
//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
// ]
func (f *filter) And(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$and", Value: filter.A()})
	return f

}

// Not { field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (f *filter) Not(key string, not interface{}) *filter {
	v := bson.M{
		"$not": not,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (f *filter) Nor(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$nor", Value: filter.A()})
	return f
}

// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
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

// Type { field: { $type: <BSON type> } }
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

// Expr Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (f *filter) Expr(filter Condition) *filter {
	f.d = append(f.d, bson.E{Key: "$expr", Value: filter.A()})
	return f
}

// Regex todo 简单实现，后续增加支持
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
