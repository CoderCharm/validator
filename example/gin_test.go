package example

import (
	"fmt"
	"github.com/CoderCharm/validator"
	"github.com/gin-gonic/gin"
	"testing"
)

type person struct {
	Name   string `json:"name" form:"name" binding:"required" required:"用户名不能为空" minLen:"5:用户名最小长度不能小于5" maxLen:"10:用户名最大长度不能超过10" `
	Age    int64  `json:"age" form:"age" binding:"required" required:"年龄不能为空" gte:"18:年龄应当大于18岁" lte:"100:年龄应该小于等于100岁"`
	Other  string `json:"other" form:"other" regx:"^hello:\\d{3,5}$\\:其他字段正则校验失败"`
	Gender *bool  `json:"gender" form:"gender" binding:"required" required:"性别不能为空"`
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
