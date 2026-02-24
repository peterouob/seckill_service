package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peterouob/seckill_service/services/seckill-service/internal/service"
	"github.com/peterouob/seckill_service/services/seckill-service/pkg/model"
)

type SeckillController struct {
	seckillService service.SeckillService
}

func NewSeckillController(svc service.SeckillService) *SeckillController {
	return &SeckillController{
		seckillService: svc,
	}
}

func (ctl *SeckillController) Buy(c *gin.Context) {
	var seckillReq model.SeckillReq

	if err := c.ShouldBindJSON(&seckillReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "should bind error"})
		return
	}

	err := ctl.seckillService.Buy(c.Request.Context(), seckillReq.UserID, seckillReq.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
