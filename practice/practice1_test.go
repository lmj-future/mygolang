package practice_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// PrintLetterNumber 交替打印数字和字母
func TestPrintLetterNumber(t *testing.T) {
	letter := make(chan bool)
	number := make(chan bool)
	wait := sync.WaitGroup{}
	num := 0
	str := "abcdefghijklmnopqrstuvwxyz"

	// 该协程负责打印数字
	go func() {
		for {
			select {
			case <-number:
				fmt.Printf("%d", num)
				num++
				time.Sleep(time.Second)
				letter <- true
			}
		}
	}()
	// wait用来等待打印结束退出循环
	wait.Add(1)
	// 该协程负责打印字母
	go func(wait *sync.WaitGroup) {
		for {
			select {
			case <-letter:
				if num == len(str)+1 {
					wait.Done()
					return
				}
				fmt.Printf("%s", str[num-1:num])
				time.Sleep(time.Second)
				number <- true
			}
		}
	}(&wait)
	number <- true
	wait.Wait()
}
