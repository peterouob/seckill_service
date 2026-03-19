package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/pkg/middleware"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/controller"
)

func InitRouter(r *gin.Engine, ctl *controller.SeckillController) {
	r.Use(middleware.AuthByJWT(), middleware.Cors())
	r.POST("/buy", ctl.Buy)
}
