package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/peterouob/seckill_service/app/seckill-service/internal/controller"
	"github.com/peterouob/seckill_service/app/seckill-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/app/seckill-service/internal/router"
	"github.com/peterouob/seckill_service/app/seckill-service/internal/service"
	"github.com/peterouob/seckill_service/pkg/database"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/kafka"
	"github.com/peterouob/seckill_service/pkg/logger"
)

func main() {
	_ = godotenv.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.InitLogger("seckill")
	var DbDsn string
	if DbDsn = os.Getenv("DB_DSB"); DbDsn == "" {
		DbDsn = "root:123456@tcp(localhost:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db := database.ConnMysql(DbDsn)
	produce := kafka.NewProducer()
	defer produce.Close()
	etcd := etcdregister.NewEtcdRegister([]string{"localhost:2379"}, 5*time.Second.Milliseconds())
	rdb := database.ConnRedis()
	repo := repository.NewSeckillRepo(ctx, rdb, db)
	srv := service.NewSeckillService(ctx, repo, produce, etcd)
	ctl := controller.NewSeckillController(ctx, srv)

	// TODO: groupId use product ID instead
	consumer := kafka.NewConsumer("seckill", []string{"order"}, 1000, 1*time.Second, db)
	defer consumer.Close()
	go func() {
		consumer.StartConsume()
	}()

	r := router.InitRouter(ctl)

	server := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Log("Starting server ...")
		serverErrors <- server.ListenAndServe()
	}()

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErrors:
		logger.Logf("Error starting server ... %v\n", err)
	case sig := <-shutDown:
		logger.ErrorMsgF("Server is shutting due to the %v signal\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.ErrorMsgF("Could n ot stdio the server gracefully %v\n", err)
			_ = server.Close()
		}
	}
}
