package udp

import (
	"fmt"
	"net"

	"github.com/h3c/mygolang/metric"
	"github.com/prometheus/client_golang/prometheus"
)

func UDP() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:3501")
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 65507)
	for {
		i, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		data := make([]byte, i)
		copy(data, buf[:i])
		fmt.Println(addr)
		fmt.Println(data)
		// metric.UDPReceiveCount.With(prometheus.Labels{"description": "this is a test", "address": addr.IP.String(), "data": string(data)}).Inc()
		// metric.UDPReceiveCount.With(prometheus.Labels{"description": "this is a test", "address": "", "data": string(data)}).Inc()
		// metric.UDPReceiveCount.With(prometheus.Labels{"description": "this is a test", "address": addr.IP.String(), "data": ""}).Inc()
		metric.UDPReceiveCount.With(prometheus.Labels{}).Inc()
	}
}
