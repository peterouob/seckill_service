package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/internal/controller"
)

func InitRouter(user *controller.UserController) *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})
	r.Group("user")
	{
		r.POST("login", user.Login)
		r.POST("register")
	}
	r.Group("seckill")
	{
		r.POST("order")
		//r.POST("purchase") 模擬
	}
	return r
}
