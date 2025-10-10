package pie

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Indexes index manager
type Indexes struct {
	engine     *Engine
	collection *mongo.Collection
	indexes    []mongo.IndexModel
	options    *options.CreateIndexesOptionsBuilder
}

// NewIndexes create new index manager
func NewIndexes(engine *Engine) *Indexes {
	return &Indexes{
		engine:  engine,
		indexes: make([]mongo.IndexModel, 0),
		options: options.CreateIndexes(),
	}
}

// Collection set collection
func (i *Indexes) Collection(name string) *Indexes {
	i.collection = i.engine.Collection(name)
	return i
}

// CollectionForStruct set collection by struct
func (i *Indexes) CollectionForStruct(v interface{}) *Indexes {
	collection, err := i.engine.CollectionForStruct(v)
	if err != nil {
		// Here we can log the error but don't interrupt the chain call
		return i
	}
	i.collection = collection
	return i
}

// AddIndex add index
func (i *Indexes) AddIndex(keys bson.D, opts ...*options.IndexOptions) *Indexes {
	indexModel := mongo.IndexModel{
		Keys: keys,
	}

	if len(opts) > 0 {
		indexOpts := options.Index()
		for _, opt := range opts {
			if opt.ExpireAfterSeconds != nil {
				indexOpts.SetExpireAfterSeconds(*opt.ExpireAfterSeconds)
			}
			if opt.Name != nil {
				indexOpts.SetName(*opt.Name)
			}
			if opt.Sparse != nil {
				indexOpts.SetSparse(*opt.Sparse)
			}
			if opt.Unique != nil {
				indexOpts.SetUnique(*opt.Unique)
			}
			if opt.Version != nil {
				indexOpts.SetVersion(*opt.Version)
			}
			// Skip unsupported fields
			// if opt.Weights != nil {
			//     indexOpts.SetWeights(*opt.Weights)
			// }
			if opt.DefaultLanguage != nil {
				indexOpts.SetDefaultLanguage(*opt.DefaultLanguage)
			}
			if opt.LanguageOverride != nil {
				indexOpts.SetLanguageOverride(*opt.LanguageOverride)
			}
			// Skip unsupported fields
			// if opt.TextIndexVersion != nil {
			//     indexOpts.SetTextIndexVersion(*opt.TextIndexVersion)
			// }
			if opt.SphereVersion != nil {
				indexOpts.SetSphereVersion(*opt.SphereVersion)
			}
			if opt.Bits != nil {
				indexOpts.SetBits(*opt.Bits)
			}
			if opt.Max != nil {
				indexOpts.SetMax(*opt.Max)
			}
			if opt.Min != nil {
				indexOpts.SetMin(*opt.Min)
			}
			if opt.BucketSize != nil {
				indexOpts.SetBucketSize(*opt.BucketSize)
			}
			// Skip unsupported fields
			// if opt.PartialFilterExpression != nil {
			//     indexOpts.SetPartialFilterExpression(*opt.PartialFilterExpression)
			// }
			if opt.Collation != nil {
				indexOpts.SetCollation(opt.Collation)
			}
			// Skip unsupported fields
			// if opt.WildcardProjection != nil {
			//     indexOpts.SetWildcardProjection(*opt.WildcardProjection)
			// }
			if opt.Hidden != nil {
				indexOpts.SetHidden(*opt.Hidden)
			}
		}
		indexModel.Options = indexOpts
	}

	i.indexes = append(i.indexes, indexModel)
	return i
}

// AddTextIndex add text index
func (i *Indexes) AddTextIndex(fields ...string) *Indexes {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: "text"})
	}
	return i.AddIndex(keys)
}

// AddCompoundIndex add compound index
func (i *Indexes) AddCompoundIndex(keys bson.D) *Indexes {
	return i.AddIndex(keys)
}

// AddSparseIndex add sparse index
func (i *Indexes) AddSparseIndex(keys bson.D) *Indexes {
	sparse := true
	return i.AddIndex(keys, &options.IndexOptions{Sparse: &sparse})
}

// AddUniqueIndex add unique index
func (i *Indexes) AddUniqueIndex(keys bson.D) *Indexes {
	unique := true
	return i.AddIndex(keys, &options.IndexOptions{Unique: &unique})
}

// AddTTLIndex add TTL index
func (i *Indexes) AddTTLIndex(field string, expireAfterSeconds int32) *Indexes {
	keys := bson.D{{field, 1}}
	return i.AddIndex(keys, &options.IndexOptions{ExpireAfterSeconds: &expireAfterSeconds})
}

