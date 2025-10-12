package pie

import (
	"math/rand/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ScopeFunc scope function type
type ScopeFunc func(*Query) *Query

// Scopes apply multiple query scopes
func (s *Session[T]) Scopes(scopes ...ScopeFunc) *Session[T] {
	for _, scope := range scopes {
		s.query = scope(s.query)
	}
	return s
}

// Latest get latest N records
func (s *Session[T]) Latest(field string, limit int) *Session[T] {
	s.query.OrderByDesc(field).Limit(int64(limit))
	return s
}

// Oldest get earliest N records
func (s *Session[T]) Oldest(field string, limit int) *Session[T] {
	s.query.OrderBy(field).Limit(int64(limit))
	return s
}

// Random get N random records
func (s *Session[T]) Random(limit int) *Session[T] {
	// Use $sample aggregation stage to implement random sampling
	// Note: This needs to use aggregation pipeline in actual execution
	// Here we mark it first, handle it specially in execution
	s.query.Limit(int64(limit))

	// Add random sort mark
	s.query.Sort(bson.D{{Key: "_random", Value: rand.Float64()}})
	return s
}

// Common predefined scopes

// ActiveScope only query active status records
func ActiveScope(statusField string) ScopeFunc {
	return func(q *Query) *Query {
		return q.Where(statusField, "active")
	}
}

// NotDeletedScope only query non-deleted records
func NotDeletedScope(deletedAtField string) ScopeFunc {
	return func(q *Query) *Query {
		return q.WhereNull(deletedAtField)
	}
}

// PublishedScope only query published records
func PublishedScope(publishedField string) ScopeFunc {
	return func(q *Query) *Query {
		return q.Where(publishedField, true)
	}
}

// VisibleScope only query visible records
func VisibleScope(visibilityField string, visibleValues ...string) ScopeFunc {
	return func(q *Query) *Query {
		if len(visibleValues) == 0 {
			visibleValues = []string{"public"}
		}
		return q.WhereIn(visibilityField, visibleValues)
	}
}

// DateRangeScope date range scope
func DateRangeScope(field string, start, end any) ScopeFunc {
	return func(q *Query) *Query {
		return q.WhereBetween(field, start, end)
	}
}

// SearchScope search scope (multi-field fuzzy match)
func SearchScope(keyword string, fields ...string) ScopeFunc {
	return func(q *Query) *Query {
		if keyword == "" || len(fields) == 0 {
			return q
		}

		// Build $or conditions
		var orConditions []bson.D
		for _, field := range fields {
			orConditions = append(orConditions, bson.D{
				{Key: field, Value: bson.D{
					{Key: "$regex", Value: keyword},
					{Key: "$options", Value: "i"},
				}},
			})
		}

		q.filter = append(q.filter, bson.E{Key: "$or", Value: orConditions})
		return q
	}
}

// PriceRangeScope price range scope
func PriceRangeScope(minPrice, maxPrice *float64) ScopeFunc {
	return func(q *Query) *Query {
		if minPrice != nil {
			q.filter = append(q.filter, bson.E{Key: "price", Value: bson.D{{Key: "$gte", Value: *minPrice}}})
		}
		if maxPrice != nil {
			q.filter = append(q.filter, bson.E{Key: "price", Value: bson.D{{Key: "$lte", Value: *maxPrice}}})
		}
		return q
	}
}

// CategoryScope category scope
func CategoryScope(categories ...string) ScopeFunc {
	return func(q *Query) *Query {
		if len(categories) == 0 {
			return q
		}
		if len(categories) == 1 {
			return q.Where("category", categories[0])
		}
		return q.WhereIn("category", categories)
	}
}

// TagsScope tag scope (contains any tag)
func TagsScope(tags ...string) ScopeFunc {
	return func(q *Query) *Query {
		if len(tags) == 0 {
			return q
		}
		return q.WhereIn("tags", tags)
	}
}

// AllTagsScope tag scope (contains all tags)
func AllTagsScope(tags ...string) ScopeFunc {
	return func(q *Query) *Query {
		if len(tags) == 0 {
			return q
		}
		return q.WhereArrayAll("tags", tags)
	}
}

// UserScope user scope
func UserScope(userID string) ScopeFunc {
	return func(q *Query) *Query {
		return q.Where("user_id", userID)
	}
}

// TenantScope tenant scope
func TenantScope(tenantID string) ScopeFunc {
	return func(q *Query) *Query {
		return q.Where("tenant_id", tenantID)
	}
}

// FeaturedScope featured scope
func FeaturedScope() ScopeFunc {
	return func(q *Query) *Query {
		return q.Where("featured", true).OrderByDesc("priority")
	}
}

// PopularScope popular scope (sort by view count)
func PopularScope(viewCountField string, limit int) ScopeFunc {
	return func(q *Query) *Query {
		return q.OrderByDesc(viewCountField).Limit(int64(limit))
	}
}

// RecentScope recently created records
func RecentScope(createdAtField string, days int) ScopeFunc {
	return func(q *Query) *Query {
		return q.WhereRecentDays(createdAtField, days)
	}
}

// VerifiedScope verified records
func VerifiedScope(verifiedField string) ScopeFunc {
	return func(q *Query) *Query {
		return q.Where(verifiedField, true)
	}
}

// ApprovedScope approved records
func ApprovedScope(approvalStatus string) ScopeFunc {
	return func(q *Query) *Query {
		if approvalStatus == "" {
			approvalStatus = "approved"
		}
		return q.Where("approval_status", approvalStatus)
	}
}
