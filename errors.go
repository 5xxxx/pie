package pie

import "errors"

// Define library error types
var (
	// Connection related errors
	ErrNotConnected = errors.New("pie: not connected to database")
	ErrInvalidURI   = errors.New("pie: invalid MongoDB URI")

	// Query related errors
	ErrInvalidQuery    = errors.New("pie: invalid query")
	ErrEmptyResult     = errors.New("pie: no documents found")
	ErrMultipleResults = errors.New("pie: multiple documents found")
	ErrInvalidDocument = errors.New("pie: invalid document")

	// Type related errors
	ErrInvalidType     = errors.New("pie: invalid type")
	ErrTypeMismatch    = errors.New("pie: type mismatch")
	ErrUnsupportedType = errors.New("pie: unsupported type")

	// Configuration related errors
	ErrInvalidConfig   = errors.New("pie: invalid configuration")
	ErrMissingRequired = errors.New("pie: missing required field")

	// Transaction related errors
	ErrTransactionFailed = errors.New("pie: transaction failed")
	ErrSessionRequired   = errors.New("pie: session required for transaction")

	// Cache related errors
	ErrCacheNotFound  = errors.New("pie: cache not found")
	ErrCacheExpired   = errors.New("pie: cache expired")
	ErrCacheDisabled  = errors.New("pie: cache disabled")
	ErrCacheOperation = errors.New("pie: cache operation failed")
)