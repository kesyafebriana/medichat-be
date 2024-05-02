package domain

import (
	"context"
	"mime/multipart"
)

type Payment struct {
	ID            int64
	InvoiceNumber string
	FileURL       string
	IsConfirmed   bool
	Amount        int
}

type PaymentListDetails struct {
	IsConfirmed *bool
}

type PaymentRepository interface {
	GetPageInfo(ctx context.Context, dets PaymentListDetails) (PageInfo, error)
	List(ctx context.Context, dets PaymentListDetails) ([]Payment, error)
	GetByID(ctx context.Context, id int64) (Payment, error)
	GetByInvoiceNumber(ctx context.Context, num string) (Payment, error)
	Add(ctx context.Context, p Payment) (Payment, error)
	Update(ctx context.Context, p Payment) (Payment, error)
}

type PaymentService interface {
	List(ctx context.Context, id int64) ([]Payment, PageInfo, error)
	GetByID(ctx context.Context, id int64) (Payment, error)
	GetByInvoiceNumber(ctx context.Context, num string) (Payment, error)

	UploadPayment(ctx context.Context, num string, file multipart.File) error
	ConfirmPayment(ctx context.Context, num string)
}
