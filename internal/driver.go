package internal

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/NSObjects/pie/driver"
	"github.com/NSObjects/pie/names"
	"github.com/NSObjects/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type defaultDriver struct {
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
	driver := &defaultDriver{
		clientOpts: opts,
		parser:     parser,
		client:     client,
		db:         db,
	}
	return driver, nil
}

func (d *defaultDriver) Connect(ctx context.Context) (err error) {
	d.client, err = mongo.Connect(ctx, d.clientOpts...)
	return err
}

func (d *defaultDriver) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

func (d *defaultDriver) BulkWrite(ctx context.Context, docs interface{}) (*mongo.BulkWriteResult, error) {
	session := d.NewSession()
	return session.BulkWrite(ctx, docs)
}

func (d *defaultDriver) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	session := d.NewSession()
	return session.Distinct(ctx, doc, columns)
}

func (d *defaultDriver) ReplaceOne(ctx context.Context, doc interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.ReplaceOne(ctx, doc)
}

func (d defaultDriver) UpdateOneBson(ctx context.Context, coll interface{}, bson interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.UpdateOneBson(ctx, coll, bson)
}

func (d defaultDriver) FindOneAndUpdateBson(ctx context.Context, coll interface{}, bson interface{}) (*mongo.SingleResult, error) {
	session := d.NewSession()
	return session.FindOneAndUpdateBson(ctx, coll, bson)
}

func (d defaultDriver) UpdateManyBson(ctx context.Context, coll interface{}, bson interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.UpdateManyBson(ctx, coll, bson)
}

func (d *defaultDriver) FindOneAndReplace(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOneAndReplace(ctx, doc)
}

func (d *defaultDriver) FindOneAndUpdate(ctx context.Context, doc interface{}) (*mongo.SingleResult, error) {
	session := d.NewSession()
	return session.FindOneAndUpdate(ctx, doc)
}

func (d *defaultDriver) FindAndDelete(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindAndDelete(ctx, doc)
}

// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
func (d *defaultDriver) FindOne(ctx context.Context, doc interface{}) error {
	session := d.NewSession()
	return session.FindOne(ctx, doc)
}

func (d *defaultDriver) FindAll(ctx context.Context, docs interface{}) error {
	session := d.NewSession()
	return session.FindAll(ctx, docs)
}

func (d *defaultDriver) RegexFilter(key, pattern string) driver.Session {
	session := d.NewSession()
	return session.RegexFilter(key, pattern)
}

func (d *defaultDriver) Asc(colNames ...string) driver.Session {
	session := d.NewSession()
	return session.Asc(colNames...)
}

func (d *defaultDriver) Eq(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Eq(key, value)
}

func (d *defaultDriver) Ne(key string, ne interface{}) driver.Session {
	session := d.NewSession()
	return session.Gt(key, ne)
}

func (d *defaultDriver) Nin(key string, nin interface{}) driver.Session {
	session := d.NewSession()
	return session.Nin(key, nin)
}

func (d *defaultDriver) Nor(c driver.Condition) driver.Session {
	session := d.NewSession()
	return session.Nor(c)
}

func (d *defaultDriver) Exists(key string, exists bool, filter ...driver.Condition) driver.Session {
	session := d.NewSession()
	return session.Exists(key, exists, filter...)
}

func (d *defaultDriver) Type(key string, t interface{}) driver.Session {
	session := d.NewSession()
	return session.Gt(key, t)
}

func (d *defaultDriver) Expr(filter driver.Condition) driver.Session {
	session := d.NewSession()
	return session.Expr(filter)
}

func (d *defaultDriver) Regex(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Regex(key, value)
}

func (d *defaultDriver) DataBase() *mongo.Database {
	return d.client.Database(d.db)
}

func (d *defaultDriver) Collection(name string) *mongo.Collection {
	return d.client.Database(d.db).Collection(name)
}

func (d *defaultDriver) Ping() error {
	return d.client.Ping(context.TODO(), readpref.Primary())
}

func (d *defaultDriver) Filter(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Filter(key, value)
}

func (d *defaultDriver) ID(id interface{}) driver.Session {
	session := d.NewSession()
	return session.ID(id)
}

