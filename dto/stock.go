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
		Name string `json:"name"`
	} `json:"pharmacy"`

	Stock int `json:"stock"`
	Price int `json:"price"`
}

func NewStockJoinedResponse(s domain.StockJoined) StockJoinedResponse {
	return StockJoinedResponse(s)
}

type StockCreateRequest struct {
	ProductSlug string `json:"product_slug" binding:"required"`
	PharmacyID  int64  `json:"pharmacy_id" binding:"required"`

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

	Timestamp time.Time `json:"timestamp"`
}

func NewStockMutationResponse(s domain.StockMutation) StockMutationResponse {
	return StockMutationResponse(s)
}

type StockMutationJoinedResponse struct {
	ID int64 `json:"id"`

	Source struct {
		ID           int64  `json:"id"`
		PharmacyID   int64  `json:"pharmacy_id"`
		PharmacyName string `json:"pharmacy_name"`
	} `json:"source"`
	Target struct {
		ID           int64  `json:"id"`
		PharmacyID   int64  `json:"pharmacy_id"`
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

	Timestamp time.Time `json:"timestamp"`
}

func NewStockMutationJoinedResponse(s domain.StockMutationJoined) StockMutationJoinedResponse {
	return StockMutationJoinedResponse{
		ID: s.ID,
		Source: struct {
			ID           int64  "json:\"id\""
			PharmacyID   int64  "json:\"pharmacy_id\""
			PharmacyName string "json:\"pharmacy_name\""
		}{
			ID:           s.Source.ID,
			PharmacyID:   s.Source.PharmacyID,
			PharmacyName: s.Source.PharmacyName,
		},
		Target: struct {
			ID           int64  "json:\"id\""
			PharmacyID   int64  "json:\"pharmacy_id\""
			PharmacyName string "json:\"pharmacy_name\""
		}{
			ID:           s.Target.ID,
			PharmacyID:   s.Target.PharmacyID,
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
	SourcePharmacyID int64  `json:"source_pharmacy_id"`
	TargetPharmacyID int64  `json:"target_pharmacy_id"`
	ProductSlug      string `json:"product_slug"`
	Amount           int    `json:"amount"`
}

func (r StockTransferRequest) ToRequest() domain.StockTransferRequest {
	return domain.StockTransferRequest(r)
}

type StockListQuery struct {
	ProductSlug *string `form:"product_slug"`
	ProductName *string `form:"product_name"`
	PharmacyID  *int64  `form:"pharmacy_id"`

	SortBy  string `form:"sort_by"`
	SortAsc bool   `form:"sort_asc"`

	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (q StockListQuery) ToDetails() domain.StockListDetails {
	return domain.StockListDetails(q)
}

type StockMutationListQuery struct {
	ProductSlug *string `form:"product_slug"`
	ProductName *string `form:"product_name"`

	SourcePharmacyID *int64 `form:"source_pharmacy_id"`
	TargetPharmacyID *int64 `form:"target_pharmacy_id"`

	Method *string `form:"method"`
	Status *string `form:"status"`

	SortBy  string `form:"sort_by"`
	SortAsc bool   `form:"sort_asc"`

	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (q StockMutationListQuery) ToDetails() domain.StockMutationListDetails {
	return domain.StockMutationListDetails(q)
}
