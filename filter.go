package pie

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Operator query operator interface
type Operator interface {
	ToBSON() bson.D
}

// ComparisonOperator comparison operator
type ComparisonOperator struct {
	Field string
	Op    string
	Value any
}

func (op ComparisonOperator) ToBSON() bson.D {
	return bson.D{{op.Field, bson.D{{op.Op, op.Value}}}}
}

// LogicalOperator logical operator
type LogicalOperator struct {
	Op    string
	Value []Operator
}

func (op LogicalOperator) ToBSON() bson.D {
	var conditions []bson.D
	for _, cond := range op.Value {
		conditions = append(conditions, cond.ToBSON())
	}
	return bson.D{{op.Op, conditions}}
}

// FieldOperator field operator
type FieldOperator struct {
	Field string
	Op    string
	Value any
}

func (op FieldOperator) ToBSON() bson.D {
	return bson.D{{op.Field, bson.D{{op.Op, op.Value}}}}
}

// Predefined operator functions

// Eq equal
func Eq(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$eq", Value: value}
}

// Ne not equal
func Ne(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$ne", Value: value}
}

// Gt greater than
func Gt(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$gt", Value: value}
}

// Gte greater than or equal
func Gte(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$gte", Value: value}
}

// Lt less than
func Lt(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$lt", Value: value}
}

// Lte less than or equal
func Lte(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$lte", Value: value}
}

// In in range
func In(field string, values any) FieldOperator {
	return FieldOperator{Field: field, Op: "$in", Value: values}
}

// Nin not in range
func Nin(field string, values any) FieldOperator {
	return FieldOperator{Field: field, Op: "$nin", Value: values}
}

// Exists field exists
func Exists(field string, exists bool) FieldOperator {
	return FieldOperator{Field: field, Op: "$exists", Value: exists}
}

// Regex regular expression
func Regex(field string, pattern string) FieldOperator {
	return FieldOperator{Field: field, Op: "$regex", Value: pattern}
}

// RegexWithOptions regular expression with options
func RegexWithOptions(field string, pattern string, options string) FieldOperator {
	return FieldOperator{Field: field, Op: "$regex", Value: bson.D{
		{"$regex", pattern},
		{"$options", options},
	}}
}

// And logical and
func And(operators ...Operator) LogicalOperator {
	return LogicalOperator{Op: "$and", Value: operators}
}

// Or logical or
func Or(operators ...Operator) LogicalOperator {
	return LogicalOperator{Op: "$or", Value: operators}
}

// Nor logical nor
func Nor(operators ...Operator) LogicalOperator {
	return LogicalOperator{Op: "$nor", Value: operators}
}

// Not logical not
func Not(field string, value any) FieldOperator {
	return FieldOperator{Field: field, Op: "$not", Value: value}
}

// All array contains all elements
func All(field string, values any) FieldOperator {
	return FieldOperator{Field: field, Op: "$all", Value: values}
}

// ElemMatch array element match
func ElemMatch(field string, query bson.D) FieldOperator {
	return FieldOperator{Field: field, Op: "$elemMatch", Value: query}
}

// Size array size
func Size(field string, size int) FieldOperator {
	return FieldOperator{Field: field, Op: "$size", Value: size}
}

// Type field type
func Type(field string, bsonType string) FieldOperator {
	return FieldOperator{Field: field, Op: "$type", Value: bsonType}
}

// Mod modulo operation
func Mod(field string, divisor, remainder int) FieldOperator {
	return FieldOperator{Field: field, Op: "$mod", Value: []int{divisor, remainder}}
}

// Text text search
func Text(search string) FieldOperator {
	return FieldOperator{Field: "$text", Op: "$search", Value: search}
}

// Where use JavaScript expression
func Where(expression string) FieldOperator {
	return FieldOperator{Field: "$where", Op: "$expr", Value: expression}
}

// ID according to ID query
func ID(id any) FieldOperator {
	var objectID bson.ObjectID

	switch v := id.(type) {
	case string:
		var err error
		objectID, err = bson.ObjectIDFromHex(v)
		if err != nil {
			// If not a valid ObjectID, use the string directly
			return FieldOperator{Field: "_id", Op: "$eq", Value: v}
		}
	case bson.ObjectID:
		objectID = v
	default:
		return FieldOperator{Field: "_id", Op: "$eq", Value: id}
	}

	return FieldOperator{Field: "_id", Op: "$eq", Value: objectID}
}

// Between range query
func Between(field string, min, max any) LogicalOperator {
	return And(
		Gte(field, min),
		Lte(field, max),
	)
}

// Like fuzzy query (using regular expression)
func Like(field string, pattern string) FieldOperator {
	return Regex(field, pattern)
}

// ILike case-insensitive fuzzy query
func ILike(field string, pattern string) FieldOperator {
	return RegexWithOptions(field, pattern, "i")
}

// IsNull field is null
func IsNull(field string) FieldOperator {
	return FieldOperator{Field: field, Op: "$eq", Value: nil}
}

// IsNotNull field is not null
func IsNotNull(field string) FieldOperator {
	return FieldOperator{Field: field, Op: "$ne", Value: nil}
}

// IsEmpty field is empty string or empty array
func IsEmpty(field string) LogicalOperator {
	return Or(
		FieldOperator{Field: field, Op: "$eq", Value: ""},
		FieldOperator{Field: field, Op: "$size", Value: 0},
	)
}

// IsNotEmpty field is not empty
func IsNotEmpty(field string) LogicalOperator {
	return And(
		FieldOperator{Field: field, Op: "$ne", Value: ""},
		FieldOperator{Field: field, Op: "$ne", Value: nil},
	)
}