func (d *defaultDriver) Gt(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Gt(key, value)
}

func (d *defaultDriver) Gte(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Gte(key, value)
}

func (d *defaultDriver) Lt(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Lt(key, value)
}

func (d *defaultDriver) Lte(key string, value interface{}) driver.Session {
	session := d.NewSession()
	return session.Lte(key, value)
}

func (d *defaultDriver) In(key string, value interface{}) driver.Session {
	session := d.NewSession()
	session.In(key, value)
	return session
}

func (d *defaultDriver) And(filter driver.Condition) driver.Session {
	session := d.NewSession()
	session.And(filter)
	return session
}

func (d *defaultDriver) Not(key string, value interface{}) driver.Session {
	session := d.NewSession()
	session.Not(key, value)
	return session
}

func (d *defaultDriver) Or(filter driver.Condition) driver.Session {
	session := d.NewSession()
	session.Or(filter)
	return session
}

func (d *defaultDriver) InsertOne(ctx context.Context, v interface{}) (primitive.ObjectID, error) {
	session := d.NewSession()
	return session.InsertOne(ctx, v)
}

func (d *defaultDriver) InsertMany(ctx context.Context, v interface{}) (*mongo.InsertManyResult, error) {
	session := d.NewSession()
	return session.InsertMany(ctx, v)
}

func (d *defaultDriver) Limit(limit int64) driver.Session {
	session := d.NewSession()
	return session.Limit(limit)
}

func (d *defaultDriver) Skip(skip int64) driver.Session {
	session := d.NewSession()
	return session.Limit(skip)
}

func (d *defaultDriver) Count(i interface{}) (int64, error) {
	session := d.NewSession()
	return session.Count(i)
}

func (d *defaultDriver) Desc(s2 ...string) driver.Session {
	session := d.NewSession()
	return session.Desc(s2...)
}

func (d *defaultDriver) Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.UpdateOne(ctx, bean)
}

//The following operation updates all of the documents with quantity value less than 50.
func (d *defaultDriver) UpdateMany(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := d.NewSession()
	return session.UpdateMany(ctx, bean)
}

func (d *defaultDriver) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := d.NewSession()
	return session.DeleteOne(ctx, filter)
}

func (d *defaultDriver) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := d.NewSession()
	return session.DeleteMany(ctx, filter)
}
func (d *defaultDriver) SoftDeleteOne(ctx context.Context, filter interface{}) error {
	session := d.NewSession()
	return session.SoftDeleteOne(ctx, filter)
}

func (d *defaultDriver) SoftDeleteMany(ctx context.Context, filter interface{}) error {
	session := d.NewSession()
	return session.SoftDeleteMany(ctx, filter)
}

func (d *defaultDriver) FilterBy(object interface{}) driver.Session {
	session := d.NewSession()
	return session.FilterBy(object)
}

func (d *defaultDriver) DropAll(ctx context.Context, doc interface{}) error {
	indexes := d.NewIndexes()
	return indexes.DropAll(ctx, doc)
}

func (d *defaultDriver) DropOne(ctx context.Context, doc interface{}, name string) error {
	indexes := d.NewIndexes()
	return indexes.DropOne(ctx, doc, name)
}

func (d *defaultDriver) AddIndex(keys interface{}, opt ...*options.IndexOptions) driver.Indexes {
	indexes := d.NewIndexes()
	return indexes.AddIndex(keys, opt...)
}

func (d *defaultDriver) NewIndexes() driver.Indexes {
	return NewIndexes(d)
}

func (d *defaultDriver) NewSession() driver.Session {
	return NewSession(d)
}

func (d *defaultDriver) Aggregate() driver.Aggregate {
	return NewAggregate(d)
}

func (d *defaultDriver) CollectionNameForStruct(doc interface{}) (*schemas.Collection, error) {
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

func (d *defaultDriver) SetDatabase(string string) driver.Client {
	d.db = string
	return d
}

func (d *defaultDriver) CollectionNameForSlice(doc interface{}) (*schemas.Collection, error) {
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

func (d *defaultDriver) getStructCollAndSetKey(doc interface{}) (*schemas.Collection, error) {
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
