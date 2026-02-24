package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/internal/service"
)

type UserController struct {
	srv service.UserService
}

func NewUserController(srv service.UserService) *UserController {
	return &UserController{
		srv: srv,
	}
}

func (u *UserController) Login(c *gin.Context) {

}
