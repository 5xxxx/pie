/*
 *
 * find_one_session.go
 * tugrik
 *
 * Created by lintao on 2020/7/17 10:16 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/NSObjects/pie/schemas"

	"github.com/NSObjects/pie/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	db                    string
	filter                Condition
	engine                *Driver
	doc                   interface{}
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

func NewSession(engine *Driver) *Session {
	return &Session{engine: engine, filter: DefaultCondition()}
}

func (s *Session) BulkWrite(ctx context.Context, docs interface{}) (*mongo.BulkWriteResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}
	values := reflect.ValueOf(docs)
	var mods []mongo.WriteModel
	for i := 0; i < values.Len(); i++ {
		mods = append(mods, mongo.NewInsertOneModel().SetDocument(docs))
	}

	return coll.BulkWrite(ctx, mods, s.bulkWriteOptions...)
}

func (s *Session) FilterBy(object interface{}) *Session {

	if m, ok := object.(bson.M); ok {
		for key, value := range m {
			s.Filter(key, value)
		}
		return s
	}

	if d, ok := object.(bson.D); ok {
		for _, v := range d {
			s.Filter(v.Key, v.Value)
		}
		return s
	}

	beanValue := reflect.ValueOf(object)
	if beanValue.Kind() != reflect.Struct ||
		//Todo how to fix array?
		beanValue.Kind() == reflect.Array {
		panic(errors.New("needs a struct"))
	}

	docType := reflect.TypeOf(object)
	for index := 0; index < docType.NumField(); index++ {
		fieldTag := docType.Field(index).Tag.Get("filter")
		if fieldTag != "" && fieldTag != "-" {
			split := strings.Split(fieldTag, ",")
			if len(split) > 0 {
				s.makeFilterValue(split[0], beanValue.Field(index).Interface())
			}
		}
	}

	return s
}

func (s *Session) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	coll, err := s.collectionForSlice(doc)
	if err != nil {
		return nil, err
	}
	return coll.Distinct(ctx, columns, s.filter.Filters(), s.distinctOpts...)
}

func (s *Session) ReplaceOne(ctx context.Context, doc interface{}) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}
	return coll.ReplaceOne(ctx, s.filter.Filters(), doc, s.replaceOpts...)
}

func (s *Session) FindOneAndReplace(ctx context.Context, doc interface{}) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	return coll.FindOneAndReplace(ctx, s.filter.Filters(), doc, s.findOneAndReplaceOpts...).Decode(&doc)
}

func (s *Session) FindOneAndUpdate(ctx context.Context, doc interface{}) error {
	var coll *mongo.Collection
	var err error
	if m, ok := doc.(bson.M); ok {
		if s.doc == nil {
			return errors.New("update bson.M must call Collection(struct{}) first")
		}
		coll, err = s.collectionForStruct(s.doc)
		if err != nil {
			return err
		}
		return coll.FindOneAndUpdate(ctx, s.filter.Filters(), m, s.findOneAndUpdateOpts...).Decode(&doc)
	} else if d, ok := doc.(bson.D); ok {
		if s.doc == nil {
			return errors.New("update bson.D must call Collection(struct{}) first")
		}
		coll, err = s.collectionForStruct(s.doc)
		if err != nil {
			return err
		}

		return coll.FindOneAndUpdate(ctx, s.filter.Filters(), d, s.findOneAndUpdateOpts...).Decode(&doc)
	}

	coll, err = s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	return coll.FindOneAndUpdate(ctx, s.filter.Filters(), doc, s.findOneAndUpdateOpts...).Decode(&doc)
}

func (s *Session) FindAndDelete(ctx context.Context, doc interface{}) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}
	return coll.FindOneAndDelete(ctx, s.filter.Filters(), s.findOneAndDeleteOpts...).Decode(&doc)
}

// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
func (s *Session) FindOne(ctx context.Context, doc interface{}) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	result := coll.FindOne(ctx, s.filter.Filters(), s.findOneOptions...)
	if err = result.Err(); err != nil {
		return err
	}

	if err = result.Decode(doc); err != nil {
		return err
	}

	return nil
}

// Find executes a find command and returns a Cursor over the matching documents in the collectionByName.
func (s *Session) FindAll(ctx context.Context, rowsSlicePtr interface{}) error {
	coll, err := s.collectionForSlice(rowsSlicePtr)
	if err != nil {
		return err
	}
	cursor, err := coll.Find(ctx, s.filter.Filters(), s.findOptions...)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, rowsSlicePtr); err != nil {
		return err
	}

	return nil
}

// InsertOne executes an insert command to insert a single document into the collectionByName.
func (s *Session) InsertOne(ctx context.Context, doc interface{}) (primitive.ObjectID, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return [12]byte{}, err
	}
	result, err := coll.InsertOne(ctx, doc, s.insertOneOpts...)
	if err != nil {
		return [12]byte{}, err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, err
	}
	return [12]byte{}, err
}

// InsertMany executes an insert command to insert multiple documents into the collectionByName.
func (s *Session) InsertMany(ctx context.Context, docs interface{}) (*mongo.InsertManyResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(docs)
	var many []interface{}
	for index := 0; index < value.Len(); index++ {
		many = append(many, value.Index(index).Interface())
	}
	return coll.InsertMany(ctx, many, s.insertManyOpts...)
}

// DeleteOne executes a delete command to delete at most one document from the collectionByName.
func (s *Session) DeleteOne(ctx context.Context, doc interface{}) (*mongo.DeleteResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	return coll.DeleteOne(ctx, s.filter.Filters(), s.deleteOpts...)
}

// DeleteMany executes a delete command to delete documents from the collectionByName.
func (s *Session) DeleteMany(ctx context.Context, doc interface{}) (*mongo.DeleteResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	return coll.DeleteMany(ctx, s.filter.Filters(), s.deleteOpts...)
}

func (s *Session) Collection(doc interface{}) *Session {
	s.doc = doc
	return s
}

func (s *Session) Clone() *Session {
	var sess = *s
	return &sess
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
	return coll.CountDocuments(context.Background(), s.filter.Filters(), s.countOpts...)
}

func (s *Session) Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	var coll *mongo.Collection
	var err error
	if m, ok := bean.(bson.M); ok {
		if s.doc == nil {
			return nil, errors.New("update bson.M must call Collection(struct{}) first")
		}
		coll, err = s.collectionForStruct(s.doc)
		if err != nil {
			return nil, err
		}
		return coll.UpdateOne(ctx, s.filter.Filters(), m, s.updateOpts...)
	} else if d, ok := bean.(bson.D); ok {
		if s.doc == nil {
			return nil, errors.New("update bson.D must call Collection(struct{}) first")
		}
		coll, err = s.collectionForStruct(s.doc)
		if err != nil {
			return nil, err
		}
		return coll.UpdateOne(ctx, s.filter.Filters(), d, s.updateOpts...)
	}

	coll, err = s.collectionForStruct(bean)

	if err != nil {
		return nil, err
	}
	return coll.UpdateOne(ctx, s.filter.Filters(), bson.M{"$set": insertOmitemptyTag(bean)}, s.updateOpts...)
}

func (s *Session) UpdateMany(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	var coll *mongo.Collection
	var err error
	if m, ok := bean.(bson.M); ok {
		if s.doc == nil {
			return nil, errors.New("update bson.M must call Collection(struct{}) first")
		}
		coll, err = s.collectionForStruct(s.doc)
		if err != nil {
			return nil, err
		}
		return coll.UpdateMany(ctx, s.filter.Filters(), m, s.updateOpts...)
	} else if d, ok := bean.(bson.D); ok {
		if s.doc == nil {
			return nil, errors.New("update bson.D must call Collection(struct{}) first")
		}
		coll, err = s.collectionForStruct(s.doc)
		if err != nil {
			return nil, err
		}
		return coll.UpdateMany(ctx, s.filter.Filters(), d, s.updateOpts...)
	}

	coll, err = s.collectionForSlice(bean)
	if err != nil {
		return nil, err
	}
	return coll.UpdateMany(ctx, s.filter.Filters(), bson.M{"$set": bean}, s.updateOpts...)

}

func (s *Session) RegexFilter(key, pattern string) *Session {
	s.filter.RegexFilter(key, pattern)
	return s
}

func (s *Session) ID(id interface{}) *Session {
	s.filter.ID(id)
	return s
}

func (s *Session) Asc(colNames ...string) *Session {
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

func (s *Session) Desc(colNames ...string) *Session {
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
	s.filter.Eq(key, value)
	return s
}

//{field: {$gt: value} } >
func (s *Session) Gt(key string, gt interface{}) *Session {
	s.filter.Gt(key, gt)
	return s
}

//{ qty: { $gte: 20 } } >=
func (s *Session) Gte(key string, gte interface{}) *Session {
	s.filter.Gte(key, gte)
	return s
}

//{ field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (s *Session) In(key string, in interface{}) *Session {
	s.filter.In(key, in)
	return s
}

//{field: {$lt: value} } <
func (s *Session) Lt(key string, lt interface{}) *Session {
	s.filter.Lt(key, lt)
	return s
}

//{ field: { $lte: value} } <=
func (s *Session) Lte(key string, lte interface{}) *Session {
	s.filter.Lte(key, lte)
	return s
}

//{field: {$ne: value} } !=
func (s *Session) Ne(key string, ne interface{}) *Session {
	s.filter.Ne(key, ne)
	return s
}

//{ field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (s *Session) Nin(key string, nin interface{}) *Session {
	s.filter.Nin(key, nin)
	return s
}

//{ $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
//$and: [
//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
// ]
func (s *Session) And(c Condition) *Session {
	s.filter.And(c)
	return s

}

//{ field: { $not: { <operator-expression> } } }
//not and Regular Expressions
//{ item: { $not: /^p.*/ } }
func (s *Session) Not(key string, not interface{}) *Session {
	s.filter.Not(key, not)
	return s
}

// { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (s *Session) Nor(c Condition) *Session {
	s.filter.Nor(c)
	return s
}

// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (s *Session) Or(c Condition) *Session {
	s.filter.Or(c)
	return s
}

func (s *Session) Exists(key string, exists bool, filter ...Condition) *Session {
	s.filter.Exists(key, exists, filter...)
	return s
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (s *Session) SetArrayFilters(filters options.ArrayFilters) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetArrayFilters(filters))
	s.updateOpts = append(s.updateOpts, options.Update().SetArrayFilters(filters))
	return s
}

// SetOrdered sets the value for the Ordered field.
func (s *Session) SetOrdered(ordered bool) *Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetOrdered(ordered))
	return s
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (s *Session) SetBypassDocumentValidation(b bool) *Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetBypassDocumentValidation(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetBypassDocumentValidation(b))
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts, options.FindOneAndUpdate().SetBypassDocumentValidation(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetBypassDocumentValidation(b))

	return s
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (s *Session) SetReturnDocument(rd options.ReturnDocument) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetReturnDocument(rd))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetReturnDocument(rd))
	return s
}

// SetUpsert sets the value for the Upsert field.
func (s *Session) SetUpsert(b bool) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetUpsert(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetUpsert(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetUpsert(b))
	return s
}

// SetCollation sets the value for the Collation field.
func (s *Session) SetCollation(collation *options.Collation) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetCollation(collation))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetCollation(collation))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetCollation(collation))
	s.updateOpts = append(s.updateOpts, options.Update().SetCollation(collation))
	return s
}

