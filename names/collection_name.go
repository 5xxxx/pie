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

// CollectionName interface to define custom collection name
type CollectionName interface {
	CollectionName() string
}

var (
	tpCollectionName = reflect.TypeOf((*CollectionName)(nil)).Elem()
	tvCache          sync.Map
)

func implementsCollectionName(v reflect.Value) (string, bool) {
	if v.Type().Implements(tpCollectionName) {
		return v.Interface().(CollectionName).CollectionName(), true
	}

	return "", false
}

func checkAndStoreInCache(v reflect.Value) string {
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

	return ""
}

func GetCollectionName(mapper Mapper, v reflect.Value) string {
	if result, ok := implementsCollectionName(v); ok {
		return result
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if result, ok := implementsCollectionName(v); ok {
			return result
		}
	} else if v.CanAddr() {
		v1 := v.Addr()
		if result, ok := implementsCollectionName(v1); ok {
			return result
		}
	}

	name := checkAndStoreInCache(v)
	if name != "" {
		return name
	}

	return mapper.Obj2Collection(v.Type().Name())
}
