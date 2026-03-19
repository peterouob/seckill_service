package main

import (
	"github.com/joho/godotenv"
	"github.com/peterouob/seckill_service/pkg/cache"
	"github.com/peterouob/seckill_service/pkg/config"
	"github.com/peterouob/seckill_service/pkg/database"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/kafka"
	"github.com/peterouob/seckill_service/pkg/netutil"
	"github.com/peterouob/seckill_service/service/seckill-service/cmd/module"
	"go.uber.org/fx"
)

func main() {
	_ = godotenv.Load()
	app := fx.New(
		fx.Provide(func() *config.Config {
			return config.ProvideConfig(
				netutil.FormatIP("50052"),
				"8084",
				"seckill-svc",
			)
		}),
		database.MySQLModule,
		cache.RedisModule,
		etcdregister.EtcdServerModule,
		kafka.KafkaProducerModule,
		module.SeckillModule,
	)
	app.Run()
}
