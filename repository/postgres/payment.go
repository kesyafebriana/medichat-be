package postgres

import (
	"context"
	"medichat-be/domain"
)

type paymentRepository struct {
	querier Querier
}

func (r *paymentRepository) GetPageInfo(ctx context.Context, dets domain.PaymentListDetails) (domain.PageInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (r *paymentRepository) List(ctx context.Context, dets domain.PaymentListDetails) ([]domain.Payment, error) {
	panic("not implemented") // TODO: Implement
}

func (r *paymentRepository) GetByID(ctx context.Context, id int64) (domain.Payment, error) {
	panic("not implemented") // TODO: Implement
}

func (r *paymentRepository) GetByInvoiceNumber(ctx context.Context, num string) (domain.Payment, error) {
	panic("not implemented") // TODO: Implement
}

func (r *paymentRepository) Add(ctx context.Context, p domain.Payment) (domain.Payment, error) {
	panic("not implemented") // TODO: Implement
}

func (r *paymentRepository) Update(ctx context.Context, p domain.Payment) (domain.Payment, error) {
	panic("not implemented") // TODO: Implement
}
