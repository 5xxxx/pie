package pie

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// OptionsManager options manager
type OptionsManager struct {
	engine *Engine
}

// NewOptionsManager create options manager
func NewOptionsManager(engine *Engine) *OptionsManager {
	return &OptionsManager{
		engine: engine,
	}
}

// ArrayFiltersManager array filters manager
type ArrayFiltersManager struct {
	manager *OptionsManager
	filters []interface{}
}

// LetManager Let variables manager
type LetManager struct {
	manager   *OptionsManager
	variables bson.D
}

// CommentManager comment manager
type CommentManager struct {
	manager *OptionsManager
	comment string
}

// CollationManager collation manager
type CollationManager struct {
	manager   *OptionsManager
	collation bson.D
}

// HintManager index hint manager
type HintManager struct {
	manager *OptionsManager
	hint    interface{}
}

// MaxTimeManager max time manager
type MaxTimeManager struct {
	manager *OptionsManager
	maxTime time.Duration
}

// BatchSizeManager batch size manager
type BatchSizeManager struct {
	manager   *OptionsManager
	batchSize int32
}

// LimitManager limit manager
type LimitManager struct {
	manager *OptionsManager
	limit   int64
}

// SkipManager skip manager
type SkipManager struct {
	manager *OptionsManager
	skip    int64
}

// SortManager sort manager
type SortManager struct {
	manager *OptionsManager
	sort    bson.D
}

// ProjectionManager projection manager
type ProjectionManager struct {
	manager    *OptionsManager
	projection bson.D
}

// UpsertManager upsert manager
type UpsertManager struct {
	manager *OptionsManager
	upsert  bool
}

// ReturnDocumentManager return document manager
type ReturnDocumentManager struct {
	manager        *OptionsManager
	returnDocument string
}

// NewArrayFiltersManager create array filters manager
func (om *OptionsManager) NewArrayFiltersManager() *ArrayFiltersManager {
	return &ArrayFiltersManager{
		manager: om,
		filters: make([]interface{}, 0),
	}
}

// NewLetManager create Let variables manager
func (om *OptionsManager) NewLetManager() *LetManager {
	return &LetManager{
		manager:   om,
		variables: bson.D{},
	}
}

// NewCommentManager create comment manager
func (om *OptionsManager) NewCommentManager() *CommentManager {
	return &CommentManager{
		manager: om,
	}
}

// NewCollationManager create collation manager
func (om *OptionsManager) NewCollationManager() *CollationManager {
	return &CollationManager{
		manager: om,
	}
}

// NewHintManager create index hint manager
func (om *OptionsManager) NewHintManager() *HintManager {
	return &HintManager{
		manager: om,
	}
}

// NewMaxTimeManager create max time manager
func (om *OptionsManager) NewMaxTimeManager() *MaxTimeManager {
	return &MaxTimeManager{
		manager: om,
	}
}

// NewBatchSizeManager create batch size manager
func (om *OptionsManager) NewBatchSizeManager() *BatchSizeManager {
	return &BatchSizeManager{
		manager: om,
	}
}

// NewLimitManager create limit manager
func (om *OptionsManager) NewLimitManager() *LimitManager {
	return &LimitManager{
		manager: om,
	}
}

// NewSkipManager create skip manager
func (om *OptionsManager) NewSkipManager() *SkipManager {
	return &SkipManager{
		manager: om,
	}
}

// NewSortManager create sort manager
func (om *OptionsManager) NewSortManager() *SortManager {
	return &SortManager{
		manager: om,
	}
}

// NewProjectionManager create projection manager
func (om *OptionsManager) NewProjectionManager() *ProjectionManager {
	return &ProjectionManager{
		manager: om,
	}
}

// NewUpsertManager create upsert manager
func (om *OptionsManager) NewUpsertManager() *UpsertManager {
	return &UpsertManager{
		manager: om,
	}
}

// NewReturnDocumentManager create return document manager
func (om *OptionsManager) NewReturnDocumentManager() *ReturnDocumentManager {
	return &ReturnDocumentManager{
		manager: om,
	}
}

// ArrayFiltersManager methods

