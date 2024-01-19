package internal

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/5xxxx/pie/driver"
	"github.com/5xxxx/pie/names"
	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// defaultClient is a struct that represents a default client.
type defaultClient struct {
	client     *mongo.Client
	parser     *driver.Parser
	db         string
	clientOpts []*options.ClientOptions
}

// NewClient creates a new client with the specified database name and options.
// It returns a driver.Client interface and an error.
func NewClient(db string, opts ...*options.ClientOptions) (driver.Client, error) {
	mapper := names.NewCacheMapper(new(names.SnakeMapper))
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	parser := driver.NewParser(mapper, mapper)
	d := defaultClient{
		clientOpts: opts,
		parser:     parser,
		client:     client,
		db:         db,
	}
	return &d, nil
}

// Connect connects to the MongoDB server using the provided context.
// If a context is not provided, it uses a background context.
// It sets the client field of the defaultClient instance to the connected client.
// It returns any error that occurs during the connection process.
func (d *defaultClient) Connect(ctx ...context.Context) (err error) {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	d.client, err = mongo.Connect(c, d.clientOpts...)
	return err
}

func (d *defaultClient) Disconnect(ctx ...context.Context) error {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}
	return d.client.Disconnect(c)
}

// FindPagination executes a find command with pagination and returns an error.
// It takes the page number, the number of documents per page, the document type,
// and an optional context as parameters. It then calls the FindPagination method
// of the defaultClient's underlying session with the given parameters and returns
// the result.
func (d *defaultClient) FindPagination(page, count int64, doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindPagination(page, count, doc, ctx...)
}

// BulkWrite executes multiple write operations in bulk and returns a BulkWriteResult.
// It takes in a slice of documents (docs) and optional context(s) (ctx).
// The function creates a new session and calls the BulkWrite method on the session passing the provided parameters.
// It returns the BulkWriteResult and an error (if any).
func (d *defaultClient) BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error) {
	return d.NewSession().BulkWrite(docs, ctx...)
}

// Distinct executes the distinct command and returns an array of distinct values for the specified column(s) in the collection.
// It takes the document as input, which specifies the query criteria, and the columns string, which specifies the column(s) to retrieve distinct values from.
// It also takes an optional context.Context parameter for additional context options.
// The method returns an array of distinct values as []interface{} and an error if any.
// Example usage:
//
//	collection := client.Database("mydb").Collection("mycollection")
//	distinctValues, err := client.Distinct(context.TODO(), bson.D{{"status", "A"}}, "name")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(distinctValues)
//	// Output: [John, Jane, Lisa]
//
// Note: This method requires the MongoDB server version 3.2 or above.
func (d *defaultClient) Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error) {
	return d.NewSession().Distinct(doc, columns, ctx...)
}

// ReplaceOne replaces a single document in the collection.
// It takes a document and an optional context as parameters.
// The method returns an UpdateResult and an error, indicating the success of the operation.
// The UpdateResult contains information about the modification, such as the number of documents matched and modified.
// The ReplaceOne method uses the NewSession method of the defaultClient to create a new session and call the ReplaceOne method on it.
func (d *defaultClient) ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().ReplaceOne(doc, ctx...)
}

// UpdateOneBson updates a single document in the collection specified by 'coll' with the values provided in 'bson'.
// 'coll' is the collection where the document needs to be updated.
// 'bson' is the new values to update in the document.
// 'ctx' is an optional list of context.Context for cancellation or timeout.
// It returns the *mongo.UpdateResult containing the update information and any error that occurred during the operation.
func (d *defaultClient) UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateOneBson(coll, bson, ctx...)
}

// FindOneAndUpdateBson executes a find and update command using the provided bson filter and updates a single document in the specified collection.
// It returns a SingleResult containing the updated document if found, and an error if the command fails.
// The optional context parameter can be used to pass a context.Context instance for cancellation or deadline.
// Example usage:
//
//	coll := "mycollection"
//	filter := bson.M{"name": "John Doe"}
//	update := bson.M{"$set": bson.M{"age": 30}}
//	result, err := client.FindOneAndUpdateBson(coll, filter, update, ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	var doc bson.M
//	if err := result.Decode(&doc); err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(doc)
//	// Output: {"_id": ObjectId("60aaa5d3f91d0d23d2895e11"), "name": "John Doe", "age": 30}
func (d *defaultClient) FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {
	return d.NewSession().FindOneAndUpdateBson(coll, bson, ctx...)
}

