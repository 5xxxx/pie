package internal

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/5xxxx/pie/driver"
	"github.com/5xxxx/pie/names"
	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type defaultClient struct {
	client     *mongo.Client
	parser     *driver.Parser
	db         string
	clientOpts []*options.ClientOptions
}

func NewClient(db string, opts ...*options.ClientOptions) (driver.Client, error) {
	mapper := names.NewCacheMapper(new(names.SnakeMapper))
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	parser := driver.NewParser(mapper, mapper)
	d := defaultClient{
		clientOpts: opts,
		parser:     parser,
		client:     client,
		db:         db,
	}
	return &d, nil
}

func (d *defaultClient) Connect(ctx ...context.Context) (err error) {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	d.client, err = mongo.Connect(c, d.clientOpts...)
	return err
}

func (d *defaultClient) Disconnect(ctx ...context.Context) error {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return d.client.Disconnect(c)
}

func (d *defaultClient) FindPagination(page, count int64, doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindPagination(page, count, doc, ctx...)
}

func (d *defaultClient) BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error) {
	return d.NewSession().BulkWrite(docs, ctx...)
}

func (d *defaultClient) Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error) {
	return d.NewSession().Distinct(doc, columns, ctx...)
}

func (d *defaultClient) ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().ReplaceOne(doc, ctx...)
}

func (d *defaultClient) UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateOneBson(coll, bson, ctx...)
}

func (d *defaultClient) FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {
	return d.NewSession().FindOneAndUpdateBson(coll, bson, ctx...)
}

func (d *defaultClient) UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateManyBson(coll, bson, ctx...)
}

func (d *defaultClient) FindOneAndReplace(doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindOneAndReplace(doc, ctx...)
}

func (d *defaultClient) FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {
	return d.NewSession().FindOneAndUpdate(doc, ctx...)
}

func (d *defaultClient) FindAndDelete(doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindAndDelete(doc, ctx...)
}

// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
func (d *defaultClient) FindOne(doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindOne(doc, ctx...)
}

func (d *defaultClient) FindAll(docs interface{}, ctx ...context.Context) error {
	return d.NewSession().FindAll(docs, ctx...)
}

func (d *defaultClient) FilterBson(x bson.D) driver.Session {
	return d.NewSession().FilterBson(x)
}

func (d *defaultClient) Soft(s bool) driver.Session {
	return d.NewSession().Soft(s)
}

func (d *defaultClient) RegexFilter(key, pattern string) driver.Session {
	return d.NewSession().RegexFilter(key, pattern)
}

func (d *defaultClient) Asc(colNames ...string) driver.Session {
	return d.NewSession().Asc(colNames...)
}

func (d *defaultClient) Eq(key string, value interface{}) driver.Session {
	return d.NewSession().Eq(key, value)
}

func (d *defaultClient) Ne(key string, ne interface{}) driver.Session {
	return d.NewSession().Gt(key, ne)
}

func (d *defaultClient) Nin(key string, nin interface{}) driver.Session {
	return d.NewSession().Nin(key, nin)
}

func (d *defaultClient) Nor(c driver.Condition) driver.Session {
	return d.NewSession().Nor(c)
}

func (d *defaultClient) Exists(key string, exists bool, filter ...driver.Condition) driver.Session {
	return d.NewSession().Exists(key, exists, filter...)
}

func (d *defaultClient) Type(key string, t interface{}) driver.Session {
	return d.NewSession().Gt(key, t)
}

func (d *defaultClient) Expr(filter driver.Condition) driver.Session {
	return d.NewSession().Expr(filter)
}

func (d *defaultClient) Regex(key string, value string) driver.Session {
	return d.NewSession().Regex(key, value)
}

func (d *defaultClient) DataBase() *mongo.Database {
	return d.client.Database(d.db)
}

func (d *defaultClient) Collection(name string, collOpts []*options.CollectionOptions, db ...string) *mongo.Collection {
	var database = d.db
	if len(db) > 0 && len(db[0]) > 0 {
		database = db[0]
	}

	return d.client.Database(database).Collection(name, collOpts...)
}

func (d *defaultClient) Ping() error {
	return d.client.Ping(context.TODO(), readpref.Primary())
}