// AddFilter add array filter
func (afm *ArrayFiltersManager) AddFilter(filter bson.D) *ArrayFiltersManager {
	afm.filters = append(afm.filters, filter)
	return afm
}

// AddFilters add multiple array filters
func (afm *ArrayFiltersManager) AddFilters(filters ...bson.D) *ArrayFiltersManager {
	for _, filter := range filters {
		afm.filters = append(afm.filters, filter)
	}
	return afm
}

// Clear clear all filters
func (afm *ArrayFiltersManager) Clear() *ArrayFiltersManager {
	afm.filters = make([]interface{}, 0)
	return afm
}

// GetFilters get all filters
func (afm *ArrayFiltersManager) GetFilters() []interface{} {
	return afm.filters
}

// Count get filter count
func (afm *ArrayFiltersManager) Count() int {
	return len(afm.filters)
}

// LetManager methods

// SetVariable set Let variable
func (lm *LetManager) SetVariable(name string, value interface{}) *LetManager {
	lm.variables = append(lm.variables, bson.E{Key: name, Value: value})
	return lm
}

// SetVariables set multiple Let variables
func (lm *LetManager) SetVariables(variables bson.D) *LetManager {
	lm.variables = append(lm.variables, variables...)
	return lm
}

// GetVariable get Let variable
func (lm *LetManager) GetVariable(name string) (interface{}, bool) {
	for _, elem := range lm.variables {
		if elem.Key == name {
			return elem.Value, true
		}
	}
	return nil, false
}

// GetVariables get all Let variables
func (lm *LetManager) GetVariables() bson.D {
	return lm.variables
}

// Clear clear all variables
func (lm *LetManager) Clear() *LetManager {
	lm.variables = bson.D{}
	return lm
}

// CommentManager methods

// SetComment set comment
func (cm *CommentManager) SetComment(comment string) *CommentManager {
	cm.comment = comment
	return cm
}

// GetComment get comment
func (cm *CommentManager) GetComment() string {
	return cm.comment
}

// Clear clear comment
func (cm *CommentManager) Clear() *CommentManager {
	cm.comment = ""
	return cm
}

// CollationManager methods

// SetCollation set collation
func (cm *CollationManager) SetCollation(collation bson.D) *CollationManager {
	cm.collation = collation
	return cm
}

// SetLocale set locale
func (cm *CollationManager) SetLocale(locale string) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "locale", Value: locale})
	return cm
}

// SetCaseLevel set case level
func (cm *CollationManager) SetCaseLevel(caseLevel bool) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "caseLevel", Value: caseLevel})
	return cm
}

// SetCaseFirst set case first
func (cm *CollationManager) SetCaseFirst(caseFirst string) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "caseFirst", Value: caseFirst})
	return cm
}

// SetStrength set strength
func (cm *CollationManager) SetStrength(strength int) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "strength", Value: strength})
	return cm
}

// SetNumericOrdering set numeric ordering
func (cm *CollationManager) SetNumericOrdering(numericOrdering bool) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "numericOrdering", Value: numericOrdering})
	return cm
}

// SetAlternate set alternate
func (cm *CollationManager) SetAlternate(alternate string) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "alternate", Value: alternate})
	return cm
}

// SetMaxVariable set max variable
func (cm *CollationManager) SetMaxVariable(maxVariable string) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "maxVariable", Value: maxVariable})
	return cm
}

// SetBackwards set backwards
func (cm *CollationManager) SetBackwards(backwards bool) *CollationManager {
	cm.collation = append(cm.collation, bson.E{Key: "backwards", Value: backwards})
	return cm
}

// GetCollation get collation
func (cm *CollationManager) GetCollation() bson.D {
	return cm.collation
}

// Clear clear collation
func (cm *CollationManager) Clear() *CollationManager {
	cm.collation = bson.D{}
	return cm
}

// HintManager methods

// SetHint set index hint
func (hm *HintManager) SetHint(hint interface{}) *HintManager {
	hm.hint = hint
	return hm
}

// SetIndexHint set index hint
func (hm *HintManager) SetIndexHint(indexName string) *HintManager {
	hm.hint = indexName
	return hm
}

