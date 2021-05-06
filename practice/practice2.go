package practice

import (
	"strings"
)

// IsUniqueString 判断字符串中字符是否全都不同
func IsUniqueString() bool {
	str := "abcdefg"
	if len(str) > 3000 {
		return false
	}
	for v := range str {
		if v > 127 {
			return false
		}
		if strings.Count(str, string(v)) > 1 {
			return false
		}
	}
	return true
}
