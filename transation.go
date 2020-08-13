/*
 *
 * transation.go
 * tugrik
 *
 * Created by lintao on 2020/8/11 8:33 上午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package tugrik

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransFunc func(context.Context) error

func (t Tugrik) TransactionWithOptions(ctx context.Context, opt *options.SessionOptions, f TransFunc) error {
	session, err := t.client.StartSession(opt)
	if err != nil {
		return err
	}
	if err = f(ctx); err != nil {
		return session.AbortTransaction(ctx)
	}

	return session.CommitTransaction(ctx)
}

func (t Tugrik) Transaction(ctx context.Context, f TransFunc) error {
	opts := options.Session().
		SetDefaultReadConcern(readconcern.Majority()).
		SetDefaultWriteConcern(writeconcern.New(writeconcern.WMajority()))
	return t.TransactionWithOptions(ctx, opts, f)
}
