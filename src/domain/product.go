package domain

import "context"

type Product struct {
	Id              uint   `json:"id"`
	Name            string `json:"name"`
	AmountAvailable uint   `json:"amount_available"`
	Cost            uint   `json:"cost"`
	SellerId        uint   `json:"seller_id"`
}

type Bill struct {
	TotalSpent      uint             `json:"total_spent"`
	Items           map[Product]uint `json:"items"`
	RemainedDeposit []Coin           `json:"remained_deposit"`
}

type ProductService interface {
	Add(ctx context.Context, name string, amount, cost uint) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, id uint, name string, amount, cost uint) (*Product, error)
	Delete(ctx context.Context, id uint) error
	Buy(ctx context.Context, cart map[uint]int) (*Bill, error)
}

type ProductRepository interface {
	Insert(ctx context.Context, p Product) (uint, error)
	FindById(ctx context.Context, id uint) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uint) error
}
