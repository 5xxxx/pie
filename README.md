![Alt](https://repobeats.axiom.co/api/embed/3925555e0ce56f9d8fdfce801a830d9547718942.svg "Repobeats analytics image")

## pie
pie is a wrapper based on the official [mongo-go-driver](https://github.com/mongodb/mongo-go-driver), all operations are done through chain calls, making it easy to operate MongoDB.

### Installation

```
go get github.com/5xxxx/pie
```

### Connect to database

```go
package main

import (
    "github.com/5xxxx/pie"
    "go.mongodb.org/mongo-driver/mongo/options"
)

client, err := pie.NewClient(cfg.DataBase, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
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
```