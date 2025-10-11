package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/5xxxx/pie"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Product product model
type Product struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Price     float64       `bson:"price"`
	Stock     int           `bson:"stock"`
	Category  string        `bson:"category"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

// BeforeCreate before create hook
func (p *Product) BeforeCreate(ctx context.Context) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	fmt.Printf("📝 BeforeCreate: About to create product %s\n", p.Name)
	return nil
}

// AfterCreate after create hook
func (p *Product) AfterCreate(ctx context.Context) error {
	fmt.Printf("✅ AfterCreate: Product %s created, ID=%s\n", p.Name, p.ID.Hex())
	return nil
}

func main() {
	ctx := context.Background()

	// Create engine and enable query logging
	engine, err := pie.NewEngine(ctx, "testdb",
		pie.WithQueryLog(os.Stdout),
	)
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Disconnect(ctx)

	fmt.Println("\n========== 1. Cursor iteration demo ==========")
	demonstrateCursor(ctx, engine)

	fmt.Println("\n========== 2. CountDocuments precise count demo ==========")
	demonstrateCount(ctx, engine)

	fmt.Println("\n========== 3. ReplaceOne document replacement demo ==========")
	demonstrateReplace(ctx, engine)

	fmt.Println("\n========== 4. FindOneAndUpdate find and update demo ==========")
	demonstrateFindOneAndUpdate(ctx, engine)

	fmt.Println("\n========== 5. FindOneAndDelete find and delete demo ==========")
	demonstrateFindOneAndDelete(ctx, engine)

	fmt.Println("\n========== 6. BulkWrite batch operations demo ==========")
	demonstrateBulkWrite(ctx, engine)

	fmt.Println("\n========== 7. Cursor advanced features demo ==========")
	demonstrateAdvancedCursor(ctx, engine)
}

// 1. Cursor iteration demo
func demonstrateCursor(ctx context.Context, engine *pie.Engine) {
	session := pie.Table[Product](engine)

	// Insert test data
	products := []*Product{
		{Name: "笔记本电脑", Price: 5999, Stock: 10, Category: "电子产品"},
		{Name: "键盘", Price: 299, Stock: 50, Category: "电子产品"},
		{Name: "鼠标", Price: 99, Stock: 100, Category: "电子产品"},
		{Name: "显示器", Price: 1999, Stock: 20, Category: "电子产品"},
	}

	for _, p := range products {
		session.Insert(ctx, p)
	}

	fmt.Println("\nIterate all products using Cursor:")
	cursor, err := session.
		Where("category", pie.Eq("category", "电子产品")).
		FindCursor(ctx)
	if err != nil {
		log.Printf("Failed to get cursor: %v", err)
		return
	}

	// Method 1: Using Next() and Decode()
	fmt.Println("\nMethod 1: Next() + Decode()")
	count := 0
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			log.Printf("Decode failed: %v", err)
			continue
		}
		count++
		fmt.Printf("  %d. %s - ¥%.2f (Stock: %d)\n", count, product.Name, product.Price, product.Stock)
	}
	cursor.Close(ctx)

	// Method 2: Using All() to get all documents at once
	cursor2, _ := session.
		Where("category", pie.Eq("category", "电子产品")).
		FindCursor(ctx)

	fmt.Println("\nMethod 2: All()")
	allProducts, err := cursor2.All(ctx)
	if err != nil {
		log.Printf("Failed to get all documents: %v", err)
		return
	}
	fmt.Printf("  Found %d products\n", len(allProducts))

	// Method 3: Using Iterate() to process
	cursor3, _ := session.
		Where("price", pie.Gt("price", 1000)).
		FindCursor(ctx)

	fmt.Println("\nMethod 3: Iterate()")
	cursor3.Iterate(ctx, func(p *Product) error {
		fmt.Printf("  Processing high-price product: %s - ¥%.2f\n", p.Name, p.Price)
		return nil
	})

	// Method 4: Using Take() to get first N documents
	cursor4, _ := session.
		Where("category", pie.Eq("category", "电子产品")).
		OrderBy("price").
		FindCursor(ctx)

	fmt.Println("\nMethod 4: Take(2) - Get the 2 cheapest products")
	topProducts, err := cursor4.Take(ctx, 2)
	if err != nil {
		log.Printf("Take failed: %v", err)
		return
	}
	for i, p := range topProducts {
		fmt.Printf("  %d. %s - ¥%.2f\n", i+1, p.Name, p.Price)
	}
}

// 2. CountDocuments count demo
func demonstrateCount(ctx context.Context, engine *pie.Engine) {
	session := pie.Table[Product](engine)

	// Precise count
	totalCount, err := session.CountDocuments(ctx)
	if err != nil {
		log.Printf("Count failed: %v", err)
		return
	}
	fmt.Printf("Total products (precise): %d\n", totalCount)

	// Conditional count
	expensiveCount, err := session.
		Where("price", pie.Gt("price", 1000)).
		CountDocuments(ctx)
	if err != nil {
		log.Printf("Conditional count failed: %v", err)
		return
	}
	fmt.Printf("Products with price>1000: %d\n", expensiveCount)

	// Estimated count (fast)
	estimatedCount, err := session.EstimatedDocumentCount(ctx)
	if err != nil {
		log.Printf("Estimated count failed: %v", err)
		return
	}
	fmt.Printf("Estimated product count (fast): %d\n", estimatedCount)
}

// 3. ReplaceOne document replacement demo
func demonstrateReplace(ctx context.Context, engine *pie.Engine) {
	session := pie.Table[Product](engine)

	// Find a product
	product, err := session.
		Where("name", pie.Eq("name", "鼠标")).
		FindOne(ctx)
	if err != nil {
		log.Printf("Failed to find product: %v", err)
		return
	}

	fmt.Printf("Original product: %s - ¥%.2f\n", product.Name, product.Price)

	// Completely replace product information
	newProduct := Product{
		ID:        product.ID, // Keep ID unchanged
		Name:      "无线鼠标",
		Price:     199,
		Stock:     80,
		Category:  "电子产品",
		CreatedAt: product.CreatedAt,
		UpdatedAt: time.Now(),
	}

	result, err := session.
		Where("_id", pie.ID(product.ID)).
		ReplaceOne(ctx, &newProduct)
	if err != nil {
		log.Printf("Replace failed: %v", err)
		return
	}

	fmt.Printf("Replace successful! Matched: %d, Modified: %d\n", result.MatchedCount, result.ModifiedCount)
	fmt.Printf("New product: %s - ¥%.2f\n", newProduct.Name, newProduct.Price)
}

// 4. FindOneAndUpdate find and update demo
func demonstrateFindOneAndUpdate(ctx context.Context, engine *pie.Engine) {
	session := pie.Table[Product](engine)

	// Find and update, return document before update
	fmt.Println("\nReturn document before update:")
	oldProduct, err := session.
		Where("name", pie.Eq("name", "键盘")).
		FindOneAndUpdate(ctx, bson.D{
			{Key: "$set", Value: bson.D{{Key: "price", Value: 399}}},
			{Key: "$inc", Value: bson.D{{Key: "stock", Value: -5}}},
		}, false) // false = return document before update
	if err != nil {
		log.Printf("Find and update failed: %v", err)
		return
	}
	fmt.Printf("  Before update: %s - ¥%.2f (Stock: %d)\n", oldProduct.Name, oldProduct.Price, oldProduct.Stock)

	// Find and update, return document after update
	fmt.Println("\nReturn document after update:")
	newProduct, err := session.
		Where("name", pie.Eq("name", "键盘")).
		FindOneAndUpdate(ctx, bson.D{
			{Key: "$set", Value: bson.D{{Key: "price", Value: 499}}},
		}, true) // true = return document after update
	if err != nil {
		log.Printf("Find and update failed: %v", err)
		return
	}
	fmt.Printf("  After update: %s - ¥%.2f (Stock: %d)\n", newProduct.Name, newProduct.Price, newProduct.Stock)
}

// 5. FindOneAndDelete find and delete demo
func demonstrateFindOneAndDelete(ctx context.Context, engine *pie.Engine) {
	session := pie.Table[Product](engine)

	// Find and delete product with zero stock
	deletedProduct, err := session.
		Where("name", pie.Eq("name", "无线鼠标")).
		FindOneAndDelete(ctx)
	if err != nil {
		log.Printf("Find and delete failed: %v", err)
		return
	}

	fmt.Printf("Deleted product: %s - ¥%.2f\n", deletedProduct.Name, deletedProduct.Price)
}

// 6. BulkWrite batch operations demo
func demonstrateBulkWrite(ctx context.Context, engine *pie.Engine) {
	// Create bulk write operations
	bulkWrite := pie.NewBulkWrite[Product](engine).
		CollectionForStruct(Product{})

	// Insert new products
	bulkWrite.InsertOne(&Product{
		Name:     "平板电脑",
		Price:    3999,
		Stock:    15,
		Category: "电子产品",
	})

	bulkWrite.InsertOne(&Product{
		Name:     "耳机",
		Price:    599,
		Stock:    30,
		Category: "电子产品",
	})

	// Update product price
	bulkWrite.UpdateOne(
		bson.D{{Key: "name", Value: "显示器"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "price", Value: 2299}}}},
	)

	// Batch update stock
	bulkWrite.UpdateMany(
		bson.D{{Key: "stock", Value: bson.D{{Key: "$lt", Value: 20}}}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "stock", Value: 50}}}},
	)

	// Delete low-price products
	bulkWrite.DeleteMany(
		bson.D{{Key: "price", Value: bson.D{{Key: "$lt", Value: 100}}}},
	)

	// Execute batch operations (ordered)
	fmt.Println("\nExecuting batch operations...")
	result, err := bulkWrite.ExecuteOrdered(ctx)
	if err != nil {
		log.Printf("Batch operations failed: %v", err)
		return
	}

	fmt.Println("Batch operation results:")
	fmt.Printf("  Inserted: %d\n", result.InsertedCount)
	fmt.Printf("  Matched: %d\n", result.MatchedCount)
	fmt.Printf("  Modified: %d\n", result.ModifiedCount)
	fmt.Printf("  Deleted: %d\n", result.DeletedCount)
	fmt.Printf("  Upserted: %d\n", result.UpsertedCount)
}

// 7. Cursor advanced features demo
func demonstrateAdvancedCursor(ctx context.Context, engine *pie.Engine) {
	session := pie.Table[Product](engine)

	cursor, err := session.
		Where("category", pie.Eq("category", "电子产品")).
		FindCursor(ctx)
	if err != nil {
		log.Printf("Failed to get cursor: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Set batch size
	cursor.SetBatchSize(2)

	// Check remaining batch length
	fmt.Printf("Remaining documents in current batch: %d\n", cursor.RemainingBatchLength())

	// Use First() to get the first document
	cursor2, _ := session.
		Where("category", pie.Eq("category", "电子产品")).
		OrderByDesc("price").
		FindCursor(ctx)

	firstProduct, err := cursor2.First(ctx)
	if err != nil {
		log.Printf("Failed to get first document: %v", err)
		return
	}
	fmt.Printf("\nMost expensive product: %s - ¥%.2f\n", firstProduct.Name, firstProduct.Price)

	// Use Filter() to filter documents
	cursor3, _ := session.
		Where("category", pie.Eq("category", "电子产品")).
		FindCursor(ctx)

	affordableProducts, err := pie.Filter(ctx, cursor3, func(p *Product) bool {
		return p.Price <= 1000
	})
	if err != nil {
		log.Printf("Filter failed: %v", err)
		return
	}

	fmt.Println("\nProducts with price <= 1000:")
	for _, p := range affordableProducts {
		fmt.Printf("  %s - ¥%.2f\n", p.Name, p.Price)
	}
}
