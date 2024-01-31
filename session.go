package pie

import (
	"context"
	"errors"
	"fmt"
	"github.com/5xxxx/pie/utils"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	BulkWrite(docs any, ctx ...context.Context) (*mongo.BulkWriteResult, error)

	FilterBy(object any) Session

	Distinct(doc any, columns string, ctx ...context.Context) ([]any, error)

	ReplaceOne(doc any, ctx ...context.Context) (*mongo.UpdateResult, error)

	FindOneAndReplace(doc any, ctx ...context.Context) error

	FindOneAndUpdate(doc any, ctx ...context.Context) (*mongo.SingleResult, error)

	FindOneAndUpdateBson(coll any, bson any, ctx ...context.Context) (*mongo.SingleResult, error)

	FindPagination(needCount bool, rowsSlicePtr any, ctx ...context.Context) (int64, error)

	FindAndDelete(doc any, ctx ...context.Context) error

	// FindOne executes a find command and returns a SingleResult for one document in the collectionByName.
	FindOne(doc any, ctx ...context.Context) error

	// FindAll Find executes a find command and returns a Cursor over the matching documents in the collectionByName.
	FindAll(rowsSlicePtr any, ctx ...context.Context) error

	// InsertOne executes an insert command to insert a single document into the collectionByName.
	InsertOne(doc any, ctx ...context.Context) (primitive.ObjectID, error)

	// InsertMany executes an insert command to insert multiple documents into the collectionByName.
	InsertMany(docs any, ctx ...context.Context) (*mongo.InsertManyResult, error)

	// DeleteOne executes a delete command to delete at most one document from the collectionByName.
	DeleteOne(doc any, ctx ...context.Context) (*mongo.DeleteResult, error)
	SoftDeleteOne(doc any, ctx ...context.Context) error

	// DeleteMany executes a delete command to delete documents from the collectionByName.
	DeleteMany(doc any, ctx ...context.Context) (*mongo.DeleteResult, error)

	SoftDeleteMany(doc any, ctx ...context.Context) error

	Clone() Session
	Limit(i int64) Session

	Skip(i int64) Session
	Project(i any) Session
	Count(i any, ctx ...context.Context) (int64, error)

	UpdateOne(bean any, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateOneBson(coll any, bson any, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateManyBson(coll any, bson any, ctx ...context.Context) (*mongo.UpdateResult, error)

	UpdateMany(bean any, ctx ...context.Context) (*mongo.UpdateResult, error)

	RegexFilter(key, pattern string) Session

	ID(id any) Session

	Asc(colNames ...string) Session

	Desc(colNames ...string) Session

	Sort(colNames ...string) Session
	Soft(f bool) Session
	Filter(key string, value any) Session
	FilterBson(d bson.D) Session
	// Eq Equals a Specified Value
	//{ qty: 20 }
	//Field in Embedded Document Equals a Value
	//{"item.name": "ab" }
	// Equals an Array Value
	//{ tags: [ "A", "B" ] }
	Eq(key string, value any) Session

	// Gt {field: {$gt: value} } >
	Gt(key string, gt any) Session

	// Gte { qty: { $gte: 20 } } >=
	Gte(key string, gte any) Session

	// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
	// tags: { $in: [ /^be/, /^st/ ] } }
	// in []string []int ...
	In(key string, in any) Session

	// Lt {field: {$lt: value} } <
	Lt(key string, lt any) Session

	// Lte { field: { $lte: value} } <=
	Lte(key string, lte any) Session

	// Ne {field: {$ne: value} } !=
	Ne(key string, ne any) Session

	// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
	Nin(key string, nin any) Session

	// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	//$and: [
	//        { $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
	//        { $or: [ { sale: true }, { price : { $lt : 5 } } ] }
	// ]
	And(c Condition) Session

	// Not { field: { $not: { <operator-expression> } } }
	//not and Regular Expressions
	//{ item: { $not: /^p.*/ } }
	Not(key string, not any) Session

	// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
	// { sale: true }, { sale: { $exists: false } } ] }
	// price != 1.99 || sale != true || sale exists || sale exists
	Nor(c Condition) Session

	// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
	Or(c Condition) Session

	Exists(key string, exists bool, filter ...Condition) Session

	// SetArrayFilters sets the value for the ArrayFilters field.
	SetArrayFilters(filters options.ArrayFilters) Session

	// SetOrdered sets the value for the Ordered field.
	SetOrdered(ordered bool) Session

	// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
	SetBypassDocumentValidation(b bool) Session

	// SetReturnDocument sets the value for the ReturnDocument field.
	SetReturnDocument(rd options.ReturnDocument) Session

	// SetUpsert sets the value for the Upsert field.
	SetUpsert(b bool) Session

	// SetCollation sets the value for the Collation field.
	SetCollation(collation *options.Collation) Session

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) Session

	// SetProjection sets the value for the Projection field.
	SetProjection(projection any) Session

	// SetSort sets the value for the Sort field.
	SetSort(sort any) Session

	// SetHint sets the value for the Hint field.
	SetHint(hint any) Session

	// Type { field: { $type: <BSON type> } }
	// { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" },
	// { "_id" : 2, address: "156 Lunar Place", zipCode : 43339374 },
	// db.find( { "zipCode" : { $type : 2 } } ); or db.find( { "zipCode" : { $type : "string" } }
	// return { "_id" : 1, address : "2030 Martian Way", zipCode : "90698345" }
	Type(key string, t any) Session

	// Expr Allows the use of aggregation expressions within the query language.
	//{ $expr: { <expression> } }
	//$expr can build query expressions that compare fields from the same document in a $match stage
	//todo 没用过，不知道行不行。。https://docs.mongodb.com/manual/reference/operator/query/expr/#op._S_expr
	Expr(c Condition) Session

	// Regex todo 简单实现，后续增加支持
	Regex(key string, value string) Session

	SetDatabase(db string) Session

	SetCollRegistry(r *bsoncodec.Registry) Session

	SetCollReadPreference(rp *readpref.ReadPref) Session

	SetCollWriteConcern(wc *writeconcern.WriteConcern) Session

	SetReadConcern(rc *readconcern.ReadConcern) Session
}

type session struct {
	db                    string
	engine                Client
	filter                Condition
	findOneOptions        []*options.FindOneOptions
	findOptions           []*options.FindOptions
	insertManyOpts        []*options.InsertManyOptions
	insertOneOpts         []*options.InsertOneOptions
	deleteOpts            []*options.DeleteOptions
	updateOpts            []*options.UpdateOptions
	countOpts             []*options.CountOptions
	distinctOpts          []*options.DistinctOptions
	findOneAndDeleteOpts  []*options.FindOneAndDeleteOptions
	findOneAndReplaceOpts []*options.FindOneAndReplaceOptions
	findOneAndUpdateOpts  []*options.FindOneAndUpdateOptions
	replaceOpts           []*options.ReplaceOptions
	bulkWriteOptions      []*options.BulkWriteOptions
	collOpts              []*options.CollectionOptions
}

func (s *session) Project(i any) Session {
	s.findOptions = append(s.findOptions, options.Find().SetProjection(i))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetProjection(i))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetProjection(i))
	return s
}

