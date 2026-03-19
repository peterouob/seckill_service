package client

import (
	"github.com/peterouob/seckill_service/api/seckillproto"
	etcdclient "github.com/peterouob/seckill_service/pkg/etcd/client"
	transport "github.com/peterouob/seckill_service/pkg/transport/grpc"
	"go.uber.org/fx"
)

var SeckillClientModule = fx.Module("seckill-client",
	fx.Provide(
		func(lc fx.Lifecycle, hub *etcdclient.ServiceHub) (seckillproto.SeckillServiceClient, error) {
			conn, err := transport.NewConn(lc, hub, "seckill-svc")
			if err != nil {
				return nil, err
			}
			return seckillproto.NewSeckillServiceClient(conn), nil
		},
	),
)
