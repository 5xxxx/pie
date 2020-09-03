/*

Example:

package main

import (
	"context"
	"fmt"
	"time"
	"tugrik"
)

func main() {
	t, err := pie.NewDriver()
	t.SetURI("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	err = t.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	t.SetDatabase("xxxx")
	var user User
	err = t.filter("nickName", "淳朴的润土").FindOne(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
*/

package pie

import (
	"context"
	"errors"
	"github.com/NSObjects/pie/names"
	"github.com/NSObjects/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"strings"
)

type Driver struct {
	client     *mongo.Client
	parser     *Parser
	db         string
	clientOpts []*options.ClientOptions
}

func NewDriver(opts ...*options.ClientOptions) (*Driver, error) {
	mapper := names.NewCacheMapper(new(names.SnakeMapper))
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	parser := NewParser(mapper, mapper)
	driver := &Driver{
		clientOpts: opts,
		parser:     parser,
		client:     client,
	}
	return driver, nil
}

func (d *Driver) Connect(ctx context.Context) (err error) {
	d.client, err = mongo.Connect(ctx, d.clientOpts...)
	return err
}

func (d *Driver) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

func (d *Driver) BulkWrite(ctx context.Context, docs interface{}) (*mongo.BulkWriteResult, error) {
	session := d.NewSession()
	return session.BulkWrite(ctx, docs)
}

func (d *Driver) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	session := d.NewSession()
	return session.Distinct(ctx, doc, columns)
}

func (d *Driver) ReplaceOne(ctx context.Context, doc interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.ReplaceOne(ctx, doc)
}

func (d *Driver) FindOneAndReplace(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOneAndReplace(ctx, doc)
}

func (d *Driver) FindOneAndUpdate(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOneAndUpdate(ctx, doc)
}

func (d *Driver) FindAndDelete(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindAndDelete(ctx, doc)
}

// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
func (d *Driver) FindOne(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOne(ctx, doc)
}

func (d *Driver) FindAll(ctx context.Context, docs interface{}) error {
	session := d.NewSession()
	return session.FindAll(ctx, docs)
}

func (d *Driver) RegexFilter(key, pattern string) *Session {
	session := d.NewSession()
	return session.RegexFilter(key, pattern)
}

func (d *Driver) Asc(colNames ...string) *Session {
	session := d.NewSession()
	return session.Asc(colNames...)
}

func (d *Driver) Eq(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Eq(key, value)
}

func (d *Driver) Ne(key string, ne interface{}) *Session {
	session := d.NewSession()
	return session.Gt(key, ne)
}

func (d *Driver) Nin(key string, nin interface{}) *Session {
	session := d.NewSession()
	return session.Nin(key, nin)
}

func (d *Driver) Nor(c Condition) *Session {
	session := d.NewSession()
	return session.Nor(c)
}

func (d *Driver) Exists(key string, exists bool, filter ...Condition) *Session {
	session := d.NewSession()
	return session.Exists(key, exists, filter...)
}

func (d *Driver) Type(key string, t interface{}) *Session {
	session := d.NewSession()
	return session.Gt(key, t)
}

func (d *Driver) Expr(filter Condition) *Session {
	session := d.NewSession()
	return session.Expr(filter)
}

func (d *Driver) Regex(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Regex(key, value)
}

func (d *Driver) SetDatabase(db string) {
	d.db = db
}

func (d *Driver) DataBase() *mongo.Database {
	return d.client.Database(d.db)
}

func (d *Driver) Collection(name string) *mongo.Collection {
	return d.client.Database(d.db).Collection(name)
}

func (d *Driver) Ping() error {
	return d.client.Ping(context.TODO(), readpref.Primary())
}

func (d *Driver) Filter(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Filter(key, value)
}

func (d *Driver) ID(id interface{}) *Session {
	session := d.NewSession()
	return session.ID(id)
}

func (d *Driver) Gt(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Gt(key, value)
}

func (d *Driver) Gte(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Gte(key, value)
}

func (d *Driver) Lt(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Lt(key, value)
}

func (d *Driver) Lte(key string, value interface{}) *Session {
	session := d.NewSession()
	return session.Lte(key, value)
}

func (d *Driver) In(key string, value interface{}) *Session {
	session := d.NewSession()
	session.In(key, value)
	return session
}

func (d *Driver) And(filter Condition) *Session {
	session := d.NewSession()
	session.And(filter)
	return session
}

func (d *Driver) Not(key string, value interface{}) *Session {
	session := d.NewSession()
	session.Not(key, value)
	return session
}

func (d *Driver) Or(filter Condition) *Session {
	session := d.NewSession()
	session.Or(filter)
	return session
}

func (d *Driver) InsertOne(ctx context.Context, v interface{}) (primitive.ObjectID, error) {
	session := d.NewSession()
	return session.InsertOne(ctx, v)
}

func (d *Driver) InsertMany(ctx context.Context, v interface{}) (*mongo.InsertManyResult, error) {
	session := d.NewSession()
	return session.InsertMany(ctx, v)
}

func (d *Driver) Limit(limit int64) *Session {
	session := d.NewSession()
	return session.Limit(limit)
}

func (d *Driver) Skip(skip int64) *Session {
	session := d.NewSession()
	return session.Limit(skip)
}

func (d *Driver) Count(i interface{}) (int64, error) {
	session := d.NewSession()
	return session.Count(i)
}

func (d *Driver) Desc(s2 ...string) *Session {
	session := d.NewSession()
	return session.Desc(s2...)
}

func (d *Driver) Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.Update(ctx, bean)
}

//The following operation updates all of the documents with quantity value less than 50.
func (d *Driver) UpdateMany(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.UpdateMany(ctx, bean)
}

func (d *Driver) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := d.NewSession()
	return session.DeleteOne(ctx, filter)
}

func (d *Driver) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := d.NewSession()
	return session.DeleteMany(ctx, filter)
}
func (d *Driver) SoftDeleteOne(ctx context.Context, filter interface{}) error {
	session := d.NewSession()
	return session.SoftDeleteOne(ctx, filter)
}

func (d *Driver) SoftDeleteMany(ctx context.Context, filter interface{}) error {
	session := d.NewSession()
	return session.SoftDeleteMany(ctx, filter)
}

func (d *Driver) FilterBy(object interface{}) *Session {
	session := d.NewSession()
	return session.FilterBy(object)
}

func (d *Driver) DropAll(ctx context.Context, doc interface{}) error {
	indexes := d.NewIndexes()
	return indexes.DropAll(ctx, doc)
}

func (d *Driver) DropOne(ctx context.Context, doc interface{}, name string) error {
	indexes := d.NewIndexes()
	return indexes.DropOne(ctx, doc, name)
}

func (d *Driver) AddIndex(keys interface{}, opt ...*options.IndexOptions) *Indexes {
	indexes := d.NewIndexes()
	return indexes.AddIndex(keys, opt...)
}

func (d *Driver) NewIndexes() *Indexes {
	return NewIndexes(d)
}

func (d *Driver) NewSession() *Session {
	return NewSession(d)
}

func (d *Driver) Aggregate() *Aggregate {
	return NewAggregate(d)
}

func (d *Driver) CollectionNameForStruct(doc interface{}) (*schemas.Collection, error) {
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

func (d *Driver) CollectionNameForSlice(doc interface{}) (*schemas.Collection, error) {
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

func (d *Driver) getStructCollAndSetKey(doc interface{}) (*schemas.Collection, error) {
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
			//d.e = append(d.e, Session("_id", beanValue.Field(i).Interface()))
			break
		}
	}

	return t, nil
}
