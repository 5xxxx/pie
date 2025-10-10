package pie

import (
        "testing"

        "go.mongodb.org/mongo-driver/v2/bson"
)

func TestCacheKeyGenerator(t *testing.T) {
        gen := NewCacheKeyGenerator("prefix")
        filter := bson.D{{Key: "name", Value: "Alice"}}

        key := gen.GenerateQueryKey("users", filter, map[string]int{"limit": 10})
        if key == "" {
                t.Fatalf("expected non-empty key")
        }
        if key != gen.GenerateQueryKey("users", filter, map[string]int{"limit": 10}) {
                t.Fatalf("expected deterministic key generation")
        }

        if findOne := gen.GenerateFindOneKey("users", filter); findOne == "" {
                t.Fatalf("find one key should not be empty")
        }

        if count := gen.GenerateCountKey("users", filter); count == "" {
                t.Fatalf("count key should not be empty")
        }

        if doc := gen.GenerateDocumentKey("users", 42); doc != "prefix:doc:users:42" {
                t.Fatalf("unexpected document key: %s", doc)
        }

        if pattern := gen.GenerateCollectionPattern("users"); pattern != "prefix:*:users:*" {
                t.Fatalf("unexpected pattern: %s", pattern)
        }

        if tag := gen.GenerateTagKey("group"); tag != "prefix:tag:group" {
                t.Fatalf("unexpected tag key: %s", tag)
        }
}

