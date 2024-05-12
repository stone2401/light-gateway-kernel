package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})
	// ipv6
	app.Run(":8080")
}
