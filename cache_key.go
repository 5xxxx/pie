package pie

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// CacheKeyGenerator cache key generator
type CacheKeyGenerator struct {
	prefix string
}

// NewCacheKeyGenerator create cache key generator
func NewCacheKeyGenerator(prefix string) *CacheKeyGenerator {
	return &CacheKeyGenerator{prefix: prefix}
}

// GenerateQueryKey generate query cache key
func (ckg *CacheKeyGenerator) GenerateQueryKey(collection string, filter bson.D, options interface{}) string {
	// Serialize filter and options
	filterJSON, _ := json.Marshal(filter)
	optionsJSON, _ := json.Marshal(options)

	// Generate MD5 hash
	hasher := md5.New()
	hasher.Write([]byte(collection))
	hasher.Write(filterJSON)
	hasher.Write(optionsJSON)
	hash := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("%s:query:%s:%s", ckg.prefix, collection, hash)
}

// GenerateFindOneKey generate FindOne cache key
func (ckg *CacheKeyGenerator) GenerateFindOneKey(collection string, filter bson.D) string {
	filterJSON, _ := json.Marshal(filter)
	hasher := md5.New()
	hasher.Write([]byte(collection))
	hasher.Write(filterJSON)
	hash := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("%s:findone:%s:%s", ckg.prefix, collection, hash)
}

// GenerateCountKey generate Count cache key
func (ckg *CacheKeyGenerator) GenerateCountKey(collection string, filter bson.D) string {
	filterJSON, _ := json.Marshal(filter)
	hasher := md5.New()
	hasher.Write([]byte(collection))
	hasher.Write(filterJSON)
	hash := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("%s:count:%s:%s", ckg.prefix, collection, hash)
}

// GenerateDocumentKey generate document cache key
func (ckg *CacheKeyGenerator) GenerateDocumentKey(collection string, id interface{}) string {
	return fmt.Sprintf("%s:doc:%s:%v", ckg.prefix, collection, id)
}

// GenerateCollectionPattern generate collection pattern (for bulk deletion)
func (ckg *CacheKeyGenerator) GenerateCollectionPattern(collection string) string {
	return fmt.Sprintf("%s:*:%s:*", ckg.prefix, collection)
}

// GenerateTagKey generate tag key
func (ckg *CacheKeyGenerator) GenerateTagKey(tag string) string {
	return fmt.Sprintf("%s:tag:%s", ckg.prefix, tag)
}