package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/peterouob/seckill_service/internal/controller"
	"github.com/peterouob/seckill_service/internal/repo"
	"github.com/peterouob/seckill_service/internal/router"
	"github.com/peterouob/seckill_service/internal/service"
	"github.com/peterouob/seckill_service/pkg/database"
)

func main() {

	user := newUser()

	r := router.InitRouter(user)

	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	go func() {
		database.ConnRedis()
		database.ConnPostgresql()
		log.Println("Starting server ...")
		serverErrors <- server.ListenAndServe()
	}()

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErrors:
		log.Printf("Error starting server ... %v\n", err)
	case sig := <-shutDown:
		log.Printf("Server is shutting due to the %v signal\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could n ot stdio the server gracefully %v\n", err)
			_ = server.Close()
		}
	}
}

func newUser() *controller.UserController {
	db := database.ConnPostgresql()
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	return userController
}
