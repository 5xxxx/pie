package pie

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestSessionScopes 用于测试Session方法
type TestSessionScopes struct {
	query *Query
}

func (s *TestSessionScopes) Scopes(scopes ...ScopeFunc) *TestSessionScopes {
	for _, scope := range scopes {
		s.query = scope(s.query)
	}
	return s
}

func (s *TestSessionScopes) Latest(field string, limit int) *TestSessionScopes {
	s.query.OrderByDesc(field).Limit(int64(limit))
	return s
}

func (s *TestSessionScopes) Oldest(field string, limit int) *TestSessionScopes {
	s.query.OrderBy(field).Limit(int64(limit))
	return s
}

func (s *TestSessionScopes) Random(limit int) *TestSessionScopes {
	s.query.Limit(int64(limit))
	s.query.Sort(bson.D{{Key: "_random", Value: 0.5}})
	return s
}

func TestSessionScopesMethod(t *testing.T) {
	s := &TestSessionScopes{query: NewQuery()}
	
	// 测试多个作用域
	s.Scopes(
		ActiveScope("status"),
		NotDeletedScope("deleted_at"),
		PublishedScope("published"),
	)
	
	filter := s.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected conditions to be generated")
	}
	
	// 验证作用域被应用
	foundFields := make(map[string]bool)
	for _, cond := range filter {
		foundFields[cond.Key] = true
	}
	
	// 应该包含status字段（来自ActiveScope）
	if !foundFields["status"] {
		t.Error("Expected 'status' field from ActiveScope")
	}
}

func TestSessionLatest(t *testing.T) {
	s := &TestSessionScopes{query: NewQuery()}
	
	s.Latest("created_at", 10)
	
	sort := s.query.GetSort()
	if len(sort) == 0 {
		t.Fatal("Expected sort conditions")
	}
	
	// 验证排序
	found := false
	for _, sortCond := range sort {
		if sortCond.Key == "created_at" && sortCond.Value == -1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected descending sort on created_at")
	}
	
	// 验证限制
	limit := s.query.GetLimit()
	if limit == nil || *limit != 10 {
		t.Error("Expected limit 10")
	}
}

func TestSessionOldest(t *testing.T) {
	s := &TestSessionScopes{query: NewQuery()}
	
	s.Oldest("created_at", 5)
	
	sort := s.query.GetSort()
	if len(sort) == 0 {
		t.Fatal("Expected sort conditions")
	}
	
	// 验证排序
	found := false
	for _, sortCond := range sort {
		if sortCond.Key == "created_at" && sortCond.Value == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ascending sort on created_at")
	}
	
	// 验证限制
	limit := s.query.GetLimit()
	if limit == nil || *limit != 5 {
		t.Error("Expected limit 5")
	}
}

func TestSessionRandom(t *testing.T) {
	s := &TestSessionScopes{query: NewQuery()}
	
	s.Random(3)
	
	// 验证限制
	limit := s.query.GetLimit()
	if limit == nil || *limit != 3 {
		t.Error("Expected limit 3")
	}
	
	// 验证随机排序
	sort := s.query.GetSort()
	if len(sort) == 0 {
		t.Fatal("Expected sort conditions")
	}
	
	found := false
	for _, sortCond := range sort {
		if sortCond.Key == "_random" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected random sort condition")
	}
}

func TestActiveScope(t *testing.T) {
	q := NewQuery()
	
	scope := ActiveScope("status")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证status字段
	found := false
	for _, cond := range filter {
		if cond.Key == "status" && cond.Value == "active" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected status='active' condition")
	}
}

func TestNotDeletedScope(t *testing.T) {
	q := NewQuery()
	
	scope := NotDeletedScope("deleted_at")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证$or条件
	found := false
	for _, cond := range filter {
		if cond.Key == "$or" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected $or condition for null check")
	}
}

func TestPublishedScope(t *testing.T) {
	q := NewQuery()
	
	scope := PublishedScope("published")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证published字段
	found := false
	for _, cond := range filter {
		if cond.Key == "published" && cond.Value == true {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected published=true condition")
	}
}

func TestVisibleScope(t *testing.T) {
	q := NewQuery()
	
	// 测试默认值
	scope := VisibleScope("visibility")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证$in条件
	found := false
	for _, cond := range filter {
		if cond.Key == "visibility" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$in" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected visibility $in condition")
	}
	
	// 测试自定义值
	q = NewQuery()
	scope = VisibleScope("visibility", "public", "internal")
	result = scope(q)
	
	filter = q.GetFilter()
	found = false
	for _, cond := range filter {
		if cond.Key == "visibility" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$in" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected visibility $in condition with custom values")
	}
}

