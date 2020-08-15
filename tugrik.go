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
	err = t.filter("nickName", "淳朴的润土").FindOne(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
*/

package tugrik

import (
	"context"
	"errors"
	"github.com/NSObjects/tugrik/names"
	"github.com/NSObjects/tugrik/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"

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

func (e *Tugrik) Distinct(ctx context.Context, doc interface{}, columns string) ([]interface{}, error) {
	session := e.NewSession()
	return session.Distinct(ctx, doc, columns)
}

func (e *Tugrik) FindOne(ctx context.Context, doc interface{}) error {
	session := e.NewSession()
	return session.FindOne(ctx, doc)
}

func (e *Tugrik) FindAll(ctx context.Context, docs interface{}) error {
	session := e.NewSession()
	return session.FindAll(ctx, docs)
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

func (e *Tugrik) Nor(c Condition) *Session {
	session := e.NewSession()
	return session.Nor(c)
}

func (e *Tugrik) Exists(key string, exists bool, filter ...Condition) *Session {
	session := e.NewSession()
	return session.Exists(key, exists, filter...)
}

func (e *Tugrik) Type(key string, t interface{}) *Session {
	session := e.NewSession()
	return session.Gt(key, t)
}

func (e *Tugrik) Expr(filter Condition) *Session {
	session := e.NewSession()
	return session.Expr(filter)
}

func (e *Tugrik) Regex(key string, value interface{}) *Session {
	session := e.NewSession()
	return session.Regex(key, value)
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

func (e *Tugrik) And(filter Condition) *Session {
	session := e.NewSession()
	session.And(filter)
	return session
}

func (e *Tugrik) Not(key string, value interface{}) *Session {
	session := e.NewSession()
	session.Not(key, value)
	return session
}

func (e *Tugrik) Or(filter Condition) *Session {
	session := e.NewSession()
	session.Or(filter)
	return session
}

func (e *Tugrik) InsertOne(ctx context.Context, v interface{}) (primitive.ObjectID, error) {
	session := e.NewSession()
	return session.InsertOne(ctx, v)
}

func (e *Tugrik) InsertMany(ctx context.Context, v []interface{}) (*mongo.InsertManyResult, error) {
	session := e.NewSession()
	return session.InsertMany(ctx, v)
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

func (e *Tugrik) Update(ctx context.Context, bean interface{}) (*mongo.UpdateResult, error) {
	session := e.NewSession()
	return session.Update(ctx, bean)
}

//The following operation updates all of the documents with quantity value less than 50.
func (e *Tugrik) UpdateMany(bean interface{}) error {
	session := e.NewSession()
	return session.UpdateMany(bean)
}

func (e *Tugrik) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := e.NewSession()
	return session.DeleteOne(ctx, filter)
}

func (e *Tugrik) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	session := e.NewSession()
	return session.DeleteMany(ctx, filter)
}

func (e *Tugrik) FilterBy(object interface{}) *Session {
	session := e.NewSession()
	return session.FilterBy(object)
}

func (e *Tugrik) NewSession() *Session {
	return NewSession(e)
}

func (e *Tugrik) Aggregate() *Aggregate {
	return NewAggregate(e)
}

func (s *Tugrik) getStructColl(doc interface{}) (*mongo.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := s.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}
	coll := s.Collection(t.Name)
	return coll, nil
}

func (s *Tugrik) getSliceColl(doc interface{}) (*mongo.Collection, error) {
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
			t, err = s.parser.Parse(pv)
		}
	} else {
		t, err = s.parser.Parse(sliceValue)
	}

	if err != nil {
		return nil, err
	}

	coll := s.Collection(t.Name)
	return coll, nil
}

func (s *Tugrik) getStructCollAndSetKey(doc interface{}) (*mongo.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := s.parser.Parse(beanValue)
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

	coll := s.Collection(t.Name)
	return coll, nil
}
