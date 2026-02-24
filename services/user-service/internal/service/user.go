package service

import (
	"context"

	"github.com/peterouob/seckill_service/services/user-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/services/user-service/pkg/model"
)

type UserService interface {
	Login(context.Context, model.UserLoginReq)
	Register(context.Context, model.UserRegisterReq)
}

type userServiceImpl struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (u *userServiceImpl) Login(ctx context.Context, req model.UserLoginReq) {

}

func (u *userServiceImpl) Register(ctx context.Context, req model.UserRegisterReq) {

}
