package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Rdb *redis.Client
)

// TODO: use Errors.New replace panic

func ConnPostgresql() *gorm.DB {
	dsn := "host=localhost user=root password=password dbname=seckill port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("failed to connect database %v\n", err))
	}

	return db
}

func ConnRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer func() {
		_ = rdb.Close()
	}()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Errorf("failed to connect redis %v\n", err))
	}

	Rdb = rdb
}
