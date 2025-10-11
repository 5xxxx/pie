package pie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Aggregate aggregation operation builder
type Aggregate[T any] struct {
	engine     *Engine
	collection *mongo.Collection
	pipeline   bson.A
	options    *options.AggregateOptionsBuilder
}

// NewAggregate create new aggregation operation
func NewAggregate[T any](engine *Engine) *Aggregate[T] {
	return &Aggregate[T]{
		engine:   engine,
		pipeline: bson.A{},
		options:  options.Aggregate(),
	}
}

// Collection set collection
func (a *Aggregate[T]) Collection(name string) *Aggregate[T] {
	a.collection = a.engine.Collection(name)
	return a
}

// CollectionForStruct set collection by struct
func (a *Aggregate[T]) CollectionForStruct(v interface{}) *Aggregate[T] {
	collection, err := a.engine.CollectionForStruct(v)
	if err != nil {
		// Here we can log the error but don't interrupt the chain call
		return a
	}
	a.collection = collection
	return a
}

// ========== 阶段方法 ==========

// MatchStage 创建匹配阶段构建器
func (a *Aggregate[T]) MatchStage() *MatchStage[T] {
	return &MatchStage[T]{
		agg:     a,
		filters: []bson.D{},
	}
}

// Match 添加匹配条件（便捷方法）
func (a *Aggregate[T]) Match(filter bson.D) *Aggregate[T] {
	if len(filter) == 0 {
		return a
	}
	a.pipeline = append(a.pipeline, bson.D{{"$match", filter}})
	return a
}

// AddFieldsStage 创建添加字段阶段构建器
func (a *Aggregate[T]) AddFieldsStage() *AddFieldsStage[T] {
	return &AddFieldsStage[T]{
		agg:    a,
		fields: bson.D{},
	}
}

// SetStage 创建设置字段阶段构建器(AddFields的别名)
func (a *Aggregate[T]) SetStage() *AddFieldsStage[T] {
	return a.AddFieldsStage()
}

// UnsetStage 移除字段阶段
func (a *Aggregate[T]) UnsetStage(fields ...string) *Aggregate[T] {
	unsetDoc := bson.D{}
	for _, field := range fields {
		unsetDoc = append(unsetDoc, bson.E{Key: field, Value: ""})
	}
	a.pipeline = append(a.pipeline, bson.D{{"$unset", unsetDoc}})
	return a
}

// ReplaceRootStage 替换根文档阶段
func (a *Aggregate[T]) ReplaceRootStage(newRoot any) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$replaceRoot", bson.D{{"newRoot", newRoot}}}})
	return a
}

// ReplaceWithStage 替换文档阶段
func (a *Aggregate[T]) ReplaceWithStage(replacement any) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$replaceWith", replacement}})
	return a
}

// UnwindStage 创建展开阶段构建器
func (a *Aggregate[T]) UnwindStage(path string) *UnwindStage[T] {
	return &UnwindStage[T]{
		agg:     a,
		path:    path,
		options: bson.D{},
	}
}

// GroupStage 创建分组阶段构建器
func (a *Aggregate[T]) GroupStage() *GroupStage[T] {
	return &GroupStage[T]{
		agg:          a,
		id:           bson.D{},
		accumulators: bson.D{},
	}
}

// Group 添加分组阶段（便捷方法）
func (a *Aggregate[T]) Group(id bson.D, accumulators ...bson.E) *Aggregate[T] {
	groupDoc := bson.D{{"_id", id}}
	groupDoc = append(groupDoc, accumulators...)
	a.pipeline = append(a.pipeline, bson.D{{"$group", groupDoc}})
	return a
}

// ProjectStage 创建投影阶段构建器
func (a *Aggregate[T]) ProjectStage() *ProjectStage[T] {
	return &ProjectStage[T]{
		agg:    a,
		fields: bson.D{},
	}
}

// SortStage 创建排序阶段构建器
func (a *Aggregate[T]) SortStage() *SortStage[T] {
	return &SortStage[T]{
		agg:   a,
		sorts: bson.D{},
	}
}

// Sort 添加排序阶段（便捷方法）
func (a *Aggregate[T]) Sort(sort bson.D) *Aggregate[T] {
	if len(sort) == 0 {
		return a
	}
	a.pipeline = append(a.pipeline, bson.D{{"$sort", sort}})
	return a
}

