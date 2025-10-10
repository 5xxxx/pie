package pie

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// WhereIn checks if field value is in specified range
func (s *Session[T]) WhereIn(field string, values interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$in", Value: values}}})
	return s
}

// WhereNotIn checks if field value is not in specified range
func (s *Session[T]) WhereNotIn(field string, values interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$nin", Value: values}}})
	return s
}

// WhereBetween checks if field value is between min and max
func (s *Session[T]) WhereBetween(field string, min, max interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$gte", Value: min},
			{Key: "$lte", Value: max},
		},
	})
	return s
}

// WhereNotBetween checks if field value is not between min and max
func (s *Session[T]) WhereNotBetween(field string, min, max interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$or",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$lt", Value: min}}}},
			{{Key: field, Value: bson.D{{Key: "$gt", Value: max}}}},
		},
	})
	return s
}

// WhereNull checks if field value is null or does not exist
func (s *Session[T]) WhereNull(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$or",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: false}}}},
			{{Key: field, Value: nil}},
		},
	})
	return s
}

// WhereNotNull checks if field value is not null and exists
func (s *Session[T]) WhereNotNull(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$and",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: true}}}},
			{{Key: field, Value: bson.D{{Key: "$ne", Value: nil}}}},
		},
	})
	return s
}

// WhereExists adds a filter ensuring the given field exists in the document.
func (s *Session[T]) WhereExists(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: true}}})
	return s
}

// WhereNotExists adds a filter requiring the given field to be missing.
func (s *Session[T]) WhereNotExists(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: false}}})
	return s
}

// WhereDate performs date comparisons on the specified field. When a string is
// supplied it is parsed using several common layouts before falling back to the
// original value.
func (s *Session[T]) WhereDate(field string, operator string, value interface{}) *Session[T] {
	var op string
	switch operator {
	case ">", "gt":
		op = "$gt"
	case ">=", "gte":
		op = "$gte"
	case "<", "lt":
		op = "$lt"
	case "<=", "lte":
		op = "$lte"
	case "=", "==", "eq":
		op = "$eq"
	default:
		op = "$eq"
	}

	// Convert string values into time.Time when possible so comparisons are
	// executed using native date operators.
	var dateValue time.Time
	switch v := value.(type) {
	case string:
		// Attempt multiple layouts to support common timestamp formats
		// without requiring callers to pre-parse values.
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			time.RFC3339,
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				dateValue = t
				break
			}
		}
	case time.Time:
		dateValue = v
	default:
		// Fall back to storing the raw value when conversion is not
		// possible. This allows callers to provide pre-parsed values or
		// alternative types.
		s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: op, Value: value}}})
		return s
	}

	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: op, Value: dateValue}}})
	return s
}

// WhereDateBetween constrains a date field to lie within the provided range.
func (s *Session[T]) WhereDateBetween(field string, start, end interface{}) *Session[T] {
	return s.WhereBetween(field, start, end)
}

// WhereMonth filters documents whose date field resolves to the given month
// number, leveraging MongoDB aggregation expressions.
func (s *Session[T]) WhereMonth(field string, month int) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$expr",
		Value: bson.D{{
			Key: "$eq",
			Value: []interface{}{
				bson.D{{Key: "$month", Value: "$" + field}},
				month,
			},
		}},
	})
	return s
}

// WhereYear filters documents whose date field resolves to the provided year.
func (s *Session[T]) WhereYear(field string, year int) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$expr",
		Value: bson.D{{
			Key: "$eq",
			Value: []interface{}{
				bson.D{{Key: "$year", Value: "$" + field}},
				year,
			},
		}},
	})
	return s
}

// WhereRecentDays keeps documents whose date field is within the last N days.
func (s *Session[T]) WhereRecentDays(field string, days int) *Session[T] {
	startTime := time.Now().AddDate(0, 0, -days)
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$gte", Value: startTime}}})
	return s
}

// WhereLike builds a case-insensitive regular-expression match derived from a
// SQL-like pattern string (using % wildcards).
func (s *Session[T]) WhereLike(field string, pattern string) *Session[T] {
	// Strip SQL-like wildcard markers and convert them into regular
	// expression anchors so the resulting filter matches the caller's
	// expectations.
	regexPattern := pattern
	if len(pattern) > 0 {
		if pattern[0] == '%' && pattern[len(pattern)-1] == '%' {
			// %keyword% -> keyword (contains)
			regexPattern = pattern[1 : len(pattern)-1]
		} else if pattern[0] == '%' {
			// %keyword -> keyword$ (suffix match)
			regexPattern = pattern[1:] + "$"
		} else if pattern[len(pattern)-1] == '%' {
			// keyword% -> ^keyword (prefix match)
			regexPattern = "^" + pattern[:len(pattern)-1]
		}
	}

	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$regex", Value: regexPattern},
			{Key: "$options", Value: "i"}, // Case-insensitive matching
		},
	})
	return s
}

