package pie

import (
	"context"
	"time"

	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Indexes represents the interface for managing indexes in a database.
type Indexes interface {
	CreateIndexes(doc any, ctx ...context.Context) ([]string, error)
	DropAll(doc any, ctx ...context.Context) error
	DropOne(doc any, name string, ctx ...context.Context) error
	AddIndex(keys any, opt ...*options.IndexOptions) Indexes
	SetMaxTime(d time.Duration) Indexes
	SetCommitQuorumInt(quorum int32) Indexes
	SetCommitQuorumString(quorum string) Indexes
	SetCommitQuorumMajority() Indexes
	SetCommitQuorumVotingMembers() Indexes
	SetDatabase(db string) Indexes
	Collection(doc any) Indexes
}

type index struct {
	db                 string
	doc                any
	engine             Client
	indexes            []mongo.IndexModel
	createIndexOpts    []*options.CreateIndexesOptionsBuilder
	dropIndexesOptions []*options.DropIndexesOptionsBuilder
}

func NewIndexes(driver Client) Indexes {
	return &index{engine: driver}
}

func (i *index) CreateIndexes(doc any, ctx ...context.Context) ([]string, error) {
	coll, err := i.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	var opts []options.Lister[options.CreateIndexesOptions]
	for _, opt := range i.createIndexOpts {
		opts = append(opts, opt)
	}
	return coll.Indexes().CreateMany(c, i.indexes, opts...)
}

func (i *index) DropAll(doc any, ctx ...context.Context) error {
	coll, err := i.collectionForStruct(doc)
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	var opts []options.Lister[options.DropIndexesOptions]
	for _, opt := range i.dropIndexesOptions {
		opts = append(opts, opt)
	}
	return coll.Indexes().DropAll(c, opts...)
}

func (i *index) DropOne(doc any, name string, ctx ...context.Context) error {
	coll, err := i.collectionForStruct(doc)
	if err != nil {
		return err
	}
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	var opts []options.Lister[options.DropIndexesOptions]
	for _, opt := range i.dropIndexesOptions {
		opts = append(opts, opt)
	}
	return coll.Indexes().DropOne(c, name, opts...)
}

func (i *index) AddIndex(keys any, opt ...*options.IndexOptions) Indexes {
	m := mongo.IndexModel{
		Keys: keys,
	}

	if len(opt) > 0 {
		// Convert IndexOptions to IndexOptionsBuilder
		builder := options.Index()
		for _, v := range opt {
			if v.ExpireAfterSeconds != nil {
				builder.SetExpireAfterSeconds(*v.ExpireAfterSeconds)
			}
			if v.Name != nil {
				builder.SetName(*v.Name)
			}
			if v.Sparse != nil {
				builder.SetSparse(*v.Sparse)
			}
			if v.StorageEngine != nil {
				builder.SetStorageEngine(v.StorageEngine)
			}
			if v.Unique != nil {
				builder.SetUnique(*v.Unique)
			}
			if v.Version != nil {
				builder.SetVersion(*v.Version)
			}
			if v.DefaultLanguage != nil {
				builder.SetDefaultLanguage(*v.DefaultLanguage)
			}
			if v.LanguageOverride != nil {
				builder.SetLanguageOverride(*v.LanguageOverride)
			}
			if v.TextVersion != nil {
				builder.SetTextVersion(*v.TextVersion)
			}
			if v.Weights != nil {
				builder.SetWeights(v.Weights)
			}
			if v.SphereVersion != nil {
				builder.SetSphereVersion(*v.SphereVersion)
			}
			if v.Bits != nil {
				builder.SetBits(*v.Bits)
			}
			if v.Max != nil {
				builder.SetMax(*v.Max)
			}
			if v.Min != nil {
				builder.SetMin(*v.Min)
			}
			if v.BucketSize != nil {
				builder.SetBucketSize(*v.BucketSize)
			}
			if v.PartialFilterExpression != nil {
				builder.SetPartialFilterExpression(v.PartialFilterExpression)
			}
			if v.Collation != nil {
				builder.SetCollation(v.Collation)
			}
			if v.WildcardProjection != nil {
				builder.SetWildcardProjection(v.WildcardProjection)
			}
			if v.Hidden != nil {
				builder.SetHidden(*v.Hidden)
			}
		}
		m.Options = builder
	}

	i.indexes = append(i.indexes, m)
	return i
}

// SetMaxTime sets the value for the MaxTime field.
func (i *index) SetMaxTime(d time.Duration) Indexes {
	// Note: MaxTime is not available in CreateIndexes/DropIndexes options in v2
	// This method is kept for compatibility but does nothing
	return i
}

// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (i *index) SetCommitQuorumInt(quorum int32) Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumInt(quorum))
	return i
}

// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
func (i *index) SetCommitQuorumString(quorum string) Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumString(quorum))
	return i
}

// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (i *index) SetCommitQuorumMajority() Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumMajority())
	return i
}

// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (i *index) SetCommitQuorumVotingMembers() Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumVotingMembers())
	return i
}

func (i *index) SetDatabase(db string) Indexes {
	i.db = db
	return i
}

func (i *index) collectionForStruct(doc any) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if i.doc != nil {
		coll, err = i.engine.CollectionNameForStruct(i.doc)
	} else {
		coll, err = i.engine.CollectionNameForStruct(doc)
	}
	if err != nil {
		return nil, err
	}
	return i.collectionByName(coll.Name), nil
}

func (i *index) collectionForSlice(doc any) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if i.doc != nil {
		coll, err = i.engine.CollectionNameForStruct(i.doc)
	} else {
		coll, err = i.engine.CollectionNameForSlice(doc)
	}
	if err != nil {
		return nil, err
	}
	return i.collectionByName(coll.Name), nil
}

func (i *index) collectionByName(name string) *mongo.Collection {
	return i.engine.Collection(name, nil, i.db)
}

func (i *index) Collection(doc any) Indexes {
	i.doc = doc
	return i
}