// LimitStage 限制阶段
func (a *Aggregate[T]) LimitStage(n int64) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$limit", n}})
	return a
}

// SkipStage 跳过阶段
func (a *Aggregate[T]) SkipStage(n int64) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$skip", n}})
	return a
}

// LookupStage 创建关联查询阶段构建器
func (a *Aggregate[T]) LookupStage(from, localField, foreignField, as string) *LookupStage[T] {
	return &LookupStage[T]{
		agg:  a,
		from: from,
		options: bson.D{
			{"localField", localField},
			{"foreignField", foreignField},
			{"as", as},
		},
	}
}

// GraphLookupStage 创建图查找阶段构建器
func (a *Aggregate[T]) GraphLookupStage(from string) *GraphLookupStage[T] {
	return &GraphLookupStage[T]{
		agg:     a,
		options: bson.D{{"from", from}},
	}
}

// UnionWithStage 创建联合阶段构建器
func (a *Aggregate[T]) UnionWithStage(collection string) *UnionWithStage[T] {
	return &UnionWithStage[T]{
		agg:     a,
		options: bson.D{{"coll", collection}},
	}
}

// FacetStage 创建分面阶段构建器
func (a *Aggregate[T]) FacetStage() *FacetStage[T] {
	return &FacetStage[T]{
		agg:    a,
		facets: bson.D{},
	}
}

// SampleStage 采样阶段
func (a *Aggregate[T]) SampleStage(size int64) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$sample", bson.D{{"size", size}}}})
	return a
}

// CountStage 计数阶段
func (a *Aggregate[T]) CountStage(field string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$count", field}})
	return a
}

// OutStage 输出阶段
func (a *Aggregate[T]) OutStage(collection string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$out", collection}})
	return a
}

// MergeStage 创建合并阶段构建器
func (a *Aggregate[T]) MergeStage(into string) *MergeStage[T] {
	return &MergeStage[T]{
		agg:     a,
		into:    into,
		options: bson.D{},
	}
}

// BucketStage 创建分桶阶段构建器
func (a *Aggregate[T]) BucketStage() *BucketStage[T] {
	return &BucketStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// BucketAutoStage 创建自动分桶阶段构建器
func (a *Aggregate[T]) BucketAutoStage() *BucketAutoStage[T] {
	return &BucketAutoStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// SortByCountStage 按计数排序阶段
func (a *Aggregate[T]) SortByCountStage(field string) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$sortByCount", field}})
	return a
}

// RedactStage 编辑阶段
func (a *Aggregate[T]) RedactStage(expression any) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$redact", expression}})
	return a
}

// GeoNearStage 创建地理邻近阶段构建器
func (a *Aggregate[T]) GeoNearStage() *GeoNearStage[T] {
	return &GeoNearStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// SetWindowFieldsStage 创建窗口字段阶段构建器
func (a *Aggregate[T]) SetWindowFieldsStage() *SetWindowFieldsStage[T] {
	return &SetWindowFieldsStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// DocumentsStage 文档阶段
func (a *Aggregate[T]) DocumentsStage(documents ...any) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$documents", documents}})
	return a
}

// SearchStage 创建搜索阶段构建器
func (a *Aggregate[T]) SearchStage() *SearchStage[T] {
	return &SearchStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// SearchMetaStage 创建搜索元数据阶段构建器
func (a *Aggregate[T]) SearchMetaStage() *SearchMetaStage[T] {
	return &SearchMetaStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// VectorSearchStage 创建向量搜索阶段构建器
func (a *Aggregate[T]) VectorSearchStage() *VectorSearchStage[T] {
	return &VectorSearchStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// DensifyStage 创建密度化阶段构建器
func (a *Aggregate[T]) DensifyStage() *DensifyStage[T] {
	return &DensifyStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// FillStage 创建填充阶段构建器
func (a *Aggregate[T]) FillStage() *FillStage[T] {
	return &FillStage[T]{
		agg:     a,
		options: bson.D{},
	}
}

// CollStatsStage 集合统计阶段
func (a *Aggregate[T]) CollStatsStage(options M) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$collStats", options}})
	return a
}

// IndexStatsStage 索引统计阶段
func (a *Aggregate[T]) IndexStatsStage() *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$indexStats", bson.D{}}})
	return a
}

// PlanCacheStatsStage 查询计划缓存统计阶段
func (a *Aggregate[T]) PlanCacheStatsStage() *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$planCacheStats", bson.D{}}})
	return a
}

// CurrentOpStage 当前操作阶段
func (a *Aggregate[T]) CurrentOpStage(options M) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$currentOp", options}})
	return a
}

// ListSessionsStage 会话列表阶段
func (a *Aggregate[T]) ListSessionsStage(options M) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$listSessions", options}})
	return a
}

