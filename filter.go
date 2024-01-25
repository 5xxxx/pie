package pie

import (
	"errors"
	"fmt"
	"github.com/5xxxx/pie/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type Condition interface {
	RegexFilter(key, pattern string) Condition
	ID(id interface{}) Condition
	// Eq Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value interface{}) Condition

	// Gt {field: {$gt: value} } >
	Gt(key string, gt interface{}) Condition

	// Gte { qty: { $gte: 20 } } >=
	Gte(key string, gte interface{}) Condition

	// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in interface{}) Condition

	// Lt {field: {$lt: value} } <
	Lt(key string, lt interface{}) Condition

	// Lte { field: { $lte: value} } <=
	Lte(key string, lte interface{}) Condition

	// Ne {field: {$ne: value} } !=
	Ne(key string, ne interface{}) Condition

	// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin interface{}) Condition

	// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(filter Condition) Condition

	// Not { field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not interface{}) Condition

	// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(filter Condition) Condition
	// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(filter Condition) Condition

	Exists(key string, exists bool, filter ...Condition) Condition

	// Type { field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t interface{}) Condition

	// Expr Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(filter Condition) Condition

	// Regex todo 简单实现，后续增加支持
	Regex(key string, value string) Condition

	Filters() (bson.D, error)
	A() bson.A
	Err() error
	FilterBson(d interface{}) Condition
	FilterBy(object interface{}) Condition

	Clone() Condition
}

// filter is a type used to build query filters for MongoDB.
// It contains a bson.D field to store the filter conditions and an err field to handle any errors that may occur during filter creation.
// Clone creates a copy of the filter.
// It returns a new filter with the same conditions and error value as the original.
type filter struct {
	d   bson.D
	err error
}

// Clone creates a new instance of filter with the same values.
// The new filter is a shallow copy of the original filter, sharing the same slice and error.
// Changes made to the original filter will not affect the cloned filter.
// Returns the cloned filter.
//
// Example usage:
//
//	filter := &filter{d: []bson.E{{Key: "key1", Value: "value1"}}}
//	clonedFilter := filter.Clone()
//	// clonedFilter is a new filter instance with the same values as the original filter
//	// filter.d and clonedFilter.d now point to the same underlying slice.
//	clonedFilter.d = append(clonedFilter.d, bson.E{Key: "key2", Value: 123})
//	// filter.d remains unchanged
func (f *filter) Clone() Condition {
	return &filter{
		d:   f.d,
		err: f.err,
	}
}

// FilterBson accepts either a bson.M or bson.D object and adds its key-value pairs to the filter condition.
// If the object is of type bson.M, each key-value pair is added as a bson.E element in the filter.
// If the object is of type bson.D, the entire bson.D object is appended to the filter.
// Returns the updated filter condition.
//
// Example usage:
//
//	filter := &filter{}
//	obj := bson.M{"key1": "value1", "key2": 123}
//	filter.FilterBson(obj)
//	// filter.d is now []bson.E{bson.E{Key: "key1", Value: "value1"}, bson.E{Key: "key2", Value: 123}}
//
//	obj2 := bson.D{{"key3", "value3"}, {"key4", true}}
//	filter.FilterBson(obj2)
//	// filter.d is now []bson.E{bson.E{Key: "key1", Value: "value1"}, bson.E{Key: "key2", Value: 123}, bson.E{Key: "key3", Value: "value3"}, bson.E{Key: "key4", Value: true}}
//
// Note: The filter is maintained internally by the `filter` object to build a MongoDB query condition.
func (f *filter) FilterBson(object interface{}) Condition {
	if m, ok := object.(bson.M); ok {
		for key, value := range m {
			f.d = append(f.d, bson.E{Key: key, Value: value})
		}
	} else if d, ok := object.(bson.D); ok {
		f.d = append(f.d, d...)
	}
	return f
}

// FilterBy applies filtering based on the fields of the provided object
func (f *filter) FilterBy(object interface{}) Condition {
	beanValue := reflect.ValueOf(object)
	if beanValue.Kind() != reflect.Struct {
		f.err = errors.New("needs a struct")
		return f
	}
	f.constructFilterFromFields(beanValue)
	return f
}

func (f *filter) constructFilterFromFields(beanValue reflect.Value) {
	docType := reflect.TypeOf(beanValue.Interface())
	fieldCount := docType.NumField()
	for idx := 0; idx < fieldCount; idx++ {
		field := docType.Field(idx)
		if f.shouldMakeFilterValue(field) {
			split := strings.Split(field.Tag.Get("bson"), ",")
			f.makeFilterValue(split[0], beanValue.Field(idx).Interface())
		}
	}
}

