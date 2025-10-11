package pie

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestNewCacheKeyGenerator(t *testing.T) {
	prefix := "test:"
	ckg := NewCacheKeyGenerator(prefix)
	if ckg == nil {
		t.Error("Expected CacheKeyGenerator to be created")
	}
	if ckg.prefix != prefix {
		t.Errorf("Expected prefix to be '%s', got '%s'", prefix, ckg.prefix)
	}

	// 测试空前缀
	ckg2 := NewCacheKeyGenerator("")
	if ckg2.prefix != "" {
		t.Errorf("Expected empty prefix, got '%s'", ckg2.prefix)
	}
}

func TestGenerateQueryKey(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"
	filter := bson.D{{"name", "John"}, {"age", 25}}
	options := map[string]interface{}{"limit": 10, "skip": 0}

	key := ckg.GenerateQueryKey(collection, filter, options)

	// 验证键格式
	expectedPrefix := "pie::query:users:"
	if len(key) <= len(expectedPrefix) {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}
	if key[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}

	// 验证相同输入产生相同键
	key2 := ckg.GenerateQueryKey(collection, filter, options)
	if key != key2 {
		t.Error("Expected same input to generate same key")
	}

	// 验证不同输入产生不同键
	filter2 := bson.D{{"name", "Jane"}}
	key3 := ckg.GenerateQueryKey(collection, filter2, options)
	if key == key3 {
		t.Error("Expected different input to generate different key")
	}

	// 验证不同集合产生不同键
	collection2 := "products"
	key4 := ckg.GenerateQueryKey(collection2, filter, options)
	if key == key4 {
		t.Error("Expected different collection to generate different key")
	}

	// 验证不同选项产生不同键
	options2 := map[string]interface{}{"limit": 20}
	key5 := ckg.GenerateQueryKey(collection, filter, options2)
	if key == key5 {
		t.Error("Expected different options to generate different key")
	}
}

func TestGenerateFindOneKey(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"
	filter := bson.D{{"name", "John"}, {"age", 25}}

	key := ckg.GenerateFindOneKey(collection, filter)

	// 验证键格式
	expectedPrefix := "pie:findone:users:"
	if len(key) <= len(expectedPrefix) {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}
	if key[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}

	// 验证相同输入产生相同键
	key2 := ckg.GenerateFindOneKey(collection, filter)
	if key != key2 {
		t.Error("Expected same input to generate same key")
	}

	// 验证不同输入产生不同键
	filter2 := bson.D{{"name", "Jane"}}
	key3 := ckg.GenerateFindOneKey(collection, filter2)
	if key == key3 {
		t.Error("Expected different input to generate different key")
	}

	// 验证不同集合产生不同键
	collection2 := "products"
	key4 := ckg.GenerateFindOneKey(collection2, filter)
	if key == key4 {
		t.Error("Expected different collection to generate different key")
	}
}

func TestGenerateCountKey(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"
	filter := bson.D{{"status", "active"}}

	key := ckg.GenerateCountKey(collection, filter)

	// 验证键格式
	expectedPrefix := "pie:count:users:"
	if len(key) <= len(expectedPrefix) {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}
	if key[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}

	// 验证相同输入产生相同键
	key2 := ckg.GenerateCountKey(collection, filter)
	if key != key2 {
		t.Error("Expected same input to generate same key")
	}

	// 验证不同输入产生不同键
	filter2 := bson.D{{"status", "inactive"}}
	key3 := ckg.GenerateCountKey(collection, filter2)
	if key == key3 {
		t.Error("Expected different input to generate different key")
	}

	// 验证空过滤器
	filter3 := bson.D{}
	key4 := ckg.GenerateCountKey(collection, filter3)
	if key == key4 {
		t.Error("Expected different filter to generate different key")
	}
}

func TestGenerateDocumentKey(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	// 测试字符串ID
	id1 := "507f1f77bcf86cd799439011"
	key1 := ckg.GenerateDocumentKey(collection, id1)
	expected1 := "pie:doc:users:507f1f77bcf86cd799439011"
	if key1 != expected1 {
		t.Errorf("Expected '%s', got '%s'", expected1, key1)
	}

	// 测试数字ID
	id2 := 12345
	key2 := ckg.GenerateDocumentKey(collection, id2)
	expected2 := "pie:doc:users:12345"
	if key2 != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, key2)
	}

	// 测试不同集合
	collection2 := "products"
	key3 := ckg.GenerateDocumentKey(collection2, id1)
	expected3 := "pie:doc:products:507f1f77bcf86cd799439011"
	if key3 != expected3 {
		t.Errorf("Expected '%s', got '%s'", expected3, key3)
	}

	// 测试不同前缀
	ckg2 := NewCacheKeyGenerator("custom:")
	key4 := ckg2.GenerateDocumentKey(collection, id1)
	expected4 := "custom:doc:users:507f1f77bcf86cd799439011"
	if key4 != expected4 {
		t.Errorf("Expected '%s', got '%s'", expected4, key4)
	}
}

func TestGenerateCollectionPattern(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	pattern := ckg.GenerateCollectionPattern(collection)
	expected := "pie:*:users:*"
	if pattern != expected {
		t.Errorf("Expected '%s', got '%s'", expected, pattern)
	}

	// 测试不同集合
	collection2 := "products"
	pattern2 := ckg.GenerateCollectionPattern(collection2)
	expected2 := "pie:*:products:*"
	if pattern2 != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, pattern2)
	}

	// 测试不同前缀
	ckg2 := NewCacheKeyGenerator("custom:")
	pattern3 := ckg2.GenerateCollectionPattern(collection)
	expected3 := "custom:*:users:*"
	if pattern3 != expected3 {
		t.Errorf("Expected '%s', got '%s'", expected3, pattern3)
	}

	// 测试空集合名
	pattern4 := ckg.GenerateCollectionPattern("")
	expected4 := "pie:*::*"
	if pattern4 != expected4 {
		t.Errorf("Expected '%s', got '%s'", expected4, pattern4)
	}
}

