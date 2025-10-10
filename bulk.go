package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// BulkWrite bulk write operation wrapper
type BulkWrite[T any] struct {
	engine     *Engine
	collection *mongo.Collection
	models     []mongo.WriteModel
	skipHooks  bool
}

// NewBulkWrite create new bulk write operation
func NewBulkWrite[T any](engine *Engine) *BulkWrite[T] {
	return &BulkWrite[T]{
		engine: engine,
		models: make([]mongo.WriteModel, 0),
	}
}

// CollectionForStruct set target collection (by struct)
func (b *BulkWrite[T]) CollectionForStruct(v interface{}) *BulkWrite[T] {
	collection, err := b.engine.CollectionForStruct(v)
	if err == nil {
		b.collection = collection
	}
	return b
}

// Collection set target collection (by name)
func (b *BulkWrite[T]) Collection(name string) *BulkWrite[T] {
	b.collection = b.engine.Collection(name)
	return b
}

// SkipHooks skip all hooks
func (b *BulkWrite[T]) SkipHooks() *BulkWrite[T] {
	b.skipHooks = true
	return b
}

// InsertOne add insert operation
func (b *BulkWrite[T]) InsertOne(doc *T) *BulkWrite[T] {
	model := mongo.NewInsertOneModel().SetDocument(doc)
	b.models = append(b.models, model)
	return b
}

// UpdateOne add update single document operation
func (b *BulkWrite[T]) UpdateOne(filter interface{}, update interface{}) *BulkWrite[T] {
	model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
	b.models = append(b.models, model)
	return b
}

// UpdateMany add update multiple documents operation
func (b *BulkWrite[T]) UpdateMany(filter interface{}, update interface{}) *BulkWrite[T] {
	model := mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update)
	b.models = append(b.models, model)
	return b
}

// ReplaceOne add replace document operation
func (b *BulkWrite[T]) ReplaceOne(filter interface{}, replacement *T) *BulkWrite[T] {
	model := mongo.NewReplaceOneModel().SetFilter(filter).SetReplacement(replacement)
	b.models = append(b.models, model)
	return b
}

// DeleteOne add delete single document operation
func (b *BulkWrite[T]) DeleteOne(filter interface{}) *BulkWrite[T] {
	model := mongo.NewDeleteOneModel().SetFilter(filter)
	b.models = append(b.models, model)
	return b
}

// DeleteMany add delete multiple documents operation
func (b *BulkWrite[T]) DeleteMany(filter interface{}) *BulkWrite[T] {
	model := mongo.NewDeleteManyModel().SetFilter(filter)
	b.models = append(b.models, model)
	return b
}

// Upsert set upsert option for last added Update operation
func (b *BulkWrite[T]) Upsert(upsert bool) *BulkWrite[T] {
	if len(b.models) > 0 {
		switch model := b.models[len(b.models)-1].(type) {
		case *mongo.UpdateOneModel:
			model.SetUpsert(upsert)
		case *mongo.UpdateManyModel:
			model.SetUpsert(upsert)
		case *mongo.ReplaceOneModel:
			model.SetUpsert(upsert)
		}
	}
	return b
}

// ArrayFilters set array filters for last added Update operation
func (b *BulkWrite[T]) ArrayFilters(filters []interface{}) *BulkWrite[T] {
	if len(b.models) > 0 {
		// Convert to []any for v2
		anyFilters := make([]any, len(filters))
		for i, filter := range filters {
			anyFilters[i] = filter
		}

		switch model := b.models[len(b.models)-1].(type) {
		case *mongo.UpdateOneModel:
			model.SetArrayFilters(anyFilters)
		case *mongo.UpdateManyModel:
			model.SetArrayFilters(anyFilters)
		}
	}
	return b
}

// Execute execute bulk write operation
func (b *BulkWrite[T]) Execute(ctx context.Context, ordered bool) (*BulkWriteResult, error) {
	start := time.Now()

	if b.collection == nil {
		collection, err := b.engine.CollectionForStruct((*T)(nil))
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}
		b.collection = collection
	}

	if len(b.models) == 0 {
		return &BulkWriteResult{}, nil
	}

	// Execute hooks (for each insert operation)
	if !b.skipHooks {
		for _, model := range b.models {
			switch m := model.(type) {
			case *mongo.InsertOneModel:
				doc := m.Document
				if err := b.engine.hooks.executeBeforeSave(ctx, doc); err != nil {
					return nil, err
				}
				if err := b.engine.hooks.executeBeforeCreate(ctx, doc); err != nil {
					return nil, err
				}
				if err := b.engine.hooks.executeModelBeforeCreate(ctx, doc); err != nil {
					return nil, err
				}
			case *mongo.UpdateOneModel, *mongo.UpdateManyModel:
				// For update operations, execute global hooks
				if err := b.engine.hooks.executeBeforeSave(ctx, nil); err != nil {
					return nil, err
				}
				if err := b.engine.hooks.executeBeforeUpdate(ctx, nil); err != nil {
					return nil, err
				}
			case *mongo.DeleteOneModel, *mongo.DeleteManyModel:
				// For delete operations, execute global hooks
				if err := b.engine.hooks.executeBeforeDelete(ctx, nil); err != nil {
					return nil, err
				}
			}
		}
	}

	// Execute bulk write
	isOrdered := ordered
	opts := options.BulkWrite().SetOrdered(isOrdered)
	result, err := b.collection.BulkWrite(ctx, b.models, opts)

	// Record query log
	if b.engine.queryLogger.IsEnabled() {
		b.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: b.collection.Name(),
			Operation:  "bulkWrite",
			Document:   fmt.Sprintf("%d operations", len(b.models)),
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// Execute After hooks
	if err == nil && !b.skipHooks {
		for _, model := range b.models {
			switch m := model.(type) {
			case *mongo.InsertOneModel:
				doc := m.Document
				b.engine.hooks.executeModelAfterCreate(ctx, doc)
				b.engine.hooks.executeAfterCreate(ctx, doc)
				b.engine.hooks.executeAfterSave(ctx, doc)
			case *mongo.UpdateOneModel, *mongo.UpdateManyModel:
				b.engine.hooks.executeAfterUpdate(ctx, nil)
				b.engine.hooks.executeAfterSave(ctx, nil)
			case *mongo.DeleteOneModel, *mongo.DeleteManyModel:
				b.engine.hooks.executeAfterDelete(ctx, nil)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("bulk write failed: %w", err)
	}

	return FromMongoBulkWriteResult(result), nil
}

// ExecuteOrdered execute ordered bulk write
func (b *BulkWrite[T]) ExecuteOrdered(ctx context.Context) (*BulkWriteResult, error) {
	return b.Execute(ctx, true)
}

// ExecuteUnordered execute unordered bulk write
func (b *BulkWrite[T]) ExecuteUnordered(ctx context.Context) (*BulkWriteResult, error) {
	return b.Execute(ctx, false)
}

// BulkWriteResult bulk write result
type BulkWriteResult struct {
	InsertedCount int64
	MatchedCount  int64
	ModifiedCount int64
	DeletedCount  int64
	UpsertedCount int64
	UpsertedIDs   map[int64]interface{}
}

// FromMongoBulkWriteResult convert from MongoDB's BulkWriteResult
func FromMongoBulkWriteResult(result *mongo.BulkWriteResult) *BulkWriteResult {
	if result == nil {
		return &BulkWriteResult{}
	}

	return &BulkWriteResult{
		InsertedCount: result.InsertedCount,
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		DeletedCount:  result.DeletedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedIDs:   result.UpsertedIDs,
	}
}
