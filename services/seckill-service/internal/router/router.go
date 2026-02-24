package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/services/seckill-service/internal/controller"
	"github.com/peterouob/seckill_service/utils"
)

func InitRouter(ctl *controller.SeckillController) *gin.Engine {
	r := gin.Default()
	r.Use(utils.AuthByJWT(), utils.Cors())
	r.POST("/buy", ctl.Buy)
	return r
}
