package domain

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

// SeckillProduct store to redis not in db
type SeckillProduct struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}