func (f *filter) shouldMakeFilterValue(field reflect.StructField) bool {
	fieldTag := field.Tag.Get("bson")
	return fieldTag != "" && fieldTag != "-"
}

// Rest of the code...

func (f *filter) Err() error {
	return f.err
}

// DefaultCondition returns a default condition initialized with an empty bson.D slice. This condition can be used to chain query filter methods and construct complex filter expressions
func DefaultCondition() Condition {
	return &filter{d: bson.D{}}
}

// Filters returns the filter conditions as a bson.D and any error that occurred during filtering.
func (f *filter) Filters() (bson.D, error) {
	return f.d, f.err
}

// RegexFilter applies a regular expression filter to a specified key
// and pattern. The "i" option is used to make the regular expression case-insensitive.
// Example:
// { name: /^J/i }
// Field that Starts with a Case-Insensitive Letter "J"
// { 'name.first': /^J/i }
// Field in an Embedded Document Starts with a Case-Insensitive Letter "J"
// { 'contact.phone': /^1\d/i }
// Address Field Exists and Starts with a Case-Insensitive Digit "1"
// { 'contact.phone': /^1\d{10}$/ }
// Address Field Exists and is a String of 11 Digits that Starts with a Case-Insensitive Digit "1"
// { 'contact.email': /^a-zA-Z\d+@[a-zA-Z\d]+\.[a-zA-Z]{2,}$/ }
// Address Field Exists and is a String that Contains only ASCII Characters.
func (f *filter) RegexFilter(key, pattern string) Condition {
	v := primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

func (f *filter) ID(id interface{}) Condition {
	if id == nil {
		return f
	}
	switch id.(type) {
	case string:
		f.processStringID(id.(string))
	case primitive.ObjectID:
		f.processObjectID(id.(primitive.ObjectID))
	default:
		f.err = f.generateError("id type must be string or primitive.ObjectID", id)
	}
	return f
}

func (f *filter) processStringID(id string) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		f.err = f.generateError("id can't parse", id, err)
		return
	}
	f.processObjectID(objectId)
}

func (f *filter) processObjectID(id primitive.ObjectID) {
	if id == primitive.NilObjectID {
		f.err = f.generateError("id can't be nil", id)
		return
	}
	f.addFilter("_id", id)
}

func (f *filter) generateError(format string, a ...interface{}) error {
	return fmt.Errorf(format+" %v", append(a, f.err)...)
}

func (f *filter) addFilter(key string, value interface{}) {
	f.d = append(f.d, bson.E{Key: key, Value: value})
}

// Eq adds an equality expression to the filter condition.
// It appends a bson.E{Key: key, Value: value} to the filter's d field.
// The returned Condition is the filter itself.
// Example usage: f.Eq("field", value) adds an equality expression "field: value" to the filter.
func (f *filter) Eq(key string, value interface{}) Condition {
	f.d = append(f.d, bson.E{Key: key, Value: value})
	return f
}

