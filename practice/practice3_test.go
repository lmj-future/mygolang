package practice_test

import (
	"fmt"
	"testing"
)

// ReverString 字符串翻转
func TestReverString(t *testing.T) {
	str := "abcdefg"
	reverStr := []rune(str)
	for i, j := 0, len(reverStr)-1; i < j; i, j = i+1, j-1 {
		reverStr[i], reverStr[j] = reverStr[j], reverStr[i]
	}
	fmt.Println("befor: " + str)
	fmt.Println("after: " + string(reverStr))
}
