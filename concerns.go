package pie

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ConcernManager read/write concern manager
type ConcernManager struct {
	engine       *Engine
	readConcern  *ReadConcernManager
	writeConcern *WriteConcernManager
	readPref     *ReadPreferenceManager
}

// ReadConcernManager read concern manager
type ReadConcernManager struct {
	engine      *Engine
	collections map[string]*mongo.Collection
}

// WriteConcernManager write concern manager
type WriteConcernManager struct {
	engine      *Engine
	collections map[string]*mongo.Collection
}

// ReadPreferenceManager read preference manager
type ReadPreferenceManager struct {
	engine      *Engine
	collections map[string]*mongo.Collection
}

// NewConcernManager create read/write concern manager
func NewConcernManager(engine *Engine) *ConcernManager {
	return &ConcernManager{
		engine:       engine,
		readConcern:  NewReadConcernManager(engine),
		writeConcern: NewWriteConcernManager(engine),
		readPref:     NewReadPreferenceManager(engine),
	}
}

// NewReadConcernManager create read concern manager
func NewReadConcernManager(engine *Engine) *ReadConcernManager {
	return &ReadConcernManager{
		engine:      engine,
		collections: make(map[string]*mongo.Collection),
	}
}

// NewWriteConcernManager create write concern manager
func NewWriteConcernManager(engine *Engine) *WriteConcernManager {
	return &WriteConcernManager{
		engine:      engine,
		collections: make(map[string]*mongo.Collection),
	}
}

// NewReadPreferenceManager create read preference manager
func NewReadPreferenceManager(engine *Engine) *ReadPreferenceManager {
	return &ReadPreferenceManager{
		engine:      engine,
		collections: make(map[string]*mongo.Collection),
	}
}

// ReadConcernManager methods

// SetReadConcern set collection's read concern
func (rcm *ReadConcernManager) SetReadConcern(collectionName string, level string) error {
	collection := rcm.engine.Collection(collectionName)

	// Note: ReadConcern is typically set at collection level in v2
	// This method is provided for API compatibility
	rcm.collections[collectionName] = collection
	return nil
}

// GetReadConcern get collection's read concern
func (rcm *ReadConcernManager) GetReadConcern(collectionName string) (string, error) {
	// Note: ReadConcern retrieval is not directly supported in v2
	// This method is provided for API compatibility
	return "local", nil
}

// SetReadConcernForStruct set read concern for struct
func (rcm *ReadConcernManager) SetReadConcernForStruct(v any, level string) error {
	collection, err := rcm.engine.CollectionForStruct(v)
	if err != nil {
		return fmt.Errorf("failed to get collection for struct: %w", err)
	}

	rcm.collections[collection.Name()] = collection
	return nil
}

// SetGlobalReadConcern set global read concern
func (rcm *ReadConcernManager) SetGlobalReadConcern(level string) error {
	// Note: Global ReadConcern is typically set at client level in v2
	// This method is provided for API compatibility
	return nil
}

// GetGlobalReadConcern get global read concern
func (rcm *ReadConcernManager) GetGlobalReadConcern() (string, error) {
	// Note: Global ReadConcern retrieval is not directly supported in v2
	// This method is provided for API compatibility
	return "local", nil
}

// WriteConcernManager methods

// SetWriteConcern set collection's write concern
func (wcm *WriteConcernManager) SetWriteConcern(collectionName string, level string) error {
	collection := wcm.engine.Collection(collectionName)

	// Note: WriteConcern is typically set at collection level in v2
	// This method is provided for API compatibility
	wcm.collections[collectionName] = collection
	return nil
}

// GetWriteConcern get collection's write concern
func (wcm *WriteConcernManager) GetWriteConcern(collectionName string) (string, error) {
	// Note: WriteConcern retrieval is not directly supported in v2
	// This method is provided for API compatibility
	return "majority", nil
}

