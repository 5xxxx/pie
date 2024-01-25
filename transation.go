package pie

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/5xxxx/pie/schemas"

	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionWithOptions executes a transaction.
//
// TransactionWithOptions starts a session with the provided options and begins a transaction.
// It then calls the provided TransFunc with the created session, passing the transactional
// session context. If the TransFunc returns an error, the transaction is aborted.
// Otherwise, the transaction is committed automatically when the TransFunc completes without error.
// If an error occurs during the session creation or transaction execution, it is returned.
//
// Note: This method uses the default read preference of PrimaryPreferred.
//
// Parameters:
// - ctx: The context.Context object to use for the session and transaction.
// - f: The TransFunc that will be executed within the transaction.
// - opt: Optional session options to use when starting the session.
//
// Returns:
// - An error if the session creation or transaction execution fails.
func (d *defaultClient) TransactionWithOptions(ctx context.Context, f schemas.TransFunc, opt ...*options.SessionOptions) error {
	transaction, err := d.client.StartSession(opt...)
	if err != nil {
		return err
	}
	defer transaction.EndSession(context.Background())

	txnOpts := options.Transaction().
		SetReadPreference(readpref.PrimaryPreferred())
	_, err = transaction.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, f(sessCtx)
	}, txnOpts)

	return err
}

// Transaction is a method that executes a transaction with the provided TransFunc.
// It sets the default read concern to `readconcern.Majority()` and then calls the
// TransactionWithOptions method to perform the transaction using the specified context
// and session options.
//
// Parameters:
// - ctx: The context.Context for the transaction.
// - f: The TransFunc that will be executed within the transaction.
//
// Returns:
// - error: If an error occurs during the transaction, it will be returned.
//     Otherwise, it will return nil.
//
// Example usage:
//     err := client.Transaction(ctx, func(ctx context.Context) error {
//         // Perform the transaction logic here
//         return nil
//     })
//     if err != nil {
//         // Handle the error
//     }
//
// See the schemas.TransFunc documentation for more details on how to define the
// transaction function.
func (d *defaultClient) Transaction(ctx context.Context, f schemas.TransFunc) error {
	opts := options.Session().
		SetDefaultReadConcern(readconcern.Majority())
	return d.TransactionWithOptions(ctx, f, []*options.SessionOptions{opts}...)
}
