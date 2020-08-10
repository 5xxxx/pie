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

func NewTugrik(opts ...*options.ClientOptions) (*Tugrik, error) {
	mapper := names.NewCacheMapper(new(names.SnakeMapper))

	parser := NewParser(mapper, mapper)
	tugrik := &Tugrik{
		clientOpts: opts,
		parser:     parser,
	}
	return tugrik, nil
}

func (e *Tugrik) Connect(ctx context.Context) (err error) {
	e.client, err = mongo.Connect(ctx, e.clientOpts...)
	return err
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

func (e *Tugrik) Count(i ...interface{}) (int64, error) {
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
