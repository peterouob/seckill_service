package usergrpc

import (
	"context"

	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/services/user-service/internal/service"
	"github.com/peterouob/seckill_service/services/user-service/pkg/model"
)

type UserGrpcHandler struct {
	userproto.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserGrpcHandlers(srv service.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{
		userService: srv,
	}
}

func (u *UserGrpcHandler) UserLogin(ctx context.Context, in *userproto.UserLoginReq) (*userproto.UserLoginResp, error) {
	user := model.UserLoginReq{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	}
	u.userService.Login(ctx, user)
	return &userproto.UserLoginResp{
		Msg: "success",
	}, nil
}

func (u *UserGrpcHandler) UserRegister(ctx context.Context, in *userproto.UserRegisterReq) (*userproto.UserRegisterResp, error) {
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
