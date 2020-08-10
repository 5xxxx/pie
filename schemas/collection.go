/*
 *
 * table.go
 * schemas
 *
 * Created by lintao on 2020/5/18 3:13 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package schemas

import "reflect"

type Collection struct {
	Name string
	Type reflect.Type
}

func NewEmptyCollection() *Collection {
	return NewCollection("", nil)
}

// NewCollection creates a new Collection object
func NewCollection(name string, t reflect.Type) *Collection {
	return &Collection{Name: name, Type: t}
}
