package pie

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SessionTestUser 测试用的用户结构
type SessionTestUser struct {
	ID      bson.ObjectID      `bson:"_id,omitempty" pie:"primary_key"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Age     int                `bson:"age"`
	Active  bool               `bson:"active"`
	Tags    []string           `bson:"tags"`
	Profile SessionTestProfile `bson:"profile"`
	Created time.Time          `bson:"created"`
	Updated time.Time          `bson:"updated"`
}

// SessionTestProfile 测试用的用户资料结构
type SessionTestProfile struct {
	Bio      string `bson:"bio"`
	Location string `bson:"location"`
	Website  string `bson:"website"`
}

// SessionTestProduct 测试用的产品结构
type SessionTestProduct struct {
	ID       bson.ObjectID `bson:"_id,omitempty" pie:"primary_key"`
	Name     string        `bson:"name"`
	Price    float64       `bson:"price"`
	Category string        `bson:"category"`
	InStock  bool          `bson:"in_stock"`
	Tags     []string      `bson:"tags"`
	Created  time.Time     `bson:"created"`
}

// SessionTestHookUser 实现hooks的测试用户
type SessionTestHookUser struct {
	ID        bson.ObjectID `bson:"_id,omitempty" pie:"primary_key"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	HookCalls []string      `bson:"hook_calls"`
}

// BeforeCreate hook
func (u *SessionTestHookUser) BeforeCreate(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "BeforeCreate")
	return nil
}

// AfterCreate hook
func (u *SessionTestHookUser) AfterCreate(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "AfterCreate")
	return nil
}

// BeforeUpdate hook
func (u *SessionTestHookUser) BeforeUpdate(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "BeforeUpdate")
	return nil
}

// AfterUpdate hook
func (u *SessionTestHookUser) AfterUpdate(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "AfterUpdate")
	return nil
}

// BeforeDelete hook
func (u *SessionTestHookUser) BeforeDelete(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "BeforeDelete")
	return nil
}

// AfterDelete hook
func (u *SessionTestHookUser) AfterDelete(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "AfterDelete")
	return nil
}

// AfterFind hook
func (u *SessionTestHookUser) AfterFind(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "AfterFind")
	return nil
}

// BeforeSave hook
func (u *SessionTestHookUser) BeforeSave(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "BeforeSave")
	return nil
}

// AfterSave hook
func (u *SessionTestHookUser) AfterSave(ctx context.Context) error {
	u.HookCalls = append(u.HookCalls, "AfterSave")
	return nil
}

// SessionTestErrorUser 用于测试错误场景的用户
type SessionTestErrorUser struct {
	ID   bson.ObjectID `bson:"_id,omitempty" pie:"primary_key"`
	Name string        `bson:"name"`
}

// BeforeCreate hook that returns error
func (u *SessionTestErrorUser) BeforeCreate(ctx context.Context) error {
	return fmt.Errorf("test error in BeforeCreate")
}

// 测试环境设置
var (
	sessionTestEngine   *Engine
	sessionTestDatabase *mongo.Database
	sessionTestClient   *mongo.Client
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) {
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

	sessionTestClient = client
	sessionTestDatabase = client.Database("pie_test_" + fmt.Sprintf("%d", time.Now().UnixNano()))

	// 创建Engine
	sessionTestEngine, err = NewEngine(ctx, sessionTestDatabase.Name(), WithURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
}

// teardownTestDB 清理测试数据库
func teardownTestDB(t *testing.T) {
	if sessionTestClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 删除测试数据库
		if sessionTestDatabase != nil {
			sessionTestDatabase.Drop(ctx)
		}

		sessionTestClient.Disconnect(ctx)
	}
}

