package pie

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// WhereStruct generate query conditions based on struct
func (s *Session[T]) WhereStruct(filter any) *Session[T] {
	conditions := parseStructToConditions(filter)
	for _, cond := range conditions {
		s.query.filter = append(s.query.filter, cond)
	}
	return s
}

// parseStructToConditions parse struct to query conditions
func parseStructToConditions(filter any) []bson.E {
	var conditions []bson.E

	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return conditions
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get pie tag
		pieTag := field.Tag.Get("pie")
		if pieTag == "" {
			continue
		}

		// Skip "-" tag
		if pieTag == "-" {
			continue
		}

		// Parse tag
		parts := strings.Split(pieTag, ",")
		if len(parts) == 0 {
			continue
		}

		fieldName := parts[0]
		operators := parts[1:]

		// Check omitempty option
		hasOmitEmpty := false
		for _, op := range operators {
			if op == "omitempty" {
				hasOmitEmpty = true
				break
			}
		}

		// If field value is zero value and omitempty is set, skip
		if hasOmitEmpty && isZeroValue(value) {
			continue
		}

		// Generate conditions based on operators
		if len(operators) > 0 {
			operator := ""
			for _, op := range operators {
				if op == "" || op == "omitempty" {
					continue
				}
				operator = op
				break
			}

			if operator != "" {
				condition := buildCondition(fieldName, operator, value)
				if condition.Key != "" {
					conditions = append(conditions, condition)
				}
				continue
			}
		}

		// No operator, use exact match by default
		if !isZeroValue(value) {
			conditions = append(conditions, bson.E{Key: fieldName, Value: value.Interface()})
		}
	}

	return conditions
}

// buildCondition build query condition based on operator
func buildCondition(field string, operator string, value reflect.Value) bson.E {
	// 对于null和notNull操作符，总是生成条件，不管值是什么
	switch operator {
	case "null":
		// Field is null - always generate condition
		return bson.E{Key: "$or", Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: false}}}},
			{{Key: field, Value: nil}},
		}}

	case "notNull":
		// Field is not null - always generate condition
		return bson.E{Key: "$and", Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: true}}}},
			{{Key: field, Value: bson.D{{Key: "$ne", Value: nil}}}},
		}}
	}

	// 检查value是否有效
	if !value.IsValid() {
		return bson.E{}
	}

	val := value.Interface()

	switch operator {
	case "eq":
		// Skip zero value for exact match
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: val}

	case "ne":
		// Skip zero value for not equal
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$ne", Value: val}}}

	case "gt":
		// Skip zero value for greater than
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$gt", Value: val}}}

	case "gte":
		// Skip zero value for greater than or equal
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$gte", Value: val}}}

	case "lt":
		// Skip zero value for less than
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$lt", Value: val}}}

	case "lte":
		// Skip zero value for less than or equal
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$lte", Value: val}}}

	case "in":
		// Skip zero value for in
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$in", Value: val}}}

	case "nin":
		// Skip zero value for not in
		if isZeroValue(value) {
			return bson.E{}
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$nin", Value: val}}}

	case "like":
		// Fuzzy query
		if strVal, ok := val.(string); ok {
			return bson.E{Key: field, Value: bson.D{
				{Key: "$regex", Value: strVal},
				{Key: "$options", Value: "i"},
			}}
		}

	case "ilike":
		// Case-insensitive fuzzy query
		if strVal, ok := val.(string); ok {
			return bson.E{Key: field, Value: bson.D{
				{Key: "$regex", Value: strVal},
				{Key: "$options", Value: "i"},
			}}
		}

	case "prefix":
		// Prefix match
		if strVal, ok := val.(string); ok {
			return bson.E{Key: field, Value: bson.D{
				{Key: "$regex", Value: "^" + strVal},
				{Key: "$options", Value: "i"},
			}}
		}

	case "suffix":
		// Suffix match
		if strVal, ok := val.(string); ok {
			return bson.E{Key: field, Value: bson.D{
				{Key: "$regex", Value: strVal + "$"},
				{Key: "$options", Value: "i"},
			}}
		}

	case "regex":
		// Regular expression
		if strVal, ok := val.(string); ok {
			return bson.E{Key: field, Value: bson.D{{Key: "$regex", Value: strVal}}}
		}

	case "contains":
		// Array contains
		return bson.E{Key: field, Value: val}

	case "all":
		// Array contains all
		return bson.E{Key: field, Value: bson.D{{Key: "$all", Value: val}}}

	case "size":
		// Array size
		return bson.E{Key: field, Value: bson.D{{Key: "$size", Value: val}}}

	case "elemMatch":
		// Array element match
		return bson.E{Key: field, Value: bson.D{{Key: "$elemMatch", Value: val}}}

	case "exists":
		// Field exists - always generate condition regardless of value
		exists := true
		if boolVal, ok := val.(bool); ok {
			exists = boolVal
		}
		return bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: exists}}}

	case "between":
		// Range query (requires array type)
		if value.Kind() == reflect.Slice && value.Len() == 2 {
			return bson.E{Key: field, Value: bson.D{
				{Key: "$gte", Value: value.Index(0).Interface()},
				{Key: "$lte", Value: value.Index(1).Interface()},
			}}
		}

	default:
		// Default exact match
		if !isZeroValue(value) {
			return bson.E{Key: field, Value: val}
		}
	}

	return bson.E{}
}

// isZeroValue check if value is zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	case reflect.Struct:
		// Special handling for time.Time
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		// Other structs are considered non-zero values
		return false
	default:
		return false
	}
}

// Query method extensions

// WhereStruct generate query conditions based on struct
func (q *Query) WhereStruct(filter any) *Query {
	conditions := parseStructToConditions(filter)
	for _, cond := range conditions {
		q.filter = append(q.filter, cond)
	}
	return q
}

// StructFilter struct filter interface (optional)
type StructFilter interface {
	ToFilter() bson.D
}

// WhereStructFilter use struct that implements StructFilter interface
func (s *Session[T]) WhereStructFilter(filter StructFilter) *Session[T] {
	conditions := filter.ToFilter()
	s.query.filter = append(s.query.filter, conditions...)
	return s
}

// validateStructFilter validate struct filter
func validateStructFilter(filter any) error {
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("filter must be a struct, got %s", v.Kind())
	}

	return nil
}

// ParseStructToBSON parse struct to bson.D (exported for external use)
func ParseStructToBSON(filter any) (bson.D, error) {
	if err := validateStructFilter(filter); err != nil {
		return nil, err
	}

	conditions := parseStructToConditions(filter)
	result := make(bson.D, 0, len(conditions))
	result = append(result, conditions...)

	return result, nil
}