// SetMaxTime sets the value for the MaxTime field.
func (s *Session) SetMaxTime(d time.Duration) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetMaxTime(d))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetMaxTime(d))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetMaxTime(d))
	return s
}

// SetProjection sets the value for the Projection field.
func (s *Session) SetProjection(projection interface{}) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetProjection(projection))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetProjection(projection))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetProjection(projection))
	return s
}

// SetSort sets the value for the Sort field.
func (s *Session) SetSort(sort interface{}) *Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetSort(sort))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetSort(sort))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetSort(sort))
	return s
}

// SetHint sets the value for the Hint field.
func (s *Session) SetHint(hint interface{}) *Session {
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
func (s *Session) Type(key string, t interface{}) *Session {
	s.filter.Type(key, t)
	return s
}

//Allows the use of aggregation expressions within the query language.
//{ $expr: { <expression> } }
//$expr can build query expressions that compare fields from the same document in a $match stage
//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
func (s *Session) Expr(c Condition) *Session {
	s.filter.Expr(c)
	return s
}

//todo 简单实现，后续增加支持
func (s *Session) Regex(key string, value interface{}) *Session {
	s.filter.Regex(key, value)
	return s
}

func (s *Session) SetDatabase(db string) *Session {
	s.db = db
	return s
}

func (s *Session) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if s.doc != nil {
		coll, err = s.engine.CollectionNameForStruct(doc)
	}
	coll, err = s.engine.CollectionNameForStruct(doc)
	if err != nil {
		return nil, err
	}
	return s.collectionByName(coll.Name), nil
}

func (s *Session) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if s.doc != nil {
		coll, err = s.engine.CollectionNameForSlice(doc)
	}
	coll, err = s.engine.CollectionNameForSlice(doc)
	if err != nil {
		return nil, err
	}
	return s.collectionByName(coll.Name), nil
}

func (s Session) collectionByName(name string) *mongo.Collection {
	var db string
	if s.db != "" {
		db = s.db
	} else {
		db = s.engine.db
	}
	return s.engine.client.Database(db).Collection(name)
}

func (s *Session) makeFilterValue(field string, value interface{}) {
	if utils.IsZero(value) {
		return
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Struct:
		s.makeStructValue(field, v)
		return
	case reflect.Array:
		return
	}
	s.Filter(field, value)
}

func (s *Session) makeStructValue(field string, value reflect.Value) {
	for index := 0; index < value.NumField(); index++ {
		docType := reflect.TypeOf(value.Interface())
		tag := docType.Field(index).Tag.Get("bson")
		if tag != "" {
			if !utils.IsZero(value.Field(index)) {
				fieldTags := fmt.Sprintf("%s.%s", field, tag)
				s.makeFilterValue(fieldTags, value.Field(index).Interface())
			}
		}
	}
}
