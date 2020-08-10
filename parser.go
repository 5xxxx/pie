/*
 *
 * parser.go
 * tugrik
 *
 * Created by lintao on 2020/5/18 3:12 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package tugrik

import (
	"errors"
	"reflect"
	"sync"
	"tugrik/names"
	"tugrik/schemas"
)

var (
	ErrUnsupportedType = errors.New("Unsupported type")
)

type Parser struct {
	identifier       string
	collectionMapper names.Mapper
	collectionCache  sync.Map // map[reflect.Type]*schemas.Collection
	columnMapper     names.Mapper
}

func NewParser(collectionMapper, columnMapper names.Mapper) *Parser {
	return &Parser{
		collectionMapper: collectionMapper,
		columnMapper:     columnMapper,
	}
}

func (parser *Parser) Parse(v reflect.Value) (*schemas.Collection, error) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrUnsupportedType
	}

	collection := schemas.NewEmptyCollection()
	collection.Type = t
	collection.Name = names.GetCollectionName(parser.collectionMapper, v)
	return collection, nil
}
