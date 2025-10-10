package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SoftDeleteable soft delete interface
type SoftDeleteable interface {
	SetDeletedAt(deletedAt time.Time)
	GetDeletedAt() *time.Time
	IsDeleted() bool
}

// SoftDeleteManager soft delete manager
type SoftDeleteManager struct {
	engine         *Engine
	collection     *mongo.Collection
	deletedAtField string
}

// NewSoftDeleteManager create soft delete manager
func NewSoftDeleteManager(engine *Engine) *SoftDeleteManager {
	return &SoftDeleteManager{
		engine:         engine,
		deletedAtField: "deleted_at",
	}
}

// Collection set target collection
func (sdm *SoftDeleteManager) Collection(name string) *SoftDeleteManager {
	sdm.collection = sdm.engine.Collection(name)
	return sdm
}

// CollectionForStruct set target collection by struct
func (sdm *SoftDeleteManager) CollectionForStruct(v interface{}) *SoftDeleteManager {
	collection, err := sdm.engine.CollectionForStruct(v)
	if err == nil {
		sdm.collection = collection
	}
	return sdm
}

// SetDeletedAtField set deleted time field name
func (sdm *SoftDeleteManager) SetDeletedAtField(fieldName string) *SoftDeleteManager {
	sdm.deletedAtField = fieldName
	return sdm
}

// SoftDelete soft delete document
func (sdm *SoftDeleteManager) SoftDelete(ctx context.Context, filter bson.D) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: sdm.deletedAtField, Value: time.Now()},
		}},
	}

	_, err := sdm.collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDeleteMany soft delete multiple documents
func (sdm *SoftDeleteManager) SoftDeleteMany(ctx context.Context, filter bson.D) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: sdm.deletedAtField, Value: time.Now()},
		}},
	}

	_, err := sdm.collection.UpdateMany(ctx, filter, update)
	return err
}

// SoftDeleteByID soft delete document by ID
func (sdm *SoftDeleteManager) SoftDeleteByID(ctx context.Context, id bson.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	return sdm.SoftDelete(ctx, filter)
}

// Restore restore soft deleted document
func (sdm *SoftDeleteManager) Restore(ctx context.Context, filter bson.D) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	update := bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: sdm.deletedAtField, Value: ""},
		}},
	}

	_, err := sdm.collection.UpdateOne(ctx, filter, update)
	return err
}

// RestoreMany restore multiple soft deleted documents
func (sdm *SoftDeleteManager) RestoreMany(ctx context.Context, filter bson.D) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	update := bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: sdm.deletedAtField, Value: ""},
		}},
	}

	_, err := sdm.collection.UpdateMany(ctx, filter, update)
	return err
}

// RestoreByID restore soft deleted document by ID
func (sdm *SoftDeleteManager) RestoreByID(ctx context.Context, id bson.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	return sdm.Restore(ctx, filter)
}

// FindWithDeleted find documents including deleted ones
func (sdm *SoftDeleteManager) FindWithDeleted(ctx context.Context, filter bson.D) ([]bson.M, error) {
	if sdm.collection == nil {
		return nil, fmt.Errorf("collection not set")
	}

	cursor, err := sdm.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, doc)
	}

	return results, nil
}

// FindOnlyDeleted only find deleted documents
func (sdm *SoftDeleteManager) FindOnlyDeleted(ctx context.Context, filter bson.D) ([]bson.M, error) {
	if sdm.collection == nil {
		return nil, fmt.Errorf("collection not set")
	}

	// Add delete condition
	deletedFilter := bson.D{
		{Key: "$and", Value: []bson.D{
			filter,
			{{Key: sdm.deletedAtField, Value: bson.D{{Key: "$exists", Value: true}}}},
		}},
	}

	return sdm.FindWithDeleted(ctx, deletedFilter)
}

