package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/peterouob/seckill_service/services/user-service/pkg/model"
	"github.com/peterouob/seckill_service/utils/logs"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByName(context.Context, string) (*model.User, error)
	Register(context.Context, model.UserRegisterReq)
	Login(ctx context.Context, username string, password string) (*model.User, error)
}

type userRepoImpl struct {
	db *gorm.DB
}

var _ (UserRepo) = (*userRepoImpl)(nil)

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepoImpl{
		db: db,
	}
}

func (u *userRepoImpl) GetUserByName(_ context.Context, username string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepoImpl) Register(ctx context.Context, req model.UserRegisterReq) {
	if _, err := u.GetUserByName(ctx, req.Username); errors.Is(err, gorm.ErrRecordNotFound) {
		u7, _ := uuid.NewV7()
		var user model.User
		user.UserId = u7.String()
		user.Username = req.Username
		user.Password = req.Password
		if err := u.db.Create(&user).Error; err != nil {
			logs.Error("create user error", err)
		}
	}
}

func (u *userRepoImpl) Login(ctx context.Context, username string, password string) (*model.User, error) {
	user, err := u.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("password error")
	}

	return user, nil
}
