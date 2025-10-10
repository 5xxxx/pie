package pie

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Result query result
type Result[T any] struct {
	Data  T
	Error error
}

// InsertResult insert result
type InsertResult struct {
	InsertedID bson.ObjectID
	Error      error
}

// InsertManyResult bulk insert result
type InsertManyResult struct {
	InsertedIDs []bson.ObjectID
	Error       error
}

// UpdateResult update result
type UpdateResult struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    bson.ObjectID
	Error         error
}

// DeleteResult delete result
type DeleteResult struct {
	DeletedCount int64
	Error        error
}

// AggregateResult aggregation result
type AggregateResult[T any] struct {
	Data  []T
	Error error
}

// CountResult count result
type CountResult struct {
	Count int64
	Error error
}

// DistinctResult distinct result
type DistinctResult[T any] struct {
	Values []T
	Error  error
}

// PaginationResult pagination result
type PaginationResult[T any] struct {
	Data  []T
	Total int64
	Page  int64
	Size  int64
	Error error
}

// TransactionResult transaction result
type TransactionResult[T any] struct {
	Data  T
	Error error
}

// Convert from MongoDB official results

// FromMongoInsertResult convert from MongoDB insert result
func FromMongoInsertResult(result *mongo.InsertOneResult, err error) InsertResult {
	if err != nil {
		return InsertResult{Error: err}
	}

	var insertedID bson.ObjectID
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		insertedID = oid
	}

	return InsertResult{
		InsertedID: insertedID,
		Error:      nil,
	}
}

// FromMongoInsertManyResult convert from MongoDB bulk insert result
func FromMongoInsertManyResult(result *mongo.InsertManyResult, err error) InsertManyResult {
	if err != nil {
		return InsertManyResult{Error: err}
	}

	var insertedIDs []bson.ObjectID
	for _, id := range result.InsertedIDs {
		if oid, ok := id.(bson.ObjectID); ok {
			insertedIDs = append(insertedIDs, oid)
		}
	}

	return InsertManyResult{
		InsertedIDs: insertedIDs,
		Error:       nil,
	}
}

// FromMongoUpdateResult convert from MongoDB update result
func FromMongoUpdateResult(result *mongo.UpdateResult, err error) UpdateResult {
	if err != nil {
		return UpdateResult{Error: err}
	}

	var upsertedID bson.ObjectID
	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(bson.ObjectID); ok {
			upsertedID = oid
		}
	}

	return UpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    upsertedID,
		Error:         nil,
	}
}

// FromMongoDeleteResult convert from MongoDB delete result
func FromMongoDeleteResult(result *mongo.DeleteResult, err error) DeleteResult {
	if err != nil {
		return DeleteResult{Error: err}
	}

	return DeleteResult{
		DeletedCount: result.DeletedCount,
		Error:        nil,
	}
}