// FindOnlyActive only find non-deleted documents
func (sdm *SoftDeleteManager) FindOnlyActive(ctx context.Context, filter bson.D) ([]bson.M, error) {
	if sdm.collection == nil {
		return nil, fmt.Errorf("collection not set")
	}

	// Add non-delete condition
	activeFilter := bson.D{
		{Key: "$and", Value: []bson.D{
			filter,
			{{Key: sdm.deletedAtField, Value: bson.D{{Key: "$exists", Value: false}}}},
		}},
	}

	return sdm.FindWithDeleted(ctx, activeFilter)
}

// CountWithDeleted count documents including deleted ones
func (sdm *SoftDeleteManager) CountWithDeleted(ctx context.Context, filter bson.D) (int64, error) {
	if sdm.collection == nil {
		return 0, fmt.Errorf("collection not set")
	}

	return sdm.collection.CountDocuments(ctx, filter)
}

// CountOnlyDeleted count deleted documents
func (sdm *SoftDeleteManager) CountOnlyDeleted(ctx context.Context, filter bson.D) (int64, error) {
	if sdm.collection == nil {
		return 0, fmt.Errorf("collection not set")
	}

	deletedFilter := bson.D{
		{Key: "$and", Value: []bson.D{
			filter,
			{{Key: sdm.deletedAtField, Value: bson.D{{Key: "$exists", Value: true}}}},
		}},
	}

	return sdm.collection.CountDocuments(ctx, deletedFilter)
}

// CountOnlyActive count non-deleted documents
func (sdm *SoftDeleteManager) CountOnlyActive(ctx context.Context, filter bson.D) (int64, error) {
	if sdm.collection == nil {
		return 0, fmt.Errorf("collection not set")
	}

	activeFilter := bson.D{
		{Key: "$and", Value: []bson.D{
			filter,
			{{Key: sdm.deletedAtField, Value: bson.D{{Key: "$exists", Value: false}}}},
		}},
	}

	return sdm.collection.CountDocuments(ctx, activeFilter)
}

// HardDelete hard delete document (permanent delete)
func (sdm *SoftDeleteManager) HardDelete(ctx context.Context, filter bson.D) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	_, err := sdm.collection.DeleteOne(ctx, filter)
	return err
}

// HardDeleteMany hard delete multiple documents (permanent delete)
func (sdm *SoftDeleteManager) HardDeleteMany(ctx context.Context, filter bson.D) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	_, err := sdm.collection.DeleteMany(ctx, filter)
	return err
}

// HardDeleteByID hard delete document by ID (permanent delete)
func (sdm *SoftDeleteManager) HardDeleteByID(ctx context.Context, id bson.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	return sdm.HardDelete(ctx, filter)
}

// CleanupDeleted clean up deleted documents older than specified time
func (sdm *SoftDeleteManager) CleanupDeleted(ctx context.Context, olderThan time.Duration) error {
	if sdm.collection == nil {
		return fmt.Errorf("collection not set")
	}

	cutoffTime := time.Now().Add(-olderThan)
	filter := bson.D{
		{Key: sdm.deletedAtField, Value: bson.D{{Key: "$lt", Value: cutoffTime}}},
	}

	_, err := sdm.collection.DeleteMany(ctx, filter)
	return err
}

// SoftDeleteSession soft delete session
type SoftDeleteSession[T SoftDeleteable] struct {
	session *Session[T]
	manager *SoftDeleteManager
}

// NewSoftDeleteSession create soft delete session
func NewSoftDeleteSession[T SoftDeleteable](session *Session[T]) *SoftDeleteSession[T] {
	return &SoftDeleteSession[T]{
		session: session,
		manager: NewSoftDeleteManager(session.engine).Collection(session.collection.Name()),
	}
}

// SoftDelete soft delete current query document
func (sds *SoftDeleteSession[T]) SoftDelete(ctx context.Context) error {
	filter := sds.session.query.Build()
	return sds.manager.SoftDelete(ctx, filter)
}

// SoftDeleteMany soft delete current query multiple documents
func (sds *SoftDeleteSession[T]) SoftDeleteMany(ctx context.Context) error {
	filter := sds.session.query.Build()
	return sds.manager.SoftDeleteMany(ctx, filter)
}

