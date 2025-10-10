package pie

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Transaction transaction manager
type Transaction struct {
	engine *Engine
}

// NewTransaction create new transaction manager
func NewTransaction(engine *Engine) *Transaction {
	return &Transaction{
		engine: engine,
	}
}

// TransactionFunc transaction function type
type TransactionFunc func(context.Context) error

// TransactionWithOptions execute transaction with options
func (t *Transaction) TransactionWithOptions(ctx context.Context, fn TransactionFunc, opts ...options.Lister[options.SessionOptions]) error {
	session, err := t.engine.Client().StartSession(opts...)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// Transaction execute transaction
func (t *Transaction) Transaction(ctx context.Context, fn TransactionFunc) error {
	return t.TransactionWithOptions(ctx, fn)
}

// TransactionWithResult execute transaction and return result
func TransactionWithResult[T any](t *Transaction, ctx context.Context, fn func(context.Context) (T, error)) (*TransactionResult[T], error) {
	session, err := t.engine.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result T
	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		var err error
		result, err = fn(sessCtx)
		return result, err
	})

	return &TransactionResult[T]{
		Data:  result,
		Error: err,
	}, err
}

// TransactionWithResultAndOptions execute transaction with options and return result
func TransactionWithResultAndOptions[T any](t *Transaction, ctx context.Context, fn func(context.Context) (T, error), opts ...options.Lister[options.SessionOptions]) (*TransactionResult[T], error) {
	session, err := t.engine.Client().StartSession(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result T
	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		var err error
		result, err = fn(sessCtx)
		return result, err
	})

	return &TransactionResult[T]{
		Data:  result,
		Error: err,
	}, err
}

// Session get transaction session
func (t *Transaction) Session(ctx context.Context) (*mongo.Session, error) {
	session, err := t.engine.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	return session, nil
}

// SessionWithOptions get transaction session with options
func (t *Transaction) SessionWithOptions(ctx context.Context, opts ...options.Lister[options.SessionOptions]) (*mongo.Session, error) {
	session, err := t.engine.Client().StartSession(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	return session, nil
}

// TransactionOptions transaction options builder
type TransactionOptions struct {
	options *options.TransactionOptionsBuilder
}

// NewTransactionOptions create new transaction options
func NewTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		options: options.Transaction(),
	}
}

// Build build transaction options
func (to *TransactionOptions) Build() *options.TransactionOptionsBuilder {
	return to.options
}

// TransactionSession transaction session wrapper
type TransactionSession struct {
	session *mongo.Session
	ctx     context.Context
}

// NewTransactionSession create new transaction session
func NewTransactionSession(session *mongo.Session, ctx context.Context) *TransactionSession {
	return &TransactionSession{
		session: session,
		ctx:     ctx,
	}
}

// WithTransaction execute transaction
func (ts *TransactionSession) WithTransaction(fn TransactionFunc) error {
	_, err := ts.session.WithTransaction(ts.ctx, func(sessCtx context.Context) (interface{}, error) {
		return nil, fn(sessCtx)
	})
	return err
}

// WithTransactionResult execute transaction and return result
func WithTransactionResult[T any](ts *TransactionSession, fn func(context.Context) (T, error)) (*TransactionResult[T], error) {
	var result T
	_, err := ts.session.WithTransaction(ts.ctx, func(sessCtx context.Context) (interface{}, error) {
		var err error
		result, err = fn(sessCtx)
		return result, err
	})

	return &TransactionResult[T]{
		Data:  result,
		Error: err,
	}, err
}

// AbortTransaction abort transaction
func (ts *TransactionSession) AbortTransaction(ctx context.Context) error {
	return ts.session.AbortTransaction(ctx)
}

// CommitTransaction commit transaction
func (ts *TransactionSession) CommitTransaction(ctx context.Context) error {
	return ts.session.CommitTransaction(ctx)
}

// EndSession end session
func (ts *TransactionSession) EndSession(ctx context.Context) {
	ts.session.EndSession(ctx)
}

// GetSession get original session
func (ts *TransactionSession) GetSession() *mongo.Session {
	return ts.session
}

// GetContext get context
func (ts *TransactionSession) GetContext() context.Context {
	return ts.ctx
}
