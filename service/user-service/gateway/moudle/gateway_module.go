package moudle

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/service/user-service/internal/controller"
	"github.com/peterouob/seckill_service/service/user-service/internal/router"
	"go.uber.org/fx"
)

func provideHTTPServer(r *gin.Engine, c *controller.UserController) {
	router.InitRouter(r, c)
}

var UserGatewayModule = fx.Module("user-gateway",
	fx.Invoke(provideHTTPServer),
)
