package main

import (
	"github.com/joho/godotenv"
	"github.com/peterouob/seckill_service/pkg/cache"
	"github.com/peterouob/seckill_service/pkg/config"
	"github.com/peterouob/seckill_service/pkg/database"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/netutil"
	transport "github.com/peterouob/seckill_service/pkg/transport/grpc"
	"github.com/peterouob/seckill_service/service/user-service/cmd/module"
	"go.uber.org/fx"
)

func main() {
	_ = godotenv.Load()
	app := fx.New(
		fx.Provide(func() *config.Config {
			return config.ProvideConfig(
				netutil.FormatIP("50051"),
				"8083",
				"user-svc",
			)
		}),
		database.MySQLModule,
		cache.RedisModule,
		etcdregister.EtcdServerModule,
		transport.GrpcServerModule,

		module.UserModule,
	)
	app.Run()
}
