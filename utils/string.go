package em_utils

import (
	"encoding/json"
	"math/rand"
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
