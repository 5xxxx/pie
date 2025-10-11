package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/5xxxx/pie"
)

func main() {
	// 创建引擎
	engine, err := pie.NewEngine(context.Background(), "test_db")
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close(context.Background())

	// 示例1: 基本聚合查询
	basicAggregationExample(engine)

	// 示例2: 复杂聚合查询(你的原始示例优化后)
	complexAggregationExample(engine)

	// 示例3: 多阶段聚合
	multiStageAggregationExample(engine)
}

// 示例1: 基本聚合查询
func basicAggregationExample(engine *pie.Engine) {
	fmt.Println("=== 基本聚合查询示例 ===")

	// 统计每个类别的数量和平均价格
	agg := pie.NewAggregate[map[string]any](engine).
		Collection("products")

	result, err := agg.MatchStage().
		Where("status", "active").
		GroupStage().
		By("category", "$category").
		Count("total").
		Avg("avg_price", "$price").
		Done().
		SortStage().
		Desc("total").
		LimitStage(10).
		Exec(context.Background())

	if err != nil {
		log.Printf("聚合查询失败: %v", err)
		return
	}

	fmt.Printf("查询结果: %+v\n", result.Data)
}

// 示例2: 复杂聚合查询(你的原始示例优化后)
func complexAggregationExample(engine *pie.Engine) {
	fmt.Println("=== 复杂聚合查询示例 ===")

	startUTC := time.Now().AddDate(0, -12, 0) // 12个月前
	endUTC := time.Now()
	MongoYearFormat := "%Y"

	agg := pie.NewAggregate[map[string]any](engine).
		Collection("alerts")

	result, err := agg.MatchStage().
		Between("month", startUTC, endUTC).
		AddFieldsStage().
		Add("year", pie.DateFromString(pie.DateToString("$month", MongoYearFormat, "UTC"))).
		Done().
		UnwindStage("$counts").
		PreserveNullAndEmptyArrays(false).
		Done().
		GroupStage().
		By("year", "$year").
		By("device_id", "$device_id").
		By("device_type", "$device_type").
		By("package_tier", "$counts.tier").
		Sum("alerts_total", "$counts.alerts_total").
		Done().
		GroupStage().
		By("year", "$_id.year").
		By("device_id", "$_id.device_id").
		By("device_type", "$_id.device_type").
		AddToSet("package_tier", "$_id.package_tier").
		Push("counts", pie.M{
			"tier":         "$_id.package_tier",
			"alerts_total": "$alerts_total",
		}).
		Sum("alerts_total", "$alerts_total").
		Done().
		ProjectStage().
		Exclude("_id").
		Field("year", "$_id.year").
		Field("device_id", "$_id.device_id").
		Field("device_type", "$_id.device_type").
		Include("package_tier", "counts", "alerts_total").
		Field("create_time", pie.Now()).
		Field("update_time", pie.Now()).
		Done().
		MergeStage("yearly_alerts").
		On("year", "device_id", "device_type").
		WhenMatched("replace").
		WhenNotMatched("insert").
		Done().
		Exec(context.Background())

	if err != nil {
		log.Printf("复杂聚合查询失败: %v", err)
		return
	}

	fmt.Printf("复杂查询结果: %+v\n", result.Data)
}

// 示例3: 多阶段聚合
func multiStageAggregationExample(engine *pie.Engine) {
	fmt.Println("=== 多阶段聚合示例 ===")

	// 用户行为分析
	agg := pie.NewAggregate[map[string]any](engine).
		Collection("user_activities")

	result, err := agg.MatchStage().
		Where("timestamp", pie.GteExpr("$timestamp", time.Now().AddDate(0, 0, -30))).
		AddFieldsStage().
		Add("day_of_week", pie.DayOfWeek("$timestamp")).
		Add("hour", pie.Hour("$timestamp")).
		Done().
		GroupStage().
		By("user_id", "$user_id").
		By("day_of_week", "$day_of_week").
		Count("activity_count").
		Push("activities", pie.M{
			"action":    "$action",
			"timestamp": "$timestamp",
		}).
		Done().
		GroupStage().
		By("user_id", "$_id.user_id").
		Avg("avg_daily_activities", "$activity_count").
		Max("max_daily_activities", "$activity_count").
		Push("daily_stats", pie.M{
			"day_of_week": "$_id.day_of_week",
			"count":       "$activity_count",
		}).
		Done().
		ProjectStage().
		Exclude("_id").
		Field("user_id", "$_id.user_id").
		Include("avg_daily_activities", "max_daily_activities", "daily_stats").
		Field("user_level", pie.Cond(
			pie.GteExpr("$avg_daily_activities", 10),
			"active",
			pie.Cond(
				pie.GteExpr("$avg_daily_activities", 5),
				"moderate",
				"inactive",
			),
		)).
		Done().
		SortStage().
		Desc("avg_daily_activities").
		LimitStage(20).
		Exec(context.Background())

	if err != nil {
		log.Printf("多阶段聚合查询失败: %v", err)
		return
	}

	fmt.Printf("多阶段查询结果: %+v\n", result.Data)
}

// 示例4: 使用原始阶段
func rawStageExample(engine *pie.Engine) {
	fmt.Println("=== 原始阶段示例 ===")

	agg := pie.NewAggregate[map[string]any](engine).
		Collection("products")

	// 使用原始阶段处理复杂逻辑
	result, err := agg.MatchStage().
		Where("category", "electronics").
		RawStage(pie.M{
			"$addFields": pie.M{
				"discount_price": pie.Cond(
					pie.GteExpr("$price", 100),
					pie.Multiply("$price", 0.9), // 10% 折扣
					"$price",
				),
			},
		}).
		ProjectStage().
		Include("name", "price", "discount_price").
		Done().
		Exec(context.Background())

	if err != nil {
		log.Printf("原始阶段查询失败: %v", err)
		return
	}

	fmt.Printf("原始阶段查询结果: %+v\n", result.Data)
}

// 示例5: 关联查询
func lookupExample(engine *pie.Engine) {
	fmt.Println("=== 关联查询示例 ===")

	agg := pie.NewAggregate[map[string]any](engine).
		Collection("orders")

	result, err := agg.MatchStage().
		Where("status", "completed").
		LookupStage("customers", "customer_id", "_id", "customer_info").
		Done().
		LookupStage("products", "product_id", "_id", "product_info").
		Done().
		AddFieldsStage().
		Add("total_amount", pie.Multiply("$quantity", "$product_info.price")).
		Done().
		GroupStage().
		By("customer_id", "$customer_id").
		Sum("total_orders", 1).
		Sum("total_spent", "$total_amount").
		Avg("avg_order_value", "$total_amount").
		Done().
		ProjectStage().
		Exclude("_id").
		Field("customer_id", "$_id.customer_id").
		Include("total_orders", "total_spent", "avg_order_value").
		Field("customer_level", pie.Cond(
			pie.GteExpr("$total_spent", 1000),
			"VIP",
			pie.Cond(
				pie.GteExpr("$total_spent", 500),
				"Premium",
				"Regular",
			),
		)).
		Done().
		SortStage().
		Desc("total_spent").
		Exec(context.Background())

	if err != nil {
		log.Printf("关联查询失败: %v", err)
		return
	}

	fmt.Printf("关联查询结果: %+v\n", result.Data)
}
