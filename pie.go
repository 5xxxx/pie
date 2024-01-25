/*

Example:

package main

import (
	"context"
	"fmt"
	"time"


	"github.com/5xxxx/pie"
)

func main() {
	t, err := pie.NewClient("demo")
	t.SetURI("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	err = t.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	var user User
	err = t.filter("nickName", "淳朴的润土").FindOne(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
*/

package pie

func NewCondition() Condition {
	return DefaultCondition()
}
