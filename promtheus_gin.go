package promgin

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var ginLables = []string{"method", "uri", "code"}

var ginHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "gin_duration_millseconds",
		Help:    "gin duration millseconds distribution",
		Buckets: []float64{10, 50, 100, 200, 480, 1000, 2000, 5000},
	},
	ginLables,
)

var ginRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "gin_http_request_counter",
	Help: "gin http request counter",
}, ginLables)

var ginRespAvgPerMintue = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "gin_resp_time_avg_per_min",
	Help: "gin resp time avg per min",
}, ginLables)

var ginRespMaxPerMintue = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "gin_resp_time_max_per_min",
	Help: "gin_resp_time_max_per_min",
}, ginLables)

func resetMetrics() {
	ginHistogram.Reset()
	ginRequestCounter.Reset()
	ginRespAvgPerMintue.Reset()
	ginRespMaxPerMintue.Reset()
}

// reset store every minute
func cleaner() {
	tick1m := time.NewTicker(time.Minute)
	defer tick1m.Stop()
	tick1h := time.NewTicker(time.Hour)
	defer tick1h.Stop()
	for {
		select {
		case <-tick1m.C:
			resetCache()
		case <-tick1h.C:
			resetMetrics()
		}
	}
}

func init() {
	prometheus.MustRegister(ginHistogram)
	prometheus.MustRegister(ginRequestCounter)
	prometheus.MustRegister(ginRespAvgPerMintue)
	prometheus.MustRegister(ginRespMaxPerMintue)
	go cleaner()
}
