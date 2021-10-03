package main

//还是用框架搭服务简单

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test ok")
	})
	r.GET("/api/v1", func(c *gin.Context) {
		c.String(200, "api v1 ok")
	})

	r.Run(":8081")
}
