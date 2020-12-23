/*

Example:

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NSObjects/pie"
)

func main() {
	t, err := pie.NewClient("demo")
	t.SetURI("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	err = t.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	var user User
	err = t.filter("nickName", "淳朴的润土").FindOne(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
*/

package driver

import (
	"context"
	"errors"
	"github.com/NSObjects/pie"
	"github.com/NSObjects/pie/names"
	"github.com/NSObjects/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"strings"
)

type driver struct {
	client     *mongo.Client
	parser     *pie.Parser
	db         string
	clientOpts []*options.ClientOptions
}

func NewClient(db string, opts ...*options.ClientOptions) (pie.Client, error) {
	mapper := names.NewCacheMapper(new(names.SnakeMapper))
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	parser := pie.NewParser(mapper, mapper)
	driver := &driver{
		clientOpts: opts,
		parser:     parser,
		client:     client,
		db:         db,
	}
	return driver, nil
}

func (d *driver) Connect(ctx context.Context) (err error) {
	d.client, err = mongo.Connect(ctx, d.clientOpts...)
	return err
}

func (d *driver) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

func (d *driver) BulkWrite(ctx context.Context, docs interface{}) (*mongo.BulkWriteResult, error) {
	session := d.NewSession()
	return session.BulkWrite(ctx, docs)
}

func (d *driver) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	session := d.NewSession()
	return session.Distinct(ctx, doc, columns)
}

func (d *driver) ReplaceOne(ctx context.Context, doc interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.ReplaceOne(ctx, doc)
}

func (d *driver) FindOneAndReplace(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOneAndReplace(ctx, doc)
}

func (d *driver) FindOneAndUpdate(ctx context.Context, doc interface{}) (*mongo.SingleResult, error) {
	session := d.NewSession()
	return session.FindOneAndUpdate(ctx, doc)
}

func (d *driver) FindAndDelete(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindAndDelete(ctx, doc)
}

// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
func (d *driver) FindOne(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOne(ctx, doc)
}

func (d *driver) FindAll(ctx context.Context, docs interface{}) error {
	session := d.NewSession()
	return session.FindAll(ctx, docs)
}

func (d *driver) RegexFilter(key, pattern string) pie.Session {
	session := d.NewSession()
	return session.RegexFilter(key, pattern)
}

func (d *driver) Asc(colNames ...string) pie.Session {
	session := d.NewSession()
	return session.Asc(colNames...)
}

func (d *driver) Eq(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Eq(key, value)
}

func (d *driver) Ne(key string, ne interface{}) pie.Session {
	session := d.NewSession()
	return session.Gt(key, ne)
}

func (d *driver) Nin(key string, nin interface{}) pie.Session {
	session := d.NewSession()
	return session.Nin(key, nin)
}

func (d *driver) Nor(c pie.Condition) pie.Session {
	session := d.NewSession()
	return session.Nor(c)
}

func (d *driver) Exists(key string, exists bool, filter ...pie.Condition) pie.Session {
	session := d.NewSession()
	return session.Exists(key, exists, filter...)
}

func (d *driver) Type(key string, t interface{}) pie.Session {
	session := d.NewSession()
	return session.Gt(key, t)
}

func (d *driver) Expr(filter pie.Condition) pie.Session {
	session := d.NewSession()
	return session.Expr(filter)
}

func (d *driver) Regex(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Regex(key, value)
}

func (d *driver) DataBase() *mongo.Database {
	return d.client.Database(d.db)
}

func (d *driver) Collection(name string) *mongo.Collection {
	return d.client.Database(d.db).Collection(name)
}

func (d *driver) Ping() error {
	return d.client.Ping(context.TODO(), readpref.Primary())
}

func (d *driver) Filter(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Filter(key, value)
}

func (d *driver) ID(id interface{}) pie.Session {
	session := d.NewSession()
	return session.ID(id)
}

func (d *driver) Gt(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Gt(key, value)
}

func (d *driver) Gte(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Gte(key, value)
}

func (d *driver) Lt(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Lt(key, value)
}

