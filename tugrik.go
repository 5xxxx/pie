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
	t, err := tugrik.NewTugrik()
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
	err = t.Filter("nickName", "淳朴的润土").FindOne(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
*/

package tugrik

import (
	"context"
	"tugrik/names"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Tugrik struct {
	client     *mongo.Client
	parser     *Parser
	db         string
	clientOpts []*options.ClientOptions
}

func (e *Tugrik) Distinct(doc interface{}, columns string) ([]interface{}, error) {
	session := e.NewSession()
	return session.Distinct(doc, columns)
}

func (e *Tugrik) FindOne(doc interface{}) error {
	session := e.NewSession()
	return session.FindOne(doc)
}

func (e *Tugrik) FindAll(docs interface{}) error {
	session := e.NewSession()
	return session.FindAll(docs)
}

func (e *Tugrik) RegexFilter(key, pattern string) *Session {
	session := e.NewSession()
	return session.RegexFilter(key, pattern)
}

func (e *Tugrik) Asc(colNames ...string) *Session {
	session := e.NewSession()
	return session.Asc(colNames...)
}

func (e *Tugrik) Eq(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Eq(key, value)
}

func (e *Tugrik) Ne(key string, ne interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, ne)
}

func (e *Tugrik) Nin(key string, nin interface{}) *Session {
	session := e.NewSession()
	return session.Nin(key, nin)
}

func (e *Tugrik) Nor(filter Session) *Session {
	session := e.NewSession()
	return session.Nor(filter)
}

func (e *Tugrik) Exists(key string, exists bool, filter ...*Session) *Session {
	session := e.NewSession()
	return session.Exists(key, exists, filter...)
}

func (e *Tugrik) Type(key string, t interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, t)
}

func (e *Tugrik) Expr(filter Session) *Session {
	session := e.NewSession()
	return session.Expr(filter)
}

func (e *Tugrik) Regex(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Regex(key, value)
}

func NewTugrik(opts ...*options.ClientOptions) (*Tugrik, error) {
	mapper := names.NewCacheMapper(new(names.SnakeMapper))
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	parser := NewParser(mapper, mapper)
	tugrik := &Tugrik{
		clientOpts: opts,
		parser:     parser,
		client:     client,
	}
	return tugrik, nil
}

func (e *Tugrik) Connect(ctx context.Context) (err error) {
	e.client, err = mongo.Connect(ctx, e.clientOpts...)
	return err
}

func (e *Tugrik) Disconnect(ctx context.Context) error {
	return e.client.Disconnect(ctx)
}

func (e *Tugrik) SetDatabase(db string) {
	e.db = db
}

func (e *Tugrik) Collection(name string) *mongo.Collection {
	return e.client.Database(e.db).Collection(name)
}

func (e *Tugrik) Ping() error {
	return e.client.Ping(context.TODO(), readpref.Primary())
}

func (e *Tugrik) Filter(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Filter(key, value)
}

func (e *Tugrik) ID(id interface{}) *Session {
	session := e.NewSession()
	return session.ID(id)
}

func (e *Tugrik) Gt(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, value)
}

func (e *Tugrik) Gte(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Gte(key, value)
}

func (e *Tugrik) Lt(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Lt(key, value)
}

func (e *Tugrik) Lte(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Lte(key, value)
}

func (e *Tugrik) In(key string, value interface{}) *Session {
	session := e.NewSession()
	session.In(key, value)
	return session
}

func (e *Tugrik) And(filter Session) *Session {
	session := e.NewSession()
	session.And(filter)
	return session
}

func (e *Tugrik) Not(key string, value interface{}) *Session {
	session := e.NewSession()
	session.Not(key, value)
	return session
}

func (e *Tugrik) Or(filter Session) *Session {
	session := e.NewSession()
	session.Or(filter)
	return session
}

func (e *Tugrik) InsertOne(v interface{}) error {
	session := e.NewSession()
	return session.InsertOne(v)
}

func (e *Tugrik) InsertMany(v []interface{}) error {
	session := e.NewSession()
	return session.InsertMany(v)
}

func (e *Tugrik) Limit(limit int64) *Session {
	session := e.NewSession()
	return session.Limit(limit)
}

func (e *Tugrik) Skip(skip int64) *Session {
	session := e.NewSession()
	return session.Limit(skip)
}

func (e *Tugrik) Count(i interface{}) (int64, error) {
	session := e.NewSession()
	return session.Count(i)
}

func (e *Tugrik) Desc(s2 ...string) *Session {
	session := e.NewSession()
	return session.Desc(s2...)
}

func (e *Tugrik) Update(bean interface{}) error {
	session := e.NewSession()
	return session.Update(bean)
}

//The following operation updates all of the documents with quantity value less than 50.
func (e *Tugrik) UpdateMany(bean interface{}) error {
	session := e.NewSession()
	return session.UpdateMany(bean)
}

func (e *Tugrik) DeleteOne(filter interface{}) error {
	session := e.NewSession()
	return session.DeleteOne(filter)
}

func (e *Tugrik) DeleteMany(filter interface{}) error {
	session := e.NewSession()
	return session.DeleteMany(filter)
}

func (e *Tugrik) NewSession() *Session {
	return &Session{engine: e, m: bson.M{}}
}
