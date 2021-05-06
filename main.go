package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/h3c/mygolang/metric"
	"github.com/h3c/mygolang/udp"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("main err : ", err)
		}
	}()
	go func() {
		http.ListenAndServe("0.0.0.0:80", nil)
	}()
	fmt.Printf("********************************Hello LMJ********************************\n")

	go udp.UDP()

	go metric.PrometheusStart("1", "2", "3")

	time.Sleep(time.Hour)
}