// TestSessionBasicOperations 测试基础CRUD操作
func TestSessionBasicOperations(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 测试插入
	user := &SessionTestUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Tags:    []string{"admin", "user"},
		Profile: SessionTestProfile{Bio: "Software developer", Location: "NYC"},
		Created: time.Now(),
		Updated: time.Now(),
	}

	result, err := session.Insert(ctx, user)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if result.InsertedID.IsZero() {
		t.Error("Expected InsertedID to be set")
	}

	// 测试查找单条
	foundUser, err := session.Where("name", "John Doe").FindOne(ctx)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if foundUser == nil {
		t.Fatal("Expected to find user")
	}
	if foundUser.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", foundUser.Name)
	}

	// 测试查找多条
	users, err := session.Where("active", true).Find(ctx)
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}
	if len(users) == 0 {
		t.Error("Expected to find users")
	}

	// 测试更新
	updateResult, err := session.Where("name", "John Doe").Update(ctx, bson.D{
		{Key: "$set", Value: bson.D{{Key: "age", Value: 31}}},
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updateResult.ModifiedCount == 0 {
		t.Error("Expected to modify at least one document")
	}

	// 测试计数
	count, err := session.Where("active", true).Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count == 0 {
		t.Error("Expected count to be greater than 0")
	}

	// 测试删除
	deleteResult, err := session.Where("name", "John Doe").Delete(ctx)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteResult.DeletedCount == 0 {
		t.Error("Expected to delete at least one document")
	}
}

// TestSessionInsertMany 测试批量插入
func TestSessionInsertMany(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	users := []SessionTestUser{
		{Name: "Alice", Email: "alice@example.com", Age: 25, Active: true},
		{Name: "Bob", Email: "bob@example.com", Age: 30, Active: true},
		{Name: "Charlie", Email: "charlie@example.com", Age: 35, Active: false},
	}

	result, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("InsertMany failed: %v", err)
	}
	if len(result.InsertedIDs) != 3 {
		t.Errorf("Expected 3 inserted IDs, got %d", len(result.InsertedIDs))
	}

	// 验证插入的数据
	count, err := session.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

// TestSessionUpdateMany 测试批量更新
func TestSessionUpdateMany(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: false},
		{Name: "Bob", Age: 30, Active: false},
		{Name: "Charlie", Age: 35, Active: true},
	}
	insertResult, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("InsertMany failed: %v", err)
	}
	t.Logf("InsertMany result: %+v", insertResult)

	// 批量更新
	updateResult, err := session.Where("active", false).UpdateMany(ctx, bson.D{
		{Key: "$set", Value: bson.D{{Key: "active", Value: true}}},
	})
	if err != nil {
		t.Fatalf("UpdateMany failed: %v", err)
	}
	if updateResult.ModifiedCount != 2 {
		t.Errorf("Expected to modify 2 documents, got %d", updateResult.ModifiedCount)
	}

	// 验证更新
	// 先检查总文档数（重置查询条件）
	totalCount, err := Table[SessionTestUser](sessionTestEngine).Count(ctx)
	if err != nil {
		t.Fatalf("Total count failed: %v", err)
	}
	t.Logf("Total documents: %d", totalCount)

	activeCount, err := Table[SessionTestUser](sessionTestEngine).Where("active", true).Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	t.Logf("Active documents: %d", activeCount)
	if activeCount != 3 {
		t.Errorf("Expected 3 active users, got %d", activeCount)
	}
}

// TestSessionDeleteMany 测试批量删除
func TestSessionDeleteMany(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: false},
		{Name: "Bob", Age: 30, Active: false},
		{Name: "Charlie", Age: 35, Active: true},
	}
	_, err := session.InsertMany(ctx, users)
	if err != nil {
		t.Fatalf("InsertMany failed: %v", err)
	}

	// 批量删除
	result, err := session.Where("active", false).DeleteMany(ctx)
	if err != nil {
		t.Fatalf("DeleteMany failed: %v", err)
	}
	if result.DeletedCount != 2 {
		t.Errorf("Expected to delete 2 documents, got %d", result.DeletedCount)
	}

	// 验证删除（重置查询条件）
	count, err := Table[SessionTestUser](sessionTestEngine).Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 remaining user, got %d", count)
	}
}

