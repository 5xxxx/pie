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

// ConvenienceTestUser 用于便捷方法测试的用户结构
type ConvenienceTestUser struct {
	ID      bson.ObjectID `bson:"_id,omitempty" pie:"primary_key"`
	Name    string        `bson:"name"`
	Email   string        `bson:"email"`
	Age     int           `bson:"age"`
	Score   float64       `bson:"score"`
	Active  bool          `bson:"active"`
	Tags    []string      `bson:"tags"`
	Created time.Time     `bson:"created"`
	Updated time.Time     `bson:"updated"`
}

// 全局测试变量
var (
	convenienceTestClient   *mongo.Client
	convenienceTestDatabase *mongo.Database
	convenienceTestEngine   *Engine
)

func setupConvenienceTestDB(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 从环境变量获取MongoDB连接地址，默认为本地地址
	mongoURI := os.Getenv("MONGO_TEST_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27018/pie-test"
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

	convenienceTestClient = client
	convenienceTestDatabase = client.Database("pie_convenience_test_" + fmt.Sprintf("%d", time.Now().UnixNano()))

	// 创建Engine
	convenienceTestEngine, err = NewEngine(ctx, convenienceTestDatabase.Name(), WithURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
}

func teardownConvenienceTestDB(t *testing.T) {
	if convenienceTestClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 删除测试数据库
		if convenienceTestDatabase != nil {
			convenienceTestDatabase.Drop(ctx)
		}

		convenienceTestClient.Disconnect(ctx)
	}
}

func createTestUser(t *testing.T, name string) ConvenienceTestUser {
	return ConvenienceTestUser{
		Name:    name,
		Email:   fmt.Sprintf("%s@example.com", name),
		Age:     25,
		Score:   85.5,
		Active:  true,
		Tags:    []string{"test", "user"},
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func TestFindByID(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	result, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试FindByID
	foundUser, err := session.FindByID(ctx, result.InsertedID)
	if err != nil {
		t.Fatalf("Failed to find by ID: %v", err)
	}

	if foundUser == nil {
		t.Fatal("Expected user to be found")
	}
	if foundUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, foundUser.Name)
	}

	// 测试不存在的ID
	nonExistentID := bson.NewObjectID()
	_, err = session.FindByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent ID")
	}
}

func TestFindByIDs(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}

	var ids []bson.ObjectID
	for _, user := range users {
		result, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
		ids = append(ids, result.InsertedID)
	}

	// 测试FindByIDs
	foundUsers, err := session.FindByIDs(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to find by IDs: %v", err)
	}

	if len(foundUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(foundUsers))
	}
}

func TestFirstOne(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试FirstOne
	foundUser, err := session.FirstOne(ctx)
	if err != nil {
		t.Fatalf("Failed to find first one: %v", err)
	}

	if foundUser == nil {
		t.Fatal("Expected user to be found")
	}
	if foundUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, foundUser.Name)
	}

	// 测试空结果
	session = Table[ConvenienceTestUser](convenienceTestEngine)
	foundUser, err = session.Where("name", "nonexistent").FirstOne(ctx)
	if err != nil {
		t.Fatalf("Failed to find first one (empty): %v", err)
	}
	if foundUser != nil {
		t.Error("Expected nil for empty result")
	}
}

func TestExists(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试Exists
	exists, err := session.Where("name", "testuser").Exists(ctx)
	if err != nil {
		t.Fatalf("Failed to check exists: %v", err)
	}
	if !exists {
		t.Error("Expected user to exist")
	}

	// 测试不存在
	exists, err = session.Where("name", "nonexistent").Exists(ctx)
	if err != nil {
		t.Fatalf("Failed to check exists (false): %v", err)
	}
	if exists {
		t.Error("Expected user to not exist")
	}
}

