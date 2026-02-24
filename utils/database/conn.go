package database

import (
	"context"
	"fmt"

	"github.com/peterouob/seckill_service/services/seckill-service/pkg/model"
	usermodel "github.com/peterouob/seckill_service/services/user-service/pkg/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO: use Errors.New replace panic

func ConnPostgresql() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("failed to connect database %v\n", err))
	}
	if err := db.AutoMigrate(&model.Order{}, &usermodel.User{}); err != nil {
		panic(fmt.Errorf("failed to migrate database %v\n", err))
	}
	return db
}

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
