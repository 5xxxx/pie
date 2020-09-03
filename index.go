/*
 *
 * index.go
 * pie
 *
 * Created by lintao on 2020/8/15 4:29 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Indexes struct {
	collection
	indexes            []mongo.IndexModel
	createIndexOpts    []*options.CreateIndexesOptions
	dropIndexesOptions []*options.DropIndexesOptions
}

func NewIndexes(driver *Driver) *Indexes {
	return &Indexes{collection: collection{engine: driver}}
}

func (i *Indexes) CreateIndexes(ctx context.Context, doc interface{}) ([]string, error) {
	coll, err := i.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	return coll.Indexes().CreateMany(ctx, i.indexes, i.createIndexOpts...)
}

func (i *Indexes) DropAll(ctx context.Context, doc interface{}) error {
	coll, err := i.collectionForStruct(doc)
	if err != nil {
		return err
	}
	_, err = coll.Indexes().DropAll(ctx, i.dropIndexesOptions...)
	return err
}

func (i *Indexes) DropOne(ctx context.Context, doc interface{}, name string) error {
	coll, err := i.collectionForStruct(doc)
	if err != nil {
		return err
	}
	_, err = coll.Indexes().DropOne(ctx, name, i.dropIndexesOptions...)
	return err
}

func (i *Indexes) AddIndex(keys interface{}, opt ...*options.IndexOptions) *Indexes {
	m := mongo.IndexModel{
		Keys: keys,
	}
	if len(opt) > 0 {
		m.Options = opt[0]
	}
	i.indexes = append(i.indexes, m)
	return i
}

// SetMaxTime sets the value for the MaxTime field.
func (i *Indexes) SetMaxTime(d time.Duration) *Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetMaxTime(d))
	i.dropIndexesOptions = append(i.dropIndexesOptions, options.DropIndexes().SetMaxTime(d))
	return i
}

// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (i *Indexes) SetCommitQuorumInt(quorum int32) *Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumInt(quorum))
	return i
}

// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
func (i *Indexes) SetCommitQuorumString(quorum string) *Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumString(quorum))
	return i
}

// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (i *Indexes) SetCommitQuorumMajority() *Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumMajority())
	return i
}

// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (i *Indexes) SetCommitQuorumVotingMembers() *Indexes {
	i.createIndexOpts = append(i.createIndexOpts, options.CreateIndexes().SetCommitQuorumVotingMembers())
	return i
}
