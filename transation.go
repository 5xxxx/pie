/*
 *
 * transation.go
 * tugrik
 *
 * Created by lintao on 2020/8/11 8:33 上午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransFunc func(context.Context) error

func (t Driver) TransactionWithOptions(ctx context.Context, f TransFunc, opt ...*options.SessionOptions) error {
	session, err := t.client.StartSession(opt...)
	if err != nil {
		return err
	}
	if err = f(ctx); err != nil {
		return session.AbortTransaction(ctx)
	}
	defer session.EndSession(context.Background())
	return session.CommitTransaction(ctx)
}

func (t Driver) Transaction(ctx context.Context, f TransFunc) error {
	opts := options.Session().
		SetDefaultReadConcern(readconcern.Majority()).
		SetDefaultWriteConcern(writeconcern.New(writeconcern.WMajority()))
	return t.TransactionWithOptions(ctx, f, []*options.SessionOptions{opts}...)
}
