package cache

import (
	"context"
	"fmt"

	"github.com/peterouob/seckill_service/pkg/injection"
	"github.com/peterouob/seckill_service/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func connRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Errorf("failed to connect redis %v\n", err))
	}

	return rdb
}

var RedisModule = fx.Module("redis", fx.Provide(
	func(lc fx.Lifecycle, cfg *injection.Config) *redis.Client {
		rdb := connRedis()
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Log("redis connect closing ...")
				if err := rdb.Close(); err != nil {
					return err
				}
				return nil
			},
		})
		return rdb
	},
))
