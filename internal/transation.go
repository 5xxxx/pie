package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransFunc func(context.Context) error

func (d defaultClient) TransactionWithOptions(ctx context.Context, f TransFunc, opt ...*options.SessionOptions) error {
	session, err := d.client.StartSession(opt...)
	if err != nil {
		return err
	}
	if err = f(ctx); err != nil {
		return session.AbortTransaction(ctx)
	}
	defer session.EndSession(context.Background())
	return session.CommitTransaction(ctx)
}

func (d defaultClient) Transaction(ctx context.Context, f TransFunc) error {
	opts := options.Session().
		SetDefaultReadConcern(readconcern.Majority()).
		SetDefaultWriteConcern(writeconcern.New(writeconcern.WMajority()))
	return d.TransactionWithOptions(ctx, f, []*options.SessionOptions{opts}...)
}
