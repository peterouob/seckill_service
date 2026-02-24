package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/peterouob/seckill_service/services/seckill-service/internal/controller"
	"github.com/peterouob/seckill_service/services/seckill-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/services/seckill-service/internal/router"
	"github.com/peterouob/seckill_service/services/seckill-service/internal/service"
	"github.com/peterouob/seckill_service/utils/database"
	"github.com/peterouob/seckill_service/utils/logs"
)

func main() {
	logs.InitLogger("seckill")
	//db := database.ConnPostgresql()
	rdb := database.ConnRedis()

	repo := repository.NewSeckillRepo(rdb)
	srv := service.NewSeckillService(repo)
	ctl := controller.NewSeckillController(srv)

	r := router.InitRouter(ctl)

	server := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logs.Log("Starting server ...")
		serverErrors <- server.ListenAndServe()
	}()

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErrors:
		logs.Logf("Error starting server ... %v\n", err)
	case sig := <-shutDown:
		logs.ErrorMsgF("Server is shutting due to the %v signal\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logs.ErrorMsgF("Could n ot stdio the server gracefully %v\n", err)
			_ = server.Close()
		}
	}
}
