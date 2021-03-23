/*
 *
 * find_one_session.go
 * pie
 *
 * Created by lintao on 2020/7/17 10:16 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package internal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/NSObjects/pie/driver"
	"github.com/NSObjects/pie/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type session struct {
	db                    string
	engine                driver.Client
	filter                driver.Condition
	findOneOptions        []*options.FindOneOptions
	findOptions           []*options.FindOptions
	insertManyOpts        []*options.InsertManyOptions
	insertOneOpts         []*options.InsertOneOptions
	deleteOpts            []*options.DeleteOptions
	findOneAndDeleteOpts  []*options.FindOneAndDeleteOptions
	updateOpts            []*options.UpdateOptions
	countOpts             []*options.CountOptions
	distinctOpts          []*options.DistinctOptions
	findOneAndReplaceOpts []*options.FindOneAndReplaceOptions
	findOneAndUpdateOpts  []*options.FindOneAndUpdateOptions
	replaceOpts           []*options.ReplaceOptions
	bulkWriteOptions      []*options.BulkWriteOptions
}

func (s *session) FilterBson(d bson.D) driver.Session {
	s.filter.FilterBson(d)
	return s
}

func NewSession(engine driver.Client) driver.Session {
	return &session{engine: engine, filter: DefaultCondition()}
}

func (s *session) FindPagination(page, count int64, rowsSlicePtr interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForSlice(rowsSlicePtr)
	if err != nil {
		return err
	}
	if page == 0 {
		return errors.New("page must be greater that 0")
	}

	if count == 0 {
		count = 10
	}

	s.Skip((page - 1) * count).Limit(count)
	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}

	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	cursor, err := coll.Find(c, filters, s.findOptions...)
	if err != nil {
		return err
	}

	if err = cursor.All(c, rowsSlicePtr); err != nil {
		return err
	}
	return nil
}

func (s *session) BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}
	values := reflect.ValueOf(docs)
	var mods []mongo.WriteModel
	for i := 0; i < values.Len(); i++ {
		mods = append(mods, mongo.NewInsertOneModel().SetDocument(values.Index(i).Interface()))
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return coll.BulkWrite(c, mods, s.bulkWriteOptions...)
}

func (s *session) FilterBy(object interface{}) driver.Session {
	s.filter.FilterBy(object)
	return s
}

func (s *session) Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error) {
	coll, err := s.collectionForSlice(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return coll.Distinct(c, columns, filters, s.distinctOpts...)
}

func (s *session) ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}

	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return coll.ReplaceOne(c, filters, doc, s.replaceOpts...)
}

func (s *session) FindOneAndReplace(doc interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}

	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return coll.FindOneAndReplace(c, filters, doc, s.findOneAndReplaceOpts...).Decode(&doc)
}

func (s *session) FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {
	c, err := s.collectionForStruct(coll)
	if err != nil {
		return nil, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}

	cc := context.Background()
	if len(ctx) > 0 {
		cc = ctx[0]
	}
	return c.FindOneAndUpdate(cc, filters, bson, s.findOneAndUpdateOpts...), nil
}

func (s *session) FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {

	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.FindOneAndUpdate(c, filters, bson.M{"$set": doc}, s.findOneAndUpdateOpts...), nil
}

func (s *session) FindAndDelete(doc interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.FindOneAndDelete(c, filters, s.findOneAndDeleteOpts...).Decode(&doc)
}

// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
func (s *session) FindOne(doc interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	result := coll.FindOne(c, filters, s.findOneOptions...)
	if err = result.Err(); err != nil {
		return err
	}

	if err = result.Decode(doc); err != nil {
		return err
	}

	return nil
}

// Find executes a find command and returns a Cursor over the matching documents in the collectionByName.
func (s *session) FindAll(rowsSlicePtr interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForSlice(rowsSlicePtr)
	if err != nil {
		return err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	cursor, err := coll.Find(c, filters, s.findOptions...)
	if err != nil {
		return err
	}

	if err = cursor.All(c, rowsSlicePtr); err != nil {
		return err
	}

	return nil
}

// InsertOne executes an insert command to insert a single document into the collectionByName.
func (s *session) InsertOne(doc interface{}, ctx ...context.Context) (primitive.ObjectID, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return [12]byte{}, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	result, err := coll.InsertOne(c, doc, s.insertOneOpts...)
	if err != nil {
		return [12]byte{}, err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, err
	}
	return [12]byte{}, err
}

// InsertMany executes an insert command to insert multiple documents into the collectionByName.
func (s *session) InsertMany(docs interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(docs)
	var many []interface{}
	for index := 0; index < value.Len(); index++ {
		many = append(many, value.Index(index).Interface())
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.InsertMany(c, many, s.insertManyOpts...)
}

// DeleteOne executes a delete command to delete at most one document from the collectionByName.
func (s *session) DeleteOne(doc interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.DeleteOne(c, filters, s.deleteOpts...)
}

func (s *session) SoftDeleteOne(doc interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	_, err = coll.UpdateOne(c, filters, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})

	return err
}

// DeleteMany executes a delete command to delete documents from the collectionByName.
func (s *session) DeleteMany(doc interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.DeleteMany(c, filters, s.deleteOpts...)
}

func (s *session) SoftDeleteMany(doc interface{}, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	_, err = coll.UpdateMany(c, filters, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})

	return err
}

func (s *session) Clone() driver.Session {
	var sess = *s
	return &sess
}

func (s *session) Limit(i int64) driver.Session {
	s.findOptions = append(s.findOptions, options.Find().SetLimit(i))
	return s
}

func (s *session) Skip(i int64) driver.Session {
	s.findOptions = append(s.findOptions, options.Find().SetSkip(i))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSkip(i))
	return s
}

func (s *session) Count(i interface{}, ctx ...context.Context) (int64, error) {
	kind := reflect.TypeOf(i).Kind()
	if kind == reflect.Ptr {
		kind = reflect.TypeOf(reflect.Indirect(reflect.ValueOf(i)).Interface()).Kind()
	}
	var coll *mongo.Collection
	var err error
	switch kind {
	case reflect.Slice:
		coll, err = s.collectionForSlice(i)
	case reflect.Struct:
		coll, err = s.collectionForStruct(i)
	default:
		return 0, errors.New("neet slice or struct")
	}

	if err != nil {
		return 0, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return 0, err
	}

	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return coll.CountDocuments(c, filters, s.countOpts...)
}

func (s *session) UpdateOne(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForStruct(bean)

	if err != nil {
		return nil, err
	}

	if utils.IsStructZero(reflect.ValueOf(bean).Elem()) {
		return nil, nil
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.UpdateOne(c, filters, bson.M{"$set": bean}, s.updateOpts...)
}

func (s *session) UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	c, err := s.collectionForStruct(coll)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	cc := context.Background()
	if len(ctx) > 0 {
		cc = ctx[0]
	}
	return c.UpdateOne(cc, filters, bson, s.updateOpts...)
}

func (s *session) UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	c, err := s.collectionForStruct(coll)
	if err != nil {
		return nil, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	cc := context.Background()
	if len(ctx) > 0 {
		cc = ctx[0]
	}
	return c.UpdateMany(cc, filters, bson, s.updateOpts...)
}

func (s *session) toBson(obj interface{}) bson.M {
	beanValue := reflect.ValueOf(obj).Elem()
	if beanValue.Kind() != reflect.Struct ||
		//Todo how to fix array?
		beanValue.Kind() == reflect.Array {
		panic(errors.New("needs a struct"))
	}

	ret := bson.M{}
	docType := reflect.TypeOf(obj).Elem()
	for index := 0; index < docType.NumField(); index++ {
		fieldTag := docType.Field(index).Tag.Get("bson")
		if fieldTag != "" && fieldTag != "-" {
			s.makeValue(fieldTag, beanValue.Field(index).Interface(), ret)
		}
	}
	return ret
}

func (s *session) makeValue(field string, value interface{}, ret bson.M) {
	split := strings.Split(field, ",")
	if len(split) <= 0 {
		return
	}
	if strings.Contains(field, "omitempty") {
		if utils.IsZero(value) {
			return
		}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Struct:
		s.makeStruct(field, v, ret)
		return
	case reflect.Array:
		return
	}
	ret[split[0]] = value
}

func (s *session) makeStruct(field string, value reflect.Value, ret bson.M) {
	for index := 0; index < value.NumField(); index++ {
		docType := reflect.TypeOf(value.Interface())
		tag := docType.Field(index).Tag.Get("bson")
		split := strings.Split(tag, ",")
		if len(split) > 0 {
			if tag != "" {
				if strings.Contains(tag, "omitempty") {
					if !utils.IsZero(value.Field(index)) {
						fieldTags := fmt.Sprintf("%s.%s", field, split[0])
						s.makeValue(fieldTags, value.Field(index).Interface(), ret)
					}
				} else {
					fieldTags := fmt.Sprintf("%s.%s", field, split[0])
					s.makeValue(fieldTags, value.Field(index).Interface(), ret)
				}

			}
		}
	}
}

func (s *session) UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForSlice(bean)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return coll.UpdateMany(c, filters, bson.M{"$set": bean}, s.updateOpts...)

}

func (s *session) RegexFilter(key, pattern string) driver.Session {
	s.filter.RegexFilter(key, pattern)
	return s
}

func (s *session) ID(id interface{}) driver.Session {
	s.filter.ID(id)
	return s
}

func (s *session) Asc(colNames ...string) driver.Session {
	if len(colNames) == 0 {
		return s
	}

	es := bson.M{}
	for _, c := range colNames {
		es[c] = -1
	}
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	return s
}

func (s *session) Desc(colNames ...string) driver.Session {
	if len(colNames) == 0 {
		return s
	}

	es := bson.M{}
	for _, c := range colNames {
		es[c] = 1
	}

	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	return s
}

func (s *session) Sort(colNames ...string) driver.Session {
	if len(colNames) == 0 {
		return s
	}
	es := bson.D{}
	for _, field := range colNames {
		if field != "" {
			switch field[0] {
			case '-':
				es = append(es, bson.E{field[1:], -1})
			default:
				es = append(es, bson.E{field, 1})
			}
		}
	}
	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	return s
}

func (s *session) Filter(key string, value interface{}) driver.Session {
	return s.Eq(key, value)
}

//Equals a Specified Value
//{ qty: 20 }
//Field in Embedded Document Equals a Value
//{"item.name": "ab" }
// Equals an Array Value
//{ tags: [ "A", "B" ] }
func (s *session) Eq(key string, value interface{}) driver.Session {
	s.filter.Eq(key, value)
	return s
}

//{field: {$gt: value} } >
func (s *session) Gt(key string, gt interface{}) driver.Session {
	s.filter.Gt(key, gt)
	return s
}

//{ qty: { $gte: 20 } } >=
func (s *session) Gte(key string, gte interface{}) driver.Session {
	s.filter.Gte(key, gte)
	return s
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (s *session) In(key string, in interface{}) driver.Session {
	s.filter.In(key, in)
	return s
}

//{field: {$lt: value} } <
func (s *session) Lt(key string, lt interface{}) driver.Session {
	s.filter.Lt(key, lt)
	return s
}

//{ field: { $lte: value} } <=
func (s *session) Lte(key string, lte interface{}) driver.Session {
	s.filter.Lte(key, lte)
	return s
}

//{field: {$ne: value} } !=
func (s *session) Ne(key string, ne interface{}) driver.Session {
	s.filter.Ne(key, ne)
	return s
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (s *session) Nin(key string, nin interface{}) driver.Session {
	s.filter.Nin(key, nin)
	return s
}

//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
//$and: [
//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
// ]
func (s *session) And(c driver.Condition) driver.Session {
	s.filter.And(c)
	return s

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (s *session) Not(key string, not interface{}) driver.Session {
	s.filter.Not(key, not)
	return s
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (s *session) Nor(c driver.Condition) driver.Session {
	s.filter.Nor(c)
	return s
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (s *session) Or(c driver.Condition) driver.Session {
	s.filter.Or(c)
	return s
}

func (s *session) Exists(key string, exists bool, filter ...driver.Condition) driver.Session {
	s.filter.Exists(key, exists, filter...)
	return s
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (s *session) SetArrayFilters(filters options.ArrayFilters) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetArrayFilters(filters))
	s.updateOpts = append(s.updateOpts, options.Update().SetArrayFilters(filters))
	return s
}

// SetOrdered sets the value for the Ordered field.
func (s *session) SetOrdered(ordered bool) driver.Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetOrdered(ordered))
	return s
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (s *session) SetBypassDocumentValidation(b bool) driver.Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetBypassDocumentValidation(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetBypassDocumentValidation(b))
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts, options.FindOneAndUpdate().SetBypassDocumentValidation(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetBypassDocumentValidation(b))

	return s
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (s *session) SetReturnDocument(rd options.ReturnDocument) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetReturnDocument(rd))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetReturnDocument(rd))
	return s
}

// SetUpsert sets the value for the Upsert field.
func (s *session) SetUpsert(b bool) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetUpsert(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetUpsert(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetUpsert(b))
	return s
}

// SetCollation sets the value for the Collation field.
func (s *session) SetCollation(collation *options.Collation) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetCollation(collation))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetCollation(collation))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetCollation(collation))
	s.updateOpts = append(s.updateOpts, options.Update().SetCollation(collation))
	return s
}

// SetMaxTime sets the value for the MaxTime field.
func (s *session) SetMaxTime(d time.Duration) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetMaxTime(d))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetMaxTime(d))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetMaxTime(d))
	return s
}

// SetProjection sets the value for the Projection field.
func (s *session) SetProjection(projection interface{}) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetProjection(projection))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetProjection(projection))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetProjection(projection))
	return s
}

// SetSort sets the value for the Sort field.
func (s *session) SetSort(sort interface{}) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetSort(sort))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetSort(sort))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetSort(sort))
	return s
}

// SetHint sets the value for the Hint field.
func (s *session) SetHint(hint interface{}) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetHint(hint))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetHint(hint))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetHint(hint))
	s.updateOpts = append(s.updateOpts, options.Update().SetHint(hint))
	return s
}

//{ field: { $type: <BSON type> } }
// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
func (s *session) Type(key string, t interface{}) driver.Session {
	s.filter.Type(key, t)
	return s
}

//Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (s *session) Expr(c driver.Condition) driver.Session {
	s.filter.Expr(c)
	return s
}

//todo 简单实现，后续增加支持
func (s *session) Regex(key string, value interface{}) driver.Session {
	s.filter.Regex(key, value)
	return s
}

func (c *session) SetDatabase(db string) driver.Session {
	c.db = db
	return c
}

func (c *session) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
	coll, err := c.engine.CollectionNameForStruct(doc)
	if err != nil {
		return nil, err
	}

	return c.collectionByName(coll.Name), nil
}

func (c *session) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
	coll, err := c.engine.CollectionNameForSlice(doc)
	if err != nil {
		return nil, err
	}
	return c.collectionByName(coll.Name), nil
}

func (c *session) collectionByName(name string) *mongo.Collection {
	return c.engine.Collection(name, c.db)
}

//func (s *session) makeFilterValue(field string, value interface{}) {
//	if utils.IsZero(value) {
//		return
//	}
//	v := reflect.ValueOf(value)
//	switch v.Kind() {
//	case reflect.Struct:
//		s.makeStructValue(field, v)
//	case reflect.Array:
//		return
//	}
//	s.Filter(field, value)
//}
//
//func (s *session) makeStructValue(field string, value reflect.Value) {
//	for index := 0; index < value.NumField(); index++ {
//		docType := reflect.TypeOf(value.Interface())
//		tag := docType.Field(index).Tag.Get("bson")
//		if tag != "" {
//			if !utils.IsZero(value.Field(index)) {
//				fieldTags := fmt.Sprintf("%s.%s", field, tag)
//				s.makeFilterValue(fieldTags, value.Field(index).Interface())
//			}
//		}
//	}
//}
