package validator

import (
	"github.com/go-playground/assert/v2"
	"strings"
	"testing"
)

type person struct {
	Name   string `json:"name" required:"用户名不能为空" minLen:"5:用户名最小长度不能小于5" maxLen:"10:用户名最大长度不能超过10" `
	Age    int64  `json:"age" required:"年龄不能为空" gte:"18:年龄应当大于18岁" lte:"100:年龄应该小于等于100岁"`
	Other  string `json:"other" regx:"^hello:\\d{3,5}$\\:其他字段正则校验失败"`
	Gender *bool  `json:"gender" required:"性别不能为空"`
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

func Test_Verify_Age(t *testing.T) {

	// 判断所有字段 都有数据时是否校验
	var user1 = person{Name: "nick123", Age: 24, Other: "hello:123", Gender: &gender}

	user1.Age = 15
	if err := Verify(user1); err != nil {
		assert.Equal(t, err.Error(), "年龄应当大于18岁")
	}

	user1.Age = 101
	if err := Verify(user1); err != nil {
		assert.Equal(t, err.Error(), "年龄应该小于等于100岁")
	}
}

func Test_Verify_Other(t *testing.T) {

	// 判断所有字段 都有数据时是否校验
	var user1 = person{Name: "nick123", Age: 24, Other: "hello:123", Gender: &gender}

	user1.Other = "hello:123456"
	if err := Verify(user1); err != nil {
		assert.Equal(t, err.Error(), "其他字段正则校验失败")
	}
}

func Test_Verify_Default_Field(t *testing.T) {
	// Other 字段可以为空测试
	user2 := person{Name: "nick123", Age: 23, Gender: &gender}

	if err := Verify(user2); err != nil {
		panic(err.Error())
	}
}

func Test_sp(t *testing.T) {
	//z := "2021-11-10 11:10:10"
	z := "2021-11-10"
	arr := strings.Split(z, " ")
	t.Log(arr[0])

}