// Soft sets the `deleted_at` field of the session's filter object to the given value.
// If `f` is true, it indicates that the session is soft-deleted.
// If `f` is false, it indicates that the session is not soft-deleted.
// The updated filter conditions are applied to the session.
// The method then returns the session object itself for method chaining.
func (s *session) Soft(f bool) Session {
	s.filter.Exists("deleted_at", f)
	return s
}

// FilterBson applies a BSON filter to the session's filter object.
// The provided BSON filter is added to the existing filter conditions.
// The method then returns the session object itself for method chaining.
func (s *session) FilterBson(d bson.D) Session {
	s.filter.FilterBson(d)
	return s
}

// NewSession creates a new session with the specified engine and default condition filter.
// It returns a Session interface which can be used to interact with the engine.
func NewSession(engine Client) Session {
	return &session{engine: engine, filter: DefaultCondition()}
}

func (s *session) prepareContext(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		return ctx[0]
	}
	return context.Background()
}

// FindPagination retrieves a paginated set of documents from a collection and optionally returns the total number of matching documents.
// The method takes a boolean parameter "needCount" to indicate whether or not to count the total number of matching documents.
// It also accepts a pointer to a slice which will contain the retrieved documents.
// Optionally, a context.Context object can be passed as a variadic parameter to customize the behavior of the method.
//
// The method first obtains the collection associated with the provided rowsSlicePtr using the "collectionForSlice" method.
// If there is an error obtaining the collection, an error is returned.
//
// Next, the method retrieves the filter conditions from the filter object associated with the session.
// If there is an error obtaining the filter conditions, an error is returned.
//
// A context object is prepared by calling the "prepareContext" method with the provided context.Context object(s).
//
// The "coll.Find" method is then called with the obtained filter conditions and any additional find options specified.
// If there is an
func (s *session) FindPagination(needCount bool, rowsSlicePtr any, ctx ...context.Context) (int64, error) {
	coll, err := s.collectionForSlice(rowsSlicePtr)
	if err != nil {
		return 0, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return 0, err
	}
	c := s.prepareContext(ctx...)

	cursor, err := coll.Find(c, filters, s.findOptions...)
	if err != nil {
		return 0, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		if err = cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}(cursor, c)

	var rowCount int64
	if needCount {
		rowCount, err = coll.CountDocuments(c, filters, s.countOpts...)
		if err != nil {
			return 0, err
		}
	}
	if err = cursor.All(c, rowsSlicePtr); err != nil {
		return 0, err
	}
	return rowCount, nil
}

