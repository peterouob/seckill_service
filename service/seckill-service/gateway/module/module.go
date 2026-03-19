package module

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/controller"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/router"
	"go.uber.org/fx"
)

func provideHTTPServer(r *gin.Engine, c *controller.SeckillController) {
	router.InitRouter(r, c)
}

var SeckillClientModule = fx.Module("seckill-gateway",
	fx.Invoke(provideHTTPServer),
)
