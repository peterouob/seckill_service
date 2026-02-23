package router

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})
	r.Group("user")
	{
		r.POST("login")
		r.POST("register")
	}
	r.Group("seckill")
	{
		r.POST("order")
		//r.POST("purchase") 模擬
	}
}
