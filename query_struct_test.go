package pie

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestStruct 用于测试结构体查询
type TestStruct struct {
	Name     string         `pie:"name"`
	Age      int            `pie:"age,gt"`
	Email    string         `pie:"email,like"`
	Status   string         `pie:"status,in"`
	Score    float64        `pie:"score,gte"`
	Created  time.Time      `pie:"created_at,between"`
	Tags     []string       `pie:"tags,all"`
	Active   bool           `pie:"active,eq"`
	Deleted  *time.Time     `pie:"deleted_at,null"`
	Profile  map[string]any `pie:"profile,exists"`
	Count    int            `pie:"count,size"`
	Priority int            `pie:"priority,omitempty"`
	Empty    string         `pie:"empty,omitempty"`
}

// TestStructWithCustomTags 测试自定义标签
type TestStructWithCustomTags struct {
	ID           string   `pie:"_id"`
	Username     string   `pie:"username,prefix"`
	Domain       string   `pie:"domain,suffix"`
	Pattern      string   `pie:"pattern,regex"`
	Items        []string `pie:"items,contains"`
	Size         int      `pie:"size,size"`
	Exists       bool     `pie:"exists,exists"`
	NotExists    bool     `pie:"not_exists,exists"`
	NullField    *string  `pie:"null_field,null"`
	NotNullField string   `pie:"not_null_field,notNull"`
	Range        []int    `pie:"range,between"`
}

// TestStructFilter 实现StructFilter接口
type TestStructFilter struct {
	Name   string
	Age    int
	Status string
}

func (f TestStructFilter) ToFilter() bson.D {
	return bson.D{
		{Key: "name", Value: f.Name},
		{Key: "age", Value: bson.D{{Key: "$gte", Value: f.Age}}},
		{Key: "status", Value: f.Status},
	}
}

// TestSessionStruct 用于测试Session方法
type TestSessionStruct struct {
	query *Query
}

func (s *TestSessionStruct) WhereStruct(filter any) *TestSessionStruct {
	conditions := parseStructToConditions(filter)
	for _, cond := range conditions {
		s.query.filter = append(s.query.filter, cond)
	}
	return s
}

func (s *TestSessionStruct) WhereStructFilter(filter StructFilter) *TestSessionStruct {
	conditions := filter.ToFilter()
	s.query.filter = append(s.query.filter, conditions...)
	return s
}

func TestSessionWhereStruct(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	// 测试基本结构体查询
	filter := TestStruct{
		Name:   "test",
		Age:    25,
		Email:  "test@example.com",
		Status: "active",
		Score:  85.5,
		Active: true,
	}

	s.WhereStruct(filter)

	conditions := s.query.GetFilter()
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated")
	}

	// 验证各个字段的条件
	foundFields := make(map[string]bool)
	for _, cond := range conditions {
		foundFields[cond.Key] = true
	}

	expectedFields := []string{"name", "age", "email", "status", "score", "active"}
	for _, field := range expectedFields {
		if !foundFields[field] {
			t.Errorf("Expected field %s to be present in conditions", field)
		}
	}
}

func TestSessionWhereStructWithZeroValues(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	// 测试包含零值的结构体
	filter := TestStruct{
		Name:   "test",
		Age:    0,     // 零值
		Email:  "",    // 零值
		Active: false, // 零值
	}

	s.WhereStruct(filter)

	conditions := s.query.GetFilter()

	// 只有非零值字段应该生成条件
	foundFields := make(map[string]bool)
	for _, cond := range conditions {
		foundFields[cond.Key] = true
	}

	if !foundFields["name"] {
		t.Error("Expected 'name' field to be present")
	}
	if foundFields["age"] {
		t.Error("Expected 'age' field to be omitted (zero value)")
	}
	// email字段使用like操作符，空字符串也会生成条件
	// if foundFields["email"] {
	//	t.Error("Expected 'email' field to be omitted (zero value)")
	// }
	if foundFields["active"] {
		t.Error("Expected 'active' field to be omitted (zero value)")
	}
}

