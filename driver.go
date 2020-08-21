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
	"github.com/NSObjects/pie/names"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	tugrik := &Driver{
		clientOpts: opts,
		parser:     parser,
		client:     client,
	}
	return tugrik, nil
}

func (e *Driver) Connect(ctx context.Context) (err error) {
	e.client, err = mongo.Connect(ctx, e.clientOpts...)
	return err
}

func (e *Driver) Disconnect(ctx context.Context) error {
	return e.client.Disconnect(ctx)
}

func (d *Driver) BulkWrite(ctx context.Context, docs interface{}) (*mongo.BulkWriteResult, error) {
	session := d.NewSession()
	return session.BulkWrite(ctx, docs)
}

func (e *Driver) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	session := e.NewSession()
	return session.Distinct(ctx, doc, columns)
}

func (e *Driver) ReplaceOne(ctx context.Context, doc interface{}) (*mongo.UpdateResult, error) {
	session := e.NewSession()
	return session.ReplaceOne(ctx, doc)
}

func (s *Driver) FindOneAndReplace(ctx context.Context, doc interface{}) error {
	session := s.NewSession()
	return session.FindOneAndReplace(ctx, doc)
}

func (s *Driver) FindOneAndUpdate(ctx context.Context, doc interface{}) error {
	session := s.NewSession()
	return session.FindOneAndUpdate(ctx, doc)
}

func (s *Driver) FindAndDelete(ctx context.Context, doc interface{}) error {
	session := s.NewSession()
	return session.FindAndDelete(ctx, doc)
}

// FindOne executes a find command and returns a SingleResult for one document in the collection.
func (e *Driver) FindOne(ctx context.Context, doc interface{}) error {
	session := e.NewSession()
	return session.FindOne(ctx, doc)
}

func (e *Driver) FindAll(ctx context.Context, docs interface{}) error {
	session := e.NewSession()
	return session.FindAll(ctx, docs)
}

func (e *Driver) RegexFilter(key, pattern string) *Session {
	session := e.NewSession()
	return session.RegexFilter(key, pattern)
}

func (e *Driver) Asc(colNames ...string) *Session {
	session := e.NewSession()
	return session.Asc(colNames...)
}

func (e *Driver) Eq(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Eq(key, value)
}

func (e *Driver) Ne(key string, ne interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, ne)
}

func (e *Driver) Nin(key string, nin interface{}) *Session {
	session := e.NewSession()
	return session.Nin(key, nin)
}

func (e *Driver) Nor(c Condition) *Session {
	session := e.NewSession()
	return session.Nor(c)
}

func (e *Driver) Exists(key string, exists bool, filter ...Condition) *Session {
	session := e.NewSession()
	return session.Exists(key, exists, filter...)
}

func (e *Driver) Type(key string, t interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, t)
}

func (e *Driver) Expr(filter Condition) *Session {
	session := e.NewSession()
	return session.Expr(filter)
}

func (e *Driver) Regex(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Regex(key, value)
}

func (e *Driver) SetDatabase(db string) {
	e.db = db
}

func (e *Driver) DataBase() *mongo.Database {
	return e.client.Database(e.db)
}

func (e *Driver) Collection(name string) *mongo.Collection {
	return e.client.Database(e.db).Collection(name)
}

func (e *Driver) Ping() error {
	return e.client.Ping(context.TODO(), readpref.Primary())
}

func (e *Driver) Filter(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Filter(key, value)
}

func (e *Driver) ID(id interface{}) *Session {
	session := e.NewSession()
	return session.ID(id)
}

func (e *Driver) Gt(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, value)
}

func (e *Driver) Gte(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Gte(key, value)
}

func (e *Driver) Lt(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Lt(key, value)
}

func (e *Driver) Lte(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Lte(key, value)
}

func (e *Driver) In(key string, value interface{}) *Session {
	session := e.NewSession()
	session.In(key, value)
	return session
}

func (e *Driver) And(filter Condition) *Session {
	session := e.NewSession()
	session.And(filter)
	return session
}

func (e *Driver) Not(key string, value interface{}) *Session {
	session := e.NewSession()
	session.Not(key, value)
	return session
}

func (e *Driver) Or(filter Condition) *Session {
	session := e.NewSession()
	session.Or(filter)
	return session
}

func (e *Driver) InsertOne(ctx context.Context, v interface{}) (primitive.ObjectID, error) {
	session := e.NewSession()
	return session.InsertOne(ctx, v)
}

func (e *Driver) InsertMany(ctx context.Context, v interface{}) (*mongo.InsertManyResult, error) {
	session := e.NewSession()
	return session.InsertMany(ctx, v)
}

func (e *Driver) Limit(limit int64) *Session {
	session := e.NewSession()
	return session.Limit(limit)
}

func (e *Driver) Skip(skip int64) *Session {
	session := e.NewSession()
	return session.Limit(skip)
}

func (e *Driver) Count(i interface{}) (int64, error) {
	session := e.NewSession()
	return session.Count(i)
}

func (e *Driver) Desc(s2 ...string) *Session {
	session := e.NewSession()
	return session.Desc(s2...)
}

func (e *Driver) Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := e.NewSession()
	return session.Update(ctx, bean)
}

//The following operation updates all of the documents with quantity value less than 50.
func (e *Driver) UpdateMany(bean interface{}) error {
	session := e.NewSession()
	return session.UpdateMany(bean)
}

func (e *Driver) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := e.NewSession()
	return session.DeleteOne(ctx, filter)
}

func (e *Driver) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := e.NewSession()
	return session.DeleteMany(ctx, filter)
}

func (e *Driver) FilterBy(object interface{}) *Session {
	session := e.NewSession()
	return session.FilterBy(object)
}

func (s *Driver) DropAll(ctx context.Context, doc interface{}) error {
	indexes := s.NewIndexes()
	return indexes.DropAll(ctx, doc)
}

func (s *Driver) DropOne(ctx context.Context, doc interface{}, name string) error {
	indexes := s.NewIndexes()
	return indexes.DropOne(ctx, doc, name)
}

func (e *Driver) AddIndex(keys interface{}, opt ...*options.IndexOptions) *Indexes {
	indexes := e.NewIndexes()
	return indexes.AddIndex(keys, opt...)
}

func (e *Driver) NewIndexes() *Indexes {
	return NewIndexes(e)
}

func (e *Driver) NewSession() *Session {
	return NewSession(e)
}

func (e *Driver) Aggregate() *Aggregate {
	return NewAggregate(e)
}
