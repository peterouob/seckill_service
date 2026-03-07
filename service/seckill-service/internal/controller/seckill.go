package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/api/seckillproto"
	"github.com/peterouob/seckill_service/service/seckill-service/pkg/model"
)

type SeckillController struct {
	grpcClient seckillproto.SeckillServiceClient
}

func NewSeckillController(grpcClient seckillproto.SeckillServiceClient) *SeckillController {
	return &SeckillController{
		grpcClient: grpcClient,
	}
}

func (s *SeckillController) Buy(c *gin.Context) {
	var seckillReq model.SeckillReq

	if err := c.ShouldBindJSON(&seckillReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "should bind error"})
		return
	}

	grpcReq := &seckillproto.SeckillRequest{
		UserId:    seckillReq.UserID,
		ProductId: seckillReq.ProductID,
	}

	resp, err := s.grpcClient.Seckill(context.Background(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}