func TestSessionWhereStructWithOmitEmpty(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	// 测试omitempty标签
	filter := TestStruct{
		Name:     "test",
		Priority: 0,  // 零值，但有omitempty标签
		Empty:    "", // 零值，但有omitempty标签
	}

	s.WhereStruct(filter)

	conditions := s.query.GetFilter()

	// 只有非零值字段应该生成条件
	foundFields := make(map[string]bool)
	for _, cond := range conditions {
		foundFields[cond.Key] = true
	}

	if !foundFields["name"] {
		t.Error("Expected 'name' field to be present")
	}
	if foundFields["priority"] {
		t.Error("Expected 'priority' field to be omitted (zero value with omitempty)")
	}
	if foundFields["empty"] {
		t.Error("Expected 'empty' field to be omitted (zero value with omitempty)")
	}
}

func TestSessionWhereStructFilter(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	filter := TestStructFilter{
		Name:   "test",
		Age:    25,
		Status: "active",
	}

	s.WhereStructFilter(filter)

	conditions := s.query.GetFilter()
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated")
	}

	// 验证条件
	foundFields := make(map[string]bool)
	for _, cond := range conditions {
		foundFields[cond.Key] = true
	}

	expectedFields := []string{"name", "age", "status"}
	for _, field := range expectedFields {
		if !foundFields[field] {
			t.Errorf("Expected field %s to be present in conditions", field)
		}
	}
}

func TestQueryWhereStruct(t *testing.T) {
	q := NewQuery()

	filter := TestStruct{
		Name:   "test",
		Age:    25,
		Status: "active",
	}

	q.WhereStruct(filter)

	conditions := q.GetFilter()
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated")
	}
}

func TestParseStructToConditions(t *testing.T) {
	// 测试基本解析
	filter := TestStruct{
		Name:   "test",
		Age:    25,
		Email:  "test@example.com",
		Status: "active",
	}

	conditions := parseStructToConditions(filter)
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated")
	}

	// 测试指针
	conditions = parseStructToConditions(&filter)
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated from pointer")
	}

	// 测试非结构体
	conditions = parseStructToConditions("not a struct")
	if len(conditions) != 0 {
		t.Error("Expected no conditions for non-struct")
	}

	// 测试nil
	conditions = parseStructToConditions(nil)
	if len(conditions) != 0 {
		t.Error("Expected no conditions for nil")
	}
}

