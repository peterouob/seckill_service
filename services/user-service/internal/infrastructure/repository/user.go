package repository

import (
	"context"

	"github.com/peterouob/seckill_service/services/user-service/pkg/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByName(context.Context, model.UserLoginReq) (*model.User, error)
	Register(context.Context, model.UserRegisterReq)
}

type userRepoImpl struct {
	db *gorm.DB
}

func (u *userRepoImpl) Register(ctx context.Context, req model.UserRegisterReq) {
	//TODO implement me
	panic("implement me")
}

var _ (UserRepo) = (*userRepoImpl)(nil)

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepoImpl{
		db: db,
	}
}

func (u *userRepoImpl) GetUserByName(ctx context.Context, user model.UserLoginReq) (*model.User, error) {
	return nil, nil
}
