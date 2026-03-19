package main

import (
	"github.com/peterouob/seckill_service/pkg/config"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/netutil"
	transport "github.com/peterouob/seckill_service/pkg/transport/http"
	"github.com/peterouob/seckill_service/service/seckill-service/gateway/module"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/controller"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/infrastructure/seckillgrpc/client"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(func() *config.Config {
			return &config.Config{
				HttpAddr: netutil.FormatIP("8082"),
				EtcdConfig: &config.EtcdConfig{
					Endpoints:   []string{"127.0.0.1:2379"},
					ServiceName: "seckill-svc",
				},
			}
		}),

		fx.Provide(controller.NewSeckillController),
		transport.HTTPServerModule,
		etcdregister.EtcdClientModule,
		client.SeckillClientModule,
		module.SeckillClientModule,
	)
	app.Run()
}
