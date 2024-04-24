package util

import (
	"strconv"
	"strings"
)

// 生成重复字符串
func GenerateDuplicates(count int, src, suffix string) string {
	repeat := strings.Repeat(src, count)
	return strings.TrimSuffix(repeat, suffix)
}

func FirstLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(s[:1]) + s[1:]
	}
	return s
}

func FirstUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

// 驼峰命名
func HumpNaming(s string) string {
	split := strings.Split(s, "_")
	for i := 0; i < len(split); i++ {
		str := split[i]
		if len(str) > 0 {
			split[i] = strings.ToUpper(str[:1]) + str[1:]
		}
	}
	return strings.Join(split, "")
}

// 下划线命名
func SlideNaming(s string) string {
	var chars []byte
	for i := 0; i < len(s); i++ {
		if 'A' <= s[i] && s[i] <= 'Z' {
			if i != 0 {
				chars = append(chars, '_')
			}
			chars = append(chars, s[i]+32)
		} else {
			chars = append(chars, s[i])
		}
	}
	return string(chars)
}

func StrToIntArray(str string, split string) (array []int, err error) {
	if len(str) <= 0 {
		return
	}
	list := strings.Split(str, split)
	array = make([]int, len(list))
	for i := 0; i < len(list); i++ {
		atoi, err := strconv.Atoi(list[i])
		if err != nil {
			return nil, err
		}
		array[i] = atoi
	}
	return
}
