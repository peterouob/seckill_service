package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnRedis() *redis.Client {
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
