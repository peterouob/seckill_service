package repo

import (
	"context"

	"github.com/peterouob/seckill_service/internal/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByName(context.Context, model.UserLoginReq) (*model.User, error)
	Register(context.Context, model.UserRegisterReq)
}

type UserRepoImpl struct {
	db *gorm.DB
}

func (u *UserRepoImpl) Register(ctx context.Context, req model.UserRegisterReq) {
	//TODO implement me
	panic("implement me")
}

var _ (UserRepo) = (*UserRepoImpl)(nil)

func NewUserRepo(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{
		db: db,
	}
}

func (u *UserRepoImpl) GetUserByName(ctx context.Context, user model.UserLoginReq) (*model.User, error) {
	return nil, nil
}
