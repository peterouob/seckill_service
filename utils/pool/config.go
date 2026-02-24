// Package configs grpc pool 主要參考: https://github.com/shimingyah/pool
package pool

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	DialTimeout = 5 * time.Second

	KeepAliveTime = time.Duration(10) * time.Second

	KeepAliveTimeout = time.Duration(3) * time.Second

	InitialWindowSize = 1 << 30

	InitialConnWindowSize = 1 << 30

	MaxSendMsgSize = 4 << 30

	MaxRecvMsgSize = 4 << 30
)

var MaxBackoffDelay = grpc.ConnectParams{
	Backoff:           backoff.Config{BaseDelay: 100 * time.Millisecond, MaxDelay: 10 * time.Second, Multiplier: 1.6, Jitter: 0.1},
	MinConnectTimeout: 3 * time.Second,
}

type Option struct {
	Dial                 func(addr string) (*grpc.ClientConn, error)
	MaxIdle              int32
	MaxActive            int32
	MaxConcurrentStreams int32
	Reuse                bool
}

var DefaultOption = Option{
	Dial:                 Dial,
	MaxIdle:              10,
	MaxActive:            10,
	MaxConcurrentStreams: 64,
	Reuse:                true,
}

func Dial(addr string) (*grpc.ClientConn, error) {
	g, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(MaxBackoffDelay),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	return g, nil
}
