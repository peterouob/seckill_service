package injection

import (
	"context"
	"fmt"
	"time"

	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	etcdclient "github.com/peterouob/seckill_service/pkg/etcd/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
)

var etcdClientHub = fx.Provide(
	func(cfg *Config) *etcdclient.ServiceHub {
		return etcdclient.GetService(cfg.EtcdConfig.Endpoints)
	},
)

var EtcdClientModule = fx.Module("etcd_client",
	fx.Provide(
		func(cfg *Config) (*clientv3.Client, error) {
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   cfg.EtcdConfig.Endpoints,
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				return nil, fmt.Errorf("etcd client init failed: %w", err)
			}
			return client, nil
		},
	),
	fx.Options(etcdClientHub),
)

var EtcdModule = fx.Module("etcd",
	fx.Options(EtcdClientModule),
	fx.Invoke(
		func(lc fx.Lifecycle, cfg *Config, client *clientv3.Client) error {
			etcd := etcdregister.NewEtcdRegister(cfg.EtcdConfig.Endpoints, 5*time.Second.Milliseconds())
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return etcd.Register(cfg.EtcdConfig.ServiceName, cfg.GrpcAddr)
				},
				OnStop: func(ctx context.Context) error {
					return etcd.UnRegister(cfg.EtcdConfig.ServiceName, cfg.GrpcAddr)
				},
			})
			return nil
		},
	))
