package em

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// Generate random strings
// 生成随机字符串
func GenerateRandomString(l int) string {
	var code = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/~!@#$%^&*()_="

	data := make([]byte, l)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < l; i++ {
		idx := rand.Intn(len(code))
		data[i] = code[idx]
	}
	return string(data)
}

func MustConvertJson(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}

// Note: json to map int format will be converted to float
// 注意：json转map int格式会转换为float
func StructToMap(v interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m = make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Check if the slice contains elements
// 检查切片是否包含元素
func CheckIfSliceContainsInt(search int, ints []int) bool {
	for _, v := range ints {
		if v == search {
			return true
		}
	}

	return false
}



func GetUrlPath(urlStr string, trim bool) string {
	u, _ :=url.Parse(urlStr)
	if trim {
		return strings.TrimLeft(u.Path, "/")
	} else {
		return u.Path
	}
}

// Deprecated: Use ChangeTypeV2
func ChangeType(in interface{}, out interface{}) (error) {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		return err
	}
	return nil
}

func ChangeTypeV2(in interface{}, out interface{}) (error) {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}
	return nil
}

// Generate errors with both custom messages and error messages
// 生成同时带有自定义信息和错误信息的错误
func GenerateErrorWithMessage(msg string, err error) error {
	return errors.New(msg + err.Error())
}

// Debug errors and custom errors are used as parameters at the same time, and appropriate errors are output according to environment variables.
// 把Debug错误和自定义错误同时作为参数，根据环境变量输出适合的错误。
func GetErrorByIfDebug(debug bool, err error, msg string) error {
	if debug {
		return err
	}
	return errors.New(msg)
}

// Get random number
// 获取随机数
// Example: len=3 return: 0/1/2
func GetRandomNumberByLength(len int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len)
}

// Request : "abc"  Return: "abc"
// Request : ""  Return: "[]"
func IfEmptyChangeJsonFormat(str string) string {
	if len(str) == 0 {
		return "[]"
	}
	return str
}

