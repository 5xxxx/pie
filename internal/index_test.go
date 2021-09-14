package internal

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/5xxxx/pie/driver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewIndexes(t *testing.T) {
	type args struct {
		driver driver.Client
	}
	tests := []struct {
		name string
		args args
		want driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIndexes(tt.args.driver); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_AddIndex(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		keys interface{}
		opt  []*options.IndexOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := i.AddIndex(tt.args.keys, tt.args.opt...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_Collection(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := c.Collection(tt.args.doc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_CreateIndexes(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			got, err := i.CreateIndexes(tt.args.doc, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateIndexes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_DropAll(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
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
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if err := i.DropAll(tt.args.doc, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DropAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_index_DropOne(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		ctx  context.Context
		doc  interface{}
		name string
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
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if err := i.DropOne(tt.args.doc, tt.args.name, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DropOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_index_SetCommitQuorumInt(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		quorum int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := i.SetCommitQuorumInt(tt.args.quorum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCommitQuorumInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_SetCommitQuorumMajority(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	tests := []struct {
		name   string
		fields fields
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := i.SetCommitQuorumMajority(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCommitQuorumMajority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_SetCommitQuorumString(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		quorum string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := i.SetCommitQuorumString(tt.args.quorum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCommitQuorumString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_SetCommitQuorumVotingMembers(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	tests := []struct {
		name   string
		fields fields
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := i.SetCommitQuorumVotingMembers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCommitQuorumVotingMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_SetDatabase(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		db string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := c.SetDatabase(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_SetMaxTime(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
	}
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Indexes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := i.SetMaxTime(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetMaxTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_collectionByName(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
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
			c := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
			}
			if got := c.collectionByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectionByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_collectionForSlice(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
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
			c := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
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

func Test_index_collectionForStruct(t *testing.T) {
	type fields struct {
		db                 string
		doc                interface{}
		engine             driver.Client
		indexes            []mongo.IndexModel
		createIndexOpts    []*options.CreateIndexesOptions
		dropIndexesOptions []*options.DropIndexesOptions
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
			c := &index{
				db:                 tt.fields.db,
				doc:                tt.fields.doc,
				engine:             tt.fields.engine,
				indexes:            tt.fields.indexes,
				createIndexOpts:    tt.fields.createIndexOpts,
				dropIndexesOptions: tt.fields.dropIndexesOptions,
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
