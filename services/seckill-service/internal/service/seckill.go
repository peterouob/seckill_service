package service

import (
	"context"
	"errors"
	"time"

	"github.com/peterouob/seckill_service/services/seckill-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/services/seckill-service/pkg/model"
	"github.com/peterouob/seckill_service/utils/mq"
)

type SeckillService interface {
	Buy(ctx context.Context, userId, productId string) error
}

type seckillServiceImpl struct {
	repo repository.SeckillRepo
	pd   *mq.Producer
}

func NewSeckillService(repo repository.SeckillRepo, pd *mq.Producer) SeckillService {
	return &seckillServiceImpl{
		repo: repo,
		pd:   pd,
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
		order := model.Order{
			UserId:    userId,
			ProductId: productId,
			CreateAt:  time.Now(),
		}
		s.pd.Send("order", order)
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
