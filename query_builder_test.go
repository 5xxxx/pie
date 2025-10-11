package pie

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestUserBuilder 用于测试的结构体
type TestUserBuilder struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func TestSessionWhereIn(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereIn("status", []string{"active", "pending"})

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "status" {
		t.Errorf("Expected key 'status', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$in" {
			t.Errorf("Expected $in operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereNotIn(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereNotIn("status", []string{"deleted", "archived"})

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "status" {
		t.Errorf("Expected key 'status', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$nin" {
			t.Errorf("Expected $nin operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereBetween(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereBetween("age", 18, 65)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "age" {
		t.Errorf("Expected key 'age', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(valueDoc))
		}
		hasGte := false
		hasLte := false
		for _, cond := range valueDoc {
			if cond.Key == "$gte" && cond.Value == 18 {
				hasGte = true
			}
			if cond.Key == "$lte" && cond.Value == 65 {
				hasLte = true
			}
		}
		if !hasGte || !hasLte {
			t.Errorf("Expected $gte: 18 and $lte: 65, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereNotBetween(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereNotBetween("age", 18, 65)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$or" {
		t.Errorf("Expected key '$or', got %s", filter[0].Key)
	}
}

func TestSessionWhereNull(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereNull("deleted_at")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$or" {
		t.Errorf("Expected key '$or', got %s", filter[0].Key)
	}
}

func TestSessionWhereNotNull(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereNotNull("name")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$and" {
		t.Errorf("Expected key '$and', got %s", filter[0].Key)
	}
}

func TestSessionWhereExists(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereExists("profile")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "profile" {
		t.Errorf("Expected key 'profile', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$exists" || valueDoc[0].Value != true {
			t.Errorf("Expected $exists: true, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereNotExists(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereNotExists("deleted_at")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "deleted_at" {
		t.Errorf("Expected key 'deleted_at', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$exists" || valueDoc[0].Value != false {
			t.Errorf("Expected $exists: false, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereDate(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	now := time.Now()
	session.WhereDate("created_at", ">", now)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "created_at" {
		t.Errorf("Expected key 'created_at', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$gt" {
			t.Errorf("Expected $gt operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereDateBetween(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()
	session.WhereDateBetween("created_at", start, end)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "created_at" {
		t.Errorf("Expected key 'created_at', got %s", filter[0].Key)
	}
}

func TestSessionWhereMonth(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereMonth("created_at", 12)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$expr" {
		t.Errorf("Expected key '$expr', got %s", filter[0].Key)
	}
}

func TestSessionWhereYear(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereYear("created_at", 2023)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$expr" {
		t.Errorf("Expected key '$expr', got %s", filter[0].Key)
	}
}

func TestSessionWhereRecentDays(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereRecentDays("created_at", 7)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "created_at" {
		t.Errorf("Expected key 'created_at', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$gte" {
			t.Errorf("Expected $gte operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereLike(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereLike("name", "%john%")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "name" {
		t.Errorf("Expected key 'name', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(valueDoc))
		}
		hasRegex := false
		hasOptions := false
		for _, cond := range valueDoc {
			if cond.Key == "$regex" && cond.Value == "john" {
				hasRegex = true
			}
			if cond.Key == "$options" && cond.Value == "i" {
				hasOptions = true
			}
		}
		if !hasRegex || !hasOptions {
			t.Errorf("Expected $regex: 'john' and $options: 'i', got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereStartsWith(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereStartsWith("name", "john")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "name" {
		t.Errorf("Expected key 'name', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(valueDoc))
		}
		hasRegex := false
		hasOptions := false
		for _, cond := range valueDoc {
			if cond.Key == "$regex" && cond.Value == "^john" {
				hasRegex = true
			}
			if cond.Key == "$options" && cond.Value == "i" {
				hasOptions = true
			}
		}
		if !hasRegex || !hasOptions {
			t.Errorf("Expected $regex: '^john' and $options: 'i', got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereEndsWith(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereEndsWith("name", "smith")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "name" {
		t.Errorf("Expected key 'name', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(valueDoc))
		}
		hasRegex := false
		hasOptions := false
		for _, cond := range valueDoc {
			if cond.Key == "$regex" && cond.Value == "smith$" {
				hasRegex = true
			}
			if cond.Key == "$options" && cond.Value == "i" {
				hasOptions = true
			}
		}
		if !hasRegex || !hasOptions {
			t.Errorf("Expected $regex: 'smith$' and $options: 'i', got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereArrayContains(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereArrayContains("tags", "golang")

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "tags" {
		t.Errorf("Expected key 'tags', got %s", filter[0].Key)
	}
	if filter[0].Value != "golang" {
		t.Errorf("Expected value 'golang', got %v", filter[0].Value)
	}
}

func TestSessionWhereArraySize(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereArraySize("tags", 3)

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "tags" {
		t.Errorf("Expected key 'tags', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$size" || valueDoc[0].Value != 3 {
			t.Errorf("Expected $size: 3, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionWhereArrayAll(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.WhereArrayAll("tags", []string{"golang", "mongodb"})

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "tags" {
		t.Errorf("Expected key 'tags', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$all" {
			t.Errorf("Expected $all operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestSessionOrWhere(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.OrWhere(func(q *Query) *Query {
		return q.Where("name", "john").Where("age", 25)
	})

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$or" {
		t.Errorf("Expected key '$or', got %s", filter[0].Key)
	}
}

func TestSessionAndWhere(t *testing.T) {
	engine, _ := NewEngine(context.Background(), "pie-test", WithURI("mongodb://admin:password@localhost:27017/"))
	session := Table[TestUserBuilder](engine)

	session.AndWhere(func(q *Query) *Query {
		return q.Where("status", "active").Where("verified", true)
	})

	filter := session.query.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$and" {
		t.Errorf("Expected key '$and', got %s", filter[0].Key)
	}
}

// Query方法测试
func TestQueryWhereIn(t *testing.T) {
	q := NewQuery()
	q.WhereIn("status", []string{"active", "pending"})

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "status" {
		t.Errorf("Expected key 'status', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$in" {
			t.Errorf("Expected $in operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestQueryWhereNotIn(t *testing.T) {
	q := NewQuery()
	q.WhereNotIn("status", []string{"deleted", "archived"})

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "status" {
		t.Errorf("Expected key 'status', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$nin" {
			t.Errorf("Expected $nin operator, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestQueryWhereBetween(t *testing.T) {
	q := NewQuery()
	q.WhereBetween("age", 18, 65)

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "age" {
		t.Errorf("Expected key 'age', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(valueDoc))
		}
		hasGte := false
		hasLte := false
		for _, cond := range valueDoc {
			if cond.Key == "$gte" && cond.Value == 18 {
				hasGte = true
			}
			if cond.Key == "$lte" && cond.Value == 65 {
				hasLte = true
			}
		}
		if !hasGte || !hasLte {
			t.Errorf("Expected $gte: 18 and $lte: 65, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestQueryWhereNull(t *testing.T) {
	q := NewQuery()
	q.WhereNull("deleted_at")

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$or" {
		t.Errorf("Expected key '$or', got %s", filter[0].Key)
	}
}

func TestQueryWhereNotNull(t *testing.T) {
	q := NewQuery()
	q.WhereNotNull("name")

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "$and" {
		t.Errorf("Expected key '$and', got %s", filter[0].Key)
	}
}

func TestQueryWhereExists(t *testing.T) {
	q := NewQuery()
	q.WhereExists("profile")

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "profile" {
		t.Errorf("Expected key 'profile', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$exists" || valueDoc[0].Value != true {
			t.Errorf("Expected $exists: true, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

func TestQueryWhereNotExists(t *testing.T) {
	q := NewQuery()
	q.WhereNotExists("deleted_at")

	filter := q.GetFilter()
	if len(filter) == 0 {
		t.Fatal("Expected filter to be generated")
	}

	if filter[0].Key != "deleted_at" {
		t.Errorf("Expected key 'deleted_at', got %s", filter[0].Key)
	}
	// 验证值结构
	if valueDoc, ok := filter[0].Value.(bson.D); ok {
		if len(valueDoc) != 1 || valueDoc[0].Key != "$exists" || valueDoc[0].Value != false {
			t.Errorf("Expected $exists: false, got %v", valueDoc)
		}
	} else {
		t.Errorf("Expected bson.D value, got %T", filter[0].Value)
	}
}

// 辅助函数测试
func TestBool(t *testing.T) {
	result := Bool(true)
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}
	if *result != true {
		t.Errorf("Expected true, got %v", *result)
	}

	result = Bool(false)
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}
	if *result != false {
		t.Errorf("Expected false, got %v", *result)
	}
}

func TestFloat64(t *testing.T) {
	result := Float64(3.14)
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}
	if *result != 3.14 {
		t.Errorf("Expected 3.14, got %v", *result)
	}
}

func TestInt(t *testing.T) {
	result := Int(42)
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}
	if *result != 42 {
		t.Errorf("Expected 42, got %v", *result)
	}
}

func TestString(t *testing.T) {
	result := String("hello")
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}
	if *result != "hello" {
		t.Errorf("Expected 'hello', got %v", *result)
	}
}
