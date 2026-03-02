package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/services/user-service/internal/controller"
	"github.com/peterouob/seckill_service/utils"
)

func InitRouter(r *gin.Engine, user *controller.UserController) {
	r.Use(utils.Cors())
	r.POST("login", user.Login)
	u := r.RouterGroup
	{
		u.Use()
		u.POST("register", user.Register)
	}
}
