package pie

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Cursor cursor wrapper, providing type-safe iteration functionality
type Cursor[T any] struct {
	cursor *mongo.Cursor
	ctx    context.Context
}

// NewCursor create new cursor wrapper
func NewCursor[T any](ctx context.Context, cursor *mongo.Cursor) *Cursor[T] {
	return &Cursor[T]{
		cursor: cursor,
		ctx:    ctx,
	}
}

// Next move to next document
// return true means successfully moved to next document, false means no more documents
func (c *Cursor[T]) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

// TryNext try to move to next document
// return true means successfully moved to next document
// return false means no more documents in current batch (but may have more batches)
func (c *Cursor[T]) TryNext(ctx context.Context) bool {
	return c.cursor.TryNext(ctx)
}

// Decode decode current document to specified struct
func (c *Cursor[T]) Decode(val *T) error {
	return c.cursor.Decode(val)
}

// All get all remaining documents
func (c *Cursor[T]) All(ctx context.Context) ([]T, error) {
	var results []T
	err := c.cursor.All(ctx, &results)
	return results, err
}

// Close close cursor
func (c *Cursor[T]) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}

// Err return cursor error (if any)
func (c *Cursor[T]) Err() error {
	return c.cursor.Err()
}

// ID return cursor ID
func (c *Cursor[T]) ID() int64 {
	return c.cursor.ID()
}

// RemainingBatchLength return remaining document count in current batch
func (c *Cursor[T]) RemainingBatchLength() int {
	return c.cursor.RemainingBatchLength()
}

// SetBatchSize set document count per batch
func (c *Cursor[T]) SetBatchSize(batchSize int32) {
	c.cursor.SetBatchSize(batchSize)
}

// SetComment set cursor comment
func (c *Cursor[T]) SetComment(comment string) {
	c.cursor.SetComment(comment)
}

// Iterate iterate all documents, execute specified function for each document
func (c *Cursor[T]) Iterate(ctx context.Context, fn func(*T) error) error {
	defer c.Close(ctx)

	for c.Next(ctx) {
		var doc T
		if err := c.Decode(&doc); err != nil {
			return fmt.Errorf("failed to decode document: %w", err)
		}

		if err := fn(&doc); err != nil {
			return fmt.Errorf("iterator function failed: %w", err)
		}
	}

	if err := c.Err(); err != nil {
		return fmt.Errorf("cursor error: %w", err)
	}

	return nil
}

// ToSlice convert all documents to slice
func (c *Cursor[T]) ToSlice(ctx context.Context) ([]T, error) {
	return c.All(ctx)
}

// Count count remaining document count (requires full iteration)
func (c *Cursor[T]) Count(ctx context.Context) (int, error) {
	defer c.Close(ctx)

	count := 0
	for c.Next(ctx) {
		count++
	}

	if err := c.Err(); err != nil {
		return 0, fmt.Errorf("cursor error: %w", err)
	}

	return count, nil
}

// First get first document
func (c *Cursor[T]) First(ctx context.Context) (*T, error) {
	defer c.Close(ctx)

	if !c.Next(ctx) {
		if err := c.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
		return nil, ErrEmptyResult
	}

	var doc T
	if err := c.Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	return &doc, nil
}

// Map apply transformation function to each document, return slice of new type
func Map[T any, R any](ctx context.Context, cursor *Cursor[T], fn func(*T) (*R, error)) ([]R, error) {
	defer cursor.Close(ctx)

	var results []R
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		result, err := fn(&doc)
		if err != nil {
			return nil, fmt.Errorf("map function failed: %w", err)
		}

		if result != nil {
			results = append(results, *result)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}

// Filter filter documents, only keep documents that satisfy condition
func Filter[T any](ctx context.Context, cursor *Cursor[T], fn func(*T) bool) ([]T, error) {
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		if fn(&doc) {
			results = append(results, doc)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}

// ForEach execute specified operation for each document
func (c *Cursor[T]) ForEach(ctx context.Context, fn func(*T) error) error {
	return c.Iterate(ctx, fn)
}

// Skip skip specified number of documents
func (c *Cursor[T]) Skip(ctx context.Context, n int) error {
	for i := 0; i < n && c.Next(ctx); i++ {
		// 跳过文档
	}
	return c.Err()
}

// Take get at most n documents
func (c *Cursor[T]) Take(ctx context.Context, n int) ([]T, error) {
	defer c.Close(ctx)

	results := make([]T, 0, n)
	count := 0

	for c.Next(ctx) && count < n {
		var doc T
		if err := c.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, doc)
		count++
	}

	if err := c.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}
