package module

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/app/user-service/internal/controller"
	"go.uber.org/fx"
)

var UserGatewayModule = fx.Module("user-gateway",
	fx.Provide(controller.NewUserController),
	fx.Invoke(func(r *gin.Engine, c *controller.UserController) {
		r.POST("/login", c.Login)
		r.POST("/register", c.Register)
	}),
)
