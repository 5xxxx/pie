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
	t, err := pie.NewClient(cfg.DataBase, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}

	err = t.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	var user User
	err = t.filter("nickName", "frank").FindOne(&user)
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
