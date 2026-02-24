package usergrpc

import (
	"context"

	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/services/user-service/internal/domain"
	"github.com/peterouob/seckill_service/services/user-service/internal/service"
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
	user := domain.UserLoginReq{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	}
	u.userService.Login(ctx, user)
	return &userproto.UserLoginResp{
		Msg: "success",
	}, nil
}