// SoftDeleteByID soft delete document by ID
func (sds *SoftDeleteSession[T]) SoftDeleteByID(ctx context.Context, id bson.ObjectID) error {
	return sds.manager.SoftDeleteByID(ctx, id)
}

// Restore restore current query document
func (sds *SoftDeleteSession[T]) Restore(ctx context.Context) error {
	filter := sds.session.query.Build()
	return sds.manager.Restore(ctx, filter)
}

// RestoreMany restore current query multiple documents
func (sds *SoftDeleteSession[T]) RestoreMany(ctx context.Context) error {
	filter := sds.session.query.Build()
	return sds.manager.RestoreMany(ctx, filter)
}

// RestoreByID restore document by ID
func (sds *SoftDeleteSession[T]) RestoreByID(ctx context.Context, id bson.ObjectID) error {
	return sds.manager.RestoreByID(ctx, id)
}

// FindWithDeleted find documents including deleted ones
func (sds *SoftDeleteSession[T]) FindWithDeleted(ctx context.Context, results *[]T) error {
	filter := sds.session.query.Build()
	docs, err := sds.manager.FindWithDeleted(ctx, filter)
	if err != nil {
		return err
	}

	*results = make([]T, 0, len(docs))
	for _, doc := range docs {
		var result T
		bsonBytes, err := bson.Marshal(doc)
		if err != nil {
			return err
		}
		if err := bson.Unmarshal(bsonBytes, &result); err != nil {
			return err
		}
		*results = append(*results, result)
	}

	return nil
}

// FindOnlyDeleted only find deleted documents
func (sds *SoftDeleteSession[T]) FindOnlyDeleted(ctx context.Context, results *[]T) error {
	filter := sds.session.query.Build()
	docs, err := sds.manager.FindOnlyDeleted(ctx, filter)
	if err != nil {
		return err
	}

	*results = make([]T, 0, len(docs))
	for _, doc := range docs {
		var result T
		bsonBytes, err := bson.Marshal(doc)
		if err != nil {
			return err
		}
		if err := bson.Unmarshal(bsonBytes, &result); err != nil {
			return err
		}
		*results = append(*results, result)
	}

	return nil
}

// FindOnlyActive only find non-deleted documents
func (sds *SoftDeleteSession[T]) FindOnlyActive(ctx context.Context, results *[]T) error {
	filter := sds.session.query.Build()
	docs, err := sds.manager.FindOnlyActive(ctx, filter)
	if err != nil {
		return err
	}

	*results = make([]T, 0, len(docs))
	for _, doc := range docs {
		var result T
		bsonBytes, err := bson.Marshal(doc)
		if err != nil {
			return err
		}
		if err := bson.Unmarshal(bsonBytes, &result); err != nil {
			return err
		}
		*results = append(*results, result)
	}

	return nil
}

// CountWithDeleted count documents including deleted ones
func (sds *SoftDeleteSession[T]) CountWithDeleted(ctx context.Context) (int64, error) {
	filter := sds.session.query.Build()
	return sds.manager.CountWithDeleted(ctx, filter)
}

// CountOnlyDeleted count deleted documents
func (sds *SoftDeleteSession[T]) CountOnlyDeleted(ctx context.Context) (int64, error) {
	filter := sds.session.query.Build()
	return sds.manager.CountOnlyDeleted(ctx, filter)
}

// CountOnlyActive count non-deleted documents
func (sds *SoftDeleteSession[T]) CountOnlyActive(ctx context.Context) (int64, error) {
	filter := sds.session.query.Build()
	return sds.manager.CountOnlyActive(ctx, filter)
}

// HardDelete hard delete current query document
func (sds *SoftDeleteSession[T]) HardDelete(ctx context.Context) error {
	filter := sds.session.query.Build()
	return sds.manager.HardDelete(ctx, filter)
}

