package usergrpc

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/services/user-service/pkg/verify"
)

type TokenValid struct {
	userproto.UnimplementedUserServiceServer
}

func NewTokenValidServer() *TokenValid {
	return &TokenValid{}
}

func (auth TokenValid) TokenValid(ctx context.Context, req *userproto.TokenValidRequest) (*userproto.TokenValidResponse, error) {
	tokenString := req.GetToken()
	token := verify.TokenVerify(tokenString)

	if !token.Valid {
		return &userproto.TokenValidResponse{
			Msg: "Token is invalid",
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &userproto.TokenValidResponse{
			Msg: "Token claims are not valid",
		}, nil
	}
	userID := int64(claims["userId"].(float64))
	return &userproto.TokenValidResponse{
		Id:  userID,
		Msg: "Token is valid",
	}, nil
}
