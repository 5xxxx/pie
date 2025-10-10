package pie

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NameMapper name mapper interface
type NameMapper interface {
	TableName(structName string) string
	FieldName(fieldName string) string
}

// SnakeMapper snake case name mapper
type SnakeMapper struct{}

func (m SnakeMapper) TableName(structName string) string {
	return toSnakeCase(structName)
}

func (m SnakeMapper) FieldName(fieldName string) string {
	return toSnakeCase(fieldName)
}

// CamelMapper camel case name mapper
type CamelMapper struct{}

func (m CamelMapper) TableName(structName string) string {
	return toCamelCase(structName)
}

func (m CamelMapper) FieldName(fieldName string) string {
	return toCamelCase(fieldName)
}

// SameMapper same name mapper
type SameMapper struct{}

func (m SameMapper) TableName(structName string) string {
	return structName
}

func (m SameMapper) FieldName(fieldName string) string {
	return fieldName
}

// TagParser tag parser
type TagParser struct {
	nameMapper NameMapper
}

// NewTagParser create tag parser
func NewTagParser(mapper NameMapper) *TagParser {
	return &TagParser{
		nameMapper: mapper,
	}
}

// ParseStruct parse struct information
func (p *TagParser) ParseStruct(v interface{}) (*CollectionInfo, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid type: expected struct")
	}

	rt := rv.Type()
	info := &CollectionInfo{
		Name:       p.nameMapper.TableName(rt.Name()),
		StructType: rt,
		Fields:     make([]FieldInfo, 0),
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldInfo := p.parseField(field)
		if fieldInfo != nil {
			info.Fields = append(info.Fields, *fieldInfo)

			// Check primary key
			if fieldInfo.IsPrimary {
				info.PrimaryKey = fieldInfo.BSONName
			}

			// Check soft delete field
			if fieldInfo.IsSoftDelete {
				info.SoftDelete = true
				info.DeletedAt = fieldInfo.BSONName
			}

			// Check time fields
			if fieldInfo.IsTimeField {
				if strings.Contains(fieldInfo.Name, "Created") {
					info.CreatedAt = fieldInfo.BSONName
				} else if strings.Contains(fieldInfo.Name, "Updated") {
					info.UpdatedAt = fieldInfo.BSONName
				}
			}
		}
	}

	return info, nil
}

// parseField parse single field
func (p *TagParser) parseField(field reflect.StructField) *FieldInfo {
	// Skip private fields
	if !field.IsExported() {
		return nil
	}

	tag := field.Tag.Get("bson")
	if tag == "" {
		tag = field.Tag.Get("pie")
	}

	// Parse tag
	bsonName := p.parseBSONTag(tag, field.Name)
	if bsonName == "-" {
		return nil // Skip this field
	}

	fieldInfo := &FieldInfo{
		Name:     field.Name,
		BSONName: bsonName,
		Type:     field.Type,
		Tag:      field.Tag,
	}

	// Check primary key
	if fieldInfo.BSONName == "_id" || strings.Contains(tag, "primary") {
		fieldInfo.IsPrimary = true
	}

	// Check index
	if strings.Contains(tag, "index") {
		fieldInfo.IsIndex = true
	}

	// Check unique index
	if strings.Contains(tag, "unique") {
		fieldInfo.IsUnique = true
	}

	// Check soft delete field
	if strings.Contains(tag, "soft_delete") || strings.Contains(field.Name, "DeletedAt") {
		fieldInfo.IsSoftDelete = true
	}

	// Check time fields
	if field.Type == reflect.TypeOf(bson.DateTime(0)) ||
		field.Type == reflect.TypeOf(bson.Timestamp{}) ||
		strings.Contains(field.Name, "At") {
		fieldInfo.IsTimeField = true
	}

	return fieldInfo
}

// parseBSONTag parse BSON tag
func (p *TagParser) parseBSONTag(tag, fieldName string) string {
	if tag == "" {
		return p.nameMapper.FieldName(fieldName)
	}

	// Handle "field_name,omitempty" format
	parts := strings.Split(tag, ",")
	if parts[0] == "" {
		return p.nameMapper.FieldName(fieldName)
	}

	return parts[0]
}

// CollectionInfo collection detailed information
type CollectionInfo struct {
	Name       string
	Database   string
	StructType reflect.Type
	Fields     []FieldInfo
	PrimaryKey string
	SoftDelete bool
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}

// FieldInfo field information
type FieldInfo struct {
	Name         string
	BSONName     string
	Type         reflect.Type
	Tag          reflect.StructTag
	IsPrimary    bool
	IsIndex      bool
	IsUnique     bool
	IsSoftDelete bool
	IsTimeField  bool
}

// Helper functions

// toSnakeCase convert to snake case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// toCamelCase convert to camel case
func toCamelCase(str string) string {
	parts := strings.Split(str, "_")
	if len(parts) == 1 {
		return str
	}

	result := parts[0]
	for _, part := range parts[1:] {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return result
}
