package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Aggregate aggregation operation builder
type Aggregate[T any] struct {
	engine     *Engine
	collection *mongo.Collection
	pipeline   bson.A
	options    *options.AggregateOptionsBuilder
}

// NewAggregate create new aggregation operation
func NewAggregate[T any](engine *Engine) *Aggregate[T] {
	return &Aggregate[T]{
		engine:   engine,
		pipeline: bson.A{},
		options:  options.Aggregate(),
	}
}

// Collection set collection
func (a *Aggregate[T]) Collection(name string) *Aggregate[T] {
	a.collection = a.engine.Collection(name)
	return a
}

// CollectionForStruct set collection by struct
func (a *Aggregate[T]) CollectionForStruct(v interface{}) *Aggregate[T] {
	collection, err := a.engine.CollectionForStruct(v)
	if err != nil {
		// Here we can log the error but don't interrupt the chain call
		return a
	}
	a.collection = collection
	return a
}

// Match add match stage
func (a *Aggregate[T]) Match(filter bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$match", filter}})
	return a
}

// MatchOperator add match stage using operator
func (a *Aggregate[T]) MatchOperator(op Operator) *Aggregate[T] {
	return a.Match(op.ToBSON())
}

// Group add group stage
func (a *Aggregate[T]) Group(group bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$group", group}})
	return a
}

// Sort add sort stage
func (a *Aggregate[T]) Sort(sort bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$sort", sort}})
	return a
}

// Limit add limit stage
func (a *Aggregate[T]) Limit(limit int64) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$limit", limit}})
	return a
}

// Skip add skip stage
func (a *Aggregate[T]) Skip(skip int64) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$skip", skip}})
	return a
}

// Project add projection stage
func (a *Aggregate[T]) Project(project bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$project", project}})
	return a
}

// AddFields add fields stage
func (a *Aggregate[T]) AddFields(fields bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$addFields", fields}})
	return a
}

// Lookup add lookup stage
func (a *Aggregate[T]) Lookup(from, localField, foreignField, as string) *Aggregate[T] {
	lookup := bson.D{
		{"from", from},
		{"localField", localField},
		{"foreignField", foreignField},
		{"as", as},
	}
	a.pipeline = append(a.pipeline, bson.D{{"$lookup", lookup}})
	return a
}

// Unwind add unwind stage
func (a *Aggregate[T]) Unwind(path string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$unwind", path}})
	return a
}

// Facet add facet stage
func (a *Aggregate[T]) Facet(facet bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$facet", facet}})
	return a
}

// Count add count stage
func (a *Aggregate[T]) Count(field string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$count", field}})
	return a
}

// Sample add sample stage
func (a *Aggregate[T]) Sample(size int64) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$sample", bson.D{{"size", size}}}})
	return a
}

// ReplaceRoot add replace root stage
func (a *Aggregate[T]) ReplaceRoot(newRoot bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$replaceRoot", bson.D{{"newRoot", newRoot}}}})
	return a
}

// ReplaceWith add replace stage
func (a *Aggregate[T]) ReplaceWith(replacement bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$replaceWith", replacement}})
	return a
}

// Out add out stage
func (a *Aggregate[T]) Out(collection string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$out", collection}})
	return a
}

// Merge add merge stage
func (a *Aggregate[T]) Merge(into string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$merge", bson.D{{"into", into}}}})
	return a
}

// Pipeline add custom pipeline stage
func (a *Aggregate[T]) Pipeline(stage bson.D) *Aggregate[T] {
	a.pipeline = append(a.pipeline, stage)
	return a
}

// Option setting methods

// SetAllowDiskUse set allow disk use
func (a *Aggregate[T]) SetAllowDiskUse(allow bool) *Aggregate[T] {
	a.options.SetAllowDiskUse(allow)
	return a
}

// SetBatchSize set batch size
func (a *Aggregate[T]) SetBatchSize(size int32) *Aggregate[T] {
	a.options.SetBatchSize(size)
	return a
}

// SetBypassDocumentValidation set bypass document validation
func (a *Aggregate[T]) SetBypassDocumentValidation(bypass bool) *Aggregate[T] {
	a.options.SetBypassDocumentValidation(bypass)
	return a
}

// SetCollation set collation
func (a *Aggregate[T]) SetCollation(collation *options.Collation) *Aggregate[T] {
	a.options.SetCollation(collation)
	return a
}

// SetMaxAwaitTime set max await time
func (a *Aggregate[T]) SetMaxAwaitTime(duration int64) *Aggregate[T] {
	a.options.SetMaxAwaitTime(time.Duration(duration))
	return a
}

// SetComment set comment
func (a *Aggregate[T]) SetComment(comment string) *Aggregate[T] {
	a.options.SetComment(comment)
	return a
}

// SetHint set hint
func (a *Aggregate[T]) SetHint(hint interface{}) *Aggregate[T] {
	a.options.SetHint(hint)
	return a
}

// Exec execute aggregation operation
func (a *Aggregate[T]) Exec(ctx context.Context) (*AggregateResult[T], error) {
	if a.collection == nil {
		// Try to get collection based on generic type
		var zero T
		collection, err := a.engine.CollectionForStruct(zero)
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}
		a.collection = collection
	}

	cursor, err := a.collection.Aggregate(ctx, a.pipeline, a.options)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	return &AggregateResult[T]{
		Data:  results,
		Error: nil,
	}, nil
}

// ExecOne execute aggregation operation and return single result
func (a *Aggregate[T]) ExecOne(ctx context.Context, result *T) error {
	if a.collection == nil {
		var zero T
		collection, err := a.engine.CollectionForStruct(zero)
		if err != nil {
			return fmt.Errorf("failed to get collection: %w", err)
		}
		a.collection = collection
	}

	cursor, err := a.collection.Aggregate(ctx, a.pipeline, a.options)
	if err != nil {
		return fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(result); err != nil {
			return fmt.Errorf("failed to decode aggregation result: %w", err)
		}
		return nil
	}

	return ErrEmptyResult
}

// GetPipeline get pipeline
func (a *Aggregate[T]) GetPipeline() bson.A {
	return a.pipeline
}

// Clone clone aggregation operation
func (a *Aggregate[T]) Clone() *Aggregate[T] {
	newPipeline := make(bson.A, len(a.pipeline))
	copy(newPipeline, a.pipeline)

	return &Aggregate[T]{
		engine:     a.engine,
		collection: a.collection,
		pipeline:   newPipeline,
		options:    a.options,
	}
}

// Clear clear pipeline
func (a *Aggregate[T]) Clear() *Aggregate[T] {
	a.pipeline = bson.A{}
	return a
}
