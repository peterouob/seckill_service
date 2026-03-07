package main

import (
	"github.com/peterouob/seckill_service/pkg/config"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/netutil"
	transport "github.com/peterouob/seckill_service/pkg/transport/http"
	"github.com/peterouob/seckill_service/service/user-service/gateway/moudle"
	"github.com/peterouob/seckill_service/service/user-service/internal/controller"
	"github.com/peterouob/seckill_service/service/user-service/internal/infrastructure/usergrpc/client"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(func() *config.Config {
			return &config.Config{
				HttpAddr: netutil.FormatIP("8083"),
				EtcdConfig: &config.EtcdConfig{
					Endpoints:   []string{"127.0.0.1:2379"},
					ServiceName: "user-svc",
				},
			}
		}),
		fx.Provide(controller.NewUserController),
		transport.HTTPServerModule,
		etcdregister.EtcdClientModule,
		client.UserClientModule,
		moudle.UserGatewayModule,
	)
	app.Run()
}
