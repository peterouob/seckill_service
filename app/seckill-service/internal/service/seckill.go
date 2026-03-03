package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/peterouob/seckill_service/app/seckill-service/internal/infrastructure/repository"
	etcdregister "github.com/peterouob/seckill_service/pkg/etcd"
	"github.com/peterouob/seckill_service/pkg/kafka"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type SeckillService interface {
	Buy(ctx context.Context, userId, productId string) error
}

type seckillServiceImpl struct {
	repo    repository.SeckillRepo
	pd      *kafka.Producer
	etcd    *etcdregister.EtcdRegister
	ctx     context.Context
	session *concurrency.Session
}

func NewSeckillService(ctx context.Context, repo repository.SeckillRepo, pd *kafka.Producer, etcd *etcdregister.EtcdRegister) SeckillService {
	session, _ := concurrency.NewSession(etcd.Client.Client(), concurrency.WithTTL(5))
	s := &seckillServiceImpl{
		ctx:     ctx,
		repo:    repo,
		pd:      pd,
		etcd:    etcd,
		session: session,
	}
	return s
}

func (s *seckillServiceImpl) Buy(ctx context.Context, userId, productId string) error {
	result, err := s.repo.DeductStock(productId, userId)
	if err != nil {
		return errors.New("error in Buy service")
	}

	orderId := fmt.Sprintf("%s_%s", userId, productId)
	switch result {
	case 1:
		s.reduceStock(ctx, userId, productId, orderId)
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

// reduceStock TOOD:use etcd distributed lock make sure only one can reduce for the truth stock in mysql
func (s *seckillServiceImpl) reduceStock(ctx context.Context, userId, productId, orderId string) error {

	pLock := fmt.Sprintf("/seckill/locks/%s", productId)
	mutex := concurrency.NewMutex(s.session, pLock)

	if err := mutex.Lock(ctx); err != nil {
		return fmt.Errorf("could not acquire lock for %s: %w", productId, err)
	}
	defer mutex.Unlock(ctx)

	if err := s.repo.ReduceStock(productId); err != nil {
		return errors.New("error in reduce stock service")
	}

	return nil
}