// BulkWrite is a method that performs a bulk write operation on the session's collection.
// It takes in a slice of documents and an optional context.
// The passed documents can be of any type that can be converted into a BSON document.
// The method returns a *mongo.BulkWriteResult and an error.
//
// The method first retrieves the collection from the session using the 'collectionForSlice' method.
// If an error occurs during the retrieval, the error is returned.
//
// Then, it uses reflection to iterate over the values in the 'docs' slice and creates an array of mongo.WriteModel.
// Each value in the 'docs' slice is converted into a BSON document and added as an insert one model to the 'mods' array.
//
// After that, the method prepares the context using the 'prepareContext
func (s *session) BulkWrite(docs any, ctx ...context.Context) (*mongo.BulkWriteResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}
	values := reflect.ValueOf(docs)
	var mods []mongo.WriteModel
	for i := 0; i < values.Len(); i++ {
		mods = append(mods, mongo.NewInsertOneModel().SetDocument(values.Index(i).Interface()))
	}
	c := s.prepareContext(ctx...)

	return coll.BulkWrite(c, mods, s.bulkWriteOptions...)
}

// FilterBy sets the filter for the session to be used in the subsequent database operations.
// The filter is specified by the given object.
// The function returns the session itself to allow method chaining.
func (s *session) FilterBy(object any) Session {
	s.filter.FilterBy(object)
	return s
}

// Distinct retrieves distinct values for the specified columns from the collection associated with the session.
// The provided 'doc' object is used to determine the collection to perform the distinct operation on.
// It returns a slice of distinct any values for the specified columns.
// The 'columns' parameter specifies the columns from which to retrieve distinct values.
// The optional 'ctx' parameter allows you to provide a context.Context object to control the execution of the operation.
// It returns an error if any error occurs during the operation.
// The method internally retrieves the collection based on the provided 'doc' object.
// It then retrieves the filters associated with the session's filter object.
// Finally, it performs the distinct operation on the collection using the retrieved context, columns, filters, and distinct options.
// It returns the result of the distinct operation and any error that occurred.
func (s *session) Distinct(doc any, columns string, ctx ...context.Context) ([]any, error) {
	coll, err := s.collectionForSlice(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := s.prepareContext(ctx...)

	return coll.Distinct(c, columns, filters, s.distinctOpts...)
}

// ReplaceOne replaces a single document in the collection that matches the specified filters with the provided document.
// It returns the *mongo.UpdateResult, indicating the number of documents matched and modified, and an error if any.
// The method first retrieves the collection associated with the document's struct using the s.collectionForStruct method.
// If an error occurs during this retrieval, it is returned along with nil as the *mongo.UpdateResult.
// Then, it retrieves the filters from the session's filter object using the s.filter.Filters method.
// If an error occurs during this retrieval, it is returned along with nil as the *mongo.UpdateResult.
// It prepares the context by creating a new context if ctx is not provided or using the provided context otherwise.
// Finally, it calls the coll.ReplaceOne method to perform the replacement and returns the result or any error encountered.
func (s *session) ReplaceOne(doc any, ctx ...context.Context) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}

	c := s.prepareContext(ctx...)

	return coll.ReplaceOne(c, filters, doc, s.replaceOpts...)
}

// FindOneAndReplace executes a find and replace command for one document in the collection.
// It takes the specified document and replaces the first document that matches the given filters.
// If the document is not found, an error is returned.
// The function returns an error if there was an issue finding the collection or if there was an error
// decoding the replaced document into the original document variable.
// If a context is provided, it is used for the operation. Otherwise, a default background context is used.
func (s *session) FindOneAndReplace(doc any, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}

	c := s.prepareContext(ctx...)

	return coll.FindOneAndReplace(c, filters, doc, s.findOneAndReplaceOpts...).Decode(&doc)
}

// FindOneAndUpdateBson executes a find and update command on the collection.
// It takes in the following parameters:
// - coll: the collection on which to execute the command
// - bson: the update document in BSON format
// - ctx: optional context.Context object for cancellation and deadline propagation
//
// It returns a *mongo.SingleResult and an error. The SingleResult contains
// the result of the operation, and the error indicates any encountered errors
// during the execution of the command.
func (s *session) FindOneAndUpdateBson(coll any, bson any, ctx ...context.Context) (*mongo.SingleResult, error) {
	c, err := s.collectionForStruct(coll)
	if err != nil {
		return nil, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}

	cc := s.prepareContext(ctx...)
	return c.FindOneAndUpdate(cc, filters, bson, s.findOneAndUpdateOpts...), nil
}

// FindOneAndUpdate updates a single document in the given collection based on the specified filter conditions.
// The method takes a document interface and an optional context as parameters.
// It first determines the collection for the given document using the 'collectionForStruct' method.
// If an error occurs during the determination, it returns nil and the error.
// Next, it retrieves the filter conditions using the 'Filters' method of the session's filter object.
// If an error occurs during the retrieval, it returns nil and the error.
// The method then prepares the context using the 'prepareContext' method.
// It calls the 'FindOneAndUpdate' method on the collection with the prepared context, filters, and the update document, which is formatted as a '$set' operator.
// The method also includes the session's 'findOneAndUpdateOpts' as additional options.
// Finally, it returns the single result and any error that occurred during the update process.
func (s *session) FindOneAndUpdate(doc any, ctx ...context.Context) (*mongo.SingleResult, error) {

	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := s.prepareContext(ctx...)
	return coll.FindOneAndUpdate(c, filters, bson.M{"$set": doc}, s.findOneAndUpdateOpts...), nil
}

