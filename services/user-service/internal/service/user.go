package service

import (
	"context"

	"github.com/peterouob/seckill_service/services/user-service/internal/domain"
	"github.com/peterouob/seckill_service/services/user-service/internal/infrastructure/repository"
)

type UserService interface {
	Login(context.Context, domain.UserLoginReq)
}

type UserServiceImpl struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (u *UserServiceImpl) Login(ctx context.Context, req domain.UserLoginReq) {

}