// UpdateManyBson executes an update command and returns a mongo.UpdateResult,
// which provides details about the update operation.
// It updates multiple documents in the specified collection with the specified BSON.
// The BSON must adhere to the MongoDB update format.
// The method also accepts optional context.Context parameter(s) to allow cancellation or timeout of the operation.
func (d *defaultClient) UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateManyBson(coll, bson, ctx...)
}

// FindOneAndReplace executes a findAndModify command that replaces one document in the collectionByName.
// The method takes in the document to be used for replacement and an optional context.
// It returns an error if the findAndModify command encounters any issues or if no document is found to replace.
func (d *defaultClient) FindOneAndReplace(doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindOneAndReplace(doc, ctx...)
}

// FindOneAndUpdate executes a find and update command on the collectionByName and returns a SingleResult for the updated document.
// The 'doc' parameter specifies the filter for finding the document to update.
// The 'ctx' parameter provides optional context for the operation.
// It returns a pointer to SingleResult that represents the updated document and an error if any occurred.
func (d *defaultClient) FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {
	return d.NewSession().FindOneAndUpdate(doc, ctx...)
}

// FindAndDelete executes a find command to delete one document in the collectionByName
// and returns an error if any occurred during the operation.
func (d *defaultClient) FindAndDelete(doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindAndDelete(doc, ctx...)
}

// FindOne retrieves a single document from the database that matches the specified filter criteria.
// The retrieved document is stored in the `doc` parameter.
// The `ctx` parameter is an optional context.Context object for controlling the operation.
// It returns an error if the operation fails.
// Example:
//
//	client := NewDefaultClient()
//	var result User
//	err := client.FindOne(&result, context.Background())
//	if err != nil {
//	    fmt.Println("Failed to find document:", err)
//	} else {
//	    fmt.Println("Found document:", result)
//	}
//
// Note: This method is a shorthand for `NewSession().FindOne(doc, ctx...)`.
func (d *defaultClient) FindOne(doc interface{}, ctx ...context.Context) error {
	return d.NewSession().FindOne(doc, ctx...)
}

// FindAll executes a find command and populates the provided interface with multiple documents from the collectionByName.
func (d *defaultClient) FindAll(docs interface{}, ctx ...context.Context) error {
	return d.NewSession().FindAll(docs, ctx...)
}

// FilterBson filters the BSON document `x` and returns a new session with the filter applied.
func (d *defaultClient) FilterBson(x bson.D) driver.Session {
	return d.NewSession().FilterBson(x)
}

// Soft sets the soft session flag, which enables/disables the use of soft session.
// When soft session is enabled, the session is marked as "soft" and sensitive operations like
// database deletions are temporarily disabled. This can be helpful for simulating safe
// operations without affecting the actual data.
// The method returns a new session with the soft flag set to the specified value.
func (d *defaultClient) Soft(s bool) driver.Session {
	return d.NewSession().Soft(s)
}

// RegexFilter applies a regular expression filter to the query by matching the given key with the provided pattern.
// It returns a new session with the applied filter.
// The key parameter specifies the field on which the regular expression filter is applied.
// The pattern parameter specifies the regular expression pattern to match against the key.
func (d *defaultClient) RegexFilter(key, pattern string) driver.Session {
	return d.NewSession().RegexFilter(key, pattern)
}

// Asc returns a new session with sort ascending order applied to the collectionByName.
// It accepts one or more colNames (string) as arguments to specify the columns to sort.
// Example usage:
//
//	client := &defaultClient{}
//	session := client.Asc("name", "age")
//	// Use the session for further operations
//
// Note: The order of columns provided determines the priority of sorting.
//
//	The first column specified will be sorted first, followed by the second column, and so on.
//	If multiple documents have the same value for the first column,
//	the second column will be used to determine the order, and so on.
func (d *defaultClient) Asc(colNames ...string) driver.Session {
	return d.NewSession().Asc(colNames...)
}

