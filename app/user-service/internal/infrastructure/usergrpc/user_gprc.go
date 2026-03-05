package usergrpc

import (
	"context"

	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/app/user-service/internal/service"
	"github.com/peterouob/seckill_service/app/user-service/pkg/model"
)

type UserHandler struct {
	userproto.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserGrpcHandlers(srv service.UserService) *UserHandler {
	return &UserHandler{
		userService: srv,
	}
}

func (u *UserHandler) UserLogin(ctx context.Context, in *userproto.UserLoginReq) (*userproto.UserLoginResp, error) {
	user := model.UserLoginReq{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	}
	token, err := u.userService.Login(ctx, user)
	if err != nil {
		return nil, err
	}
	return &userproto.UserLoginResp{
		Msg:   "success",
		Token: token,
	}, nil
}

func (u *UserHandler) UserRegister(ctx context.Context, in *userproto.UserRegisterReq) (*userproto.UserRegisterResp, error) {
	user := model.UserRegisterReq{
		Username:      in.GetUsername(),
		Password:      in.GetPassword(),
		CheckPassword: in.GetCheckPassword(),
	}
	_ = u.userService.Register(ctx, user)
	return &userproto.UserRegisterResp{
		Msg: "success",
	}, nil
}
