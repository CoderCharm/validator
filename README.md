# Validator

功能如名字一样，这是一个`golang`结构体字段验证的库，算是一个非常重复的轮子。和 https://github.com/go-playground/validator 库功能类似,
但是功能远远没有那个库强大。

> 那么为什么要写这个库了?

- 1 我发现上面 [validator](https://github.com/go-playground/validator) 有些地方不够友好，比如自定义错误消息提示，见
  [issues](https://github.com/go-playground/validator/issues/559) 。
- 2 自己也是一个初学者，顺带练习一下反射，花点业余时间倒腾下。

## 声明

这个库只是我业余时间写的，不保证代码质量，生产环境使用，可自行查看代码更改，主要核心校验参数的代码，也就几十行左右。

## 如何使用

额，因为代码不多，所以可以直接复制`validator.go`文件到你到项目中。当然也可以选择安装 `go get github.com/CoderCharm/validator`。

我写这个主要就是解决自定义错误提示。所以我这样规定了结构体`tag`的用法:

以冒号`:`分割，第一个为匹配规则，第二个为错误提示信息(没有默认值，必填)。如果使用的正则里面包含`:`可以把分割的冒号用`\\:`表示。具体也有例子。

```
字段名 类型 `验证类型1:"格式:错误提示信息1" 验证类型2:"格式:错误提示信息2"`


// 主要定义了以下 验证的tag字段
const (
	required = "required"  // 表示必填
	regx     = "regx"      // 正则校验
	maxLen   = "maxLen"    // 字符串最大长度
	minLen   = "minLen"    // 字符串最小长度
	lt       = "lt"        // 小于
	lte      = "lte"       // 小于等于
	gt       = "gt"        // 大于
	gte      = "gte"       // 大于等于
)
```

如下: 作用是验证`Age`这个字段必须小于35, 否则返回 `errors.New("年龄必须小于35岁")` 错误。
```
Age int64 `json:"age" lt:"35:年龄必须小于35岁"`
```

见`vaildator_test`文件示例

```go
type person struct {
	Name   string `json:"name" required:"用户名不能为空" minLen:"5:用户名最小长度不能小于5" maxLen:"10:用户名最大长度不能超过10" `
	Age    int64  `json:"age" required:"年龄不能为空" gte:"18:年龄应当大于18岁" lte:"100:年龄应该小于等于100岁"`
	Other  string `json:"other" regx:"^hello:\\d{3,5}$\\:其他字段正则校验失败"`
	Gender *bool   `json:"gender" required:"性别不能为空"`
    Email  string `json:"email" regx:"^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$:邮箱格式错误"`
}

// 由于普通bool类型 无法判断 false 还是未赋值，所以用*bool 见: https://stackoverflow.com/questions/43351216/check-if-boolean-value-is-set-in-go
var gender = false // 当然现在一般都用 整数 表示性别 我这里只是用作演示

func Test_Verify_Name(t *testing.T) {

	// 判断所有字段 都有数据时是否校验
	var user1 = person{Name: "nick123", Age: 24, Other: "hello:123", Gender: &gender}

	if err := Verify(user1); err != nil {
		panic(err.Error())
	}

	user1.Name = "nick"
	if err := Verify(user1); err != nil {
		assert.Equal(t, err.Error(), "用户名最小长度不能小于5")
	}

	user1.Name = "nick1234567890"
	if err := Verify(user1); err != nil {
		assert.Equal(t, err.Error(), "用户名最大长度不能超过10")
	}
}
```

## 如何在 golang Web框架中使用

因为我这库，只对结构体`struct`校验，所以无法校验非结构体之类的接收参数方式。

<details>
<summary>gin 和 iris 框架示例</summary>

- gin 示例

```go
type person struct {
	Name   string `json:"name" form:"name" binding:"required" required:"用户名不能为空" minLen:"5:用户名最小长度不能小于5" maxLen:"10:用户名最大长度不能超过10" `
	Age    int64  `json:"age" form:"age" binding:"required" required:"年龄不能为空" gte:"18:年龄应当大于18岁" lte:"100:年龄应该小于等于100岁"`
	Other  string `json:"other" form:"other" regx:"^hello:\\d{3,5}$\\:其他字段正则校验失败"`
	Gender *bool   `json:"gender" form:"gender" binding:"required" required:"性别不能为空"`
}

func Test_Gin_Params(t *testing.T) {

	r := gin.Default()

	// url参数校验
	r.GET("/", func(ctx *gin.Context) {

		p := person{}

		if err := ctx.ShouldBindQuery(&p); err != nil {
			ctx.JSON(401, gin.H{
				"code": 4001,
				"msg":  "参数缺失",
				"tip":  err.Error(),
			})
			return
		}

		if errs := validator.Verify(p); errs != nil {
			ctx.JSON(402, gin.H{
				"code": 4002,
				"msg":  "参数规则错误",
				"tip":  errs.Error(),
			})
			return
		}

		ctx.JSON(200,
			gin.H{
				"code": 200,
				"msg":  "ok",
				"data": p,
			},
		)
	})

	// json参数校验
	r.POST("/", func(ctx *gin.Context) {
		p := person{}

		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(401, gin.H{
				"code": 4001,
				"msg":  "参数缺失",
				"tip":  err.Error(),
			})
			return
		}

		if errs := validator.Verify(p); errs != nil {
			ctx.JSON(402, gin.H{
				"code": 4002,
				"msg":  "参数规则错误",
				"tip":  errs.Error(),
			})
			return
		}

		ctx.JSON(200,
			gin.H{
				"code": 200,
				"msg":  "ok",
				"data": p,
			},
		)
	})

	fmt.Println("http://127.0.0.1:8082")
	r.Run(":8082")
}
```

- iris 示例

```go
type user struct {
	Name   string `json:"name" required:"用户名不能为空" minLen:"5:用户名最小长度不能小于5" maxLen:"10:用户名最大长度不能超过10" `
	Age    int64  `json:"age" required:"年龄不能为空" gte:"18:年龄应当大于18岁" lte:"100:年龄应该小于等于100岁"`
	Other  string `json:"other" regx:"^hello:\\d{3,5}$\\:其他字段正则校验失败"`
	Gender *bool   `json:"gender" required:"性别不能为空"`
}


func Test_Iris_Params(t *testing.T) {
	app := iris.New()

	var u1 user

	// url参数格式校验
	app.Get("/", func(ctx iris.Context) {
		if err := ctx.ReadQuery(&u1); err != nil{
			ctx.JSON(iris.Map{
				"code": 4001,
				"msg": "参数错误",
				"tip": err.Error(),
			})
			return
		}

		if errs := validator.Verify(u1); errs != nil {
			ctx.JSON(iris.Map{
				"code": 4002,
				"msg":  "参数规则错误",
				"tip":  errs.Error(),
			})
			return
		}

		ctx.JSON(iris.Map{
			"code": 200,
			"msg":  "ok",
			"data": u1,
		})

	})
	// json参数校验
	app.Post("/", func(ctx iris.Context) {
		if err := ctx.ReadJSON(&u1); err != nil{
			ctx.JSON(iris.Map{
				"code": 4001,
				"msg": "参数错误",
				"tip": err.Error(),
			})
			return
		}

		if errs := validator.Verify(u1); errs != nil {
			ctx.JSON(iris.Map{
				"code": 4002,
				"msg":  "参数规则错误",
				"tip":  errs.Error(),
			})
			return
		}

		ctx.JSON(iris.Map{
			"code": 200,
			"msg":  "ok",
			"data": u1,
		})

	})
	
	app.Listen(":8081")
}

```

</details>


## 不足之处

由于时间仓促，没有写验证复杂的数据类型的处理，比如结构体嵌套， slice等，还有就是自定义的类型，
比如这个库 [null.xx](https://github.com/guregu/null) ，`gorm`的`sql.NullString`之类的也不支持 。具体实现应该也不难。