// TestSessionQueryChaining 测试链式查询
func TestSessionQueryChaining(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: true, Tags: []string{"admin"}},
		{Name: "Bob", Age: 30, Active: true, Tags: []string{"user"}},
		{Name: "Charlie", Age: 35, Active: true, Tags: []string{"admin", "user"}},
		{Name: "David", Age: 40, Active: false, Tags: []string{"user"}},
	}
	session.InsertMany(ctx, users)

	// 测试链式查询
	results, err := session.
		Where("active", true).
		Where("age", bson.D{{Key: "$gte", Value: 30}}).
		OrderBy("age").
		Limit(2).
		Find(ctx)
	if err != nil {
		t.Fatalf("Chained query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
	if results[0].Name != "Bob" {
		t.Errorf("Expected first result to be Bob, got %s", results[0].Name)
	}
	if results[1].Name != "Charlie" {
		t.Errorf("Expected second result to be Charlie, got %s", results[1].Name)
	}
}

// TestSessionAdvancedQueries 测试高级查询
func TestSessionAdvancedQueries(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: true, Tags: []string{"admin"}},
		{Name: "Bob", Age: 30, Active: true, Tags: []string{"user"}},
		{Name: "Charlie", Age: 35, Active: true, Tags: []string{"admin", "user"}},
	}
	session.InsertMany(ctx, users)

	// 测试And查询
	results, err := session.
		And(
			Eq("active", true),
			Gte("age", 30),
		).
		Find(ctx)
	if err != nil {
		t.Fatalf("And query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 测试Or查询
	results, err = session.
		Or(
			Eq("name", "Alice"),
			Gte("age", 35),
		).
		Find(ctx)
	if err != nil {
		t.Fatalf("Or query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 测试OrderByDesc
	results, err = session.
		OrderByDesc("age").
		Find(ctx)
	if err != nil {
		t.Fatalf("OrderByDesc query failed: %v", err)
	}
	if results[0].Name != "Charlie" {
		t.Errorf("Expected first result to be Charlie, got %s", results[0].Name)
	}

	// 测试Skip
	results, err = session.
		OrderBy("age").
		Skip(1).
		Limit(1).
		Find(ctx)
	if err != nil {
		t.Fatalf("Skip query failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Bob" {
		t.Errorf("Expected result to be Bob, got %s", results[0].Name)
	}
}

// TestSessionProjection 测试投影
func TestSessionProjection(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Profile: SessionTestProfile{Bio: "Software developer"},
	}
	session.Insert(ctx, user)

	// 测试Select
	results, err := session.
		Where("name", "John Doe").
		Select("name", "email").
		Find(ctx)
	if err != nil {
		t.Fatalf("Select query failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatal("Expected 1 result")
	}
	if results[0].Name != "John Doe" {
		t.Error("Expected name to be preserved")
	}
	if results[0].Email != "john@example.com" {
		t.Error("Expected email to be preserved")
	}
	// Age should be zero value due to projection
	if results[0].Age != 0 {
		t.Error("Expected age to be zero due to projection")
	}

	// 测试Exclude
	results, err = session.
		Where("name", "John Doe").
		Exclude("profile").
		Find(ctx)
	if err != nil {
		t.Fatalf("Exclude query failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatal("Expected 1 result")
	}
	// Profile should be zero value due to exclusion
	if results[0].Profile.Bio != "" {
		t.Error("Expected profile to be empty due to exclusion")
	}
}

// TestSessionDistinct 测试去重查询
func TestSessionDistinct(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: true},
		{Name: "Bob", Age: 30, Active: true},
		{Name: "Charlie", Age: 25, Active: false},
	}
	session.InsertMany(ctx, users)

	// 测试Distinct
	ages, err := session.Where("active", true).Distinct(ctx, "age")
	if err != nil {
		t.Fatalf("Distinct failed: %v", err)
	}
	if len(ages) != 2 {
		t.Errorf("Expected 2 distinct ages, got %d", len(ages))
	}
}

// TestSessionCountDocuments 测试精确计数
func TestSessionCountDocuments(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: true},
		{Name: "Bob", Age: 30, Active: true},
		{Name: "Charlie", Age: 35, Active: false},
	}
	session.InsertMany(ctx, users)

	// 测试CountDocuments
	count, err := session.Where("active", true).CountDocuments(ctx)
	if err != nil {
		t.Fatalf("CountDocuments failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// 测试EstimatedDocumentCount
	estimatedCount, err := session.EstimatedDocumentCount(ctx)
	if err != nil {
		t.Fatalf("EstimatedDocumentCount failed: %v", err)
	}
	if estimatedCount != 3 {
		t.Errorf("Expected estimated count 3, got %d", estimatedCount)
	}
}

// TestSessionFindCursor 测试游标查询
func TestSessionFindCursor(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}
	session.InsertMany(ctx, users)

	// 测试FindCursor
	cursor, err := session.Where("age", bson.D{{Key: "$gte", Value: 30}}).FindCursor(ctx)
	if err != nil {
		t.Fatalf("FindCursor failed: %v", err)
	}
	defer cursor.Close(ctx)

	var results []SessionTestUser
	for cursor.Next(ctx) {
		var user SessionTestUser
		if err := cursor.Decode(&user); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		results = append(results, user)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

// TestSessionReplaceOne 测试替换文档
func TestSessionReplaceOne(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Age:  30,
	}
	session.Insert(ctx, user)

	// 替换文档
	newUser := &SessionTestUser{
		Name: "Jane Doe",
		Age:  25,
	}
	result, err := session.Where("name", "John Doe").ReplaceOne(ctx, newUser)
	if err != nil {
		t.Fatalf("ReplaceOne failed: %v", err)
	}
	if result.ModifiedCount == 0 {
		t.Error("Expected to modify at least one document")
	}

	// 验证替换
	foundUser, err := session.Where("name", "Jane Doe").FindOne(ctx)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if foundUser == nil {
		t.Fatal("Expected to find replaced user")
	}
	if foundUser.Age != 25 {
		t.Errorf("Expected age 25, got %d", foundUser.Age)
	}
}

// TestSessionFindOneAndDelete 测试查找并删除
func TestSessionFindOneAndDelete(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Age:  30,
	}
	session.Insert(ctx, user)

	// 查找并删除
	deletedUser, err := session.Where("name", "John Doe").FindOneAndDelete(ctx)
	if err != nil {
		t.Fatalf("FindOneAndDelete failed: %v", err)
	}
	if deletedUser == nil {
		t.Fatal("Expected to find deleted user")
	}
	if deletedUser.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", deletedUser.Name)
	}

	// 验证删除
	_, err = session.Where("name", "John Doe").FindOne(ctx)
	if err != ErrEmptyResult {
		t.Errorf("Expected ErrEmptyResult, got %v", err)
	}
}

// TestSessionFindOneAndReplace 测试查找并替换
func TestSessionFindOneAndReplace(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Age:  30,
	}
	_, err := session.Insert(ctx, user)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// 查找并替换（返回替换前的文档）
	newUser := &SessionTestUser{
		Name: "Jane Doe",
		Age:  25,
	}
	// 先验证用户是否存在
	existingUser, err := session.Where("name", "John Doe").FindOne(ctx)
	if err != nil {
		t.Fatalf("Failed to find existing user before replacement: %v", err)
	}
	if existingUser == nil {
		t.Fatal("Expected to find existing user before replacement")
	}

	oldUser, err := session.Where("name", "John Doe").FindOneAndReplace(ctx, newUser, false)
	if err != nil {
		t.Fatalf("FindOneAndReplace failed: %v", err)
	}
	if oldUser == nil {
		t.Fatal("Expected to find old user")
	}
	if oldUser.Name != "John Doe" {
		t.Errorf("Expected old name 'John Doe', got '%s'", oldUser.Name)
	}

	// 验证替换是否成功 - 先检查旧用户是否还存在
	_, err = Table[SessionTestUser](sessionTestEngine).Where("name", "John Doe").FindOne(ctx)
	if err == nil {
		t.Error("Old user should not exist after replacement")
	}

	// 验证替换
	foundUser, err := Table[SessionTestUser](sessionTestEngine).Where("name", "Jane Doe").FindOne(ctx)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if foundUser == nil {
		t.Fatal("Expected to find replaced user")
	}
}

// TestSessionFindOneAndUpdate 测试查找并更新
func TestSessionFindOneAndUpdate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Age:  30,
	}
	session.Insert(ctx, user)

	// 查找并更新（返回更新后的文档）
	updatedUser, err := session.Where("name", "John Doe").FindOneAndUpdate(ctx, bson.D{
		{"$set", bson.D{{Key: "age", Value: 31}}},
	}, true)
	if err != nil {
		t.Fatalf("FindOneAndUpdate failed: %v", err)
	}
	if updatedUser == nil {
		t.Fatal("Expected to find updated user")
	}
	if updatedUser.Age != 31 {
		t.Errorf("Expected age 31, got %d", updatedUser.Age)
	}
}

