package pie

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

type person struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}

func TestFilterByPointer(t *testing.T) {
	Convey("FilterBy should accept struct pointer", t, func() {
		f := DefaultCondition()
		p := &person{ID: "123"}
		f.FilterBy(p)
		So(f.Err(), ShouldBeNil)
		d, err := f.Filters()
		So(err, ShouldBeNil)
		So(d, ShouldResemble, bson.D{{Key: "_id", Value: "123"}})
	})
}
