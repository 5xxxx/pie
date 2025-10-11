package pie

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestNewQuery(t *testing.T) {
	q := NewQuery()
	if q == nil {
		t.Fatal("NewQuery should not return nil")
	}
	
	// 验证初始状态
	if len(q.GetFilter()) != 0 {
		t.Error("New query should have empty filter")
	}
	if len(q.GetSort()) != 0 {
		t.Error("New query should have empty sort")
	}
	if q.GetLimit() != nil {
		t.Error("New query should have nil limit")
	}
	if q.GetSkip() != nil {
		t.Error("New query should have nil skip")
	}
	if len(q.GetProject()) != 0 {
		t.Error("New query should have empty project")
	}
}

func TestQueryWhere(t *testing.T) {
	q := NewQuery()
	
	// 测试基本Where条件
	q.Where("name", "test")
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	expected := bson.E{Key: "name", Value: "test"}
	if filter[0] != expected {
		t.Errorf("Expected %v, got %v", expected, filter[0])
	}
	
	// 测试多个Where条件
	q.Where("age", 25)
	filter = q.GetFilter()
	if len(filter) != 2 {
		t.Fatalf("Expected 2 filter conditions, got %d", len(filter))
	}
	
	expected2 := bson.E{Key: "age", Value: 25}
	if filter[1] != expected2 {
		t.Errorf("Expected %v, got %v", expected2, filter[1])
	}
}

