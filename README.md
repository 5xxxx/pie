## pie

pie 是对[mongo-go-driver](https://github.com/mongodb/mongo-go-driver) 二次开发的操作库

### 安装

```
go get github.com/NSObjects/pie
```

### 连接到数据库

```
driver := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

or

t.SetURI("mongodb://127.0.0.1:27017")

if err != nil {
    fmt.Println(err)
}
driver.SetDatabase("xxx")

err = driver.Connect(context.Background())
if err != nil {
    panic(err)
}

var user User
err = driver.Filter("nickName", "淳朴的润土").FindOne(&user)
if err != nil {
    panic(err)
}

fmt.Println(user)
```

### 数据定义

```
type User struct {
	ID             string              `bson:"_id"`
	Username       string              `bson:"username"`
	MobileNumber   string              `bson:"mobileNumber"`
}
```

### Crete

```
driver := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

var user User
user.Username = "小明"
user.MobileNumber = "138xxxx"

err := driver.Insert(&user)
if err != nil {
    fmt.Println(err)
}
```

### FindOne

```
driver := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

var u User
_, err := driver.Filter("username", "小明").Get(&u)
if err != nil {
    panic(err)
}

or

var user User
user.Id, _ = primitive.ObjectIDFromHex("5f0ace734e2d4100013d8797")
err = driver.FilterBy(user).FindOne(context.Background(), &user)
if err != nil {
    panic(err)
}

```

### FindAll

```
driver := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

var user []User
driver.Gt("age",10).Skip(10).Limit(100).FindAll(&user)
if err != nil {
    panic(err)
}

or

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

```
driver := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

u := new(User)
u.Username = "hahhahahah"
err := driver.Filter("mobileNumber", "138xxxxx").Id(u.Id).Update(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteOne

```
t := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

u := new(User)
u.MobileNumber = "138xxxxx"
err := driver.Id(u.Id).Delete(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteMany

```
driver := pie.NewDriver(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

u := new(User)
u.MobileNumber = "138xxxxx"
err := driver.DeleteMany(u)
if err != nil {
    fmt.Println(err)
}
```

### Aggregate
```
driver, err := pie.NewDriver()
driver.SetURI("mongodb://127.0.0.1:27017")
if err != nil {
    panic(err)
}
if err = driver.Connect(context.Background()); err != nil {
    panic(err)
}

driver.SetDatabase("jishimao_local")
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