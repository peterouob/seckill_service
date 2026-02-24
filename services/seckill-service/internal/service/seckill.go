package service

import (
	"context"
	"errors"

	"github.com/peterouob/seckill_service/services/seckill-service/internal/infrastructure/repository"
)

type SeckillService interface {
	Buy(ctx context.Context, userId, productId string) error
}

type seckillServiceImpl struct {
	repo repository.SeckillRepo
}

func NewSeckillService(repo repository.SeckillRepo) SeckillService {
	return &seckillServiceImpl{
		repo: repo,
	}
}

func (s *seckillServiceImpl) Buy(ctx context.Context, userId, productId string) error {
	result, err := s.repo.DeductStock(ctx, productId, userId)
	if err != nil {
		return errors.New("error in Buy service")
	}

	switch result {
	case 1:
		// TODO: kafka async write in db
		return nil
	case 2:
		return errors.New("您已經搶購過此商品了")
	case 0:
		return errors.New("沒貨了")
	case -1:
		return errors.New("活動尚未開始")
	default:
		return errors.New("error")
	}
}
