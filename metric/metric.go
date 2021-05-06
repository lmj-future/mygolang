package metric

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var UDPReceiveCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "UDP_send_aa",
		Help: "total UDP count",
	},
	// []string{"description", "address", "data"},
	[]string{},
)

// 初始化 web_reqeust_total， counter类型指标， 表示接收http请求总次数
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_reqeust_total",
		Help: "Number of hello requests in total",
	},
	// 设置两个标签 请求方法和 路径 对请求总次数在两个
	[]string{"method", "endpoint"},
)

// web_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "web_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
)

func init() {
	// 注册监控指标
	prometheus.MustRegister(UDPReceiveCount)
	prometheus.MustRegister(WebRequestTotal)
	prometheus.MustRegister(WebRequestDuration)
}

// 包装 handler function,不侵入业务逻辑
func Monitor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		duration := time.Since(start)
		// counter类型 metric的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 meric的记录方式
		WebRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
	}
}

// PrometheusStart prometheus Start
func PrometheusStart(arg ...string) {
	fmt.Println(arg)
	// expose prometheus metrics接口
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/query", Monitor(Query))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// query
func Query(w http.ResponseWriter, r *http.Request) {
	//模拟业务查询耗时0~1s
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	_, _ = io.WriteString(w, "some results")
}