// FindAndDelete deletes a single document from the collection based on the provided filters.
// The document to be deleted is specified with the "doc" parameter.
// The method returns an error if the collection for the provided document cannot be found.
// If filters cannot be fetched from the session's filter object, an error is returned.
// The method also accepts an optional context.Context parameter.
// The context is used to control the behavior of the operation, such as timeouts or cancellation.
// If no context is provided, a default context is used.
// The function queries the collection for a document that matches the provided filters
// and deletes that document using the FindOneAndDelete method.
// The deleted document is then decoded into the "doc" parameter.
// If decoding fails, an error is returned.
// The method
func (s *session) FindAndDelete(doc any, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := s.prepareContext(ctx...)
	return coll.FindOneAndDelete(c, filters, s.findOneAndDeleteOpts...).Decode(&doc)
}

// FindOne finds a single document in the collection that matches the specified filters.
// The method takes a pointer to a document struct and an optional context as parameters.
// It returns an error if any error occurs during the process.
// The document struct must be provided as a pointer, and it will be populated with the found document's data.
// The method first determines the appropriate collection for the provided document using the collectionForStruct method.
// If an error occurs during this process, it is
func (s *session) FindOne(doc any, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := s.prepareContext(ctx...)
	result := coll.FindOne(c, filters, s.findOneOptions...)
	if err = result.Err(); err != nil {
		return err
	}

	if err = result.Decode(doc); err != nil {
		return err
	}

	return nil
}

// FindAll retrieves all documents from the collection specified by the session's filter
// and populates the provided slice pointer with the results.
// The slice pointer should be of type []<document-type>.
// Optionally, a context can be passed to customize the operation.
// If no context
func (s *session) FindAll(rowsSlicePtr any, ctx ...context.Context) error {
	coll, err := s.collectionForSlice(rowsSlicePtr)
	if err != nil {
		return err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := s.prepareContext(ctx...)

	cursor, err := coll.Find(c, filters, s.findOptions...)
	if err != nil {
		return err
	}

	if err = cursor.All(c, rowsSlicePtr); err != nil {
		return err
	}

	return nil
}

// InsertOne inserts a single document into the collection.
// It returns the inserted document's ObjectID and any error that occurred during the insertion.
// If an error occurs during the insertion, the returned ObjectID will be [12]byte{} and the error will be non-nil.
// The inserted document's ObjectID can be retrieved by performing a type assertion on the InsertedID field of the returned result.
// Example:
//
//	insertedID, err := session.InsertOne(document)
//	if err != nil {
//	  // handle error
//	} else {
//	  // handle success
//	}
func (s *session) InsertOne(doc any, ctx ...context.Context) (primitive.ObjectID, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return [12]byte{}, err
	}
	c := s.prepareContext(ctx...)
	result, err := coll.InsertOne(c, doc, s.insertOneOpts...)
	if err != nil {
		return [12]byte{}, err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, err
	}
	return [12]byte{}, err
}

// InsertMany inserts multiple documents into the collection.
// It takes a pointer to a slice of documents as the first argument,
// and an optional context as the second argument.
// The documents are converted to a slice of interfaces using reflection.
// The function then retrieves the collection associated with the documents,
// and inserts the documents into the collection using the InsertMany method.
// The function returns the InsertManyResult and an error, if any.
func (s *session) InsertMany(docs any, ctx ...context.Context) (*mongo.InsertManyResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(docs)
	var many []any
	for index := 0; index < value.Len(); index++ {
		many = append(many, value.Index(index).Interface())
	}
	c := s.prepareContext(ctx...)
	return coll.InsertMany(c, many, s.insertManyOpts...)
}

// DeleteOne executes a delete command and returns a DeleteResult for one document in the collection.
// It takes in the document to be deleted and allows for an optional context.
// If an error occurs during the operation, it returns nil for DeleteResult and the error.
// Example usage:
//
//	result, err := session.DeleteOne(&user, context.Background())
//	if err != nil {
//	  // handle error
//	}
//	fmt.Printf("Deleted %d documents\n", result.DeletedCount)
//	// Output: Deleted 1 documents
func (s *session) DeleteOne(doc any, ctx ...context.Context) (*mongo.DeleteResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := s.prepareContext(ctx...)
	return coll.DeleteOne(c, filters, s.deleteOpts...)
}

// SoftDeleteOne soft deletes one document in the collection.
// It updates the document by setting the "deleted_at" field to the current time.
// The document is identified by the provided filters.
//
// The method first retrieves the collection for the given document.
// It then prepares the filters to be used in the update operation.
// If a context is provided, it uses that context; otherwise, it uses the background context.
// Finally, it performs the update operation by calling UpdateOne on the collection.
// The "deleted_at" field is updated with the current time using the $set operator.
// If any error occurs during the update operation, it is returned.
func (s *session) SoftDeleteOne(doc any, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := s.prepareContext(ctx...)
	_, err = coll.UpdateOne(c, filters, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})

	return err
}