// SetWriteConcernForStruct set write concern for struct
func (wcm *WriteConcernManager) SetWriteConcernForStruct(v any, level string) error {
	collection, err := wcm.engine.CollectionForStruct(v)
	if err != nil {
		return fmt.Errorf("failed to get collection for struct: %w", err)
	}

	wcm.collections[collection.Name()] = collection
	return nil
}

// SetGlobalWriteConcern set global write concern
func (wcm *WriteConcernManager) SetGlobalWriteConcern(level string) error {
	// Note: Global WriteConcern is typically set at client level in v2
	// This method is provided for API compatibility
	return nil
}

// GetGlobalWriteConcern get global write concern
func (wcm *WriteConcernManager) GetGlobalWriteConcern() (string, error) {
	// Note: Global WriteConcern retrieval is not directly supported in v2
	// This method is provided for API compatibility
	return "majority", nil
}

// ReadPreferenceManager methods

// SetReadPreference set collection's read preference
func (rpm *ReadPreferenceManager) SetReadPreference(collectionName string, mode string) error {
	collection := rpm.engine.Collection(collectionName)

	// Note: ReadPreference is typically set at collection level in v2
	// This method is provided for API compatibility
	rpm.collections[collectionName] = collection
	return nil
}

// GetReadPreference get collection's read preference
func (rpm *ReadPreferenceManager) GetReadPreference(collectionName string) (string, error) {
	// Note: ReadPreference retrieval is not directly supported in v2
	// This method is provided for API compatibility
	return "primary", nil
}

// SetReadPreferenceForStruct set read preference for struct
func (rpm *ReadPreferenceManager) SetReadPreferenceForStruct(v any, mode string) error {
	collection, err := rpm.engine.CollectionForStruct(v)
	if err != nil {
		return fmt.Errorf("failed to get collection for struct: %w", err)
	}

	rpm.collections[collection.Name()] = collection
	return nil
}

// SetGlobalReadPreference set global read preference
func (rpm *ReadPreferenceManager) SetGlobalReadPreference(mode string) error {
	// Note: Global ReadPreference is typically set at client level in v2
	// This method is provided for API compatibility
	return nil
}

// GetGlobalReadPreference get global read preference
func (rpm *ReadPreferenceManager) GetGlobalReadPreference() (string, error) {
	// Note: Global ReadPreference retrieval is not directly supported in v2
	// This method is provided for API compatibility
	return "primary", nil
}

// ConcernManager methods

// ReadConcern get read concern manager
func (cm *ConcernManager) ReadConcern() *ReadConcernManager {
	return cm.readConcern
}

// WriteConcern get write concern manager
func (cm *ConcernManager) WriteConcern() *WriteConcernManager {
	return cm.writeConcern
}

// ReadPreference get read preference manager
func (cm *ConcernManager) ReadPreference() *ReadPreferenceManager {
	return cm.readPref
}

// SetCollectionConcerns set all concerns for collection
func (cm *ConcernManager) SetCollectionConcerns(collectionName string, readConcern, writeConcern, readPreference string) error {
	if err := cm.readConcern.SetReadConcern(collectionName, readConcern); err != nil {
		return fmt.Errorf("failed to set read concern: %w", err)
	}

	if err := cm.writeConcern.SetWriteConcern(collectionName, writeConcern); err != nil {
		return fmt.Errorf("failed to set write concern: %w", err)
	}

	if err := cm.readPref.SetReadPreference(collectionName, readPreference); err != nil {
		return fmt.Errorf("failed to set read preference: %w", err)
	}

	return nil
}

// SetStructConcerns set all concerns for struct
func (cm *ConcernManager) SetStructConcerns(v any, readConcern, writeConcern, readPreference string) error {
	if err := cm.readConcern.SetReadConcernForStruct(v, readConcern); err != nil {
		return fmt.Errorf("failed to set read concern: %w", err)
	}

	if err := cm.writeConcern.SetWriteConcernForStruct(v, writeConcern); err != nil {
		return fmt.Errorf("failed to set write concern: %w", err)
	}

	if err := cm.readPref.SetReadPreferenceForStruct(v, readPreference); err != nil {
		return fmt.Errorf("failed to set read preference: %w", err)
	}

	return nil
}

