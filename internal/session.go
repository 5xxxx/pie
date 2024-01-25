package internal

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"reflect"
	"strings"
	"time"

	"github.com/5xxxx/pie/driver"
	"github.com/5xxxx/pie/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type session struct {
	db                    string
	engine                driver.Client
	filter                driver.Condition
	findOneOptions        []*options.FindOneOptions
	findOptions           []*options.FindOptions
	insertManyOpts        []*options.InsertManyOptions
	insertOneOpts         []*options.InsertOneOptions
	deleteOpts            []*options.DeleteOptions
	findOneAndDeleteOpts  []*options.FindOneAndDeleteOptions
	updateOpts            []*options.UpdateOptions
	countOpts             []*options.CountOptions
	distinctOpts          []*options.DistinctOptions
	findOneAndReplaceOpts []*options.FindOneAndReplaceOptions
	findOneAndUpdateOpts  []*options.FindOneAndUpdateOptions
	replaceOpts           []*options.ReplaceOptions
	bulkWriteOptions      []*options.BulkWriteOptions
	collOpts              []*options.CollectionOptions
}

// Soft sets the "deleted_at" filter condition on the session's filter object.
// The provided boolean value is used to determine if the filter should include or exclude deleted records.
// When 'true' is passed, the "deleted_at" field exists condition is added to the filter, including deleted records.
// When 'false' is passed, the "deleted_at" field does not exist condition is added to the filter, excluding deleted records.
// The method then returns the session object itself for method chaining.
func (s *session) Soft(f bool) driver.Session {
	s.filter.Exists("deleted_at", f)
	return s
}

// FilterBson applies a BSON filter to the session's filter object.
// The provided BSON filter is added to the existing filter conditions.
// The method then returns the session object itself for method chaining.
func (s *session) FilterBson(d bson.D) driver.Session {
	s.filter.FilterBson(d)
	return s
}

// NewSession creates a new session with the specified engine and default condition filter.
// It returns a driver.Session interface which can be used to interact with the engine.
func NewSession(engine driver.Client) driver.Session {
	return &session{engine: engine, filter: DefaultCondition()}
}

func (s *session) prepareContext(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		return ctx[0]
	}
	return context.Background()
}

