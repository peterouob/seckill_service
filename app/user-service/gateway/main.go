package main

import (
	"github.com/peterouob/seckill_service/app/user-service/gateway/module"
	"github.com/peterouob/seckill_service/app/user-service/internal/infrastructure/usergrpc"
	"github.com/peterouob/seckill_service/pkg/injection"
	"github.com/peterouob/seckill_service/pkg/netutil"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(func() *injection.Config {
			return &injection.Config{
				HttpAddr: netutil.FormatIP("8083"),
				EtcdConfig: &injection.EtcdConfig{
					Endpoints:   []string{"127.0.0.1:2379"},
					ServiceName: "user-svc",
				},
			}
		}),

		injection.HTTPServerModule,
		injection.EtcdClientModule,
		usergrpc.UserClientModule,
		module.UserGatewayModule,
	)
	app.Run()
}
