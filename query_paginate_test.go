package pie

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// PaginateTestUser 用于分页测试的用户结构
type PaginateTestUser struct {
	ID      bson.ObjectID `bson:"_id,omitempty" pie:"primary_key"`
	Name    string        `bson:"name"`
	Email   string        `bson:"email"`
	Age     int           `bson:"age"`
	Score   float64       `bson:"score"`
	Created time.Time     `bson:"created"`
	Updated time.Time     `bson:"updated"`
}

// 全局测试变量
var (
	paginateTestClient   *mongo.Client
	paginateTestDatabase *mongo.Database
	paginateTestEngine   *Engine
)

func setupPaginateTestDB(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 从环境变量获取MongoDB连接地址，默认为本地地址
	mongoURI := os.Getenv("MONGO_TEST_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@localhost:27017/pie-test?authSource=admin"
	}

	// 连接到MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
		return
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("MongoDB not available: %v", err)
		return
	}

	paginateTestClient = client
	paginateTestDatabase = client.Database("pie_paginate_test_" + fmt.Sprintf("%d", time.Now().UnixNano()))

	// 创建Engine
	paginateTestEngine, err = NewEngine(ctx, paginateTestDatabase.Name(), WithURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
}

func teardownPaginateTestDB(t *testing.T) {
	if paginateTestClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 删除测试数据库
		if paginateTestDatabase != nil {
			paginateTestDatabase.Drop(ctx)
		}

		paginateTestClient.Disconnect(ctx)
	}
}

func createTestUsers(t *testing.T, count int) []PaginateTestUser {
	users := make([]PaginateTestUser, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		users[i] = PaginateTestUser{
			Name:    fmt.Sprintf("User%d", i+1),
			Email:   fmt.Sprintf("user%d@example.com", i+1),
			Age:     20 + (i % 50),
			Score:   float64(50 + (i % 100)),
			Created: now.Add(-time.Duration(i) * time.Hour),
			Updated: now.Add(-time.Duration(i) * time.Hour),
		}
	}

	return users
}

func TestPaginate(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 25)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试第一页
	result, err := session.Paginate(ctx, PaginateParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to paginate: %v", err)
	}

	if result.Total != 25 {
		t.Errorf("Expected total 25, got %d", result.Total)
	}
	if result.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Page)
	}
	if result.PageSize != 10 {
		t.Errorf("Expected page size 10, got %d", result.PageSize)
	}
	if result.TotalPages != 3 {
		t.Errorf("Expected total pages 3, got %d", result.TotalPages)
	}
	if len(result.Data) != 10 {
		t.Errorf("Expected 10 items on first page, got %d", len(result.Data))
	}
	if !result.HasNext {
		t.Error("Expected has next page")
	}
	if result.HasPrev {
		t.Error("Expected no previous page")
	}

	// 测试第二页
	result, err = session.Paginate(ctx, PaginateParams{
		Page:     2,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to paginate page 2: %v", err)
	}

	if result.Page != 2 {
		t.Errorf("Expected page 2, got %d", result.Page)
	}
	if len(result.Data) != 10 {
		t.Errorf("Expected 10 items on second page, got %d", len(result.Data))
	}
	if !result.HasNext {
		t.Error("Expected has next page")
	}
	if !result.HasPrev {
		t.Error("Expected has previous page")
	}

	// 测试最后一页
	result, err = session.Paginate(ctx, PaginateParams{
		Page:     3,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to paginate page 3: %v", err)
	}

	if result.Page != 3 {
		t.Errorf("Expected page 3, got %d", result.Page)
	}
	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items on last page, got %d", len(result.Data))
	}
	if result.HasNext {
		t.Error("Expected no next page")
	}
	if !result.HasPrev {
		t.Error("Expected has previous page")
	}
}

func TestPaginateWithDefaultValues(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 5)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试默认值
	result, err := session.Paginate(ctx, PaginateParams{
		Page:     0, // 应该默认为1
		PageSize: 0, // 应该默认为20
	})
	if err != nil {
		t.Fatalf("Failed to paginate with defaults: %v", err)
	}

	if result.Page != 1 {
		t.Errorf("Expected page 1 (default), got %d", result.Page)
	}
	if result.PageSize != 20 {
		t.Errorf("Expected page size 20 (default), got %d", result.PageSize)
	}
}

func TestPaginateSimple(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 15)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试简单分页
	result, err := session.PaginateSimple(ctx, PaginateParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to paginate simple: %v", err)
	}

	if result.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Page)
	}
	if result.PageSize != 10 {
		t.Errorf("Expected page size 10, got %d", result.PageSize)
	}
	if len(result.Data) != 10 {
		t.Errorf("Expected 10 items, got %d", len(result.Data))
	}
	if !result.HasNext {
		t.Error("Expected has next page")
	}

	// 测试最后一页
	result, err = session.PaginateSimple(ctx, PaginateParams{
		Page:     2,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to paginate simple page 2: %v", err)
	}

	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items on last page, got %d", len(result.Data))
	}
	if result.HasNext {
		t.Error("Expected no next page")
	}
}

