package util

import (
	"strings"
)

// 下划线式
// date_at => date_at
// dateAt => date_at
// DateAt => date_at
func LineCase(str string) string {
	if len(str) <= 0 {
		return ""
	}
	var res strings.Builder
	// 首字母特殊处理
	if str[0] >= 'A' && str[0] <= 'Z' {
		res.WriteByte(str[0] + 32)
	} else {
		res.WriteByte(str[0])
	}
	str = str[1:]
	for k, val := range str {
		if val >= 'A' && val <= 'Z' {
			res.WriteByte('_')
			res.WriteByte(str[k] + 32)
		} else {
			res.WriteByte(str[k])
		}
	}
	return res.String()
}

// 小写驼峰式
// date_at => dateAt
// dateAt => dateAt
// DateAt => dateAt
func LowerCamelCase(str string) string {
	strs := strings.Split(str, "_")
	if len(strs) > 1 {
		var res strings.Builder
		// 第一个单词首字母小写
		res.WriteString(LowerWord(strs[0]))
		strs = strs[1:]
		for _, val := range strs {
			res.WriteString(UpperWord(val))
		}
		return res.String()
	} else {
		res := LowerWord(strs[0])
		return res
	}
}

// 大写驼峰式
// date_at => DateAt
// dateAt => DateAt
// DateAt => DateAt
func UpperCamelCase(str string) string {
	strs := strings.Split(str, "_")
	if len(strs) > 1 {
		var res strings.Builder
		for _, val := range strs {
			res.WriteString(UpperWord(val))
		}
		return res.String()
	} else {
		res := UpperWord(strs[0])
		return res
	}
}

// 单词首字母大写
func UpperWord(str string) string {
	str = strings.ToUpper(str[0:1]) + str[1:]
	return str
}

// 单词首字母小写
func LowerWord(str string) string {
	str = strings.ToLower(str[0:1]) + str[1:]
	return str
}