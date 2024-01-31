package pie

import (
	"context"
	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
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
	createIndexOpts    []*options.CreateIndexesOptions
	dropIndexesOptions []*options.DropIndexesOptions
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
	return coll.Indexes().CreateMany(c, i.indexes, i.createIndexOpts...)
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
	_, err = coll.Indexes().DropAll(c, i.dropIndexesOptions...)
	return err
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
	_, err = coll.Indexes().DropOne(c, name, i.dropIndexesOptions...)
	return err
}

func (i *index) AddIndex(keys any, opt ...*options.IndexOptions) Indexes {
	m := mongo.IndexModel{
		Keys: keys,
	}
	op := new(options.IndexOptions)
	for _, v := range opt {
		if v.Background != nil {
			op.Background = v.Background
		}
		if v.ExpireAfterSeconds != nil {
			op.ExpireAfterSeconds = v.ExpireAfterSeconds
		}

		if v.Name != nil {
			op.Name = v.Name
		}

		if v.Sparse != nil {
			op.Sparse = v.Sparse
		}
		if v.StorageEngine != nil {
			op.StorageEngine = v.StorageEngine
		}
		if v.Unique != nil {
			op.Unique = v.Unique
		}

		if v.Version != nil {
			op.Version = v.Version
		}

		if v.DefaultLanguage != nil {
			op.DefaultLanguage = v.DefaultLanguage
		}

		if v.LanguageOverride != nil {
			op.LanguageOverride = v.LanguageOverride
		}

		if v.TextVersion != nil {
			op.TextVersion = v.TextVersion
		}

		if v.Weights != nil {
			op.Weights = v.Weights
		}

		if v.SphereVersion != nil {
			op.SphereVersion = v.SphereVersion
		}

		if v.Bits != nil {
			op.Bits = v.Bits
		}

		if v.Max != nil {
			op.Max = v.Max
		}

		if v.Min != nil {
			op.Min = v.Min
		}

		if v.BucketSize != nil {
			op.BucketSize = v.BucketSize
		}

		if v.PartialFilterExpression != nil {
			op.PartialFilterExpression = v.PartialFilterExpression
		}

		if v.Collation != nil {
			op.Collation = v.Collation
		}

		if v.WildcardProjection != nil {
			op.WildcardProjection = v.WildcardProjection
		}

		if v.Hidden != nil {
			op.Hidden = v.Hidden
		}

	}
	m.Options = op
	i.indexes = append(i.indexes, m)
	return i
}

// SetMaxTime sets the value for the MaxTime field.
func (i *index) SetMaxTime(d time.Duration) Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetMaxTime(d))
	i.dropIndexesOptions = append(i.dropIndexesOptions, options.DropIndexes().SetMaxTime(d))
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
