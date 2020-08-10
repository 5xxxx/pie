package tugrik

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_insertOmitemptyTag(t *testing.T) {
	type Person struct {
		ID   string `bson:"_id,omitempty"`
		Name string `bson:"name"`
		Age  int    `bson:"age"`
	}

	Convey("test insertOmitemptyTag", t, func() {
		u := Person{ID: "12345678"}

		ru := insertOmitemptyTag(&u)
		ruType := reflect.TypeOf(ru).Elem()
		for i := 0; i < ruType.NumField(); i++ {
			field := ruType.Field(i)
			switch field.Name {
			case "ID":
				So(field.Tag.Get("bson"), ShouldEqual, "_id,omitempty")
			case "Name":
				So(field.Tag.Get("bson"), ShouldEqual, "name,omitempty")
			case "Age":
				So(field.Tag.Get("bson"), ShouldEqual, "age,omitempty")
			}
		}
	})
}