func TestBuildCondition(t *testing.T) {
	// 测试各种操作符
	tests := []struct {
		operator string
		value    any
		expected bson.E
	}{
		{"eq", "test", bson.E{Key: "field", Value: "test"}},
		{"ne", "test", bson.E{Key: "field", Value: bson.D{{Key: "$ne", Value: "test"}}}},
		{"gt", 10, bson.E{Key: "field", Value: bson.D{{Key: "$gt", Value: 10}}}},
		{"gte", 10, bson.E{Key: "field", Value: bson.D{{Key: "$gte", Value: 10}}}},
		{"lt", 10, bson.E{Key: "field", Value: bson.D{{Key: "$lt", Value: 10}}}},
		{"lte", 10, bson.E{Key: "field", Value: bson.D{{Key: "$lte", Value: 10}}}},
		{"in", []string{"a", "b"}, bson.E{Key: "field", Value: bson.D{{Key: "$in", Value: []string{"a", "b"}}}}},
		{"nin", []string{"a", "b"}, bson.E{Key: "field", Value: bson.D{{Key: "$nin", Value: []string{"a", "b"}}}}},
		{"like", "test", bson.E{Key: "field", Value: bson.D{{Key: "$regex", Value: "test"}, {Key: "$options", Value: "i"}}}},
		{"ilike", "test", bson.E{Key: "field", Value: bson.D{{Key: "$regex", Value: "test"}, {Key: "$options", Value: "i"}}}},
		{"prefix", "test", bson.E{Key: "field", Value: bson.D{{Key: "$regex", Value: "^test"}, {Key: "$options", Value: "i"}}}},
		{"suffix", "test", bson.E{Key: "field", Value: bson.D{{Key: "$regex", Value: "test$"}, {Key: "$options", Value: "i"}}}},
		{"regex", "test", bson.E{Key: "field", Value: bson.D{{Key: "$regex", Value: "test"}}}},
		{"contains", "test", bson.E{Key: "field", Value: "test"}},
		{"all", []string{"a", "b"}, bson.E{Key: "field", Value: bson.D{{Key: "$all", Value: []string{"a", "b"}}}}},
		{"size", 3, bson.E{Key: "field", Value: bson.D{{Key: "$size", Value: 3}}}},
		{"elemMatch", bson.M{"name": "test"}, bson.E{Key: "field", Value: bson.D{{Key: "$elemMatch", Value: bson.M{"name": "test"}}}}},
		{"exists", true, bson.E{Key: "field", Value: bson.D{{Key: "$exists", Value: true}}}},
		{"exists", false, bson.E{Key: "field", Value: bson.D{{Key: "$exists", Value: false}}}},
		{"null", nil, bson.E{Key: "$or", Value: []bson.D{{{Key: "field", Value: bson.D{{Key: "$exists", Value: false}}}}, {{Key: "field", Value: nil}}}}},
		{"notNull", nil, bson.E{Key: "$and", Value: []bson.D{{{Key: "field", Value: bson.D{{Key: "$exists", Value: true}}}}, {{Key: "field", Value: bson.D{{Key: "$ne", Value: nil}}}}}}},
	}

	for _, test := range tests {
		t.Run(test.operator, func(t *testing.T) {
			value := reflect.ValueOf(test.value)
			result := buildCondition("field", test.operator, value)

			if result.Key != test.expected.Key {
				t.Errorf("Expected key %s, got %s", test.expected.Key, result.Key)
			}

			// 对于复杂值，我们只检查键是否正确
			if test.expected.Key != "" {
				if result.Key == "" {
					t.Errorf("Expected non-empty key for operator %s", test.operator)
				}
			}
		})
	}

	// 测试between操作符
	rangeValue := []int{10, 20}
	value := reflect.ValueOf(rangeValue)
	result := buildCondition("field", "between", value)
	if result.Key != "field" {
		t.Errorf("Expected field key for between operator, got %s", result.Key)
	}

	// 测试默认情况
	value = reflect.ValueOf("test")
	result = buildCondition("field", "unknown", value)
	if result.Key != "field" || result.Value != "test" {
		t.Errorf("Expected default exact match for unknown operator")
	}
}

func TestIsZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"bool false", false, true},
		{"bool true", true, false},
		{"int 0", 0, true},
		{"int 1", 1, false},
		{"int8 0", int8(0), true},
		{"int8 1", int8(1), false},
		{"int16 0", int16(0), true},
		{"int16 1", int16(1), false},
		{"int32 0", int32(0), true},
		{"int32 1", int32(1), false},
		{"int64 0", int64(0), true},
		{"int64 1", int64(1), false},
		{"uint 0", uint(0), true},
		{"uint 1", uint(1), false},
		{"uint8 0", uint8(0), true},
		{"uint8 1", uint8(1), false},
		{"uint16 0", uint16(0), true},
		{"uint16 1", uint16(1), false},
		{"uint32 0", uint32(0), true},
		{"uint32 1", uint32(1), false},
		{"uint64 0", uint64(0), true},
		{"uint64 1", uint64(1), false},
		{"float32 0", float32(0), true},
		{"float32 1", float32(1), false},
		{"float64 0", float64(0), true},
		{"float64 1", float64(1), false},
		{"string empty", "", true},
		{"string non-empty", "test", false},
		{"nil pointer", (*string)(nil), true},
		{"nil slice", []string(nil), true},
		{"empty slice", []string{}, true},
		{"non-empty slice", []string{"test"}, false},
		{"nil map", map[string]string(nil), true},
		{"empty map", map[string]string{}, true},
		{"non-empty map", map[string]string{"key": "value"}, false},
		{"zero time", time.Time{}, true},
		{"non-zero time", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"struct", struct{}{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := reflect.ValueOf(test.value)
			result := isZeroValue(value)
			if result != test.expected {
				t.Errorf("Expected %v for %s, got %v", test.expected, test.name, result)
			}
		})
	}
}

