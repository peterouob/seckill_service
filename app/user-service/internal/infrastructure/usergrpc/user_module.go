package usergrpc

import (
	"github.com/peterouob/seckill_service/api/userproto"
	etcdclient "github.com/peterouob/seckill_service/pkg/etcd/client"
	"go.uber.org/fx"
)

var UserClientModule = fx.Module("user_grpc_client",
	fx.Provide(
		func(lc fx.Lifecycle, hub *etcdclient.ServiceHub) (userproto.UserServiceClient, error) {
			conn, err := NewConn(lc, hub, "user-svc")
			if err != nil {
				return nil, err
			}
			return userproto.NewUserServiceClient(conn), nil
		},
	),
)
