/*
 *
 * index.go
 * pie
 *
 * Created by lintao on 2020/8/15 4:29 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package internal

import (
	"context"
	"time"

	"github.com/NSObjects/pie/driver"
	"github.com/NSObjects/pie/schemas"
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
	if len(opt) > 0 {
		m.Options = opt[0]
	}
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

func (c *index) SetDatabase(db string) driver.Indexes {
	c.db = db
	return c
}

func (c *index) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if c.doc != nil {
		coll, err = c.engine.CollectionNameForStruct(c.doc)
	} else {
		coll, err = c.engine.CollectionNameForStruct(doc)
	}
	if err != nil {
		return nil, err
	}
	return c.collectionByName(coll.Name), nil
}

func (c *index) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
	var coll *schemas.Collection
	var err error
	if c.doc != nil {
		coll, err = c.engine.CollectionNameForStruct(c.doc)
	} else {
		coll, err = c.engine.CollectionNameForSlice(doc)
	}
	if err != nil {
		return nil, err
	}
	return c.collectionByName(coll.Name), nil
}

func (c *index) collectionByName(name string) *mongo.Collection {
	return c.engine.Collection(name, c.db)
}

func (c *index) Collection(doc interface{}) driver.Indexes {
	c.doc = doc
	return c
}
