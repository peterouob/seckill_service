package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/peterouob/seckill_service/app/seckill-service/pkg/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SeckillRepo interface {
	DeductStock(productID, userID string) (int, error)
	ReduceStock(productID string) error
}

type seckillRepoImpl struct {
	ctx         context.Context
	redisClient *redis.Client
	db          *gorm.DB
}

func NewSeckillRepo(ctx context.Context, rdb *redis.Client, db *gorm.DB) SeckillRepo {
	return &seckillRepoImpl{
		ctx:         ctx,
		redisClient: rdb,
		db:          db,
	}
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

func (r *seckillRepoImpl) DeductStock(productID string, userID string) (int, error) {
	stockKey := fmt.Sprintf("secKill:{%s}:stock", productID)
	usersKey := fmt.Sprintf("secKill:{%s}:users", productID)

	keys := []string{stockKey, usersKey}
	result, err := seckillScript.Run(r.ctx, r.redisClient, keys, userID).Result()
	if err != nil {
		return 0, fmt.Errorf("redis腳本執行失敗: %v", err)
	}

	return int(result.(int64)), nil
}

func (r *seckillRepoImpl) ReduceStock(productId string) error {
	return r.db.WithContext(r.ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Stock{}).
			Where("product_id = ? AND stock > 0", productId).
			Update("stock", gorm.Expr("stock - 1"))

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("out of stock")
		}

		return nil
	})
}
