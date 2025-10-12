package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Session type-safe MongoDB session
type Session[T any] struct {
	engine      *Engine
	collection  *mongo.Collection
	query       *Query
	options     *SessionOptions
	skipHooks   bool
	cacheConfig *SessionCacheConfig
	initErr     error // Added: save initialization error
}

// SessionOptions session options
type SessionOptions struct {
	FindOptions       *options.FindOptionsBuilder
	FindOneOptions    *options.FindOneOptionsBuilder
	UpdateOneOptions  *options.UpdateOneOptionsBuilder
	UpdateManyOptions *options.UpdateManyOptionsBuilder
	DeleteOneOptions  *options.DeleteOneOptionsBuilder
	DeleteManyOptions *options.DeleteManyOptionsBuilder
	InsertOptions     *options.InsertOneOptionsBuilder
}

// NewSessionOptions create new session options
func NewSessionOptions() *SessionOptions {
	return &SessionOptions{
		FindOptions:       options.Find(),
		FindOneOptions:    options.FindOne(),
		UpdateOneOptions:  options.UpdateOne(),
		UpdateManyOptions: options.UpdateMany(),
		DeleteOneOptions:  options.DeleteOne(),
		DeleteManyOptions: options.DeleteMany(),
		InsertOptions:     options.InsertOne(),
	}
}

// Find finds multiple documents
func (s *Session[T]) Find(ctx context.Context) ([]T, error) {
	start := time.Now()

	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()
	opts := s.query.BuildFindOptions()

	// Execute query
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "find",
			Filter:     filter,
			Options:    opts,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// Parse results
	var results []T
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	// Execute AfterFind hooks (for each document)
	if !s.skipHooks && err == nil {
		for i := range results {
			s.engine.hooks.executeModelAfterFind(ctx, &results[i])
			s.engine.hooks.executeAfterFind(ctx, &results[i])
		}
	}

	return results, nil
}

// FindOne find a single document
func (s *Session[T]) FindOne(ctx context.Context) (*T, error) {
	start := time.Now()

	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()
	opts := s.query.BuildFindOneOptions()

	// Execute query
	singleResult := s.collection.FindOne(ctx, filter, opts)
	err := singleResult.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrEmptyResult
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "findOne",
			Filter:     filter,
			Options:    opts,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// Parse results
	var result T
	err = singleResult.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	// Execute AfterFind hook
	if !s.skipHooks && err == nil {
		s.engine.hooks.executeModelAfterFind(ctx, &result)
		s.engine.hooks.executeAfterFind(ctx, &result)
	}

	return &result, nil
}

// Insert insert a single document
func (s *Session[T]) Insert(ctx context.Context, doc *T) (InsertResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return InsertResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	// 1. Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeSave(ctx, doc); err != nil {
			return InsertResult{}, err
		}
		if err := s.engine.hooks.executeBeforeCreate(ctx, doc); err != nil {
			return InsertResult{}, err
		}
		if err := s.engine.hooks.executeModelBeforeCreate(ctx, doc); err != nil {
			return InsertResult{}, err
		}
	}

	// 2. Execute insert
	result, err := s.collection.InsertOne(ctx, doc, s.options.InsertOptions)

	// 3. Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "insertOne",
			Document:   doc,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// 4. Execute After hooks
	if err == nil && !s.skipHooks {
		s.engine.hooks.executeModelAfterCreate(ctx, doc)
		s.engine.hooks.executeAfterCreate(ctx, doc)
		s.engine.hooks.executeAfterSave(ctx, doc)
	}

	return FromMongoInsertResult(result, err), err
}

