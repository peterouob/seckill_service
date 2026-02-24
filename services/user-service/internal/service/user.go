package service

import (
	"context"
	"errors"

	"github.com/peterouob/seckill_service/services/user-service/internal/infrastructure/repository"
	"github.com/peterouob/seckill_service/services/user-service/pkg/model"
	"github.com/peterouob/seckill_service/services/user-service/pkg/verify"
	"github.com/peterouob/seckill_service/utils/logs"
)

type UserService interface {
	Login(context.Context, model.UserLoginReq) (string, error)
	Register(context.Context, model.UserRegisterReq) error
}

type userServiceImpl struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (u *userServiceImpl) Login(ctx context.Context, req model.UserLoginReq) (string, error) {
	user, err := u.userRepo.Login(ctx, req.Username, req.Password)
	if err != nil {
		logs.Error("login error", err)
		return "", err
	}
	token := verify.NewToken(int64(user.ID))
	token.CreateToken()
	return token.AccessToken, nil
}

var ErrNotSamePassword = errors.New("password and check password not the same")

func (u *userServiceImpl) Register(ctx context.Context, req model.UserRegisterReq) error {
	if req.Password != req.CheckPassword {
		logs.Error("password not same", ErrNotSamePassword)
		return ErrNotSamePassword
	}
	u.userRepo.Register(ctx, req)
	return nil
}