// Eq generates a new session with an equality filter applied to the specified key and value.
// The generated session can be further used to execute find commands with the equality filter.
func (d *defaultClient) Eq(key string, value interface{}) driver.Session {
	return d.NewSession().Eq(key, value)
}

// Ne is a method that constructs a "not equal" query filter
// using the given key and value. It returns a driver.Session
// that can be used to execute the query.
//
// Example:
//
//	session := client.Ne("age", 30)
//	session.Execute()
//
// Parameters:
// - key: the key to query on
// - ne: the value that the key should not be equal to
//
// Returns:
// - driver.Session: a session that can be used to execute the query.
func (d *defaultClient) Ne(key string, ne interface{}) driver.Session {
	return d.NewSession().Gt(key, ne)
}

// Nin returns a new driver.Session where the specified key does not match the specified value(s).
// It executes a nin command and returns the resulting session.
// This can be used to exclude documents from a query based on the values of a specific field.
// The key parameter specifies the field key on which the nin operation will be performed.
// The nin parameter specifies the value(s) that should not match the field key.
// Example usage:
//
//	session := client.Nin("status", []string{"completed", "cancelled"})
//	// Use the session object for further operations
//	...
//	// Execute the session
//	err := session.Execute()
//	if err != nil {
//	  // Handle error
//	}
//	...
//
// Note: The value provided for the nin parameter can be of any type, as long as it matches the key's type in MongoDB.
// If the value is of a different type, an error may occur during the execution of the session.
func (d *defaultClient) Nin(key string, nin interface{}) driver.Session {
	return d.NewSession().Nin(key, nin)
}

// Nor constructs a negation condition and returns a new driver.Session with the negation condition applied.
func (d *defaultClient) Nor(c driver.Condition) driver.Session {
	return d.NewSession().Nor(c)
}

// Exists checks whether a key exists in the specified collection. It returns a Session object for further operations.
// The key parameter specifies the key to check for existence.
// The exists parameter specifies whether the key should exist or not.
// The filter parameter (optional) specifies additional conditions to apply when checking for existence.
// The Session object returned can be used to perform further operations on the collection.
// Example usage:
//
//	client := NewDefaultClient()
//	session := client.Exists("myKey", true)
//	// Perform operations on the session
//	...
//	session.Close()
//	// Close the session to release resources
func (d *defaultClient) Exists(key string, exists bool, filter ...driver.Condition) driver.Session {
	return d.NewSession().Exists(key, exists, filter...)
}

// Type executes a GT command with the given key and value and returns a driver.Session.
// This method is used to filter the results of a find command based on a type field in the documents.
func (d *defaultClient) Type(key string, t interface{}) driver.Session {
	return d.NewSession().Gt(key, t)
}

// Expr creates and returns a new session with the given filter expression.
// The session can be used to execute operations using the given filter condition.
// The filter parameter specifies the condition to be used for the session.
// Returns a Session object that allows executing operations using the provided filter.
func (d *defaultClient) Expr(filter driver.Condition) driver.Session {
	return d.NewSession().Expr(filter)
}

// Regex constructs a regular expression using the specified key and value
// and returns a Session with the regular expression applied.
func (d *defaultClient) Regex(key string, value string) driver.Session {
	return d.NewSession().Regex(key, value)
}

// DataBase returns a reference to the MongoDB database that is being used by the default client.
func (d *defaultClient) DataBase() *mongo.Database {
	return d.client.Database(d.db)
}

// Collection returns a new Collection with the specified name and options.
// If a database name is provided, it will use that database; otherwise, it will use the default database of the defaultClient.
// The collOpts parameter is optional and allows for specifying additional collection options.
// It returns a *mongo.Collection.
func (d *defaultClient) Collection(name string, collOpts []*options.CollectionOptions, db ...string) *mongo.Collection {
	var database = d.db
	if len(db) > 0 && len(db[0]) > 0 {
		database = db[0]
	}

	return d.client.Database(database).Collection(name, collOpts...)
}

