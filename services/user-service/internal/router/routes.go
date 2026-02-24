package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/services/user-service/internal/controller"
)

func InitRouter(user *controller.UserController) *gin.Engine {
	r := gin.Default()
	r.Group("user")
	{
		r.POST("login", user.Login)
		r.POST("register")
	}
	return r
}
