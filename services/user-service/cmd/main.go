package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/services/user-service/internal/controller"
	"github.com/peterouob/seckill_service/services/user-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/services/user-service/internal/infrastructure/usergrpc"
	"github.com/peterouob/seckill_service/services/user-service/internal/router"
	"github.com/peterouob/seckill_service/services/user-service/internal/service"
	"github.com/peterouob/seckill_service/utils/database"
	"google.golang.org/grpc"
)

func main() {

	db := database.ConnPostgresql()
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := router.InitRouter(userController)

	userGrpc := usergrpc.NewUserGrpcHandlers(userService)
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Grpc Listen Error: %v", err)
		}

		grpcServer := grpc.NewServer()
		userproto.RegisterUserServiceServer(grpcServer, userGrpc)

		log.Println("gRPC server start on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC 伺服器啟動失敗: %v", err)
		}
	}()

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