// DeleteMany deletes multiple documents from the collection associated with the session.
// The method takes the document interface as the first argument, which represents the filter conditions for deleting documents.
// It also accepts a variadic argument ctx of type context.Context, which allows passing additional context options.
// The method returns a *mongo.DeleteResult, which contains information about the deletion operation, and an
func (s *session) DeleteMany(doc any, ctx ...context.Context) (*mongo.DeleteResult, error) {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return nil, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := s.prepareContext(ctx...)
	return coll.DeleteMany(c, filters, s.deleteOpts...)
}

// SoftDeleteMany executes an update command to "soft delete" multiple documents in the collection.
// It adds a "deleted_at" field with the current timestamp to the matching documents.
// The method takes an any type parameter (doc) representing the document(s) to be soft deleted.
// It returns an error if any error occurs during the update operation.
func (s *session) SoftDeleteMany(doc any, ctx ...context.Context) error {
	coll, err := s.collectionForStruct(doc)
	if err != nil {
		return err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return err
	}
	c := s.prepareContext(ctx...)
	_, err = coll.UpdateMany(c, filters, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})

	return err
}

// Clone creates a new instance of the session and returns it as a Session.
// The new session has the same values for the db, engine, filter, and various options
// as the original session.
// The cloned session is independent from the original session and can be used separately.
func (s *session) Clone() Session {
	sess := session{
		db:                    s.db,
		engine:                s.engine,
		filter:                s.filter.Clone(),
		findOneOptions:        s.findOneOptions,
		findOptions:           s.findOptions,
		insertManyOpts:        s.insertManyOpts,
		insertOneOpts:         s.insertOneOpts,
		deleteOpts:            s.deleteOpts,
		findOneAndDeleteOpts:  s.findOneAndDeleteOpts,
		updateOpts:            s.updateOpts,
		countOpts:             s.countOpts,
		distinctOpts:          s.distinctOpts,
		findOneAndReplaceOpts: s.findOneAndReplaceOpts,
		findOneAndUpdateOpts:  s.findOneAndUpdateOpts,
		replaceOpts:           s.replaceOpts,
		bulkWriteOptions:      s.bulkWriteOptions,
	}

	return &sess
}

// Limit sets the maximum number of documents to be returned by a find operation in the session.
// The limit value must be a positive integer.
func (s *session) Limit(i int64) Session {
	s.findOptions = append(s.findOptions, options.Find().SetLimit(i))
	return s
}

// SetReadConcern sets the value for the ReadConcern field.
func (s *session) SetReadConcern(rc *readconcern.ReadConcern) Session {
	s.collOpts = append(s.collOpts, options.Collection().SetReadConcern(rc))

	return s
}

// SetCollWriteConcern sets the write concern for the collection in the current session.
// It appends the options.Collection().SetWriteConcern(wc) to the s.collOpts field.
// The updated session is returned.
func (s *session) SetCollWriteConcern(wc *writeconcern.WriteConcern) Session {
	s.collOpts = append(s.collOpts, options.Collection().SetWriteConcern(wc))
	return s
}

// SetCollReadPreference sets the read preference for the session's collection options.
// The read preference determines how the driver routes read operations to replica set members or shards.
// The parameter rp is a pointer to a readpref.ReadPref instance representing the desired read preference.
// It modifies the session's collection options and returns the modified session.
// Example usage:
//
//	session := &session{}
//	readPref := readpref.Primary()
//	session.SetCollReadPreference(readPref)
//	// Continue using the session with the updated collection options.
func (s *session) SetCollReadPreference(rp *readpref.ReadPref) Session {
	s.collOpts = append(s.collOpts, options.Collection().SetReadPreference(rp))
	return s
}

// SetCollRegistry sets the bsoncodec.Registry for the session's collection.
// It appends the options.Collection().SetRegistry() to the session's collOpts
// and returns the updated session.
func (s *session) SetCollRegistry(r *bsoncodec.Registry) Session {
	s.collOpts = append(s.collOpts, options.Collection().SetRegistry(r))
	return s
}

// Skip sets the number of documents to skip before returning results.
// It adds the skip option to the find and findOne options in the session.
// The skip value is specified by the i parameter.
// It returns the session.
func (s *session) Skip(i int64) Session {
	s.findOptions = append(s.findOptions, options.Find().SetSkip(i))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSkip(i))
	return s
}

func (s *session) Count(i any, ctx ...context.Context) (int64, error) {
	kind := reflect.TypeOf(i).Kind()
	if kind == reflect.Ptr {
		kind = reflect.TypeOf(reflect.Indirect(reflect.ValueOf(i)).Interface()).Kind()
	}
	var coll *mongo.Collection
	var err error
	switch kind {
	case reflect.Slice:
		coll, err = s.collectionForSlice(i)
	case reflect.Struct:
		coll, err = s.collectionForStruct(i)
	default:
		return 0, errors.New("need slice or struct")
	}

	if err != nil {
		return 0, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return 0, err
	}

	c := s.prepareContext(ctx...)

	return coll.CountDocuments(c, filters, s.countOpts...)
}

