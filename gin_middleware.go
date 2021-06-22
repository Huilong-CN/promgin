package promgin

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type StatsCache struct {
	*sync.Map
}

func (sc *StatsCache) Get(api string) (*apistatic, bool) {
	val, ok := sc.Load(api)
	if !ok {
		return nil, false
	}
	apistat, ok := val.(*apistatic)
	return apistat, ok
}

var statsCache = &StatsCache{&sync.Map{}}

type apistatic struct {
	uri     string
	count   int64
	totalMs int64
	maxMs   int64
	minMs   int64
	avgMs   int64
}

//Prometheus API调用统计
func Prometheus(c *gin.Context) {
	start := time.Now().UnixNano()
	c.Next()
	uri := c.Request.URL.Path
	ms := int64(time.Now().UnixNano()-start) / int64(time.Millisecond)
	stat, ok := statsCache.Get(uri)
	if !ok {
		stat = &apistatic{
			uri:     uri,
			count:   1,
			totalMs: ms,
			maxMs:   ms,
			minMs:   ms,
			avgMs:   ms,
		}
		statsCache.Store(uri, stat)
	} else {
		atomic.AddInt64(&stat.count, 1)
		// stat.totalMs += ms
		atomic.AddInt64(&stat.totalMs, ms)
		if stat.maxMs < ms {
			// stat.maxMs = ms
			atomic.SwapInt64(&stat.maxMs, ms)
		}
		if stat.minMs > ms {
			// stat.minMs = ms
			atomic.SwapInt64(&stat.minMs, ms)
		}
		newAvg := stat.totalMs / stat.count
		atomic.SwapInt64(&stat.avgMs, newAvg)
	}
	// GinHistogram.With()
	labels := prometheus.Labels{
		"method": c.Request.Method,
		"uri":    uri,
		"code":   c.Request.Response.Status,
	}
	ginHistogram.With(labels).Observe(float64(ms))
	ginRequestCounter.With(labels).Inc()
	ginRespAvgPerMintue.With(labels).Set(float64(stat.avgMs))
	ginRespMaxPerMintue.With(labels).Set(float64(stat.maxMs))
}

func resetCache() {
	statsCache = &StatsCache{&sync.Map{}}
}
