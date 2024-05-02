package postgres

import (
	"context"
	"medichat-be/domain"
)

type orderRepository struct {
	querier Querier
}

func (r *orderRepository) GetPageInfo(ctx context.Context, dets domain.OrderListDetails) (domain.PageInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) List(ctx context.Context, dets domain.OrderListDetails) ([]domain.Order, error) {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) GetByID(ctx context.Context, id int64) (domain.Order, error) {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) GetByIDAndLock(ctx context.Context, id int64) (domain.Order, error) {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) Add(ctx context.Context, order domain.Order) (domain.Order, error) {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) UpdateStatusByID(ctx context.Context, id int64, status string) error {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) ListItemsByOrderID(ctx context.Context, id int64) ([]domain.OrderItem, error) {
	panic("not implemented") // TODO: Implement
}

func (r *orderRepository) AddItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	panic("not implemented") // TODO: Implement
}
