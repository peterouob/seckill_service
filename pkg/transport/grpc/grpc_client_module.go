package transport

import (
	"context"
	"fmt"

	etcdclient "github.com/peterouob/seckill_service/pkg/etcd/client"
	"github.com/peterouob/seckill_service/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewConn(
	lc fx.Lifecycle,
	hub *etcdclient.ServiceHub,
	serviceName string,
) (*grpc.ClientConn, error) {

	endpoints := hub.GetServiceEndPoint(serviceName)
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoint found for service: %s", serviceName)
	}

	addr := endpoints[0]
	logger.Logf("connecting to %s at %s", serviceName, addr)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s failed: %w", addr, err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			conn.Connect()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return conn, nil
}
