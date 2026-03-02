package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/peterouob/seckill_service/services/user-service/pkg/module"
	"github.com/peterouob/seckill_service/utils/injection"
	"go.uber.org/fx"
)

func main() {
	_ = godotenv.Load()
	app := fx.New(
		fx.StopTimeout(30*time.Second),
		fx.Provide(func() *injection.Config {
			return injection.ProvideConfig(
				"50051",
				"8083",
				"seckill-svc",
			)
		}),
		injection.MySQLModule,
		injection.RedisModule,
		injection.EtcdModule,
		injection.GrpcServerModule,
		injection.HTTPServerModule,

		module.UserModule,
	)
	app.Run()
}
