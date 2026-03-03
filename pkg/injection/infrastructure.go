package injection

import (
	"context"

	"github.com/peterouob/seckill_service/pkg/cache"
	"github.com/peterouob/seckill_service/pkg/database"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type MigrationParams struct {
	fx.In
	DB     *gorm.DB
	Models []any `group:"gorm_models"`
}

func RunMigration(m MigrationParams) error {
	if len(m.Models) > 0 {
		logger.Log("migration model table ...")

		for _, model := range m.Models {
			if err := m.DB.AutoMigrate(model); err != nil {
				return err
			}
		}
	}
	return nil
}

var MySQLModule = fx.Module("mysql",
	fx.Provide(
		func(lc fx.Lifecycle, cfg *Config) *gorm.DB {
			db := database.ConnMysql(cfg.DbDSN)
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					logger.Log("database connection closing...")
					sqlDb, err := db.DB()
					if err != nil {
						return err
					}
					return sqlDb.Close()
				},
			})
			return db
		},
	),
	fx.Invoke(RunMigration),
)

var RedisModule = fx.Module("redis", fx.Provide(
	func(lc fx.Lifecycle, cfg *Config) *redis.Client {
		rdb := cache.ConnRedis()
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

var EtcdModule = fx.Module("etcd", fx.Provide(
	func(lc fx.Lifecycle, cfg *Config) *etcdregister.EtcdRegister {
		etcdServiceName := cfg.EtcdConfig.ServiceName
		etcdEndpoints := cfg.EtcdConfig.Endpoints
		etcd := etcdregister.NewEtcdRegister(etcdEndpoints, 3)

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := etcd.Register(etcdServiceName, cfg.GrpcAddr); err != nil {
					return err
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Log("etcd service stopped")
				if err := etcd.UnRegister(cfg.EtcdConfig.ServiceName, cfg.GrpcAddr); err != nil {
					logger.Error("unregister etcd service err: ", err)
					return err
				}
				return nil
			},
		})
		return etcd
	},
))
