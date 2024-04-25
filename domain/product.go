package domain

import "context"

const (
)

type Product struct {
	ID            int64
	Name string
	Picture string
	ProductDetailId int64
	ProductCategoryId int64
	IsActive bool
}

type ProductDetails struct{
	ID int64
	GenericName string
	Content string
	Manufacturer string
	Description string
	ProductClassification string
	ProductForm string
	UnitInPack string
	SellingUnit string
	Weight float64
	Height float64
	Length float64
	Width float64
}

type AddProductRequest struct{
	Product Product
	Category ProductCategories
	Details ProductDetails
}

type ProductCategories struct{
	ID int64
	Name string
	ParentId int64
}

type ProductRepository interface {
	GetByName(ctx context.Context, name string) (Product, error)
	AddProduct(ctx context.Context, product AddProductRequest) (Product,error)
}

type ProductService interface {
	GetByName(ctx context.Context) (Product, error)
}