// SetDocumentHint set document hint
func (hm *HintManager) SetDocumentHint(hint bson.D) *HintManager {
	hm.hint = hint
	return hm
}

// GetHint get index hint
func (hm *HintManager) GetHint() interface{} {
	return hm.hint
}

// Clear clear index hint
func (hm *HintManager) Clear() *HintManager {
	hm.hint = nil
	return hm
}

// MaxTimeManager methods

// SetMaxTime set max time
func (mtm *MaxTimeManager) SetMaxTime(duration time.Duration) *MaxTimeManager {
	mtm.maxTime = duration
	return mtm
}

// SetMaxTimeMS set max time (milliseconds)
func (mtm *MaxTimeManager) SetMaxTimeMS(milliseconds int32) *MaxTimeManager {
	mtm.maxTime = time.Duration(milliseconds) * time.Millisecond
	return mtm
}

// GetMaxTime get max time
func (mtm *MaxTimeManager) GetMaxTime() time.Duration {
	return mtm.maxTime
}

// GetMaxTimeMS get max time (milliseconds)
func (mtm *MaxTimeManager) GetMaxTimeMS() int32 {
	return int32(mtm.maxTime.Milliseconds())
}

// Clear clear max time
func (mtm *MaxTimeManager) Clear() *MaxTimeManager {
	mtm.maxTime = 0
	return mtm
}

// BatchSizeManager methods

// SetBatchSize set batch size
func (bsm *BatchSizeManager) SetBatchSize(size int32) *BatchSizeManager {
	bsm.batchSize = size
	return bsm
}

// GetBatchSize get batch size
func (bsm *BatchSizeManager) GetBatchSize() int32 {
	return bsm.batchSize
}

// Clear clear batch size
func (bsm *BatchSizeManager) Clear() *BatchSizeManager {
	bsm.batchSize = 0
	return bsm
}

// LimitManager methods

// SetLimit set limit
func (lm *LimitManager) SetLimit(limit int64) *LimitManager {
	lm.limit = limit
	return lm
}

// GetLimit get limit
func (lm *LimitManager) GetLimit() int64 {
	return lm.limit
}

// Clear clear limit
func (lm *LimitManager) Clear() *LimitManager {
	lm.limit = 0
	return lm
}

// SkipManager methods

// SetSkip set skip
func (sm *SkipManager) SetSkip(skip int64) *SkipManager {
	sm.skip = skip
	return sm
}

// GetSkip get skip
func (sm *SkipManager) GetSkip() int64 {
	return sm.skip
}

// Clear clear skip
func (sm *SkipManager) Clear() *SkipManager {
	sm.skip = 0
	return sm
}

// SortManager methods

// SetSort set sort
func (sm *SortManager) SetSort(sort bson.D) *SortManager {
	sm.sort = sort
	return sm
}

// AddSort add sort field
func (sm *SortManager) AddSort(field string, order int) *SortManager {
	sm.sort = append(sm.sort, bson.E{Key: field, Value: order})
	return sm
}

// AddAscending add ascending sort
func (sm *SortManager) AddAscending(field string) *SortManager {
	sm.sort = append(sm.sort, bson.E{Key: field, Value: 1})
	return sm
}

// AddDescending add descending sort
func (sm *SortManager) AddDescending(field string) *SortManager {
	sm.sort = append(sm.sort, bson.E{Key: field, Value: -1})
	return sm
}

// GetSort get sort
func (sm *SortManager) GetSort() bson.D {
	return sm.sort
}

// Clear clear sort
func (sm *SortManager) Clear() *SortManager {
	sm.sort = bson.D{}
	return sm
}

// ProjectionManager methods

// SetProjection set projection
func (pm *ProjectionManager) SetProjection(projection bson.D) *ProjectionManager {
	pm.projection = projection
	return pm
}

// Include include field
func (pm *ProjectionManager) Include(field string) *ProjectionManager {
	pm.projection = append(pm.projection, bson.E{Key: field, Value: 1})
	return pm
}