// HardDeleteMany hard delete current query multiple documents
func (sds *SoftDeleteSession[T]) HardDeleteMany(ctx context.Context) error {
	filter := sds.session.query.Build()
	return sds.manager.HardDeleteMany(ctx, filter)
}

// HardDeleteByID hard delete document by ID
func (sds *SoftDeleteSession[T]) HardDeleteByID(ctx context.Context, id bson.ObjectID) error {
	return sds.manager.HardDeleteByID(ctx, id)
}

// CleanupDeleted clean up deleted documents older than specified time
func (sds *SoftDeleteSession[T]) CleanupDeleted(ctx context.Context, olderThan time.Duration) error {
	return sds.manager.CleanupDeleted(ctx, olderThan)
}

// SoftDeleteOptions soft delete options
type SoftDeleteOptions struct {
	DeletedAtField string
	AutoCleanup    bool
	CleanupAfter   time.Duration
}

// NewSoftDeleteOptions create soft delete options
func NewSoftDeleteOptions() *SoftDeleteOptions {
	return &SoftDeleteOptions{
		DeletedAtField: "deleted_at",
		AutoCleanup:    false,
		CleanupAfter:   30 * 24 * time.Hour, // Retain soft-deleted documents for 30 days by default
	}
}

// WithDeletedAtField set deleted time field name
func (sdo *SoftDeleteOptions) WithDeletedAtField(fieldName string) *SoftDeleteOptions {
	sdo.DeletedAtField = fieldName
	return sdo
}

// WithAutoCleanup enable auto cleanup
func (sdo *SoftDeleteOptions) WithAutoCleanup() *SoftDeleteOptions {
	sdo.AutoCleanup = true
	return sdo
}

// WithCleanupAfter set cleanup time
func (sdo *SoftDeleteOptions) WithCleanupAfter(duration time.Duration) *SoftDeleteOptions {
	sdo.CleanupAfter = duration
	return sdo
}

// Build build soft delete options
func (sdo *SoftDeleteOptions) Build() *SoftDeleteOptions {
	return sdo
}

// SoftDeleteHook soft delete hook
type SoftDeleteHook struct {
	manager *SoftDeleteManager
	options *SoftDeleteOptions
}

// NewSoftDeleteHook create soft delete hook
func NewSoftDeleteHook(manager *SoftDeleteManager, options *SoftDeleteOptions) *SoftDeleteHook {
	return &SoftDeleteHook{
		manager: manager,
		options: options,
	}
}

// BeforeDelete delete before hook
func (sdh *SoftDeleteHook) BeforeDelete(ctx context.Context, doc interface{}) error {
	// Can add validation logic before delete here
	return nil
}

// AfterDelete delete after hook
func (sdh *SoftDeleteHook) AfterDelete(ctx context.Context, doc interface{}) error {
	// Can add cleanup logic after delete here
	if sdh.options.AutoCleanup {
		go func() {
			time.Sleep(sdh.options.CleanupAfter)
			sdh.manager.CleanupDeleted(ctx, sdh.options.CleanupAfter)
		}()
	}
	return nil
}

// SoftDeleteValidator soft delete validator
type SoftDeleteValidator struct{}

// NewSoftDeleteValidator create soft delete validator
func NewSoftDeleteValidator() *SoftDeleteValidator {
	return &SoftDeleteValidator{}
}

// ValidateSoftDeleteable validate object implements soft delete interface
func (sdv *SoftDeleteValidator) ValidateSoftDeleteable(obj interface{}) error {
	if _, ok := obj.(SoftDeleteable); !ok {
		return fmt.Errorf("object does not implement SoftDeleteable interface")
	}
	return nil
}

// ValidateDeletedAtField validate deleted time field
func (sdv *SoftDeleteValidator) ValidateDeletedAtField(fieldName string) error {
	if fieldName == "" {
		return fmt.Errorf("deletedAt field name cannot be empty")
	}
	return nil
}