// Ping pings the database server and returns an error if the ping fails.
// The ping is executed using the underlying MongoClient's Ping method.
// It uses the context.TODO() as the context and readpref.Primary() as the read preference.
// An error is returned if the ping fails.
func (d *defaultClient) Ping() error {
	return d.client.Ping(context.TODO(), readpref.Primary())
}

// Filter adds a filter to the current session's command.
// This filter will be applied when executing a command that reads from the database.
// The filter is specified by providing a key and a value.
// Only documents matching the filter will be returned.
// Example:
//
//	client.Filter("age", 30)
//
// This will only return documents where the "age" field is equal to 30.
func (d *defaultClient) Filter(key string, value interface{}) driver.Session {
	return d.NewSession().Filter(key, value)
}

// ID returns a driver.Session with the specified ID. The driver.Session object
// represents a session with a specific ID that can be used to perform various
// operations on the MongoDB database.
//
// Parameters:
// - id: The ID of the session.
//
// Returns:
// A driver.Session object with the specified ID.
func (d *defaultClient) ID(id interface{}) driver.Session {
	return d.NewSession().ID(id)
}

// Gt executes a find command with a greater than condition and returns a new driver.Session.
// The condition is applied to the specified key in the document.
// The value parameter specifies the value that the key should be greater than.
// Example usage:
//
//	err := client.Gt("age", 18).FindOne(&result)
//
// This will find a document where the "age" key is greater than 18
func (d *defaultClient) Gt(key string, value interface{}) driver.Session {
	return d.NewSession().Gt(key, value)
}

// Gte creates a query filter for the "greater than or equal to" comparison operator.
// It takes a key string and a value interface{} as parameters.
// The key represents the field to compare with, and the value is the value to compare against.
// The method returns a Session object with the query filter applied.
// Example usage:
//
//	client.Gte("age", 30)
//
// This will create a query filter where the "age" field must be greater than or equal to 30.
func (d *defaultClient) Gte(key string, value interface{}) driver.Session {
	return d.NewSession().Gte(key, value)
}

// Lt returns a new session with a query filter that matches documents where the value of the specified key is less than the given value.
func (d *defaultClient) Lt(key string, value interface{}) driver.Session {
	return d.NewSession().Lt(key, value)
}

// Lte returns a new session with the query filter that matches documents where the value of the specified key is less than or equal to the given value.
// The returned session can be used to execute further operations using the Lte operator in the query filter.
func (d *defaultClient) Lte(key string, value interface{}) driver.Session {
	return d.NewSession().Lte(key, value)
}

// In sets a key-value pair in the session context and returns the updated session.
// It is used to pass arbitrary data among method calls in the session chain.
// The key is a string that identifies the data, and the value is the value to be set.
// The updated session is returned to allow for method chaining.
// Example:
//
//	client := &defaultClient{}
//	session := client.In("userID", "123").In("isAdmin", true)
//
// In this example, "userID" is the key and "123" is the value.
// "isAdmin" is another key and the value is a boolean true.
// The In method is called twice to set two different key-value pairs.
// The updated session is then stored in the session variable for further method chaining.
func (d *defaultClient) In(key string, value interface{}) driver.Session {
	return d.NewSession().In(key, value)
}

// And appends an additional filter to the existing filter list in the session
// and returns the updated session.
//
// filter is the driver.Condition to be added to the existing filter list in the session.
//
// Example:
// addFilter := driver.Eq("name", "John")
// newSession := d.And(addFilter)
//
// In the above example, the "name" field should be equal to "John" for the document
// to match the filter. The newSession variable will contain the updated filter list.
//
// Note:
// The And method does not modify the existing filter list in the session,
// but creates a new session with the augmented filter list.
//
// Returns:
// The updated session with the additional filter applied.
func (d *defaultClient) And(filter driver.Condition) driver.Session {
	return d.NewSession().And(filter)
}

// Not excludes documents from a find command that have the specified key-value pair.
// This method returns a new Session with the exclusion applied.
// The new Session will return results excluding documents with the given key-value pair.
func (d *defaultClient) Not(key string, value interface{}) driver.Session {
	return d.NewSession().Not(key, value)
}

// Or adds an additional filter condition to the current session using the OR logical operator.
// It returns a new session with the added filter condition.
func (d *defaultClient) Or(filter driver.Condition) driver.Session {
	return d.NewSession().Or(filter)
}