func TestDateRangeScope(t *testing.T) {
	q := NewQuery()
	
	start := "2023-01-01"
	end := "2023-12-31"
	
	scope := DateRangeScope("created_at", start, end)
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证范围条件
	found := false
	for _, cond := range filter {
		if cond.Key == "created_at" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) >= 2 {
				hasGte := false
				hasLte := false
				for _, v := range valueDoc {
					if v.Key == "$gte" && v.Value == start {
						hasGte = true
					}
					if v.Key == "$lte" && v.Value == end {
						hasLte = true
					}
				}
				if hasGte && hasLte {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Error("Expected date range condition")
	}
}

func TestSearchScope(t *testing.T) {
	q := NewQuery()
	
	// 测试正常搜索
	scope := SearchScope("test", "name", "description", "content")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证$or条件
	found := false
	for _, cond := range filter {
		if cond.Key == "$or" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected $or condition for search")
	}
	
	// 测试空关键词
	q = NewQuery()
	scope = SearchScope("", "name", "description")
	result = scope(q)
	
	filter = q.GetFilter()
	if len(filter) != 0 {
		t.Error("Expected no conditions for empty keyword")
	}
	
	// 测试空字段
	q = NewQuery()
	scope = SearchScope("test")
	result = scope(q)
	
	filter = q.GetFilter()
	if len(filter) != 0 {
		t.Error("Expected no conditions for empty fields")
	}
}

func TestPriceRangeScope(t *testing.T) {
	q := NewQuery()
	
	// 测试只有最小价格
	minPrice := 10.0
	scope := PriceRangeScope(&minPrice, nil)
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证$gte条件
	found := false
	for _, cond := range filter {
		if cond.Key == "price" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$gte" && valueDoc[0].Value == minPrice {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected price $gte condition")
	}
	
	// 测试只有最大价格
	q = NewQuery()
	maxPrice := 100.0
	scope = PriceRangeScope(nil, &maxPrice)
	result = scope(q)
	
	filter = q.GetFilter()
	found = false
	for _, cond := range filter {
		if cond.Key == "price" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$lte" && valueDoc[0].Value == maxPrice {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected price $lte condition")
	}
	
	// 测试价格范围
	q = NewQuery()
	scope = PriceRangeScope(&minPrice, &maxPrice)
	result = scope(q)
	
	filter = q.GetFilter()
	if len(filter) != 2 {
		t.Error("Expected 2 price conditions")
	}
}

func TestCategoryScope(t *testing.T) {
	q := NewQuery()
	
	// 测试单个分类
	scope := CategoryScope("electronics")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证分类条件
	found := false
	for _, cond := range filter {
		if cond.Key == "category" && cond.Value == "electronics" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected category condition")
	}
	
	// 测试多个分类
	q = NewQuery()
	scope = CategoryScope("electronics", "books", "clothing")
	result = scope(q)
	
	filter = q.GetFilter()
	found = false
	for _, cond := range filter {
		if cond.Key == "category" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$in" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected category $in condition")
	}
	
	// 测试空分类
	q = NewQuery()
	scope = CategoryScope()
	result = scope(q)
	
	filter = q.GetFilter()
	if len(filter) != 0 {
		t.Error("Expected no conditions for empty categories")
	}
}

func TestTagsScope(t *testing.T) {
	q := NewQuery()
	
	// 测试标签
	scope := TagsScope("golang", "mongodb", "testing")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证标签条件
	found := false
	for _, cond := range filter {
		if cond.Key == "tags" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$in" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected tags $in condition")
	}
	
	// 测试空标签
	q = NewQuery()
	scope = TagsScope()
	result = scope(q)
	
	filter = q.GetFilter()
	if len(filter) != 0 {
		t.Error("Expected no conditions for empty tags")
	}
}

