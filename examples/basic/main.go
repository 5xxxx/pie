package main

import (
	"context"
	"fmt"
	"log"

	"github.com/5xxxx/pie"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// User user structure
type User struct {
	ID    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Email string `bson:"email"`
}

func main() {
	// Create MongoDB engine
	engine, err := pie.NewEngine(
		context.Background(),
		"testdb",
		pie.WithURI("mongodb://localhost:27017"),
		pie.WithMapper(&pie.SnakeMapper{}),
	)
	if err != nil {
		log.Fatal("Failed to create engine:", err)
	}
	defer engine.Disconnect(context.Background())

	// Create type-safe session
	session := pie.Table[User](engine)

	// Test insert
	user := &User{
		Name:  "张三",
		Age:   25,
		Email: "zhangsan@example.com",
	}

	result, err := session.Insert(context.Background(), user)
	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}
	fmt.Printf("Inserted user with ID: %s\n", result.InsertedID)

	// Test query
	users, err := session.
		Where("age", pie.Gte("age", 20)).
		OrderBy("name").
		Limit(10).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users:", err)
	}
	fmt.Printf("Found %d users\n", len(users))

	// Test update
	updateResult, err := session.
		Where("name", "张三").
		Update(context.Background(), bson.D{{"$set", bson.D{{"age", 26}}}})
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}
	fmt.Printf("Updated %d users\n", updateResult.ModifiedCount)

	// Test aggregation
	aggregateResult, err := pie.NewAggregate[User](engine).
		CollectionForStruct(User{}).
		Match(bson.D{{"age", bson.D{{"$gte", 20}}}}).
		Group(bson.D{
			{"_id", "$age"},
			{"count", bson.D{{"$sum", 1}}},
		}).
		Sort(bson.D{{"_id", 1}}).
		Exec(context.Background())
	if err != nil {
		log.Fatal("Failed to execute aggregation:", err)
	}
	fmt.Printf("Aggregation result: %+v\n", aggregateResult.Data)

	fmt.Println("All tests completed successfully!")
}