// UpdateOne executes an update command and returns a *mongo.UpdateResult for the matched document in the collection.
// Parameters:
// - `bean` : The document to update. Must be a non-zero value.
// - `ctx` : Optional context.Context. If provided, it will be used as the context for the operation.
// Returns:
// - *mongo.UpdateResult : The result of the update operation, including the number of documents matched and modified.
// - error : Any error that occurred during the update operation.
// Example usage:
//
//	result, err := session.UpdateOne(bean)
//	if err != nil {
//	  // Handle error
//	}
//	fmt.Printf("Updated documents: %v\n", result.ModifiedCount)
func (s *session) UpdateOne(bean any, ctx ...context.Context) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForStruct(bean)

	if err != nil {
		return nil, err
	}

	if utils.IsStructZero(reflect.ValueOf(bean).Elem()) {
		return nil, nil
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := s.prepareContext(ctx...)
	return coll.UpdateOne(c, filters, bson.M{"$set": bean}, s.updateOpts...)
}

// UpdateOneBson updates a single document in the collection corresponding to the given struct
// using the provided BSON filter.
// It returns a pointer to mongo.UpdateResult, which contains information about the update operation,
// and an error if any.
// Parameters:
// - coll: the collection or struct for which the update operation needs to be performed.
// - bson: the BSON filter to identify the document to be updated.
// - ctx: optional context.Context for the update operation.
// Usage:
// result, err := session.UpdateOneBson(collection, bsonFilter)
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//
// fmt.Println("Matched Count:", result.MatchedCount)
// fmt.Println("Modified Count:", result.ModifiedCount)
func (s *session) UpdateOneBson(coll any, bson any, ctx ...context.Context) (*mongo.UpdateResult, error) {
	c, err := s.collectionForStruct(coll)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	cc := s.prepareContext(ctx...)
	return c.UpdateOne(cc, filters, bson, s.updateOpts...)
}

// UpdateManyBson updates multiple documents in the collection
// based on the specified BSON document.
// It returns the UpdateResult that contains information about
// the update operation and any errors encountered.
//
// Parameters:
// - coll: The collection or struct to update.
// - bson: The BSON document specifying the update.
// - ctx: Optional context for the operation.
//
// Returns:
// - *mongo.UpdateResult: The result of the update operation.
// - error: Any error encountered during the operation.
func (s *session) UpdateManyBson(coll any, bson any, ctx ...context.Context) (*mongo.UpdateResult, error) {
	c, err := s.collectionForStruct(coll)
	if err != nil {
		return nil, err
	}
	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	cc := s.prepareContext(ctx...)
	return c.UpdateMany(cc, filters, bson, s.updateOpts...)
}

func (s *session) toBson(obj any) bson.M {
	beanValue := reflect.ValueOf(obj).Elem()
	if beanValue.Kind() != reflect.Struct ||
		//Todo how to fix array?
		beanValue.Kind() == reflect.Array {
		panic(errors.New("needs a struct"))
	}

	ret := bson.M{}
	docType := reflect.TypeOf(obj).Elem()
	for index := 0; index < docType.NumField(); index++ {
		fieldTag := docType.Field(index).Tag.Get("bson")
		if fieldTag != "" && fieldTag != "-" {
			s.makeValue(fieldTag, beanValue.Field(index).Interface(), ret)
		}
	}
	return ret
}

func (s *session) makeValue(field string, value any, ret bson.M) {
	split := strings.Split(field, ",")
	if len(split) <= 0 {
		return
	}
	if strings.Contains(field, "omitempty") {
		if utils.IsZero(value) {
			return
		}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Struct:
		s.makeStruct(field, v, ret)
		return
	case reflect.Array:
		return
	}
	ret[split[0]] = value
}

// makeStruct iterates through the fields of a struct and populates a bson.M map with the field values.
func (s *session) makeStruct(field string, value reflect.Value, ret bson.M) {
	for index := 0; index < value.NumField(); index++ {
		docType := reflect.TypeOf(value.Interface())
		tag := docType.Field(index).Tag.Get("bson")
		split := strings.Split(tag, ",")
		if len(split) > 0 {
			if tag != "" {
				if strings.Contains(tag, "omitempty") {
					if !utils.IsZero(value.Field(index)) {
						fieldTags := fmt.Sprintf("%s.%s", field, split[0])
						s.makeValue(fieldTags, value.Field(index).Interface(), ret)
					}
				} else {
					fieldTags := fmt.Sprintf("%s.%s", field, split[0])
					s.makeValue(fieldTags, value.Field(index).Interface(), ret)
				}

			}
		}
	}
}

func (s *session) UpdateMany(bean any, ctx ...context.Context) (*mongo.UpdateResult, error) {
	coll, err := s.collectionForSlice(bean)
	if err != nil {
		return nil, err
	}

	filters, err := s.filter.Filters()
	if err != nil {
		return nil, err
	}
	c := s.prepareContext(ctx...)
	return coll.UpdateMany(c, filters, bson.M{"$set": bean}, s.updateOpts...)

}

func (s *session) RegexFilter(key, pattern string) Session {
	s.filter.RegexFilter(key, pattern)
	return s
}

