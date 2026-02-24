package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type SeckillRepo interface {
	DeductStock(ctx context.Context, productID, userID string) (int, error)
}

type seckillRepoImpl struct {
	redisClient *redis.Client
}

func NewSeckillRepo(rdb *redis.Client) SeckillRepo {
	return &seckillRepoImpl{redisClient: rdb}
}

var seckillScript = redis.NewScript(`
    local stock_key = KEYS[1]
    local users_key = KEYS[2]
    
    local user_id = ARGV[1]

    local user_exist = redis.call('sismember', users_key, user_id)
    if tonumber(user_exist) == 1 then
        return 2
    end

    local stock_val = redis.call('get', stock_key)
    if not stock_val then
        return -1
    end

    if tonumber(stock_val) <= 0 then
        return 0
    else
        redis.call('decr', stock_key)
        redis.call('sadd', users_key, user_id)
        return 1
    end
`)

func (r *seckillRepoImpl) DeductStock(ctx context.Context, productID string, userID string) (int, error) {
	stockKey := fmt.Sprintf("secKill:{%s}:stock", productID)
	usersKey := fmt.Sprintf("secKill:{%s}:users", productID)

	keys := []string{stockKey, usersKey}

	log.Println(stockKey, usersKey, keys)
	result, err := seckillScript.Run(ctx, r.redisClient, keys, userID).Result()
	if err != nil {
		return 0, fmt.Errorf("redis腳本執行失敗: %v", err)
	}

	return int(result.(int64)), nil
}
