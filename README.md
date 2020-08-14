## tugrik

tugrik 是对[mongo-go-driver](https://github.com/mongodb/mongo-go-driver) 二次开发的操作库

### 安装

```
go get github.com/NSObjects/tugrik
```

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
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

var user User
user.Username = "小明"
user.MobileNumber = "138xxxx"

err := t.Insert(&user)
if err != nil {
    fmt.Println(err)
}
```

### FindOne

```
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

var u User
_, err := t.Filter("username", "小明").Get(&u)
if err != nil {
    panic(err)
}
```

### FindAll

```
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

var user []User
t.Gt("age",10).Skip(10).Limit(100).FindAll(&user)
if err != nil {
    panic(err)
}
```

### UpdateOne

```
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

u := new(User)
u.Username = "hahhahahah"
err := t.Filter("mobileNumber", "138xxxxx").Id(u.Id).Update(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteOne

```
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

u := new(User)
u.MobileNumber = "138xxxxx"
err := t.Id(u.Id).Delete(u)
if err != nil {
    fmt.Println(err)
}
```

### DeleteMany

```
t := tugrik.NewTugrik(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

u := new(User)
u.MobileNumber = "138xxxxx"
err := t.DeleteMany(u)
if err != nil {
    fmt.Println(err)
}
```

### Aggregate
```
    t, err := tugrik.NewTugrik()
	t.SetURI("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	if err = t.Connect(context.Background()); err != nil {
		panic(err)
	}

	t.SetDatabase("jishimao_local")
	var user []User

	err = t.Aggregate().
		Match(tugrik.DefaultCondition().
			Eq("nick_name", "黄晶晶").
			Eq("mobile_number", "c5b013cb2e102e0e743f117220b2acd1")).All(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
```