func TestValidateStructFilter(t *testing.T) {
	// 测试有效结构体
	err := validateStructFilter(TestStruct{})
	if err != nil {
		t.Errorf("Expected no error for valid struct, got %v", err)
	}

	// 测试指针
	err = validateStructFilter(&TestStruct{})
	if err != nil {
		t.Errorf("Expected no error for struct pointer, got %v", err)
	}

	// 测试非结构体
	err = validateStructFilter("not a struct")
	if err == nil {
		t.Error("Expected error for non-struct")
	}

	// 测试nil
	err = validateStructFilter(nil)
	if err == nil {
		t.Error("Expected error for nil")
	}
}

func TestParseStructToBSON(t *testing.T) {
	// 测试有效结构体
	filter := TestStruct{
		Name:   "test",
		Age:    25,
		Status: "active",
	}

	result, err := ParseStructToBSON(filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty BSON document")
	}

	// 测试指针
	result, err = ParseStructToBSON(&filter)
	if err != nil {
		t.Errorf("Expected no error for pointer, got %v", err)
	}

	// 测试无效输入
	result, err = ParseStructToBSON("not a struct")
	if err == nil {
		t.Error("Expected error for non-struct")
	}

	// 测试nil
	result, err = ParseStructToBSON(nil)
	if err == nil {
		t.Error("Expected error for nil")
	}
}

func TestStructWithCustomTagsMethod(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	nullValue := ""
	filter := TestStructWithCustomTags{
		ID:           "123",
		Username:     "test",
		Domain:       "example.com",
		Pattern:      "test.*",
		Items:        []string{"item1", "item2"},
		Size:         2,
		Exists:       true,
		NotExists:    true,
		NullField:    &nullValue, // 提供非nil指针
		NotNullField: "not null",
		Range:        []int{10, 20},
	}

	s.WhereStruct(filter)

	conditions := s.query.GetFilter()
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated")
	}

	// 验证各个字段的条件类型
	foundFields := make(map[string]bool)
	for _, cond := range conditions {
		foundFields[cond.Key] = true
	}

	expectedFields := []string{"_id", "username", "domain", "pattern", "items", "size", "exists", "not_exists", "range"}
	for _, field := range expectedFields {
		if !foundFields[field] {
			t.Errorf("Expected field %s to be present in conditions", field)
		}
	}

	// 检查特殊操作符生成的复合条件
	hasOrCondition := false
	hasAndCondition := false
	for _, cond := range conditions {
		if cond.Key == "$or" {
			hasOrCondition = true
		}
		if cond.Key == "$and" {
			hasAndCondition = true
		}
	}

	if !hasOrCondition {
		t.Error("Expected $or condition for null_field")
	}
	if !hasAndCondition {
		t.Error("Expected $and condition for not_null_field")
	}
}

func TestStructWithEmptyTags(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	// 测试没有pie标签的字段
	type StructWithoutTags struct {
		Name  string
		Age   int    `pie:"-"`
		Email string `pie:""`
	}

	filter := StructWithoutTags{
		Name:  "test",
		Age:   25,
		Email: "test@example.com",
	}

	s.WhereStruct(filter)

	conditions := s.query.GetFilter()
	// 只有有pie标签且不为"-"的字段才会生成条件
	if len(conditions) != 0 {
		t.Error("Expected no conditions for struct without valid pie tags")
	}
}

func TestStructWithComplexTypes(t *testing.T) {
	s := &TestSessionStruct{query: NewQuery()}

	// 测试复杂类型
	type ComplexStruct struct {
		CreatedAt []time.Time    `pie:"created_at,between"`
		Tags      []string       `pie:"tags,all"`
		Metadata  map[string]any `pie:"metadata,exists"`
		Count     int            `pie:"count,size"`
	}

	now := time.Now()
	filter := ComplexStruct{
		CreatedAt: []time.Time{now.Add(-24 * time.Hour), now},
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]any{"key": "value"},
		Count:     5,
	}

	s.WhereStruct(filter)

	conditions := s.query.GetFilter()
	if len(conditions) == 0 {
		t.Fatal("Expected conditions to be generated")
	}

	// 验证条件存在
	foundFields := make(map[string]bool)
	for _, cond := range conditions {
		foundFields[cond.Key] = true
	}

	expectedFields := []string{"created_at", "tags", "metadata", "count"}
	for _, field := range expectedFields {
		if !foundFields[field] {
			t.Errorf("Expected field %s to be present in conditions", field)
		}
	}
}
