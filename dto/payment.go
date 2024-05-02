package dto

import (
	"medichat-be/domain"
	"mime/multipart"
)

type PaymentResponse struct {
	ID            int64  `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	User          struct {
		ID   int64  `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"user"`
	FileURL     *string `json:"file_url"`
	IsConfirmed bool    `json:"is_confirmed"`
	Amount      int     `json:"amount"`
}

func NewPaymentResponse(p domain.Payment) PaymentResponse {
	return PaymentResponse(p)
}

type PaymentListQuery struct {
	IsConfirmed *bool `form:"is_confirmed"`

	Page  *int `form:"page" binding:"omitempty,min=1"`
	Limit *int `form:"limit" binding:"omitempty,min=1"`
}

func (q PaymentListQuery) ToDetails() domain.PaymentListDetails {
	ret := domain.PaymentListDetails{
		IsConfirmed: q.IsConfirmed,
		UserID:      nil,
		Page:        1,
		Limit:       10,
	}

	if q.Page != nil {
		ret.Page = *q.Page
	}
	if q.Limit != nil {
		ret.Limit = *q.Limit
	}

	return ret
}

type PaymentInvoiceNumberURI struct {
	InvoiceNumber string `json:"invoice_number" binding:"required"`
}

type PaymentUploadRequest = MultipartForm[
	struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	},
	struct{},
]