// InsertOne inserts a single document into the collectionByName.
// It returns the generated ObjectID for the inserted document and an error, if any.
// The document to be inserted is passed as the "v" parameter.
// An optional context can be provided using the "ctx" parameter.
// Example usage:
//
//	doc := bson.M{"name": "John Doe", "age": 30}
//	objID, err := client.InsertOne(doc)
func (d *defaultClient) InsertOne(v interface{}, ctx ...context.Context) (primitive.ObjectID, error) {
	return d.NewSession().InsertOne(v, ctx...)
}

// InsertMany inserts multiple documents into the collection.
//
// It takes a slice of documents and an optional context. If a context is not provided,
// the default context will be used.
//
// The method returns an InsertManyResult, which contains information about the
// success or failure of the operation along with the inserted document IDs.
//
// Example usage:
//
//	result, err := client.InsertMany(docs)
//	if err != nil {
//	    fmt.Println("Failed to insert documents:", err)
//	    return
//	}
//	fmt.Println("Inserted", len(result.InsertedIDs), "documents")
//
// Parameters:
//   - v: A slice of documents to be inserted into the collection.
//   - ctx: A variadic parameter of type context.Context that allows passing
//     additional context options to the operation.
func (d *defaultClient) InsertMany(v interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error) {
	return d.NewSession().InsertMany(v, ctx...)
}

// Limit sets the maximum number of documents that the query will return.
// It takes an integer parameter `limit` and returns a `driver.Session`
// that allows further query operations with the specified limit.
func (d *defaultClient) Limit(limit int64) driver.Session {
	return d.NewSession().Limit(limit)
}

// Skip returns a new session with the number of documents to skip set to the provided value.
// The skip parameter determines the number of documents to skip before starting to return documents.
// It creates a new session using NewSession() method, and then sets the limit on the session using Limit() method with the skip value.
// The session is returned as a driver.Session interface.
// Example usage:
//
//	session := d.Skip(10)
//	// use the session to perform operations on the database
//
// Note: The Skip method is specific to the defaultClient type and can only be called on instances of that type.
func (d *defaultClient) Skip(skip int64) driver.Session {
	return d.NewSession().Limit(skip)
}

// Count executes a count command and returns the number of documents that match the provided filter in the collection.
// It takes an interface{} variable `i` which represents the filter for counting documents.
// The function optionally accepts a variable number of context.Context `ctx` arguments for customizing the count operation.
// It returns the count of matched documents as an int64 and an error, if any.
// The count operation is performed using a NewSession method of the defaultClient struct.
// If there is an error executing the count command, the error is returned.
func (d *defaultClient) Count(i interface{}, ctx ...context.Context) (int64, error) {
	return d.NewSession().Count(i, ctx...)
}

// Desc creates a new session with the default client and calls the Desc method on it.
// It accepts optional session options as variadic input arguments.
// It returns a driver.Session instance that has the Desc method called on it.
func (d *defaultClient) Desc(s2 ...string) driver.Session {
	return d.NewSession().Desc(s2...)
}

// Update executes an update command and returns an UpdateResult
// for the modified document in the collectionByName.
//
// The 'bean' argument represents the document to be updated.
// The optional 'ctx' argument allows specifying additional context.
// It returns the UpdateResult and any error encountered during the update.
func (d *defaultClient) Update(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateOne(bean, ctx...)
}

// UpdateMany updates multiple documents in the collection.
// It takes a bean as the first parameter, which represents the document(s) to be updated.
// The optional ctx parameter can be used to pass a context.Context for cancellation or timeouts.
// It returns an UpdateResult which provides information about the performed update operation.
func (d *defaultClient) UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
	return d.NewSession().UpdateMany(bean, ctx...)
}

// DeleteOne executes a delete command and returns a DeleteResult for one document in the collection.
func (d *defaultClient) DeleteOne(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
	return d.NewSession().DeleteOne(filter, ctx...)
}

