package server

import (
	"context"

	"github.com/peterouob/seckill_service/api/seckillproto"
	"github.com/peterouob/seckill_service/service/seckill-service/internal/service"
	"github.com/peterouob/seckill_service/service/seckill-service/pkg/model"
)

type SeckillHandler struct {
	seckillproto.UnimplementedSeckillServiceServer
	seckillService service.SeckillService
}

func NewSeckillGrpcHandler(seckillService service.SeckillService) *SeckillHandler {
	return &SeckillHandler{
		seckillService: seckillService,
	}
}

func (s *SeckillHandler) Seckill(ctx context.Context, in *seckillproto.SeckillRequest) (*seckillproto.SeckillResponse, error) {
	req := model.SeckillReq{
		UserID:    in.UserId,
		ProductID: in.ProductId,
	}

	if err := s.seckillService.Buy(ctx, req.UserID, req.ProductID); err != nil {
		return nil, err
	}

	resp := &seckillproto.SeckillResponse{
		Msg: "success",
	}
	return resp, nil
}
