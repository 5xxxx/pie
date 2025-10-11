package pie

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ========== 阶段构建器类型 ==========

// MatchStage 匹配阶段构建器
type MatchStage[T any] struct {
	agg     *Aggregate[T]
	filters []bson.D
}

// AddFieldsStage 添加字段阶段构建器
type AddFieldsStage[T any] struct {
	agg    *Aggregate[T]
	fields bson.D
}

// UnwindStage 展开阶段构建器
type UnwindStage[T any] struct {
	agg     *Aggregate[T]
	path    string
	options bson.D
}

// GroupStage 分组阶段构建器
type GroupStage[T any] struct {
	agg          *Aggregate[T]
	id           bson.D
	accumulators bson.D
}

// ProjectStage 投影阶段构建器
type ProjectStage[T any] struct {
	agg    *Aggregate[T]
	fields bson.D
}

// SortStage 排序阶段构建器
type SortStage[T any] struct {
	agg   *Aggregate[T]
	sorts bson.D
}

// LookupStage 关联查询阶段构建器
type LookupStage[T any] struct {
	agg     *Aggregate[T]
	from    string
	options bson.D
}

// MergeStage 合并阶段构建器
type MergeStage[T any] struct {
	agg     *Aggregate[T]
	into    string
	options bson.D
}

// FacetStage 分面阶段构建器
type FacetStage[T any] struct {
	agg    *Aggregate[T]
	facets bson.D
}

// BucketStage 分桶阶段构建器
type BucketStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// BucketAutoStage 自动分桶阶段构建器
type BucketAutoStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// SetWindowFieldsStage 窗口字段阶段构建器
type SetWindowFieldsStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// SearchStage 搜索阶段构建器
type SearchStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// SearchMetaStage 搜索元数据阶段构建器
type SearchMetaStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// VectorSearchStage 向量搜索阶段构建器
type VectorSearchStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// DensifyStage 密度化阶段构建器
type DensifyStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// FillStage 填充阶段构建器
type FillStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// UnionWithStage 联合阶段构建器
type UnionWithStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// GraphLookupStage 图查找阶段构建器
type GraphLookupStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// GeoNearStage 地理邻近阶段构建器
type GeoNearStage[T any] struct {
	agg     *Aggregate[T]
	options bson.D
}

// ========== MatchStage 构建器方法 ==========

// Where 添加相等条件并自动完成阶段
func (m *MatchStage[T]) Where(field string, value any) *Aggregate[T] {
	m.filters = append(m.filters, bson.D{{Key: field, Value: value}})
	return m.Done()
}

// Operator 添加操作符条件
func (m *MatchStage[T]) Operator(op Operator) *MatchStage[T] {
	m.filters = append(m.filters, op.ToBSON())
	return m
}

// Between 添加范围条件并自动完成阶段
func (m *MatchStage[T]) Between(field string, min, max any) *Aggregate[T] {
	m.filters = append(m.filters, bson.D{{Key: field, Value: bson.D{
		{Key: "$gte", Value: min},
		{Key: "$lte", Value: max},
	}}})
	return m.Done()
}

// In 添加包含条件
func (m *MatchStage[T]) In(field string, values ...any) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{field, bson.D{{"$in", values}}}})
	return m
}

// Regex 添加正则表达式条件
func (m *MatchStage[T]) Regex(field, pattern string) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{field, bson.D{{"$regex", pattern}}}})
	return m
}

// Exists 添加存在条件
func (m *MatchStage[T]) Exists(field string, exists bool) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{field, bson.D{{"$exists", exists}}}})
	return m
}

// Text 添加文本搜索条件
func (m *MatchStage[T]) Text(search string) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{"$text", bson.D{{"$search", search}}}})
	return m
}

// And 添加AND条件
func (m *MatchStage[T]) And(conditions ...bson.D) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{"$and", conditions}})
	return m
}

// Or 添加OR条件
func (m *MatchStage[T]) Or(conditions ...bson.D) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{"$or", conditions}})
	return m
}

// Nor 添加NOR条件
func (m *MatchStage[T]) Nor(conditions ...bson.D) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{"$nor", conditions}})
	return m
}

// Not 添加NOT条件
func (m *MatchStage[T]) Not(condition bson.D) *MatchStage[T] {
	m.filters = append(m.filters, bson.D{{"$not", condition}})
	return m
}

// Raw 添加原始条件
func (m *MatchStage[T]) Raw(filter bson.D) *MatchStage[T] {
	m.filters = append(m.filters, filter)
	return m
}

// Done 完成匹配阶段构建
func (m *MatchStage[T]) Done() *Aggregate[T] {
	if len(m.filters) == 0 {
		return m.agg
	}

	// 合并所有过滤器
	var combinedFilter bson.D
	if len(m.filters) == 1 {
		combinedFilter = m.filters[0]
	} else {
		combinedFilter = bson.D{{"$and", m.filters}}
	}

	m.agg.pipeline = append(m.agg.pipeline, bson.D{{"$match", combinedFilter}})
	return m.agg
}