func TestAllTagsScope(t *testing.T) {
	q := NewQuery()
	
	// 测试所有标签
	scope := AllTagsScope("golang", "mongodb", "testing")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证标签条件
	found := false
	for _, cond := range filter {
		if cond.Key == "tags" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$all" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected tags $all condition")
	}
	
	// 测试空标签
	q = NewQuery()
	scope = AllTagsScope()
	result = scope(q)
	
	filter = q.GetFilter()
	if len(filter) != 0 {
		t.Error("Expected no conditions for empty tags")
	}
}

func TestUserScope(t *testing.T) {
	q := NewQuery()
	
	scope := UserScope("user123")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证用户条件
	found := false
	for _, cond := range filter {
		if cond.Key == "user_id" && cond.Value == "user123" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected user_id condition")
	}
}

func TestTenantScope(t *testing.T) {
	q := NewQuery()
	
	scope := TenantScope("tenant456")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证租户条件
	found := false
	for _, cond := range filter {
		if cond.Key == "tenant_id" && cond.Value == "tenant456" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected tenant_id condition")
	}
}

func TestFeaturedScope(t *testing.T) {
	q := NewQuery()
	
	scope := FeaturedScope()
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	sort := q.GetSort()
	
	// 验证featured条件
	found := false
	for _, cond := range filter {
		if cond.Key == "featured" && cond.Value == true {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected featured condition")
	}
	
	// 验证优先级排序
	found = false
	for _, sortCond := range sort {
		if sortCond.Key == "priority" && sortCond.Value == -1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected priority descending sort")
	}
}

func TestPopularScope(t *testing.T) {
	q := NewQuery()
	
	scope := PopularScope("view_count", 20)
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	sort := q.GetSort()
	limit := q.GetLimit()
	
	// 验证排序
	found := false
	for _, sortCond := range sort {
		if sortCond.Key == "view_count" && sortCond.Value == -1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected view_count descending sort")
	}
	
	// 验证限制
	if limit == nil || *limit != 20 {
		t.Error("Expected limit 20")
	}
}

func TestRecentScope(t *testing.T) {
	q := NewQuery()
	
	scope := RecentScope("created_at", 7)
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证最近天数条件
	found := false
	for _, cond := range filter {
		if cond.Key == "created_at" {
			valueDoc, ok := cond.Value.(bson.D)
			if ok && len(valueDoc) > 0 && valueDoc[0].Key == "$gte" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected created_at $gte condition")
	}
}

func TestVerifiedScope(t *testing.T) {
	q := NewQuery()
	
	scope := VerifiedScope("verified")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证verified条件
	found := false
	for _, cond := range filter {
		if cond.Key == "verified" && cond.Value == true {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected verified condition")
	}
}

func TestApprovedScope(t *testing.T) {
	q := NewQuery()
	
	// 测试默认状态
	scope := ApprovedScope("")
	result := scope(q)
	
	if result != q {
		t.Error("Scope should return the same query instance")
	}
	
	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter conditions")
	}
	
	// 验证approval_status条件
	found := false
	for _, cond := range filter {
		if cond.Key == "approval_status" && cond.Value == "approved" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected approval_status='approved' condition")
	}
	
	// 测试自定义状态
	q = NewQuery()
	scope = ApprovedScope("pending")
	result = scope(q)
	
	filter = q.GetFilter()
	found = false
	for _, cond := range filter {
		if cond.Key == "approval_status" && cond.Value == "pending" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected approval_status='pending' condition")
	}
}

func TestScopeChaining(t *testing.T) {
	s := &TestSessionScopes{query: NewQuery()}
	
	// 测试作用域链式调用
	result := s.Scopes(
		ActiveScope("status"),
		NotDeletedScope("deleted_at"),
		PublishedScope("published"),
	).Latest("created_at", 10)
	
	if result != s {
		t.Error("Chaining should return the same instance")
	}
	
	// 验证所有条件都被应用
	filter := s.query.GetFilter()
	sort := s.query.GetSort()
	limit := s.query.GetLimit()
	
	if len(filter) == 0 {
		t.Error("Expected filter conditions from scopes")
	}
	if len(sort) == 0 {
		t.Error("Expected sort conditions from Latest")
	}
	if limit == nil {
		t.Error("Expected limit from Latest")
	}
}
