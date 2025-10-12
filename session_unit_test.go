package pie

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestSessionMethodCalls 测试Session方法调用
func TestSessionMethodCalls(t *testing.T) {
	// 创建一个基本的Session实例用于测试方法调用
	session := &Session[SessionTestUser]{
		query:   NewQuery(),
		options: NewSessionOptions(),
	}

	// 测试Where方法
	whereSession := session.Where("name", "test")
	if whereSession == nil {
		t.Fatal("Where should return a session")
	}
	if whereSession != session {
		t.Error("Where should return the same session instance")
	}

	// 测试WhereOperator方法
	whereOpSession := session.WhereOperator(Eq("age", 30))
	if whereOpSession == nil {
		t.Fatal("WhereOperator should return a session")
	}

	// 测试And方法
	andSession := session.And(Eq("active", true), Gte("age", 18))
	if andSession == nil {
		t.Fatal("And should return a session")
	}

	// 测试Or方法
	orSession := session.Or(Eq("name", "Alice"), Eq("name", "Bob"))
	if orSession == nil {
		t.Fatal("Or should return a session")
	}

	// 测试OrderBy方法
	orderSession := session.OrderBy("name")
	if orderSession == nil {
		t.Fatal("OrderBy should return a session")
	}

	// 测试OrderByDesc方法
	orderDescSession := session.OrderByDesc("age")
	if orderDescSession == nil {
		t.Fatal("OrderByDesc should return a session")
	}

	// 测试Limit方法
	limitSession := session.Limit(10)
	if limitSession == nil {
		t.Fatal("Limit should return a session")
	}

	// 测试Skip方法
	skipSession := session.Skip(5)
	if skipSession == nil {
		t.Fatal("Skip should return a session")
	}

	// 测试Select方法
	selectSession := session.Select("name", "email")
	if selectSession == nil {
		t.Fatal("Select should return a session")
	}

	// 测试Exclude方法
	excludeSession := session.Exclude("profile")
	if excludeSession == nil {
		t.Fatal("Exclude should return a session")
	}

	// 测试Project方法
	projectSession := session.Project(bson.D{{Key: "name", Value: 1}})
	if projectSession == nil {
		t.Fatal("Project should return a session")
	}

	// 测试Clone方法
	clonedSession := session.Clone()
	if clonedSession == nil {
		t.Fatal("Clone should return a session")
	}
	if clonedSession == session {
		t.Error("Clone should return a different instance")
	}

	// 测试Clear方法
	clearedSession := session.Where("name", "test").Clear()
	if clearedSession == nil {
		t.Fatal("Clear should return a session")
	}

	// 测试GetQuery方法
	query := session.GetQuery()
	if query == nil {
		t.Fatal("GetQuery should return a query")
	}

	// 测试GetOptions方法
	options := session.GetOptions()
	if options == nil {
		t.Fatal("GetOptions should return options")
	}

	// 测试SkipHooks方法
	skipHooksSession := session.SkipHooks()
	if skipHooksSession == nil {
		t.Fatal("SkipHooks should return a session")
	}
	if !skipHooksSession.skipHooks {
		t.Error("SkipHooks should set skipHooks to true")
	}
}

// TestSessionAdvancedOptionsUnit 测试高级查询选项方法
func TestSessionAdvancedOptionsUnit(t *testing.T) {
	session := &Session[SessionTestUser]{
		query:   NewQuery(),
		options: NewSessionOptions(),
	}

	// 测试Hint方法
	hintSession := session.Hint("age_1")
	if hintSession == nil {
		t.Fatal("Hint should return a session")
	}

	// 测试Comment方法
	commentSession := session.Comment("test query")
	if commentSession == nil {
		t.Fatal("Comment should return a session")
	}

	// 测试BatchSize方法
	batchSizeSession := session.BatchSize(10)
	if batchSizeSession == nil {
		t.Fatal("BatchSize should return a session")
	}

	// 测试NoCursorTimeout方法
	noTimeoutSession := session.NoCursorTimeout(true)
	if noTimeoutSession == nil {
		t.Fatal("NoCursorTimeout should return a session")
	}

	// 测试ReturnKey方法
	returnKeySession := session.ReturnKey(true)
	if returnKeySession == nil {
		t.Fatal("ReturnKey should return a session")
	}

	// 测试ShowRecordId方法
	showRecordIdSession := session.ShowRecordId(true)
	if showRecordIdSession == nil {
		t.Fatal("ShowRecordId should return a session")
	}

	// 测试Min方法
	minSession := session.Min(bson.D{{Key: "age", Value: 18}})
	if minSession == nil {
		t.Fatal("Min should return a session")
	}

	// 测试Max方法
	maxSession := session.Max(bson.D{{Key: "age", Value: 65}})
	if maxSession == nil {
		t.Fatal("Max should return a session")
	}

	// 测试ArrayFilters方法
	arrayFiltersSession := session.ArrayFilters([]any{bson.D{{Key: "tag", Value: "admin"}}})
	if arrayFiltersSession == nil {
		t.Fatal("ArrayFilters should return a session")
	}

	// 测试Let方法
	letSession := session.Let(bson.D{{Key: "newAge", Value: 30}})
	if letSession == nil {
		t.Fatal("Let should return a session")
	}

	// 测试Upsert方法
	upsertSession := session.Upsert(true)
	if upsertSession == nil {
		t.Fatal("Upsert should return a session")
	}
}