// TestSessionAdvancedOptions 测试高级查询选项
func TestSessionAdvancedOptions(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}
	session.InsertMany(ctx, users)

	// 测试Hint
	results, err := session.
		Where("age", bson.D{{Key: "$gte", Value: 30}}).
		Hint("age_1").
		Find(ctx)
	if err != nil {
		t.Fatalf("Hint query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 测试Comment
	results, err = session.
		Where("age", bson.D{{Key: "$gte", Value: 30}}).
		Comment("test query").
		Find(ctx)
	if err != nil {
		t.Fatalf("Comment query failed: %v", err)
	}

	// 测试BatchSize
	results, err = session.
		BatchSize(1).
		Find(ctx)
	if err != nil {
		t.Fatalf("BatchSize query failed: %v", err)
	}

	// 测试NoCursorTimeout
	results, err = session.
		NoCursorTimeout(true).
		Find(ctx)
	if err != nil {
		t.Fatalf("NoCursorTimeout query failed: %v", err)
	}

	// 测试ReturnKey
	results, err = session.
		ReturnKey(true).
		Find(ctx)
	if err != nil {
		t.Fatalf("ReturnKey query failed: %v", err)
	}

	// 测试ShowRecordId
	results, err = session.
		ShowRecordId(true).
		Find(ctx)
	if err != nil {
		t.Fatalf("ShowRecordId query failed: %v", err)
	}
}

// TestSessionMinMax 测试Min/Max选项
func TestSessionMinMax(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}
	session.InsertMany(ctx, users)

	// 测试Min
	_, err := session.
		Min(bson.D{{Key: "age", Value: 30}}).
		Find(ctx)
	if err != nil {
		t.Fatalf("Min query failed: %v", err)
	}

	// 测试Max
	_, err = session.
		Max(bson.D{{Key: "age", Value: 30}}).
		Find(ctx)
	if err != nil {
		t.Fatalf("Max query failed: %v", err)
	}
}

