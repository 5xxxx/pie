/*
 *
 * collection_name.go
 * names
 *
 * Created by lintao on 2020/8/8 4:20 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package names

import (
	"reflect"
	"sync"
)

// CollectionName collection name driver to define customerize collection name
type CollectionName interface {
	CollectionName() string
}

var (
	tpCollectionName = reflect.TypeOf((*CollectionName)(nil)).Elem()
	tvCache          sync.Map
)

func GetCollectionName(mapper Mapper, v reflect.Value) string {
	if v.Type().Implements(tpCollectionName) {
		return v.Interface().(CollectionName).CollectionName()
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Type().Implements(tpCollectionName) {
			return v.Interface().(CollectionName).CollectionName()
		}
	} else if v.CanAddr() {
		v1 := v.Addr()
		if v1.Type().Implements(tpCollectionName) {
			return v1.Interface().(CollectionName).CollectionName()
		}
	} else {
		name, ok := tvCache.Load(v.Type())
		if ok {
			if name.(string) != "" {
				return name.(string)
			}
		} else {
			v2 := reflect.New(v.Type())
			if v2.Type().Implements(tpCollectionName) {
				tableName := v2.Interface().(CollectionName).CollectionName()
				tvCache.Store(v.Type(), tableName)
				return tableName
			}

			tvCache.Store(v.Type(), "")
		}
	}

	return mapper.Obj2Collection(v.Type().Name())
}
