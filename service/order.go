package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
	"time"
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

	_, profile, err := util.GetProfileFromContext(ctx, s.dataRepository)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	if user, ok := profile.(domain.User); ok {
		dets.UserID = &user.ID
	}

	if manager, ok := profile.(domain.PharmacyManager); ok {
		dets.PharmacyManagerID = &manager.ID
	}

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

	_, profile, err := util.GetProfileFromContext(ctx, s.dataRepository)
	if err != nil {
		return domain.Order{}, apperror.Wrap(err)
	}

	order, err := orderRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Order{}, apperror.Wrap(err)
	}

	items, err := orderRepo.ListItemsByOrderID(ctx, id)
	if err != nil {
		return domain.Order{}, apperror.Wrap(err)
	}

	if user, ok := profile.(domain.User); ok {
		if order.User.ID != user.ID {
			return domain.Order{}, apperror.NewForbidden(nil)
		}
	}

	if manager, ok := profile.(domain.PharmacyManager); ok {
		if order.Pharmacy.ManagerID != manager.ID {
			return domain.Order{}, apperror.NewForbidden(nil)
		}
	}

	order.Items = items

	return order, err
}

func (s *orderService) getOrders(dr domain.DataRepository, ctx context.Context, dets []domain.OrderCreateDetails) (domain.Orders, error) {
	productRepo := dr.ProductRepository()
	pharmacyRepo := dr.PharmacyRepository()
	userRepo := dr.UserRepository()
	stockRepo := dr.StockRepository()
	shipmentRepo := dr.ShipmentMethodRepository()

	accountID, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
		return domain.Orders{}, apperror.Wrap(err)
	}

	user, err := userRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return domain.Orders{}, apperror.Wrap(err)
	}

	orders := domain.Orders{
		Orders: []domain.Order{},
		Total:  0,
	}

	var orderID int64 = 1
	var itemID int64 = 1

	for _, det := range dets {
		pharmacy, err := pharmacyRepo.GetBySlug(ctx, det.PharmacySlug)
		if err != nil {
			return domain.Orders{}, apperror.Wrap(err)
		}

		shipment, err := shipmentRepo.GetShipmentMethodById(ctx, det.ShipmentMethodID)
		if err != nil {
			return domain.Orders{}, apperror.Wrap(err)
		}

		order := domain.Order{
			ID: orderID,
			User: struct {
				ID   int64
				Name string
			}{
				ID:   user.ID,
				Name: user.Account.Name,
			},
			Pharmacy: struct {
				ID        int64
				Slug      string
				Name      string
				ManagerID int64
			}{
				ID:        pharmacy.ID,
				Slug:      pharmacy.Slug,
				Name:      pharmacy.Name,
				ManagerID: pharmacy.ManagerID,
			},
			ShipmentMethod: shipment,
			Address:        det.Address,
			Coordinate:     det.Coordinate,
			NItems:         0,
			Subtotal:       0,
			ShipmentFee:    0,
			Total:          0,
			Status:         domain.OrderStatusWaitingPayment,
			OrderedAt:      time.Now(),
			FinishedAt:     nil,
			Items:          []domain.OrderItem{},
		}

		for _, it := range det.Items {
			product, err := productRepo.GetBySlug(ctx, it.ProductSlug)
			if err != nil {
				return domain.Orders{}, apperror.Wrap(err)
			}

			stock, err := stockRepo.GetByPharmacyAndProduct(ctx, pharmacy.ID, product.ID)
			if err != nil {
				return domain.Orders{}, apperror.Wrap(err)
			}
			if stock.Stock < it.Amount {
				return domain.Orders{}, apperror.NewStockNotEnough(nil)
			}
			price := stock.Price

			picture := ""
			if product.Picture != nil {
				picture = *product.Picture
			}

			order.Items = append(order.Items, domain.OrderItem{
				ID:      itemID,
				OrderID: order.ID,
				Product: struct {
					ID       int64
					Slug     string
					Name     string
					PhotoURL string
				}{
					ID:       product.ID,
					Slug:     product.Slug,
					Name:     product.Name,
					PhotoURL: picture,
				},
				Price:  price,
				Amount: it.Amount,
			})

			order.Subtotal += price * it.Amount
			order.NItems += it.Amount
			itemID++
		}

		distance, err := dr.GetDistance(ctx, det.Coordinate, pharmacy.Coordinate)
		if err != nil {
			return domain.Orders{}, apperror.Wrap(err)
		}
		var shipmentFee int
		switch det.ShipmentMethodID {
		case domain.ShipmentOfficialInstantID:
			shipmentFee = int(distance / 1000 * domain.ShipmentOfficialInstantFee)
		case domain.ShipmentOfficialSameDayID:
			shipmentFee = int(distance / 1000 * domain.ShipmentOfficialSameDayFee)
		default:
			shipmentFee = int(distance / 1000 * domain.ShipmentOfficialSameDayFee)
		}

		order.ShipmentFee = shipmentFee
		order.Total = order.Subtotal + order.ShipmentFee

		orders.Total += order.Total
		orders.Orders = append(orders.Orders, order)

		orderID++
	}

	return orders, nil
}

func (s *orderService) GetCartInfo(ctx context.Context, dets []domain.OrderCreateDetails) (domain.Orders, error) {
	return s.getOrders(s.dataRepository, ctx, dets)
}