// SetGlobalConcerns set global concerns
func (cm *ConcernManager) SetGlobalConcerns(readConcern, writeConcern, readPreference string) error {
	if err := cm.readConcern.SetGlobalReadConcern(readConcern); err != nil {
		return fmt.Errorf("failed to set global read concern: %w", err)
	}

	if err := cm.writeConcern.SetGlobalWriteConcern(writeConcern); err != nil {
		return fmt.Errorf("failed to set global write concern: %w", err)
	}

	if err := cm.readPref.SetGlobalReadPreference(readPreference); err != nil {
		return fmt.Errorf("failed to set global read preference: %w", err)
	}

	return nil
}

// GetCollectionConcerns get all concerns for collection
func (cm *ConcernManager) GetCollectionConcerns(collectionName string) (readConcern, writeConcern, readPreference string, err error) {
	readConcern, err = cm.readConcern.GetReadConcern(collectionName)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get read concern: %w", err)
	}

	writeConcern, err = cm.writeConcern.GetWriteConcern(collectionName)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get write concern: %w", err)
	}

	readPreference, err = cm.readPref.GetReadPreference(collectionName)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get read preference: %w", err)
	}

	return readConcern, writeConcern, readPreference, nil
}

// GetGlobalConcerns get global concerns
func (cm *ConcernManager) GetGlobalConcerns() (readConcern, writeConcern, readPreference string, err error) {
	readConcern, err = cm.readConcern.GetGlobalReadConcern()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get global read concern: %w", err)
	}

	writeConcern, err = cm.writeConcern.GetGlobalWriteConcern()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get global write concern: %w", err)
	}

	readPreference, err = cm.readPref.GetGlobalReadPreference()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get global read preference: %w", err)
	}

	return readConcern, writeConcern, readPreference, nil
}

// ConcernOptions concern options
type ConcernOptions struct {
	ReadConcern    string
	WriteConcern   string
	ReadPreference string
}

// NewConcernOptions create concern options
func NewConcernOptions() *ConcernOptions {
	return &ConcernOptions{
		ReadConcern:    "local",
		WriteConcern:   "majority",
		ReadPreference: "primary",
	}
}

// WithReadConcern set read concern
func (co *ConcernOptions) WithReadConcern(level string) *ConcernOptions {
	co.ReadConcern = level
	return co
}

// WithWriteConcern set write concern
func (co *ConcernOptions) WithWriteConcern(level string) *ConcernOptions {
	co.WriteConcern = level
	return co
}

// WithReadPreference set read preference
func (co *ConcernOptions) WithReadPreference(mode string) *ConcernOptions {
	co.ReadPreference = mode
	return co
}

// Build build concern options
func (co *ConcernOptions) Build() *ConcernOptions {
	return co
}

// ConcernLevel concern level constants
const (
	// ReadConcern levels
	ReadConcernLocal        = "local"
	ReadConcernMajority     = "majority"
	ReadConcernLinearizable = "linearizable"
	ReadConcernSnapshot     = "snapshot"
	ReadConcernAvailable    = "available"

	// WriteConcern levels
	WriteConcernUnacknowledged = "unacknowledged"
	WriteConcernAcknowledged   = "acknowledged"
	WriteConcernMajority       = "majority"
	WriteConcernJournaled      = "journaled"

	// ReadPreference modes
	ReadPreferencePrimary            = "primary"
	ReadPreferencePrimaryPreferred   = "primaryPreferred"
	ReadPreferenceSecondary          = "secondary"
	ReadPreferenceSecondaryPreferred = "secondaryPreferred"
	ReadPreferenceNearest            = "nearest"
)

