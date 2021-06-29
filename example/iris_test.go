package example

import (
	"github.com/CoderCharm/validator"
	"github.com/kataras/iris/v12"
	"testing"
)

type user struct {
	Name   string `json:"name" required:"用户名不能为空" minLen:"5:用户名最小长度不能小于5" maxLen:"10:用户名最大长度不能超过10" `
	Age    int64  `json:"age" required:"年龄不能为空" gte:"18:年龄应当大于18岁" lte:"100:年龄应该小于等于100岁"`
	Other  string `json:"other" regx:"^hello:\\d{3,5}$\\:其他字段正则校验失败"`
	Gender *bool  `json:"gender" required:"性别不能为空"`
}

func Test_Iris_Params(t *testing.T) {
	app := iris.New()

	var u1 user

	// url参数格式校验
	app.Get("/", func(ctx iris.Context) {
		if err := ctx.ReadQuery(&u1); err != nil {
			ctx.JSON(iris.Map{
				"code": 4001,
				"msg":  "参数错误",
				"tip":  err.Error(),
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
		if err := ctx.ReadJSON(&u1); err != nil {
			ctx.JSON(iris.Map{
				"code": 4001,
				"msg":  "参数错误",
				"tip":  err.Error(),
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
