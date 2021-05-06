package practice

import (
	"fmt"
	"time"
)

func proc() {
	panic("ok")
}

// Practice6 Practice6
func Practice6() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				go func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(err)
						}
					}()
					proc()
				}()
			}
		}
	}()
	select {}
}