func TestQueryWhereOperator(t *testing.T) {
	q := NewQuery()
	
	// 测试WhereOperator
	op := Eq("status", "active")
	q.WhereOperator(op)
	
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	if filter[0].Key != "status" {
		t.Errorf("Expected key 'status', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$eq" || valueDoc[0].Value != "active" {
			t.Errorf("Expected $eq: 'active', got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestQueryAnd(t *testing.T) {
	q := NewQuery()
	
	// 测试And条件
	q.And(Eq("active", true), Gte("age", 18))
	
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	if filter[0].Key != "$and" {
		t.Errorf("Expected $and key, got %s", filter[0].Key)
	}
	
	andValue, ok := filter[0].Value.([]bson.D)
	if !ok {
		t.Fatal("Expected $and value to be []bson.D")
	}
	
	if len(andValue) != 2 {
		t.Fatalf("Expected 2 AND conditions, got %d", len(andValue))
	}
}

func TestQueryOr(t *testing.T) {
	q := NewQuery()
	
	// 测试Or条件
	q.Or(Eq("name", "Alice"), Eq("name", "Bob"))
	
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	if filter[0].Key != "$or" {
		t.Errorf("Expected $or key, got %s", filter[0].Key)
	}
	
	orValue, ok := filter[0].Value.([]bson.D)
	if !ok {
		t.Fatal("Expected $or value to be []bson.D")
	}
	
	if len(orValue) != 2 {
		t.Fatalf("Expected 2 OR conditions, got %d", len(orValue))
	}
}

func TestQueryNor(t *testing.T) {
	q := NewQuery()
	
	// 测试Nor条件
	q.Nor(Eq("deleted", true), Eq("archived", true))
	
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	if filter[0].Key != "$nor" {
		t.Errorf("Expected $nor key, got %s", filter[0].Key)
	}
	
	norValue, ok := filter[0].Value.([]bson.D)
	if !ok {
		t.Fatal("Expected $nor value to be []bson.D")
	}
	
	if len(norValue) != 2 {
		t.Fatalf("Expected 2 NOR conditions, got %d", len(norValue))
	}
}

func TestQueryOrderBy(t *testing.T) {
	q := NewQuery()
	
	// 测试OrderBy
	q.OrderBy("name")
	sort := q.GetSort()
	if len(sort) != 1 {
		t.Fatalf("Expected 1 sort condition, got %d", len(sort))
	}
	
	expected := bson.E{Key: "name", Value: 1}
	if sort[0] != expected {
		t.Errorf("Expected %v, got %v", expected, sort[0])
	}
	
	// 测试多个OrderBy
	q.OrderBy("age")
	sort = q.GetSort()
	if len(sort) != 2 {
		t.Fatalf("Expected 2 sort conditions, got %d", len(sort))
	}
}

func TestQueryOrderByDesc(t *testing.T) {
	q := NewQuery()
	
	// 测试OrderByDesc
	q.OrderByDesc("created_at")
	sort := q.GetSort()
	if len(sort) != 1 {
		t.Fatalf("Expected 1 sort condition, got %d", len(sort))
	}
	
	expected := bson.E{Key: "created_at", Value: -1}
	if sort[0] != expected {
		t.Errorf("Expected %v, got %v", expected, sort[0])
	}
}

func TestQuerySort(t *testing.T) {
	q := NewQuery()
	
	// 测试Sort
	sortDoc := bson.D{
		{Key: "priority", Value: -1},
		{Key: "created_at", Value: 1},
	}
	q.Sort(sortDoc)
	
	sort := q.GetSort()
	if len(sort) != 2 {
		t.Fatalf("Expected 2 sort conditions, got %d", len(sort))
	}
	
	expected1 := bson.E{Key: "priority", Value: -1}
	expected2 := bson.E{Key: "created_at", Value: 1}
	
	if sort[0] != expected1 {
		t.Errorf("Expected %v, got %v", expected1, sort[0])
	}
	if sort[1] != expected2 {
		t.Errorf("Expected %v, got %v", expected2, sort[1])
	}
}

func TestQueryLimit(t *testing.T) {
	q := NewQuery()
	
	// 测试Limit
	q.Limit(10)
	limit := q.GetLimit()
	if limit == nil {
		t.Fatal("Expected limit to be set")
	}
	if *limit != 10 {
		t.Errorf("Expected limit 10, got %d", *limit)
	}
}

func TestQuerySkip(t *testing.T) {
	q := NewQuery()
	
	// 测试Skip
	q.Skip(5)
	skip := q.GetSkip()
	if skip == nil {
		t.Fatal("Expected skip to be set")
	}
	if *skip != 5 {
		t.Errorf("Expected skip 5, got %d", *skip)
	}
}

func TestQueryProject(t *testing.T) {
	q := NewQuery()
	
	// 测试Project
	projectDoc := bson.D{
		{Key: "name", Value: 1},
		{Key: "email", Value: 1},
	}
	q.Project(projectDoc)
	
	project := q.GetProject()
	if len(project) != 2 {
		t.Fatalf("Expected 2 project fields, got %d", len(project))
	}
	
	expected1 := bson.E{Key: "name", Value: 1}
	expected2 := bson.E{Key: "email", Value: 1}
	
	if project[0] != expected1 {
		t.Errorf("Expected %v, got %v", expected1, project[0])
	}
	if project[1] != expected2 {
		t.Errorf("Expected %v, got %v", expected2, project[1])
	}
}

func TestQuerySelect(t *testing.T) {
	q := NewQuery()
	
	// 测试Select
	q.Select("name", "email", "age")
	
	project := q.GetProject()
	if len(project) != 3 {
		t.Fatalf("Expected 3 project fields, got %d", len(project))
	}
	
	expected := []bson.E{
		{Key: "name", Value: 1},
		{Key: "email", Value: 1},
		{Key: "age", Value: 1},
	}
	
	for i, exp := range expected {
		if project[i] != exp {
			t.Errorf("Expected %v, got %v", exp, project[i])
		}
	}
}

func TestQueryExclude(t *testing.T) {
	q := NewQuery()
	
	// 测试Exclude
	q.Exclude("password", "secret")
	
	project := q.GetProject()
	if len(project) != 2 {
		t.Fatalf("Expected 2 project fields, got %d", len(project))
	}
	
	expected := []bson.E{
		{Key: "password", Value: 0},
		{Key: "secret", Value: 0},
	}
	
	for i, exp := range expected {
		if project[i] != exp {
			t.Errorf("Expected %v, got %v", exp, project[i])
		}
	}
}

func TestQueryClone(t *testing.T) {
	q := NewQuery()
	q.Where("name", "test")
	q.OrderBy("age")
	q.Limit(10)
	q.Skip(5)
	q.Select("name", "email")
	
	cloned := q.Clone()
	if cloned == q {
		t.Error("Clone should return a different instance")
	}
	
	// 验证克隆的内容
	if len(cloned.GetFilter()) != len(q.GetFilter()) {
		t.Error("Cloned filter should have same length")
	}
	if len(cloned.GetSort()) != len(q.GetSort()) {
		t.Error("Cloned sort should have same length")
	}
	if cloned.GetLimit() == nil || *cloned.GetLimit() != *q.GetLimit() {
		t.Error("Cloned limit should be same")
	}
	if cloned.GetSkip() == nil || *cloned.GetSkip() != *q.GetSkip() {
		t.Error("Cloned skip should be same")
	}
	if len(cloned.GetProject()) != len(q.GetProject()) {
		t.Error("Cloned project should have same length")
	}
	
	// 验证修改克隆不影响原对象
	cloned.Where("status", "active")
	if len(q.GetFilter()) != 1 {
		t.Error("Original query should not be affected by clone modification")
	}
}

func TestQueryClear(t *testing.T) {
	q := NewQuery()
	q.Where("name", "test")
	q.OrderBy("age")
	q.Limit(10)
	q.Skip(5)
	q.Select("name", "email")
	
	cleared := q.Clear()
	if cleared != q {
		t.Error("Clear should return the same instance")
	}
	
	// 验证清空
	if len(q.GetFilter()) != 0 {
		t.Error("Filter should be empty after clear")
	}
	if len(q.GetSort()) != 0 {
		t.Error("Sort should be empty after clear")
	}
	if q.GetLimit() != nil {
		t.Error("Limit should be nil after clear")
	}
	if q.GetSkip() != nil {
		t.Error("Skip should be nil after clear")
	}
	if len(q.GetProject()) != 0 {
		t.Error("Project should be empty after clear")
	}
}

func TestQueryBuildFindOptions(t *testing.T) {
	q := NewQuery()
	q.Where("name", "test")
	q.OrderBy("age")
	q.Limit(10)
	q.Skip(5)
	q.Select("name", "email")
	
	opts := q.BuildFindOptions()
	if opts == nil {
		t.Fatal("BuildFindOptions should not return nil")
	}
	
	// 验证选项构建器是否正确设置
	// 注意：这里我们主要测试方法调用不报错，具体选项验证需要更复杂的测试
}

func TestQueryBuildFindOneOptions(t *testing.T) {
	q := NewQuery()
	q.Where("name", "test")
	q.OrderBy("age")
	q.Skip(5)
	q.Select("name", "email")
	
	opts := q.BuildFindOneOptions()
	if opts == nil {
		t.Fatal("BuildFindOneOptions should not return nil")
	}
	
	// 验证选项构建器是否正确设置
}

func TestQueryBuild(t *testing.T) {
	q := NewQuery()
	q.Where("name", "test")
	q.Where("age", 25)
	
	filter := q.Build()
	if len(filter) != 2 {
		t.Fatalf("Expected 2 filter conditions, got %d", len(filter))
	}
	
	expected1 := bson.E{Key: "name", Value: "test"}
	expected2 := bson.E{Key: "age", Value: 25}
	
	if filter[0] != expected1 {
		t.Errorf("Expected %v, got %v", expected1, filter[0])
	}
	if filter[1] != expected2 {
		t.Errorf("Expected %v, got %v", expected2, filter[1])
	}
}

func TestQueryWhereArrayAll(t *testing.T) {
	q := NewQuery()
	
	// 测试WhereArrayAll
	values := []string{"tag1", "tag2", "tag3"}
	q.WhereArrayAll("tags", values)
	
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	if filter[0].Key != "tags" {
		t.Errorf("Expected key 'tags', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$all" {
			t.Errorf("Expected $all operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestQueryWhereRecentDays(t *testing.T) {
	q := NewQuery()
	
	// 测试WhereRecentDays
	q.WhereRecentDays("created_at", 7)
	
	filter := q.GetFilter()
	if len(filter) != 1 {
		t.Fatalf("Expected 1 filter condition, got %d", len(filter))
	}
	
	if filter[0].Key != "created_at" {
		t.Errorf("Expected field 'created_at', got %s", filter[0].Key)
	}
	
	// 验证值是一个时间范围
	valueDoc, ok := filter[0].Value.(bson.D)
	if !ok {
		t.Fatal("Expected value to be bson.D")
	}
	
	if len(valueDoc) != 1 {
		t.Fatalf("Expected 1 condition in value, got %d", len(valueDoc))
	}
	
	if valueDoc[0].Key != "$gte" {
		t.Errorf("Expected $gte operator, got %s", valueDoc[0].Key)
	}
	
	// 验证时间值
	timeValue, ok := valueDoc[0].Value.(time.Time)
	if !ok {
		t.Fatal("Expected time.Time value")
	}
	
	// 验证时间在合理范围内（7天前到现在）
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)
	if timeValue.Before(sevenDaysAgo.Add(-time.Hour)) || timeValue.After(now) {
		t.Errorf("Time value %v is not within expected range", timeValue)
	}
}

func TestQueryEmptyOperators(t *testing.T) {
	q := NewQuery()
	
	// 测试空操作符
	q.And() // 空And
	q.Or()  // 空Or
	q.Nor() // 空Nor
	
	filter := q.GetFilter()
	if len(filter) != 0 {
		t.Errorf("Expected no filter conditions with empty operators, got %d", len(filter))
	}
}

func TestQueryChaining(t *testing.T) {
	q := NewQuery()
	
	// 测试方法链式调用
	result := q.Where("name", "test").
		Where("age", 25).
		OrderBy("created_at").
		OrderByDesc("priority").
		Limit(10).
		Skip(5).
		Select("name", "email")
	
	if result != q {
		t.Error("Chaining should return the same instance")
	}
	
	// 验证所有条件都被正确设置
	if len(q.GetFilter()) != 2 {
		t.Error("Expected 2 filter conditions")
	}
	if len(q.GetSort()) != 2 {
		t.Error("Expected 2 sort conditions")
	}
	if q.GetLimit() == nil || *q.GetLimit() != 10 {
		t.Error("Expected limit 10")
	}
	if q.GetSkip() == nil || *q.GetSkip() != 5 {
		t.Error("Expected skip 5")
	}
	if len(q.GetProject()) != 2 {
		t.Error("Expected 2 project fields")
	}
}
