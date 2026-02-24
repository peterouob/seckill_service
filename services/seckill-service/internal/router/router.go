package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/services/seckill-service/internal/controller"
)

func InitRouter(ctl *controller.SeckillController) *gin.Engine {
	r := gin.Default()
	r.POST("/buy", ctl.Buy)
	return r
}
