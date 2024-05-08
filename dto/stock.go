package dto

import (
	"medichat-be/domain"
	"time"
)

type StockResponse struct {
	ID int64 `json:"id"`

	ProductID  int64 `json:"product_id"`
	PharmacyID int64 `json:"pharmacy_id"`

	Stock int `json:"stock"`
	Price int `json:"price"`
}

func NewStockResponse(s domain.Stock) StockResponse {
	return StockResponse(s)
}

type StockJoinedResponse struct {
	ID int64 `json:"id"`

	Product struct {
		ID   int64  `json:"id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"product"`
	Pharmacy struct {
		ID   int64  `json:"id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"pharmacy"`

	Stock int `json:"stock"`
	Price int `json:"price"`
}

func NewStockJoinedResponse(s domain.StockJoined) StockJoinedResponse {
	return StockJoinedResponse(s)
}

type StockCreateRequest struct {
	ProductSlug  string `json:"product_slug" binding:"required"`
	PharmacySlug string `json:"pharmacy_slug" binding:"required"`

	Stock int `json:"stock" binding:"min=0"`
	Price int `json:"price" binding:"min=0"`
}

func (r StockCreateRequest) ToDetails() domain.StockCreateDetail {
	return domain.StockCreateDetail(r)
}

type StockUpdateRequest struct {
	ID int64 `json:"id" binding:"required"`

	Stock *int `json:"stock" binding:"omitempty,min=0"`
	Price *int `json:"price" binding:"omitempty,min=0"`
}

func (r StockUpdateRequest) ToDetails() domain.StockUpdateDetail {
	return domain.StockUpdateDetail(r)
}

type StockMutationResponse struct {
	ID int64 `json:"id"`

	SourceID int64 `json:"source_id"`
	TargetID int64 `json:"target_id"`

	Method string `json:"method"`
	Status string `json:"status"`

	Amount int `json:"amount"`

	Timestamp time.Time `json:"created_at"`
}

func NewStockMutationResponse(s domain.StockMutation) StockMutationResponse {
	return StockMutationResponse(s)
}

type StockMutationJoinedResponse struct {
	ID int64 `json:"id"`

	Source struct {
		ID           int64  `json:"id"`
		PharmacyID   int64  `json:"pharmacy_id"`
		PharmacySlug string "json:\"pharmacy_slug\""
		PharmacyName string `json:"pharmacy_name"`
	} `json:"source"`
	Target struct {
		ID           int64  `json:"id"`
		PharmacyID   int64  `json:"pharmacy_id"`
		PharmacySlug string "json:\"pharmacy_slug\""
		PharmacyName string `json:"pharmacy_name"`
	} `json:"target"`
	Product struct {
		ID   int64  `json:"id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"product"`

	Method string `json:"method"`
	Status string `json:"status"`

	Amount int `json:"amount"`

	Timestamp time.Time `json:"created_at"`
}

func NewStockMutationJoinedResponse(s domain.StockMutationJoined) StockMutationJoinedResponse {
	return StockMutationJoinedResponse{
		ID: s.ID,
		Source: struct {
			ID           int64  "json:\"id\""
			PharmacyID   int64  "json:\"pharmacy_id\""
			PharmacySlug string "json:\"pharmacy_slug\""
			PharmacyName string "json:\"pharmacy_name\""
		}{
			ID:           s.Source.ID,
			PharmacyID:   s.Source.PharmacyID,
			PharmacySlug: s.Source.PharmacySlug,
			PharmacyName: s.Source.PharmacyName,
		},
		Target: struct {
			ID           int64  "json:\"id\""
			PharmacyID   int64  "json:\"pharmacy_id\""
			PharmacySlug string "json:\"pharmacy_slug\""
			PharmacyName string "json:\"pharmacy_name\""
		}{
			ID:           s.Target.ID,
			PharmacyID:   s.Target.PharmacyID,
			PharmacySlug: s.Target.PharmacySlug,
			PharmacyName: s.Target.PharmacyName,
		},
		Product: struct {
			ID   int64  "json:\"id\""
			Slug string "json:\"slug\""
			Name string "json:\"name\""
		}{
			ID:   s.Product.ID,
			Slug: s.Product.Slug,
			Name: s.Product.Name,
		},
		Method:    s.Method,
		Status:    s.Status,
		Amount:    s.Amount,
		Timestamp: s.Timestamp,
	}
}

type StockTransferRequest struct {
	SourcePharmacySlug string `json:"source_pharmacy_slug" binding:"required"`
	TargetPharmacySlug string `json:"target_pharmacy_slug" binding:"required"`
	ProductSlug        string `json:"product_slug" binding:"required"`
	Amount             int    `json:"amount" binding:"min=1"`
}

func (r StockTransferRequest) ToRequest() domain.StockTransferRequest {
	return domain.StockTransferRequest(r)
}

type StockListQuery struct {
	ProductSlug  *string `form:"product_slug"`
	ProductName  *string `form:"product_name"`
	PharmacySlug *string `form:"pharmacy_slug"`

	SortBy *string `form:"sort_by"`
	Sort   *string `form:"sort"`

	Page  *int `form:"page"`
	Limit *int `form:"limit"`
}

func (q StockListQuery) ToDetails() domain.StockListDetails {
	ret := domain.StockListDetails{
		ProductSlug:  q.ProductSlug,
		ProductName:  q.ProductName,
		PharmacySlug: q.PharmacySlug,
		SortBy:       domain.StockSortByProductName,
		SortAsc:      true,
		Page:         1,
		Limit:        10,
	}

	if q.SortBy != nil {
		ret.SortBy = *q.SortBy
	}
	if q.Sort != nil {
		ret.SortAsc = *q.Sort == "asc"
	}
	if q.Page != nil {
		ret.Page = *q.Page
	}
	if q.Limit != nil {
		ret.Limit = *q.Limit
	}

	return ret
}

type StockMutationListQuery struct {
	ProductSlug *string `form:"product_slug"`
	ProductName *string `form:"product_name"`

	SourcePharmacySlug *string `form:"source_pharmacy_slug"`
	TargetPharmacySlug *string `form:"target_pharmacy_slug"`

	Method *string `form:"method"`
	Status *string `form:"status"`

	SortBy *string `form:"sort_by"`
	Sort   *string `form:"sort"`

	Page  *int `form:"page"`
	Limit *int `form:"limit"`
}

func (q StockMutationListQuery) ToDetails() domain.StockMutationListDetails {
	ret := domain.StockMutationListDetails{
		ProductSlug:        q.ProductSlug,
		ProductName:        q.ProductName,
		SourcePharmacySlug: q.SourcePharmacySlug,
		TargetPharmacySlug: q.TargetPharmacySlug,
		Method:             q.Method,
		Status:             q.Status,
		SortBy:             "created_at",
		SortAsc:            false,
		Page:               1,
		Limit:              10,
	}

	if q.SortBy != nil {
		ret.SortBy = *q.SortBy
	}
	if q.Sort != nil {
		ret.SortAsc = *q.Sort == "asc"
	}
	if q.Page != nil {
		ret.Page = *q.Page
	}
	if q.Limit != nil {
		ret.Limit = *q.Limit
	}

	return ret
}