// TestSessionErrorHandlingUnit 测试错误处理
func TestSessionErrorHandlingUnit(t *testing.T) {
	// 测试初始化错误
	session := &Session[SessionTestUser]{
		query:   NewQuery(),
		options: NewSessionOptions(),
		initErr: fmt.Errorf("test initialization error"),
	}

	ctx := context.Background()

	// 测试所有方法都应该返回初始化错误
	methods := []struct {
		name string
		fn   func() error
	}{
		{"FindOne", func() error {
			_, err := session.FindOne(ctx)
			return err
		}},
		{"Find", func() error {
			_, err := session.Find(ctx)
			return err
		}},
		{"Insert", func() error {
			_, err := session.Insert(ctx, &SessionTestUser{})
			return err
		}},
		{"InsertMany", func() error {
			_, err := session.InsertMany(ctx, []SessionTestUser{})
			return err
		}},
		{"Update", func() error {
			_, err := session.Update(ctx, bson.D{})
			return err
		}},
		{"UpdateMany", func() error {
			_, err := session.UpdateMany(ctx, bson.D{})
			return err
		}},
		{"Delete", func() error {
			_, err := session.Delete(ctx)
			return err
		}},
		{"DeleteMany", func() error {
			_, err := session.DeleteMany(ctx)
			return err
		}},
		{"Count", func() error {
			_, err := session.Count(ctx)
			return err
		}},
		{"CountDocuments", func() error {
			_, err := session.CountDocuments(ctx)
			return err
		}},
		{"EstimatedDocumentCount", func() error {
			_, err := session.EstimatedDocumentCount(ctx)
			return err
		}},
		{"Distinct", func() error {
			_, err := session.Distinct(ctx, "name")
			return err
		}},
		{"FindCursor", func() error {
			_, err := session.FindCursor(ctx)
			return err
		}},
		{"ReplaceOne", func() error {
			_, err := session.ReplaceOne(ctx, &SessionTestUser{})
			return err
		}},
		{"FindOneAndDelete", func() error {
			_, err := session.FindOneAndDelete(ctx)
			return err
		}},
		{"FindOneAndReplace", func() error {
			_, err := session.FindOneAndReplace(ctx, &SessionTestUser{}, false)
			return err
		}},
		{"FindOneAndUpdate", func() error {
			_, err := session.FindOneAndUpdate(ctx, bson.D{}, false)
			return err
		}},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			err := method.fn()
			if err == nil {
				t.Errorf("%s should return initialization error", method.name)
			}
			if err.Error() != "session initialization failed: test initialization error" {
				t.Errorf("%s should return initialization error, got: %v", method.name, err)
			}
		})
	}
}

