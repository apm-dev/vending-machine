package domain

import "context"

type Product struct {
	Id              int64
	Name            string
	AmountAvailable int64
	Cost            int64
	SellerId        int64
}

type ProductService interface {
	
}

type ProductRepository interface {
	FindById(ctx context.Context, id int64) *Product
	List(ctx context.Context) []Product
	Insert(ctx context.Context, p Product) int64
	Update(ctx context.Context, p *Product)
	Delete(ctx context.Context, id int64)
}