// Exclude exclude field
func (pm *ProjectionManager) Exclude(field string) *ProjectionManager {
	pm.projection = append(pm.projection, bson.E{Key: field, Value: 0})
	return pm
}

// IncludeFields include multiple fields
func (pm *ProjectionManager) IncludeFields(fields ...string) *ProjectionManager {
	for _, field := range fields {
		pm.projection = append(pm.projection, bson.E{Key: field, Value: 1})
	}
	return pm
}

// ExcludeFields exclude multiple fields
func (pm *ProjectionManager) ExcludeFields(fields ...string) *ProjectionManager {
	for _, field := range fields {
		pm.projection = append(pm.projection, bson.E{Key: field, Value: 0})
	}
	return pm
}

// GetProjection get projection
func (pm *ProjectionManager) GetProjection() bson.D {
	return pm.projection
}

// Clear clear projection
func (pm *ProjectionManager) Clear() *ProjectionManager {
	pm.projection = bson.D{}
	return pm
}

// UpsertManager methods

// SetUpsert set upsert
func (um *UpsertManager) SetUpsert(upsert bool) *UpsertManager {
	um.upsert = upsert
	return um
}

// EnableUpsert enable upsert
func (um *UpsertManager) EnableUpsert() *UpsertManager {
	um.upsert = true
	return um
}

// DisableUpsert disable upsert
func (um *UpsertManager) DisableUpsert() *UpsertManager {
	um.upsert = false
	return um
}

// GetUpsert get upsert status
func (um *UpsertManager) GetUpsert() bool {
	return um.upsert
}

// ReturnDocumentManager methods

// SetReturnDocument set return document
func (rdm *ReturnDocumentManager) SetReturnDocument(after bool) *ReturnDocumentManager {
	if after {
		rdm.returnDocument = "after"
	} else {
		rdm.returnDocument = "before"
	}
	return rdm
}

// SetReturnAfter set return updated document
func (rdm *ReturnDocumentManager) SetReturnAfter() *ReturnDocumentManager {
	rdm.returnDocument = "after"
	return rdm
}

// SetReturnBefore set return original document
func (rdm *ReturnDocumentManager) SetReturnBefore() *ReturnDocumentManager {
	rdm.returnDocument = "before"
	return rdm
}

// GetReturnDocument get return document setting
func (rdm *ReturnDocumentManager) GetReturnDocument() string {
	return rdm.returnDocument
}

// Clear clear return document setting
func (rdm *ReturnDocumentManager) Clear() *ReturnDocumentManager {
	rdm.returnDocument = ""
	return rdm
}

// OptionsBuilder options builder
type OptionsBuilder struct {
	manager        *OptionsManager
	arrayFilters   *ArrayFiltersManager
	letVars        *LetManager
	comment        *CommentManager
	collation      *CollationManager
	hint           *HintManager
	maxTime        *MaxTimeManager
	batchSize      *BatchSizeManager
	limit          *LimitManager
	skip           *SkipManager
	sort           *SortManager
	projection     *ProjectionManager
	upsert         *UpsertManager
	returnDocument *ReturnDocumentManager
}

// NewOptionsBuilder create options builder
func (om *OptionsManager) NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{
		manager:        om,
		arrayFilters:   om.NewArrayFiltersManager(),
		letVars:        om.NewLetManager(),
		comment:        om.NewCommentManager(),
		collation:      om.NewCollationManager(),
		hint:           om.NewHintManager(),
		maxTime:        om.NewMaxTimeManager(),
		batchSize:      om.NewBatchSizeManager(),
		limit:          om.NewLimitManager(),
		skip:           om.NewSkipManager(),
		sort:           om.NewSortManager(),
		projection:     om.NewProjectionManager(),
		upsert:         om.NewUpsertManager(),
		returnDocument: om.NewReturnDocumentManager(),
	}
}

// ArrayFilters get array filters manager
func (ob *OptionsBuilder) ArrayFilters() *ArrayFiltersManager {
	return ob.arrayFilters
}

// Let get Let variables manager
func (ob *OptionsBuilder) Let() *LetManager {
	return ob.letVars
}