// ConcernBuilder concern builder
type ConcernBuilder struct {
	manager *ConcernManager
	options *ConcernOptions
}

// NewConcernBuilder create concern builder
func NewConcernBuilder(manager *ConcernManager) *ConcernBuilder {
	return &ConcernBuilder{
		manager: manager,
		options: NewConcernOptions(),
	}
}

// ReadConcern set read concern
func (cb *ConcernBuilder) ReadConcern(level string) *ConcernBuilder {
	cb.options.ReadConcern = level
	return cb
}

// WriteConcern set write concern
func (cb *ConcernBuilder) WriteConcern(level string) *ConcernBuilder {
	cb.options.WriteConcern = level
	return cb
}

// ReadPreference set read preference
func (cb *ConcernBuilder) ReadPreference(mode string) *ConcernBuilder {
	cb.options.ReadPreference = mode
	return cb
}

// ApplyToCollection apply to collection
func (cb *ConcernBuilder) ApplyToCollection(collectionName string) error {
	return cb.manager.SetCollectionConcerns(
		collectionName,
		cb.options.ReadConcern,
		cb.options.WriteConcern,
		cb.options.ReadPreference,
	)
}

// ApplyToStruct apply to struct
func (cb *ConcernBuilder) ApplyToStruct(v any) error {
	return cb.manager.SetStructConcerns(
		v,
		cb.options.ReadConcern,
		cb.options.WriteConcern,
		cb.options.ReadPreference,
	)
}

// ApplyGlobally apply globally
func (cb *ConcernBuilder) ApplyGlobally() error {
	return cb.manager.SetGlobalConcerns(
		cb.options.ReadConcern,
		cb.options.WriteConcern,
		cb.options.ReadPreference,
	)
}

// Build build concern options
func (cb *ConcernBuilder) Build() *ConcernOptions {
	return cb.options
}

// ConcernValidator concern validator
type ConcernValidator struct{}

// NewConcernValidator create concern validator
func NewConcernValidator() *ConcernValidator {
	return &ConcernValidator{}
}

// ValidateReadConcern validate read concern
func (cv *ConcernValidator) ValidateReadConcern(level string) error {
	validLevels := []string{
		ReadConcernLocal,
		ReadConcernMajority,
		ReadConcernLinearizable,
		ReadConcernSnapshot,
		ReadConcernAvailable,
	}

	for _, valid := range validLevels {
		if level == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid read concern level: %s", level)
}

// ValidateWriteConcern validate write concern
func (cv *ConcernValidator) ValidateWriteConcern(level string) error {
	validLevels := []string{
		WriteConcernUnacknowledged,
		WriteConcernAcknowledged,
		WriteConcernMajority,
		WriteConcernJournaled,
	}

	for _, valid := range validLevels {
		if level == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid write concern level: %s", level)
}

// ValidateReadPreference validate read preference
func (cv *ConcernValidator) ValidateReadPreference(mode string) error {
	validModes := []string{
		ReadPreferencePrimary,
		ReadPreferencePrimaryPreferred,
		ReadPreferenceSecondary,
		ReadPreferenceSecondaryPreferred,
		ReadPreferenceNearest,
	}

	for _, valid := range validModes {
		if mode == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid read preference mode: %s", mode)
}

// ValidateConcerns validate all concerns
func (cv *ConcernValidator) ValidateConcerns(options *ConcernOptions) error {
	if err := cv.ValidateReadConcern(options.ReadConcern); err != nil {
		return fmt.Errorf("read concern validation failed: %w", err)
	}

	if err := cv.ValidateWriteConcern(options.WriteConcern); err != nil {
		return fmt.Errorf("write concern validation failed: %w", err)
	}

	if err := cv.ValidateReadPreference(options.ReadPreference); err != nil {
		return fmt.Errorf("read preference validation failed: %w", err)
	}

	return nil
}
