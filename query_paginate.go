package pie

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PaginateParams pagination parameters
type PaginateParams struct {
	Page     int // Page number (starting from 1)
	PageSize int // Page size
}

// PaginateResult pagination result
type PaginateResult[T any] struct {
	Data       []T   // Current page data
	Total      int64 // Total record count
	Page       int   // Current page number
	PageSize   int   // Page size
	TotalPages int   // Total page count
	HasNext    bool  // Has next page
	HasPrev    bool  // Has previous page
}

// SimplePaginateResult simple pagination result (without total count)
type SimplePaginateResult[T any] struct {
	Data     []T  // Current page data
	Page     int  // Current page number
	PageSize int  // Page size
	HasNext  bool // Has next page
}

// CursorPaginateParams cursor pagination parameters
type CursorPaginateParams struct {
	Cursor     string   // Cursor (base64 encoded)
	PageSize   int      // Page size
	SortField  string   // Sort field
	SortFields []string // Multi-field sort (composite cursor)
}

// CursorPaginateResult cursor pagination result
type CursorPaginateResult[T any] struct {
	Data       []T    // Current page data
	NextCursor string // Next page cursor
	PrevCursor string // Previous page cursor
	HasNext    bool   // Has next page
	HasPrev    bool   // Has previous page
}

// IDCursorParams ID-based cursor pagination parameters
type IDCursorParams struct {
	AfterID  interface{} // Last ID of previous page
	PageSize int         // Page size
}

// Paginate offset-based pagination (with total count)
func (s *Session[T]) Paginate(ctx context.Context, params PaginateParams) (*PaginateResult[T], error) {
	// Set default values
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	// Calculate skip
	skip := int64((params.Page - 1) * params.PageSize)

	// First count total
	total, err := s.CountDocuments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	// Set pagination parameters
	s.query.Skip(skip).Limit(int64(params.PageSize))

	// Query data
	data, err := s.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	result := &PaginateResult[T]{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}

	return result, nil
}

// PaginateSimple offset-based pagination (without total count, faster)
func (s *Session[T]) PaginateSimple(ctx context.Context, params PaginateParams) (*SimplePaginateResult[T], error) {
	// Set default values
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	// Calculate skip
	skip := int64((params.Page - 1) * params.PageSize)

	// Query one more to determine if there's next page
	s.query.Skip(skip).Limit(int64(params.PageSize + 1))

	// Query data
	data, err := s.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Check if there's next page
	hasNext := len(data) > params.PageSize
	if hasNext {
		data = data[:params.PageSize] // Remove extra one
	}

	result := &SimplePaginateResult[T]{
		Data:     data,
		Page:     params.Page,
		PageSize: params.PageSize,
		HasNext:  hasNext,
	}

	return result, nil
}

// PaginateCursor cursor-based pagination
func (s *Session[T]) PaginateCursor(ctx context.Context, params CursorPaginateParams) (*CursorPaginateResult[T], error) {
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	// Use single field or multi-field sort
	sortFields := params.SortFields
	if len(sortFields) == 0 && params.SortField != "" {
		sortFields = []string{params.SortField}
	}
	if len(sortFields) == 0 {
		return nil, fmt.Errorf("sort field is required for cursor pagination")
	}

	// Parse cursor
	var cursorValue map[string]interface{}
	if params.Cursor != "" {
		decoded, err := base64.StdEncoding.DecodeString(params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		if err := json.Unmarshal(decoded, &cursorValue); err != nil {
			return nil, fmt.Errorf("invalid cursor format: %w", err)
		}

		// Add query conditions based on cursor value
		// Assume ascending sort, data after cursor should be greater than cursor value
		for _, field := range sortFields {
			if val, ok := cursorValue[field]; ok {
				s.query.Where(field, bson.D{{Key: "$gt", Value: val}})
			}
		}
	}

	// Query one more to determine if there's next page
	s.query.Limit(int64(params.PageSize + 1))

	// Query data
	data, err := s.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Check if there's next page
	hasNext := len(data) > params.PageSize
	if hasNext {
		data = data[:params.PageSize]
	}

	// Generate next page cursor
	var nextCursor string
	if hasNext && len(data) > 0 {
		lastDoc := data[len(data)-1]
		cursorMap := make(map[string]interface{})

		// Use reflection to get sort field values
		docValue := bson.M{}
		bsonBytes, _ := bson.Marshal(lastDoc)
		bson.Unmarshal(bsonBytes, &docValue)

		for _, field := range sortFields {
			if val, ok := docValue[field]; ok {
				cursorMap[field] = val
			}
		}

		cursorJSON, _ := json.Marshal(cursorMap)
		nextCursor = base64.StdEncoding.EncodeToString(cursorJSON)
	}

	result := &CursorPaginateResult[T]{
		Data:       data,
		NextCursor: nextCursor,
		PrevCursor: "", // Previous page cursor needs additional logic to implement
		HasNext:    hasNext,
		HasPrev:    params.Cursor != "",
	}

	return result, nil
}

// PaginateCursorByID ID-based cursor pagination (most simple)
func (s *Session[T]) PaginateCursorByID(ctx context.Context, params IDCursorParams) (*CursorPaginateResult[T], error) {
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	// If there's afterID, add query condition
	if params.AfterID != nil {
		s.query.Where("_id", bson.D{{Key: "$gt", Value: params.AfterID}})
	}

	// Sort by _id
	s.query.OrderBy("_id")

	// Query one more to determine if there's next page
	s.query.Limit(int64(params.PageSize + 1))

	// Query data
	data, err := s.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Check if there's next page
	hasNext := len(data) > params.PageSize
	if hasNext {
		data = data[:params.PageSize]
	}

	// Generate next page cursor (using last document's _id)
	var nextCursor string
	if hasNext && len(data) > 0 {
		lastDoc := data[len(data)-1]

		// Use reflection to get _id field value
		docValue := bson.M{}
		bsonBytes, _ := bson.Marshal(lastDoc)
		bson.Unmarshal(bsonBytes, &docValue)

		if id, ok := docValue["_id"]; ok {
			cursorMap := map[string]interface{}{"_id": id}
			cursorJSON, _ := json.Marshal(cursorMap)
			nextCursor = base64.StdEncoding.EncodeToString(cursorJSON)
		}
	}

	result := &CursorPaginateResult[T]{
		Data:       data,
		NextCursor: nextCursor,
		PrevCursor: "",
		HasNext:    hasNext,
		HasPrev:    params.AfterID != nil,
	}

	return result, nil
}
