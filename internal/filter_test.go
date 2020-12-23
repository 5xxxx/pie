package internal

import (
	"reflect"
	"testing"

	"github.com/NSObjects/pie/driver"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDefaultCondition(t *testing.T) {
	tests := []struct {
		name string
		want driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultCondition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_A(t *testing.T) {
	type fields struct {
		d bson.D
	}
	tests := []struct {
		name   string
		fields fields
		want   bson.A
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.A(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("A() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_And(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		filter driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.And(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Eq(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Eq(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Exists(t *testing.T) {
	type fields struct {
		d bson.D
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
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Exists(tt.args.key, tt.args.exists, tt.args.filter...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Expr(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		filter driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Expr(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Filters(t *testing.T) {
	type fields struct {
		d bson.D
	}
	tests := []struct {
		name   string
		fields fields
		want   bson.D
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Filters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Gt(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		gt  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Gt(tt.args.key, tt.args.gt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Gte(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		gte interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Gte(tt.args.key, tt.args.gte); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_ID(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		id interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.ID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_In(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		in  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.In(tt.args.key, tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Lt(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		lt  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Lt(tt.args.key, tt.args.lt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Lte(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		lte interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Lte(tt.args.key, tt.args.lte); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Ne(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		ne  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Ne(tt.args.key, tt.args.ne); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Nin(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		nin interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Nin(tt.args.key, tt.args.nin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Nor(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		filter driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Nor(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Not(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		not interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Not(tt.args.key, tt.args.not); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Not() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Or(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		filter driver.Condition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Or(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Regex(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Regex(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Regex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_RegexFilter(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key     string
		pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.RegexFilter(tt.args.key, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_Type(t *testing.T) {
	type fields struct {
		d bson.D
	}
	type args struct {
		key string
		t   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   driver.Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				d: tt.fields.d,
			}
			if got := f.Type(tt.args.key, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}