// InsertMany insert multiple documents
func (s *Session[T]) InsertMany(ctx context.Context, docs []T) (InsertManyResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return InsertManyResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	// 1. Execute Before hooks (for each document)
	if !s.skipHooks {
		for i := range docs {
			if err := s.engine.hooks.executeBeforeSave(ctx, &docs[i]); err != nil {
				return InsertManyResult{}, err
			}
			if err := s.engine.hooks.executeBeforeCreate(ctx, &docs[i]); err != nil {
				return InsertManyResult{}, err
			}
			if err := s.engine.hooks.executeModelBeforeCreate(ctx, &docs[i]); err != nil {
				return InsertManyResult{}, err
			}
		}
	}

	// 2. Convert to any slice
	interfaceDocs := make([]any, len(docs))
	for i, doc := range docs {
		interfaceDocs[i] = doc
	}

	// 3. Execute insert
	result, err := s.collection.InsertMany(ctx, interfaceDocs)

	// 4. Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "insertMany",
			Document:   docs,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// 5. Execute After hooks (for each document)
	if err == nil && !s.skipHooks {
		for i := range docs {
			s.engine.hooks.executeModelAfterCreate(ctx, &docs[i])
			s.engine.hooks.executeAfterCreate(ctx, &docs[i])
			s.engine.hooks.executeAfterSave(ctx, &docs[i])
		}
	}

	return FromMongoInsertManyResult(result, err), err
}

// SkipHooks skip all hooks
func (s *Session[T]) SkipHooks() *Session[T] {
	s.skipHooks = true
	return s
}

// Update update document
func (s *Session[T]) Update(ctx context.Context, update bson.D) (UpdateResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return UpdateResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// 1. Execute Before hooks (for update operation, we cannot get the original document, so skip model hooks)
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeSave(ctx, nil); err != nil {
			return UpdateResult{}, err
		}
		if err := s.engine.hooks.executeBeforeUpdate(ctx, nil); err != nil {
			return UpdateResult{}, err
		}
	}

	// 2. Execute update
	result, err := s.collection.UpdateOne(ctx, filter, update, s.options.UpdateOneOptions)

	// 3. Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "updateOne",
			Filter:     filter,
			Update:     update,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// 4. Execute After hooks
	if err == nil && !s.skipHooks {
		s.engine.hooks.executeAfterUpdate(ctx, nil)
		s.engine.hooks.executeAfterSave(ctx, nil)
	}

	return FromMongoUpdateResult(result, err), err
}

// UpdateMany update multiple documents
func (s *Session[T]) UpdateMany(ctx context.Context, update bson.D) (UpdateResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return UpdateResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// 1. Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeSave(ctx, nil); err != nil {
			return UpdateResult{}, err
		}
		if err := s.engine.hooks.executeBeforeUpdate(ctx, nil); err != nil {
			return UpdateResult{}, err
		}
	}

	// 2. Execute update
	result, err := s.collection.UpdateMany(ctx, filter, update, s.options.UpdateManyOptions)

	// 3. Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "updateMany",
			Filter:     filter,
			Update:     update,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// 4. Execute After hooks
	if err == nil && !s.skipHooks {
		s.engine.hooks.executeAfterUpdate(ctx, nil)
		s.engine.hooks.executeAfterSave(ctx, nil)
	}

	return FromMongoUpdateResult(result, err), err
}

// Delete delete document
func (s *Session[T]) Delete(ctx context.Context) (DeleteResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return DeleteResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// 1. Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeDelete(ctx, nil); err != nil {
			return DeleteResult{}, err
		}
	}

	// 2. Execute delete
	result, err := s.collection.DeleteOne(ctx, filter, s.options.DeleteOneOptions)

	// 3. Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "deleteOne",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// 4. Execute After hooks
	if err == nil && !s.skipHooks {
		s.engine.hooks.executeAfterDelete(ctx, nil)
	}

	return FromMongoDeleteResult(result, err), err
}

// DeleteMany delete multiple documents
func (s *Session[T]) DeleteMany(ctx context.Context) (DeleteResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return DeleteResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// 1. Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeDelete(ctx, nil); err != nil {
			return DeleteResult{}, err
		}
	}

	// 2. Execute delete
	result, err := s.collection.DeleteMany(ctx, filter, s.options.DeleteManyOptions)

	// 3. Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "deleteMany",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// 4. Execute After hooks
	if err == nil && !s.skipHooks {
		s.engine.hooks.executeAfterDelete(ctx, nil)
	}

	return FromMongoDeleteResult(result, err), err
}

