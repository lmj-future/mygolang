package practice

import (
	"fmt"
)

// ReverString 字符串翻转
func ReverString() {
	str := "abcdefg"
	reverStr := []rune(str)
	for i, j := 0, len(reverStr)-1; i < j; i, j = i+1, j-1 {
		reverStr[i], reverStr[j] = reverStr[j], reverStr[i]
	}
	fmt.Println("befor: " + str)
	fmt.Println("after: " + string(reverStr))
}