// FindPagination returns a paginated result of documents from the MongoDB collection associated with the session.
// It takes a boolean parameter 'needCount' to determine if the count of documents should be returned as well.
// The 'rowsSlicePtr' parameter is a pointer to the slice where the documents will be stored.
// Optionally, it accepts a 'ctx' parameter of type 'context.Context'.
// It returns the count of documents returned and an error, if any.
// First, it gets the collection object for the provided 'rowsSlicePtr' using the 'collectionForSlice' method.
// If there is an error in getting the collection, it returns 0 and the error.
// Then, it retrieves the filter conditions from the session's filter object using the 'Filters' method.
// If there is an error in getting the filters, it returns 0 and the error.
// It prepares the 'context.Context' using the 'prepareContext' method with the optional 'ctx' parameter.
// The method queries the collection using the obtained filters and the session's 'findOptions' with the 'Find' method.
// If there is an error in querying the collection, it returns 0 and the error.
// The cursor is stored in the 'cursor' variable.
// A defer statement is used to ensure that the cursor is closed when the method returns.
// If there is an error in closing the cursor, it prints the error.
// If the 'needCount' parameter is true, it calls the 'CountDocuments' method on the collection to get the total count of matching documents.
// If there is an error in getting the count, it returns 0 and the error.
// The count is stored in the 'rowCount' variable.
// Then, it retrieves all the documents from the cursor using the 'All' method and stores them in the 'rowsSlicePtr' slice.
// If there is an error in retrieving the documents, it returns 0 and the error.
// Finally, it returns the count of documents and a nil error.
func (s *session) FindPagination(needCount bool, rowsSlicePtr interface{}, ctx ...context.Context) (int64, error) {
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

// BulkWrite executes a bulk write operation on the collection for the given slice of documents.
// It takes an optional context.Context as an argument.
// The method returns a *mongo.BulkWriteResult and an error.
// The method first checks if the collection exists for the given slice of documents.
// If the collection does not exist, it returns an error.
// It then iterates over each document in the slice and creates an insert model using mongo.NewInsertOneModel.
// The insert model is set with the document and appended to the mods slice.
// A context is created, using the default context if no optional context is provided.
// The method calls coll.BulkWrite with the created context, the mods slice, and s.bulkWriteOptions.
// Finally, it returns the *mongo.BulkWriteResult and any error that occurred during the operation.
func (s *session) BulkWrite(docs interface{}, ctx ...context.Context) (*mongo.BulkWriteResult, error) {
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
func (s *session) FilterBy(object interface{}) driver.Session {
	s.filter.FilterBy(object)
	return s
}

// Distinct performs a distinct operation on the collection associated with the session object.
// It takes the following parameters:
//   - doc: The document or struct to be used for determining the collection to perform the distinct operation on.
//   - columns: The name of the field(s) for which the distinct values should be returned.
//   - ctx: Optional context.Context object(s) to use for the operation.
//
// The method first retrieves the collection associated with the provided document or struct using the collectionForSlice method.
// If an error occurs during the collection retrieval, it is returned immediately.
// It then retrieves the filters from the session's filter object using the Filters method.
// If an error occurs during the filters retrieval, it is returned immediately.
// The method creates a new context object using context.Background() and assigns it to the variable c.
// If any optional context.Context objects were provided, the first one is used instead of context.Background().
// Finally, the method invokes the Distinct method on the retrieved collection with the provided context c, columns, filters, and session's distinctOpts.
// The result of the operation is returned as a slice of interface{} and any error encountered is also returned.
func (s *session) Distinct(doc interface{}, columns string, ctx ...context.Context) ([]interface{}, error) {
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

// ReplaceOne executes a replace command and returns the UpdateResult for the document in the collection.
// It replaces a single document in the collection that matches the specified filter with the given replacement document.
// The replacement document can be of any type that can be encoded into BSON.
// If the replacement document does not contain an "_id" field, one will be generated and added to the replacement document.
// The method accepts an optional context.Context as the first parameter to allow for cancellation or deadline.
// If no context is provided, a background context will be used.
// The method returns the UpdateResult, which contains information about the operation such as the number of matched documents and modified documents.
// If any error occurs during the operation, it will be returned along with the nil UpdateResult.
func (s *session) ReplaceOne(doc interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
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
func (s *session) FindOneAndReplace(doc interface{}, ctx ...context.Context) error {
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
func (s *session) FindOneAndUpdateBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {
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

// FindOneAndUpdate executes a find and update command on the collection and returns the SingleResult for the updated document.
// It takes the document to be updated and an optional context.
// If the operation is successful, it returns the SingleResult containing the updated document and nil error.
// If there is an error during the operation, it returns nil for the SingleResult and the error.
func (s *session) FindOneAndUpdate(doc interface{}, ctx ...context.Context) (*mongo.SingleResult, error) {

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

// FindAndDelete executes a findAndDelete command and returns a SingleResult for the deleted document in the collection.
// It takes in a document and an optional context as parameters. The document is used to determine the collection.
// If an error occurs during the collection retrieval, it is returned immediately.
// The filters for the find command are retrieved using the filter object of the session.
// If an error occurs during the retrieval of filters, it is returned immediately.
// The context is either obtained from the parameter or from a newly created background context if no parameter is provided.
// The findAndDelete command is executed using the obtained collection, filters, and findOneAndDelete options of the session.
// If an error occurs during the execution of the command, it is returned immediately.
// The retrieved document is decoded and stored in the provided document interface{}.
// If an error occurs during the decoding process, it is returned immediately.
// Finally, if all operations succeed without any errors, nil is returned to indicate success.
func (s *session) FindAndDelete(doc interface{}, ctx ...context.Context) error {
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

// FindOne FindOne executes a find command and returns a single document from the collectionByName.
//
// Parameters:
// - "doc" : a pointer to a struct where the result will be decoded into.
// - "ctx" (optional) : context.Context containing any additional options or settings for the operation.
//
// Returns:
// - "error" : the error encountered during the operation, if any. If the operation is successful, it returns nil.
//
// Example usage:
//
//	type User struct {
//	    Name  string `bson:"name"`
//	    Email string `bson:"email"`
//	}
//
//	session := NewSession()
//	var user User
//	err := session.FindOne(&user)
//	if err != nil {
//	    // handle error
//	}
//
//	// use the user data
//	fmt.Println(user.Name, user.Email)
//
// Note:
// The function uses the "collectionForStruct" method to get the collection associated with "doc".
// It then applies the provided filters using the "filter.Filters" method.
// If "ctx" is provided, it uses it as the context for the operation. Otherwise, it uses the default context.
// The function retrieves the first document that matches the filters using the "FindOne" method of the collection.
// If there is any error during the operation, it returns the error.
// Otherwise, it decodes the document into the provided "doc" using the "Decode" method of the result.
func (s *session) FindOne(doc interface{}, ctx ...context.Context) error {
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

// FindAll retrieves all documents from the collection and stores them in the provided slice.
// The slice pointer should point to a slice of structs that match the documents schema.
// The retrieved documents are unmarshalled and stored in the slice.
// If the provided context is not empty, it is used for the database operation.
// If an error occurs during the execution, it is returned. Otherwise, nil is returned.
func (s *session) FindAll(rowsSlicePtr interface{}, ctx ...context.Context) error {
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
func (s *session) InsertOne(doc interface{}, ctx ...context.Context) (primitive.ObjectID, error) {
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
func (s *session) InsertMany(docs interface{}, ctx ...context.Context) (*mongo.InsertManyResult, error) {
	coll, err := s.collectionForSlice(docs)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(docs)
	var many []interface{}
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
func (s *session) DeleteOne(doc interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
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
func (s *session) SoftDeleteOne(doc interface{}, ctx ...context.Context) error {
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

// DeleteMany deletes multiple documents from the collection.
// It takes the following parameters:
// - doc: The document or value representing the document to delete.
// - ctx: Optional context(s) for the operation.
// It returns:
// - *mongo.DeleteResult: The result of the delete operation.
// - error: If an error occurs during the delete operation.
// This method retrieves the collection for the specified document,
// converts the filter into a MongoDB filter, executes the delete operation
// using the specified context and delete options, and returns
// the result of the delete operation or an error if any.
func (s *session) DeleteMany(doc interface{}, ctx ...context.Context) (*mongo.DeleteResult, error) {
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
// The method takes an interface{} type parameter (doc) representing the document(s) to be soft deleted.
// It returns an error if any error occurs during the update operation.
func (s *session) SoftDeleteMany(doc interface{}, ctx ...context.Context) error {
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

// Clone creates a new instance of the session and returns it as a driver.Session.
// The new session has the same values for the db, engine, filter, and various options
// as the original session.
// The cloned session is independent from the original session and can be used separately.
func (s *session) Clone() driver.Session {
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
func (s *session) Limit(i int64) driver.Session {
	s.findOptions = append(s.findOptions, options.Find().SetLimit(i))
	return s
}

// SetReadConcern sets the value for the ReadConcern field.
func (s *session) SetReadConcern(rc *readconcern.ReadConcern) driver.Session {
	s.collOpts = append(s.collOpts, options.Collection().SetReadConcern(rc))

	return s
}

// SetCollWriteConcern sets the write concern for the collection in the current session.
// It appends the options.Collection().SetWriteConcern(wc) to the s.collOpts field.
// The updated session is returned.
func (s *session) SetCollWriteConcern(wc *writeconcern.WriteConcern) driver.Session {
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
func (s *session) SetCollReadPreference(rp *readpref.ReadPref) driver.Session {
	s.collOpts = append(s.collOpts, options.Collection().SetReadPreference(rp))
	return s
}

// SetCollRegistry sets the bsoncodec.Registry for the session's collection.
// It appends the options.Collection().SetRegistry() to the session's collOpts
// and returns the updated session.
func (s *session) SetCollRegistry(r *bsoncodec.Registry) driver.Session {
	s.collOpts = append(s.collOpts, options.Collection().SetRegistry(r))
	return s
}

// Skip sets the number of documents to skip before returning results.
// It adds the skip option to the find and findOne options in the session.
// The skip value is specified by the i parameter.
// It returns the session.
func (s *session) Skip(i int64) driver.Session {
	s.findOptions = append(s.findOptions, options.Find().SetSkip(i))
	s.findOneOptions = append(s.findOneOptions, options.FindOne().SetSkip(i))
	return s
}

func (s *session) Count(i interface{}, ctx ...context.Context) (int64, error) {
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
func (s *session) UpdateOne(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
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
func (s *session) UpdateOneBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
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
func (s *session) UpdateManyBson(coll interface{}, bson interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
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

func (s *session) toBson(obj interface{}) bson.M {
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

func (s *session) makeValue(field string, value interface{}, ret bson.M) {
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

func (s *session) UpdateMany(bean interface{}, ctx ...context.Context) (*mongo.UpdateResult, error) {
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

func (s *session) RegexFilter(key, pattern string) driver.Session {
	s.filter.RegexFilter(key, pattern)
	return s
}

// ID sets the filter condition on the session's filter object to filter records by their ID.
// The provided 'id' parameter specifies the ID value to filter by.
// The method then returns the session object itself for method chaining.
func (s *session) ID(id interface{}) driver.Session {
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
func (s *session) Asc(colNames ...string) driver.Session {
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
func (s *session) Desc(colNames ...string) driver.Session {
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
func (s *session) Sort(colNames ...string) driver.Session {
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

func (s *session) Filter(key string, value interface{}) driver.Session {
	return s.Eq(key, value)
}

// Eq Equals a Specified Value
// { qty: 20 }
// Field in Embedded Document Equals a Value
// {"item.name": "ab" }
// Equals an Array Value
// { tags: [ "A", "B" ] }
func (s *session) Eq(key string, value interface{}) driver.Session {
	s.filter.Eq(key, value)
	return s
}

// Gt {field: {$gt: value} } >
func (s *session) Gt(key string, gt interface{}) driver.Session {
	s.filter.Gt(key, gt)
	return s
}

// Gte { qty: { $gte: 20 } } >=
func (s *session) Gte(key string, gte interface{}) driver.Session {
	s.filter.Gte(key, gte)
	return s
}

// In { field: { $in: [<value1>, <value2>, ... <valueN> ] } }
// tags: { $in: [ /^be/, /^st/ ] } }
// in []string []int ...
func (s *session) In(key string, in interface{}) driver.Session {
	s.filter.In(key, in)
	return s
}

// Lt {field: {$lt: value} } <
func (s *session) Lt(key string, lt interface{}) driver.Session {
	s.filter.Lt(key, lt)
	return s
}

// Lte { field: { $lte: value} } <=
func (s *session) Lte(key string, lte interface{}) driver.Session {
	s.filter.Lte(key, lte)
	return s
}

// Ne {field: {$ne: value} } !=
func (s *session) Ne(key string, ne interface{}) driver.Session {
	s.filter.Ne(key, ne)
	return s
}

// Nin { field: { $nin: [ <value1>, <value2> ... <valueN> ]} } the field does not exist.
func (s *session) Nin(key string, nin interface{}) driver.Session {
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
func (s *session) And(c driver.Condition) driver.Session {
	s.filter.And(c)
	return s

}

// Not { field: { $not: { <operator-expression> } } }
// not and Regular Expressions
// { item: { $not: /^p.*/ } }
func (s *session) Not(key string, not interface{}) driver.Session {
	s.filter.Not(key, not)
	return s
}

// Nor { $nor: [ { price: 1.99 }, { price: { $exists: false } },
// { sale: true }, { sale: { $exists: false } } ] }
// price != 1.99 || sale != true || sale exists || sale exists
func (s *session) Nor(c driver.Condition) driver.Session {
	s.filter.Nor(c)
	return s
}

// Or { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
func (s *session) Or(c driver.Condition) driver.Session {
	s.filter.Or(c)
	return s
}

func (s *session) Exists(key string, exists bool, filter ...driver.Condition) driver.Session {
	s.filter.Exists(key, exists, filter...)
	return s
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (s *session) SetArrayFilters(filters options.ArrayFilters) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetArrayFilters(filters))
	s.updateOpts = append(s.updateOpts, options.Update().SetArrayFilters(filters))
	return s
}

// SetOrdered sets the value for the Ordered field.
func (s *session) SetOrdered(ordered bool) driver.Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetOrdered(ordered))
	return s
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (s *session) SetBypassDocumentValidation(b bool) driver.Session {
	s.bulkWriteOptions = append(s.bulkWriteOptions, options.BulkWrite().SetBypassDocumentValidation(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetBypassDocumentValidation(b))
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts, options.FindOneAndUpdate().SetBypassDocumentValidation(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetBypassDocumentValidation(b))

	return s
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (s *session) SetReturnDocument(rd options.ReturnDocument) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetReturnDocument(rd))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetReturnDocument(rd))
	return s
}

// SetUpsert sets the value for the Upsert field.
func (s *session) SetUpsert(b bool) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetUpsert(b))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetUpsert(b))
	s.updateOpts = append(s.updateOpts, options.Update().SetUpsert(b))
	return s
}

// SetCollation sets the value for the Collation field.
func (s *session) SetCollation(collation *options.Collation) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetCollation(collation))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetCollation(collation))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetCollation(collation))
	s.updateOpts = append(s.updateOpts, options.Update().SetCollation(collation))
	return s
}

// SetMaxTime sets the value for the MaxTime field.
func (s *session) SetMaxTime(d time.Duration) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetMaxTime(d))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetMaxTime(d))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetMaxTime(d))
	return s
}

// SetProjection sets the projection for findOneAndUpdate, findOneAndReplace, and findOneAndDelete operations in the session.
// The projection parameter specifies which fields to include or exclude in the result documents.
// The projection value should be a document with field names as keys and a value of 1 to include the field in the result,
// or a value of 0 to exclude the field from the result.
// This method returns the session itself to allow for method chaining.
func (s *session) SetProjection(projection interface{}) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetProjection(projection))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetProjection(projection))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetProjection(projection))
	return s
}

// SetSort sets the sort order for the find one and update, find one and replace,
// and find one and delete operations in the session.
//
// The sort parameter specifies the sorting criteria in the form of a document.
// The document should contain field-value pairs, where the field is the name of the field
// by which to sort and the value is either 1 for ascending sort or -1 for descending sort.
// Example:
//
//	session.SetSort(bson.D{{"name", 1}}) // sort by name in ascending order
//	session.SetSort(bson.D{{"age", -1}}) // sort by age in descending order
//
// This method appends the specified sort option to the find one and update, find one and replace,
// and find one and delete options in the session. The options are used when executing the operations
// on the collection.
//
// This method returns the session itself, allowing for method chaining.
func (s *session) SetSort(sort interface{}) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetSort(sort))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetSort(sort))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetSort(sort))
	return s
}

// SetHint appends the hint option to the `findOneAndUpdateOpts`, `findOneAndReplaceOpts`, `findOneAndDeleteOpts`, and `updateOpts` slices of the session and returns the modified session
func (s *session) SetHint(hint interface{}) driver.Session {
	s.findOneAndUpdateOpts = append(s.findOneAndUpdateOpts,
		options.FindOneAndUpdate().SetHint(hint))
	s.findOneAndReplaceOpts = append(s.findOneAndReplaceOpts,
		options.FindOneAndReplace().SetHint(hint))
	s.findOneAndDeleteOpts = append(s.findOneAndDeleteOpts, options.FindOneAndDelete().SetHint(hint))
	s.updateOpts = append(s.updateOpts, options.Update().SetHint(hint))
	return s
}

// Type sets a type filter for the specified key in the driver session.
// The type filter is used to restrict the types of documents retrieved
// during a database query.
// This method takes in two parameters: key, which is the key used to specify
// the type filter, and t, which is the interface{} representing the type filter.
// The function returns the driver.Session object.
func (s *session) Type(key string, t interface{}) driver.Session {
	s.filter.Type(key, t)
	return s
}

// Expr applies the given condition to the current filter of the session.
// It returns the session itself to allow method chaining.
func (s *session) Expr(c driver.Condition) driver.Session {
	s.filter.Expr(c)
	return s
}

// Regex todo 简单实现，后续增加支持
func (s *session) Regex(key string, value string) driver.Session {
	s.filter.Regex(key, value)
	return s
}

// SetDatabase sets the database name for the session.
// It takes a string argument representing the name of the database.
// It updates the session's `db` field with the provided name and returns the updated session object.
func (s *session) SetDatabase(db string) driver.Session {
	s.db = db
	return s
}

func (s *session) collectionForStruct(doc interface{}) (*mongo.Collection, error) {
	coll, err := s.engine.CollectionNameForStruct(doc)
	if err != nil {
		return nil, err
	}

	return s.collectionByName(coll.Name), nil
}

func (s *session) collectionForSlice(doc interface{}) (*mongo.Collection, error) {
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

//func (s *session) makeFilterValue(field string, value interface{}) {
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
