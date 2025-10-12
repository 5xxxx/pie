package pie

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	testEngine *Engine
	testDB     string = "pie_test"
	testCtx    context.Context
)

// TestMain 设置测试环境
func TestMain(m *testing.M) {
	// 从环境变量获取 MongoDB URI，如果未设置则使用默认值
	mongoURI := os.Getenv("MONGODB_TEST_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@localhost:27017/"
	}

	testCtx = context.Background()

	// 创建测试引擎
	var err error
	testEngine, err = Connect(testCtx, mongoURI, testDB)
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	// 运行测试
	code := m.Run()

	// 清理测试数据库
	if testEngine != nil && testEngine.client != nil {
		_ = testEngine.client.Database(testDB).Drop(testCtx)
		_ = testEngine.Disconnect(testCtx)
	}

	os.Exit(code)
}

// 测试数据结构
type TestUser struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	Age       int           `bson:"age"`
	Score     float64       `bson:"score"`
	Active    bool          `bson:"active"`
	Tags      []string      `bson:"tags"`
	CreatedAt time.Time     `bson:"created_at"`
	Address   TestAddress   `bson:"address"`
	Metadata  bson.M        `bson:"metadata"`
}

type TestAddress struct {
	City    string  `bson:"city"`
	Country string  `bson:"country"`
	Lat     float64 `bson:"lat"`
	Lng     float64 `bson:"lng"`
}

type TestOrder struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	ProductID string        `bson:"product_id"`
	Amount    float64       `bson:"amount"`
	Quantity  int           `bson:"quantity"`
	Status    string        `bson:"status"`
	OrderDate time.Time     `bson:"order_date"`
}

type TestProduct struct {
	ID       string   `bson:"_id"`
	Name     string   `bson:"name"`
	Category string   `bson:"category"`
	Price    float64  `bson:"price"`
	Stock    int      `bson:"stock"`
	Tags     []string `bson:"tags"`
}

// setupTestData 设置测试数据
func setupTestData(t *testing.T, collectionName string) []TestUser {
	t.Helper()

	coll := testEngine.client.Database(testDB).Collection(collectionName)
	_ = coll.Drop(testCtx)

	users := []TestUser{
		{
			ID:        bson.NewObjectID(),
			Name:      "Alice",
			Email:     "alice@example.com",
			Age:       25,
			Score:     95.5,
			Active:    true,
			Tags:      []string{"vip", "premium"},
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
			Address: TestAddress{
				City:    "New York",
				Country: "USA",
				Lat:     40.7128,
				Lng:     -74.0060,
			},
			Metadata: bson.M{"level": 5, "points": 1000},
		},
		{
			ID:        bson.NewObjectID(),
			Name:      "Bob",
			Email:     "bob@example.com",
			Age:       30,
			Score:     88.0,
			Active:    true,
			Tags:      []string{"regular"},
			CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
			Address: TestAddress{
				City:    "San Francisco",
				Country: "USA",
				Lat:     37.7749,
				Lng:     -122.4194,
			},
			Metadata: bson.M{"level": 3, "points": 500},
		},
		{
			ID:        bson.NewObjectID(),
			Name:      "Charlie",
			Email:     "charlie@example.com",
			Age:       35,
			Score:     72.5,
			Active:    false,
			Tags:      []string{"vip"},
			CreatedAt: time.Now().Add(-90 * 24 * time.Hour),
			Address: TestAddress{
				City:    "London",
				Country: "UK",
				Lat:     51.5074,
				Lng:     -0.1278,
			},
			Metadata: bson.M{"level": 2, "points": 200},
		},
		{
			ID:        bson.NewObjectID(),
			Name:      "David",
			Email:     "david@example.com",
			Age:       28,
			Score:     91.0,
			Active:    true,
			Tags:      []string{"premium", "early-adopter"},
			CreatedAt: time.Now().Add(-15 * 24 * time.Hour),
			Address: TestAddress{
				City:    "Tokyo",
				Country: "Japan",
				Lat:     35.6762,
				Lng:     139.6503,
			},
			Metadata: bson.M{"level": 4, "points": 750},
		},
		{
			ID:        bson.NewObjectID(),
			Name:      "Eve",
			Email:     "eve@example.com",
			Age:       22,
			Score:     85.5,
			Active:    true,
			Tags:      []string{"regular", "student"},
			CreatedAt: time.Now().Add(-45 * 24 * time.Hour),
			Address: TestAddress{
				City:    "Berlin",
				Country: "Germany",
				Lat:     52.5200,
				Lng:     13.4050,
			},
			Metadata: bson.M{"level": 3, "points": 400},
		},
	}

	// 插入测试数据
	docs := make([]any, len(users))
	for i, u := range users {
		docs[i] = u
	}
	_, err := coll.InsertMany(testCtx, docs)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	return users
}

// setupTestOrders 设置订单测试数据
func setupTestOrders(t *testing.T, collectionName string, users []TestUser) {
	t.Helper()

	coll := testEngine.client.Database(testDB).Collection(collectionName)
	_ = coll.Drop(testCtx)

	orders := []TestOrder{
		{
			ID:        bson.NewObjectID(),
			UserID:    users[0].ID,
			ProductID: "prod_1",
			Amount:    100.0,
			Quantity:  2,
			Status:    "completed",
			OrderDate: time.Now().Add(-10 * 24 * time.Hour),
		},
		{
			ID:        bson.NewObjectID(),
			UserID:    users[0].ID,
			ProductID: "prod_2",
			Amount:    50.0,
			Quantity:  1,
			Status:    "completed",
			OrderDate: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID:        bson.NewObjectID(),
			UserID:    users[1].ID,
			ProductID: "prod_1",
			Amount:    100.0,
			Quantity:  2,
			Status:    "pending",
			OrderDate: time.Now().Add(-3 * 24 * time.Hour),
		},
		{
			ID:        bson.NewObjectID(),
			UserID:    users[2].ID,
			ProductID: "prod_3",
			Amount:    200.0,
			Quantity:  1,
			Status:    "completed",
			OrderDate: time.Now().Add(-20 * 24 * time.Hour),
		},
	}

	docs := make([]any, len(orders))
	for i, o := range orders {
		docs[i] = o
	}
	_, err := coll.InsertMany(testCtx, docs)
	if err != nil {
		t.Fatalf("Failed to insert test orders: %v", err)
	}
}

// setupTestProducts 设置产品测试数据
func setupTestProducts(t *testing.T, collectionName string) {
	t.Helper()

	coll := testEngine.client.Database(testDB).Collection(collectionName)
	_ = coll.Drop(testCtx)

	products := []TestProduct{
		{ID: "prod_1", Name: "Laptop", Category: "Electronics", Price: 999.99, Stock: 50, Tags: []string{"computer", "portable"}},
		{ID: "prod_2", Name: "Mouse", Category: "Electronics", Price: 29.99, Stock: 200, Tags: []string{"computer", "accessory"}},
		{ID: "prod_3", Name: "Desk", Category: "Furniture", Price: 299.99, Stock: 30, Tags: []string{"office", "home"}},
		{ID: "prod_4", Name: "Chair", Category: "Furniture", Price: 199.99, Stock: 45, Tags: []string{"office", "home"}},
	}

	docs := make([]any, len(products))
	for i, p := range products {
		docs[i] = p
	}
	_, err := coll.InsertMany(testCtx, docs)
	if err != nil {
		t.Fatalf("Failed to insert test products: %v", err)
	}
}

// cleanupCollection 清理集合
func cleanupCollection(t *testing.T, collectionName string) {
	t.Helper()
	coll := testEngine.client.Database(testDB).Collection(collectionName)
	_ = coll.Drop(testCtx)
}

// ========== aggregate.go 基础方法测试 ==========

func TestNewAggregate(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	if agg == nil {
		t.Fatal("NewAggregate should not return nil")
	}
	if agg.engine != testEngine {
		t.Error("Engine should be set correctly")
	}
	if agg.pipeline == nil {
		t.Error("Pipeline should be initialized")
	}
	if agg.options == nil {
		t.Error("Options should be initialized")
	}
}

func TestAggregateCollection(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.Collection("test_users")

	if result.collection == nil {
		t.Error("Collection should be set")
	}
	if result.collection.Name() != "test_users" {
		t.Errorf("Expected collection name 'test_users', got '%s'", result.collection.Name())
	}
}

func TestAggregateCollectionForStruct(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.CollectionForStruct(TestUser{})

	if result.collection == nil {
		t.Error("Collection should be set")
	}
}

func TestAggregateGetPipeline(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	pipeline := agg.GetPipeline()

	if pipeline == nil {
		t.Error("Pipeline should not be nil")
	}
	if len(pipeline) != 0 {
		t.Error("Initial pipeline should be empty")
	}
}

func TestAggregateClone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	agg.LimitStage(10)

	cloned := agg.Clone()
	if cloned == nil {
		t.Fatal("Clone should not return nil")
	}
	if cloned == agg {
		t.Error("Clone should return a different instance")
	}
	if len(cloned.GetPipeline()) != len(agg.GetPipeline()) {
		t.Error("Cloned pipeline should have same length")
	}
}

func TestAggregateClear(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	agg.LimitStage(10)
	agg.SkipStage(5)

	if len(agg.GetPipeline()) == 0 {
		t.Error("Pipeline should not be empty before clear")
	}

	result := agg.Clear()
	if result != agg {
		t.Error("Clear should return self")
	}
	if len(agg.GetPipeline()) != 0 {
		t.Error("Pipeline should be empty after clear")
	}
}

// ========== aggregate.go 阶段创建方法测试 ==========

func TestAggregateMatchStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage()

	if matchStage == nil {
		t.Error("MatchStage should not be nil")
	}
	if matchStage.agg != agg {
		t.Error("MatchStage should reference the aggregate")
	}
}

func TestAggregateAddFieldsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	addFieldsStage := agg.AddFieldsStage()

	if addFieldsStage == nil {
		t.Error("AddFieldsStage should not be nil")
	}
	if addFieldsStage.agg != agg {
		t.Error("AddFieldsStage should reference the aggregate")
	}
}

func TestAggregateSetStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	setStage := agg.SetStage()

	if setStage == nil {
		t.Error("SetStage should not be nil")
	}
	if setStage.agg != agg {
		t.Error("SetStage should reference the aggregate")
	}
}

func TestAggregateUnsetStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.UnsetStage("field1", "field2")

	if result != agg {
		t.Error("UnsetStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateReplaceRootStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ReplaceRootStage("$newRoot")

	if result != agg {
		t.Error("ReplaceRootStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateReplaceWithStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ReplaceWithStage("$replacement")

	if result != agg {
		t.Error("ReplaceWithStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateUnwindStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	unwindStage := agg.UnwindStage("$items")

	if unwindStage == nil {
		t.Error("UnwindStage should not be nil")
	}
	if unwindStage.agg != agg {
		t.Error("UnwindStage should reference the aggregate")
	}
	if unwindStage.path != "$items" {
		t.Error("UnwindStage should set path correctly")
	}
}

func TestAggregateGroupStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage()

	if groupStage == nil {
		t.Error("GroupStage should not be nil")
	}
	if groupStage.agg != agg {
		t.Error("GroupStage should reference the aggregate")
	}
}

func TestAggregateProjectStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	projectStage := agg.ProjectStage()

	if projectStage == nil {
		t.Error("ProjectStage should not be nil")
	}
	if projectStage.agg != agg {
		t.Error("ProjectStage should reference the aggregate")
	}
}

func TestAggregateSortStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	sortStage := agg.SortStage()

	if sortStage == nil {
		t.Error("SortStage should not be nil")
	}
	if sortStage.agg != agg {
		t.Error("SortStage should reference the aggregate")
	}
}

func TestAggregateLimitStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.LimitStage(10)

	if result != agg {
		t.Error("LimitStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateSkipStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SkipStage(5)

	if result != agg {
		t.Error("SkipStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateLookupStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	lookupStage := agg.LookupStage("orders", "user_id", "_id", "user_orders")

	if lookupStage == nil {
		t.Error("LookupStage should not be nil")
	}
	if lookupStage.agg != agg {
		t.Error("LookupStage should reference the aggregate")
	}
	if lookupStage.from != "orders" {
		t.Error("LookupStage should set from correctly")
	}
}

func TestAggregateGraphLookupStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	graphLookupStage := agg.GraphLookupStage("users")

	if graphLookupStage == nil {
		t.Error("GraphLookupStage should not be nil")
	}
	if graphLookupStage.agg != agg {
		t.Error("GraphLookupStage should reference the aggregate")
	}
}

func TestAggregateUnionWithStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	unionWithStage := agg.UnionWithStage("other_collection")

	if unionWithStage == nil {
		t.Error("UnionWithStage should not be nil")
	}
	if unionWithStage.agg != agg {
		t.Error("UnionWithStage should reference the aggregate")
	}
}

func TestAggregateFacetStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	facetStage := agg.FacetStage()

	if facetStage == nil {
		t.Error("FacetStage should not be nil")
	}
	if facetStage.agg != agg {
		t.Error("FacetStage should reference the aggregate")
	}
}

func TestAggregateSampleStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SampleStage(10)

	if result != agg {
		t.Error("SampleStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateCountStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.CountStage("total")

	if result != agg {
		t.Error("CountStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateOutStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.OutStage("output_collection")

	if result != agg {
		t.Error("OutStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateMergeStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	mergeStage := agg.MergeStage("target_collection")

	if mergeStage == nil {
		t.Error("MergeStage should not be nil")
	}
	if mergeStage.agg != agg {
		t.Error("MergeStage should reference the aggregate")
	}
	if mergeStage.into != "target_collection" {
		t.Error("MergeStage should set into correctly")
	}
}

func TestAggregateBucketStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketStage := agg.BucketStage()

	if bucketStage == nil {
		t.Error("BucketStage should not be nil")
	}
	if bucketStage.agg != agg {
		t.Error("BucketStage should reference the aggregate")
	}
}

func TestAggregateBucketAutoStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketAutoStage := agg.BucketAutoStage()

	if bucketAutoStage == nil {
		t.Error("BucketAutoStage should not be nil")
	}
	if bucketAutoStage.agg != agg {
		t.Error("BucketAutoStage should reference the aggregate")
	}
}

func TestAggregateSortByCountStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SortByCountStage("category")

	if result != agg {
		t.Error("SortByCountStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateRedactStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.RedactStage("$$KEEP")

	if result != agg {
		t.Error("RedactStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateGeoNearStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	geoNearStage := agg.GeoNearStage()

	if geoNearStage == nil {
		t.Error("GeoNearStage should not be nil")
	}
	if geoNearStage.agg != agg {
		t.Error("GeoNearStage should reference the aggregate")
	}
}

func TestAggregateSetWindowFieldsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	setWindowFieldsStage := agg.SetWindowFieldsStage()

	if setWindowFieldsStage == nil {
		t.Error("SetWindowFieldsStage should not be nil")
	}
	if setWindowFieldsStage.agg != agg {
		t.Error("SetWindowFieldsStage should reference the aggregate")
	}
}

func TestAggregateDocumentsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.DocumentsStage(bson.M{"name": "test"})

	if result != agg {
		t.Error("DocumentsStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateSearchStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	searchStage := agg.SearchStage()

	if searchStage == nil {
		t.Error("SearchStage should not be nil")
	}
	if searchStage.agg != agg {
		t.Error("SearchStage should reference the aggregate")
	}
}

func TestAggregateSearchMetaStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	searchMetaStage := agg.SearchMetaStage()

	if searchMetaStage == nil {
		t.Error("SearchMetaStage should not be nil")
	}
	if searchMetaStage.agg != agg {
		t.Error("SearchMetaStage should reference the aggregate")
	}
}

func TestAggregateVectorSearchStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	vectorSearchStage := agg.VectorSearchStage()

	if vectorSearchStage == nil {
		t.Error("VectorSearchStage should not be nil")
	}
	if vectorSearchStage.agg != agg {
		t.Error("VectorSearchStage should reference the aggregate")
	}
}

func TestAggregateDensifyStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	densifyStage := agg.DensifyStage()

	if densifyStage == nil {
		t.Error("DensifyStage should not be nil")
	}
	if densifyStage.agg != agg {
		t.Error("DensifyStage should reference the aggregate")
	}
}

func TestAggregateFillStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	fillStage := agg.FillStage()

	if fillStage == nil {
		t.Error("FillStage should not be nil")
	}
	if fillStage.agg != agg {
		t.Error("FillStage should reference the aggregate")
	}
}

func TestAggregateCollStatsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.CollStatsStage(bson.M{"latencyStats": true})

	if result != agg {
		t.Error("CollStatsStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateIndexStatsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.IndexStatsStage()

	if result != agg {
		t.Error("IndexStatsStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregatePlanCacheStatsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.PlanCacheStatsStage()

	if result != agg {
		t.Error("PlanCacheStatsStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateCurrentOpStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.CurrentOpStage(bson.M{"allUsers": true})

	if result != agg {
		t.Error("CurrentOpStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateListSessionsStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ListSessionsStage(bson.M{"allUsers": true})

	if result != agg {
		t.Error("ListSessionsStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateListSampledQueriesStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ListSampledQueriesStage(bson.M{"allUsers": true})

	if result != agg {
		t.Error("ListSampledQueriesStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateChangeStreamStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ChangeStreamStage(bson.M{"fullDocument": "updateLookup"})

	if result != agg {
		t.Error("ChangeStreamStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateChangeStreamSplitLargeEventStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ChangeStreamSplitLargeEventStage()

	if result != agg {
		t.Error("ChangeStreamSplitLargeEventStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestAggregateRawStage(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.RawStage(bson.M{"$customStage": bson.M{"field": "value"}})

	if result != agg {
		t.Error("RawStage should return self")
	}

	pipeline := agg.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

// ========== aggregate.go 选项方法测试 ==========

func TestAggregateSetAllowDiskUse(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SetAllowDiskUse(true)

	if result != agg {
		t.Error("SetAllowDiskUse should return self")
	}
}

func TestAggregateSetBatchSize(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SetBatchSize(100)

	if result != agg {
		t.Error("SetBatchSize should return self")
	}
}

func TestAggregateSetBypassDocumentValidation(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SetBypassDocumentValidation(true)

	if result != agg {
		t.Error("SetBypassDocumentValidation should return self")
	}
}

func TestAggregateSetCollation(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	collation := &options.Collation{Locale: "en_US"}
	result := agg.SetCollation(collation)

	if result != agg {
		t.Error("SetCollation should return self")
	}
}

func TestAggregateSetMaxAwaitTime(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SetMaxAwaitTime(5000) // 5 seconds

	if result != agg {
		t.Error("SetMaxAwaitTime should return self")
	}
}

func TestAggregateSetComment(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SetComment("test aggregation")

	if result != agg {
		t.Error("SetComment should return self")
	}
}

func TestAggregateSetHint(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SetHint("name_1")

	if result != agg {
		t.Error("SetHint should return self")
	}
}

// ========== aggregate.go 执行方法测试 ==========

func TestAggregateExec(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_exec")
	defer cleanupCollection(t, "test_users_exec")

	// 测试基本聚合执行
	agg := NewAggregate[TestUser](testEngine)
	result, err := agg.Collection("test_users_exec").
		MatchStage().Where("active", true).
		Exec(testCtx)

	if err != nil {
		t.Fatalf("Exec should not return error: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Error != nil {
		t.Fatalf("Result should not have error: %v", result.Error)
	}
	if len(result.Data) == 0 {
		t.Error("Should return some active users")
	}
}

func TestAggregateExecOne(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_exec_one")
	defer cleanupCollection(t, "test_users_exec_one")

	// 测试ExecOne
	agg := NewAggregate[TestUser](testEngine)
	var user TestUser
	err := agg.Collection("test_users_exec_one").
		MatchStage().Where("name", "Alice").
		ExecOne(testCtx, &user)

	if err != nil {
		t.Fatalf("ExecOne should not return error: %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("Expected Alice, got %s", user.Name)
	}
}

func TestAggregateExecOneNotFound(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_exec_one_not_found")
	defer cleanupCollection(t, "test_users_exec_one_not_found")

	// 测试ExecOne找不到文档的情况
	agg := NewAggregate[TestUser](testEngine)
	var user TestUser
	err := agg.Collection("test_users_exec_one_not_found").
		MatchStage().Where("name", "NonExistent").
		ExecOne(testCtx, &user)

	if err != ErrEmptyResult {
		t.Errorf("Expected ErrEmptyResult, got %v", err)
	}
}

// ========== aggregate_expr.go 日期表达式测试 ==========

func TestDateToString(t *testing.T) {
	expr := DateToString("$created_at", "%Y-%m-%d", "UTC")
	if expr == nil {
		t.Error("DateToString should not be nil")
	}
	if expr["$dateToString"] == nil {
		t.Error("DateToString should contain $dateToString operator")
	}
}

func TestDateFromString(t *testing.T) {
	expr := DateFromString("$dateString")
	if expr == nil {
		t.Error("DateFromString should not be nil")
	}
	if expr["$dateFromString"] == nil {
		t.Error("DateFromString should contain $dateFromString operator")
	}
}

func TestDateToParts(t *testing.T) {
	expr := DateToParts("$date")
	if expr == nil {
		t.Error("DateToParts should not be nil")
	}
	if expr["$dateToParts"] == nil {
		t.Error("DateToParts should contain $dateToParts operator")
	}
}

func TestDateToPartsWithTimezone(t *testing.T) {
	expr := DateToParts("$date", "UTC")
	if expr == nil {
		t.Error("DateToParts should not be nil")
	}
	if expr["$dateToParts"] == nil {
		t.Error("DateToParts should contain $dateToParts operator")
	}
}

func TestDateFromParts(t *testing.T) {
	parts := M{"year": 2023, "month": 12, "day": 25}
	expr := DateFromParts(parts)
	if expr == nil {
		t.Error("DateFromParts should not be nil")
	}
	if expr["$dateFromParts"] == nil {
		t.Error("DateFromParts should contain $dateFromParts operator")
	}
}

func TestDateAdd(t *testing.T) {
	expr := DateAdd("$date", 1, "day")
	if expr == nil {
		t.Error("DateAdd should not be nil")
	}
	if expr["$dateAdd"] == nil {
		t.Error("DateAdd should contain $dateAdd operator")
	}
}

func TestDateSubtract(t *testing.T) {
	expr := DateSubtract("$date", 1, "day")
	if expr == nil {
		t.Error("DateSubtract should not be nil")
	}
	if expr["$dateSubtract"] == nil {
		t.Error("DateSubtract should contain $dateSubtract operator")
	}
}

func TestDateDiff(t *testing.T) {
	expr := DateDiff("$startDate", "$endDate", "day")
	if expr == nil {
		t.Error("DateDiff should not be nil")
	}
	if expr["$dateDiff"] == nil {
		t.Error("DateDiff should contain $dateDiff operator")
	}
}

func TestDateTrunc(t *testing.T) {
	expr := DateTrunc("$date", "day")
	if expr == nil {
		t.Error("DateTrunc should not be nil")
	}
	if expr["$dateTrunc"] == nil {
		t.Error("DateTrunc should contain $dateTrunc operator")
	}
}

func TestDateTruncWithBinSize(t *testing.T) {
	expr := DateTrunc("$date", "day", 2)
	if expr == nil {
		t.Error("DateTrunc should not be nil")
	}
	if expr["$dateTrunc"] == nil {
		t.Error("DateTrunc should contain $dateTrunc operator")
	}
}

func TestYear(t *testing.T) {
	expr := Year("$date")
	if expr == nil {
		t.Error("Year should not be nil")
	}
	if expr["$year"] == nil {
		t.Error("Year should contain $year operator")
	}
}

func TestMonth(t *testing.T) {
	expr := Month("$date")
	if expr == nil {
		t.Error("Month should not be nil")
	}
	if expr["$month"] == nil {
		t.Error("Month should contain $month operator")
	}
}

func TestWeek(t *testing.T) {
	expr := Week("$date")
	if expr == nil {
		t.Error("Week should not be nil")
	}
	if expr["$week"] == nil {
		t.Error("Week should contain $week operator")
	}
}

func TestDayOfMonth(t *testing.T) {
	expr := DayOfMonth("$date")
	if expr == nil {
		t.Error("DayOfMonth should not be nil")
	}
	if expr["$dayOfMonth"] == nil {
		t.Error("DayOfMonth should contain $dayOfMonth operator")
	}
}

func TestDayOfWeek(t *testing.T) {
	expr := DayOfWeek("$date")
	if expr == nil {
		t.Error("DayOfWeek should not be nil")
	}
	if expr["$dayOfWeek"] == nil {
		t.Error("DayOfWeek should contain $dayOfWeek operator")
	}
}

func TestDayOfYear(t *testing.T) {
	expr := DayOfYear("$date")
	if expr == nil {
		t.Error("DayOfYear should not be nil")
	}
	if expr["$dayOfYear"] == nil {
		t.Error("DayOfYear should contain $dayOfYear operator")
	}
}

func TestHour(t *testing.T) {
	expr := Hour("$date")
	if expr == nil {
		t.Error("Hour should not be nil")
	}
	if expr["$hour"] == nil {
		t.Error("Hour should contain $hour operator")
	}
}

func TestMinute(t *testing.T) {
	expr := Minute("$date")
	if expr == nil {
		t.Error("Minute should not be nil")
	}
	if expr["$minute"] == nil {
		t.Error("Minute should contain $minute operator")
	}
}

func TestSecond(t *testing.T) {
	expr := Second("$date")
	if expr == nil {
		t.Error("Second should not be nil")
	}
	if expr["$second"] == nil {
		t.Error("Second should contain $second operator")
	}
}

func TestMillisecond(t *testing.T) {
	expr := Millisecond("$date")
	if expr == nil {
		t.Error("Millisecond should not be nil")
	}
	if expr["$millisecond"] == nil {
		t.Error("Millisecond should contain $millisecond operator")
	}
}

func TestISOWeek(t *testing.T) {
	expr := ISOWeek("$date")
	if expr == nil {
		t.Error("ISOWeek should not be nil")
	}
	if expr["$isoWeek"] == nil {
		t.Error("ISOWeek should contain $isoWeek operator")
	}
}

func TestISOWeekYear(t *testing.T) {
	expr := ISOWeekYear("$date")
	if expr == nil {
		t.Error("ISOWeekYear should not be nil")
	}
	if expr["$isoWeekYear"] == nil {
		t.Error("ISOWeekYear should contain $isoWeekYear operator")
	}
}

func TestIsoDayOfWeek(t *testing.T) {
	expr := IsoDayOfWeek("$date")
	if expr == nil {
		t.Error("IsoDayOfWeek should not be nil")
	}
	if expr["$isoDayOfWeek"] == nil {
		t.Error("IsoDayOfWeek should contain $isoDayOfWeek operator")
	}
}

// ========== aggregate_expr.go 算术表达式测试 ==========

func TestAdd(t *testing.T) {
	expr := Add("$field1", "$field2", 10)
	if expr == nil {
		t.Error("Add should not be nil")
	}
	if expr["$add"] == nil {
		t.Error("Add should contain $add operator")
	}
}

func TestSubtract(t *testing.T) {
	expr := Subtract("$field1", "$field2")
	if expr == nil {
		t.Error("Subtract should not be nil")
	}
	if expr["$subtract"] == nil {
		t.Error("Subtract should contain $subtract operator")
	}
}

func TestMultiply(t *testing.T) {
	expr := Multiply("$field1", "$field2", 2)
	if expr == nil {
		t.Error("Multiply should not be nil")
	}
	if expr["$multiply"] == nil {
		t.Error("Multiply should contain $multiply operator")
	}
}

func TestDivide(t *testing.T) {
	expr := Divide("$field1", "$field2")
	if expr == nil {
		t.Error("Divide should not be nil")
	}
	if expr["$divide"] == nil {
		t.Error("Divide should contain $divide operator")
	}
}

func TestModExpr(t *testing.T) {
	expr := ModExpr("$field1", "$field2")
	if expr == nil {
		t.Error("ModExpr should not be nil")
	}
	if expr["$mod"] == nil {
		t.Error("ModExpr should contain $mod operator")
	}
}

func TestAbs(t *testing.T) {
	expr := Abs("$field")
	if expr == nil {
		t.Error("Abs should not be nil")
	}
	if expr["$abs"] == nil {
		t.Error("Abs should contain $abs operator")
	}
}

func TestCeil(t *testing.T) {
	expr := Ceil("$field")
	if expr == nil {
		t.Error("Ceil should not be nil")
	}
	if expr["$ceil"] == nil {
		t.Error("Ceil should contain $ceil operator")
	}
}

func TestFloor(t *testing.T) {
	expr := Floor("$field")
	if expr == nil {
		t.Error("Floor should not be nil")
	}
	if expr["$floor"] == nil {
		t.Error("Floor should contain $floor operator")
	}
}

func TestRound(t *testing.T) {
	expr := Round("$field")
	if expr == nil {
		t.Error("Round should not be nil")
	}
	if expr["$round"] == nil {
		t.Error("Round should contain $round operator")
	}
}

func TestRoundWithPlace(t *testing.T) {
	expr := Round("$field", 2)
	if expr == nil {
		t.Error("Round should not be nil")
	}
	if expr["$round"] == nil {
		t.Error("Round should contain $round operator")
	}
}

func TestTrunc(t *testing.T) {
	expr := Trunc("$field")
	if expr == nil {
		t.Error("Trunc should not be nil")
	}
	if expr["$trunc"] == nil {
		t.Error("Trunc should contain $trunc operator")
	}
}

func TestTruncWithPlace(t *testing.T) {
	expr := Trunc("$field", 2)
	if expr == nil {
		t.Error("Trunc should not be nil")
	}
	if expr["$trunc"] == nil {
		t.Error("Trunc should contain $trunc operator")
	}
}

func TestSqrt(t *testing.T) {
	expr := Sqrt("$field")
	if expr == nil {
		t.Error("Sqrt should not be nil")
	}
	if expr["$sqrt"] == nil {
		t.Error("Sqrt should contain $sqrt operator")
	}
}

func TestPow(t *testing.T) {
	expr := Pow("$base", "$exponent")
	if expr == nil {
		t.Error("Pow should not be nil")
	}
	if expr["$pow"] == nil {
		t.Error("Pow should contain $pow operator")
	}
}

func TestExp(t *testing.T) {
	expr := Exp("$field")
	if expr == nil {
		t.Error("Exp should not be nil")
	}
	if expr["$exp"] == nil {
		t.Error("Exp should contain $exp operator")
	}
}

func TestLn(t *testing.T) {
	expr := Ln("$field")
	if expr == nil {
		t.Error("Ln should not be nil")
	}
	if expr["$ln"] == nil {
		t.Error("Ln should contain $ln operator")
	}
}

func TestLog(t *testing.T) {
	expr := Log("$number", "$base")
	if expr == nil {
		t.Error("Log should not be nil")
	}
	if expr["$log"] == nil {
		t.Error("Log should contain $log operator")
	}
}

func TestLog10(t *testing.T) {
	expr := Log10("$field")
	if expr == nil {
		t.Error("Log10 should not be nil")
	}
	if expr["$log10"] == nil {
		t.Error("Log10 should contain $log10 operator")
	}
}

func TestSin(t *testing.T) {
	expr := Sin("$field")
	if expr == nil {
		t.Error("Sin should not be nil")
	}
	if expr["$sin"] == nil {
		t.Error("Sin should contain $sin operator")
	}
}

func TestCos(t *testing.T) {
	expr := Cos("$field")
	if expr == nil {
		t.Error("Cos should not be nil")
	}
	if expr["$cos"] == nil {
		t.Error("Cos should contain $cos operator")
	}
}

func TestTan(t *testing.T) {
	expr := Tan("$field")
	if expr == nil {
		t.Error("Tan should not be nil")
	}
	if expr["$tan"] == nil {
		t.Error("Tan should contain $tan operator")
	}
}

func TestAsin(t *testing.T) {
	expr := Asin("$field")
	if expr == nil {
		t.Error("Asin should not be nil")
	}
	if expr["$asin"] == nil {
		t.Error("Asin should contain $asin operator")
	}
}

func TestAcos(t *testing.T) {
	expr := Acos("$field")
	if expr == nil {
		t.Error("Acos should not be nil")
	}
	if expr["$acos"] == nil {
		t.Error("Acos should contain $acos operator")
	}
}

func TestAtan(t *testing.T) {
	expr := Atan("$field")
	if expr == nil {
		t.Error("Atan should not be nil")
	}
	if expr["$atan"] == nil {
		t.Error("Atan should contain $atan operator")
	}
}

func TestAtan2(t *testing.T) {
	expr := Atan2("$y", "$x")
	if expr == nil {
		t.Error("Atan2 should not be nil")
	}
	if expr["$atan2"] == nil {
		t.Error("Atan2 should contain $atan2 operator")
	}
}

func TestDegreesToRadians(t *testing.T) {
	expr := DegreesToRadians("$field")
	if expr == nil {
		t.Error("DegreesToRadians should not be nil")
	}
	if expr["$degreesToRadians"] == nil {
		t.Error("DegreesToRadians should contain $degreesToRadians operator")
	}
}

func TestRadiansToDegrees(t *testing.T) {
	expr := RadiansToDegrees("$field")
	if expr == nil {
		t.Error("RadiansToDegrees should not be nil")
	}
	if expr["$radiansToDegrees"] == nil {
		t.Error("RadiansToDegrees should contain $radiansToDegrees operator")
	}
}

// ========== aggregate_expr.go 字符串表达式测试 ==========

func TestConcat(t *testing.T) {
	expr := Concat("$field1", " ", "$field2")
	if expr == nil {
		t.Error("Concat should not be nil")
	}
	if expr["$concat"] == nil {
		t.Error("Concat should contain $concat operator")
	}
}

func TestSubStr(t *testing.T) {
	expr := SubStr("$field", 0, 5)
	if expr == nil {
		t.Error("SubStr should not be nil")
	}
	if expr["$substr"] == nil {
		t.Error("SubStr should contain $substr operator")
	}
}

func TestSubStrBytes(t *testing.T) {
	expr := SubStrBytes("$field", 0, 5)
	if expr == nil {
		t.Error("SubStrBytes should not be nil")
	}
	if expr["$substrBytes"] == nil {
		t.Error("SubStrBytes should contain $substrBytes operator")
	}
}

func TestSubStrCP(t *testing.T) {
	expr := SubStrCP("$field", 0, 5)
	if expr == nil {
		t.Error("SubStrCP should not be nil")
	}
	if expr["$substrCP"] == nil {
		t.Error("SubStrCP should contain $substrCP operator")
	}
}

func TestToUpper(t *testing.T) {
	expr := ToUpper("$field")
	if expr == nil {
		t.Error("ToUpper should not be nil")
	}
	if expr["$toUpper"] == nil {
		t.Error("ToUpper should contain $toUpper operator")
	}
}

func TestToLower(t *testing.T) {
	expr := ToLower("$field")
	if expr == nil {
		t.Error("ToLower should not be nil")
	}
	if expr["$toLower"] == nil {
		t.Error("ToLower should contain $toLower operator")
	}
}

func TestStrCaseCmp(t *testing.T) {
	expr := StrCaseCmp("$field1", "$field2")
	if expr == nil {
		t.Error("StrCaseCmp should not be nil")
	}
	if expr["$strcasecmp"] == nil {
		t.Error("StrCaseCmp should contain $strcasecmp operator")
	}
}

func TestStrLenBytes(t *testing.T) {
	expr := StrLenBytes("$field")
	if expr == nil {
		t.Error("StrLenBytes should not be nil")
	}
	if expr["$strLenBytes"] == nil {
		t.Error("StrLenBytes should contain $strLenBytes operator")
	}
}

func TestStrLenCP(t *testing.T) {
	expr := StrLenCP("$field")
	if expr == nil {
		t.Error("StrLenCP should not be nil")
	}
	if expr["$strLenCP"] == nil {
		t.Error("StrLenCP should contain $strLenCP operator")
	}
}

func TestSplit(t *testing.T) {
	expr := Split("$field", ",")
	if expr == nil {
		t.Error("Split should not be nil")
	}
	if expr["$split"] == nil {
		t.Error("Split should contain $split operator")
	}
}

func TestTrim(t *testing.T) {
	expr := Trim("$field")
	if expr == nil {
		t.Error("Trim should not be nil")
	}
	if expr["$trim"] == nil {
		t.Error("Trim should contain $trim operator")
	}
}

func TestLTrim(t *testing.T) {
	expr := LTrim("$field")
	if expr == nil {
		t.Error("LTrim should not be nil")
	}
	if expr["$ltrim"] == nil {
		t.Error("LTrim should contain $ltrim operator")
	}
}

func TestRTrim(t *testing.T) {
	expr := RTrim("$field")
	if expr == nil {
		t.Error("RTrim should not be nil")
	}
	if expr["$rtrim"] == nil {
		t.Error("RTrim should contain $rtrim operator")
	}
}

func TestReplaceOne(t *testing.T) {
	expr := ReplaceOne("$field", "old", "new")
	if expr == nil {
		t.Error("ReplaceOne should not be nil")
	}
	if expr["$replaceOne"] == nil {
		t.Error("ReplaceOne should contain $replaceOne operator")
	}
}

func TestReplaceAll(t *testing.T) {
	expr := ReplaceAll("$field", "old", "new")
	if expr == nil {
		t.Error("ReplaceAll should not be nil")
	}
	if expr["$replaceAll"] == nil {
		t.Error("ReplaceAll should contain $replaceAll operator")
	}
}

func TestIndexOfBytes(t *testing.T) {
	expr := IndexOfBytes("$field", "substring")
	if expr == nil {
		t.Error("IndexOfBytes should not be nil")
	}
	if expr["$indexOfBytes"] == nil {
		t.Error("IndexOfBytes should contain $indexOfBytes operator")
	}
}

func TestIndexOfCP(t *testing.T) {
	expr := IndexOfCP("$field", "substring")
	if expr == nil {
		t.Error("IndexOfCP should not be nil")
	}
	if expr["$indexOfCP"] == nil {
		t.Error("IndexOfCP should contain $indexOfCP operator")
	}
}

// ========== aggregate_expr.go 数组表达式测试 ==========

func TestArrayElemAt(t *testing.T) {
	expr := ArrayElemAt("$array", 0)
	if expr == nil {
		t.Error("ArrayElemAt should not be nil")
	}
	if expr["$arrayElemAt"] == nil {
		t.Error("ArrayElemAt should contain $arrayElemAt operator")
	}
}

func TestArrayToObject(t *testing.T) {
	expr := ArrayToObject("$array")
	if expr == nil {
		t.Error("ArrayToObject should not be nil")
	}
	if expr["$arrayToObject"] == nil {
		t.Error("ArrayToObject should contain $arrayToObject operator")
	}
}

func TestConcatArrays(t *testing.T) {
	expr := ConcatArrays("$array1", "$array2")
	if expr == nil {
		t.Error("ConcatArrays should not be nil")
	}
	if expr["$concatArrays"] == nil {
		t.Error("ConcatArrays should contain $concatArrays operator")
	}
}

func TestFilterArray(t *testing.T) {
	expr := FilterArray("$array", bson.M{"$gt": []any{"$$item", 0}})
	if expr == nil {
		t.Error("FilterArray should not be nil")
	}
	if expr["$filter"] == nil {
		t.Error("FilterArray should contain $filter operator")
	}
}

func TestFirst(t *testing.T) {
	expr := First("$array")
	if expr == nil {
		t.Error("First should not be nil")
	}
	if expr["$first"] == nil {
		t.Error("First should contain $first operator")
	}
}

func TestLast(t *testing.T) {
	expr := Last("$array")
	if expr == nil {
		t.Error("Last should not be nil")
	}
	if expr["$last"] == nil {
		t.Error("Last should contain $last operator")
	}
}

func TestInArray(t *testing.T) {
	expr := InArray("$field", "$array")
	if expr == nil {
		t.Error("InArray should not be nil")
	}
	if expr["$in"] == nil {
		t.Error("InArray should contain $in operator")
	}
}

func TestIndexOfArray(t *testing.T) {
	expr := IndexOfArray("$array", "value")
	if expr == nil {
		t.Error("IndexOfArray should not be nil")
	}
	if expr["$indexOfArray"] == nil {
		t.Error("IndexOfArray should contain $indexOfArray operator")
	}
}

func TestIsArray(t *testing.T) {
	expr := IsArray("$field")
	if expr == nil {
		t.Error("IsArray should not be nil")
	}
	if expr["$isArray"] == nil {
		t.Error("IsArray should contain $isArray operator")
	}
}

func TestMapArray(t *testing.T) {
	expr := MapArray("$array", "item", bson.M{"$multiply": []any{"$$item", 2}})
	if expr == nil {
		t.Error("MapArray should not be nil")
	}
	if expr["$map"] == nil {
		t.Error("MapArray should contain $map operator")
	}
}

func TestObjectToArray(t *testing.T) {
	expr := ObjectToArray("$object")
	if expr == nil {
		t.Error("ObjectToArray should not be nil")
	}
	if expr["$objectToArray"] == nil {
		t.Error("ObjectToArray should contain $objectToArray operator")
	}
}

func TestRange(t *testing.T) {
	expr := Range(0, 10, 2)
	if expr == nil {
		t.Error("Range should not be nil")
	}
	if expr["$range"] == nil {
		t.Error("Range should contain $range operator")
	}
}

func TestReduce(t *testing.T) {
	expr := Reduce("$array", 0, bson.M{"$add": []any{"$$value", "$$this"}})
	if expr == nil {
		t.Error("Reduce should not be nil")
	}
	if expr["$reduce"] == nil {
		t.Error("Reduce should contain $reduce operator")
	}
}

func TestReverseArray(t *testing.T) {
	expr := ReverseArray("$array")
	if expr == nil {
		t.Error("ReverseArray should not be nil")
	}
	if expr["$reverseArray"] == nil {
		t.Error("ReverseArray should contain $reverseArray operator")
	}
}

func TestSizeArray(t *testing.T) {
	expr := SizeArray("$array")
	if expr == nil {
		t.Error("SizeArray should not be nil")
	}
	if expr["$size"] == nil {
		t.Error("SizeArray should contain $size operator")
	}
}

func TestSlice(t *testing.T) {
	expr := Slice("$array", 5)
	if expr == nil {
		t.Error("Slice should not be nil")
	}
	if expr["$slice"] == nil {
		t.Error("Slice should contain $slice operator")
	}
}

func TestSliceWithPosition(t *testing.T) {
	expr := Slice("$array", 5, 2)
	if expr == nil {
		t.Error("Slice should not be nil")
	}
	if expr["$slice"] == nil {
		t.Error("Slice should contain $slice operator")
	}
}

func TestZip(t *testing.T) {
	expr := Zip(true, "$array1", "$array2")
	if expr == nil {
		t.Error("Zip should not be nil")
	}
	if expr["$zip"] == nil {
		t.Error("Zip should contain $zip operator")
	}
}

func TestMaxN(t *testing.T) {
	expr := MaxN("$array", 3)
	if expr == nil {
		t.Error("MaxN should not be nil")
	}
	if expr["$maxN"] == nil {
		t.Error("MaxN should contain $maxN operator")
	}
}

func TestMinN(t *testing.T) {
	expr := MinN("$array", 3)
	if expr == nil {
		t.Error("MinN should not be nil")
	}
	if expr["$minN"] == nil {
		t.Error("MinN should contain $minN operator")
	}
}

func TestSortArray(t *testing.T) {
	expr := SortArray("$array")
	if expr == nil {
		t.Error("SortArray should not be nil")
	}
	if expr["$sortArray"] == nil {
		t.Error("SortArray should contain $sortArray operator")
	}
}

func TestSortArrayWithSortBy(t *testing.T) {
	expr := SortArray("$array", bson.M{"$field": 1})
	if expr == nil {
		t.Error("SortArray should not be nil")
	}
	if expr["$sortArray"] == nil {
		t.Error("SortArray should contain $sortArray operator")
	}
}

// ========== aggregate_expr.go 条件和比较表达式测试 ==========

func TestCond(t *testing.T) {
	expr := Cond(true, "true_value", "false_value")
	if expr == nil {
		t.Error("Cond should not be nil")
	}
	if expr["$cond"] == nil {
		t.Error("Cond should contain $cond operator")
	}
}

func TestIfNull(t *testing.T) {
	expr := IfNull("$field", "default")
	if expr == nil {
		t.Error("IfNull should not be nil")
	}
	if expr["$ifNull"] == nil {
		t.Error("IfNull should contain $ifNull operator")
	}
}

func TestSwitch(t *testing.T) {
	branches := []M{
		{"case": bson.M{"$eq": []any{"$field", 1}}, "then": "one"},
		{"case": bson.M{"$eq": []any{"$field", 2}}, "then": "two"},
	}
	expr := Switch(branches, "default")
	if expr == nil {
		t.Error("Switch should not be nil")
	}
	if expr["$switch"] == nil {
		t.Error("Switch should contain $switch operator")
	}
}

func TestEqExpr(t *testing.T) {
	expr := EqExpr("$field1", "$field2")
	if expr == nil {
		t.Error("EqExpr should not be nil")
	}
	if expr["$eq"] == nil {
		t.Error("EqExpr should contain $eq operator")
	}
}

func TestNeExpr(t *testing.T) {
	expr := NeExpr("$field1", "$field2")
	if expr == nil {
		t.Error("NeExpr should not be nil")
	}
	if expr["$ne"] == nil {
		t.Error("NeExpr should contain $ne operator")
	}
}

func TestGtExpr(t *testing.T) {
	expr := GtExpr("$field1", "$field2")
	if expr == nil {
		t.Error("GtExpr should not be nil")
	}
	if expr["$gt"] == nil {
		t.Error("GtExpr should contain $gt operator")
	}
}

func TestGteExpr(t *testing.T) {
	expr := GteExpr("$field1", "$field2")
	if expr == nil {
		t.Error("GteExpr should not be nil")
	}
	if expr["$gte"] == nil {
		t.Error("GteExpr should contain $gte operator")
	}
}

func TestLtExpr(t *testing.T) {
	expr := LtExpr("$field1", "$field2")
	if expr == nil {
		t.Error("LtExpr should not be nil")
	}
	if expr["$lt"] == nil {
		t.Error("LtExpr should contain $lt operator")
	}
}

func TestLteExpr(t *testing.T) {
	expr := LteExpr("$field1", "$field2")
	if expr == nil {
		t.Error("LteExpr should not be nil")
	}
	if expr["$lte"] == nil {
		t.Error("LteExpr should contain $lte operator")
	}
}

func TestCmp(t *testing.T) {
	expr := Cmp("$field1", "$field2")
	if expr == nil {
		t.Error("Cmp should not be nil")
	}
	if expr["$cmp"] == nil {
		t.Error("Cmp should contain $cmp operator")
	}
}

func TestAndExpr(t *testing.T) {
	expr := AndExpr("$field1", "$field2", "$field3")
	if expr == nil {
		t.Error("AndExpr should not be nil")
	}
	if expr["$and"] == nil {
		t.Error("AndExpr should contain $and operator")
	}
}

func TestOrExpr(t *testing.T) {
	expr := OrExpr("$field1", "$field2", "$field3")
	if expr == nil {
		t.Error("OrExpr should not be nil")
	}
	if expr["$or"] == nil {
		t.Error("OrExpr should contain $or operator")
	}
}

func TestNotExpr(t *testing.T) {
	expr := NotExpr("$field")
	if expr == nil {
		t.Error("NotExpr should not be nil")
	}
	if expr["$not"] == nil {
		t.Error("NotExpr should contain $not operator")
	}
}

// ========== aggregate_expr.go 类型表达式测试 ==========

func TestTypeExpr(t *testing.T) {
	expr := TypeExpr("$field")
	if expr == nil {
		t.Error("TypeExpr should not be nil")
	}
	if expr["$type"] == nil {
		t.Error("TypeExpr should contain $type operator")
	}
}

func TestConvert(t *testing.T) {
	expr := Convert("$field", "string", "error", "null")
	if expr == nil {
		t.Error("Convert should not be nil")
	}
	if expr["$convert"] == nil {
		t.Error("Convert should contain $convert operator")
	}
}

func TestToBool(t *testing.T) {
	expr := ToBool("$field")
	if expr == nil {
		t.Error("ToBool should not be nil")
	}
	if expr["$toBool"] == nil {
		t.Error("ToBool should contain $toBool operator")
	}
}

func TestToDate(t *testing.T) {
	expr := ToDate("$field")
	if expr == nil {
		t.Error("ToDate should not be nil")
	}
	if expr["$toDate"] == nil {
		t.Error("ToDate should contain $toDate operator")
	}
}

func TestToDecimal(t *testing.T) {
	expr := ToDecimal("$field")
	if expr == nil {
		t.Error("ToDecimal should not be nil")
	}
	if expr["$toDecimal"] == nil {
		t.Error("ToDecimal should contain $toDecimal operator")
	}
}

func TestToDouble(t *testing.T) {
	expr := ToDouble("$field")
	if expr == nil {
		t.Error("ToDouble should not be nil")
	}
	if expr["$toDouble"] == nil {
		t.Error("ToDouble should contain $toDouble operator")
	}
}

func TestToInt(t *testing.T) {
	expr := ToInt("$field")
	if expr == nil {
		t.Error("ToInt should not be nil")
	}
	if expr["$toInt"] == nil {
		t.Error("ToInt should contain $toInt operator")
	}
}

func TestToLong(t *testing.T) {
	expr := ToLong("$field")
	if expr == nil {
		t.Error("ToLong should not be nil")
	}
	if expr["$toLong"] == nil {
		t.Error("ToLong should contain $toLong operator")
	}
}

func TestToObjectId(t *testing.T) {
	expr := ToObjectId("$field")
	if expr == nil {
		t.Error("ToObjectId should not be nil")
	}
	if expr["$toObjectId"] == nil {
		t.Error("ToObjectId should contain $toObjectId operator")
	}
}

func TestToString(t *testing.T) {
	expr := ToString("$field")
	if expr == nil {
		t.Error("ToString should not be nil")
	}
	if expr["$toString"] == nil {
		t.Error("ToString should contain $toString operator")
	}
}

func TestIsNumber(t *testing.T) {
	expr := IsNumber("$field")
	if expr == nil {
		t.Error("IsNumber should not be nil")
	}
	if expr["$isNumber"] == nil {
		t.Error("IsNumber should contain $isNumber operator")
	}
}

// ========== aggregate_expr.go 对象表达式测试 ==========

func TestMergeObjects(t *testing.T) {
	expr := MergeObjects("$obj1", "$obj2")
	if expr == nil {
		t.Error("MergeObjects should not be nil")
	}
	if expr["$mergeObjects"] == nil {
		t.Error("MergeObjects should contain $mergeObjects operator")
	}
}

func TestSetField(t *testing.T) {
	expr := SetField("$field", "path", "value")
	if expr == nil {
		t.Error("SetField should not be nil")
	}
	if expr["$setField"] == nil {
		t.Error("SetField should contain $setField operator")
	}
}

func TestGetField(t *testing.T) {
	expr := GetField("$field", "path")
	if expr == nil {
		t.Error("GetField should not be nil")
	}
	if expr["$getField"] == nil {
		t.Error("GetField should contain $getField operator")
	}
}

func TestUnsetField(t *testing.T) {
	expr := UnsetField("$field", "path")
	if expr == nil {
		t.Error("UnsetField should not be nil")
	}
	if expr["$unsetField"] == nil {
		t.Error("UnsetField should contain $unsetField operator")
	}
}

// ========== aggregate_expr.go 累加器表达式测试 ==========

func TestSum(t *testing.T) {
	expr := Sum("$field")
	if expr == nil {
		t.Error("Sum should not be nil")
	}
	if expr["$sum"] == nil {
		t.Error("Sum should contain $sum operator")
	}
}

func TestAvg(t *testing.T) {
	expr := Avg("$field")
	if expr == nil {
		t.Error("Avg should not be nil")
	}
	if expr["$avg"] == nil {
		t.Error("Avg should contain $avg operator")
	}
}

func TestMax(t *testing.T) {
	expr := Max("$field")
	if expr == nil {
		t.Error("Max should not be nil")
	}
	if expr["$max"] == nil {
		t.Error("Max should contain $max operator")
	}
}

func TestMin(t *testing.T) {
	expr := Min("$field")
	if expr == nil {
		t.Error("Min should not be nil")
	}
	if expr["$min"] == nil {
		t.Error("Min should contain $min operator")
	}
}

func TestStdDevPop(t *testing.T) {
	expr := StdDevPop("$field")
	if expr == nil {
		t.Error("StdDevPop should not be nil")
	}
	if expr["$stdDevPop"] == nil {
		t.Error("StdDevPop should contain $stdDevPop operator")
	}
}

func TestStdDevSamp(t *testing.T) {
	expr := StdDevSamp("$field")
	if expr == nil {
		t.Error("StdDevSamp should not be nil")
	}
	if expr["$stdDevSamp"] == nil {
		t.Error("StdDevSamp should contain $stdDevSamp operator")
	}
}

func TestFirstAccumulator(t *testing.T) {
	expr := FirstAccumulator("$field")
	if expr == nil {
		t.Error("FirstAccumulator should not be nil")
	}
	if expr["$first"] == nil {
		t.Error("FirstAccumulator should contain $first operator")
	}
}

func TestLastAccumulator(t *testing.T) {
	expr := LastAccumulator("$field")
	if expr == nil {
		t.Error("LastAccumulator should not be nil")
	}
	if expr["$last"] == nil {
		t.Error("LastAccumulator should contain $last operator")
	}
}

func TestPush(t *testing.T) {
	expr := Push("$field")
	if expr == nil {
		t.Error("Push should not be nil")
	}
	if expr["$push"] == nil {
		t.Error("Push should contain $push operator")
	}
}

func TestAddToSet(t *testing.T) {
	expr := AddToSet("$field")
	if expr == nil {
		t.Error("AddToSet should not be nil")
	}
	if expr["$addToSet"] == nil {
		t.Error("AddToSet should contain $addToSet operator")
	}
}

func TestMergeObjectsAccumulator(t *testing.T) {
	expr := MergeObjectsAccumulator("$field")
	if expr == nil {
		t.Error("MergeObjectsAccumulator should not be nil")
	}
	if expr["$mergeObjects"] == nil {
		t.Error("MergeObjectsAccumulator should contain $mergeObjects operator")
	}
}

func TestAccumulator(t *testing.T) {
	expr := Accumulator("function() { return 0; }", "function(state, value) { return state + value; }", "function(state1, state2) { return state1 + state2; }", nil, nil, nil, nil)
	if expr == nil {
		t.Error("Accumulator should not be nil")
	}
	if expr["$accumulator"] == nil {
		t.Error("Accumulator should contain $accumulator operator")
	}
}

// ========== aggregate_expr.go 其他表达式测试 ==========

func TestLet(t *testing.T) {
	vars := M{"x": 1, "y": 2}
	expr := Let(vars, bson.M{"$add": []any{"$$x", "$$y"}})
	if expr == nil {
		t.Error("Let should not be nil")
	}
	if expr["$let"] == nil {
		t.Error("Let should contain $let operator")
	}
}

func TestLiteral(t *testing.T) {
	expr := Literal("constant")
	if expr == nil {
		t.Error("Literal should not be nil")
	}
	if expr["$literal"] == nil {
		t.Error("Literal should contain $literal operator")
	}
}

func TestNow(t *testing.T) {
	now := Now()
	if now != "$$NOW" {
		t.Error("Now should return $$NOW")
	}
}

func TestNull(t *testing.T) {
	null := Null()
	if null != nil {
		t.Error("Null should return nil")
	}
}

func TestCurrentDate(t *testing.T) {
	expr := CurrentDate()
	if expr == nil {
		t.Error("CurrentDate should not be nil")
	}
	if expr["$currentDate"] == nil {
		t.Error("CurrentDate should contain $currentDate operator")
	}
}

func TestROOT(t *testing.T) {
	root := ROOT()
	if root != "$$ROOT" {
		t.Error("ROOT should return $$ROOT")
	}
}

func TestREMOVE(t *testing.T) {
	remove := REMOVE()
	if remove != "$$REMOVE" {
		t.Error("REMOVE should return $$REMOVE")
	}
}

func TestPRUNE(t *testing.T) {
	prune := PRUNE()
	if prune != "$$PRUNE" {
		t.Error("PRUNE should return $$PRUNE")
	}
}

func TestKEEP(t *testing.T) {
	keep := KEEP()
	if keep != "$$KEEP" {
		t.Error("KEEP should return $$KEEP")
	}
}

func TestDESCEND(t *testing.T) {
	descend := DESCEND()
	if descend != "$$DESCEND" {
		t.Error("DESCEND should return $$DESCEND")
	}
}

func TestField(t *testing.T) {
	field := Field("name")
	if field != "$name" {
		t.Error("Field should return $name")
	}
}

func TestVar(t *testing.T) {
	variable := Var("name")
	if variable != "$$name" {
		t.Error("Var should return $$name")
	}
}

// ========== aggregate_stages.go 阶段构建器测试 ==========

func TestMatchStageWhere(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.MatchStage().Where("active", true)

	if result == nil {
		t.Error("Where should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestMatchStageBetween(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.MatchStage().Between("age", 20, 30)

	if result == nil {
		t.Error("Between should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("Pipeline should have one stage")
	}
}

func TestMatchStageIn(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().In("status", "active", "pending")

	if matchStage == nil {
		t.Error("In should return MatchStage")
	}
}

func TestMatchStageRegex(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Regex("name", "^A")

	if matchStage == nil {
		t.Error("Regex should return MatchStage")
	}
}

func TestMatchStageExists(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Exists("email", true)

	if matchStage == nil {
		t.Error("Exists should return MatchStage")
	}
}

func TestMatchStageText(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Text("search term")

	if matchStage == nil {
		t.Error("Text should return MatchStage")
	}
}

func TestMatchStageAnd(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().And(
		bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 18}}}},
		bson.D{{Key: "active", Value: true}},
	)

	if matchStage == nil {
		t.Error("And should return MatchStage")
	}
}

func TestMatchStageOr(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Or(
		bson.D{{Key: "status", Value: "active"}},
		bson.D{{Key: "status", Value: "pending"}},
	)

	if matchStage == nil {
		t.Error("Or should return MatchStage")
	}
}

func TestMatchStageNor(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Nor(
		bson.D{{Key: "status", Value: "inactive"}},
		bson.D{{Key: "deleted", Value: true}},
	)

	if matchStage == nil {
		t.Error("Nor should return MatchStage")
	}
}

func TestMatchStageNot(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Not(bson.D{{Key: "age", Value: bson.D{{Key: "$lt", Value: 18}}}})

	if matchStage == nil {
		t.Error("Not should return MatchStage")
	}
}

func TestMatchStageRaw(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	matchStage := agg.MatchStage().Raw(bson.D{{Key: "custom", Value: "condition"}})

	if matchStage == nil {
		t.Error("Raw should return MatchStage")
	}
}

func TestMatchStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.MatchStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty MatchStage should not add pipeline stage")
	}
}

func TestAddFieldsStageAdd(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	addFieldsStage := agg.AddFieldsStage().Add("computed", Add("$field1", "$field2"))

	if addFieldsStage == nil {
		t.Error("Add should return AddFieldsStage")
	}
}

func TestAddFieldsStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.AddFieldsStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty AddFieldsStage should not add pipeline stage")
	}
}

func TestUnwindStagePreserveNullAndEmptyArrays(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	unwindStage := agg.UnwindStage("$items").PreserveNullAndEmptyArrays(true)

	if unwindStage == nil {
		t.Error("PreserveNullAndEmptyArrays should return UnwindStage")
	}
}

func TestUnwindStageIncludeArrayIndex(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	unwindStage := agg.UnwindStage("$items").IncludeArrayIndex("index")

	if unwindStage == nil {
		t.Error("IncludeArrayIndex should return UnwindStage")
	}
}

func TestUnwindStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.UnwindStage("$items").Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("UnwindStage should add one pipeline stage")
	}
}

func TestGroupStageBy(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().By("category", "$category")

	if groupStage == nil {
		t.Error("By should return GroupStage")
	}
}

func TestGroupStageByRaw(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().ByRaw("$category")

	if groupStage == nil {
		t.Error("ByRaw should return GroupStage")
	}
}

func TestGroupStageSum(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Sum("total", "$amount")

	if groupStage == nil {
		t.Error("Sum should return GroupStage")
	}
}

func TestGroupStageAvg(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Avg("average", "$score")

	if groupStage == nil {
		t.Error("Avg should return GroupStage")
	}
}

func TestGroupStageCount(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Count("count")

	if groupStage == nil {
		t.Error("Count should return GroupStage")
	}
}

func TestGroupStageMax(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Max("max", "$score")

	if groupStage == nil {
		t.Error("Max should return GroupStage")
	}
}

func TestGroupStageMin(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Min("min", "$score")

	if groupStage == nil {
		t.Error("Min should return GroupStage")
	}
}

func TestGroupStageFirst(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().First("first", "$name")

	if groupStage == nil {
		t.Error("First should return GroupStage")
	}
}

func TestGroupStageLast(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Last("last", "$name")

	if groupStage == nil {
		t.Error("Last should return GroupStage")
	}
}

func TestGroupStagePush(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().Push("names", "$name")

	if groupStage == nil {
		t.Error("Push should return GroupStage")
	}
}

func TestGroupStageAddToSet(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().AddToSet("uniqueNames", "$name")

	if groupStage == nil {
		t.Error("AddToSet should return GroupStage")
	}
}

func TestGroupStageStdDevPop(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().StdDevPop("stdDev", "$score")

	if groupStage == nil {
		t.Error("StdDevPop should return GroupStage")
	}
}

func TestGroupStageStdDevSamp(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().StdDevSamp("stdDev", "$score")

	if groupStage == nil {
		t.Error("StdDevSamp should return GroupStage")
	}
}

func TestGroupStageMergeObjects(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	groupStage := agg.GroupStage().MergeObjects("merged", "$metadata")

	if groupStage == nil {
		t.Error("MergeObjects should return GroupStage")
	}
}

func TestGroupStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.GroupStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("GroupStage should add one pipeline stage")
	}
}

func TestProjectStageInclude(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	projectStage := agg.ProjectStage().Include("name", "email")

	if projectStage == nil {
		t.Error("Include should return ProjectStage")
	}
}

func TestProjectStageExclude(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	projectStage := agg.ProjectStage().Exclude("password", "secret")

	if projectStage == nil {
		t.Error("Exclude should return ProjectStage")
	}
}

func TestProjectStageField(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	projectStage := agg.ProjectStage().Field("fullName", Concat("$firstName", " ", "$lastName"))

	if projectStage == nil {
		t.Error("Field should return ProjectStage")
	}
}

func TestProjectStageSlice(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	projectStage := agg.ProjectStage().Slice("firstThree", "$tags", 3)

	if projectStage == nil {
		t.Error("Slice should return ProjectStage")
	}
}

func TestProjectStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.ProjectStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty ProjectStage should not add pipeline stage")
	}
}

func TestSortStageAsc(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	sortStage := agg.SortStage().Asc("name", "email")

	if sortStage == nil {
		t.Error("Asc should return SortStage")
	}
}

func TestSortStageDesc(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SortStage().Desc("score", "age")

	if result == nil {
		t.Error("Desc should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("SortStage should add one pipeline stage")
	}
}

func TestSortStageField(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	sortStage := agg.SortStage().Field("name", 1)

	if sortStage == nil {
		t.Error("Field should return SortStage")
	}
}

func TestSortStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.SortStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty SortStage should not add pipeline stage")
	}
}

func TestLookupStageLet(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	lookupStage := agg.LookupStage("orders", "user_id", "_id", "user_orders").Let(M{"userId": "$_id"})

	if lookupStage == nil {
		t.Error("Let should return LookupStage")
	}
}

func TestLookupStagePipeline(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	lookupStage := agg.LookupStage("orders", "user_id", "_id", "user_orders").Pipeline(
		bson.M{"$match": bson.M{"status": "completed"}},
		bson.M{"$limit": 10},
	)

	if lookupStage == nil {
		t.Error("Pipeline should return LookupStage")
	}
}

func TestLookupStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.LookupStage("orders", "user_id", "_id", "user_orders").Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("LookupStage should add one pipeline stage")
	}
}

func TestMergeStageOn(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	mergeStage := agg.MergeStage("target").On("userId", "id")

	if mergeStage == nil {
		t.Error("On should return MergeStage")
	}
}

func TestMergeStageWhenMatched(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	mergeStage := agg.MergeStage("target").WhenMatched("replace")

	if mergeStage == nil {
		t.Error("WhenMatched should return MergeStage")
	}
}

func TestMergeStageWhenNotMatched(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	mergeStage := agg.MergeStage("target").WhenNotMatched("insert")

	if mergeStage == nil {
		t.Error("WhenNotMatched should return MergeStage")
	}
}

func TestMergeStageLet(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	mergeStage := agg.MergeStage("target").Let(M{"userId": "$_id"})

	if mergeStage == nil {
		t.Error("Let should return MergeStage")
	}
}

func TestMergeStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.MergeStage("target").Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 1 {
		t.Error("MergeStage should add one pipeline stage")
	}
}

func TestFacetStageFacet(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	facetStage := agg.FacetStage().Facet("activeUsers",
		bson.M{"$match": bson.M{"active": true}},
		bson.M{"$count": "count"},
	)

	if facetStage == nil {
		t.Error("Facet should return FacetStage")
	}
}

func TestFacetStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.FacetStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty FacetStage should not add pipeline stage")
	}
}

func TestBucketStageGroupBy(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketStage := agg.BucketStage().GroupBy("$age")

	if bucketStage == nil {
		t.Error("GroupBy should return BucketStage")
	}
}

func TestBucketStageBoundaries(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketStage := agg.BucketStage().Boundaries(0, 18, 30, 50, 100)

	if bucketStage == nil {
		t.Error("Boundaries should return BucketStage")
	}
}

func TestBucketStageDefault(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketStage := agg.BucketStage().Default("unknown")

	if bucketStage == nil {
		t.Error("Default should return BucketStage")
	}
}

func TestBucketStageOutput(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketStage := agg.BucketStage().Output(M{"count": M{"$sum": 1}})

	if bucketStage == nil {
		t.Error("Output should return BucketStage")
	}
}

func TestBucketStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.BucketStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty BucketStage should not add pipeline stage")
	}
}

func TestBucketAutoStageGroupBy(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketAutoStage := agg.BucketAutoStage().GroupBy("$age")

	if bucketAutoStage == nil {
		t.Error("GroupBy should return BucketAutoStage")
	}
}

func TestBucketAutoStageBuckets(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketAutoStage := agg.BucketAutoStage().Buckets(5)

	if bucketAutoStage == nil {
		t.Error("Buckets should return BucketAutoStage")
	}
}

func TestBucketAutoStageGranularity(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketAutoStage := agg.BucketAutoStage().Granularity("R5")

	if bucketAutoStage == nil {
		t.Error("Granularity should return BucketAutoStage")
	}
}

func TestBucketAutoStageOutput(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	bucketAutoStage := agg.BucketAutoStage().Output(M{"count": M{"$sum": 1}})

	if bucketAutoStage == nil {
		t.Error("Output should return BucketAutoStage")
	}
}

func TestBucketAutoStageDone(t *testing.T) {
	agg := NewAggregate[TestUser](testEngine)
	result := agg.BucketAutoStage().Done()

	if result == nil {
		t.Error("Done should return aggregate")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 0 {
		t.Error("Empty BucketAutoStage should not add pipeline stage")
	}
}

// ========== 集成测试 ==========

func TestComplexAggregationPipeline(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_complex")
	defer cleanupCollection(t, "test_users_complex")

	// 构建复杂的聚合管道
	agg := NewAggregate[bson.M](testEngine)
	result, err := agg.Collection("test_users_complex").
		MatchStage().Where("active", true).
		AddFieldsStage().Add("ageGroup", Cond(
		GteExpr("$age", 30),
		"adult",
		"young",
	)).Done().
		GroupStage().
		By("ageGroup", "$ageGroup").
		Count("total").
		Avg("avgScore", "$score").
		Max("maxScore", "$score").
		Min("minScore", "$score").
		Done().
		ProjectStage().
		Include("ageGroup", "total", "avgScore", "maxScore", "minScore").
		Field("scoreRange", Subtract("$maxScore", "$minScore")).
		Done().
		SortStage().Desc("total").
		LimitStage(10).
		Exec(testCtx)

	if err != nil {
		t.Fatalf("Complex aggregation should not fail: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Error != nil {
		t.Fatalf("Result should not have error: %v", result.Error)
	}
	if len(result.Data) == 0 {
		t.Error("Should return some results")
	}
}

func TestAggregationWithExpressions(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_expressions")
	defer cleanupCollection(t, "test_users_expressions")

	// 测试各种表达式在聚合中的使用
	agg := NewAggregate[bson.M](testEngine)
	result, err := agg.Collection("test_users_expressions").
		AddFieldsStage().
		Add("fullName", Concat("$name", " (", "$email", ")")).
		Add("ageInDays", Multiply("$age", 365)).
		Add("scoreRounded", Round("$score", 1)).
		Add("isHighScore", GteExpr("$score", 90)).
		Add("tagCount", SizeArray("$tags")).
		Done().
		MatchStage().Where("isHighScore", true).
		ProjectStage().
		Include("fullName", "ageInDays", "scoreRounded", "tagCount").
		Done().
		Exec(testCtx)

	if err != nil {
		t.Fatalf("Aggregation with expressions should not fail: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Error != nil {
		t.Fatalf("Result should not have error: %v", result.Error)
	}
}

func TestAggregationWithLookup(t *testing.T) {
	// 设置测试数据
	users := setupTestData(t, "test_users_lookup")
	defer cleanupCollection(t, "test_users_lookup")
	setupTestOrders(t, "test_orders_lookup", users)
	defer cleanupCollection(t, "test_orders_lookup")

	// 测试Lookup关联查询
	agg := NewAggregate[bson.M](testEngine)
	result, err := agg.Collection("test_users_lookup").
		LookupStage("test_orders_lookup", "_id", "user_id", "orders").Done().
		AddFieldsStage().
		Add("orderCount", SizeArray("$orders")).
		Add("totalAmount", Sum("$orders.amount")).
		Done().
		MatchStage().Where("orderCount", GtExpr(0, 0)).
		ProjectStage().
		Include("name", "email", "orderCount", "totalAmount").
		Done().
		Exec(testCtx)

	if err != nil {
		t.Fatalf("Aggregation with lookup should not fail: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Error != nil {
		t.Fatalf("Result should not have error: %v", result.Error)
	}
}

func TestAggregationWithUnwind(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_unwind")
	defer cleanupCollection(t, "test_users_unwind")

	// 测试Unwind展开数组
	agg := NewAggregate[bson.M](testEngine)
	result, err := agg.Collection("test_users_unwind").
		UnwindStage("$tags").Done().
		GroupStage().
		By("tag", "$tags").
		Count("count").
		Done().
		SortStage().Desc("count").
		Exec(testCtx)

	if err != nil {
		t.Fatalf("Aggregation with unwind should not fail: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Error != nil {
		t.Fatalf("Result should not have error: %v", result.Error)
	}
}

func TestAggregationWithFacet(t *testing.T) {
	// 设置测试数据
	setupTestData(t, "test_users_facet")
	defer cleanupCollection(t, "test_users_facet")

	// 测试Facet分面分析
	agg := NewAggregate[bson.M](testEngine)
	result, err := agg.Collection("test_users_facet").
		FacetStage().
		Facet("activeUsers",
			bson.M{"$match": bson.M{"active": true}},
			bson.M{"$count": "count"},
		).
		Facet("scoreStats",
			bson.M{"$group": bson.M{
				"_id":      nil,
				"avgScore": bson.M{"$avg": "$score"},
				"maxScore": bson.M{"$max": "$score"},
				"minScore": bson.M{"$min": "$score"},
			}},
		).
		Facet("ageGroups",
			bson.M{"$bucket": bson.M{
				"groupBy":    "$age",
				"boundaries": []int{0, 25, 30, 35, 100},
				"default":    "other",
				"output":     bson.M{"count": bson.M{"$sum": 1}},
			}},
		).
		Done().
		Exec(testCtx)

	if err != nil {
		t.Fatalf("Aggregation with facet should not fail: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Error != nil {
		t.Fatalf("Result should not have error: %v", result.Error)
	}
}

// ========== 原有测试（保持兼容性） ==========

func TestAggregateStages(t *testing.T) {
	// 创建一个模拟的引擎
	engine := &Engine{}

	// 测试基本阶段构建
	agg := NewAggregate[bson.M](engine)

	// 测试 MatchStage
	matchStage := agg.MatchStage()
	if matchStage == nil {
		t.Error("MatchStage should not be nil")
	}

	// 测试 AddFieldsStage
	addFieldsStage := agg.AddFieldsStage()
	if addFieldsStage == nil {
		t.Error("AddFieldsStage should not be nil")
	}

	// 测试 GroupStage
	groupStage := agg.GroupStage()
	if groupStage == nil {
		t.Error("GroupStage should not be nil")
	}

	// 测试 ProjectStage
	projectStage := agg.ProjectStage()
	if projectStage == nil {
		t.Error("ProjectStage should not be nil")
	}

	// 测试 SortStage
	sortStage := agg.SortStage()
	if sortStage == nil {
		t.Error("SortStage should not be nil")
	}

	// 测试 UnwindStage
	unwindStage := agg.UnwindStage("$items")
	if unwindStage == nil {
		t.Error("UnwindStage should not be nil")
	}

	// 测试 MergeStage
	mergeStage := agg.MergeStage("target_collection")
	if mergeStage == nil {
		t.Error("MergeStage should not be nil")
	}
}

func TestAggregateExpressions(t *testing.T) {
	// 测试日期表达式
	dateExpr := DateToString("$date", "%Y-%m-%d", "UTC")
	if dateExpr == nil {
		t.Error("DateToString should not be nil")
	}

	// 测试算术表达式
	addExpr := Add("$field1", "$field2", 10)
	if addExpr == nil {
		t.Error("Add should not be nil")
	}

	// 测试字符串表达式
	concatExpr := Concat("$field1", " ", "$field2")
	if concatExpr == nil {
		t.Error("Concat should not be nil")
	}

	// 测试条件表达式
	condExpr := Cond(true, "true_value", "false_value")
	if condExpr == nil {
		t.Error("Cond should not be nil")
	}

	// 测试常量
	now := Now()
	if now != "$$NOW" {
		t.Error("Now should return $$NOW")
	}
}

func TestAggregatePipeline(t *testing.T) {
	engine := &Engine{}
	agg := NewAggregate[bson.M](engine)

	// 构建一个简单的管道
	result := agg.MatchStage().
		Where("status", "active").
		GroupStage().
		By("category", "$category").
		Count("total").
		Done().
		ProjectStage().
		Include("category", "total").
		Done().
		SortStage().
		Desc("total")

	if result == nil {
		t.Error("Pipeline should not be nil")
	}

	// 检查管道是否被正确构建
	pipeline := result.GetPipeline()
	if len(pipeline) == 0 {
		t.Error("Pipeline should not be empty")
	}

	// 应该有 4 个阶段: match, group, project, sort
	if len(pipeline) != 4 {
		t.Errorf("Expected 4 stages, got %d", len(pipeline))
	}
}

func TestAggregateChain(t *testing.T) {
	engine := &Engine{}
	agg := NewAggregate[bson.M](engine)

	// 测试链式调用
	result := agg.
		MatchStage().Where("active", true).
		AddFieldsStage().Add("computed_field", Add("$field1", "$field2")).Done().
		GroupStage().By("group_field", "$group").Count("count").Done().
		ProjectStage().Include("group_field", "count").Done().
		SortStage().Desc("count").
		LimitStage(10)

	if result == nil {
		t.Error("Chained pipeline should not be nil")
	}

	pipeline := result.GetPipeline()
	if len(pipeline) != 6 {
		t.Errorf("Expected 6 stages, got %d", len(pipeline))
	}
}
