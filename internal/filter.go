package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/NSObjects/pie/utils"

	"github.com/NSObjects/pie/driver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filter struct {
	d   bson.D
	err error
}

func (f *filter) FilterBy(object interface{}) driver.Condition {
	beanValue := reflect.ValueOf(object)
	if beanValue.Kind() != reflect.Struct {
		if m, ok := object.(bson.M); ok {
			for key, value := range m {
				f.Eq(key, value)
			}
			return f
		}

		if d, ok := object.(bson.D); ok {
			for _, v := range d {
				f.Eq(v.Key, v.Value)
			}
			return f
		}

		f.err = errors.New("needs a struct")
		return f
	}

	docType := reflect.TypeOf(object)
	for index := 0; index < docType.NumField(); index++ {
		fieldTag := docType.Field(index).Tag.Get("filter")
		if fieldTag != "" && fieldTag != "-" {
			split := strings.Split(fieldTag, ",")
			if len(split) > 0 {
				f.makeFilterValue(split[0], beanValue.Field(index).Interface())
			}
		}
	}

	return f
}

func (f *filter) Err() error {
	return f.err
}

func DefaultCondition() driver.Condition {
	return &filter{d: bson.D{}}
}

func (f *filter) Filters() (bson.D, error) {
	return f.d, f.err
}

func (f *filter) RegexFilter(key, pattern string) driver.Condition {
	v := primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

func (f *filter) ID(id interface{}) driver.Condition {
	if id == nil {
		return f
	}
	switch id.(type) {
	case string:
		objectId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			f.err = fmt.Errorf("id can't parse %v %v", id, err)
		}
		if objectId == primitive.NilObjectID {
			f.err = fmt.Errorf("id type must be string or primitive.ObjectID %v", id)
		}
		f.d = append(f.d, bson.E{Key: "_id", Value: objectId})
	case primitive.ObjectID:
		if id == primitive.NilObjectID {
			f.err = fmt.Errorf("id can't be nil %v", id)
			return f

		}
		f.d = append(f.d, bson.E{Key: "_id", Value: id})
	default:
		f.err = fmt.Errorf("id type must be string or primitive.ObjectID %v", id)
	}

	return f
}

//Equals a Specified Value
//{ qty: 20 }
//Field in Embedded Document Equals a Value
//{"item.name": "ab" }
// Equals an Array Value
//{ tags: [ "A", "B" ] }
func (f *filter) Eq(key string, value interface{}) driver.Condition {
	f.d = append(f.d, bson.E{Key: key, Value: value})
	return f
}

//{field: {$gt: value} } >
func (f *filter) Gt(key string, gt interface{}) driver.Condition {
	v := bson.M{
		"$gt": gt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ qty: { $gte: 20 } } >=
func (f *filter) Gte(key string, gte interface{}) driver.Condition {
	v := bson.M{
		"$gte": gte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (f *filter) In(key string, in interface{}) driver.Condition {
	v := bson.M{
		"$in": in,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{field: {$lt: value} } <
func (f *filter) Lt(key string, lt interface{}) driver.Condition {
	v := bson.M{
		"$lt": lt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ field: { $lte: value} } <=
func (f *filter) Lte(key string, lte interface{}) driver.Condition {
	v := bson.M{
		"$lte": lte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{field: {$ne: value} } !=
func (f *filter) Ne(key string, ne interface{}) driver.Condition {
	v := bson.M{
		"$ne": ne,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (f *filter) Nin(key string, nin interface{}) driver.Condition {
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
func (f *filter) And(filter driver.Condition) driver.Condition {
	f.d = append(f.d, bson.E{Key: "$and", Value: filter.A()})
	return f

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (f *filter) Not(key string, not interface{}) driver.Condition {
	v := bson.M{
		"$not": not,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (f *filter) Nor(filter driver.Condition) driver.Condition {
	f.d = append(f.d, bson.E{Key: "$nor", Value: filter.A()})
	return f
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (f *filter) Or(filter driver.Condition) driver.Condition {
	f.d = append(f.d, bson.E{Key: "$or", Value: filter.A()})
	return f
}

func (f *filter) Exists(key string, exists bool, filter ...driver.Condition) driver.Condition {
	m := bson.M{
		"$exists": exists,
	}

	for _, v := range filter {
		filters, err := v.Filters()
		if err != nil {
			f.err = err
			return f
		}
		for _, fv := range filters {
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
func (f *filter) Type(key string, t interface{}) driver.Condition {
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
func (f *filter) Expr(filter driver.Condition) driver.Condition {
	f.d = append(f.d, bson.E{Key: "$expr", Value: filter.A()})
	return f
}

//todo 简单实现，后续增加支持
func (f *filter) Regex(key string, value interface{}) driver.Condition {
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

func (f *filter) makeFilterValue(field string, value interface{}) {
	if utils.IsZero(value) {
		return
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Struct:
		f.makeStructValue(field, v)
	case reflect.Array:
		return
	}
	f.Eq(field, value)
}

func (f *filter) makeStructValue(field string, value reflect.Value) {
	for index := 0; index < value.NumField(); index++ {
		docType := reflect.TypeOf(value.Interface())
		tag := docType.Field(index).Tag.Get("bson")
		if tag != "" {
			if !utils.IsZero(value.Field(index)) {
				fieldTags := fmt.Sprintf("%s.%s", field, tag)
				f.makeFilterValue(fieldTags, value.Field(index).Interface())
			}
		}
	}
}