// TestSessionArrayFilters 测试数组过滤
func TestSessionArrayFilters(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Tags: []string{"admin", "user"},
	}
	session.Insert(ctx, user)

	// 测试ArrayFilters
	result, err := session.
		Where("name", "John Doe").
		ArrayFilters([]interface{}{bson.D{{"tag", "admin"}}}).
		Update(ctx, bson.D{
			{"$set", bson.D{{"tags.$[tag]", "superadmin"}}},
		})
	if err != nil {
		t.Fatalf("ArrayFilters update failed: %v", err)
	}
	if result.ModifiedCount == 0 {
		t.Error("Expected to modify at least one document")
	}
}

// TestSessionLet 测试Let变量
func TestSessionLet(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Age:  30,
	}
	session.Insert(ctx, user)

	// 测试Let
	result, err := session.
		Where("name", "John Doe").
		Let(bson.D{{"newAge", 31}}).
		Update(ctx, bson.D{
			{"$set", bson.D{{"age", "$$newAge"}}},
		})
	if err != nil {
		t.Fatalf("Let update failed: %v", err)
	}
	if result.ModifiedCount == 0 {
		t.Error("Expected to modify at least one document")
	}
}

// TestSessionUpsert 测试Upsert
func TestSessionUpsert(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 测试Upsert - 更新不存在的文档
	result, err := session.
		Where("name", "John Doe").
		Upsert(true).
		Update(ctx, bson.D{
			{"$set", bson.D{{Key: "age", Value: 30}}},
		})
	if err != nil {
		t.Fatalf("Upsert update failed: %v", err)
	}
	if result.UpsertedID.IsZero() {
		t.Error("Expected UpsertedID to be set")
	}

	// 验证插入
	foundUser, err := session.Where("name", "John Doe").FindOne(ctx)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if foundUser == nil {
		t.Fatal("Expected to find upserted user")
	}
	if foundUser.Age != 30 {
		t.Errorf("Expected age 30, got %d", foundUser.Age)
	}
}

// TestSessionUtilityMethods 测试工具方法
func TestSessionUtilityMethods(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	session := Table[SessionTestUser](sessionTestEngine)

	// 测试Clone
	clonedSession := session.Clone()
	if clonedSession == nil {
		t.Fatal("Expected cloned session to be non-nil")
	}
	if clonedSession == session {
		t.Error("Expected cloned session to be different instance")
	}

	// 测试Clear
	session.Where("name", "test").Clear()
	query := session.GetQuery()
	if query == nil {
		t.Fatal("Expected query to be non-nil")
	}

	// 测试GetOptions
	options := session.GetOptions()
	if options == nil {
		t.Fatal("Expected options to be non-nil")
	}

	// 测试SkipHooks
	sessionWithHooks := session.SkipHooks()
	if !sessionWithHooks.skipHooks {
		t.Error("Expected skipHooks to be true")
	}
}

