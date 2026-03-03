package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/app/seckill-service/internal/controller"
	"github.com/peterouob/seckill_service/pkg/middleware"
)

func InitRouter(ctl *controller.SeckillController) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthByJWT(), middleware.Cors())
	r.POST("/buy", ctl.Buy)
	return r
}
