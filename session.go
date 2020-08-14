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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	filter         Condition
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

func NewSession(engine *Tugrik) *Session {
	return &Session{engine: engine, filter: DefaultCondition()}
}

func (s *Session) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	coll, err := s.engine.getSliceColl(doc)
	if err != nil {
		return nil, err
	}
	return coll.Distinct(ctx, columns, s.filter.Filters(), s.distinctOpts...)
}

// FindOne executes a find command and returns a SingleResult for one document in the collection.
func (s *Session) FindOne(ctx context.Context, doc interface{}) error {
	coll, err := s.engine.getStructColl(doc)
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

// Find executes a find command and returns a Cursor over the matching documents in the collection.
func (s *Session) FindAll(ctx context.Context, rowsSlicePtr interface{}) error {
	coll, err := s.engine.getSliceColl(rowsSlicePtr)
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

// InsertOne executes an insert command to insert a single document into the collection.
func (s *Session) InsertOne(ctx context.Context, doc interface{}) error {
	coll, err := s.engine.getStructColl(doc)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, doc, s.insertOneOpts...)
	return err
}

// InsertMany executes an insert command to insert multiple documents into the collection.
func (s *Session) InsertMany(ctx context.Context, docs []interface{}) error {
	coll, err := s.engine.getSliceColl(docs)
	if err != nil {
		return err
	}
	_, err = coll.InsertMany(ctx, docs, s.insertManyOpts...)

	return err
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (s *Session) DeleteOne(ctx context.Context, doc interface{}) error {
	coll, err := s.engine.getStructColl(doc)
	if err != nil {
		return err
	}
	_, err = coll.DeleteOne(ctx, s.filter.Filters(), s.deleteOpts...)
	return err
}

// DeleteMany executes a delete command to delete documents from the collection.
func (s *Session) DeleteMany(ctx context.Context, doc interface{}) error {
	coll, err := s.engine.getStructColl(doc)
	if err != nil {
		return err
	}
	_, err = coll.DeleteMany(ctx, s.filter.Filters(), s.deleteOpts...)
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
	coll, err := s.engine.getStructColl(i)
	if err != nil {
		return 0, err
	}
	return coll.CountDocuments(context.Background(), s.filter.Filters(), s.countOpts...)
}

func (s *Session) Update(bean interface{}) error {
	coll, err := s.engine.getStructColl(bean)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(context.Background(), s.filter.Filters(), bson.M{"$set": insertOmitemptyTag(bean)})

	return err
}

//todo update many
func (s *Session) UpdateMany(bean interface{}) error {
	coll, err := s.engine.getSliceColl(bean)
	if err != nil {
		return err
	}
	_, err = coll.UpdateMany(context.Background(), s.filter.Filters(), bson.M{"$set": bean})

	return err
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