// ========== AddFieldsStage 构建器方法 ==========

// Add 添加字段
func (a *AddFieldsStage[T]) Add(name string, value any) *AddFieldsStage[T] {
	a.fields = append(a.fields, bson.E{Key: name, Value: value})
	return a
}

// Done 完成添加字段阶段构建
func (a *AddFieldsStage[T]) Done() *Aggregate[T] {
	if len(a.fields) == 0 {
		return a.agg
	}
	a.agg.pipeline = append(a.agg.pipeline, bson.D{{"$addFields", a.fields}})
	return a.agg
}

// ========== UnwindStage 构建器方法 ==========

// PreserveNullAndEmptyArrays 设置是否保留空值
func (u *UnwindStage[T]) PreserveNullAndEmptyArrays(preserve bool) *UnwindStage[T] {
	u.options = append(u.options, bson.E{Key: "preserveNullAndEmptyArrays", Value: preserve})
	return u
}

// IncludeArrayIndex 设置数组索引字段
func (u *UnwindStage[T]) IncludeArrayIndex(arrayIndex string) *UnwindStage[T] {
	u.options = append(u.options, bson.E{Key: "includeArrayIndex", Value: arrayIndex})
	return u
}

// Done 完成展开阶段构建
func (u *UnwindStage[T]) Done() *Aggregate[T] {
	unwindDoc := bson.D{{"path", u.path}}
	unwindDoc = append(unwindDoc, u.options...)
	u.agg.pipeline = append(u.agg.pipeline, bson.D{{"$unwind", unwindDoc}})
	return u.agg
}

// ========== GroupStage 构建器方法 ==========

// By 添加分组字段
func (g *GroupStage[T]) By(name string, field any) *GroupStage[T] {
	g.id = append(g.id, bson.E{Key: name, Value: field})
	return g
}

// ByRaw 直接设置_id
func (g *GroupStage[T]) ByRaw(id any) *GroupStage[T] {
	g.id = bson.D{{"_id", id}}
	return g
}

// Sum 添加求和累加器
func (g *GroupStage[T]) Sum(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$sum", field}}})
	return g
}

// Avg 添加平均值累加器
func (g *GroupStage[T]) Avg(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$avg", field}}})
	return g
}

// Count 添加计数累加器
func (g *GroupStage[T]) Count(name string) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{Key: "$sum", Value: 1}}})
	return g
}

// Max 添加最大值累加器
func (g *GroupStage[T]) Max(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$max", field}}})
	return g
}

// Min 添加最小值累加器
func (g *GroupStage[T]) Min(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$min", field}}})
	return g
}

// First 添加第一个值累加器
func (g *GroupStage[T]) First(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$first", field}}})
	return g
}

// Last 添加最后一个值累加器
func (g *GroupStage[T]) Last(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$last", field}}})
	return g
}

// Push 添加推入数组累加器
func (g *GroupStage[T]) Push(name string, value any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$push", value}}})
	return g
}

// AddToSet 添加添加到集合累加器
func (g *GroupStage[T]) AddToSet(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$addToSet", field}}})
	return g
}

// StdDevPop 添加总体标准差累加器
func (g *GroupStage[T]) StdDevPop(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$stdDevPop", field}}})
	return g
}

// StdDevSamp 添加样本标准差累加器
func (g *GroupStage[T]) StdDevSamp(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$stdDevSamp", field}}})
	return g
}

// MergeObjects 添加合并对象累加器
func (g *GroupStage[T]) MergeObjects(name string, field any) *GroupStage[T] {
	g.accumulators = append(g.accumulators, bson.E{Key: name, Value: bson.D{{"$mergeObjects", field}}})
	return g
}

// Done 完成分组阶段构建
func (g *GroupStage[T]) Done() *Aggregate[T] {
	groupDoc := bson.D{{"_id", g.id}}
	groupDoc = append(groupDoc, g.accumulators...)
	g.agg.pipeline = append(g.agg.pipeline, bson.D{{"$group", groupDoc}})
	return g.agg
}

// ========== ProjectStage 构建器方法 ==========

// Include 包含字段
func (p *ProjectStage[T]) Include(names ...string) *ProjectStage[T] {
	for _, name := range names {
		p.fields = append(p.fields, bson.E{Key: name, Value: 1})
	}
	return p
}

// Exclude 排除字段
func (p *ProjectStage[T]) Exclude(names ...string) *ProjectStage[T] {
	for _, name := range names {
		p.fields = append(p.fields, bson.E{Key: name, Value: 0})
	}
	return p
}

