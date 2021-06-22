package promgin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Healthy http healthy check
func Healthy(c *gin.Context) {
	c.String(200, "success")
}

func Metrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func UsePrometheus(engin *gin.Engine) {
	engin.Use(Prometheus)
	engin.GET("/healthy", Healthy)
	engin.GET("/metrics", Metrics())
	engin.POST("/healthy", Healthy)
	engin.POST("/metrics", Metrics())
}
