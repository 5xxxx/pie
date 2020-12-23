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
)

func TestNewAggregate(t *testing.T) {
	type args struct {
		engine driver.Client
	}
	tests := []struct {
		name string
		args args
		want driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAggregate(tt.args.engine); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_All(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		ctx    context.Context
		result interface{}
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
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if err := a.All(tt.args.ctx, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_aggregate_Collection(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.Collection(tt.args.doc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_Match(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		c driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.Match(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_One(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		ctx    context.Context
		result interface{}
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
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if err := a.One(tt.args.ctx, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("One() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_aggregate_Pipeline(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		pipeline bson.A
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.Pipeline(tt.args.pipeline); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetAllowDiskUse(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetAllowDiskUse(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAllowDiskUse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetBatchSize(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		i int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetBatchSize(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetBatchSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetBypassDocumentValidation(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetBypassDocumentValidation(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetBypassDocumentValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetCollation(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		c *options.Collation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetCollation(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCollation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetComment(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetComment(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetDatabase(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		db string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetDatabase(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetHint(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		h interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetHint(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetHint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetMaxAwaitTime(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetMaxAwaitTime(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetMaxAwaitTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_SetMaxTime(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
	}
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.SetMaxTime(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetMaxTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_collectionByName(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
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
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			if got := a.collectionByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectionByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregate_collectionForSlice(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
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
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			got, err := a.collectionForSlice(tt.args.doc)
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

func Test_aggregate_collectionForStruct(t *testing.T) {
	type fields struct {
		db       string
		doc      interface{}
		engine   driver.Client
		pipeline bson.A
		opts     []*options.AggregateOptions
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
			a := &aggregate{
				db:       tt.fields.db,
				doc:      tt.fields.doc,
				engine:   tt.fields.engine,
				pipeline: tt.fields.pipeline,
				opts:     tt.fields.opts,
			}
			got, err := a.collectionForStruct(tt.args.doc)
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