// Field 添加计算字段
func (p *ProjectStage[T]) Field(name string, value any) *ProjectStage[T] {
	p.fields = append(p.fields, bson.E{Key: name, Value: value})
	return p
}

// Slice 切片数组字段
func (p *ProjectStage[T]) Slice(name, arrayField string, n int) *ProjectStage[T] {
	p.fields = append(p.fields, bson.E{Key: name, Value: bson.D{{"$slice", []any{arrayField, n}}}})
	return p
}

// Done 完成投影阶段构建
func (p *ProjectStage[T]) Done() *Aggregate[T] {
	if len(p.fields) == 0 {
		return p.agg
	}
	p.agg.pipeline = append(p.agg.pipeline, bson.D{{"$project", p.fields}})
	return p.agg
}

// ========== SortStage 构建器方法 ==========

// Asc 升序排序
func (s *SortStage[T]) Asc(fields ...string) *SortStage[T] {
	for _, field := range fields {
		s.sorts = append(s.sorts, bson.E{Key: field, Value: 1})
	}
	return s
}

// Desc 降序排序并自动完成阶段
func (s *SortStage[T]) Desc(fields ...string) *Aggregate[T] {
	for _, field := range fields {
		s.sorts = append(s.sorts, bson.E{Key: field, Value: -1})
	}
	return s.Done()
}

// Field 添加排序字段
func (s *SortStage[T]) Field(field string, order int) *SortStage[T] {
	s.sorts = append(s.sorts, bson.E{Key: field, Value: order})
	return s
}

// Done 完成排序阶段构建
func (s *SortStage[T]) Done() *Aggregate[T] {
	if len(s.sorts) == 0 {
		return s.agg
	}
	s.agg.pipeline = append(s.agg.pipeline, bson.D{{"$sort", s.sorts}})
	return s.agg
}

// ========== LookupStage 构建器方法 ==========

// Let 设置变量
func (l *LookupStage[T]) Let(vars M) *LookupStage[T] {
	l.options = append(l.options, bson.E{Key: "let", Value: vars})
	return l
}

// Pipeline 设置管道
func (l *LookupStage[T]) Pipeline(stages ...any) *LookupStage[T] {
	l.options = append(l.options, bson.E{Key: "pipeline", Value: stages})
	return l
}

// Done 完成关联查询阶段构建
func (l *LookupStage[T]) Done() *Aggregate[T] {
	lookupDoc := bson.D{
		{"from", l.from},
	}
	lookupDoc = append(lookupDoc, l.options...)
	l.agg.pipeline = append(l.agg.pipeline, bson.D{{"$lookup", lookupDoc}})
	return l.agg
}

// ========== MergeStage 构建器方法 ==========

// On 设置匹配字段
func (m *MergeStage[T]) On(fields ...string) *MergeStage[T] {
	m.options = append(m.options, bson.E{Key: "on", Value: fields})
	return m
}

// WhenMatched 设置匹配时的操作
func (m *MergeStage[T]) WhenMatched(action string) *MergeStage[T] {
	m.options = append(m.options, bson.E{Key: "whenMatched", Value: action})
	return m
}

// WhenNotMatched 设置不匹配时的操作
func (m *MergeStage[T]) WhenNotMatched(action string) *MergeStage[T] {
	m.options = append(m.options, bson.E{Key: "whenNotMatched", Value: action})
	return m
}

// Let 设置变量
func (m *MergeStage[T]) Let(vars M) *MergeStage[T] {
	m.options = append(m.options, bson.E{Key: "let", Value: vars})
	return m
}

// Done 完成合并阶段构建
func (m *MergeStage[T]) Done() *Aggregate[T] {
	mergeDoc := bson.D{{"into", m.into}}
	mergeDoc = append(mergeDoc, m.options...)
	m.agg.pipeline = append(m.agg.pipeline, bson.D{{"$merge", mergeDoc}})
	return m.agg
}

// ========== FacetStage 构建器方法 ==========

// Facet 添加分面
func (f *FacetStage[T]) Facet(name string, stages ...any) *FacetStage[T] {
	f.facets = append(f.facets, bson.E{Key: name, Value: stages})
	return f
}

// Done 完成分面阶段构建
func (f *FacetStage[T]) Done() *Aggregate[T] {
	if len(f.facets) == 0 {
		return f.agg
	}
	f.agg.pipeline = append(f.agg.pipeline, bson.D{{"$facet", f.facets}})
	return f.agg
}

// ========== BucketStage 构建器方法 ==========

// GroupBy 设置分组字段
func (b *BucketStage[T]) GroupBy(field any) *BucketStage[T] {
	b.options = append(b.options, bson.E{Key: "groupBy", Value: field})
	return b
}

// Boundaries 设置边界
func (b *BucketStage[T]) Boundaries(boundaries ...any) *BucketStage[T] {
	b.options = append(b.options, bson.E{Key: "boundaries", Value: boundaries})
	return b
}

