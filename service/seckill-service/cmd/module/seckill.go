package module

import (
	"github.com/peterouob/seckill_service/api/seckillproto"
	transport "github.com/peterouob/seckill_service/pkg/transport/grpc"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/infrastructure/seckillgrpc/server"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var SeckillModule = fx.Module("seckill-server",

	fx.Options(transport.GrpcServerModule),

	fx.Provide(
		repository.NewSeckillRepo,
		service.NewSeckillService,
		server.NewSeckillGrpcHandler,
	),

	fx.Invoke(func(srv *grpc.Server, h *server.SeckillHandler) {
		seckillproto.RegisterSeckillServiceServer(srv, h)
	}),
)
