package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/5xxxx/pie/schemas"

	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (d *defaultClient) Transaction(ctx context.Context, f schemas.TransFunc) error {
	opts := options.Session().
		SetDefaultReadConcern(readconcern.Majority())
	return d.TransactionWithOptions(ctx, f, []*options.SessionOptions{opts}...)
}
