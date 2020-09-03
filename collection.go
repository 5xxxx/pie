package pie

import (
	"github.com/NSObjects/pie/schemas"
	"go.mongodb.org/mongo-driver/mongo"
)

type collection struct {
	db     string
	doc    interface{}
	engine *Driver
}

func (c *collection) SetDatabase(db string) *collection {
	c.db = db
	return c
}

func (c *collection) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
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

func (c *collection) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
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

func (c *collection) collectionByName(name string) *mongo.Collection {
	var db string
	if c.db != "" {
		db = c.db
	} else {
		db = c.engine.db
	}
	return c.engine.client.Database(db).Collection(name)
}

func (c *collection) Collection(doc interface{}) *collection {
	c.doc = doc
	return c
}