// ValidateCleanupDuration validate cleanup time
func (sdv *SoftDeleteValidator) ValidateCleanupDuration(duration time.Duration) error {
	if duration <= 0 {
		return fmt.Errorf("cleanup duration must be positive")
	}
	return nil
}

// ValidateAll validate all soft delete options
func (sdv *SoftDeleteValidator) ValidateAll(options *SoftDeleteOptions) error {
	if err := sdv.ValidateDeletedAtField(options.DeletedAtField); err != nil {
		return fmt.Errorf("deletedAt field validation failed: %w", err)
	}

	if options.AutoCleanup {
		if err := sdv.ValidateCleanupDuration(options.CleanupAfter); err != nil {
			return fmt.Errorf("cleanup duration validation failed: %w", err)
		}
	}

	return nil
}

// SoftDeleteBuilder soft delete builder
type SoftDeleteBuilder struct {
	manager *SoftDeleteManager
	options *SoftDeleteOptions
}

// NewSoftDeleteBuilder create soft delete builder
func NewSoftDeleteBuilder(manager *SoftDeleteManager) *SoftDeleteBuilder {
	return &SoftDeleteBuilder{
		manager: manager,
		options: NewSoftDeleteOptions(),
	}
}

// WithDeletedAtField set deleted time field name
func (sdb *SoftDeleteBuilder) WithDeletedAtField(fieldName string) *SoftDeleteBuilder {
	sdb.options.DeletedAtField = fieldName
	sdb.manager.SetDeletedAtField(fieldName)
	return sdb
}

// WithAutoCleanup enable auto cleanup
func (sdb *SoftDeleteBuilder) WithAutoCleanup() *SoftDeleteBuilder {
	sdb.options.AutoCleanup = true
	return sdb
}

// WithCleanupAfter set cleanup time
func (sdb *SoftDeleteBuilder) WithCleanupAfter(duration time.Duration) *SoftDeleteBuilder {
	sdb.options.CleanupAfter = duration
	return sdb
}

// Build build soft delete options
func (sdb *SoftDeleteBuilder) Build() *SoftDeleteOptions {
	return sdb.options
}

// SoftDeleteStats soft delete stats
type SoftDeleteStats struct {
	TotalCount    int64
	ActiveCount   int64
	DeletedCount  int64
	DeletedRatio  float64
	OldestDeleted *time.Time
	NewestDeleted *time.Time
}

// GetSoftDeleteStats get soft delete stats
func (sdm *SoftDeleteManager) GetSoftDeleteStats(ctx context.Context) (*SoftDeleteStats, error) {
	if sdm.collection == nil {
		return nil, fmt.Errorf("collection not set")
	}

	stats := &SoftDeleteStats{}

	// Total count
	totalCount, err := sdm.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	stats.TotalCount = totalCount

	// Active count
	activeCount, err := sdm.CountOnlyActive(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	stats.ActiveCount = activeCount

	// Deleted count
	deletedCount, err := sdm.CountOnlyDeleted(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	stats.DeletedCount = deletedCount

	// Deleted ratio
	if totalCount > 0 {
		stats.DeletedRatio = float64(deletedCount) / float64(totalCount)
	}

	// Oldest and newest deleted time
	if deletedCount > 0 {
		pipeline := []bson.D{
			{{Key: "$match", Value: bson.D{
				{Key: sdm.deletedAtField, Value: bson.D{{Key: "$exists", Value: true}}},
			}}},
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "oldest", Value: bson.D{{Key: "$min", Value: "$" + sdm.deletedAtField}}},
				{Key: "newest", Value: bson.D{{Key: "$max", Value: "$" + sdm.deletedAtField}}},
			}}},
		}

		cursor, err := sdm.collection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		if cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err != nil {
				return nil, err
			}

			if oldest, ok := result["oldest"].(time.Time); ok {
				stats.OldestDeleted = &oldest
			}
			if newest, ok := result["newest"].(time.Time); ok {
				stats.NewestDeleted = &newest
			}
		}
	}

	return stats, nil
}
