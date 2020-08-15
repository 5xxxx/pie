/*
 *
 * session_test.go
 * tugrik
 *
 * Created by lintao on 2020/8/15 12:28 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package tugrik

import (
	"fmt"
	"testing"
	"time"

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
			tu, _ := NewTugrik()
			s := NewSession(tu)
			if err := s.FilterBy(tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("FilterBy() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println(s.filter)
		})
	}
}
