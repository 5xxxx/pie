package pie

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Query query builder
type Query struct {
	filter  bson.D
	sort    bson.D
	limit   *int64
	skip    *int64
	project bson.D
}

// NewQuery create new query builder
func NewQuery() *Query {
	return &Query{
		filter:  bson.D{},
		sort:    bson.D{},
		project: bson.D{},
	}
}

// Where add condition
func (q *Query) Where(field string, value any) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: value})
	return q
}

// WhereOperator add condition with operator
func (q *Query) WhereOperator(op Operator) *Query {
	bsonDoc := op.ToBSON()
	if len(bsonDoc) > 0 {
		q.filter = append(q.filter, bsonDoc[0])
	}
	return q
}

// And add AND condition
func (q *Query) And(operators ...Operator) *Query {
	if len(operators) == 0 {
		return q
	}

	var conditions []bson.D
	for _, op := range operators {
		conditions = append(conditions, op.ToBSON())
	}

	q.filter = append(q.filter, bson.E{Key: "$and", Value: conditions})
	return q
}

// Or add OR condition
func (q *Query) Or(operators ...Operator) *Query {
	if len(operators) == 0 {
		return q
	}

	var conditions []bson.D
	for _, op := range operators {
		conditions = append(conditions, op.ToBSON())
	}

	q.filter = append(q.filter, bson.E{Key: "$or", Value: conditions})
	return q
}

// Nor add NOR condition
func (q *Query) Nor(operators ...Operator) *Query {
	if len(operators) == 0 {
		return q
	}

	var conditions []bson.D
	for _, op := range operators {
		conditions = append(conditions, op.ToBSON())
	}

	q.filter = append(q.filter, bson.E{Key: "$nor", Value: conditions})
	return q
}

// OrderBy add sort
func (q *Query) OrderBy(field string) *Query {
	q.sort = append(q.sort, bson.E{Key: field, Value: 1})
	return q
}

// OrderByDesc add descending sort
func (q *Query) OrderByDesc(field string) *Query {
	q.sort = append(q.sort, bson.E{Key: field, Value: -1})
	return q
}

// Sort add sort (support multiple fields)
func (q *Query) Sort(sort bson.D) *Query {
	q.sort = append(q.sort, sort...)
	return q
}

// Limit set limit count
func (q *Query) Limit(limit int64) *Query {
	q.limit = &limit
	return q
}

// Skip set skip count
func (q *Query) Skip(skip int64) *Query {
	q.skip = &skip
	return q
}

// Project set projection fields
func (q *Query) Project(project bson.D) *Query {
	q.project = append(q.project, project...)
	return q
}

// Select select fields (alias method)
func (q *Query) Select(fields ...string) *Query {
	project := bson.D{}
	for _, field := range fields {
		project = append(project, bson.E{Key: field, Value: 1})
	}
	return q.Project(project)
}

// Exclude exclude fields
func (q *Query) Exclude(fields ...string) *Query {
	project := bson.D{}
	for _, field := range fields {
		project = append(project, bson.E{Key: field, Value: 0})
	}
	return q.Project(project)
}

// GetFilter get filter conditions
func (q *Query) GetFilter() bson.D {
	return q.filter
}

// GetSort get sort conditions
func (q *Query) GetSort() bson.D {
	return q.sort
}

// GetLimit get limit count
func (q *Query) GetLimit() *int64 {
	return q.limit
}

// GetSkip get skip count
func (q *Query) GetSkip() *int64 {
	return q.skip
}

// GetProject get projection conditions
func (q *Query) GetProject() bson.D {
	return q.project
}

// Clone clone query
func (q *Query) Clone() *Query {
	newQuery := &Query{
		filter:  make(bson.D, len(q.filter)),
		sort:    make(bson.D, len(q.sort)),
		project: make(bson.D, len(q.project)),
	}

	copy(newQuery.filter, q.filter)
	copy(newQuery.sort, q.sort)
	copy(newQuery.project, q.project)

	if q.limit != nil {
		limit := *q.limit
		newQuery.limit = &limit
	}

	if q.skip != nil {
		skip := *q.skip
		newQuery.skip = &skip
	}

	return newQuery
}

// Clear clear query conditions
func (q *Query) Clear() *Query {
	q.filter = bson.D{}
	q.sort = bson.D{}
	q.project = bson.D{}
	q.limit = nil
	q.skip = nil
	return q
}

// BuildFindOptions build Find options
func (q *Query) BuildFindOptions() *options.FindOptionsBuilder {
	opts := options.Find()

	if len(q.sort) > 0 {
		opts.SetSort(q.sort)
	}

	if q.limit != nil {
		opts.SetLimit(*q.limit)
	}

	if q.skip != nil {
		opts.SetSkip(*q.skip)
	}

	if len(q.project) > 0 {
		opts.SetProjection(q.project)
	}

	return opts
}

// BuildFindOneOptions build FindOne options
func (q *Query) BuildFindOneOptions() *options.FindOneOptionsBuilder {
	opts := options.FindOne()

	if len(q.sort) > 0 {
		opts.SetSort(q.sort)
	}

	if q.skip != nil {
		opts.SetSkip(*q.skip)
	}

	if len(q.project) > 0 {
		opts.SetProjection(q.project)
	}

	return opts
}

// Build build query filter
func (q *Query) Build() bson.D {
	return q.filter
}

// WhereArrayAll array contains all specified elements
func (q *Query) WhereArrayAll(field string, values any) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$all", Value: values}}})
	return q
}

// WhereRecentDays recent N days
func (q *Query) WhereRecentDays(field string, days int) *Query {
	startTime := time.Now().AddDate(0, 0, -days)
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$gte", Value: startTime}}})
	return q
}