// ListSampledQueriesStage 采样查询列表阶段
func (a *Aggregate[T]) ListSampledQueriesStage(options M) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$listSampledQueries", options}})
	return a
}

// ChangeStreamStage 更改流阶段
func (a *Aggregate[T]) ChangeStreamStage(options M) *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$changeStream", options}})
	return a
}

// ChangeStreamSplitLargeEventStage 更改流分割大事件阶段
func (a *Aggregate[T]) ChangeStreamSplitLargeEventStage() *Aggregate[T] {
	a.pipeline = append(a.pipeline, bson.D{{"$changeStreamSplitLargeEvent", bson.D{}}})
	return a
}

// RawStage 原始阶段
func (a *Aggregate[T]) RawStage(stage M) *Aggregate[T] {
	a.pipeline = append(a.pipeline, stage)
	return a
}

// Option setting methods

// SetAllowDiskUse set allow disk use
func (a *Aggregate[T]) SetAllowDiskUse(allow bool) *Aggregate[T] {
	a.options.SetAllowDiskUse(allow)
	return a
}

// SetBatchSize set batch size
func (a *Aggregate[T]) SetBatchSize(size int32) *Aggregate[T] {
	a.options.SetBatchSize(size)
	return a
}

// SetBypassDocumentValidation set bypass document validation
func (a *Aggregate[T]) SetBypassDocumentValidation(bypass bool) *Aggregate[T] {
	a.options.SetBypassDocumentValidation(bypass)
	return a
}

// SetCollation set collation
func (a *Aggregate[T]) SetCollation(collation *options.Collation) *Aggregate[T] {
	a.options.SetCollation(collation)
	return a
}

// SetMaxAwaitTime set max await time
func (a *Aggregate[T]) SetMaxAwaitTime(duration int64) *Aggregate[T] {
	a.options.SetMaxAwaitTime(time.Duration(duration))
	return a
}

// SetComment set comment
func (a *Aggregate[T]) SetComment(comment string) *Aggregate[T] {
	a.options.SetComment(comment)
	return a
}

// SetHint set hint
func (a *Aggregate[T]) SetHint(hint interface{}) *Aggregate[T] {
	a.options.SetHint(hint)
	return a
}

// Exec execute aggregation operation
func (a *Aggregate[T]) Exec(ctx context.Context) (*AggregateResult[T], error) {
	if a.collection == nil {
		// Try to get collection based on generic type
		var zero T
		collection, err := a.engine.CollectionForStruct(zero)
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}
		a.collection = collection
	}

	cursor, err := a.collection.Aggregate(ctx, a.pipeline, a.options)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	return &AggregateResult[T]{
		Data:  results,
		Error: nil,
	}, nil
}

// ExecOne execute aggregation operation and return single result
func (a *Aggregate[T]) ExecOne(ctx context.Context, result *T) error {
	if a.collection == nil {
		var zero T
		collection, err := a.engine.CollectionForStruct(zero)
		if err != nil {
			return fmt.Errorf("failed to get collection: %w", err)
		}
		a.collection = collection
	}

	cursor, err := a.collection.Aggregate(ctx, a.pipeline, a.options)
	if err != nil {
		return fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(result); err != nil {
			return fmt.Errorf("failed to decode aggregation result: %w", err)
		}
		return nil
	}

	return ErrEmptyResult
}

// GetPipeline get pipeline
func (a *Aggregate[T]) GetPipeline() bson.A {
	return a.pipeline
}

// Clone clone aggregation operation
func (a *Aggregate[T]) Clone() *Aggregate[T] {
	newPipeline := make(bson.A, len(a.pipeline))
	copy(newPipeline, a.pipeline)

	return &Aggregate[T]{
		engine:     a.engine,
		collection: a.collection,
		pipeline:   newPipeline,
		options:    a.options,
	}
}

// Clear clear pipeline
func (a *Aggregate[T]) Clear() *Aggregate[T] {
	a.pipeline = bson.A{}
	return a
}
