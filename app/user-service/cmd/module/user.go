package module

import (
	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/app/user-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/app/user-service/internal/infrastructure/usergrpc"
	"github.com/peterouob/seckill_service/app/user-service/internal/service"
	"github.com/peterouob/seckill_service/app/user-service/pkg/model"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var UserModule = fx.Module("user",

	fx.Provide(
		repository.NewUserRepo,
		service.NewUserService,
		usergrpc.NewUserGrpcHandlers,

		fx.Annotate(
			func() []any {
				return []any{
					&model.User{},
				}
			},
			fx.ResultTags(`group:"gorm_models,flatten"`),
		),
	),

	fx.Invoke(func(srv *grpc.Server, h *usergrpc.UserHandler) {
		userproto.RegisterUserServiceServer(srv, h)
	}),
)