// Gt specifies that the value of the field must be greater than the specified value.
// Example usage:
// { qty: { $gt: 20 } }
// Field value is greater than 20.
//
// { "item.name": { $gt: "ab" } }
// Field value in an embedded document is greater than "ab".
//
// { tags: { $gt: [ "A", "B" ] } }
// Field value is greater than the array ["A", "B"].
func (f *filter) Gt(key string, gt interface{}) Condition {
	v := bson.M{
		"$gt": gt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Gte specifies that the value of a field must be greater than or equal to the specified value.
//
// Example usage:
//
//	filter.Gte("qty", 20)
//	filter.Gte("item.name", "ab")
//	filter.Gte("tags", []string{"A", "B"})
//
// This will append a condition to the filter that matches documents where the value of the specified field is greater than or equal to the given value.
// If the field is nested in an embedded document, you can specify the path using dot notation (e.g., "item.name").
// If the field contains an array, you can specify an array of values and the condition will match if any of the values in the array are greater than or equal to the specified value
func (f *filter) Gte(key string, gte interface{}) Condition {
	v := bson.M{
		"$gte": gte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// In Matches any of the specified values.
// { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// Matches any value in the specified array.
// { field: { $in: [ <value1>, <value2> ... <valueN> ], $nin: [ <value1>, <value2> ... <valueN> ] } }
func (f *filter) In(key string, in interface{}) Condition {
	v := bson.M{
		"$in": in,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Lt Represents a Less Than condition
// { qty: { $lt: 20 } }
// Field in Embedded Document Less Than a Value
// {"item.price": { $lt: 9.99 } }
// Less Than an Array Value
// { tags: { $lt: [ "C", "D" ] } }
func (f *filter) Lt(key string, lt interface{}) Condition {
	v := bson.M{
		"$lt": lt,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Lte Less than or Equal to a Specified Value
// {qty: {$lte: 20}}
// Field in Embedded Document Less than or Equal to a Value
// {"item.name": {$lte: "ab"}}
// Less than or Equal to an Array Value
// {tags: {$lte: ["A", "B"]}}
func (f *filter) Lte(key string, lte interface{}) Condition {
	v := bson.M{
		"$lte": lte,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Ne Not Equals a Specified Value
// { qty: { $ne: 20 } }
// Field in Embedded Document Not Equals a Value
// { "item.name": { $ne: "ab" } }
// Not Equals an Array Value
// { tags: { $ne: [ "A", "B" ] } }
func (f *filter) Ne(key string, ne interface{}) Condition {
	v := bson.M{
		"$ne": ne,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Nin Field Does Not Match any Value in the Specified Array
// Example:
//
//	`{"qty": bson.M{"$nin": [5, 15]}}`
//	Field qty does not match any value in the array [5, 15]
//	`{"item.name": bson.M{"$nin": ["ab", "cd"]}}`
//	Field item.name does not match any value in the array ["ab", "cd"]
//	`{"tags": bson.M{"$nin": ["A", "B", "C"]}}`
//	Field tags does not match any value in the array ["A", "B", "C"]
func (f *filter) Nin(key string, nin interface{}) Condition {
	v := bson.M{
		"$nin": nin,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// And accepts a Condition object and adds it to the filter condition using the "$and" operator.
// The filter condition is updated by appending a bson.E element with key "$and" and value filter.A().
// Returns the updated filter condition.
func (f *filter) And(filter Condition) Condition {
	f.d = append(f.d, bson.E{Key: "$and", Value: filter.A()})
	return f
}

// Not adds a condition to the filter where the specified key's value is not equal to the given value.
// The key is used as the field to apply the not condition to.
// The not interface{} parameter represents the value that the key's value is not equal to.
// It creates a bson.M map with the "$not" operator and the given value, and appends it as a bson.E element in the filter.
func (f *filter) Not(key string, not interface{}) Condition {
	v := bson.M{
		"$not": not,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

func (f *filter) Nor(filter Condition) Condition {
	f.d = append(f.d, bson.E{Key: "$nor", Value: filter.A()})
	return f
}

// Or performs a logical OR operation on a filter condition.
// This method appends the filter condition to the existing filter
// using the $or operator.
//
// Example:
//
//	{
//	   $or: [
//	       { qty: 20 },
//	       { "item.name": "ab" },
//	       { tags: [ "A", "B" ] }
//	   ]
//	}
//
// Parameters:
// - filter: The filter condition to be logically OR'ed with.
//
// Returns:
// The updated filter condition after the logical OR operation.
func (f *filter) Or(filter Condition) Condition {
	f.d = append(f.d, bson.E{Key: "$or", Value: filter.A()})
	return f
}

// Exists specifies whether a field exists or not in a document.
// The `key` parameter specifies the field name to be checked.
// The `exists` parameter specifies whether the field should exist (true) or not (false).
// Additional filters can be passed as variadic arguments to further refine the query.
// The returned result is a filter condition.
//
// Example:
// To check if the "qty" field exists and has a value, use the following code:
//
//	filter := &filter{}
//	condition := filter.Exists("qty", true)
//
// To check if the "item.name" field exists and has a value, use the following code:
//
//	filter := &filter{}
//	condition := filter.Exists("item.name", true)
//
// To check if the "tags" field exists and has an array value, use the following code:
//
//	filter := &filter{}
//	condition := filter.Exists("tags", true)
//
// Multiple conditions can be chained together by passing them as variadic arguments:
//
//	filter := &filter{}
//	condition := filter.Exists("qty", true, filter.Eq("status", "active"))
//
// If an error occurs while processing the provided filters, the function sets
// the error in the filter itself and returns the filter with the error:
//
//	filter := &filter{}
//	condition := filter.Exists("qty", true, filter.Eq("status", "active"), filter.Lt("price", 10))
//
// The returned condition can be used in query operations.
// The filter can be applied later to filter the query results.
//
//	query := collection.Find(condition)
//	result, err := query.All()
//	if err != nil {
//	    // Handle error
//	}
//	// Process the query result
//
// Note: When using the `exists` parameter, if it is set to false and the field exists but has a value of null or an array is empty, it will still match the condition as true.
func (f *filter) Exists(key string, exists bool, filter ...Condition) Condition {
	m := bson.M{
		"$exists": exists,
	}

	for _, v := range filter {
		filters, err := v.Filters()
		if err != nil {
			f.err = err
			return f
		}
		for _, fv := range filters {
			m[fv.Key] = fv.Value
		}
	}

	f.d = append(f.d, bson.E{Key: key, Value: m})
	return f
}

// Type Specifies the BSON type for a specific field
// Example:
//
//	f.Type("age", bson.TypeInt32) // field "age" must be of type int32
//	f.Type("name", bson.TypeString) // field "name" must be of type string
func (f *filter) Type(key string, t interface{}) Condition {
	v := bson.M{
		"$type": t,
	}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

// Expr appends a bson.E element with the key "$expr" and the value of the filter's A() method to the filter condition.
// The filter's A() method is expected to return a Condition object.
// Returns the updated filter condition.
// Example usage:
//
//	filter := &filter{}
//	cond := bson.M{"$gt": bson.A{"$field1", "$field2"}}
//	filter.Expr(cond)
//	// filter.d is now []bson.E{bson.E{Key: "$expr", Value: bson.M{"$gt": bson.A{"$field1", "$field2"}}}}
func (f *filter) Expr(filter Condition) Condition {
	f.d = append(f.d, bson.E{Key: "$expr", Value: filter.A()})
	return f
}

// Regex adds a regular expression condition to the filter.
//
// The key parameter is the name of the field to match against.
//
// The value parameter is the regular expression pattern to match.
//
// Returns the updated filter condition.
//
// Example usage:
//
//	filter := &filter{}
//	filter.Regex("name", "^J") // Matches names that start with "J"
//	// filter.d is now []bson.E{bson.E{Key: "name", Value: primitive.Regex{Pattern: "^J", Options: ""}}}
func (f *filter) Regex(key string, value string) Condition {
	v := primitive.Regex{Pattern: value, Options: ""}
	f.d = append(f.d, bson.E{Key: key, Value: v})
	return f
}

func (f *filter) A() bson.A {
	var fs bson.A

	for _, value := range f.d {
		fs = append(fs, value)
	}
	return fs
}

// makeFilterValue checks the kind of the value and calls the appropriate function to handle it.
// If the value is a struct, it calls makeStructValue function to handle it.
// If the value is an array, it returns without doing anything.
// Otherwise, it calls the Eq function to add the field-value pair to the filter.
// Extract function refactoring applied
func (f *filter) makeFilterValue(field string, value interface{}) {
	if utils.IsZero(value) {
		return
	}
	v := reflect.ValueOf(value)
	f.processFilterValueKind(field, v, value)
}

// New extracted function
func (f *filter) processFilterValueKind(field string, v reflect.Value, value interface{}) {
	switch v.Kind() {
	case reflect.Struct:
		f.makeStructValue(field, v)
	case reflect.Array:
		// Skipping for array kind of dimensions
		return
	default:
		// Put the default action in the default clause of the switch statement
		f.Eq(field, value)
	}
}

// makeStructValue iterates over the fields of a struct value and processes each field with a "bson" struct tag.
// It takes a `field` string, which is the parent field or parent prefix for the current field, and a `value` of type `reflect.Value` which is the struct value.
//
// The function loops over each field using `value.NumField()` method and gets the type and tag information for each field.
// If the field has a non-empty "bson" tag, the function checks if the field value is non-zero using `utils.IsZero` function.
// If the field value is non-zero, a new field tag is created by appending the current field tag to the parent field, separated by a dot.
// The function then calls `f.makeFilterValue` with the new field tag and the field value.
func (f *filter) makeStructValue(field string, value reflect.Value) {
	for index := 0; index < value.NumField(); index++ {
		docType := reflect.TypeOf(value.Interface())
		tag := docType.Field(index).Tag.Get("bson")
		if tag != "" {
			if !utils.IsZero(value.Field(index)) {
				fieldTags := fmt.Sprintf("%s.%s", field, tag)
				f.makeFilterValue(fieldTags, value.Field(index).Interface())
			}
		}
	}
}
