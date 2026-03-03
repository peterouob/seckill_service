package module

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/app/user-service/internal/controller"
	"github.com/peterouob/seckill_service/app/user-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/app/user-service/internal/infrastructure/usergrpc"
	"github.com/peterouob/seckill_service/app/user-service/internal/router"
	"github.com/peterouob/seckill_service/app/user-service/internal/service"
	"github.com/peterouob/seckill_service/app/user-service/pkg/model"
	"github.com/peterouob/seckill_service/pkg/injection"
	"github.com/peterouob/seckill_service/pkg/netutil"
	"github.com/peterouob/seckill_service/pkg/pool"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type lazyClient struct {
	mu     sync.RWMutex
	client userproto.UserServiceClient
}

func (l *lazyClient) get() userproto.UserServiceClient {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.client
}

func (l *lazyClient) set(c userproto.UserServiceClient) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.client = c
}

func (l *lazyClient) UserLogin(ctx context.Context, in *userproto.UserLoginReq, opts ...grpc.CallOption) (*userproto.UserLoginResp, error) {
	return l.get().UserLogin(ctx, in, opts...)
}

func (l *lazyClient) UserRegister(ctx context.Context, in *userproto.UserRegisterReq, opts ...grpc.CallOption) (*userproto.UserRegisterResp, error) {
	return l.get().UserRegister(ctx, in, opts...)
}

var UserModule = fx.Module("user",

	fx.Provide(
		repository.NewUserRepo,
		service.NewUserService,
		usergrpc.NewUserGrpcHandlers,

		provideUserGrpcClient,

		controller.NewUserController,
		fx.Annotate(
			func() []any {
				return []any{
					&model.User{},
				}
			},
			fx.ResultTags(`group:"gorm_models,flatten"`),
		),
	),

	fx.Invoke(func(srv *grpc.Server, h *usergrpc.UserGrpcHandler) {
		userproto.RegisterUserServiceServer(srv, h)
	}),

	fx.Invoke(func(r *gin.Engine, ctl *controller.UserController) {
		router.InitRouter(r, ctl)
	}),
)

func provideUserGrpcClient(
	lc fx.Lifecycle,
	cfg *injection.Config,
	ready injection.GrpcServerReady,
) userproto.UserServiceClient {

	grpcAddr := netutil.FormatIP(cfg.GrpcAddr)
	p := pool.New(grpcAddr, pool.DefaultOption)
	var conn pool.Conn

	lazy := &lazyClient{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			select {
			case <-ready.C:
			case <-ctx.Done():
				return fmt.Errorf("wait for grpc server ready timeout: %w", ctx.Err())
			}

			var err error
			conn, err = p.Get()
			if err != nil {
				return fmt.Errorf("get grpc conn failed: %w", err)
			}
			lazy.set(userproto.NewUserServiceClient(conn.Value()))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.Put(conn)
		},
	})

	return lazy
}