func TestFindAndCount(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试FindAndCount
	foundUsers, count, err := session.FindAndCount(ctx)
	if err != nil {
		t.Fatalf("Failed to find and count: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
	if len(foundUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(foundUsers))
	}
}

func TestPluck(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试Pluck
	var names []string
	err := session.Pluck(ctx, "name", &names)
	if err != nil {
		t.Fatalf("Failed to pluck names: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 names, got %d", len(names))
	}

	// 验证名称
	expectedNames := []string{"user1", "user2", "user3"}
	for i, name := range names {
		if name != expectedNames[i] {
			t.Errorf("Expected name %s, got %s", expectedNames[i], name)
		}
	}
}

func TestValue(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试Value
	var name string
	err = session.Where("name", "testuser").Value(ctx, "name", &name)
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if name != "testuser" {
		t.Errorf("Expected name 'testuser', got %s", name)
	}
}

func TestChunk(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := make([]ConvenienceTestUser, 10)
	for i := 0; i < 10; i++ {
		users[i] = createTestUser(t, fmt.Sprintf("user%d", i+1))
		_, err := session.Insert(ctx, &users[i])
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试Chunk
	chunkCount := 0
	totalUsers := 0

	err := session.Chunk(ctx, 3, func(chunk []ConvenienceTestUser) error {
		chunkCount++
		totalUsers += len(chunk)
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to chunk: %v", err)
	}

	if chunkCount != 4 { // 3 + 3 + 3 + 1
		t.Errorf("Expected 4 chunks, got %d", chunkCount)
	}
	if totalUsers != 10 {
		t.Errorf("Expected 10 total users, got %d", totalUsers)
	}
}

func TestCreate(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 测试Create
	user := createTestUser(t, "testuser")
	result, err := session.Create(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if result.InsertedID.IsZero() {
		t.Error("Expected InsertedID")
	}
}

func TestCreateMany(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 测试CreateMany
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}

	result, err := session.CreateMany(ctx, users)
	if err != nil {
		t.Fatalf("Failed to create many users: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 users, got %d", len(result))
	}
}

func TestFirstOrCreate(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 测试FirstOrCreate - 创建新用户
	user := createTestUser(t, "newuser")
	foundUser, created, err := session.Where("name", "newuser").FirstOrCreate(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to first or create: %v", err)
	}

	if !created {
		t.Error("Expected user to be created")
	}
	if foundUser.Name != "newuser" {
		t.Errorf("Expected name 'newuser', got %s", foundUser.Name)
	}

	// 测试FirstOrCreate - 找到现有用户
	foundUser, created, err = session.Where("name", "newuser").FirstOrCreate(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to first or create (existing): %v", err)
	}

	if created {
		t.Error("Expected user to not be created")
	}
	if foundUser.Name != "newuser" {
		t.Errorf("Expected name 'newuser', got %s", foundUser.Name)
	}
}

func TestUpdateOrCreate(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 测试UpdateOrCreate
	user := createTestUser(t, "upsertuser")
	user.Age = 30

	result, err := session.Where("name", "upsertuser").UpdateOrCreate(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to update or create: %v", err)
	}

	if result.Name != "upsertuser" {
		t.Errorf("Expected name 'upsertuser', got %s", result.Name)
	}
	if result.Age != 30 {
		t.Errorf("Expected age 30, got %d", result.Age)
	}
}

func TestUpdateColumn(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试UpdateColumn
	err = session.Where("name", "testuser").UpdateColumn(ctx, "age", 35)
	if err != nil {
		t.Fatalf("Failed to update column: %v", err)
	}

	// 验证更新
	var updatedUser ConvenienceTestUser
	err = session.Where("name", "testuser").Value(ctx, "age", &updatedUser.Age)
	if err != nil {
		t.Fatalf("Failed to get updated value: %v", err)
	}

	if updatedUser.Age != 35 {
		t.Errorf("Expected age 35, got %d", updatedUser.Age)
	}
}

func TestUpdateColumns(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试UpdateColumns
	updates := map[string]interface{}{
		"age":   40,
		"score": 95.0,
	}
	err = session.Where("name", "testuser").UpdateColumns(ctx, updates)
	if err != nil {
		t.Fatalf("Failed to update columns: %v", err)
	}

	// 验证更新
	var updatedUser ConvenienceTestUser
	err = session.Where("name", "testuser").Value(ctx, "age", &updatedUser.Age)
	if err != nil {
		t.Fatalf("Failed to get updated age: %v", err)
	}
	if updatedUser.Age != 40 {
		t.Errorf("Expected age 40, got %d", updatedUser.Age)
	}
}

func TestIncrement(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	user.Score = 50.0
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试Increment
	err = session.Where("name", "testuser").Increment(ctx, "score", 10.0)
	if err != nil {
		t.Fatalf("Failed to increment: %v", err)
	}

	// 验证增量
	var updatedUser ConvenienceTestUser
	err = session.Where("name", "testuser").Value(ctx, "score", &updatedUser.Score)
	if err != nil {
		t.Fatalf("Failed to get updated score: %v", err)
	}
	if updatedUser.Score != 60.0 {
		t.Errorf("Expected score 60.0, got %f", updatedUser.Score)
	}
}

func TestDecrement(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	user.Score = 100.0
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试Decrement
	err = session.Where("name", "testuser").Decrement(ctx, "score", 20.0)
	if err != nil {
		t.Fatalf("Failed to decrement: %v", err)
	}

	// 验证减量
	var updatedUser ConvenienceTestUser
	err = session.Where("name", "testuser").Value(ctx, "score", &updatedUser.Score)
	if err != nil {
		t.Fatalf("Failed to get updated score: %v", err)
	}
	if updatedUser.Score != 80.0 {
		t.Errorf("Expected score 80.0, got %f", updatedUser.Score)
	}
}

func TestToggle(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	user.Active = true
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试Toggle
	err = session.Where("name", "testuser").Toggle(ctx, "active")
	if err != nil {
		t.Fatalf("Failed to toggle: %v", err)
	}

	// 验证切换
	var updatedUser ConvenienceTestUser
	err = session.Where("name", "testuser").Value(ctx, "active", &updatedUser.Active)
	if err != nil {
		t.Fatalf("Failed to get updated active: %v", err)
	}
	if updatedUser.Active != false {
		t.Errorf("Expected active false, got %v", updatedUser.Active)
	}
}

func TestDeleteByID(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	result, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试DeleteByID
	err = session.DeleteByID(ctx, result.InsertedID)
	if err != nil {
		t.Fatalf("Failed to delete by ID: %v", err)
	}

	// 验证删除
	exists, err := session.Where("_id", result.InsertedID).Exists(ctx)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected user to be deleted")
	}
}

func TestDeleteByIDs(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}

	var ids []bson.ObjectID
	for _, user := range users {
		result, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
		ids = append(ids, result.InsertedID)
	}

	// 测试DeleteByIDs
	count, err := session.DeleteByIDs(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to delete by IDs: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 deleted, got %d", count)
	}
}

func TestDestroy(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	user := createTestUser(t, "testuser")
	_, err := session.Insert(ctx, &user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// 测试Destroy
	err = session.Where("name", "testuser").Destroy(ctx)
	if err != nil {
		t.Fatalf("Failed to destroy: %v", err)
	}

	// 验证删除
	exists, err := session.Where("name", "testuser").Exists(ctx)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected user to be deleted")
	}
}

func TestQuickCount(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试QuickCount
	count, err := session.QuickCount(ctx)
	if err != nil {
		t.Fatalf("Failed to quick count: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestConvenienceSum(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}
	users[0].Score = 10.0
	users[1].Score = 20.0
	users[2].Score = 30.0

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试Sum
	sum, err := session.Sum(ctx, "score")
	if err != nil {
		t.Fatalf("Failed to sum: %v", err)
	}

	if sum != 60.0 {
		t.Errorf("Expected sum 60.0, got %f", sum)
	}
}

func TestConvenienceAvg(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}
	users[0].Score = 10.0
	users[1].Score = 20.0
	users[2].Score = 30.0

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试Avg
	avg, err := session.Avg(ctx, "score")
	if err != nil {
		t.Fatalf("Failed to avg: %v", err)
	}

	if avg != 20.0 {
		t.Errorf("Expected avg 20.0, got %f", avg)
	}
}

func TestMaxValue(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}
	users[0].Score = 10.0
	users[1].Score = 30.0
	users[2].Score = 20.0

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试MaxValue
	max, err := session.MaxValue(ctx, "score")
	if err != nil {
		t.Fatalf("Failed to max value: %v", err)
	}

	if max != 30.0 {
		t.Errorf("Expected max 30.0, got %v", max)
	}
}

func TestMinValue(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 创建测试用户
	users := []ConvenienceTestUser{
		createTestUser(t, "user1"),
		createTestUser(t, "user2"),
		createTestUser(t, "user3"),
	}
	users[0].Score = 10.0
	users[1].Score = 30.0
	users[2].Score = 20.0

	for _, user := range users {
		_, err := session.Insert(ctx, &user)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// 测试MinValue
	min, err := session.MinValue(ctx, "score")
	if err != nil {
		t.Fatalf("Failed to min value: %v", err)
	}

	if min != 10.0 {
		t.Errorf("Expected min 10.0, got %v", min)
	}
}

func TestConvenienceMethodsWithEmptyResult(t *testing.T) {
	setupConvenienceTestDB(t)
	defer teardownConvenienceTestDB(t)

	ctx := context.Background()
	session := Table[ConvenienceTestUser](convenienceTestEngine)

	// 测试空结果的各种方法
	_, err := session.Where("name", "nonexistent").FirstOne(ctx)
	if err != nil {
		t.Fatalf("Failed to get first one (empty): %v", err)
	}

	exists, err := session.Where("name", "nonexistent").Exists(ctx)
	if err != nil {
		t.Fatalf("Failed to check exists (empty): %v", err)
	}
	if exists {
		t.Error("Expected exists to be false")
	}

	count, err := session.Where("name", "nonexistent").QuickCount(ctx)
	if err != nil {
		t.Fatalf("Failed to quick count (empty): %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	sum, err := session.Where("name", "nonexistent").Sum(ctx, "score")
	if err != nil {
		t.Fatalf("Failed to sum (empty): %v", err)
	}
	if sum != 0 {
		t.Errorf("Expected sum 0, got %f", sum)
	}
}
