package pie

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

type result struct {
	t        reflect.Type
	changed  bool
	hasIface bool
}

type TagMaker interface {
	// MakeTag makes tag for the field the fieldIndex in the structureType.
	// Result should depends on constant parameters of creation of the TagMaker and parameters
	// passed to the MakeTag. The MakeTag should not produce side effects (like a pure function).
	MakeTag(structureType reflect.Type, fieldIndex int) reflect.StructTag
}

type cacheKey struct {
	reflect.Type
	TagMaker
}

var cache = struct {
	sync.RWMutex
	m map[cacheKey]result
}{
	m: make(map[cacheKey]result),
}

func getType(structType reflect.Type, maker TagMaker, any bool) result {
	// TODO(yar): Improve synchronization for cases when one analogue
	// is produced concurently by different goroutines in the same time
	key := cacheKey{structType, maker}
	cache.RLock()
	res, ok := cache.m[key]
	cache.RUnlock()
	if !ok || (res.hasIface && !any) {
		res = makeType(structType, maker, any)
		cache.Lock()
		cache.m[key] = res
		cache.Unlock()
	}
	return res
}

func makeStructType(structType reflect.Type, maker TagMaker, any bool) result {
	if structType.NumField() == 0 {
		return result{t: structType, changed: false}
	}
	changed := false
	hasPrivate := false
	hasIface := false
	fields := make([]reflect.StructField, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		strField := structType.Field(i)
		oldType := strField.Type
		new := getType(oldType, maker, any)
		strField.Type = new.t
		if oldType != new.t {
			changed = true
		}
		if new.hasIface {
			hasIface = true
		}
		oldTag := strField.Tag
		newTag := maker.MakeTag(structType, i)
		strField.Tag = newTag
		if oldTag != newTag {
			changed = true
		}
		fields = append(fields, strField)
	}
	if !changed {
		return result{t: structType, changed: false, hasIface: hasIface}
	} else if hasPrivate {
		panic(fmt.Sprintf("unable to change tags for type %s, because it contains unexported fields", structType))
	}
	newType := reflect.StructOf(fields)
	compareStructTypes(structType, newType)
	return result{t: newType, changed: true, hasIface: hasIface}
}

func makeType(t reflect.Type, maker TagMaker, any bool) result {
	switch t.Kind() {
	case reflect.Struct:
		return makeStructType(t, maker, any)
	case reflect.Ptr:
		res := getType(t.Elem(), maker, any)
		if !res.changed {
			return result{t: t, changed: false}
		}
		return result{t: reflect.PtrTo(res.t), changed: true}
	case reflect.Array:
		res := getType(t.Elem(), maker, any)
		if !res.changed {
			return result{t: t, changed: false}
		}
		return result{t: reflect.ArrayOf(t.Len(), res.t), changed: true}
	case reflect.Slice:
		res := getType(t.Elem(), maker, any)
		if !res.changed {
			return result{t: t, changed: false}
		}
		return result{t: reflect.SliceOf(res.t), changed: true}
	case reflect.Map:
		resKey := getType(t.Key(), maker, any)
		resElem := getType(t.Elem(), maker, any)
		if !resKey.changed && !resElem.changed {
			return result{t: t, changed: false}
		}
		return result{t: reflect.MapOf(resKey.t, resElem.t), changed: true}
	case reflect.Interface:
		if any {
			return result{t: t, changed: false, hasIface: true}
		}
		fallthrough
	case
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer:
		panic("tags.Map: Unsupported type: " + t.Kind().String())
	default:
		// don't modify type in another case
		return result{t: t, changed: false}
	}
}

func compareStructTypes(source, result reflect.Type) {
	if source.Size() != result.Size() {
		// TODO: debug
		// fmt.Println(newType.Size(), newType)
		// for i := 0; i < newType.NumField(); i++ {
		// 	fmt.Println(newType.Field(i))
		// }
		// fmt.Println(structType.Size(), structType)
		// for i := 0; i < structType.NumField(); i++ {
		// 	fmt.Println(structType.Field(i))
		// }
		panic("tags.Map: Unexpected case - type has a size different from size of original type")
	}
}

//func insertOmitemptyTag(u driver{}) driver{} {
//	strPtrVal := reflect.ValueOf(u)
//	t := strPtrVal.Type().Elem()
//	num := t.NumField()
//	fields := make([]reflect.StructField, 0, t.NumField())
//	for index := 0; index < num; index++ {
//		field := t.Field(index)
//		tag := t.Field(index).Tag.Get("bson")
//		p := (*reflect.StructField)(unsafe.Pointer(&field))
//		tag = strings.Replace(tag, " ", "", -1)
//		p.Tag = reflect.StructTag(fmt.Sprintf(`bson:"%s,omitempty"`, tag))
//		fields = append(fields, field)
//	}
//
//	newType := reflect.StructOf(fields)
//	newPtrVal := reflect.NewAt(newType, unsafe.Pointer(strPtrVal.Pointer()))
//	return newPtrVal.Interface()
//}

func insertOmitemptyTag(u any) any {
	strPtrVal := reflect.ValueOf(u)
	res := getType(strPtrVal.Type().Elem(), maker{}, true)
	newPtrVal := reflect.NewAt(res.t, unsafe.Pointer(strPtrVal.Pointer()))
	return newPtrVal.Interface()
}

var structTypeConstructorBugWasFixed bool

func init() {
	switch {
	case strings.HasPrefix(runtime.Version(), "go1.7"):
		// there is bug in reflect.StructOf
	default:
		structTypeConstructorBugWasFixed = true
	}
}

type maker struct{}

func (m maker) MakeTag(t reflect.Type, fieldIndex int) reflect.StructTag {

	tag := t.Field(fieldIndex).Tag.Get("bson")
	if tag == "" {
		return ""
	}
	tag = strings.Replace(tag, " ", "", -1)
	if strings.Index(tag, "omitempty") > 0 {
		return reflect.StructTag(fmt.Sprintf(`bson:"%s"`, tag))
	}
	return reflect.StructTag(fmt.Sprintf(`bson:"%s,omitempty"`, tag))
}