// WhereStartsWith adds a case-insensitive prefix match using a regular
// expression anchored to the beginning of the field value.
func (s *Session[T]) WhereStartsWith(field string, prefix string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$regex", Value: "^" + prefix},
			{Key: "$options", Value: "i"},
		},
	})
	return s
}

// WhereEndsWith appends a case-insensitive suffix match to the query using a
// regular expression anchored to the end of the field value.
func (s *Session[T]) WhereEndsWith(field string, suffix string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$regex", Value: suffix + "$"},
			{Key: "$options", Value: "i"},
		},
	})
	return s
}

// WhereArrayContains ensures that the provided element is present in the array
// field.
func (s *Session[T]) WhereArrayContains(field string, value interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: value})
	return s
}

// WhereArraySize enforces an exact array length using the $size operator.
func (s *Session[T]) WhereArraySize(field string, size int) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$size", Value: size}}})
	return s
}

// WhereArrayAll requires an array field to contain all elements provided by
// values.
func (s *Session[T]) WhereArrayAll(field string, values interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$all", Value: values}}})
	return s
}

// OrWhere allows callers to construct OR clauses using a callback that receives
// a transient Query instance.
func (s *Session[T]) OrWhere(callback func(*Query) *Query) *Session[T] {
	// Create a dedicated query builder for the OR clause so the callback can
	// add isolated conditions without mutating the parent builder directly.
	orQuery := NewQuery()
	callback(orQuery)

	// Retrieve the filter built by the callback and append it as an $or
	// clause when it contains predicates.
	orFilter := orQuery.GetFilter()
	if len(orFilter) > 0 {
		s.query.filter = append(s.query.filter, bson.E{
			Key:   "$or",
			Value: []bson.D{orFilter},
		})
	}

	return s
}

// AndWhere mirrors OrWhere but combines the generated conditions with $and.
func (s *Session[T]) AndWhere(callback func(*Query) *Query) *Session[T] {
	// Build the temporary query for AND conditions.
	andQuery := NewQuery()
	callback(andQuery)

	// If the callback contributed predicates, wrap them inside an $and
	// clause before attaching to the parent query.
	andFilter := andQuery.GetFilter()
	if len(andFilter) > 0 {
		s.query.filter = append(s.query.filter, bson.E{
			Key:   "$and",
			Value: []bson.D{andFilter},
		})
	}

	return s
}

// Query extensions (allowing chainable Query usage)

// WhereIn appends an $in condition for Query builders.
func (q *Query) WhereIn(field string, values interface{}) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$in", Value: values}}})
	return q
}

// WhereNotIn appends an $nin condition for Query builders.
func (q *Query) WhereNotIn(field string, values interface{}) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$nin", Value: values}}})
	return q
}

// WhereBetween adds a range constraint by combining $gte and $lte operators.
func (q *Query) WhereBetween(field string, min, max interface{}) *Query {
	q.filter = append(q.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$gte", Value: min},
			{Key: "$lte", Value: max},
		},
	})
	return q
}

// WhereNull builds a filter that matches when a field is missing or set to
// null.
func (q *Query) WhereNull(field string) *Query {
	q.filter = append(q.filter, bson.E{
		Key: "$or",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: false}}}},
			{{Key: field, Value: nil}},
		},
	})
	return q
}

// WhereNotNull ensures the field exists and its value is not null.
func (q *Query) WhereNotNull(field string) *Query {
	q.filter = append(q.filter, bson.E{
		Key: "$and",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: true}}}},
			{{Key: field, Value: bson.D{{Key: "$ne", Value: nil}}}},
		},
	})
	return q
}

// WhereExists is the Query variant of the session helper that requires field
// existence.
func (q *Query) WhereExists(field string) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: true}}})
	return q
}

// WhereNotExists ensures the field is absent in the matched documents.
func (q *Query) WhereNotExists(field string) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: false}}})
	return q
}

// Bool creates a pointer to the provided boolean literal.
func Bool(b bool) *bool {
	return &b
}

// Float64 creates a pointer to the provided float64 literal.
func Float64(f float64) *float64 {
	return &f
}

// Int creates a pointer to the provided int literal.
func Int(i int) *int {
	return &i
}

// String creates a pointer to the provided string literal.
func String(s string) *string {
	return &s
}
