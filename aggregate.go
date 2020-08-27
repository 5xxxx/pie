/*
 *
 * aggregate.go
 * tugrik
 *
 * Created by lintao on 2020/8/13 12:26 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IAggregate interface {
	AddFields() *Aggregate
	Bucket() *Aggregate
	BucketAuto() *Aggregate
	CollStats() *Aggregate
	Count() *Aggregate
	CurrentOp() *Aggregate
	Facet() *Aggregate
	GeoNear() *Aggregate
	GraphLookup() *Aggregate
	Group() *Aggregate
	IndexStats() *Aggregate
	Limit() *Aggregate
	ListLocalSession() *Aggregate
	ListSession() *Aggregate
	Lookup() *Aggregate
	Match(filter Session) *Aggregate
	Merge() *Aggregate
	Out() *Aggregate
	PlanCacheStats() *Aggregate
	Project() *Aggregate
	Redact() *Aggregate
	ReplaceRoot() *Aggregate
	ReplaceWith() *Aggregate
	Sample() *Aggregate
	Set() *Aggregate
	Skip() *Aggregate
	Sort() *Aggregate
	SortByCount() *Aggregate
	UnionWith() *Aggregate
	Unset() *Aggregate
	Unwind() *Aggregate
	All(result interface{}) error
}

type Aggregate struct {
	db       string
	doc      interface{}
	engine   *Driver
	pipeline bson.A
	opts     []*options.AggregateOptions
}

func NewAggregate(engine *Driver) *Aggregate {
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
		if err := aggregate.Decode(&result); err != nil {
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
		coll, err = a.collectionForSlice(a.doc)
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

func (s *Aggregate) SetCollection(doc interface{}) *Aggregate {
	s.doc = doc
	return s
}

func (s *Aggregate) SetDatabase(db string) *Aggregate {
	s.db = db
	return s
}

func (e *Aggregate) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
	coll, err := e.engine.CollectionNameForStruct(doc)
	if err != nil {
		return nil, err
	}
	return e.collection(coll.Name), nil
}

func (e *Aggregate) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
	coll, err := e.engine.CollectionNameForSlice(doc)
	if err != nil {
		return nil, err
	}
	return e.collection(coll.Name), nil
}

func (s Aggregate) collection(name string) *mongo.Collection {
	var db string
	if s.db != "" {
		db = s.db
	} else {
		db = s.engine.db
	}
	return s.engine.client.Database(db).Collection(name)
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (ao *Aggregate) SetAllowDiskUse(b bool) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetAllowDiskUse(b))
	return ao
}

// SetBatchSize sets the value for the BatchSize field.
func (ao *Aggregate) SetBatchSize(i int32) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetBatchSize(i))
	return ao
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ao *Aggregate) SetBypassDocumentValidation(b bool) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetBypassDocumentValidation(b))
	return ao
}

// SetCollation sets the value for the Collation field.
func (ao *Aggregate) SetCollation(c *options.Collation) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetCollation(c))
	return ao
}

// SetMaxTime sets the value for the MaxTime field.
func (ao *Aggregate) SetMaxTime(d time.Duration) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetMaxTime(d))
	return ao
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (ao *Aggregate) SetMaxAwaitTime(d time.Duration) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetMaxAwaitTime(d))
	return ao
}

// SetComment sets the value for the Comment field.
func (ao *Aggregate) SetComment(s string) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetComment(s))
	return ao
}

// SetHint sets the value for the Hint field.
func (ao *Aggregate) SetHint(h interface{}) *Aggregate {
	ao.opts = append(ao.opts, options.Aggregate().SetHint(h))
	return ao
}

func (a *Aggregate) Pipeline(pipeline bson.A) *Aggregate {
	a.pipeline = pipeline
	return a
}

func (a *Aggregate) Match(c Condition) *Aggregate {
	a.pipeline = append(a.pipeline, bson.M{
		"$match": c.Filters(),
	})
	return a
}

//func (a *Aggregate) Group() *Aggregate {
//	panic("")
//}