// Comment get comment manager
func (ob *OptionsBuilder) Comment() *CommentManager {
	return ob.comment
}

// Collation get collation manager
func (ob *OptionsBuilder) Collation() *CollationManager {
	return ob.collation
}

// Hint get index hint manager
func (ob *OptionsBuilder) Hint() *HintManager {
	return ob.hint
}

// MaxTime get max time manager
func (ob *OptionsBuilder) MaxTime() *MaxTimeManager {
	return ob.maxTime
}

// BatchSize get batch size manager
func (ob *OptionsBuilder) BatchSize() *BatchSizeManager {
	return ob.batchSize
}

// Limit get limit manager
func (ob *OptionsBuilder) Limit() *LimitManager {
	return ob.limit
}

// Skip get skip manager
func (ob *OptionsBuilder) Skip() *SkipManager {
	return ob.skip
}

// Sort get sort manager
func (ob *OptionsBuilder) Sort() *SortManager {
	return ob.sort
}

// Projection get projection manager
func (ob *OptionsBuilder) Projection() *ProjectionManager {
	return ob.projection
}

// Upsert get upsert manager
func (ob *OptionsBuilder) Upsert() *UpsertManager {
	return ob.upsert
}

// ReturnDocument get return document manager
func (ob *OptionsBuilder) ReturnDocument() *ReturnDocumentManager {
	return ob.returnDocument
}

// BuildFindOptions build Find options
func (ob *OptionsBuilder) BuildFindOptions() *options.FindOptionsBuilder {
	opts := options.Find()

	if ob.limit.GetLimit() > 0 {
		opts.SetLimit(ob.limit.GetLimit())
	}

	if ob.skip.GetSkip() > 0 {
		opts.SetSkip(ob.skip.GetSkip())
	}

	if len(ob.sort.GetSort()) > 0 {
		opts.SetSort(ob.sort.GetSort())
	}

	if len(ob.projection.GetProjection()) > 0 {
		opts.SetProjection(ob.projection.GetProjection())
	}

	if ob.batchSize.GetBatchSize() > 0 {
		opts.SetBatchSize(ob.batchSize.GetBatchSize())
	}

	if ob.hint.GetHint() != nil {
		opts.SetHint(ob.hint.GetHint())
	}

	if ob.comment.GetComment() != "" {
		// Note: SetComment not available in v2 FindOptions
		// This method is provided for API compatibility
	}

	return opts
}

// BuildUpdateOptions build Update options
func (ob *OptionsBuilder) BuildUpdateOptions() *options.UpdateOneOptionsBuilder {
	opts := options.UpdateOne()

	if ob.upsert.GetUpsert() {
		opts.SetUpsert(true)
	}

	if len(ob.arrayFilters.GetFilters()) > 0 {
		opts.SetArrayFilters(ob.arrayFilters.GetFilters())
	}

	if len(ob.letVars.GetVariables()) > 0 {
		opts.SetLet(ob.letVars.GetVariables())
	}

	if ob.collation.GetCollation() != nil {
		// Note: SetCollation not available in v2 UpdateOptions
		// This method is provided for API compatibility
	}

	return opts
}

// BuildDeleteOptions build Delete options
func (ob *OptionsBuilder) BuildDeleteOptions() *options.DeleteOneOptionsBuilder {
	opts := options.DeleteOne()

	if ob.collation.GetCollation() != nil {
		// Note: SetCollation not available in v2 DeleteOptions
		// This method is provided for API compatibility
	}

	return opts
}

// BuildInsertOptions build Insert options
func (ob *OptionsBuilder) BuildInsertOptions() *options.InsertOneOptionsBuilder {
	opts := options.InsertOne()

	if ob.comment.GetComment() != "" {
		// Note: SetComment not available in v2 InsertOptions
		// This method is provided for API compatibility
	}

	return opts
}

