package internal

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/5xxxx/pie/driver"
	"github.com/5xxxx/pie/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var URI = "mongodb://192.168.1.208:10001"

type Account struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	NickName     string             `json:"nick_name" bson:"nick_name,omitempty"`
	MobileNumber string             `json:"mobile_number" bson:"mobile_number,omitempty"`
	Gender       int                `json:"gender" bson:"gender,omitempty"`
	Birthday     string             `json:"birthday" bson:"birthday,omitempty"`
	GeoPoint     []float64          `json:"geo_point" bson:"geo_point,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

func newClient() driver.Client {
	d, err := NewClient("pie", options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	if err = d.Connect(context.Background()); err != nil {
		panic(err)
	}

	return d
}

func TestNewClient(t *testing.T) {
	type args struct {
		db   string
		opts []*options.ClientOptions
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "init client",
			args: args{
				db:   "pie",
				opts: []*options.ClientOptions{options.Client().ApplyURI(URI)},
			},
			want:    "pie",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.db, tt.args.opts...)
			if err = got.Connect(context.Background()); err != nil {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.DataBase().Name() != tt.want {
				t.Errorf("NewClient() error = %v, want %v", err, tt.want)
				return
			}
		})
	}
}

func Test_defaultDriver_InsertMany(t *testing.T) {

	type args struct {
		ctx context.Context
		v   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "InsertMany",
			args: args{
				ctx: context.Background(),
				v: []Account{
					{
						NickName:     "xiaoming",
						MobileNumber: "13888888888",
						Gender:       1,
						Birthday:     "1970-01-12",
						GeoPoint:     []float64{12.3, 12.4},
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
					{
						NickName:     "xiaohong",
						MobileNumber: "18888888888",
						Gender:       2,
						Birthday:     "1990-03-23",
						GeoPoint:     []float64{53.3, 167.4},
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newClient()
			_, err := d.InsertMany(tt.args.v, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("InsertMany() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_defaultDriver_InsertOne(t *testing.T) {
	id := primitive.NewObjectID()
	type args struct {
		ctx context.Context
		v   Account
	}
	tests := []struct {
		name    string
		args    args
		want    primitive.ObjectID
		wantErr bool
	}{
		{
			name: "Insert One",
			args: args{
				ctx: context.Background(),
				v: Account{
					Id:           id,
					NickName:     "xiaolv",
					MobileNumber: "139999999999",
					Gender:       3,
					Birthday:     "2000-03-23",
					GeoPoint:     []float64{53.3, 167.4},
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			},
			want:    id,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewClient("pie", options.Client().ApplyURI(URI))
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = d.Connect(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("InsertOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := d.InsertOne(&tt.args.v, tt.args.ctx)
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

func Test_defaultDriver_FindAll(t *testing.T) {
	type args struct {
		ctx  context.Context
		docs []Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "find all",
			args: args{
				ctx:  context.Background(),
				docs: make([]Account, 0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewClient("pie", options.Client().ApplyURI(URI))
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = d.Connect(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = d.FindAll(&tt.args.docs, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_defaultDriver_FindAndDelete(t *testing.T) {

	type args struct {
		ctx context.Context
		doc Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FindAndDelete",
			args: args{
				ctx: context.Background(),
				doc: Account{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewClient("pie", options.Client().ApplyURI(URI))
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAndDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = d.Connect(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("FindAndDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := d.FindAndDelete(&tt.args.doc, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("FindAndDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_FindOne(t *testing.T) {

	type args struct {
		ctx context.Context
		doc interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "FindOne",
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewClient("pie", options.Client().ApplyURI(URI))
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = d.Connect(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := d.FindOne(tt.args.doc, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_FindOneAndReplace(t *testing.T) {

	type args struct {
		ctx context.Context
		doc Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FindOneAndReplace",
			args: args{
				ctx: context.Background(),
				doc: Account{
					NickName: "xiaolin",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewClient("pie", options.Client().ApplyURI(URI))
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = d.Connect(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := d.FindOneAndReplace(&tt.args.doc, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("FindOneAndReplace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_FindOneAndUpdate(t *testing.T) {

	type args struct {
		ctx context.Context
		doc Account
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.SingleResult
		wantErr bool
	}{
		{
			name: "FindOneAndUpdate()",
			args: args{
				ctx: context.Background(),
				doc: Account{
					NickName:  "xiaozhang",
					Gender:    4,
					CreatedAt: time.Now(),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newClient()
			_, err := d.FindOneAndUpdate(&tt.args.doc, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneAndUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("FindOneAndUpdate() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_defaultDriver_FindOneAndUpdateBson(t *testing.T) {

	type args struct {
		ctx  context.Context
		coll Account
		bson interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.SingleResult
		wantErr bool
	}{
		{
			name: "FindOneAndUpdateBso",
			args: args{
				ctx:  context.Background(),
				coll: Account{},
				bson: bson.M{
					"$set": bson.M{"nick_name": "dilireba"},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newClient()
			_, err := d.FindOneAndUpdateBson(tt.args.bson, &tt.args.coll, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneAndUpdateBson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("FindOneAndUpdateBson() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

//func Test_defaultDriver_AddIndex(t *testing.T) {
//
//	type args struct {
//		keys interface{}
//		opt  []*options.IndexOptions
//	}
//	tests := []struct {
//		name string
//		args args
//		want driver.Indexes
//	}{
//		{
//			name: "AddIndex",
//			args: args{
//				keys: bson.M{"nick_name": 1},
//				opt:  []*options.IndexOptions{options.Index().SetBackground(true)},
//			},
//			want: []*options.IndexOptions{options.Index().SetBackground(true)},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := newClient()
//			if got := d.AddIndex(tt.args.keys, tt.args.opt...); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AddIndex() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_defaultDriver_Aggregate(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	tests := []struct {
		name   string
		fields fields
		want   driver.Aggregate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Aggregate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_And(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		filter driver.Condition
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.And(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Asc(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Asc(tt.args.colNames...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Asc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_BulkWrite(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.BulkWrite(tt.args.docs, tt.args.ctx)
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

func Test_defaultDriver_Collection(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Collection(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_CollectionNameForSlice(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schemas.Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.CollectionNameForSlice(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectionNameForSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectionNameForSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_CollectionNameForStruct(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schemas.Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.CollectionNameForStruct(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectionNameForStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectionNameForStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Connect(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		ctx context.Context
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_Count(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.Count(tt.args.i)
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

func Test_defaultDriver_DataBase(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	tests := []struct {
		name   string
		fields fields
		want   *mongo.Database
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.DataBase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_DeleteMany(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		ctx    context.Context
		filter interface{}
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.DeleteMany(tt.args.filter, tt.args.ctx)
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

func Test_defaultDriver_DeleteOne(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		ctx    context.Context
		filter interface{}
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.DeleteOne(tt.args.filter, tt.args.ctx)
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

func Test_defaultDriver_Desc(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		s2 []string
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Desc(tt.args.s2...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Desc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Disconnect(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		ctx context.Context
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.Disconnect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_Distinct(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.Distinct(tt.args.doc, tt.args.columns, tt.args.ctx)
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

func Test_defaultDriver_DropAll(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.DropAll(tt.args.doc, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DropAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_DropOne(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.DropOne(tt.args.doc, tt.args.name, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DropOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_Eq(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Eq(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Exists(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Exists(tt.args.key, tt.args.exists, tt.args.filter...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Expr(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		filter driver.Condition
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Expr(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Filter(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Filter(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_FilterBy(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.FilterBy(tt.args.object); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Gt(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Gt(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Gte(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Gte(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_ID(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.ID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_In(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.In(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Limit(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		limit int64
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Limit(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Lt(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Lt(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Lte(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Lte(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Ne(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Ne(tt.args.key, tt.args.ne); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_NewIndexes(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.NewIndexes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_NewSession(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.NewSession(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Nin(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Nin(tt.args.key, tt.args.nin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Nor(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Nor(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Not(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Not(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Not() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Or(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		filter driver.Condition
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Or(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Ping(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_Regex(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		key   string
		value string
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Regex(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Regex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_RegexFilter(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.RegexFilter(tt.args.key, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_ReplaceOne(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.ReplaceOne(tt.args.doc, tt.args.ctx)
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

func Test_defaultDriver_Skip(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		skip int64
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Skip(tt.args.skip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_SoftDeleteMany(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		ctx    context.Context
		filter interface{}
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.SoftDeleteMany(tt.args.filter, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SoftDeleteMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_SoftDeleteOne(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		ctx    context.Context
		filter interface{}
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if err := d.SoftDeleteOne(tt.args.filter, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SoftDeleteOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDriver_Type(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			if got := d.Type(tt.args.key, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_Update(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.Update(tt.args.bean, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDriver_UpdateMany(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.UpdateMany(tt.args.bean, tt.args.ctx)
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

func Test_defaultDriver_UpdateManyBson(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.UpdateManyBson(tt.args.coll, tt.args.bson, tt.args.ctx)
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

func Test_defaultDriver_UpdateOneBson(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
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
			d := defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.UpdateOneBson(tt.args.coll, tt.args.bson, tt.args.ctx)
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

func Test_defaultDriver_getStructCollAndSetKey(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		parser     *driver.Parser
		db         string
		clientOpts []*options.ClientOptions
	}
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schemas.Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultClient{
				client:     tt.fields.client,
				parser:     tt.fields.parser,
				db:         tt.fields.db,
				clientOpts: tt.fields.clientOpts,
			}
			got, err := d.getStructCollAndSetKey(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStructCollAndSetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStructCollAndSetKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