func TestPaginateCursor(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 15)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试游标分页
	result, err := session.PaginateCursor(ctx, CursorPaginateParams{
		PageSize:  5,
		SortField: "created",
	})
	if err != nil {
		t.Fatalf("Failed to paginate cursor: %v", err)
	}

	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result.Data))
	}
	if result.NextCursor == "" {
		t.Error("Expected next cursor")
	}
	if result.HasNext {
		t.Error("Expected has next page")
	}
	if result.HasPrev {
		t.Error("Expected no previous page")
	}

	// 测试下一页
	result, err = session.PaginateCursor(ctx, CursorPaginateParams{
		PageSize:  5,
		SortField: "created",
		Cursor:    result.NextCursor,
	})
	if err != nil {
		t.Fatalf("Failed to paginate cursor next page: %v", err)
	}

	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items on next page, got %d", len(result.Data))
	}
	if result.NextCursor == "" {
		t.Error("Expected next cursor for second page")
	}
	if !result.HasPrev {
		t.Error("Expected has previous page")
	}
}

func TestPaginateCursorByID(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 15)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试ID游标分页
	result, err := session.PaginateCursorByID(ctx, IDCursorParams{
		PageSize: 5,
	})
	if err != nil {
		t.Fatalf("Failed to paginate cursor by ID: %v", err)
	}

	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result.Data))
	}
	if result.NextCursor == "" {
		t.Error("Expected next cursor")
	}
	if result.HasNext {
		t.Error("Expected has next page")
	}
	if result.HasPrev {
		t.Error("Expected no previous page")
	}

	// 测试下一页
	result, err = session.PaginateCursorByID(ctx, IDCursorParams{
		PageSize: 5,
		AfterID:  result.Data[len(result.Data)-1].ID,
	})
	if err != nil {
		t.Fatalf("Failed to paginate cursor by ID next page: %v", err)
	}

	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items on next page, got %d", len(result.Data))
	}
	if !result.HasPrev {
		t.Error("Expected has previous page")
	}
}

func TestPaginateEmptyResult(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 测试空结果
	result, err := session.Paginate(ctx, PaginateParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to paginate empty result: %v", err)
	}

	if result.Total != 0 {
		t.Errorf("Expected total 0, got %d", result.Total)
	}
	if len(result.Data) != 0 {
		t.Errorf("Expected 0 items, got %d", len(result.Data))
	}
	if result.HasNext {
		t.Error("Expected no next page")
	}
	if result.HasPrev {
		t.Error("Expected no previous page")
	}
}

func TestPaginateWithFilter(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 20)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试带过滤条件的分页
	result, err := session.Where("age", bson.D{{Key: "$gte", Value: 30}}).Paginate(ctx, PaginateParams{
		Page:     1,
		PageSize: 5,
	})
	if err != nil {
		t.Fatalf("Failed to paginate with filter: %v", err)
	}

	// 验证过滤结果
	for _, user := range result.Data {
		if user.Age < 30 {
			t.Errorf("Expected age >= 30, got %d", user.Age)
		}
	}
}

func TestPaginateCursorInvalidCursor(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 测试无效游标
	_, err := session.PaginateCursor(ctx, CursorPaginateParams{
		PageSize:  5,
		SortField: "created",
		Cursor:    "invalid-cursor",
	})
	if err == nil {
		t.Error("Expected error for invalid cursor")
	}
}

func TestPaginateCursorNoSortField(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 测试没有排序字段
	_, err := session.PaginateCursor(ctx, CursorPaginateParams{
		PageSize: 5,
	})
	if err == nil {
		t.Error("Expected error for missing sort field")
	}
}

func TestPaginateCursorMultiFieldSort(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 10)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试多字段排序
	result, err := session.PaginateCursor(ctx, CursorPaginateParams{
		PageSize:   5,
		SortFields: []string{"age", "created"},
	})
	if err != nil {
		t.Fatalf("Failed to paginate cursor with multi-field sort: %v", err)
	}

	if len(result.Data) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result.Data))
	}
	if result.NextCursor == "" {
		t.Error("Expected next cursor")
	}
}

func TestPaginateCursorDecode(t *testing.T) {
	// 测试游标编码解码
	cursorMap := map[string]any{
		"created": time.Now(),
		"age":     25,
	}

	cursorJSON, err := json.Marshal(cursorMap)
	if err != nil {
		t.Fatalf("Failed to marshal cursor: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(cursorJSON)

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Failed to decode cursor: %v", err)
	}

	var decodedMap map[string]any
	err = json.Unmarshal(decoded, &decodedMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal cursor: %v", err)
	}

	if len(decodedMap) != 2 {
		t.Errorf("Expected 2 fields in decoded cursor, got %d", len(decodedMap))
	}
}

func TestPaginateBoundaryValues(t *testing.T) {
	setupPaginateTestDB(t)
	defer teardownPaginateTestDB(t)

	ctx := context.Background()
	session := Table[PaginateTestUser](paginateTestEngine)

	// 创建测试数据
	users := createTestUsers(t, 3)
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试边界值
	result, err := session.Paginate(ctx, PaginateParams{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		t.Fatalf("Failed to paginate with page size 1: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Data))
	}
	if result.TotalPages != 3 {
		t.Errorf("Expected 3 total pages, got %d", result.TotalPages)
	}

	// 测试大页面大小
	result, err = session.Paginate(ctx, PaginateParams{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		t.Fatalf("Failed to paginate with large page size: %v", err)
	}

	if len(result.Data) != 3 {
		t.Errorf("Expected 3 items, got %d", len(result.Data))
	}
	if result.TotalPages != 1 {
		t.Errorf("Expected 1 total page, got %d", result.TotalPages)
	}
}
