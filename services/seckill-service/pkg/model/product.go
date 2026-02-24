package model

type Product struct {
	Name        string `json:"name"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Img         string `json:"img"`
}

func (p Product) TableName() string {
	return "product"
}

type SeckillReq struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
}
