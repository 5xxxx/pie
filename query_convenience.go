package pie

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ========== Query convenience methods ==========

// FindByID find single document by ID
func (s *Session[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	return s.Where("_id", id).FindOne(ctx)
}

// FindByIDs find documents by multiple IDs
func (s *Session[T]) FindByIDs(ctx context.Context, ids interface{}) ([]T, error) {
	return s.WhereIn("_id", ids).Find(ctx)
}

// FirstOne find first document (simplified version, no error when not found)
func (s *Session[T]) FirstOne(ctx context.Context) (*T, error) {
	s.query.Limit(1)
	result, err := s.FindOne(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// Exists check if document exists that matches conditions
func (s *Session[T]) Exists(ctx context.Context) (bool, error) {
	s.query.Limit(1)
	count, err := s.CountDocuments(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAndCount query and count total (for pagination)
func (s *Session[T]) FindAndCount(ctx context.Context) ([]T, int64, error) {
	// First count total
	total, err := s.CountDocuments(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Then query data
	results, err := s.Find(ctx)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Pluck extract single field to slice
func (s *Session[T]) Pluck(ctx context.Context, field string, results interface{}) error {
	// Set projection, only return specified field
	s.query.Project(bson.D{{Key: field, Value: 1}})

	if s.initErr != nil {
		return fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()
	opts := s.query.BuildFindOptions()

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Use reflection to parse results
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr || resultsValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("results must be a pointer to slice")
	}

	sliceValue := resultsValue.Elem()
	sliceType := sliceValue.Type().Elem()

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		if value, ok := doc[field]; ok {
			// Create new element and set value
			elemValue := reflect.New(sliceType).Elem()
			elemValue.Set(reflect.ValueOf(value).Convert(sliceType))
			sliceValue.Set(reflect.Append(sliceValue, elemValue))
		}
	}

	return cursor.Err()
}

// Value get value of single field
func (s *Session[T]) Value(ctx context.Context, field string, result interface{}) error {
	s.query.Project(bson.D{{Key: field, Value: 1}})
	s.query.Limit(1)

	if s.initErr != nil {
		return fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()
	opts := s.query.BuildFindOneOptions()

	var doc bson.M
	err := s.collection.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		return err
	}

	if value, ok := doc[field]; ok {
		resultValue := reflect.ValueOf(result)
		if resultValue.Kind() != reflect.Ptr {
			return fmt.Errorf("result must be a pointer")
		}
		resultValue.Elem().Set(reflect.ValueOf(value).Convert(resultValue.Elem().Type()))
	}

	return nil
}

// Chunk process large dataset in chunks
func (s *Session[T]) Chunk(ctx context.Context, size int, callback func([]T) error) error {
	if s.initErr != nil {
		return fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()
	opts := s.query.BuildFindOptions()

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	chunk := make([]T, 0, size)
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		chunk = append(chunk, doc)

		if len(chunk) >= size {
			if err := callback(chunk); err != nil {
				return err
			}
			chunk = chunk[:0] // Reset slice
		}
	}

	// Process last batch of data
	if len(chunk) > 0 {
		if err := callback(chunk); err != nil {
			return err
		}
	}

	return cursor.Err()
}

// ========== Insert convenience methods ==========

// Create insert single (alias method)
func (s *Session[T]) Create(ctx context.Context, doc *T) (InsertResult, error) {
	return s.Insert(ctx, doc)
}

// CreateMany bulk insert
func (s *Session[T]) CreateMany(ctx context.Context, docs []T) ([]T, error) {
	_, err := s.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

// FirstOrCreate find or create
func (s *Session[T]) FirstOrCreate(ctx context.Context, doc *T) (*T, bool, error) {
	// First try to find
	result, err := s.FindOne(ctx)
	if err == nil {
		// Found, return existing document
		return result, false, nil
	}

	if err != mongo.ErrNoDocuments {
		// Query error
		return nil, false, err
	}

	// Not found, create new document
	_, err = s.Insert(ctx, doc)
	if err != nil {
		return nil, false, err
	}

	return doc, true, nil
}

// UpdateOrCreate update or create (Upsert)
func (s *Session[T]) UpdateOrCreate(ctx context.Context, doc *T) (*T, error) {
	if s.collection == nil {
		collection, err := s.engine.CollectionForStruct((*T)(nil))
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}
		s.collection = collection
	}

	filter := s.query.GetFilter()

	// Use $set operator
	update := bson.D{{Key: "$set", Value: doc}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result T
	err := s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ========== Update convenience methods ==========

// UpdateColumn update single field
func (s *Session[T]) UpdateColumn(ctx context.Context, field string, value interface{}) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	_, err := s.UpdateMany(ctx, update)
	return err
}

// UpdateColumns update multiple fields
func (s *Session[T]) UpdateColumns(ctx context.Context, data map[string]interface{}) error {
	update := bson.D{{Key: "$set", Value: data}}
	_, err := s.UpdateMany(ctx, update)
	return err
}

// Increment field increment
func (s *Session[T]) Increment(ctx context.Context, field string, value interface{}) error {
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: field, Value: value}}}}
	_, err := s.UpdateMany(ctx, update)
	return err
}

// Decrement field decrement
func (s *Session[T]) Decrement(ctx context.Context, field string, value interface{}) error {
	// Convert value to negative for decrement
	var negValue interface{}
	switch v := value.(type) {
	case int:
		negValue = -v
	case int32:
		negValue = -v
	case int64:
		negValue = -v
	case float32:
		negValue = -v
	case float64:
		negValue = -v
	default:
		return fmt.Errorf("unsupported type for decrement: %T", value)
	}

	update := bson.D{{Key: "$inc", Value: bson.D{{Key: field, Value: negValue}}}}
	_, err := s.UpdateMany(ctx, update)
	return err
}

// Toggle boolean toggle
func (s *Session[T]) Toggle(ctx context.Context, field string) error {
	if s.initErr != nil {
		return fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// Use aggregation pipeline to toggle
	update := mongo.Pipeline{
		{{Key: "$set", Value: bson.D{{Key: field, Value: bson.D{{Key: "$not", Value: "$" + field}}}}}},
	}

	_, err := s.collection.UpdateMany(ctx, filter, update)
	return err
}

// ========== Delete convenience methods ==========

// DeleteByID delete by ID
func (s *Session[T]) DeleteByID(ctx context.Context, id interface{}) error {
	_, err := s.Where("_id", id).Delete(ctx)
	return err
}

// DeleteByIDs delete by multiple IDs
func (s *Session[T]) DeleteByIDs(ctx context.Context, ids interface{}) (int64, error) {
	result, err := s.WhereIn("_id", ids).DeleteMany(ctx)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// Destroy delete (alias method)
func (s *Session[T]) Destroy(ctx context.Context) error {
	_, err := s.DeleteMany(ctx)
	return err
}

// ========== Count convenience methods ==========

// QuickCount quick count
func (s *Session[T]) QuickCount(ctx context.Context) (int64, error) {
	return s.CountDocuments(ctx)
}

// Sum sum
func (s *Session[T]) Sum(ctx context.Context, field string) (float64, error) {
	if s.collection == nil {
		collection, err := s.engine.CollectionForStruct((*T)(nil))
		if err != nil {
			return 0, fmt.Errorf("failed to get collection: %w", err)
		}
		s.collection = collection
	}

	filter := s.query.GetFilter()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$" + field}}},
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total float64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}

	return 0, nil
}

// Avg average
func (s *Session[T]) Avg(ctx context.Context, field string) (float64, error) {
	if s.collection == nil {
		collection, err := s.engine.CollectionForStruct((*T)(nil))
		if err != nil {
			return 0, fmt.Errorf("failed to get collection: %w", err)
		}
		s.collection = collection
	}

	filter := s.query.GetFilter()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "average", Value: bson.D{{Key: "$avg", Value: "$" + field}}},
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Average float64 `bson:"average"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Average, nil
	}

	return 0, nil
}

// MaxValue max value
func (s *Session[T]) MaxValue(ctx context.Context, field string) (interface{}, error) {
	if s.collection == nil {
		collection, err := s.engine.CollectionForStruct((*T)(nil))
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}
		s.collection = collection
	}

	filter := s.query.GetFilter()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "max", Value: bson.D{{Key: "$max", Value: "$" + field}}},
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Max interface{} `bson:"max"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		return result.Max, nil
	}

	return nil, nil
}

// MinValue min value
func (s *Session[T]) MinValue(ctx context.Context, field string) (interface{}, error) {
	if s.collection == nil {
		collection, err := s.engine.CollectionForStruct((*T)(nil))
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}
		s.collection = collection
	}

	filter := s.query.GetFilter()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "min", Value: bson.D{{Key: "$min", Value: "$" + field}}},
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Min interface{} `bson:"min"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		return result.Min, nil
	}

	return nil, nil
}
