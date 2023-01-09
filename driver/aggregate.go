package driver

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Aggregate interface {
	One(result interface{}, ctx ...context.Context) error
	All(result interface{}, ctx ...context.Context) error
	// SetAllowDiskUse sets the value for the AllowDiskUse field.
	SetAllowDiskUse(b bool) Aggregate

	// SetBatchSize sets the value for the BatchSize field.
	SetBatchSize(i int32) Aggregate

	// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
	SetBypassDocumentValidation(b bool) Aggregate

	// SetCollation sets the value for the Collation field.
	SetCollation(c *options.Collation) Aggregate

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) Aggregate

	// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
	SetMaxAwaitTime(d time.Duration) Aggregate

	// SetComment sets the value for the Comment field.
	SetComment(s string) Aggregate

	// SetHint sets the value for the Hint field.
	SetHint(h interface{}) Aggregate

	Pipeline(pipeline bson.A) Aggregate

	Match(c Condition) Aggregate

	SetDatabase(db string) Aggregate

	Collection(doc interface{}) Aggregate

	SetCollReadPreference(rp *readpref.ReadPref) Aggregate

	SetCollRegistry(r *bsoncodec.Registry) Aggregate

	SetCollWriteConcern(wc *writeconcern.WriteConcern) Aggregate
	
	SetReadConcern(rc *readconcern.ReadConcern) Aggregate
	//AddFields() Aggregate
	//Bucket() Aggregate
	//BucketAuto() Aggregate
	//CollStats() Aggregate
	//Count() Aggregate
	//CurrentOp() Aggregate
	//Facet() Aggregate
	//GeoNear() Aggregate
	//GraphLookup() Aggregate
	//Group() Aggregate
	//IndexStats() Aggregate
	//Limit() Aggregate
	//ListLocalSession() Aggregate
	//ListSession() Aggregate
	//Lookup() Aggregate
	//Match(filter Session) Aggregate
	//Merge() Aggregate
	//Out() Aggregate
	//PlanCacheStats() Aggregate
	//Project() Aggregate
	//Redact() Aggregate
	//ReplaceRoot() Aggregate
	//ReplaceWith() Aggregate
	//Sample() Aggregate
	//Set() Aggregate
	//Skip() Aggregate
	//Sort() Aggregate
	//SortByCount() Aggregate
	//UnionWith() Aggregate
	//Unset() Aggregate
	//Unwind() Aggregate
	//All(result interface{}) error
}
