package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/app/seckill-service/internal/service"
	"github.com/peterouob/seckill_service/app/seckill-service/pkg/model"
)

type SeckillController struct {
	ctx            context.Context
	seckillService service.SeckillService
}

func NewSeckillController(ctx context.Context, svc service.SeckillService) *SeckillController {
	return &SeckillController{
		ctx:            ctx,
		seckillService: svc,
	}
}

func (ctl *SeckillController) Buy(c *gin.Context) {
	var seckillReq model.SeckillReq

	if err := c.ShouldBindJSON(&seckillReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "should bind error"})
		return
	}

	err := ctl.seckillService.Buy(ctl.ctx, seckillReq.UserID, seckillReq.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
