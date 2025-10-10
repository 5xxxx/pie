package pie

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ChangeStreamWatcher change stream watcher
type ChangeStreamWatcher[T any] struct {
	engine     *Engine
	collection *mongo.Collection
	database   *mongo.Database
	stream     *mongo.ChangeStream
	pipeline   []bson.D
	options    *options.ChangeStreamOptionsBuilder
	watchType  WatchType
}

// WatchType watch type
type WatchType int

const (
	WatchCollection WatchType = iota // Watch collection
	WatchDatabase                    // Watch database
	WatchDeployment                  // Watch entire deployment
)

// ChangeEvent change event
type ChangeEvent[T any] struct {
	OperationType     string         `bson:"operationType"`
	FullDocument      *T             `bson:"fullDocument,omitempty"`
	DocumentKey       bson.D         `bson:"documentKey"`
	UpdateDescription *UpdateDesc    `bson:"updateDescription,omitempty"`
	ClusterTime       bson.Timestamp `bson:"clusterTime"`
	Namespace         Namespace      `bson:"ns"`
	ResumeToken       bson.Raw       `bson:"_id"`
}

// UpdateDesc update description
type UpdateDesc struct {
	UpdatedFields bson.D   `bson:"updatedFields"`
	RemovedFields []string `bson:"removedFields"`
}

// Namespace namespace
type Namespace struct {
	DB   string `bson:"db"`
	Coll string `bson:"coll"`
}

// Collection set the collection to watch (by name)
func (w *ChangeStreamWatcher[T]) Collection(name string) *ChangeStreamWatcher[T] {
	w.collection = w.engine.Collection(name)
	w.watchType = WatchCollection
	return w
}

// CollectionForStruct set the collection to watch (by struct)
func (w *ChangeStreamWatcher[T]) CollectionForStruct(v interface{}) *ChangeStreamWatcher[T] {
	collection, err := w.engine.CollectionForStruct(v)
	if err == nil {
		w.collection = collection
		w.watchType = WatchCollection
	}
	return w
}

// MatchOperationType match operation type
func (w *ChangeStreamWatcher[T]) MatchOperationType(operationType string) *ChangeStreamWatcher[T] {
	w.pipeline = append(w.pipeline, bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: operationType},
		}},
	})
	return w
}

// MatchInsert watch insert operations only
func (w *ChangeStreamWatcher[T]) MatchInsert() *ChangeStreamWatcher[T] {
	return w.MatchOperationType("insert")
}

// MatchUpdate watch update operations only
func (w *ChangeStreamWatcher[T]) MatchUpdate() *ChangeStreamWatcher[T] {
	return w.MatchOperationType("update")
}

// MatchReplace watch replace operations only
func (w *ChangeStreamWatcher[T]) MatchReplace() *ChangeStreamWatcher[T] {
	return w.MatchOperationType("replace")
}

// MatchDelete watch delete operations only
func (w *ChangeStreamWatcher[T]) MatchDelete() *ChangeStreamWatcher[T] {
	return w.MatchOperationType("delete")
}

// MatchFilter custom filter condition
func (w *ChangeStreamWatcher[T]) MatchFilter(filter bson.D) *ChangeStreamWatcher[T] {
	w.pipeline = append(w.pipeline, bson.D{
		{Key: "$match", Value: filter},
	})
	return w
}

// FullDocument set full document option
func (w *ChangeStreamWatcher[T]) FullDocument(option options.FullDocument) *ChangeStreamWatcher[T] {
	w.options.SetFullDocument(option)
	return w
}

// ResumeAfter resume watching from specified resume token
func (w *ChangeStreamWatcher[T]) ResumeAfter(resumeToken bson.Raw) *ChangeStreamWatcher[T] {
	w.options.SetResumeAfter(resumeToken)
	return w
}

// StartAfter start watching after specified token
func (w *ChangeStreamWatcher[T]) StartAfter(startAfter bson.Raw) *ChangeStreamWatcher[T] {
	w.options.SetStartAfter(startAfter)
	return w
}

// StartAtOperationTime start watching from specified operation time
func (w *ChangeStreamWatcher[T]) StartAtOperationTime(ts bson.Timestamp) *ChangeStreamWatcher[T] {
	w.options.SetStartAtOperationTime(&ts)
	return w
}

// BatchSize set batch size
func (w *ChangeStreamWatcher[T]) BatchSize(size int32) *ChangeStreamWatcher[T] {
	w.options.SetBatchSize(size)
	return w
}

// Comment set comment
func (w *ChangeStreamWatcher[T]) Comment(comment string) *ChangeStreamWatcher[T] {
	w.options.SetComment(comment)
	return w
}

// ShowExpandedEvents show expanded events
func (w *ChangeStreamWatcher[T]) ShowExpandedEvents(show bool) *ChangeStreamWatcher[T] {
	w.options.SetShowExpandedEvents(show)
	return w
}

