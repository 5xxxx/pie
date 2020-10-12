## pie
pie 是基于官方[mongo-go-driver](https://github.com/mongodb/mongo-go-driver) 封装，所有操作都是通过链式调用，可以方便的对MongoDB进行操作

### 1.0 目标 todo list
- [ ] Aggregate封装
- [x] CRUD全功能 e.g FindOneAndDelete
- [ ] 测试
- [ ] 文档

### 安装

```
go get github.com/NSObjects/pie
```

### 连接到数据库

```go
driver := pie.NewDriver("demo",options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

or

driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

var user User
err = driver.Filter("nickName", "淳朴的润土").FindOne(&user)
if err != nil {
    panic(err)
}

fmt.Println(user)
```

pie还处于开发阶段，如果有不能满足需求的功能可以调用DataBase()方法获取*mongo.Database进行操作

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}


base := driver.DataBase()
matchStage := bson.D{{"$match", bson.D{{"operationType", "insert"}}}}
opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)

changeStream, err := base.Watch(context.TODO(), mongo.Pipeline{matchStage}, opts)
if err != nil {
    panic(err)
}
for changeStream.Next(context.Background()) {
    fmt.Println(changeStream.Current)
}
```


### 数据定义

```go
type User struct {
	ID             string              `bson:"_id" filter:"_id"`
	Username       string              `bson:"user_name" filter:"user_name"`
	MobileNumber   string              `bson:"mobile_number" filter:"mobile_number"`
}
```

### InsertOne

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

var user User
user.Username = "小明"
user.MobileNumber = "138888888"

err := driver.Insert(&user)
if err != nil {
    fmt.Println(err)
}
```

### InsertMany

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

var users = []User{
    {NickName: "aab"},
    {NickName: "bbb"},
    {NickName: "ccc"},
    {NickName: "dd"},
    {NickName: "eee"},
}

_, err = driver.InsertMany(context.Background(), users)
if err != nil {
    panic(err)
}
```

### FindOne

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

var u User
_, err := driver.Filter("username", "小明").Get(&u)
if err != nil {
    panic(err)
}

or
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

type User struct {
	ID             string              `bson:"_id" filter:"_id"`
	Username       string              `bson:"user_name" filter:"user_name"`
	MobileNumber   string              `bson:"mobile_number" filter:"-"` //ignore filter
}

var user User
user.Id, _ = primitive.ObjectIDFromHex("5f0ace734e2d4100013d8797")
user.Username = "Jack"
user.MobileNumber = "+86 10086"
err = driver.FilterBy(user).FindOne(context.Background(), &user)
if err != nil {
    panic(err)
}

```

### FindAll

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

var user []User
driver.Gt("age",10).Skip(10).Limit(100).FindAll(&user)
if err != nil {
    panic(err)
}

or

driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

type User struct {
	ID             string              `bson:"_id" filter:"_id"`
	Username       string              `bson:"user_name" filter:"user_name"`
	MobileNumber   string              `bson:"mobile_number" filter:"-"` //ignore filter
}

var users []User
var user User
user.Age = 22
err = driver.FilterBy(user).FindAll(context.Background(), &users)
if err != nil {
    panic(err)
}
fmt.Println(users)
```



### UpdateOne

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

u := new(User)
u.Username = "xiao_ming"
err := driver.Filter("mobileNumber", "138888888").Id(u.Id).Update(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteOne

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

u := new(User)
u.MobileNumber = "138888888"
err := driver.Id(u.Id).Delete(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteMany

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

u := new(User)
u.MobileNumber = "138888888"
err := driver.DeleteMany(u)
if err != nil {
    fmt.Println(err)
}
```

### Aggregate
```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

var user []User

err = driver.Aggregate().
Match(pie.DefaultCondition().
Eq("nick_name", "黄晶晶").
Eq("mobile_number", "c5b013cb2e102e0e743f117220b2acd1")).All(&user)
if err != nil {
    panic(err)
}
fmt.Println(user)
```

### CreateIndexes

```go
driver, err := pie.NewDriver("demo")
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

indexes, err := driver.
    AddIndex(bson.M{"nickName": 1}, options.Index().SetBackground(true)).
    AddIndex(bson.M{"birthday": 1}).
    AddIndex(bson.M{"createdAt": 1, "mobileNumber": 1},
    options.Index().SetBackground(true).SetName("create_mobil_index")).
    CreateIndexes(context.Background(), &User{})

if err != nil {
    panic(err)
}
fmt.Println(indexes)
```
