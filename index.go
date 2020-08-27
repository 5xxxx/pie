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
	db 				  string
	driver             *Driver
	indexes            []mongo.IndexModel
	createIndexOpts    []*options.CreateIndexesOptions
	dropIndexesOptions []*options.DropIndexesOptions
}

func NewIndexes(driver *Driver) *Indexes {
	return &Indexes{driver: driver}
}

func (s *Indexes) CreateIndexes(ctx context.Context, doc interface{}) ([]string, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	return coll.Indexes().CreateMany(ctx, s.indexes, s.createIndexOpts...)
}

func (s *Indexes) DropAll(ctx context.Context, doc interface{}) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}
	_, err = coll.Indexes().DropAll(ctx, s.dropIndexesOptions...)
	return err
}

func (s *Indexes) DropOne(ctx context.Context, doc interface{}, name string) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}
	_, err = coll.Indexes().DropOne(ctx, name, s.dropIndexesOptions...)
	return err
}

func (s *Indexes) AddIndex(keys interface{}, opt ...*options.IndexOptions) *Indexes {
	m := mongo.IndexModel{
		Keys: keys,
	}
	if len(opt) > 0 {
		m.Options = opt[0]
	}
	s.indexes = append(s.indexes, m)
	return s
}

// SetMaxTime sets the value for the MaxTime field.
func (c *Indexes) SetMaxTime(d time.Duration) *Indexes {
	c.createIndexOpts = append(c.createIndexOpts, options.CreateIndexes().SetMaxTime(d))
	c.dropIndexesOptions = append(c.dropIndexesOptions, options.DropIndexes().SetMaxTime(d))
	return c
}

// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (c *Indexes) SetCommitQuorumInt(quorum int32) *Indexes {
	c.createIndexOpts = append(c.createIndexOpts, options.CreateIndexes().SetCommitQuorumInt(quorum))
	return c
}

// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
func (c *Indexes) SetCommitQuorumString(quorum string) *Indexes {
	c.createIndexOpts = append(c.createIndexOpts, options.CreateIndexes().SetCommitQuorumString(quorum))
	return c
}

// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (c *Indexes) SetCommitQuorumMajority() *Indexes {
	c.createIndexOpts = append(c.createIndexOpts, options.CreateIndexes().SetCommitQuorumMajority())
	return c
}

// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (c *Indexes) SetCommitQuorumVotingMembers() *Indexes {
	c.createIndexOpts = append(c.createIndexOpts, options.CreateIndexes().SetCommitQuorumVotingMembers())
	return c
}



func (s *Indexes) SetDatabase(db string) *Indexes {
	s.db = db
	return s
}

func (e *Indexes) collectionForStruct(doc interface{}) (*mongo.Collection,error) {
	coll, err := e.driver.CollectionNameForStruct(doc)
	if err != nil {
		return nil, err
	}
	return e.collection(coll.Name),nil
}

func (e *Indexes) collectionForSlice(doc interface{}) (*mongo.Collection,error) {
	coll, err := e.driver.CollectionNameForSlice(doc)
	if err != nil {
		return nil, err
	}
	return e.collection(coll.Name),nil
}

func (s *Indexes) collection(name string)*mongo.Collection {
	var db string
	if s.db != "" {
		db = s.db
	} else {
		db = s.driver.db
	}
	return s.driver.client.Database(db).Collection(name)
}