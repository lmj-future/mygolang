package lock

import (
	"fmt"
	"sync"
)

var m = sync.Mutex{}

// Lock 互斥锁上锁
func Lock() {
	m.Lock()
	fmt.Println("it is locked!")
}

// Unlock 互斥锁上锁
func Unlock() {
	m.Unlock()
	fmt.Println("it is unlocked!")
}
