package practice

import (
	"fmt"
	"sync"
	"time"
)

// Map Map
type Map struct {
	c   map[string]*entry
	rmx *sync.RWMutex
}
type entry struct {
	ch      chan struct{}
	value   interface{}
	isExist bool
}

// Out Out
func (m *Map) Out(key string, val interface{}) {
	m.rmx.Lock()
	defer m.rmx.Unlock()
	if e, ok := m.c[key]; ok {
		e.value = val
		e.isExist = true
		close(e.ch)
	} else {
		e = &entry{ch: make(chan struct{}), isExist: true, value: val}
		m.c[key] = e
		close(e.ch)
	}
}

// Rd Rd
func (m *Map) Rd(key string, timeout time.Duration) interface{} {
	m.rmx.Lock()
	if e, ok := m.c[key]; ok && e.isExist {
		m.rmx.Unlock()
		return e.value
	} else if !ok {
		e = &entry{ch: make(chan struct{}), isExist: false}
		m.c[key] = e
		m.rmx.Unlock()
		fmt.Println("协程阻塞 -> ", key)
		select {
		case <-e.ch:
			return e.value
		case <-time.After(timeout):
			fmt.Println("协程超时 -> ", key)
			return nil
		}
	} else {
		m.rmx.Unlock()
		fmt.Println("协程阻塞 -> ", key)
		select {
		case <-e.ch:
			return e.value
		case <-time.After(timeout):
			fmt.Println("协程超时 -> ", key)
			return nil
		}
	}
}

// Practice4 Practice4
func Practice4() {
	mapTest := Map{
		c:   make(map[string]*entry),
		rmx: &sync.RWMutex{},
	}
	go fmt.Println(mapTest.Rd("LMJ", 3*time.Second))
	mapTest.Out("LMJ", "hello")
	go fmt.Println(mapTest.Rd("LMJ", 3*time.Second))
	go fmt.Println(mapTest.Rd("LMJ", 3*time.Second))
	go fmt.Println(mapTest.Rd("LMJ", 3*time.Second))
}
