## tugrik

tugrik 是对[mongo-go-driver](https://github.com/mongodb/mongo-go-driver) 二次开发的操作库

### 连接到数据库

```
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

or

t.SetURI("mongodb://127.0.0.1:27017")

if err != nil {
    fmt.Println(err)
}
t.SetDatabase("xxx")

err = t.Connect(context.Background())
if err != nil {
    panic(err)
}

var user User
err = t.Filter("nickName", "淳朴的润土").FindOne(&user)
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
var user User
user.Username = "小明"
user.MobileNumber = "138xxxx"

err := tugrik.Insert(&user)
if err != nil {
    fmt.Println(err)
}
```

### FindOne

```
var u User
_, err := tugrik.Filter("username", "小明").Get(&u)
if err != nil {
    panic(err)
}
```

### UpdateOne

```
u := new(User)
u.Username = "hahhahahah"
err := tugrik.Filter("mobileNumber", "138xxxxx").Id(u.Id).Update(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteOne

```
u := new(User)
u.MobileNumber = "138xxxxx"
err := tugrik.Id(u.Id).Delete(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteMany

```
u := new(User)
u.MobileNumber = "138xxxxx"
err := tugrik.DeleteMany(u)
if err != nil {
    fmt.Println(err)
}
```