package internal

import (
	"context"
	"time"

	"github.com/5xxxx/pie/driver"
	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type index struct {
	db                 string
	doc                interface{}
	engine             driver.Client
	indexes            []mongo.IndexModel
	createIndexOpts    []*options.CreateIndexesOptions
	dropIndexesOptions []*options.DropIndexesOptions
}

func NewIndexes(driver driver.Client) driver.Indexes {
	return &index{engine: driver}
}

func (i *index) CreateIndexes(doc interface{}, ctx ...context.Context) ([]string, error) {
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

func (i *index) DropAll(doc interface{}, ctx ...context.Context) error {
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

func (i *index) DropOne(doc interface{}, name string, ctx ...context.Context) error {
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

func (i *index) AddIndex(keys interface{}, opt ...*options.IndexOptions) driver.Indexes {
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
func (i *index) SetMaxTime(d time.Duration) driver.Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetMaxTime(d))
	i.dropIndexesOptions = append(i.dropIndexesOptions, options.DropIndexes().SetMaxTime(d))
	return i
}

// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (i *index) SetCommitQuorumInt(quorum int32) driver.Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumInt(quorum))
	return i
}

// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
func (i *index) SetCommitQuorumString(quorum string) driver.Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumString(quorum))
	return i
}

// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (i *index) SetCommitQuorumMajority() driver.Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumMajority())
	return i
}

// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (i *index) SetCommitQuorumVotingMembers() driver.Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumVotingMembers())
	return i
}

func (i *index) SetDatabase(db string) driver.Indexes {
	i.db = db
	return i
}

func (i *index) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
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

func (i *index) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
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
	return i.engine.Collection(name, i.db)
}

func (i *index) Collection(doc interface{}) driver.Indexes {
	i.doc = doc
	return i
}