// TestSessionHooks 测试Hooks集成
func TestSessionHooks(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestHookUser](sessionTestEngine)

	// 测试Insert hooks
	user := &SessionTestHookUser{
		Name: "John Doe",
	}
	result, err := session.Insert(ctx, user)
	if err != nil {
		t.Fatalf("Insert with hooks failed: %v", err)
	}
	if result.InsertedID.IsZero() {
		t.Error("Expected InsertedID to be set")
	}

	// 验证hooks被调用
	if len(user.HookCalls) == 0 {
		t.Error("Expected hooks to be called")
	}

	// 测试Update hooks
	updateResult, err := session.Where("name", "John Doe").Update(ctx, bson.D{
		{"$set", bson.D{{"email", "john@example.com"}}},
	})
	if err != nil {
		t.Fatalf("Update with hooks failed: %v", err)
	}
	if updateResult.ModifiedCount == 0 {
		t.Error("Expected to modify at least one document")
	}

	// 测试Find hooks
	foundUser, err := session.Where("name", "John Doe").FindOne(ctx)
	if err != nil {
		t.Fatalf("FindOne with hooks failed: %v", err)
	}
	if foundUser == nil {
		t.Fatal("Expected to find user")
	}

	// 测试Delete hooks
	deleteResult, err := session.Where("name", "John Doe").Delete(ctx)
	if err != nil {
		t.Fatalf("Delete with hooks failed: %v", err)
	}
	if deleteResult.DeletedCount == 0 {
		t.Error("Expected to delete at least one document")
	}
}

// TestSessionErrorHandling 测试错误处理
func TestSessionErrorHandling(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()

	// 测试初始化错误
	invalidEngine := &Engine{}
	session := &Session[SessionTestUser]{
		engine:     invalidEngine,
		collection: nil,
		query:      NewQuery(),
		options:    NewSessionOptions(),
		initErr:    fmt.Errorf("test initialization error"),
	}

	// 测试FindOne with init error
	_, err := session.FindOne(ctx)
	if err == nil {
		t.Error("Expected error for initialization failure")
	}
	if err.Error() != "session initialization failed: test initialization error" {
		t.Errorf("Expected initialization error, got: %v", err)
	}

	// 测试Find with init error
	_, err = session.Find(ctx)
	if err == nil {
		t.Error("Expected error for initialization failure")
	}

	// 测试Insert with init error
	_, err = session.Insert(ctx, &SessionTestUser{})
	if err == nil {
		t.Error("Expected error for initialization failure")
	}

	// 测试Update with init error
	_, err = session.Update(ctx, bson.D{})
	if err == nil {
		t.Error("Expected error for initialization failure")
	}

	// 测试Delete with init error
	_, err = session.Delete(ctx)
	if err == nil {
		t.Error("Expected error for initialization failure")
	}

	// 测试Count with init error
	_, err = session.Count(ctx)
	if err == nil {
		t.Error("Expected error for initialization failure")
	}

	// 测试Distinct with init error
	_, err = session.Distinct(ctx, "name")
	if err == nil {
		t.Error("Expected error for initialization failure")
	}
}

// TestSessionHookErrors 测试Hook错误处理
func TestSessionHookErrors(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestErrorUser](sessionTestEngine)

	// 测试BeforeCreate hook error
	user := &SessionTestErrorUser{
		Name: "John Doe",
	}
	_, err := session.Insert(ctx, user)
	if err == nil {
		t.Error("Expected error from BeforeCreate hook")
	}
	if err.Error() != "test error in BeforeCreate" {
		t.Errorf("Expected hook error, got: %v", err)
	}
}

