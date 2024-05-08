package dto

import (
	"medichat-be/constants"
	"medichat-be/domain"
	"mime/multipart"
)

type CreateProductForm struct {
	Name                  string  `form:"name" binding:"required"`
	GenericName           string  `form:"generic_name" binding:"required"`
	Composition           string  `form:"composition" binding:"required"`
	Content               string  `form:"content" binding:"required"`
	Manufacturer          string  `form:"manufacturer" binding:"required"`
	Description           string  `form:"description" binding:"required"`
	CategoryId            int64   `form:"category_id" binding:"required"`
	ProductClassification string  `form:"product_classification" binding:"required"`
	ProductForm           string  `form:"product_form" binding:"required"`
	UnitInPack            string  `form:"unit_in_pack" binding:"required"`
	SellingUnit           string  `form:"selling_unit" binding:"required"`
	Weight                float64 `form:"weight" binding:"required"`
	Length                float64 `form:"length" binding:"required"`
	Height                float64 `form:"height" binding:"required"`
	Width                 float64 `form:"width" binding:"required"`

	Picture *multipart.FileHeader `form:"picture"`
}

type ProductSlugParams struct {
	Slug string `uri:"slug" binding:"required"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	PageInfo PageInfoResponse  `json:"page_info"`
}

type UpdateProductForm struct {
	Slug string `uri:"slug"`

	Name                  string  `form:"name"`
	GenericName           string  `form:"generic_name"`
	Composition           string  `form:"composition"`
	Content               string  `form:"content"`
	Manufacturer          string  `form:"manufacturer"`
	Description           string  `form:"description"`
	ProductClassification string  `form:"product_classification"`
	CategoryId            int64   `form:"category_id"`
	ProductForm           string  `form:"product_form"`
	UnitInPack            string  `form:"unit_in_pack"`
	SellingUnit           string  `form:"selling_unit"`
	Weight                float64 `form:"weight"`
	Length                float64 `form:"length"`
	Height                float64 `form:"height"`
	Width                 float64 `form:"width"`

	Picture *multipart.FileHeader `form:"picture"`
}

type GetProductsQuery struct {
	Page      int64    `form:"page" binding:"numeric,omitempty,min=1"`
	Limit     int64    `form:"limit" binding:"numeric,omitempty,min=1"`
	Term      string   `form:"term"`
	Longitude *float64 `form:"long"`
	Latitude  *float64 `form:"lat"`
	SortBy    string   `form:"sort_by" binding:"omitempty,oneof=name slug"`
	SortType  string   `form:"sort_type" binding:"omitempty,oneof=ASC DESC"`
}

type ProductResponse struct {
	ID              int64   `json:"id"`
	CategoryId      *int64  `json:"category_id,omitempty"`
	ProductDetailId *int64  `json:"product_detail_id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Picture         *string `json:"photo_url,omitempty"`
}

type ProductDetailResponse struct {
	ID                    int64   `json:"id"`
	GenericName           string  `json:"generic_name"`
	Content               string  `json:"content"`
	Composition           string  `json:"composition"`
	Manufacturer          string  `json:"manufacturer"`
	Description           string  `json:"description"`
	ProductClassification string  `json:"product_classification"`
	ProductForm           string  `json:"product_form"`
	UnitInPack            string  `json:"unit_in_pack"`
	SellingUnit           string  `json:"selling_unit"`
	Weight                float64 `json:"weight"`
	Height                float64 `json:"height"`
	Length                float64 `json:"length"`
	Width                 float64 `json:"width"`
}

type ProductWithDetailResponse struct {
	ID              int64                 `json:"id"`
	CategoryId      *int64                `json:"category_id,omitempty"`
	ProductDetailId *int64                `json:"product_detail_id"`
	Name            string                `json:"name"`
	Slug            string                `json:"slug"`
	Picture         *string               `json:"photo_url,omitempty"`
	ProductDetail   ProductDetailResponse `json:"product_detail"`
}

func (q *GetProductsQuery) ToProductsQuery() domain.ProductsQuery {
	var page int64 = q.Page
	var sortBy string = q.SortBy
	var sortType string = q.SortType

	if q.Page == 0 || q.Limit == 0 {
		page = 1
	}

	if q.SortBy == "" {
		sortBy = domain.ProductSortById
	}

	if q.SortType == "" {
		sortType = constants.SortAsc
	}

	if sortType == constants.SortAsc {
		sortType = constants.SortDesc
	} else {
		sortType = constants.SortAsc
	}
	return domain.ProductsQuery{
		Page:      page,
		Limit:     q.Limit,
		Term:      q.Term,
		Latitude:  q.Latitude,
		Longitude: q.Longitude,
		SortBy:    sortBy,
		SortType:  sortType,
	}
}

func NewProductResponse(c domain.Product) ProductResponse {
	picture := c.Picture
	if picture == nil {
		t := constants.DefaultCategoryImageURL
		picture = &t
	}
	return ProductResponse{
		ID:              c.ID,
		Name:            c.Name,
		Slug:            c.Slug,
		CategoryId:      &c.ProductCategoryId,
		ProductDetailId: &c.ProductDetailId,
		Picture:         picture,
	}
}

func NewProductDetail(c domain.ProductDetails) ProductDetailResponse {
	return ProductDetailResponse{
		ID:                    c.ID,
		GenericName:           c.GenericName,
		Content:               c.Content,
		Composition:           c.Composition,
		Manufacturer:          c.Manufacturer,
		Description:           c.Description,
		ProductClassification: c.ProductClassification,
		ProductForm:           c.ProductForm,
		UnitInPack:            c.UnitInPack,
		SellingUnit:           c.SellingUnit,
		Weight:                c.Weight,
		Height:                c.Height,
		Length:                c.Length,
		Width:                 c.Width,
	}
}

func NewProductwithDetailResponse(c domain.Product, d domain.ProductDetails) ProductWithDetailResponse {
	picture := c.Picture
	if picture == nil {
		t := constants.DefaultCategoryImageURL
		picture = &t
	}

	return ProductWithDetailResponse{
		ID:              c.ID,
		Name:            c.Name,
		Slug:            c.Slug,
		CategoryId:      &c.ProductCategoryId,
		ProductDetailId: &c.ProductDetailId,
		Picture:         picture,
		ProductDetail:   NewProductDetail(d),
	}
}

func NewProductsResponse(products []domain.Product, pageInfo domain.PageInfo) ProductsResponse {
	res := make([]ProductResponse, len(products))
	for i := 0; i < len(products); i++ {
		res[i] = NewProductResponse(products[i])
	}
	return ProductsResponse{
		Products: res,
		PageInfo: NewPageInfoResponse(pageInfo),
	}
}
