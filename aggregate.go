package pie

import (
	"context"
	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Aggregate represents an interface for performing aggregation operations on a MongoDB collection.
// One executes the aggregation operation and stores the result in the provided result variable.
// It returns an error if the operation fails.
type Aggregate interface {
	One(result any, ctx ...context.Context) error
	All(result any, ctx ...context.Context) error
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
	SetHint(h any) Aggregate

	Pipeline(pipeline bson.A) Aggregate

	Match(c Condition) Aggregate

	SetDatabase(db string) Aggregate

	Collection(doc any) Aggregate

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
	//All(result any) error
}

// aggregate represents an aggregation operation.
type aggregate struct {
	db       string
	doc      any
	engine   Client
	pipeline bson.A
	opts     []*options.AggregateOptions
	collOpts []*options.CollectionOptions
}

// NewAggregate creates a new instance of the Aggregate struct with the provided client as the engine.
func NewAggregate(engine Client) Aggregate {
	return &aggregate{engine: engine}
}

// One retrieves a single document from the MongoDB collection that matches the aggregation pipeline and decodes it into the provided result variable.
// If a context is provided, it will be used for the operation, otherwise a background context will be used.
// The method determines the appropriate collection based on the type of the result variable. If 'result' is a struct, it uses the `collectionForStruct` method to obtain the collection
func (a *aggregate) One(result any, ctx ...context.Context) error {

	var c context.Context
	if len(ctx) > 0 {
		c = ctx[0]
	} else {
		c = context.Background()
	}

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

	aggregate, err := coll.Aggregate(c, a.pipeline, a.opts...)
	if err != nil {
		return err
	}
	if next := aggregate.Next(c); next {
		if err := aggregate.Decode(result); err != nil {
			return err
		}
	} else {
		return mongo.ErrNoDocuments
	}
	return nil
}

// All retrieves all the documents from the MongoDB collection that match the aggregation pipeline and stores the result in the provided result variable.
// If a context is provided, it will be used for the operation, otherwise a background context will be used.
// The method determines the appropriate collection based on the type of the result variable. If 'result' is a struct, it uses the `collectionForStruct` method to obtain the collection
func (a *aggregate) All(result any, ctx ...context.Context) error {
	var c context.Context
	if len(ctx) > 0 {
		c = ctx[0]
	} else {
		c = context.Background()
	}

	var coll *mongo.Collection
	var err error
	if a.doc != nil {
		coll, err = a.collectionForStruct(a.doc)
	} else {
		coll, err = a.collectionForSlice(result)
	}
	if err != nil {
		return err
	}

	aggregate, err := coll.Aggregate(c, a.pipeline, a.opts...)
	if err != nil {
		return err
	}

	return aggregate.All(c, result)
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (a *aggregate) SetAllowDiskUse(b bool) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetAllowDiskUse(b))
	return a
}

// SetBatchSize sets the value for the BatchSize field.
func (a *aggregate) SetBatchSize(i int32) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetBatchSize(i))
	return a
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *aggregate) SetBypassDocumentValidation(b bool) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetBypassDocumentValidation(b))
	return a
}

// SetCollation sets the value for the Collation field.
func (a *aggregate) SetCollation(c *options.Collation) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetCollation(c))
	return a
}

// SetMaxTime sets the value for the MaxTime field.
func (a *aggregate) SetMaxTime(d time.Duration) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetMaxTime(d))
	return a
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *aggregate) SetMaxAwaitTime(d time.Duration) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetMaxAwaitTime(d))
	return a
}

// SetComment sets the value for the Comment field.
func (a *aggregate) SetComment(s string) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetComment(s))
	return a
}

// SetHint sets the value for the Hint field.
func (a *aggregate) SetHint(h any) Aggregate {
	a.opts = append(a.opts, options.Aggregate().SetHint(h))
	return a
}

