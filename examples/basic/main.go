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
	ID    bson.ObjectID `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Age   int           `bson:"age"`
	Email string        `bson:"email"`
}

// AgeGroupResult aggregation result for age grouping
type AgeGroupResult struct {
	ID    int `bson:"_id"`   // age value
	Count int `bson:"count"` // count of users
}

func main() {
	// Create MongoDB engine
	engine, err := pie.NewEngine(
		context.Background(),
		"pokerbot",
		pie.WithURI("mongodb://admin:password@localhost:27017/"),
		pie.WithMapper(&pie.SnakeMapper{}),
	)
	if err != nil {
		log.Fatal("Failed to create engine:", err)
	}
	defer engine.Disconnect(context.Background())

	// Create type-safe session
	session := pie.Table[User](engine)

	// Clear existing data
	_, err = session.DeleteMany(context.Background())
	if err != nil {
		log.Fatal("Failed to clear existing data:", err)
	}
	fmt.Println("Cleared existing data")

	// Test insert
	user := &User{
		Name:  "John",
		Age:   25,
		Email: "john@example.com",
	}

	result, err := session.Insert(context.Background(), user)
	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}
	fmt.Printf("Inserted user with ID: %s\n", result.InsertedID)

	// Insert more test data for AND/OR examples
	users := []*User{
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
		{Name: "Bob", Age: 22, Email: "bob@example.com"},
		{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		{Name: "John", Age: 28, Email: "john2@example.com"},
	}

	for _, u := range users {
		_, err := session.Insert(context.Background(), u)
		if err != nil {
			log.Fatal("Failed to insert user:", err)
		}
	}
	fmt.Printf("Inserted %d additional users\n", len(users))

	// Test query
	usersByAge, err := session.
		WhereOperator(pie.Gte("age", 20)).
		OrderBy("name").
		Limit(10).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users:", err)
	}
	fmt.Printf("Found %d users\n", len(usersByAge))

	// Test update
	updateResult, err := session.
		WhereOperator(pie.Eq("name", "John")).
		Update(context.Background(), bson.D{{"$set", bson.D{{"age", 26}}}})
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}
	fmt.Printf("Updated %d users\n", updateResult.ModifiedCount)

	// Test Where method (simple equality)
	usersByName, err := session.
		Where("name", "John").
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users by name:", err)
	}
	fmt.Printf("Found %d users with name 'John' using Where method\n", len(usersByName))

	// Test WhereOperator method (same result)
	usersByNameOp, err := session.
		WhereOperator(pie.Eq("name", "John")).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users by name with operator:", err)
	}
	fmt.Printf("Found %d users with name 'John' using WhereOperator method\n", len(usersByNameOp))

	// Test AND query - multiple conditions
	usersAnd, err := session.
		Where("name", "John").
		WhereOperator(pie.Gte("age", 20)).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users with AND conditions:", err)
	}
	fmt.Printf("Found %d users with name 'John' AND age >= 20\n", len(usersAnd))

	// Test OR query using WhereOperator
	usersOr, err := session.
		WhereOperator(pie.Or(
			pie.Eq("name", "John"),
			pie.Eq("name", "Alice"),
		)).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users with OR conditions:", err)
	}
	fmt.Printf("Found %d users with name 'John' OR 'Alice'\n", len(usersOr))

	// Test complex AND/OR combination
	usersComplex, err := session.
		WhereOperator(pie.And(
			pie.Gte("age", 18),
			pie.Lt("age", 30),
		)).
		WhereOperator(pie.Or(
			pie.Eq("name", "John"),
			pie.Eq("name", "Bob"),
		)).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users with complex conditions:", err)
	}
	fmt.Printf("Found %d users with age 18-30 AND (name 'John' OR 'Bob')\n", len(usersComplex))

	// Show detailed results for complex query
	fmt.Println("Complex query results:")
	for _, user := range usersComplex {
		fmt.Printf("  - Name: %s, Age: %d, Email: %s\n", user.Name, user.Age, user.Email)
	}

	// Test more AND/OR examples
	fmt.Println("\n=== More AND/OR Examples ===")

	// Example 1: Users with specific age range AND specific names
	usersAgeRange, err := session.
		WhereOperator(pie.And(
			pie.Gte("age", 25),
			pie.Lte("age", 35),
		)).
		WhereOperator(pie.Or(
			pie.Eq("name", "John"),
			pie.Eq("name", "Alice"),
			pie.Eq("name", "Charlie"),
		)).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users with age range:", err)
	}
	fmt.Printf("Users aged 25-35 with names 'John', 'Alice', or 'Charlie': %d\n", len(usersAgeRange))
	for _, user := range usersAgeRange {
		fmt.Printf("  - %s (age %d)\n", user.Name, user.Age)
	}

	// Example 2: Using And() method for multiple conditions
	usersAndMethod, err := session.
		And(
			pie.Eq("name", "John"),
			pie.Gte("age", 25),
		).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users with And method:", err)
	}
	fmt.Printf("Users found using And() method: %d\n", len(usersAndMethod))

	// Example 3: Using Or() method for multiple conditions
	usersOrMethod, err := session.
		Or(
			pie.Eq("name", "Bob"),
			pie.Eq("name", "Charlie"),
		).
		Find(context.Background())
	if err != nil {
		log.Fatal("Failed to find users with Or method:", err)
	}
	fmt.Printf("Users found using Or() method: %d\n", len(usersOrMethod))
	for _, user := range usersOrMethod {
		fmt.Printf("  - %s (age %d)\n", user.Name, user.Age)
	}

	// Test aggregation
	aggregateResult, err := pie.NewAggregate[AgeGroupResult](engine).
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