// Default 设置默认桶
func (b *BucketStage[T]) Default(defaultBucket any) *BucketStage[T] {
	b.options = append(b.options, bson.E{Key: "default", Value: defaultBucket})
	return b
}

// Output 设置输出
func (b *BucketStage[T]) Output(accumulators M) *BucketStage[T] {
	b.options = append(b.options, bson.E{Key: "output", Value: accumulators})
	return b
}

// Done 完成分桶阶段构建
func (b *BucketStage[T]) Done() *Aggregate[T] {
	if len(b.options) == 0 {
		return b.agg
	}
	b.agg.pipeline = append(b.agg.pipeline, bson.D{{"$bucket", b.options}})
	return b.agg
}

// ========== BucketAutoStage 构建器方法 ==========

// GroupBy 设置分组字段
func (b *BucketAutoStage[T]) GroupBy(field any) *BucketAutoStage[T] {
	b.options = append(b.options, bson.E{Key: "groupBy", Value: field})
	return b
}

// Buckets 设置桶数量
func (b *BucketAutoStage[T]) Buckets(n int) *BucketAutoStage[T] {
	b.options = append(b.options, bson.E{Key: "buckets", Value: n})
	return b
}

// Granularity 设置粒度
func (b *BucketAutoStage[T]) Granularity(granularity string) *BucketAutoStage[T] {
	b.options = append(b.options, bson.E{Key: "granularity", Value: granularity})
	return b
}

// Output 设置输出
func (b *BucketAutoStage[T]) Output(accumulators M) *BucketAutoStage[T] {
	b.options = append(b.options, bson.E{Key: "output", Value: accumulators})
	return b
}

// Done 完成自动分桶阶段构建
func (b *BucketAutoStage[T]) Done() *Aggregate[T] {
	if len(b.options) == 0 {
		return b.agg
	}
	b.agg.pipeline = append(b.agg.pipeline, bson.D{{"$bucketAuto", b.options}})
	return b.agg
}

// ========== 其他阶段构建器方法 ==========

// Done 完成窗口字段阶段构建
func (s *SetWindowFieldsStage[T]) Done() *Aggregate[T] {
	if len(s.options) == 0 {
		return s.agg
	}
	s.agg.pipeline = append(s.agg.pipeline, bson.D{{"$setWindowFields", s.options}})
	return s.agg
}

// Done 完成搜索阶段构建
func (s *SearchStage[T]) Done() *Aggregate[T] {
	if len(s.options) == 0 {
		return s.agg
	}
	s.agg.pipeline = append(s.agg.pipeline, bson.D{{"$search", s.options}})
	return s.agg
}

// Done 完成搜索元数据阶段构建
func (s *SearchMetaStage[T]) Done() *Aggregate[T] {
	if len(s.options) == 0 {
		return s.agg
	}
	s.agg.pipeline = append(s.agg.pipeline, bson.D{{"$searchMeta", s.options}})
	return s.agg
}

// Done 完成向量搜索阶段构建
func (v *VectorSearchStage[T]) Done() *Aggregate[T] {
	if len(v.options) == 0 {
		return v.agg
	}
	v.agg.pipeline = append(v.agg.pipeline, bson.D{{"$vectorSearch", v.options}})
	return v.agg
}

// Done 完成密度化阶段构建
func (d *DensifyStage[T]) Done() *Aggregate[T] {
	if len(d.options) == 0 {
		return d.agg
	}
	d.agg.pipeline = append(d.agg.pipeline, bson.D{{"$densify", d.options}})
	return d.agg
}

// Done 完成填充阶段构建
func (f *FillStage[T]) Done() *Aggregate[T] {
	if len(f.options) == 0 {
		return f.agg
	}
	f.agg.pipeline = append(f.agg.pipeline, bson.D{{"$fill", f.options}})
	return f.agg
}

// Done 完成联合阶段构建
func (u *UnionWithStage[T]) Done() *Aggregate[T] {
	if len(u.options) == 0 {
		return u.agg
	}
	u.agg.pipeline = append(u.agg.pipeline, bson.D{{"$unionWith", u.options}})
	return u.agg
}

// Done 完成图查找阶段构建
func (g *GraphLookupStage[T]) Done() *Aggregate[T] {
	if len(g.options) == 0 {
		return g.agg
	}
	g.agg.pipeline = append(g.agg.pipeline, bson.D{{"$graphLookup", g.options}})
	return g.agg
}

// Done 完成地理邻近阶段构建
func (g *GeoNearStage[T]) Done() *Aggregate[T] {
	if len(g.options) == 0 {
		return g.agg
	}
	g.agg.pipeline = append(g.agg.pipeline, bson.D{{"$geoNear", g.options}})
	return g.agg
}