// Pipeline appends a pipeline to the existing pipeline of the aggregate object.
// The pipeline parameter is a slice of BSON documents representing the stages of the aggregation pipeline.
// Each document in the pipeline slice represents a stage of the aggregation.
// The Pipeline method returns the aggregate object to allow for method chaining.
func (a *aggregate) Pipeline(pipeline bson.A) Aggregate {
	a.pipeline = append(a.pipeline, pipeline...)
	return a
}

// Match appends a $match stage to the pipeline based on the specified condition.
// It takes a Condition as input and retrieves the filter expressions from it using the Filters() method.
// The filter expressions are then added to the pipeline with the "$match" operator.
// The updated pipeline is stored in the aggregate struct.
// It returns the aggregate struct itself, allowing for method chaining.
func (a *aggregate) Match(c Condition) Aggregate {
	filters, err := c.Filters()
	if err != nil {
		panic(err)
	}
	a.pipeline = append(a.pipeline, bson.M{
		"$match": filters,
	})
	return a
}

// SetDatabase sets the value for the db field in the aggregate struct.
func (a *aggregate) SetDatabase(db string) Aggregate {
	a.db = db
	return a
}

// collectionForStruct returns a *mongo.Collection for the given document structure.
// It takes an input parameter 'doc' which represents the document structure.
// If 'a.doc' is not nil, it assigns the value returned by a.engine.CollectionNameForStruct(a.doc) to 'coll' and assigns the error to 'err',
// otherwise it assigns the value returned by a.engine.CollectionNameForStruct(doc) to 'coll' and assigns the error to 'err'.
// If 'err' is not nil, it returns nil and the error.
// Otherwise, it returns a.collectionByName(coll.Name) and nil.
func (a *aggregate) collectionForStruct(doc any) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if a.doc != nil {
		coll, err = a.engine.CollectionNameForStruct(a.doc)
	} else {
		coll, err = a.engine.CollectionNameForStruct(doc)
	}
	if err != nil {
		return nil, err
	}
	return a.collectionByName(coll.Name), nil
}

// collectionForSlice retrieves the collection by name for a given slice of documents.
func (a *aggregate) collectionForSlice(doc any) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if a.doc != nil {
		coll, err = a.engine.CollectionNameForStruct(a.doc)
	} else {
		coll, err = a.engine.CollectionNameForSlice(doc)
	}
	if err != nil {
		return nil, err
	}
	return a.collectionByName(coll.Name), nil
}

// collectionByName returns a *mongo.Collection for the given name. If the collection options (a.collOpts) is nil, it initializes it as an empty slice. It then calls a.engine.Collection
func (a *aggregate) collectionByName(name string) *mongo.Collection {
	if a.collOpts == nil {
		a.collOpts = make([]*options.CollectionOptions, 0)
	}
	return a.engine.Collection(name, a.collOpts, a.db)
}

// SetReadConcern sets the value for the ReadConcern field.
func (a *aggregate) SetReadConcern(rc *readconcern.ReadConcern) Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetReadConcern(rc))
	return a
}

// SetCollWriteConcern sets the value for the WriteConcern field in the CollectionOptions.
func (a *aggregate) SetCollWriteConcern(wc *writeconcern.WriteConcern) Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetWriteConcern(wc))
	return a
}

// SetCollReadPreference sets the value for the ReadPreference field.
func (a *aggregate) SetCollReadPreference(rp *readpref.ReadPref) Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetReadPreference(rp))
	return a
}

// SetCollRegistry sets the value for the Registry field in the CollectionOptions.
func (a *aggregate) SetCollRegistry(r *bsoncodec.Registry) Aggregate {
	a.collOpts = append(a.collOpts, options.Collection().SetRegistry(r))
	return a
}

// Collection sets the document to be used for the aggregate operation.
// It returns the updated aggregate object.
func (a *aggregate) Collection(doc any) Aggregate {
	a.doc = doc
	return a
}