func TestGenerateTagKey(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	tag := "user-updates"

	key := ckg.GenerateTagKey(tag)
	expected := "pie:tag:user-updates"
	if key != expected {
		t.Errorf("Expected '%s', got '%s'", expected, key)
	}

	// 测试不同标签
	tag2 := "product-changes"
	key2 := ckg.GenerateTagKey(tag2)
	expected2 := "pie:tag:product-changes"
	if key2 != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, key2)
	}

	// 测试不同前缀
	ckg2 := NewCacheKeyGenerator("custom:")
	key3 := ckg2.GenerateTagKey(tag)
	expected3 := "custom:tag:user-updates"
	if key3 != expected3 {
		t.Errorf("Expected '%s', got '%s'", expected3, key3)
	}

	// 测试空标签
	key4 := ckg.GenerateTagKey("")
	expected4 := "pie:tag:"
	if key4 != expected4 {
		t.Errorf("Expected '%s', got '%s'", expected4, key4)
	}
}

func TestKeyUniqueness(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	// 测试不同过滤器产生不同键
	filters := []bson.D{
		{{"name", "John"}},
		{{"name", "Jane"}},
		{{"age", 25}},
		{{"age", 30}},
		{{"name", "John"}, {"age", 25}},
		{{"name", "Jane"}, {"age", 30}},
		{{"status", "active"}},
		{{"status", "inactive"}},
	}

	keys := make(map[string]bool)
	for _, filter := range filters {
		key := ckg.GenerateQueryKey(collection, filter, nil)
		if keys[key] {
			t.Errorf("Duplicate key generated for filter %v: %s", filter, key)
		}
		keys[key] = true
	}
}

func TestKeyConsistency(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"
	filter := bson.D{{"name", "John"}, {"age", 25}}
	options := map[string]interface{}{"limit": 10, "sort": bson.D{{"created_at", -1}}}

	// 多次生成相同键，确保一致性
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = ckg.GenerateQueryKey(collection, filter, options)
	}

	// 验证所有键都相同
	for i := 1; i < len(keys); i++ {
		if keys[i] != keys[0] {
			t.Errorf("Inconsistent key generation: %s vs %s", keys[0], keys[i])
		}
	}
}

func TestKeyLength(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	// 测试长过滤器
	longFilter := bson.D{}
	for i := 0; i < 100; i++ {
		longFilter = append(longFilter, bson.E{Key: fmt.Sprintf("field%d", i), Value: fmt.Sprintf("value%d", i)})
	}

	key := ckg.GenerateQueryKey(collection, longFilter, nil)
	if len(key) == 0 {
		t.Error("Expected non-empty key")
	}

	// 测试长选项
	longOptions := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		longOptions[fmt.Sprintf("option%d", i)] = fmt.Sprintf("value%d", i)
	}

	key2 := ckg.GenerateQueryKey(collection, bson.D{{"name", "John"}}, longOptions)
	if len(key2) == 0 {
		t.Error("Expected non-empty key")
	}
}

func TestKeyWithSpecialCharacters(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	// 测试包含特殊字符的过滤器
	specialFilter := bson.D{
		{"email", "user@example.com"},
		{"path", "/api/v1/users"},
		{"query", "name=John&age=25"},
		{"json", `{"nested": {"value": "test"}}`},
	}

	key := ckg.GenerateQueryKey(collection, specialFilter, nil)
	if len(key) == 0 {
		t.Error("Expected non-empty key for special characters")
	}

	// 验证键的唯一性
	key2 := ckg.GenerateQueryKey(collection, bson.D{{"email", "user@example.com"}}, nil)
	if key == key2 {
		t.Error("Expected different keys for different filters")
	}
}

func TestKeyWithNilValues(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	// 测试nil选项
	filter := bson.D{{"name", "John"}}
	key1 := ckg.GenerateQueryKey(collection, filter, nil)
	key2 := ckg.GenerateQueryKey(collection, filter, map[string]interface{}{})

	// nil选项和空选项应该产生不同的键
	if key1 == key2 {
		t.Error("Expected different keys for nil and empty options")
	}

	// 测试nil过滤器
	key3 := ckg.GenerateQueryKey(collection, nil, nil)
	key4 := ckg.GenerateQueryKey(collection, bson.D{}, nil)

	// nil过滤器和空过滤器应该产生不同的键
	if key3 == key4 {
		t.Error("Expected different keys for nil and empty filter")
	}
}

func TestKeyWithComplexDataTypes(t *testing.T) {
	ckg := NewCacheKeyGenerator("pie:")
	collection := "users"

	// 测试复杂数据类型
	complexFilter := bson.D{
		{"array", []interface{}{1, 2, 3}},
		{"nested", bson.D{{"inner", "value"}}},
		{"boolean", true},
		{"float", 3.14},
		{"null", nil},
	}

	key := ckg.GenerateQueryKey(collection, complexFilter, nil)
	if len(key) == 0 {
		t.Error("Expected non-empty key for complex data types")
	}

	// 验证键的唯一性
	simpleFilter := bson.D{{"name", "John"}}
	key2 := ckg.GenerateQueryKey(collection, simpleFilter, nil)
	if key == key2 {
		t.Error("Expected different keys for different filter complexity")
	}
}
