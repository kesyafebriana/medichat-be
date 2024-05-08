package domain

import (
	"context"
	"mime/multipart"
)

const (
)

type Product struct {
	ID            int64
	Name string
	Slug string
	Picture *string
	KeyWord string
	ProductDetailId int64
	ProductCategoryId int64
	IsActive bool
}

type ProductDetails struct{
	ID int64
	GenericName string
	Content string
	Composition string
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
	Name string
	ProductCategoryId int64

	GenericName string
	Composition string
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

type UpdateProductRequest struct{
	Name string
	ProductCategoryId int64

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

type ProductsQuery struct {
	Page       int64
	Limit      int64
	Latitude   *float64
	Longitude  *float64
	Term       string
	SortBy     string
	SortType   string
}

const (
	ProductSortById     = "id"
	ProductSortByName   = "name"
)

func DefaultProductsQuery() ProductsQuery {
	return ProductsQuery{
		Page:     1,
		SortBy:   ProductSortById,
		SortType: "ASC",
	}
}

type ProductRepository interface {
	GetByName(ctx context.Context, name string) (Product, error)
	GetById(ctx context.Context, id int64) (Product, error)
	GetPageInfo(ctx context.Context, query ProductsQuery) (PageInfo, error)
	GetProducts(ctx context.Context, query ProductsQuery) ([]Product, error)
	GetProductsFromArea(ctx context.Context, query ProductsQuery) ([]Product, error)
	GetPageInfoFromArea(ctx context.Context, query ProductsQuery) (PageInfo, error)
	GetBySlug(ctx context.Context, slug string) (Product, error)
	Add(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	SoftDeleteBySlug(ctx context.Context, slug string) error
	BulkSoftDeleteBySlug(ctx context.Context, slugs []string) error
}

type ProductDetailsRepository interface {
	GetById(ctx context.Context, id int64) (ProductDetails, error)
	Add(ctx context.Context, detail ProductDetails) (ProductDetails, error)
	Update(ctx context.Context, detail ProductDetails) (ProductDetails, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, slug string) (Product, ProductDetails, error)
	GetProducts(ctx context.Context, query ProductsQuery) ([]Product, PageInfo, error)
	GetProductLocation(ctx context.Context, query ProductsQuery) ([]Product, PageInfo, error)
	CreateProduct(ctx context.Context, request AddProductRequest, file *multipart.File) (Product, error)
	DeleteProducts(ctx context.Context, slug string) error
	UpdateProduct(ctx context.Context, slug string, request UpdateProductRequest, file *multipart.File) (Product, error)
}