// DeleteMany executes a delete command to delete multiple documents in the collectionByName.
// It takes a filter parameter to specify the documents to be deleted.
// The method returns a DeleteResult, which provides information about the deletion operation,
// and an error if any error occurred during the deletion process.
func (d *defaultClient) DeleteMany(filter interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
	return d.NewSession().DeleteMany(filter, ctx...)
}

// SoftDeleteOne executes a soft delete command for a single document in the collectionByName.
// The document to be deleted is specified by the provided filter.
// The method returns an error if the delete command fails.
// An optional context can be passed to modify the behavior of the delete command.
func (d *defaultClient) SoftDeleteOne(filter interface{}, ctx ...context.Context) error {
	return d.NewSession().SoftDeleteOne(filter, ctx...)
}

// SoftDeleteMany executes a soft delete command and returns an error.
// It soft deletes multiple documents in the collectionByName based on the specified filter.
// The documents are marked as deleted, but not physically removed from the collection.
// The operation can be customized by passing optional context parameters.
// The method returns an error if the soft delete operation encounters any issues.
func (d *defaultClient) SoftDeleteMany(filter interface{}, ctx ...context.Context) error {
	return d.NewSession().SoftDeleteMany(filter, ctx...)
}

// FilterBy filters a session using the provided object.
// It creates a new session using the NewSession method of the defaultClient instance.
// It then calls the FilterBy method of the created session, passing the given object.
// It returns the filtered session.
// Example:
//
//	obj := "some object"
//	session := client.FilterBy(obj)
//	session.DoSomething()
//	// ...
//
// Note: The provided object must be compatible with the FilterBy method of session.
func (d *defaultClient) FilterBy(object interface{}) driver.Session {
	return d.NewSession().FilterBy(object)
}

// DropAll drops all indexes in the collection.
// It returns an error if there was a problem dropping the indexes.
func (d *defaultClient) DropAll(doc interface{}, ctx ...context.Context) error {
	return d.NewIndexes().DropAll(doc, ctx...)
}

// DropOne drops one index from the collectionByName using the specified name as identifier.
func (d *defaultClient) DropOne(doc interface{}, name string, ctx ...context.Context) error {
	return d.NewIndexes().DropOne(doc, name, ctx...)
}

// AddIndex adds an index to the collection using the provided keys and options.
// It delegates the operation to the AddIndex method of the NewIndexes interface
// returned by the NewIndexes method of the defaultClient instance.
// It returns the driver.Indexes interface that allows chaining additional index operations.
func (d *defaultClient) AddIndex(keys interface{}, opt ...*options.IndexOptions) driver.Indexes {
	return d.NewIndexes().AddIndex(keys, opt...)
}

// NewIndexes returns a driver.Indexes implementation.
// It creates a new instance of the index struct with the provided driver.Client.
// The index struct is used to perform index-related operations on the collection.
// Example usage:
//
//	func (d *defaultClient) DropAll(doc interface{}, ctx ...context.Context) error {
//	    return d.NewIndexes().DropAll(doc, ctx...)
//	}
//	func (d *defaultClient) DropOne(doc interface{}, name string, ctx ...context.Context) error {
//	    return d.NewIndexes().DropOne(doc, name, ctx...)
//	}
//	func (d *defaultClient) AddIndex(keys interface{}, opt ...*options.IndexOptions) driver.Indexes {
//	    return d.NewIndexes().AddIndex(keys, opt...)
//	}
//
// This method is declared as:
//
//	func NewIndexes(driver driver.Client) driver.Indexes {
//	    return &index{engine: driver}
//	}
func (d *defaultClient) NewIndexes() driver.Indexes {
	return NewIndexes(d)
}

// NewSession creates a new session using the provided defaultClient instance.
// It calls the NewSession function passing the defaultClient instance as the parameter.
// It returns a driver.Session instance which represents the new session.
// Note: The NewSession function must be implemented elsewhere and it is responsible for creating and initializing the session.
// If an error occurs during the creation of the session, it is the responsibility of the caller to handle it.
// Example usage:
//
//	session := defaultClient.NewSession()
//	// Use the session for database operations
//	session.Close()
//	// Clean up the session after use
//
// Note: It is important to always close the session after use to prevent resource leaks.
func (d *defaultClient) NewSession() driver.Session {
	return NewSession(d)
}

