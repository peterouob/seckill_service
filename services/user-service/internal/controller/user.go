package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/api/userproto"
	"github.com/peterouob/seckill_service/services/user-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	grpcClient userproto.UserServiceClient
}

func NewUserController(grpcClient userproto.UserServiceClient) *UserController {
	return &UserController{
		grpcClient: grpcClient,
	}
}

func (u *UserController) Login(c *gin.Context) {
	var loginReq domain.UserLoginReq
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userproto.UserLoginReq{
		Username: loginReq.Username,
		Password: loginReq.Password,
	}

	resp, err := u.grpcClient.UserLogin(context.Background(), grpcReq)
	if err != nil {
		state, ok := status.FromError(err)
		if ok && state.Code() == codes.Unauthenticated {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": resp})
}
