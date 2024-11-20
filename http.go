package promgin

import (
	"github.com/ziipin-server/niuhe"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Healthy http healthy check
func Healthy(c *gin.Context) {
	c.String(200, "success")
}

func Metrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// UsePrometheus middleware for gin
func UsePrometheus(engin *gin.Engine) {
	engin.Use(Prometheus)
	engin.GET("/healthy", Healthy)
	engin.GET("/metrics", Metrics())
	engin.POST("/healthy", Healthy)
	engin.POST("/metrics", Metrics())
}

// NiuhePrometheus middleware for niuhe
func NiuhePrometheus(niuheSvr *niuhe.Server) {
	niuheSvr.Use(Prometheus)
	niuheSvr.GetGinEngine().GET("/healthy", Healthy)
	niuheSvr.GetGinEngine().GET("/metrics", Metrics())
	niuheSvr.GetGinEngine().POST("/healthy", Healthy)
	niuheSvr.GetGinEngine().POST("/metrics", Metrics())
}