// TestSessionEmptyResult 测试空结果处理
func TestSessionEmptyResult(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 测试FindOne with no results
	_, err := session.Where("name", "NonExistent").FindOne(ctx)
	if err != ErrEmptyResult {
		t.Errorf("Expected ErrEmptyResult, got: %v", err)
	}

	// 测试FindOneAndDelete with no results
	_, err = session.Where("name", "NonExistent").FindOneAndDelete(ctx)
	if err != ErrEmptyResult {
		t.Errorf("Expected ErrEmptyResult, got: %v", err)
	}

	// 测试FindOneAndReplace with no results
	_, err = session.Where("name", "NonExistent").FindOneAndReplace(ctx, &SessionTestUser{}, false)
	if err != ErrEmptyResult {
		t.Errorf("Expected ErrEmptyResult, got: %v", err)
	}

	// 测试FindOneAndUpdate with no results
	_, err = session.Where("name", "NonExistent").FindOneAndUpdate(ctx, bson.D{}, false)
	if err != ErrEmptyResult {
		t.Errorf("Expected ErrEmptyResult, got: %v", err)
	}
}

// TestSessionConcurrency 测试并发操作
func TestSessionConcurrency(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 并发插入
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer func() { done <- true }()
			user := &SessionTestUser{
				Name: fmt.Sprintf("User%d", i),
				Age:  20 + i,
			}
			_, err := session.Insert(ctx, user)
			if err != nil {
				t.Errorf("Concurrent insert failed: %v", err)
			}
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有数据都被插入
	count, err := session.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 10 {
		t.Errorf("Expected count 10, got %d", count)
	}
}

// TestSessionWithCache 测试带缓存的Session操作
func TestSessionWithCache(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()

	// 设置缓存
	mockCache := &mockCache{
		data:  make(map[string][]byte),
		stats: &CacheStats{},
	}
	NewCacheManager(mockCache, nil)
	testEngine.UseCache(mockCache, nil)

	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name: "John Doe",
		Age:  30,
	}
	session.Insert(ctx, user)

	// 测试带缓存的查询
	results, err := session.
		Where("name", "John Doe").
		Cache(5 * time.Minute).
		Find(ctx)
	if err != nil {
		t.Fatalf("Cached query failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// 验证缓存被使用
	if mockCache.stats.Hits == 0 {
		t.Error("Expected cache to be used")
	}
}

// TestSessionProjectionWithBson 测试使用bson.D的投影
func TestSessionProjectionWithBson(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	user := &SessionTestUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Profile: SessionTestProfile{Bio: "Software developer"},
	}
	session.Insert(ctx, user)

	// 测试Project
	results, err := session.
		Where("name", "John Doe").
		Project(bson.D{{"name", 1}, {"email", 1}}).
		Find(ctx)
	if err != nil {
		t.Fatalf("Project query failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatal("Expected 1 result")
	}
	if results[0].Name != "John Doe" {
		t.Error("Expected name to be preserved")
	}
	if results[0].Email != "john@example.com" {
		t.Error("Expected email to be preserved")
	}
	// Age should be zero value due to projection
	if results[0].Age != 0 {
		t.Error("Expected age to be zero due to projection")
	}
}

// TestSessionWhereOperator 测试WhereOperator
func TestSessionWhereOperator(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}
	session.InsertMany(ctx, users)

	// 测试WhereOperator
	results, err := session.
		WhereOperator(Gte("age", 30)).
		Find(ctx)
	if err != nil {
		t.Fatalf("WhereOperator query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

// TestSessionComplexQueries 测试复杂查询
func TestSessionComplexQueries(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	ctx := context.Background()
	session := Table[SessionTestUser](sessionTestEngine)

	// 插入测试数据
	users := []SessionTestUser{
		{Name: "Alice", Age: 25, Active: true, Tags: []string{"admin"}},
		{Name: "Bob", Age: 30, Active: true, Tags: []string{"user"}},
		{Name: "Charlie", Age: 35, Active: false, Tags: []string{"admin", "user"}},
		{Name: "David", Age: 40, Active: true, Tags: []string{"guest"}},
	}
	session.InsertMany(ctx, users)

	// 复杂查询：活跃用户，年龄>=30，包含admin标签，按年龄降序，限制2条
	results, err := session.
		Where("active", true).
		Where("age", bson.D{{Key: "$gte", Value: 30}}).
		Where("tags", bson.D{{"$in", []string{"admin"}}}).
		OrderByDesc("age").
		Limit(2).
		Find(ctx)
	if err != nil {
		t.Fatalf("Complex query failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Charlie" {
		t.Errorf("Expected Charlie, got %s", results[0].Name)
	}
}