// Watch start watching change stream
func (w *ChangeStreamWatcher[T]) Watch(ctx context.Context) error {
	var err error

	switch w.watchType {
	case WatchCollection:
		if w.collection == nil {
			collection, cerr := w.engine.CollectionForStruct((*T)(nil))
			if cerr != nil {
				return fmt.Errorf("failed to get collection: %w", cerr)
			}
			w.collection = collection
		}

		w.stream, err = w.collection.Watch(ctx, w.pipeline, w.options)
	case WatchDatabase:
		if w.database == nil {
			w.database = w.engine.database
		}
		w.stream, err = w.database.Watch(ctx, w.pipeline, w.options)
	default:
		return fmt.Errorf("unsupported watch type: %d", w.watchType)
	}

	return err
}

// Next move to next change event
func (w *ChangeStreamWatcher[T]) Next(ctx context.Context) bool {
	if w.stream == nil {
		return false
	}
	return w.stream.Next(ctx)
}

// TryNext try to move to next event (non-blocking)
func (w *ChangeStreamWatcher[T]) TryNext(ctx context.Context) bool {
	if w.stream == nil {
		return false
	}
	return w.stream.TryNext(ctx)
}

// Decode decode current change event
func (w *ChangeStreamWatcher[T]) Decode(event *ChangeEvent[T]) error {
	if w.stream == nil {
		return fmt.Errorf("change stream not initialized")
	}
	return w.stream.Decode(event)
}

// Current get current event's raw BSON
func (w *ChangeStreamWatcher[T]) Current() bson.Raw {
	if w.stream == nil {
		return nil
	}
	return w.stream.Current
}

// ResumeToken get current resume token
func (w *ChangeStreamWatcher[T]) ResumeToken() bson.Raw {
	if w.stream == nil {
		return nil
	}
	return w.stream.ResumeToken()
}

// Err return change stream error
func (w *ChangeStreamWatcher[T]) Err() error {
	if w.stream == nil {
		return nil
	}
	return w.stream.Err()
}

// Close close change stream
func (w *ChangeStreamWatcher[T]) Close(ctx context.Context) error {
	if w.stream == nil {
		return nil
	}
	return w.stream.Close(ctx)
}

// ID return change stream ID
func (w *ChangeStreamWatcher[T]) ID() int64 {
	if w.stream == nil {
		return 0
	}
	return w.stream.ID()
}

// Listen listen to change stream and execute callback function for each event
func (w *ChangeStreamWatcher[T]) Listen(ctx context.Context, handler func(*ChangeEvent[T]) error) error {
	if err := w.Watch(ctx); err != nil {
		return fmt.Errorf("failed to start watch: %w", err)
	}
	defer w.Close(ctx)

	for w.Next(ctx) {
		var event ChangeEvent[T]
		if err := w.Decode(&event); err != nil {
			return fmt.Errorf("failed to decode event: %w", err)
		}

		if err := handler(&event); err != nil {
			return fmt.Errorf("handler error: %w", err)
		}
	}

	if err := w.Err(); err != nil {
		return fmt.Errorf("change stream error: %w", err)
	}

	return nil
}

// ListenAsync asynchronously listen to change stream
func (w *ChangeStreamWatcher[T]) ListenAsync(ctx context.Context, handler func(*ChangeEvent[T]) error, errChan chan<- error) {
	go func() {
		if err := w.Listen(ctx, handler); err != nil {
			if errChan != nil {
				errChan <- err
			}
		}
	}()
}

// ForEach iterate over all events in change stream
func (w *ChangeStreamWatcher[T]) ForEach(ctx context.Context, handler func(*ChangeEvent[T]) error) error {
	return w.Listen(ctx, handler)
}

// WatchOptions watch options builder (advanced usage)
type WatchOptions struct {
	FullDocument             options.FullDocument
	FullDocumentBeforeChange string
	ResumeAfter              bson.Raw
	StartAfter               bson.Raw
	StartAtOperationTime     *bson.Timestamp
	BatchSize                int32
	Comment                  string
	ShowExpandedEvents       bool
}

// ApplyWatchOptions apply watch options
func (w *ChangeStreamWatcher[T]) ApplyWatchOptions(opts *WatchOptions) *ChangeStreamWatcher[T] {
	if opts == nil {
		return w
	}

	if opts.FullDocument != "" {
		w.FullDocument(opts.FullDocument)
	}
	if opts.FullDocumentBeforeChange != "" {
		// FullDocumentBeforeChange not supported in v2
	}
	if opts.ResumeAfter != nil {
		w.ResumeAfter(opts.ResumeAfter)
	}
	if opts.StartAfter != nil {
		w.StartAfter(opts.StartAfter)
	}
	if opts.StartAtOperationTime != nil {
		w.StartAtOperationTime(*opts.StartAtOperationTime)
	}
	if opts.BatchSize > 0 {
		w.BatchSize(opts.BatchSize)
	}
	if opts.Comment != "" {
		w.Comment(opts.Comment)
	}
	if opts.ShowExpandedEvents {
		w.ShowExpandedEvents(true)
	}

	return w
}
