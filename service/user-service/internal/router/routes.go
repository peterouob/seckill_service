package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/pkg/middleware"
	"github.com/peterouob/seckill_service/service/user-service/internal/controller"
)

func InitRouter(r *gin.Engine, user *controller.UserController) {
	r.Use(middleware.Cors())
	r.POST("login", user.Login)
	u := r.RouterGroup
	{
		u.Use()
		u.POST("register", user.Register)
	}
}
