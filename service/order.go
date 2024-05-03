package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
)

type orderService struct {
	dataRepository domain.DataRepository
	cloudProvider  util.CloudinaryProvider
}

type OrderServiceOpts struct {
	DataRepository domain.DataRepository
	CloudProvider  util.CloudinaryProvider
}

func NewOrderService(opts OrderServiceOpts) *orderService {
	return &orderService{
		dataRepository: opts.DataRepository,
		cloudProvider:  opts.CloudProvider,
	}
}

func (s *orderService) List(ctx context.Context, dets domain.OrderListDetails) ([]domain.Order, domain.PageInfo, error) {
	orderRepo := s.dataRepository.OrderRepository()

	page, err := orderRepo.GetPageInfo(ctx, dets)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	orders, err := orderRepo.List(ctx, dets)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	return orders, page, err
}

func (s *orderService) GetByID(ctx context.Context, id int64) (domain.Order, error) {
	orderRepo := s.dataRepository.OrderRepository()

	order, err := orderRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Order{}, apperror.Wrap(err)
	}

	items, err := orderRepo.ListItemsByOrderID(ctx, id)
	if err != nil {
		return domain.Order{}, apperror.Wrap(err)
	}

	order.Items = items

	return order, err
}

func (s *orderService) GetCartInfo(ctx context.Context, dets []domain.OrderCreateDetails) (domain.Orders, error) {
	panic("not implemented") // TODO: Implement
}

func (s *orderService) AddOrdersClosure(
	ctx context.Context,
	dets []domain.OrderCreateDetails,
) domain.AtomicFunc[domain.Orders] {
	return func(dr domain.DataRepository) (domain.Orders, error) {
		panic("not implemented") // TODO: Implement
	}
}

func (s *orderService) AddOrders(
	ctx context.Context,
	dets []domain.OrderCreateDetails,
) (domain.Orders, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.AddOrdersClosure(ctx, dets),
	)
}

func (s *orderService) SendOrderClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		panic("not implemented") // TODO: Implement
	}
}

func (s *orderService) SendOrder(
	ctx context.Context,
	id int64,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.SendOrderClosure(ctx, id),
	)
	return err
}

func (s *orderService) FinishOrderClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		panic("not implemented") // TODO: Implement
	}
}

func (s *orderService) FinishOrder(
	ctx context.Context,
	id int64,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.FinishOrderClosure(ctx, id),
	)
	return err
}

func (s *orderService) CancelOrderClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		panic("not implemented") // TODO: Implement
	}
}

func (s *orderService) CancelOrder(
	ctx context.Context,
	id int64,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.CancelOrderClosure(ctx, id),
	)
	return err
}