// TestSessionCacheMethods 测试缓存方法
func TestSessionCacheMethods(t *testing.T) {
	session := &Session[SessionTestUser]{
		query:   NewQuery(),
		options: NewSessionOptions(),
	}

	// 测试Cache方法
	cacheSession := session.Cache(5 * time.Minute)
	if cacheSession == nil {
		t.Fatal("Cache should return a session")
	}
	if cacheSession.cacheConfig == nil {
		t.Error("Cache should set cacheConfig")
	}
	if cacheSession.cacheConfig.TTL != 5*time.Minute {
		t.Error("Cache should set correct TTL")
	}

	// 测试NoCache方法
	noCacheSession := session.NoCache()
	if noCacheSession == nil {
		t.Fatal("NoCache should return a session")
	}
	if noCacheSession.cacheConfig == nil {
		t.Error("NoCache should set cacheConfig")
	}
	if noCacheSession.cacheConfig.Enabled {
		t.Error("NoCache should disable caching")
	}

	// 测试CacheWithTags方法
	tagsSession := session.CacheWithTags("user", "profile")
	if tagsSession == nil {
		t.Fatal("CacheWithTags should return a session")
	}
	if tagsSession.cacheConfig == nil {
		t.Error("CacheWithTags should set cacheConfig")
	}
	if len(tagsSession.cacheConfig.Tags) != 2 {
		t.Error("CacheWithTags should set correct tags")
	}

	// 测试CacheWithJitter方法
	jitterSession := session.CacheWithJitter(5*time.Minute, 1*time.Minute)
	if jitterSession == nil {
		t.Fatal("CacheWithJitter should return a session")
	}
	if jitterSession.cacheConfig == nil {
		t.Error("CacheWithJitter should set cacheConfig")
	}
	if !jitterSession.cacheConfig.UseJitter {
		t.Error("CacheWithJitter should set UseJitter")
	}

	// 测试CacheEmpty方法
	emptySession := session.CacheEmpty(5 * time.Minute)
	if emptySession == nil {
		t.Fatal("CacheEmpty should return a session")
	}
	if emptySession.cacheConfig == nil {
		t.Error("CacheEmpty should set cacheConfig")
	}
	if !emptySession.cacheConfig.CacheEmpty {
		t.Error("CacheEmpty should enable empty caching")
	}

	// 测试CacheL1Only方法
	// l1Session := session.CacheL1Only()
	// if l1Session == nil {
	// 	t.Fatal("CacheL1Only should return a session")
	// }
	// if l1Session.cacheConfig == nil {
	// 	t.Error("CacheL1Only should set cacheConfig")
	// }
	// if !l1Session.cacheConfig.Enabled {
	// 	t.Error("CacheL1Only should enable caching")
	// }

	// 测试CacheL2Only方法
	// l2Session := session.CacheL2Only()
	// if l2Session == nil {
	// 	t.Fatal("CacheL2Only should return a session")
	// }
	// if l2Session.cacheConfig == nil {
	// 	t.Error("CacheL2Only should set cacheConfig")
	// }
	// if !l2Session.cacheConfig.Enabled {
	// 	t.Error("CacheL2Only should enable caching")
	// }
}

// TestSessionChaining 测试方法链式调用
func TestSessionChaining(t *testing.T) {
	session := &Session[SessionTestUser]{
		query:   NewQuery(),
		options: NewSessionOptions(),
	}

	// 测试复杂的链式调用
	chainedSession := session.
		Where("active", true).
		Where("age", bson.D{{Key: "$gte", Value: 18}}).
		OrderBy("name").
		OrderByDesc("age").
		Limit(10).
		Skip(5).
		Select("name", "email").
		Cache(5 * time.Minute).
		Hint("age_1").
		Comment("complex query")

	if chainedSession == nil {
		t.Fatal("Chained session should not be nil")
	}

	// 验证查询构建器状态
	query := chainedSession.GetQuery()
	if query == nil {
		t.Fatal("Query should not be nil")
	}

	// 验证选项状态
	options := chainedSession.GetOptions()
	if options == nil {
		t.Fatal("Options should not be nil")
	}

	// 验证缓存配置
	if chainedSession.cacheConfig == nil {
		t.Error("Cache config should be set")
	}
}

// TestNewSessionOptions 测试NewSessionOptions函数
func TestNewSessionOptions(t *testing.T) {
	options := NewSessionOptions()
	if options == nil {
		t.Fatal("NewSessionOptions should return options")
	}

	// 验证默认值
	if options.FindOptions == nil {
		t.Error("FindOptions should be initialized")
	}
	if options.FindOneOptions == nil {
		t.Error("FindOneOptions should be initialized")
	}
	if options.UpdateOneOptions == nil {
		t.Error("UpdateOneOptions should be initialized")
	}
	if options.UpdateManyOptions == nil {
		t.Error("UpdateManyOptions should be initialized")
	}
	if options.DeleteOneOptions == nil {
		t.Error("DeleteOneOptions should be initialized")
	}
	if options.DeleteManyOptions == nil {
		t.Error("DeleteManyOptions should be initialized")
	}
	if options.InsertOptions == nil {
		t.Error("InsertOptions should be initialized")
	}
}
