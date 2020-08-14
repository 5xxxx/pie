/*
 *
 * aggregate.go
 * tugrik
 *
 * Created by lintao on 2020/8/13 12:26 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package tugrik

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
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
	engine   *Tugrik
	ctx      context.Context
	pipeline bson.A
	opts     []*options.AggregateOptions
}

func NewAggregate(engine *Tugrik) *Aggregate {
	return &Aggregate{engine: engine}
}

func (a *Aggregate) All(result interface{}) error {
	coll, err := a.engine.getSliceColl(result)
	if err != nil {
		return err
	}

	aggregate, err := coll.Aggregate(a.ctx, a.pipeline, a.opts...)
	if err != nil {
		return err
	}

	return aggregate.All(a.ctx, result)
}

func (a *Aggregate) Match(c Condition) *Aggregate {
	a.pipeline = append(a.pipeline, bson.M{
		"$match": c.Filters(),
	})
	return a
}

func (a *Aggregate) Group() *Aggregate {
	panic("")
}