// ID sets the filter condition on the session's filter object
// to search for records with the specified ID value.
// The provided 'id' parameter is used as the value to filter by.
// The method then returns the session object itself for method chaining.
func (s *session) ID(id any) Session {
	s.filter.ID(id)
	return s
}

// Asc sets the sorting order on the session's find options.
// The provided column names are used to determine the order of sorting.
// Multiple column names can be provided to define a multi-column sort.
// The method uses the ascending (1) sorting order for the specified columns.
// If no column names are provided, the method returns the session object itself.
// The method modifies the session's find and findOne options to include the sort criteria.
// The modified options are used in subsequent find and findOne operations on the session.
// The method returns the session object itself for method chaining.
func (s *session) Asc(colNames ...string) Session {
	if len(colNames) == 0 {
		return s
	}

	es := bson.M{}
	for _, c := range colNames {
		es[c] = 1
	}
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	return s
}

// Desc sets the sort order of the session's find and findOne options based on the provided column names.
// The column names are passed as variadic arguments to the method.
// If no column names are provided, the method immediately returns the session object itself.
// Otherwise, a descending sort order is applied to each column name in the find and findOne options.
// The resulting sort order is added to the session's findOptions and findOneOptions respectively.
// Finally, the method returns the session object for method chaining.
func (s *session) Desc(colNames ...string) Session {
	if len(colNames) == 0 {
		return s
	}

	es := bson.M{}
	for _, c := range colNames {
		es[c] = -1
	}

	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	return s
}

// Sort sets the sorting options on the session's findOptions and findOneOptions objects.
// The provided column names are used to determine the sorting order.
//
// If no column names are provided, the method returns the session object itself for method chaining.
//
// The column names are passed as variadic arguments, allowing multiple columns to be sorted.
// Each column name should be a string.
//
// The sorting order for each column is based on the first character of the column name string:
// - If the first character is '-', the column is sorted in descending order (e.g., "-name").
// - If the first character is any other character, the column is sorted in ascending order (e.g., "name").
//
// The method appends the sorting options to the session's findOptions and findOneOptions objects.
// The sorting options are set using the bson.E type from the "go.mongodb.org/mongo-driver/bson" package.
// The key of the bson.E represents the column name, and the value represents the sorting order (1 for ascending, -1 for descending).
//
// Finally, the method returns the session object itself for method chaining.
func (s *session) Sort(colNames ...string) Session {
	if len(colNames) == 0 {
		return s
	}
	es := bson.D{}
	for _, field := range colNames {
		if field != "" {
			switch field[0] {
			case '-':
				es = append(es, bson.E{Key: field[1:], Value: -1})
			default:
				es = append(es, bson.E{Key: field, Value: 1})
			}
		}
	}
	s.findOptions = append(s.findOptions, options.Find().SetSort(es))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSort(es))
	return s
}

func (s *session) Filter(key string, value any) Session {
	return s.Eq(key, value)
}

// Eq Equals a Specified Value
// { qty: 20 }
// Field in Embedded Document Equals a Value
// {"item.name": "ab" }
// Equals an Array Value
// { tags: [ "A", "B" ] }
func (s *session) Eq(key string, value any) Session {
	s.filter.Eq(key, value)
	return s
}

// Gt sets the "greater than" filter condition on the session's filter object.
// The provided key parameter indicates the field name on which the condition is applied.
// The gt parameter specifies the value that the field should be greater than.
// The method then returns the session object itself for method chaining.
func (s *session) Gt(key string, gt any) Session {
	s.filter.Gt(key, gt)
	return s
}

// Gte { qty: { $gte: 20 } } >=
func (s *session) Gte(key string, gte any) Session {
	s.filter.Gte(key, gte)
	return s
}

// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (s *session) In(key string, in any) Session {
	s.filter.In(key, in)
	return s
}

// Lt {field: {$lt: value} } <
func (s *session) Lt(key string, lt any) Session {
	s.filter.Lt(key, lt)
	return s
}

// Lte { field: { $lte: value} } <=
func (s *session) Lte(key string, lte any) Session {
	s.filter.Lte(key, lte)
	return s
}

// Ne {field: {$ne: value} } !=
func (s *session) Ne(key string, ne any) Session {
	s.filter.Ne(key, ne)
	return s
}

// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (s *session) Nin(key string, nin any) Session {
	s.filter.Nin(key, nin)
	return s
}

// And { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
// $and: [
//
//	{ $or: [ { qty: { $lt : 10 } }, { qty : { $gt: 50 } } ] },
//	{ $or: [ { sale: true }, { price : { $lt : 5 } } ] }
//
// ]
func (s *session) And(c Condition) Session {
	s.filter.And(c)
	return s

}

// Not { field: { $not: { <operator-expression> } } }
// not and Regular Expressions
// { item: { $not: /^p.*/ } }
func (s *session) Not(key string, not any) Session {
	s.filter.Not(key, not)
	return s
}

// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (s *session) Nor(c Condition) Session {
	s.filter.Nor(c)
	return s
}

// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (s *session) Or(c Condition) Session {
	s.filter.Or(c)
	return s
}