func (d *driver) Lte(key string, value interface{}) pie.Session {
	session := d.NewSession()
	return session.Lte(key, value)
}

func (d *driver) In(key string, value interface{}) pie.Session {
	session := d.NewSession()
	session.In(key, value)
	return session
}

func (d *driver) And(filter pie.Condition) pie.Session {
	session := d.NewSession()
	session.And(filter)
	return session
}

func (d *driver) Not(key string, value interface{}) pie.Session {
	session := d.NewSession()
	session.Not(key, value)
	return session
}

func (d *driver) Or(filter pie.Condition) pie.Session {
	session := d.NewSession()
	session.Or(filter)
	return session
}

func (d *driver) InsertOne(ctx context.Context, v interface{}) (primitive.ObjectID, error) {
	session := d.NewSession()
	return session.InsertOne(ctx, v)
}

func (d *driver) InsertMany(ctx context.Context, v interface{}) (*mongo.InsertManyResult, error) {
	session := d.NewSession()
	return session.InsertMany(ctx, v)
}

func (d *driver) Limit(limit int64) pie.Session {
	session := d.NewSession()
	return session.Limit(limit)
}

func (d *driver) Skip(skip int64) pie.Session {
	session := d.NewSession()
	return session.Limit(skip)
}

func (d *driver) Count(i interface{}) (int64, error) {
	session := d.NewSession()
	return session.Count(i)
}

func (d *driver) Desc(s2 ...string) pie.Session {
	session := d.NewSession()
	return session.Desc(s2...)
}

func (d *driver) Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.Update(ctx, bean)
}

//The following operation updates all of the documents with quantity value less than 50.
func (d *driver) UpdateMany(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.UpdateMany(ctx, bean)
}

func (d *driver) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := d.NewSession()
	return session.DeleteOne(ctx, filter)
}

func (d *driver) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := d.NewSession()
	return session.DeleteMany(ctx, filter)
}
func (d *driver) SoftDeleteOne(ctx context.Context, filter interface{}) error {
	session := d.NewSession()
	return session.SoftDeleteOne(ctx, filter)
}

func (d *driver) SoftDeleteMany(ctx context.Context, filter interface{}) error {
	session := d.NewSession()
	return session.SoftDeleteMany(ctx, filter)
}

func (d *driver) FilterBy(object interface{}) pie.Session {
	session := d.NewSession()
	return session.FilterBy(object)
}

func (d *driver) DropAll(ctx context.Context, doc interface{}) error {
	indexes := d.NewIndexes()
	return indexes.DropAll(ctx, doc)
}

func (d *driver) DropOne(ctx context.Context, doc interface{}, name string) error {
	indexes := d.NewIndexes()
	return indexes.DropOne(ctx, doc, name)
}

func (d *driver) AddIndex(keys interface{}, opt ...*options.IndexOptions) pie.Indexes {
	indexes := d.NewIndexes()
	return indexes.AddIndex(keys, opt...)
}

func (d *driver) NewIndexes() pie.Indexes {
	return NewIndexes(d)
}

func (d *driver) NewSession() pie.Session {
	return NewSession(d)
}

func (d *driver) Aggregate() pie.Aggregate {
	return NewAggregate(d)
}

func (d *driver) CollectionNameForStruct(doc interface{}) (*schemas.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := d.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (d *driver) SetDatabase(string string) pie.Client {
	d.db = string
	return d
}

func (d *driver) CollectionNameForSlice(doc interface{}) (*schemas.Collection, error) {
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
			t, err = d.parser.Parse(pv)
		}
	} else {
		t, err = d.parser.Parse(sliceValue)
	}

	if err != nil {
		return nil, err
	}
	return t, nil
}

func (d *driver) getStructCollAndSetKey(doc interface{}) (*schemas.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := d.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}
	docTyp := t.Type
	for i := 0; i < docTyp.NumField(); i++ {
		field := docTyp.Field(i)
		if strings.Index(field.Tag.Get("bson"), "_id") > 0 {
			//d.e = append(d.e, session("_id", beanValue.Field(i).Interface()))
			break
		}
	}

	return t, nil
}
