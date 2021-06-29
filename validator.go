package validator

import (
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 常量定义 结构体需要 验证的tag字段
const (
	required = "required" // 表示必填
	regx     = "regx"     // 正则校验
	maxLen   = "maxLen"   // 字符串最大长度
	minLen   = "minLen"   // 字符串最小长度
	lt       = "lt"       // 小于
	lte      = "lte"      // 小于等于
	gt       = "gt"       // 大于
	gte      = "gte"      // 大于等于
)

func Verify(target interface{}) error {

	validatorType := []string{maxLen, minLen, lt, lte, gt, gte}

	val := reflect.ValueOf(target)
	typ := reflect.TypeOf(target)

	for i := 0; i < val.NumField(); i++ {
		tag := typ.Field(i)
		tagValue := val.Field(i)

		var needVal bool // 是否比填字段
		requireMsg, ok := tag.Tag.Lookup(required)
		if ok {
			needVal = true
			// 空判断
			if isEmpty(tagValue) {
				return errors.New(requireMsg)
			}
		}

		regRes, ok := tag.Tag.Lookup(regx)
		// 有regx tag 并且 值不为空时 校验
		if ok && !isEmpty(tagValue) {
			// 校验正则
			regTagArr, err := SplitTag(regRes)
			if err != nil {
				return err
			}

			if ok = isRegexMatch(tagValue.String(), regTagArr[0]); !ok {
				return errors.New(regTagArr[1])
			}
		}

		for _, v := range validatorType {
			tagRes, ok2 := tag.Tag.Lookup(v)
			if !ok2 {
				continue
			}
			tagArr, errs := SplitTag(tagRes)
			if errs != nil {
				return errs
			}

			// 有值 或者 tag必填时校验
			if !isEmpty(tagValue) || needVal {
				if errs = compare(tagValue, v, tagArr); errs != nil {
					return errs
				}
			}
		}
	}
	return nil
}

func compare(tagValue reflect.Value, validatorType string, tagArr []string) error {

	switch tagValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:

		VInt, VErr := strconv.ParseInt(tagArr[0], 10, 64)
		if VErr != nil {
			return errors.Errorf("校验字段不合法:%s-%s", tagArr[1], validatorType)
		}
		switch validatorType {
		case "lt":
			if !(tagValue.Int() < VInt) {
				return errors.New(tagArr[1])
			}
		case "lte":
			if !(tagValue.Int() <= VInt) {
				return errors.New(tagArr[1])
			}
		case "gt":
			if !(tagValue.Int() > VInt) {
				return errors.New(tagArr[1])
			}
		case "gte":
			if !(tagValue.Int() >= VInt) {
				return errors.New(tagArr[1])
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(tagArr[0])
		if VErr != nil {
			return errors.Errorf("校验字段不合法:%s-%s", tagArr[1], validatorType)
		}
		switch validatorType {
		case "lt":
			if !(tagValue.Uint() < uint64(VInt)) {
				return errors.New(tagArr[1])
			}
		case "lte":
			if !(tagValue.Uint() <= uint64(VInt)) {
				return errors.New(tagArr[1])
			}
		case "gt":
			if !(tagValue.Uint() > uint64(VInt)) {
				return errors.New(tagArr[1])
			}
		case "gte":
			if !(tagValue.Uint() >= uint64(VInt)) {
				return errors.New(tagArr[1])
			}
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(tagArr[0], 64)
		if VErr != nil {
			return errors.Errorf("校验字段不合法:%s-%s\n", tagArr[1], validatorType)
		}
		switch validatorType {
		case "lt":
			if !(tagValue.Float() < VFloat) {
				return errors.New(tagArr[1])
			}
		case "lte":
			if !(tagValue.Float() <= VFloat) {
				return errors.New(tagArr[1])
			}
		case "gt":
			if !(tagValue.Float() > VFloat) {
				return errors.New(tagArr[1])
			}
		case "gte":
			if !(tagValue.Float() >= VFloat) {
				return errors.New(tagArr[1])
			}
		}
	case reflect.String:

		VInt, VErr := strconv.ParseInt(tagArr[0], 10, 64)
		if VErr != nil {
			return errors.Errorf("校验字段不合法:%s-%s", tagArr[0], validatorType)
		}
		switch validatorType {
		case "minLen":
			if !(len(tagValue.String()) >= int(VInt)) {
				return errors.New(tagArr[1])
			}
		case "maxLen":
			if !(len(tagValue.String()) <= int(VInt)) {
				return errors.New(tagArr[1])
			}
		}

	}

	return nil
}

// 判断是否为空
func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

/**
以:分割字符串 \\:是为了兼容正则中包含:

e.g.:
tag: minLen:"10:最小长度为10"
return [10 最小长度为10]


tag: reg:"^hao:\\d+\\:正则规则"`
return [^hao:\d+ 正则规则]
*/
func SplitTag(tags string) (tagArr []string, err error) {

	tagArr = strings.Split(tags, ":")
	if len(tagArr) != 2 {
		tagArr = strings.Split(tags, "\\:")
		if len(tagArr) != 2 {
			return nil, errors.Errorf("format error:%s", tags)
		}
	}
	return tagArr, nil
}

// 正则校验是否合法
func isRegexMatch(value string, regexStr string) bool {
	match, _ := regexp.MatchString(regexStr, value)
	return match
}