func (s *orderService) AddOrdersClosure(
	ctx context.Context,
	dets []domain.OrderCreateDetails,
) domain.AtomicFunc[domain.Orders] {
	return func(dr domain.DataRepository) (domain.Orders, error) {
		orderRepo := dr.OrderRepository()
		paymentRepo := dr.PaymentRepository()

		if len(dets) == 0 {
			return domain.Orders{}, nil
		}

		orders, err := s.getOrders(dr, ctx, dets)
		if err != nil {
			return domain.Orders{}, apperror.Wrap(err)
		}

		payment := domain.Payment{
			InvoiceNumber: util.GenerateInvoiceNumber(),
			User:          orders.Orders[0].User,
			FileURL:       nil,
			IsConfirmed:   false,
			Amount:        orders.Total,
		}

		payment, err = paymentRepo.Add(ctx, payment)
		if err != nil {
			return domain.Orders{}, apperror.Wrap(err)
		}

		for i := range orders.Orders {
			order := &orders.Orders[i]

			order.Payment.ID = payment.ID
			order.Payment.InvoiceNumber = payment.InvoiceNumber

			newOrder, err := orderRepo.Add(ctx, *order)
			if err != nil {
				return domain.Orders{}, apperror.Wrap(err)
			}

			order.ID = newOrder.ID

			for j := range order.Items {
				item := &order.Items[j]

				item.OrderID = order.ID

				newItem, err := orderRepo.AddItem(ctx, *item)
				if err != nil {
					return domain.Orders{}, apperror.Wrap(err)
				}

				item.ID = newItem.ID
			}
		}

		return orders, err
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
		orderRepo := dr.OrderRepository()
		managerRepo := dr.PharmacyManagerRepository()
		pharmacyRepo := dr.PharmacyRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		manager, err := managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		order, err := orderRepo.GetByIDAndLock(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if order.Status != domain.OrderStatusProcessing {
			return nil, apperror.NewAppError(
				apperror.CodeBadRequest,
				"order status should be processing",
				nil,
			)
		}

		pharmacy, err := pharmacyRepo.GetByID(ctx, order.Pharmacy.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		if pharmacy.ManagerID != manager.ID {
			return nil, apperror.NewForbidden(nil)
		}

		items, err := orderRepo.ListItemsByOrderID(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		for _, item := range items {
			err = s.updateStockForOrderItem(ctx, dr, order.Pharmacy.ID, item)
			if err != nil {
				return nil, apperror.Wrap(err)
			}
		}

		err = orderRepo.UpdateStatusByID(ctx, id, domain.OrderStatusSent)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
	}
}

func (s *orderService) updateStockForOrderItem(ctx context.Context, dr domain.DataRepository, pharmacyID int64, item domain.OrderItem) error {
	stockRepo := dr.StockRepository()

	stock, err := stockRepo.GetByPharmacyAndProduct(ctx, pharmacyID, item.Product.ID)
	if err != nil {
		return err
	}

	stock, err = stockRepo.GetByIDAndLock(ctx, stock.ID)
	if err != nil {
		return err
	}

	if stock.Stock < item.Amount {
		source, err := stockRepo.GetNearestStockWithProduct(ctx, stock.PharmacyID, stock.ProductID, item.Amount)
		if err != nil {
			return err
		}

		mutation := domain.StockMutation{
			SourceID: source.ID,
			TargetID: stock.ID,
			Method:   domain.StockMutationAutomatic,
			Status:   domain.StockMutationStatusApproved,
			Amount:   item.Amount,
		}

		mutation, err = stockRepo.AddMutation(ctx, mutation)
		if err != nil {
			return err
		}

		source, stock, err = transferStock(dr, ctx, source.ID, stock.ID, item.Amount)
		if err != nil {
			return err
		}
	}

	stock.Stock -= item.Amount
	stock, err = stockRepo.Update(ctx, stock)
	if err != nil {
		return err
	}

	return nil
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
		orderRepo := dr.OrderRepository()
		userRepo := dr.UserRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		order, err := orderRepo.GetByIDAndLock(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if order.User.ID != user.ID {
			return nil, apperror.NewForbidden(nil)
		}
		if order.Status != domain.OrderStatusSent {
			return nil, apperror.NewAppError(
				apperror.CodeBadRequest,
				"order status should be sent",
				nil,
			)
		}

		err = orderRepo.UpdateStatusByID(ctx, id, domain.OrderStatusFinished)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
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
		orderRepo := dr.OrderRepository()
		userRepo := dr.UserRepository()
		accountRepo := dr.AccountRepository()
		managerRepo := dr.PharmacyManagerRepository()
		pharmacyRepo := dr.PharmacyRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByID(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		order, err := orderRepo.GetByIDAndLock(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		if account.Role == domain.AccountRoleUser {
			user, err := userRepo.GetByAccountID(ctx, accountID)
			if err != nil {
				return nil, apperror.Wrap(err)
			}

			if order.User.ID != user.ID {
				return nil, apperror.NewForbidden(nil)
			}
		}

		if account.Role == domain.AccountRolePharmacyManager {
			manager, err := managerRepo.GetByAccountID(ctx, accountID)
			if err != nil {
				return nil, apperror.Wrap(err)
			}

			pharmacy, err := pharmacyRepo.GetByID(ctx, order.Pharmacy.ID)
			if err != nil {
				return nil, apperror.Wrap(err)
			}

			if pharmacy.ManagerID != manager.ID {
				return nil, apperror.NewForbidden(nil)
			}
		}

		if order.Status != domain.OrderStatusWaitingConfirmation &&
			order.Status != domain.OrderStatusWaitingPayment &&
			order.Status != domain.OrderStatusProcessing {
			return nil, apperror.NewAppError(
				apperror.CodeBadRequest,
				"order cannot be cancelled",
				nil,
			)
		}

		err = orderRepo.UpdateStatusByID(ctx, id, domain.OrderStatusCancelled)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
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