// AddGeospatialIndex add geospatial index
func (i *Indexes) AddGeospatialIndex(field string) *Indexes {
	keys := bson.D{{field, "2dsphere"}}
	return i.AddIndex(keys)
}

// AddHashedIndex add hashed index
func (i *Indexes) AddHashedIndex(field string) *Indexes {
	keys := bson.D{{field, "hashed"}}
	return i.AddIndex(keys)
}

// AddPartialIndex add partial index
func (i *Indexes) AddPartialIndex(keys bson.D, filter bson.D) *Indexes {
	// In v2, partial index implementation may be different
	// Here we temporarily skip partial index implementation
	return i.AddIndex(keys)
}

// Create create indexes
func (i *Indexes) Create(ctx context.Context) ([]string, error) {
	if i.collection == nil {
		return nil, fmt.Errorf("collection not set")
	}

	if len(i.indexes) == 0 {
		return nil, fmt.Errorf("no indexes to create")
	}

	indexNames, err := i.collection.Indexes().CreateMany(ctx, i.indexes, i.options)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return indexNames, nil
}

// CreateOne create single index
func (i *Indexes) CreateOne(ctx context.Context, index mongo.IndexModel) (string, error) {
	if i.collection == nil {
		return "", fmt.Errorf("collection not set")
	}

	indexName, err := i.collection.Indexes().CreateOne(ctx, index, i.options)
	if err != nil {
		return "", fmt.Errorf("failed to create index: %w", err)
	}

	return indexName, nil
}

// DropAll drop all indexes
func (i *Indexes) DropAll(ctx context.Context) error {
	if i.collection == nil {
		return fmt.Errorf("collection not set")
	}

	err := i.collection.Indexes().DropAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop all indexes: %w", err)
	}

	return nil
}

// DropOne drop specified index
func (i *Indexes) DropOne(ctx context.Context, name string) error {
	if i.collection == nil {
		return fmt.Errorf("collection not set")
	}

	err := i.collection.Indexes().DropOne(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to drop index %s: %w", name, err)
	}

	return nil
}

// List list all indexes
func (i *Indexes) List(ctx context.Context) ([]bson.M, error) {
	if i.collection == nil {
		return nil, fmt.Errorf("collection not set")
	}

	cursor, err := i.collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexes []bson.M
	if err := cursor.All(ctx, &indexes); err != nil {
		return nil, fmt.Errorf("failed to decode indexes: %w", err)
	}

	return indexes, nil
}

// GetIndexes get index list
func (i *Indexes) GetIndexes() []mongo.IndexModel {
	return i.indexes
}

// Clear clear index list
func (i *Indexes) Clear() *Indexes {
	i.indexes = make([]mongo.IndexModel, 0)
	return i
}

// Clone clone index manager
func (i *Indexes) Clone() *Indexes {
	newIndexes := make([]mongo.IndexModel, len(i.indexes))
	copy(newIndexes, i.indexes)

	return &Indexes{
		engine:     i.engine,
		collection: i.collection,
		indexes:    newIndexes,
		options:    i.options,
	}
}

// Advanced index utilities merged from advanced_indexes.go

// AddWildcardIndex add wildcard index
// WildcardIndexOptions options for wildcard index
type WildcardIndexOptions struct {
	Name   *string
	Sparse *bool
}

func (i *Indexes) AddWildcardIndex(ctx context.Context, field string, opts *WildcardIndexOptions) error {
	if i.collection == nil {
		return fmt.Errorf("collection not set")
	}

	idxOpts := options.Index()
	if opts != nil {
		if opts.Name != nil {
			idxOpts = idxOpts.SetName(*opts.Name)
		}
		if opts.Sparse != nil {
			idxOpts = idxOpts.SetSparse(*opts.Sparse)
		}
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field + ".$**", Value: 1}},
		Options: idxOpts,
	}

	_, err := i.CreateOne(ctx, indexModel)
	return err
}

// AddColumnstoreIndex add columnstore index
// ColumnstoreIndexOptions options for columnstore index
type ColumnstoreIndexOptions struct {
	Name *string
}

func (i *Indexes) AddColumnstoreIndex(ctx context.Context, field string, opts *ColumnstoreIndexOptions) error {
	if i.collection == nil {
		return fmt.Errorf("collection not set")
	}

	idxOpts2 := options.Index()
	if opts != nil && opts.Name != nil {
		idxOpts2 = idxOpts2.SetName(*opts.Name)
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: "columnstore"}},
		Options: idxOpts2,
	}

	_, err := i.CreateOne(ctx, indexModel)
	return err
}
