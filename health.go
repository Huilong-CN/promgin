package promgin

import "github.com/gin-gonic/gin"

//Healthy http healthy check
func Healthy(c *gin.Context) {
	c.String(200, "success")
}
