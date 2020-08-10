/*
 *
 * find_one_session.go
 * tugrik
 *
 * Created by lintao on 2020/7/17 10:16 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package tugrik

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"tugrik/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	m              bson.M
	engine         *Tugrik
	findOneOptions []*options.FindOneOptions
	findOptions    []*options.FindOptions
	insertManyOpts []*options.InsertManyOptions
	insertOneOpts  []*options.InsertOneOptions
	deleteOpts     []*options.DeleteOptions
	updateOpts     []*options.UpdateOptions
	countOpts      []*options.CountOptions
	distinctOpts   []*options.DistinctOptions
}

func NewFilter(engine *Tugrik) *Session {
	return &Session{engine: engine}
}

func (s *Session) Distinct(doc interface{}, columns string) ([]interface{}, error) {
	coll, err := s.getSliceColl(doc)
	if err != nil {
		return nil, err
	}
	return coll.Distinct(context.TODO(), columns, s.m, s.distinctOpts...)
}

// FindOne executes a find command and returns a SingleResult for one document in the collection.
func (s Session) FindOne(doc interface{}) (bool, error) {
	coll, err := s.getStructColl(doc)
	if err != nil {
		return false, err
	}

	result := coll.FindOne(context.Background(), s.m, s.findOneOptions...)
	if err = result.Err(); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	if err = result.Decode(doc); err != nil {
		return false, err
	}

	return true, nil
}

// Find executes a find command and returns a Cursor over the matching documents in the collection.
func (s *Session) FindAll(rowsSlicePtr interface{}) error {
	coll, err := s.getSliceColl(rowsSlicePtr)
	if err != nil {
		return err
	}
	cursor, err := coll.Find(context.Background(), s.m, s.findOptions...)
	if err != nil {
		return err
	}

	if err = cursor.All(context.Background(), rowsSlicePtr); err != nil {
		return err
	}

	return nil
}

// InsertOne executes an insert command to insert a single document into the collection.
func (s *Session) InsertOne(doc interface{}) error {
	coll, err := s.getStructColl(doc)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(context.Background(), doc, s.insertOneOpts...)
	return err
}

// InsertMany executes an insert command to insert multiple documents into the collection.
func (s *Session) InsertMany(docs []interface{}) error {
	coll, err := s.getSliceColl(docs)
	if err != nil {
		return err
	}
	_, err = coll.InsertMany(context.Background(), docs, s.insertManyOpts...)

	return err
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (s *Session) DeleteOne(doc interface{}) error {
	coll, err := s.getStructColl(doc)
	if err != nil {
		return err
	}
	_, err = coll.DeleteOne(context.Background(), s.m, s.deleteOpts...)
	return err
}

// DeleteMany executes a delete command to delete documents from the collection.
func (s *Session) DeleteMany(doc interface{}) error {
	coll, err := s.getStructColl(doc)
	if err != nil {
		return err
	}
	_, err = coll.DeleteMany(context.Background(), s.m, s.deleteOpts...)
	return err
}

func (s *Session) Limit(i int64) *Session {
	s.findOptions = append(s.findOptions, options.Find().SetLimit(i))
	return s
}

func (s *Session) Skip(i int64) *Session {
	s.findOptions = append(s.findOptions, options.Find().SetSkip(i))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSkip(i))
	return s
}

func (s *Session) Count(i interface{}) (int64, error) {
	coll, err := s.getStructColl(i)
	if err != nil {
		return 0, err
	}
	return coll.CountDocuments(context.Background(), s.m, s.countOpts...)
}

func (s *Session) Update(bean interface{}) error {
	coll, err := s.getStructColl(bean)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(context.Background(), s.m, bson.M{"$set": insertOmitemptyTag(bean)})

	return err
}

//todo update many
func (s *Session) UpdateMany(bean interface{}) error {
	coll, err := s.getSliceColl(bean)
	if err != nil {
		return err
	}
	_, err = coll.UpdateMany(context.Background(), s.m, bson.M{"$set": bean})

	return err
}

func (s *Session) RegexFilter(key, pattern string) *Session {
	s.m[key] = primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	return s
}

////func (s *Session) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
////	panic("implement me")
////}
////

func (s *Session) ID(id interface{}) *Session {
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

func (s *Session) Asc(colNames ...string) *Session {

	if len(colNames) == 0 {
		return s
	}

	var es bson.M
	for _, c := range colNames {
		es[c] = -1
	}
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	return s
}

func (s *Session) Desc(colNames ...string) *Session {
	if len(colNames) == 0 {
		return s
	}
	var es bson.M
	for _, c := range colNames {
		es[c] = 1
	}

	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	return s
}

func (s *Session) Filter(key string, value interface{}) *Session {
	return s.Eq(key, value)
}

//Equals a Specified Value
//{ qty: 20 }
//Field in Embedded Document Equals a Value
//{"item.name": "ab" }
// Equals an Array Value
//{ tags: [ "A", "B" ] }
func (s *Session) Eq(key string, value interface{}) *Session {
	s.m[key] = value
	return s
}

//{field: {$gt: value} } >
func (s *Session) Gt(key string, gt interface{}) *Session {
	s.m[key] = bson.M{
		"$gt": gt,
	}
	return s
}

//{ qty: { $gte: 20 } } >=
func (s *Session) Gte(key string, gte interface{}) *Session {
	s.m[key] = bson.M{
		"$gte": gte,
	}
	return s
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (s *Session) In(key string, in interface{}) *Session {
	s.m[key] = bson.M{
		"$in": in,
	}
	return s
}

//{field: {$lt: value} } <
func (s *Session) Lt(key string, lt interface{}) *Session {
	s.m[key] = bson.M{
		"$lt": lt,
	}
	return s
}

//{ field: { $lte: value} } <=
func (s *Session) Lte(key string, lte interface{}) *Session {
	s.m[key] = bson.M{
		"$lte": lte,
	}
	return s
}

//{field: {$ne: value} } !=
func (s *Session) Ne(key string, ne interface{}) *Session {
	s.m[key] = bson.M{
		"$ne": ne,
	}
	return s
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (s *Session) Nin(key string, nin interface{}) *Session {
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
func (s *Session) And(filter Session) *Session {
	s.m["$and"] = filter.a()
	return s

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (s *Session) Not(key string, not interface{}) *Session {
	s.m[key] = bson.M{
		"$not": not,
	}
	return s
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (s *Session) Nor(filter Session) *Session {
	s.m["$nor"] = filter.a()
	return s
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (s *Session) Or(filter Session) *Session {
	s.m["$or"] = filter.a()
	return s
}

func (s *Session) Exists(key string, exists bool, filter ...*Session) *Session {
	m := bson.M{
		"$exists": exists,
	}
	for _, v := range filter {
		for fk, fv := range v.m {
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
func (s *Session) Type(key string, t interface{}) *Session {
	s.m[key] = bson.M{
		"$type": t,
	}
	return s
}

//Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (s *Session) Expr(filter Session) *Session {
	s.m["$expr"] = filter.a()
	return s
}

//todo 简单实现，后续增加支持
func (s *Session) Regex(key string, value interface{}) *Session {
	s.m[key] = bson.M{
		"$regex":   value,
		"$options": "i",
	}

	return s
}

////todo 未实现 ,玩不懂。
//func (s *Session) JsonSchema() {
//
//}
//
////{ field: { $mod: [ divisor, remainder ] } }
////Examples
////{ "_id" : 1, "item" : "abc123", "qty" : 0 }
////{ "_id" : 2, "item" : "xyz123", "qty" : 5 }
////{ "_id" : 3, "item" : "ijk123", "qty" : 12 }
////db.inventory.find( { qty: { $mod: [ 4, 0 ] } } )
////The query returns the following documents:
////{ "_id" : 1, "item" : "abc123", "qty" : 0 }
////{ "_id" : 3, "item" : "ijk123", "qty" : 12 }
//func (s *Session) Mod(key string, mod interface{}) *Session {
//	s.e = append(s.e, bson.E{Key: key, Value: bson.D{{"$mod", mod}}})
//	return s
//}
//

//
////todo 比较复杂今晚不加了
//func (s *Session) Text() *Session {
//	panic("比较复杂今晚不加")
//}

func (s *Session) getStructColl(doc interface{}) (*mongo.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := s.engine.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}
	coll := s.engine.Collection(t.Name)
	return coll, nil
}

func (s *Session) getSliceColl(doc interface{}) (*mongo.Collection, error) {
	sliceValue := reflect.Indirect(reflect.ValueOf(doc))

	if sliceValue.Kind() != reflect.Slice && reflect.Map != sliceValue.Kind() {
		return nil, errors.New("needs a pointer to a slice or a map")
	}

	var t *schemas.Collection
	var err error
	if sliceValue.Kind() == reflect.Slice {
		sliceElementType := sliceValue.Type().Elem()
		if sliceElementType.Kind() == reflect.Struct {
			pv := reflect.New(sliceElementType)
			t, err = s.engine.parser.Parse(pv)
		}
	} else {
		t, err = s.engine.parser.Parse(sliceValue)
	}

	if err != nil {
		return nil, err
	}

	coll := s.engine.Collection(t.Name)
	return coll, nil
}

func (s *Session) getStructCollAndSetKey(doc interface{}) (*mongo.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := s.engine.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}
	docTyp := t.Type
	for i := 0; i < docTyp.NumField(); i++ {
		field := docTyp.Field(i)
		if strings.Index(field.Tag.Get("bson"), "_id") > 0 {
			//s.e = append(s.e, Session("_id", beanValue.Field(i).Interface()))
			break
		}
	}

	coll := s.engine.Collection(t.Name)
	return coll, nil
}

func (s *Session) a() []bson.E {
	var fs []bson.E

	for key, value := range s.m {
		fs = append(fs, bson.E{Key: key, Value: value})
	}
	return fs
}