// Aggregate returns a new instance of the driver.Aggregate interface.
// It creates and returns a new Aggregate object using the provided defaultClient instance.
// The created Aggregate object can be used to perform aggregation operations in MongoDB.
func (d *defaultClient) Aggregate() driver.Aggregate {
	return NewAggregate(d)
}

// CollectionNameForStruct returns the Collection object that represents the name and type of the given struct.
//
// The input parameter `doc` must be a pointer to a struct.
//
// If `doc` is not a pointer, it returns an error with the message "needs a pointer to a value".
// If `doc` is a pointer to a pointer, it returns an error with the message "a pointer to a pointer is not allowed".
// If `doc` is not a struct pointer, it returns an error with the message "needs a struct pointer".
//
// It uses the `d.parser.Parse` method to parse the input value and create the Collection object.
// If the parsing is successful, it returns the Collection object and a nil error.
// If an error occurs during parsing, it returns the error.
//
// Example usage:
//
//	coll, err := CollectionNameForStruct(&MyStruct{})
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Collection Name:", coll.Name)
//	fmt.Println("Collection Type:", coll.Type)
//
// Parameters:
// - doc: A pointer to a struct
//
// Returns:
// - The Collection object that represents the name and type of the given struct
// - An error if the input is invalid or if an error occurs during parsing
func (d *defaultClient) CollectionNameForStruct(doc interface{}) (*schemas.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := d.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}

	return t, nil
}

//func (d *defaultClient) SetDatabase(string string) driver.Client {
//	d.db = string
//	return d
//}

// CollectionNameForSlice returns a *schemas.Collection and an error based on the provided document.
// It takes a document as input and checks if it is a pointer to a slice or a map.
// If it is a pointer to a slice, it gets the element type of the slice and checks if it is a struct.
// If it is a struct, it parses the element type using the parser.Parse method.
// If it is not a pointer to a slice or if the element type is not a struct, it parses the document using the parser.Parse method.
// It returns the parsed *schemas.Collection and any error encountered during parsing.
func (d *defaultClient) CollectionNameForSlice(doc interface{}) (*schemas.Collection, error) {
	sliceValue := reflect.Indirect(reflect.ValueOf(doc))

	if sliceValue.Kind() != reflect.Slice && reflect.Map != sliceValue.Kind() {
		return nil, errors.New("needs a pointer to a slice or a map")
	}

	var t *schemas.Collection
	var err error
	if sliceValue.Kind() == reflect.Slice {
		sliceElementType := sliceValue.Type().Elem()
		if sliceElementType.Kind() == reflect.Struct {
			pv := reflect.New(sliceElementType)
			t, err = d.parser.Parse(pv)
		}
	} else {
		t, err = d.parser.Parse(sliceValue)
	}

	if err != nil {
		return nil, err
	}
	return t, nil
}

// getStructCollAndSetKey receives a struct pointer `doc` and returns a pointer to a `schemas.Collection`
// and an error. It first checks if `doc` is a pointer. If not, it returns an error.
// Then it checks if the indirect value of `doc` is also a pointer. If it is, it returns an error.
// Next, it checks if the indirect value of `doc` is a struct. If it is not, it returns an error.
// After that, it calls `d.parser.Parse` passing `beanValue` to obtain the `schemas.Collection` pointer `t`.
// It then retrieves the type of `doc` from `t`, and iterates through its fields.
// If a field has a tag that contains "_id" as a value, it breaks the loop.
// Finally, it returns `t` and `nil`.
func (d *defaultClient) getStructCollAndSetKey(doc interface{}) (*schemas.Collection, error) {
	beanValue := reflect.ValueOf(doc)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}

	if beanValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("needs a struct pointer")
	}
	t, err := d.parser.Parse(beanValue)
	if err != nil {
		return nil, err
	}
	docTyp := t.Type
	for i := 0; i < docTyp.NumField(); i++ {
		field := docTyp.Field(i)
		if strings.Index(field.Tag.Get("bson"), "_id") > 0 {
			//d.e = append(d.e, session("_id", beanValue.Field(i).Interface()))
			break
		}
	}

	return t, nil
}
