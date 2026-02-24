package main

import (
	"context"
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
	etcdregister "github.com/peterouob/seckill_service/utils/etcd"
	"github.com/peterouob/seckill_service/utils/logs"
	"github.com/peterouob/seckill_service/utils/pool"
	"google.golang.org/grpc"
)

func main() {
	logs.InitLogger("user")
	db := database.ConnPostgresql()
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userGrpc := usergrpc.NewUserGrpcHandlers(userService)

	grpcChannel := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	go func() {
		lis, err := net.Listen("tcp", ":50050")
		if err != nil {
			logs.Error("failed to listen: %v\n", err)
		}
		grpcServer := grpc.NewServer()
		userproto.RegisterUserServiceServer(grpcServer, userGrpc)
		logs.Log("gRPC server start on :50050")
		close(grpcChannel)
		if err := grpcServer.Serve(lis); err != nil {
			logs.Error("start grpc server error", err)
		}
	}()

	select {
	case <-grpcChannel:
		logs.Log("grpc server ready ...")
	case <-ctx.Done():
		logs.Warn("grpc server start timeout")

	}

	p := pool.New("127.0.0.1:50050", pool.DefaultOption)
	conn, _ := p.Get()
	client := userproto.NewUserServiceClient(conn.Value())
	userController := controller.NewUserController(client)

	r := router.InitRouter(userController)

	server := &http.Server{
		Addr:    ":8083",
		Handler: r,
	}

	etcd := etcdregister.NewEtcdRegister([]string{"127.0.0.1:2379"}, 3)
	etcd.Register("user", "127.0.0.1:8083")

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