// Count count documents
func (s *Session[T]) Count(ctx context.Context) (int64, error) {
	if s.initErr != nil {
		return 0, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	count, err := s.collection.CountDocuments(ctx, s.query.GetFilter())
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return count, nil
}

// Distinct get distinct values
func (s *Session[T]) Distinct(ctx context.Context, field string) ([]any, error) {
	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	result := s.collection.Distinct(ctx, field, s.query.GetFilter())
	var values []any
	err := result.Decode(&values)
	if err != nil {
		return nil, fmt.Errorf("failed to decode distinct values: %w", err)
	}

	return values, nil
}

// Chaining methods

// Where add condition
func (s *Session[T]) Where(field string, value any) *Session[T] {
	s.query.Where(field, value)
	return s
}

// WhereOperator use operator to add condition
func (s *Session[T]) WhereOperator(op Operator) *Session[T] {
	s.query.WhereOperator(op)
	return s
}

// And add AND condition
func (s *Session[T]) And(operators ...Operator) *Session[T] {
	s.query.And(operators...)
	return s
}

// Or add OR condition
func (s *Session[T]) Or(operators ...Operator) *Session[T] {
	s.query.Or(operators...)
	return s
}

// OrderBy add sort
func (s *Session[T]) OrderBy(field string) *Session[T] {
	s.query.OrderBy(field)
	return s
}

// OrderByDesc add descending sort
func (s *Session[T]) OrderByDesc(field string) *Session[T] {
	s.query.OrderByDesc(field)
	return s
}

// Limit set limit number
func (s *Session[T]) Limit(limit int64) *Session[T] {
	s.query.Limit(limit)
	return s
}

// Skip set skip number
func (s *Session[T]) Skip(skip int64) *Session[T] {
	s.query.Skip(skip)
	return s
}

// Project set projection fields
func (s *Session[T]) Project(project bson.D) *Session[T] {
	s.query.Project(project)
	return s
}

// Select select fields
func (s *Session[T]) Select(fields ...string) *Session[T] {
	s.query.Select(fields...)
	return s
}

// Exclude exclude fields
func (s *Session[T]) Exclude(fields ...string) *Session[T] {
	s.query.Exclude(fields...)
	return s
}

// Clone clone session
func (s *Session[T]) Clone() *Session[T] {
	return &Session[T]{
		engine:     s.engine,
		collection: s.collection,
		query:      s.query.Clone(),
		options:    s.options, // share options
		initErr:    s.initErr, // copy initialization error
	}
}

// Clear clear query conditions
func (s *Session[T]) Clear() *Session[T] {
	s.query.Clear()
	return s
}

// GetQuery get query builder
func (s *Session[T]) GetQuery() *Query {
	return s.query
}

// GetOptions get session options
func (s *Session[T]) GetOptions() *SessionOptions {
	return s.options
}

// FindCursor return query cursor, support iteration access
func (s *Session[T]) FindCursor(ctx context.Context) (*Cursor[T], error) {
	start := time.Now()

	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()
	opts := s.query.BuildFindOptions()

	// Execute query
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "find",
			Filter:     filter,
			Options:    opts,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return NewCursor[T](ctx, cursor), nil
}

// CountDocuments accurately count the number of documents that match the filter
func (s *Session[T]) CountDocuments(ctx context.Context) (int64, error) {
	start := time.Now()

	if s.initErr != nil {
		return 0, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	count, err := s.collection.CountDocuments(ctx, filter)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "countDocuments",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return count, err
}

// EstimatedDocumentCount quickly estimate the total number of documents in the collection (ignoring filter conditions)
func (s *Session[T]) EstimatedDocumentCount(ctx context.Context) (int64, error) {
	start := time.Now()

	if s.initErr != nil {
		return 0, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	count, err := s.collection.EstimatedDocumentCount(ctx)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "estimatedDocumentCount",
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return count, err
}

// ReplaceOne replace a single document
func (s *Session[T]) ReplaceOne(ctx context.Context, replacement *T) (UpdateResult, error) {
	start := time.Now()

	if s.initErr != nil {
		return UpdateResult{}, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeSave(ctx, replacement); err != nil {
			return UpdateResult{}, err
		}
		if err := s.engine.hooks.executeBeforeUpdate(ctx, replacement); err != nil {
			return UpdateResult{}, err
		}
		if err := s.engine.hooks.executeModelBeforeUpdate(ctx, replacement); err != nil {
			return UpdateResult{}, err
		}
	}

	// Execute replace
	opts := options.Replace()
	result, err := s.collection.ReplaceOne(ctx, filter, replacement, opts)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "replaceOne",
			Filter:     filter,
			Document:   replacement,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	// Execute After hooks
	if err == nil && !s.skipHooks {
		s.engine.hooks.executeModelAfterUpdate(ctx, replacement)
		s.engine.hooks.executeAfterUpdate(ctx, replacement)
		s.engine.hooks.executeAfterSave(ctx, replacement)
	}

	return FromMongoUpdateResult(result, err), err
}

// FindOneAndDelete find and delete a single document, return the document before deletion
func (s *Session[T]) FindOneAndDelete(ctx context.Context) (*T, error) {
	start := time.Now()

	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	var result T

	// Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeDelete(ctx, &result); err != nil {
			return nil, err
		}
	}

	// Execute find and delete
	singleResult := s.collection.FindOneAndDelete(ctx, filter)
	err := singleResult.Decode(&result)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "findOneAndDelete",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrEmptyResult
		}
		return nil, fmt.Errorf("failed to find and delete document: %w", err)
	}

	// Execute After hooks
	if !s.skipHooks {
		s.engine.hooks.executeModelAfterDelete(ctx, &result)
		s.engine.hooks.executeAfterDelete(ctx, &result)
	}

	return &result, nil
}

