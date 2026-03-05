package main

import (
	"github.com/joho/godotenv"
	"github.com/peterouob/seckill_service/app/user-service/cmd/module"
	"github.com/peterouob/seckill_service/pkg/cache"
	"github.com/peterouob/seckill_service/pkg/database"
	"github.com/peterouob/seckill_service/pkg/injection"
	"github.com/peterouob/seckill_service/pkg/netutil"
	"go.uber.org/fx"
)

func main() {
	_ = godotenv.Load()
	app := fx.New(
		fx.Provide(func() *injection.Config {
			return injection.ProvideConfig(
				netutil.FormatIP("50051"),
				"8083",
				"user-svc",
			)
		}),
		database.MySQLModule,
		cache.RedisModule,
		injection.EtcdModule,
		injection.GrpcServerModule,

		module.UserModule,
	)
	app.Run()
}
