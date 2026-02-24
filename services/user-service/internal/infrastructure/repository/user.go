package repository

import (
	"context"

	"github.com/peterouob/seckill_service/services/user-service/internal/domain"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByName(context.Context, domain.UserLoginReq) (*domain.User, error)
	Register(context.Context, domain.UserRegisterReq)
}

type UserRepoImpl struct {
	db *gorm.DB
}

func (u *UserRepoImpl) Register(ctx context.Context, req domain.UserRegisterReq) {
	//TODO implement me
	panic("implement me")
}

var _ (UserRepo) = (*UserRepoImpl)(nil)

func NewUserRepo(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{
		db: db,
	}
}

func (u *UserRepoImpl) GetUserByName(ctx context.Context, user domain.UserLoginReq) (*domain.User, error) {
	return nil, nil
}