// FindOneAndReplace find and replace a single document, return the document before or after replacement
func (s *Session[T]) FindOneAndReplace(ctx context.Context, replacement *T, returnAfter bool) (*T, error) {
	start := time.Now()

	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeSave(ctx, replacement); err != nil {
			return nil, err
		}
		if err := s.engine.hooks.executeBeforeUpdate(ctx, replacement); err != nil {
			return nil, err
		}
		if err := s.engine.hooks.executeModelBeforeUpdate(ctx, replacement); err != nil {
			return nil, err
		}
	}

	// Set return options
	opts := options.FindOneAndReplace()
	if returnAfter {
		opts.SetReturnDocument(options.After)
	} else {
		opts.SetReturnDocument(options.Before)
	}
	// 设置upsert为false，确保只替换现有文档
	opts.SetUpsert(false)

	var result T

	// Execute find and replace
	singleResult := s.collection.FindOneAndReplace(ctx, filter, replacement, opts)
	err := singleResult.Decode(&result)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "findOneAndReplace",
			Filter:     filter,
			Document:   replacement,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrEmptyResult
		}
		return nil, fmt.Errorf("failed to find and replace document: %w", err)
	}

	// Execute After hooks
	if !s.skipHooks {
		if returnAfter {
			s.engine.hooks.executeModelAfterUpdate(ctx, &result)
			s.engine.hooks.executeAfterUpdate(ctx, &result)
			s.engine.hooks.executeAfterSave(ctx, &result)
		}
	}

	return &result, nil
}

