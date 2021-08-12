package driver

import (
	"errors"
	"reflect"
	"sync"

	"github.com/5xxxx/pie/names"
	"github.com/5xxxx/pie/schemas"
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
