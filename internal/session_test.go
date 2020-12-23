/*
 *
 * session_test.go
 * tugrik
 *
 * Created by lintao on 2020/8/15 12:28 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package internal

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/NSObjects/pie/driver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccId        string             `json:"accId" bson:"accId"`
	NickName     string             `json:"nickName" bson:"nickName"`
	HeadImgUrl   string             `json:"headImgUrl" bson:"headImgUrl"`
	MobileNumber string             `json:"mobileNumber" bson:"mobileNumber"`
	Gender       int                `json:"gender" bson:"gender"`
	Age          int                `json:"age" bson:"age"`
	Birthday     string             `json:"birthday" bson:"birthday"`
	Location     Location           `json:"location" bson:"location,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}
type Location struct {
	Province string `json:"province" bson:"province"`
	City     string `json:"city" bson:"city"`
	County   string `json:"county" bson:"county"`
}

func TestSession_FilterBy(t *testing.T) {
	var user User
	user.NickName = "linlin"
	user.Location.City = "pingj"
	user.Location.Province = "hahha"
	type args struct {
		object interface{}
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "test filter",
			args: args{object: user},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tu, _ := NewClient("")
			s := NewSession(tu)
			if err := s.FilterBy(tt.args.object); (err != nil) != tt.wantErr {
				//t.Errorf("FilterBy() error = %v, wantErr %v", err, tt.wantErr)
			}

			//fmt.Println(s.filter)
		})
	}
}

func TestNewSession(t *testing.T) {
	type args struct {
		engine driver.Client
	}
	tests := []struct {
		name string
		args args
		want driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSession(tt.args.engine); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_And(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		c driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.And(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Asc(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		colNames []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Asc(tt.args.colNames...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Asc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_BulkWrite(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		docs interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.BulkWriteResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.BulkWrite(tt.args.ctx, tt.args.docs)
			if (err != nil) != tt.wantErr {
				t.Errorf("BulkWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BulkWrite() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Clone(t *testing.T) {
	type fields struct {
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
	}
	tests := []struct {
		name   string
		fields fields
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Count(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.Count(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Count() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_DeleteMany(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.DeleteResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.DeleteMany(tt.args.ctx, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteMany() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_DeleteOne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.DeleteResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.DeleteOne(tt.args.ctx, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Desc(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		colNames []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Desc(tt.args.colNames...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Desc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Distinct(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx     context.Context
		doc     interface{}
		columns string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.Distinct(tt.args.ctx, tt.args.doc, tt.args.columns)
			if (err != nil) != tt.wantErr {
				t.Errorf("Distinct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Distinct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Eq(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Eq(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Exists(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key    string
		exists bool
		filter []driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Exists(tt.args.key, tt.args.exists, tt.args.filter...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Expr(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		c driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Expr(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Filter(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Filter(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_FilterBy(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		object interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.FilterBy(tt.args.object); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_FindAll(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx          context.Context
		rowsSlicePtr interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if err := s.FindAll(tt.args.ctx, tt.args.rowsSlicePtr); (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_session_FindAndDelete(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if err := s.FindAndDelete(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("FindAndDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_session_FindOne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if err := s.FindOne(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_session_FindOneAndReplace(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if err := s.FindOneAndReplace(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("FindOneAndReplace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_session_FindOneAndUpdate(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.SingleResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.FindOneAndUpdate(tt.args.ctx, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneAndUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOneAndUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_FindOneAndUpdateBson(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		coll interface{}
		bson interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.SingleResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.FindOneAndUpdateBson(tt.args.ctx, tt.args.coll, tt.args.bson)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneAndUpdateBson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOneAndUpdateBson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Gt(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		gt  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Gt(tt.args.key, tt.args.gt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Gte(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		gte interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Gte(tt.args.key, tt.args.gte); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_ID(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		id interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.ID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_In(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		in  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.In(tt.args.key, tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_InsertMany(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		docs interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.InsertManyResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.InsertMany(tt.args.ctx, tt.args.docs)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertMany() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_InsertOne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    primitive.ObjectID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.InsertOne(tt.args.ctx, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Limit(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		i int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Limit(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Lt(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		lt  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Lt(tt.args.key, tt.args.lt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Lte(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		lte interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Lte(tt.args.key, tt.args.lte); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Ne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		ne  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Ne(tt.args.key, tt.args.ne); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Nin(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		nin interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Nin(tt.args.key, tt.args.nin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Nor(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		c driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Nor(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Not(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		not interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Not(tt.args.key, tt.args.not); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Not() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Or(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		c driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Or(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Regex(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Regex(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Regex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_RegexFilter(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key     string
		pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.RegexFilter(tt.args.key, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_ReplaceOne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.ReplaceOne(tt.args.ctx, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetArrayFilters(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		filters options.ArrayFilters
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetArrayFilters(tt.args.filters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetArrayFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetBypassDocumentValidation(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetBypassDocumentValidation(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetBypassDocumentValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetCollation(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		collation *options.Collation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetCollation(tt.args.collation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCollation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetDatabase(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		db string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := c.SetDatabase(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetHint(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		hint interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetHint(tt.args.hint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetHint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetMaxTime(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetMaxTime(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetMaxTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetOrdered(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ordered bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetOrdered(tt.args.ordered); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetOrdered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetProjection(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		projection interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetProjection(tt.args.projection); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetProjection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetReturnDocument(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		rd options.ReturnDocument
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetReturnDocument(tt.args.rd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetReturnDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetSort(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		sort interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetSort(tt.args.sort); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SetUpsert(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.SetUpsert(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUpsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Skip(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		i int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Skip(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_SoftDeleteMany(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if err := s.SoftDeleteMany(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("SoftDeleteMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_session_SoftDeleteOne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if err := s.SoftDeleteOne(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("SoftDeleteOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_session_Sort(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		colNames []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Sort(tt.args.colNames...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_Type(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		key string
		t   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.Type(tt.args.key, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_UpdateMany(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		bean interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.UpdateMany(tt.args.ctx, tt.args.bean)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMany() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_UpdateManyBson(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		coll interface{}
		bson interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.UpdateManyBson(tt.args.ctx, tt.args.coll, tt.args.bson)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateManyBson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateManyBson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_UpdateOne(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		bean interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.UpdateOne(tt.args.ctx, tt.args.bean)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_UpdateOneBson(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		ctx  context.Context
		coll interface{}
		bson interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := s.UpdateOneBson(tt.args.ctx, tt.args.coll, tt.args.bson)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateOneBson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateOneBson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_collectionByName(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mongo.Collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := c.collectionByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectionByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_collectionForSlice(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := c.collectionForSlice(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectionForSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectionForSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_collectionForStruct(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			got, err := c.collectionForStruct(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectionForStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectionForStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session_makeFilterValue(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
		})
	}
}

func Test_session_makeStruct(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		field string
		value reflect.Value
		ret   bson.M
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
		})
	}
}

func Test_session_makeStructValue(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		field string
		value reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
		})
	}
}

func Test_session_makeValue(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		field string
		value interface{}
		ret   bson.M
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
		})
	}
}

func Test_session_toBson(t *testing.T) {
	type fields struct {
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
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bson.M
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &session{
				db:                    tt.fields.db,
				engine:                tt.fields.engine,
				filter:                tt.fields.filter,
				findOneOptions:        tt.fields.findOneOptions,
				findOptions:           tt.fields.findOptions,
				insertManyOpts:        tt.fields.insertManyOpts,
				insertOneOpts:         tt.fields.insertOneOpts,
				deleteOpts:            tt.fields.deleteOpts,
				findOneAndDeleteOpts:  tt.fields.findOneAndDeleteOpts,
				updateOpts:            tt.fields.updateOpts,
				countOpts:             tt.fields.countOpts,
				distinctOpts:          tt.fields.distinctOpts,
				findOneAndReplaceOpts: tt.fields.findOneAndReplaceOpts,
				findOneAndUpdateOpts:  tt.fields.findOneAndUpdateOpts,
				replaceOpts:           tt.fields.replaceOpts,
				bulkWriteOptions:      tt.fields.bulkWriteOptions,
			}
			if got := s.toBson(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toBson() = %v, want %v", got, tt.want)
			}
		})
	}
}
