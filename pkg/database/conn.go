package database

import (
	"context"
	"fmt"

	"github.com/peterouob/seckill_service/pkg/config"
	"github.com/peterouob/seckill_service/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connMysql(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("failed to connect database %v\n", err))
	}
	return db
}

type MigrationParams struct {
	fx.In
	Db     *gorm.DB
	Models []any `group:"gorm_models"`
}

func runMigration(m MigrationParams) error {
	if len(m.Models) > 0 {
		logger.Log("migration model table ...")

		for _, model := range m.Models {
			if err := m.Db.AutoMigrate(model); err != nil {
				return err
			}
		}
	}
	return nil
}

var MySQLModule = fx.Module("mysql",
	fx.Provide(
		func(lc fx.Lifecycle, cfg *config.Config) *gorm.DB {
			db := connMysql(cfg.DbDSN)
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
	fx.Invoke(runMigration),
)