func (d *defaultClient) Filter(key string, value interface{}) driver.Session {
	return d.NewSession().Filter(key, value)
}

func (d *defaultClient) ID(id interface{}) driver.Session {
	return d.NewSession().ID(id)
}

func (d *defaultClient) Gt(key string, value interface{}) driver.Session {
	return d.NewSession().Gt(key, value)
}

func (d *defaultClient) Gte(key string, value interface{}) driver.Session {
	return d.NewSession().Gte(key, value)
}

func (d *defaultClient) Lt(key string, value interface{}) driver.Session {
	return d.NewSession().Lt(key, value)
}

func (d *defaultClient) Lte(key string, value interface{}) driver.Session {
	return d.NewSession().Lte(key, value)
}

func (d *defaultClient) In(key string, value interface{}) driver.Session {
	return d.NewSession().In(key, value)
}

func (d *defaultClient) And(filter driver.Condition) driver.Session {
	return d.NewSession().And(filter)
}

func (d *defaultClient) Not(key string, value interface{}) driver.Session {
	return d.NewSession().Not(key, value)
}

func (d *defaultClient) Or(filter driver.Condition) driver.Session {
	return d.NewSession().Or(filter)
}

func (d *defaultClient) InsertOne(v interface{}, ctx ...context.Context) (primitive.ObjectID, error) {
	return d.NewSession().InsertOne(v, ctx...)
}

func (d *defaultClient) InsertMany(v interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error) {
	return d.NewSession().InsertMany(v, ctx...)
}

func (d *defaultClient) Limit(limit int64) driver.Session {
	return d.NewSession().Limit(limit)
}

func (d *defaultClient) Skip(skip int64) driver.Session {
	return d.NewSession().Limit(skip)
}

func (d *defaultClient) Count(i interface{}, ctx ...context.Context) (int64, error) {
	return d.NewSession().Count(i, ctx...)
}

func (d *defaultClient) Desc(s2 ...string) driver.Session {
	return d.NewSession().Desc(s2...)
}

func (d *defaultClient) Update(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateOne(bean, ctx...)
}

// UpdateMany The following operation updates all of the documents with quantity value less than 50.
func (d *defaultClient) UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateMany(bean, ctx...)
}

func (d *defaultClient) DeleteOne(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
	return d.NewSession().DeleteOne(filter, ctx...)
}

func (d *defaultClient) DeleteMany(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
	return d.NewSession().DeleteMany(filter, ctx...)
}
func (d *defaultClient) SoftDeleteOne(filter interface{}, ctx ...context.Context) error {
	return d.NewSession().SoftDeleteOne(filter, ctx...)
}

func (d *defaultClient) SoftDeleteMany(filter interface{}, ctx ...context.Context) error {
	return d.NewSession().SoftDeleteMany(filter, ctx...)
}

func (d *defaultClient) FilterBy(object interface{}) driver.Session {
	return d.NewSession().FilterBy(object)
}

func (d *defaultClient) DropAll(doc interface{}, ctx ...context.Context) error {
	return d.NewIndexes().DropAll(doc, ctx...)
}

func (d *defaultClient) DropOne(doc interface{}, name string, ctx ...context.Context) error {
	return d.NewIndexes().DropOne(doc, name, ctx...)
}

func (d *defaultClient) AddIndex(keys interface{}, opt ...*options.IndexOptions) driver.Indexes {
	return d.NewIndexes().AddIndex(keys, opt...)
}

func (d *defaultClient) NewIndexes() driver.Indexes {
	return NewIndexes(d)
}

func (d *defaultClient) NewSession() driver.Session {
	return NewSession(d)
}

func (d *defaultClient) Aggregate() driver.Aggregate {
	return NewAggregate(d)
}

func (d *defaultClient) CollectionNameForStruct(doc interface{}) (*schemas.Collection, error) {
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

//func (d *defaultClient) SetDatabase(string string) driver.Client {
//	d.db = string
//	return d
//}

func (d *defaultClient) CollectionNameForSlice(doc interface{}) (*schemas.Collection, error) {
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

func (d *defaultClient) getStructCollAndSetKey(doc interface{}) (*schemas.Collection, error) {
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