// FindOneAndUpdate find and update a single document, return the document before or after update
func (s *Session[T]) FindOneAndUpdate(ctx context.Context, update bson.D, returnAfter bool) (*T, error) {
	start := time.Now()

	if s.initErr != nil {
		return nil, fmt.Errorf("session initialization failed: %w", s.initErr)
	}

	filter := s.query.GetFilter()

	// Execute Before hooks
	if !s.skipHooks {
		if err := s.engine.hooks.executeBeforeSave(ctx, nil); err != nil {
			return nil, err
		}
		if err := s.engine.hooks.executeBeforeUpdate(ctx, nil); err != nil {
			return nil, err
		}
	}

	// Set return options
	opts := options.FindOneAndUpdate()
	if returnAfter {
		opts.SetReturnDocument(options.After)
	} else {
		opts.SetReturnDocument(options.Before)
	}

	var result T

	// Execute find and update
	singleResult := s.collection.FindOneAndUpdate(ctx, filter, update, opts)
	err := singleResult.Decode(&result)

	// Record query log
	if s.engine.queryLogger.IsEnabled() {
		s.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: s.collection.Name(),
			Operation:  "findOneAndUpdate",
			Filter:     filter,
			Update:     update,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrEmptyResult
		}
		return nil, fmt.Errorf("failed to find and update document: %w", err)
	}

	// Execute After hooks
	if !s.skipHooks {
		if returnAfter {
			s.engine.hooks.executeModelAfterUpdate(ctx, &result)
			s.engine.hooks.executeAfterUpdate(ctx, &result)
			s.engine.hooks.executeAfterSave(ctx, &result)
		}
	}

	return &result, nil
}

// Advanced query features - merged from advanced_query.go

// Hint set index hint
func (s *Session[T]) Hint(hint any) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetHint(hint)
	return s
}

// Comment set query comment
func (s *Session[T]) Comment(comment string) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetComment(comment)
	return s
}

// BatchSize set batch size
func (s *Session[T]) BatchSize(size int32) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetBatchSize(size)
	return s
}

// NoCursorTimeout set cursor timeout
func (s *Session[T]) NoCursorTimeout(noTimeout bool) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetNoCursorTimeout(noTimeout)
	return s
}

// ReturnKey set return key
func (s *Session[T]) ReturnKey(returnKey bool) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetReturnKey(returnKey)
	return s
}

// ShowRecordId set show record ID
func (s *Session[T]) ShowRecordId(show bool) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetShowRecordID(show)
	return s
}

// Min set min key
func (s *Session[T]) Min(min bson.D) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetMin(min)
	return s
}

// Max set max key
func (s *Session[T]) Max(max bson.D) *Session[T] {
	if s.options.FindOptions == nil {
		s.options.FindOptions = options.Find()
	}
	s.options.FindOptions.SetMax(max)
	return s
}

// ArrayFilters set array filters (for update operations)
func (s *Session[T]) ArrayFilters(filters []any) *Session[T] {
	// Convert to []any for v2
	anyFilters := make([]any, len(filters))
	for i, filter := range filters {
		anyFilters[i] = filter
	}

	if s.options.UpdateOneOptions == nil {
		s.options.UpdateOneOptions = options.UpdateOne()
	}
	s.options.UpdateOneOptions.SetArrayFilters(anyFilters)

	if s.options.UpdateManyOptions == nil {
		s.options.UpdateManyOptions = options.UpdateMany()
	}
	s.options.UpdateManyOptions.SetArrayFilters(anyFilters)

	return s
}

// Let set Let variables (for update operations)
func (s *Session[T]) Let(variables bson.D) *Session[T] {
	if s.options.UpdateOneOptions == nil {
		s.options.UpdateOneOptions = options.UpdateOne()
	}
	s.options.UpdateOneOptions.SetLet(variables)

	if s.options.UpdateManyOptions == nil {
		s.options.UpdateManyOptions = options.UpdateMany()
	}
	s.options.UpdateManyOptions.SetLet(variables)

	return s
}

// Upsert set Upsert (for update operations)
func (s *Session[T]) Upsert(upsert bool) *Session[T] {
	if s.options.UpdateOneOptions == nil {
		s.options.UpdateOneOptions = options.UpdateOne()
	}
	s.options.UpdateOneOptions.SetUpsert(upsert)

	if s.options.UpdateManyOptions == nil {
		s.options.UpdateManyOptions = options.UpdateMany()
	}
	s.options.UpdateManyOptions.SetUpsert(upsert)

	return s
}