// BuildAggregateOptions build Aggregate options
func (ob *OptionsBuilder) BuildAggregateOptions() *options.AggregateOptionsBuilder {
	opts := options.Aggregate()

	if ob.batchSize.GetBatchSize() > 0 {
		opts.SetBatchSize(ob.batchSize.GetBatchSize())
	}

	if ob.maxTime.GetMaxTime() > 0 {
		// Note: SetMaxTime not available in v2 AggregateOptions
		// This method is provided for API compatibility
	}

	if ob.comment.GetComment() != "" {
		opts.SetComment(ob.comment.GetComment())
	}

	if ob.hint.GetHint() != nil {
		opts.SetHint(ob.hint.GetHint())
	}

	if ob.collation.GetCollation() != nil {
		// Note: SetCollation not available in v2 AggregateOptions
		// This method is provided for API compatibility
	}

	return opts
}

// Clear clear all options
func (ob *OptionsBuilder) Clear() *OptionsBuilder {
	ob.arrayFilters.Clear()
	ob.letVars.Clear()
	ob.comment.Clear()
	ob.collation.Clear()
	ob.hint.Clear()
	ob.maxTime.Clear()
	ob.batchSize.Clear()
	ob.limit.Clear()
	ob.skip.Clear()
	ob.sort.Clear()
	ob.projection.Clear()
	ob.upsert.DisableUpsert()
	ob.returnDocument.Clear()
	return ob
}

// OptionsValidator options validator
type OptionsValidator struct{}

// NewOptionsValidator create options validator
func NewOptionsValidator() *OptionsValidator {
	return &OptionsValidator{}
}

// ValidateArrayFilters validate array filters
func (ov *OptionsValidator) ValidateArrayFilters(filters []interface{}) error {
	if len(filters) > 10 {
		return fmt.Errorf("too many array filters: %d (max 10)", len(filters))
	}
	return nil
}

// ValidateBatchSize validate batch size
func (ov *OptionsValidator) ValidateBatchSize(size int32) error {
	if size < 0 {
		return fmt.Errorf("batch size cannot be negative: %d", size)
	}
	if size > 1000 {
		return fmt.Errorf("batch size too large: %d (max 1000)", size)
	}
	return nil
}

// ValidateLimit validate limit
func (ov *OptionsValidator) ValidateLimit(limit int64) error {
	if limit < 0 {
		return fmt.Errorf("limit cannot be negative: %d", limit)
	}
	return nil
}

// ValidateSkip validate skip
func (ov *OptionsValidator) ValidateSkip(skip int64) error {
	if skip < 0 {
		return fmt.Errorf("skip cannot be negative: %d", skip)
	}
	return nil
}

// ValidateMaxTime validate max time
func (ov *OptionsValidator) ValidateMaxTime(duration time.Duration) error {
	if duration < 0 {
		return fmt.Errorf("max time cannot be negative: %v", duration)
	}
	if duration > 30*time.Minute {
		return fmt.Errorf("max time too large: %v (max 30 minutes)", duration)
	}
	return nil
}

// ValidateComment validate comment
func (ov *OptionsValidator) ValidateComment(comment string) error {
	if len(comment) > 1000 {
		return fmt.Errorf("comment too long: %d characters (max 1000)", len(comment))
	}
	return nil
}

// ValidateAll validate all options
func (ov *OptionsValidator) ValidateAll(builder *OptionsBuilder) error {
	if err := ov.ValidateArrayFilters(builder.arrayFilters.GetFilters()); err != nil {
		return fmt.Errorf("array filters validation failed: %w", err)
	}

	if err := ov.ValidateBatchSize(builder.batchSize.GetBatchSize()); err != nil {
		return fmt.Errorf("batch size validation failed: %w", err)
	}

	if err := ov.ValidateLimit(builder.limit.GetLimit()); err != nil {
		return fmt.Errorf("limit validation failed: %w", err)
	}

	if err := ov.ValidateSkip(builder.skip.GetSkip()); err != nil {
		return fmt.Errorf("skip validation failed: %w", err)
	}

	if err := ov.ValidateMaxTime(builder.maxTime.GetMaxTime()); err != nil {
		return fmt.Errorf("max time validation failed: %w", err)
	}

	if err := ov.ValidateComment(builder.comment.GetComment()); err != nil {
		return fmt.Errorf("comment validation failed: %w", err)
	}

	return nil
}
