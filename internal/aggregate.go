package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"

	"github.com/5xxxx/pie/driver"

	"github.com/5xxxx/pie/schemas"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type aggregate struct {
	db       string
	doc      interface{}
	engine   driver.Client
	pipeline bson.A
	opts     []*options.AggregateOptions
	collOpts []*options.CollectionOptions
}

func NewAggregate(engine driver.Client) driver.Aggregate {
	return &aggregate{engine: engine}
}

func (a *aggregate) One(result interface{}, ctx ...context.Context) error {

	var c context.Context
	if len(ctx) > 0 {
		c = ctx[0]
	} else {
		c = context.Background()
	}

	var coll *mongo.Collection
	var err error
	if a.doc != nil {
		coll, err = a.collectionForStruct(a.doc)
	} else {
		coll, err = a.collectionForStruct(result)
	}

	if err != nil {
		return err
	}

	aggregate, err := coll.Aggregate(c, a.pipeline, a.opts...)
	if err != nil {
		return err
	}
	if next := aggregate.Next(c); next {
		if err := aggregate.Decode(result); err != nil {
			return err
		}
	} else {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (a *aggregate) All(result interface{}, ctx ...context.Context) error {
	var c context.Context
	if len(ctx) > 0 {
		c = ctx[0]
	} else {
		c = context.Background()
	}

	var coll *mongo.Collection
	var err error
	if a.doc != nil {
		coll, err = a.collectionForStruct(a.doc)
	} else {
		coll, err = a.collectionForSlice(result)
	}
	if err != nil {
		return err
	}

	aggregate, err := coll.Aggregate(c, a.pipeline, a.opts...)
	if err != nil {
		return err
	}

	return aggregate.All(c, result)
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (a *aggregate) SetAllowDiskUse(b bool) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetAllowDiskUse(b))
	return a
}

// SetBatchSize sets the value for the BatchSize field.
func (a *aggregate) SetBatchSize(i int32) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetBatchSize(i))
	return a
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *aggregate) SetBypassDocumentValidation(b bool) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetBypassDocumentValidation(b))
	return a
}

// SetCollation sets the value for the Collation field.
func (a *aggregate) SetCollation(c *options.Collation) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetCollation(c))
	return a
}

// SetMaxTime sets the value for the MaxTime field.
func (a *aggregate) SetMaxTime(d time.Duration) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetMaxTime(d))
	return a
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *aggregate) SetMaxAwaitTime(d time.Duration) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetMaxAwaitTime(d))
	return a
}

// SetComment sets the value for the Comment field.
func (a *aggregate) SetComment(s string) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetComment(s))
	return a
}

// SetHint sets the value for the Hint field.
func (a *aggregate) SetHint(h interface{}) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetHint(h))
	return a
}

func (a *aggregate) Pipeline(pipeline bson.A) driver.Aggregate {
	a.pipeline = append(a.pipeline, pipeline...)
	return a
}

func (a *aggregate) Match(c driver.Condition) driver.Aggregate {
	filters, err := c.Filters()
	if err != nil {
		panic(err)
	}
	a.pipeline = append(a.pipeline, bson.M{
		"$match": filters,
	})
	return a
}

func (a *aggregate) SetDatabase(db string) driver.Aggregate {
	a.db = db
	return a
}

func (a *aggregate) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if a.doc != nil {
		coll, err = a.engine.CollectionNameForStruct(a.doc)
	} else {
		coll, err = a.engine.CollectionNameForStruct(doc)
	}
	if err != nil {
		return nil, err
	}
	return a.collectionByName(coll.Name), nil
}

func (a *aggregate) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if a.doc != nil {
		coll, err = a.engine.CollectionNameForStruct(a.doc)
	} else {
		coll, err = a.engine.CollectionNameForSlice(doc)
	}
	if err != nil {
		return nil, err
	}
	return a.collectionByName(coll.Name), nil
}

func (a *aggregate) collectionByName(name string) *mongo.Collection {
	if a.collOpts == nil {
		a.collOpts = make([]*options.CollectionOptions, 0)
	}
	return a.engine.Collection(name, a.collOpts, a.db)
}

// SetReadConcern sets the value for the ReadConcern field.
func (a *aggregate) SetReadConcern(rc *readconcern.ReadConcern) driver.Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetReadConcern(rc))
	return a
}

// SetCollWriteConcern sets the value for the WriteConcern field.
func (a *aggregate) SetCollWriteConcern(wc *writeconcern.WriteConcern) driver.Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetWriteConcern(wc))
	return a
}

// SetCollReadPreference sets the value for the ReadPreference field.
func (a *aggregate) SetCollReadPreference(rp *readpref.ReadPref) driver.Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetReadPreference(rp))
	return a
}

// SetCollRegistry sets the value for the Registry field.
func (a *aggregate) SetCollRegistry(r *bsoncodec.Registry) driver.Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetRegistry(r))
	return a
}

func (a *aggregate) Collection(doc interface{}) driver.Aggregate {
	a.doc = doc
	return a
}
