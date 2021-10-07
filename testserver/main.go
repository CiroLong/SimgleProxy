package main

import "github.com/gin-gonic/gin"

//还是用框架搭服务简单

func main() {
	go func() {
		r := gin.Default()

		r.GET("/test", func(c *gin.Context) {
			c.String(200, "test ok")
		})
		r.GET("/api/v1", func(c *gin.Context) {
			c.String(200, "api v1 ok")
		})
		r.GET("/api/v2", func(c *gin.Context) {
			c.String(200, "= api v2 ok")
		})

		r.GET("/api/test", func(c *gin.Context) {
			c.String(200, "8081 api test ok")
		})

		r.Run(":8081")
	}()

	go func() {
		r := gin.Default()

		r.GET("/test", func(c *gin.Context) {
			c.String(200, "test ok")
		})
		r.GET("/api/v1", func(c *gin.Context) {
			c.String(200, "api v1 ok")
		})
		r.GET("/api/v2/any", func(c *gin.Context) {
			c.String(200, "^~ api v2 ok")
		})
		r.GET("/api/test", func(c *gin.Context) {
			c.String(200, "8082 api test ok")
		})

		r.Run(":8082")
	}()

	go func() {
		r := gin.Default()

		r.GET("/test", func(c *gin.Context) {
			c.String(200, "test ok")
		})
		r.GET("/api/v1", func(c *gin.Context) {
			c.String(200, "api v1 ok")
		})
		r.GET("/api/v2/any", func(c *gin.Context) {
			c.String(200, "^~ api v2 ok")
		})

		r.GET("/api/test", func(c *gin.Context) {
			c.String(200, "8083 api test ok")
		})
		r.Run(":8083")
	}()

	select {}
}
