package model

import (
	"time"
)

type Order struct {
	UserId    string    `json:"user_id"`
	ProductId string    `json:"product_id"`
	CreateAt  time.Time `json:"create_at"`
}

func (o Order) TableName() string {
	return "order"
}
