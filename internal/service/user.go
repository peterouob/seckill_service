package service

import (
	"context"

	"github.com/peterouob/seckill_service/internal/model"
	"github.com/peterouob/seckill_service/internal/repo"
)

type UserService interface {
	Login(context.Context, model.UserLoginReq)
}

type UserServiceImpl struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (u *UserServiceImpl) Login(ctx context.Context, req model.UserLoginReq) {

}
