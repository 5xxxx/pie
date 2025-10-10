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

// UserWithHooks user model, demonstrating hook functionality
type UserWithHooks struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	Age       int           `bson:"age"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

// BeforeCreate before create hook - auto set timestamps
func (u *UserWithHooks) BeforeCreate(ctx context.Context) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	fmt.Printf("🔧 BeforeCreate: Setting timestamps for user %s\n", u.Name)
	return nil
}

// BeforeUpdate before update hook - auto update timestamp
func (u *UserWithHooks) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	fmt.Printf("🔧 BeforeUpdate: Updating timestamp for user %s\n", u.Name)
	return nil
}

// BeforeSave before save hook - data validation
func (u *UserWithHooks) BeforeSave(ctx context.Context) error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !contains(u.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	fmt.Printf("🔧 BeforeSave: Validating user %s\n", u.Name)
	return nil
}

// AfterCreate after create hook - send notification
func (u *UserWithHooks) AfterCreate(ctx context.Context) error {
	fmt.Printf("📧 AfterCreate: Sending welcome email to %s\n", u.Email)
	return nil
}

// AfterFind after find hook - record access
func (u *UserWithHooks) AfterFind(ctx context.Context) error {
	fmt.Printf("👀 AfterFind: User %s was accessed\n", u.Name)
	return nil
}

// Account account model, demonstrating password encryption
type Account struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Username  string        `bson:"username"`
	Password  string        `bson:"password"`
	CreatedAt time.Time     `bson:"created_at"`
}

// BeforeCreate password encryption hook
func (a *Account) BeforeCreate(ctx context.Context) error {
	// Simulate password encryption
	a.Password = "hashed_" + a.Password
	a.CreatedAt = time.Now()
	fmt.Printf("🔒 BeforeCreate: Encrypting password for user %s\n", a.Username)
	return nil
}

// SoftDeleteModel soft delete model
type SoftDeleteModel struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
}

// BeforeDelete soft delete hook
func (s *SoftDeleteModel) BeforeDelete(ctx context.Context) error {
	now := time.Now()
	s.DeletedAt = &now
	fmt.Printf("🗑️ BeforeDelete: Soft deleting %s\n", s.Name)
	return nil
}

func main() {
	ctx := context.Background()

	// 1. Create engine and enable query logging
	fmt.Println("=== Create engine and enable query logging ===")
	engine, err := pie.NewEngine(ctx, "testdb",
		pie.WithQueryLog(os.Stdout),                     // Enable query logging
		pie.WithSlowQueryThreshold(10*time.Millisecond), // Only log queries over 10ms
	)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Disconnect(ctx)

	// 2. Register global hooks
	fmt.Println("\n=== Register global hooks ===")
	engine.Hooks().RegisterAfterCreate(func(ctx context.Context, doc interface{}) error {
		fmt.Printf("📊 Global AfterCreate: Document created: %T\n", doc)
		return nil
	})

	engine.Hooks().RegisterAfterUpdate(func(ctx context.Context, doc interface{}) error {
		fmt.Printf("📊 Global AfterUpdate: Document updated: %T\n", doc)
		return nil
	})

	engine.Hooks().RegisterAfterDelete(func(ctx context.Context, doc interface{}) error {
		fmt.Printf("📊 Global AfterDelete: Document deleted: %T\n", doc)
		return nil
	})

	// 3. User CRUD operations demo
	fmt.Println("\n=== User CRUD operations demo ===")
	session := pie.Table[UserWithHooks](engine)

	// Insert user
	user := &UserWithHooks{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   25,
	}

	fmt.Println("\n--- Insert user ---")
	result, err := session.Insert(ctx, user)
	if err != nil {
		log.Printf("Insert failed: %v", err)
	} else {
		fmt.Printf("✅ Insert successful: %+v\n", result)
	}

	// Query user
	fmt.Println("\n--- Query user ---")
	users, err := session.Where("age", pie.Gte("age", 20)).Find(ctx)
	if err != nil {
		log.Printf("Query failed: %v", err)
	} else {
		fmt.Printf("✅ Query successful: Found %d users\n", len(users))
		for _, u := range users {
			fmt.Printf("  - %s (%s)\n", u.Name, u.Email)
		}
	}

	// Update user
	fmt.Println("\n--- Update user ---")
	updateResult, err := session.Where("name", pie.Eq("name", "Alice")).Update(ctx, bson.D{
		{Key: "$set", Value: bson.D{{Key: "age", Value: 26}}},
	})
	if err != nil {
		log.Printf("Update failed: %v", err)
	} else {
		fmt.Printf("✅ Update successful: %+v\n", updateResult)
	}

	// 4. Account password encryption demo
	fmt.Println("\n=== Account password encryption demo ===")
	accountSession := pie.Table[Account](engine)

	account := &Account{
		Username: "bob",
		Password: "secret123",
	}

	fmt.Println("\n--- Insert account (password auto encrypted) ---")
	accountResult, err := accountSession.Insert(ctx, account)
	if err != nil {
		log.Printf("Account insert failed: %v", err)
	} else {
		fmt.Printf("✅ Account insert successful: %+v\n", accountResult)
		fmt.Printf("🔒 Password encrypted: %s\n", account.Password)
	}

	// 5. Soft delete demo
	fmt.Println("\n=== Soft delete demo ===")
	softDeleteSession := pie.Table[SoftDeleteModel](engine)

	softDeleteDoc := &SoftDeleteModel{
		Name: "Test Document",
	}

	fmt.Println("\n--- Insert document ---")
	softResult, err := softDeleteSession.Insert(ctx, softDeleteDoc)
	if err != nil {
		log.Printf("Insert failed: %v", err)
	} else {
		fmt.Printf("✅ Insert successful: %+v\n", softResult)
	}

	fmt.Println("\n--- Soft delete document ---")
	deleteResult, err := softDeleteSession.Where("name", pie.Eq("name", "Test Document")).Delete(ctx)
	if err != nil {
		log.Printf("Delete failed: %v", err)
	} else {
		fmt.Printf("✅ Delete successful: %+v\n", deleteResult)
		fmt.Printf("🗑️ Soft delete time: %v\n", softDeleteDoc.DeletedAt)
	}

	// 6. Skip hooks demo
	fmt.Println("\n=== Skip hooks demo ===")
	fmt.Println("\n--- Normal insert (execute hooks) ---")
	normalUser := &UserWithHooks{Name: "Normal", Email: "normal@example.com", Age: 30}
	_, err = session.Insert(ctx, normalUser)
	if err != nil {
		log.Printf("Insert failed: %v", err)
	}

	fmt.Println("\n--- Skip hooks insert ---")
	skipUser := &UserWithHooks{Name: "Skip", Email: "skip@example.com", Age: 35}
	_, err = session.SkipHooks().Insert(ctx, skipUser)
	if err != nil {
		log.Printf("Insert failed: %v", err)
	}

	// 7. Query log output demo
	fmt.Println("\n=== Query log output demo ===")
	fmt.Println("The above operations have automatically logged query logs in the following format:")
	fmt.Println("[timestamp] [duration] db.collection.operation(parameters)")
	fmt.Println("For example: [2025-01-10 15:30:45] [5ms] db.users.insertOne({name: \"Alice\", email: \"alice@example.com\"})")

	// 8. Aggregation operations demo
	fmt.Println("\n=== Aggregation operations demo ===")
	aggregate := pie.NewAggregate[UserWithHooks](engine).CollectionForStruct(UserWithHooks{})

	aggregateResult, err := aggregate.Match(bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 20}}}}).
		Group(bson.D{{Key: "_id", Value: "$age"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}).
		Exec(ctx)

	if err != nil {
		log.Printf("Aggregation failed: %v", err)
	} else {
		fmt.Printf("✅ Aggregation successful: %+v\n", aggregateResult)
	}

	fmt.Println("\n=== Demo completed ===")
	fmt.Println("Hook functionality summary:")
	fmt.Println("✅ Auto timestamp setting")
	fmt.Println("✅ Data validation")
	fmt.Println("✅ Password encryption")
	fmt.Println("✅ Soft delete")
	fmt.Println("✅ Global hooks")
	fmt.Println("✅ Query logging")
	fmt.Println("✅ Skip hooks")
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}