/*
 *
 * aggregate.go
 * pie
 *
 * Created by lintao on 2020/8/13 12:26 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package internal

import (
	"context"
	"time"

	"github.com/NSObjects/pie/driver"
	_interface "github.com/NSObjects/pie/driver"

	"github.com/NSObjects/pie/schemas"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Aggregate struct {
	db       string
	doc      interface{}
	engine   driver.Client
	pipeline bson.A
	opts     []*options.AggregateOptions
}

func NewAggregate(engine driver.Client) driver.Aggregate {
	return &Aggregate{engine: engine}
}

func (a *Aggregate) One(ctx context.Context, result interface{}) error {
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

	aggregate, err := coll.Aggregate(ctx, a.pipeline, a.opts...)
	if err != nil {
		return err
	}
	if next := aggregate.Next(ctx); next {
		if err := aggregate.Decode(result); err != nil {
			return err
		}
	} else {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (a *Aggregate) All(ctx context.Context, result interface{}) error {
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

	aggregate, err := coll.Aggregate(ctx, a.pipeline, a.opts...)
	if err != nil {
		return err
	}

	return aggregate.All(ctx, result)
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (a *Aggregate) SetAllowDiskUse(b bool) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetAllowDiskUse(b))
	return a
}

// SetBatchSize sets the value for the BatchSize field.
func (a *Aggregate) SetBatchSize(i int32) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetBatchSize(i))
	return a
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *Aggregate) SetBypassDocumentValidation(b bool) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetBypassDocumentValidation(b))
	return a
}

// SetCollation sets the value for the Collation field.
func (a *Aggregate) SetCollation(c *options.Collation) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetCollation(c))
	return a
}

// SetMaxTime sets the value for the MaxTime field.
func (a *Aggregate) SetMaxTime(d time.Duration) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetMaxTime(d))
	return a
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *Aggregate) SetMaxAwaitTime(d time.Duration) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetMaxAwaitTime(d))
	return a
}

// SetComment sets the value for the Comment field.
func (a *Aggregate) SetComment(s string) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetComment(s))
	return a
}

// SetHint sets the value for the Hint field.
func (a *Aggregate) SetHint(h interface{}) driver.Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetHint(h))
	return a
}

func (a *Aggregate) Pipeline(pipeline bson.A) driver.Aggregate {
	a.pipeline = pipeline
	return a
}

func (a *Aggregate) Match(c _interface.Condition) driver.Aggregate {
	a.pipeline = append(a.pipeline, bson.M{
		"$match": c.Filters(),
	})
	return a
}

func (a *Aggregate) SetDatabase(db string) driver.Aggregate {
	a.db = db
	return a
}

func (a *Aggregate) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
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

func (a *Aggregate) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
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

func (a *Aggregate) collectionByName(name string) *mongo.Collection {
	var db string
	if a.db != "" {
		db = a.db
	} else {
		db = a.engine.DataBase().Name()
	}
	return a.engine.SetDatabase(db).Collection(name)
}

func (a *Aggregate) Collection(doc interface{}) driver.Aggregate {
	a.doc = doc
	return a
}