func (s *session) Exists(key string, exists bool, filter ...Condition) Session {
	s.filter.Exists(key, exists, filter...)
	return s
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (s *session) SetArrayFilters(filters options.ArrayFilters) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetArrayFilters(filters))
	s.updateOpts = append(s.updateOpts, options.Update().SetArrayFilters(filters))
	return s
}

// SetOrdered sets the value for the Ordered field.
func (s *session) SetOrdered(ordered bool) Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetOrdered(ordered))
	return s
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (s *session) SetBypassDocumentValidation(b bool) Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetBypassDocumentValidation(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetBypassDocumentValidation(b))
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts, options.FindOneAndUpdate().SetBypassDocumentValidation(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetBypassDocumentValidation(b))

	return s
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (s *session) SetReturnDocument(rd options.ReturnDocument) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetReturnDocument(rd))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetReturnDocument(rd))
	return s
}

// SetUpsert sets the value for the Upsert field.
func (s *session) SetUpsert(b bool) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetUpsert(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetUpsert(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetUpsert(b))
	return s
}

// SetCollation sets the value for the Collation field.
func (s *session) SetCollation(collation *options.Collation) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetCollation(collation))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetCollation(collation))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetCollation(collation))
	s.updateOpts = append(s.updateOpts, options.Update().SetCollation(collation))
	return s
}

// SetMaxTime sets the value for the MaxTime field.
func (s *session) SetMaxTime(d time.Duration) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetMaxTime(d))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetMaxTime(d))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetMaxTime(d))
	return s
}

// SetProjection sets the projection field on the findOneAndUpdateOpts, findOneAndReplaceOpts,
// and findOneAndDeleteOpts options objects of the session.
// The provided projection value is used to specify the fields that should be included or excluded
// in the query result.
// The method then returns the session object itself for method chaining.
func (s *session) SetProjection(projection any) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetProjection(projection))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetProjection(projection))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetProjection(projection))
	return s
}

// SetSort sets the sort option for the FindOneAndUpdate, FindOneAndReplace, and FindOneAndDelete methods on the session object.
// The provided sort parameter specifies the sorting order for the query.
// The method adds the sort option to the respective findOneAndUpdateOpts, findOneAndReplaceOpts, and findOneAndDeleteOpts options arrays.
// It then returns the session object itself for method chaining.
func (s *session) SetSort(sort any) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetSort(sort))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetSort(sort))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetSort(sort))
	return s
}

// SetHint sets the hint for the session's operations.
func (s *session) SetHint(hint any) Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetHint(hint))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetHint(hint))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetHint(hint))
	s.updateOpts = append(s.updateOpts, options.Update().SetHint(hint))
	return s
}

// Type sets the type condition on the session's filter object for the specified key.
// The type condition checks if the value of the specified key has the same type as the provided any value.
// The method accepts a string key and an any value as parameters.
// It then adds the type condition to the filter object using the specified key and value.
// The method returns the session object itself for method chaining.
func (s *session) Type(key string, t any) Session {
	s.filter.Type(key, t)
	return s
}

// Expr sets a custom filter expression on the session's filter object.
// The provided Condition value represents the custom filter expression.
// The method then returns the session object itself for method chaining.
func (s *session) Expr(c Condition) Session {
	s.filter.Expr(c)
	return s
}

// Regex todo 简单实现，后续增加支持
func (s *session) Regex(key string, value string) Session {
	s.filter.Regex(key, value)
	return s
}

// SetDatabase sets the database name for the session.
// It takes a string argument representing the name of the database.
// It updates the session's `db` field with the provided name and returns the updated session object.
func (s *session) SetDatabase(db string) Session {
	s.db = db
	return s
}

func (s *session) collectionForStruct(doc any) (*mongo.Collection, error) {
	coll, err := s.engine.CollectionNameForStruct(doc)
	if err != nil {
		return nil, err
	}

	return s.collectionByName(coll.Name), nil
}

func (s *session) collectionForSlice(doc any) (*mongo.Collection, error) {
	coll, err := s.engine.CollectionNameForSlice(doc)
	if err != nil {
		return nil, err
	}
	return s.collectionByName(coll.Name), nil
}

func (s *session) collectionByName(name string) *mongo.Collection {
	if s.collOpts == nil {
		s.collOpts = make([]*options.CollectionOptions, 0)
	}

	return s.engine.Collection(name, s.collOpts, s.db)
}

//func (s *session) makeFilterValue(field string, value any) {
//	if utils.IsZero(value) {
//		return
//	}
//	v := reflect.ValueOf(value)
//	switch v.Kind() {
//	case reflect.Struct:
//		s.makeStructValue(field, v)
//	case reflect.Array:
//		return
//	}
//	s.Filter(field, value)
//}
//
//func (s *session) makeStructValue(field string, value reflect.Value) {
//	for index := 0; index < value.NumField(); index++ {
//		docType := reflect.TypeOf(value.Interface())
//		tag := docType.Field(index).Tag.Get("bson")
//		if tag != "" {
//			if !utils.IsZero(value.Field(index)) {
//				fieldTags := fmt.Sprintf("%s.%s", field, tag)
//				s.makeFilterValue(fieldTags, value.Field(index).Interface())
//			}
//		}
//	}
//}
