package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
)

type stockService struct {
	dataRepository domain.DataRepository
}

type StockServiceOpts struct {
	DataRepository domain.DataRepository
}

func NewStockService(opts StockServiceOpts) *stockService {
	return &stockService{
		dataRepository: opts.DataRepository,
	}
}

func (s *stockService) GetByID(
	ctx context.Context,
	id int64,
) (domain.Stock, error) {
	stockRepo := s.dataRepository.StockRepository()
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	_, prof, err := util.GetProfileFromContext(ctx, s.dataRepository)
	if err != nil {
		return domain.Stock{}, apperror.Wrap(err)
	}

	stock, err := stockRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Stock{}, apperror.Wrap(err)
	}

	if manager, ok := prof.(domain.PharmacyManager); ok {
		pharmacy, err := pharmacyRepo.GetByID(ctx, id)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}
		if pharmacy.ManagerID != manager.ID {
			return domain.Stock{}, apperror.NewForbidden(nil)
		}
	}

	return stock, nil
}

func (s *stockService) List(
	ctx context.Context,
	det domain.StockListDetails,
) ([]domain.StockJoined, domain.PageInfo, error) {
	stockRepo := s.dataRepository.StockRepository()

	_, prof, err := util.GetProfileFromContext(ctx, s.dataRepository)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	if manager, ok := prof.(domain.PharmacyManager); ok {
		det.ManagerID = &manager.ID
	}

	pageInfo, err := stockRepo.GetPageInfo(ctx, det)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	stocks, err := stockRepo.List(ctx, det)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	return stocks, pageInfo, nil
}

func (s *stockService) AddClosure(
	ctx context.Context,
	det domain.StockCreateDetail,
) domain.AtomicFunc[domain.Stock] {
	return func(dr domain.DataRepository) (domain.Stock, error) {
		stockRepo := dr.StockRepository()
		productRepo := dr.ProductRepository()
		pharmacyRepo := dr.PharmacyRepository()
		managerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		manager, err := managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		product, err := productRepo.GetBySlug(ctx, det.ProductSlug)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		pharmacy, err := pharmacyRepo.GetBySlug(ctx, det.PharmacySlug)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		if pharmacy.ManagerID != manager.ID {
			return domain.Stock{}, apperror.NewForbidden(nil)
		}

		stock := domain.Stock{
			ProductID:  product.ID,
			PharmacyID: pharmacy.ID,
			Stock:      det.Stock,
			Price:      det.Price,
		}

		stock, err = stockRepo.Add(ctx, stock)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		return stock, nil
	}
}

func (s *stockService) Add(
	ctx context.Context,
	stock domain.StockCreateDetail,
) (domain.Stock, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.AddClosure(ctx, stock),
	)
}

func (s *stockService) UpdateClosure(
	ctx context.Context,
	det domain.StockUpdateDetail,
) domain.AtomicFunc[domain.Stock] {
	return func(dr domain.DataRepository) (domain.Stock, error) {
		stockRepo := dr.StockRepository()
		managerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		_, err = managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		stock, err := stockRepo.GetByIDAndLock(ctx, det.ID)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		if det.Stock != nil {
			stock.Stock = *det.Stock
		}
		if det.Price != nil {
			stock.Price = *det.Price
		}

		stock, err = stockRepo.Update(ctx, stock)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		return stock, nil
	}
}

func (s *stockService) Update(
	ctx context.Context,
	det domain.StockUpdateDetail,
) (domain.Stock, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.UpdateClosure(ctx, det),
	)
}

func (s *stockService) DeleteClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		stockRepo := dr.StockRepository()
		managerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		_, err = managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		_, err = stockRepo.GetByIDAndLock(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = stockRepo.SoftDeleteByID(ctx, id)
		if err != nil {
			return domain.Stock{}, apperror.Wrap(err)
		}

		return nil, nil
	}
}

func (s *stockService) DeleteByID(
	ctx context.Context,
	id int64,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.DeleteClosure(ctx, id),
	)

	return err
}

func (s *stockService) GetMutationByID(
	ctx context.Context,
	id int64,
) (domain.StockMutation, error) {
	stockRepo := s.dataRepository.StockRepository()

	mut, err := stockRepo.GetMutationByID(ctx, id)
	if err != nil {
		return domain.StockMutation{}, apperror.Wrap(err)
	}

	return mut, nil
}

func (s *stockService) ListMutations(
	ctx context.Context,
	det domain.StockMutationListDetails,
) ([]domain.StockMutationJoined, domain.PageInfo, error) {
	stockRepo := s.dataRepository.StockRepository()

	_, prof, err := util.GetProfileFromContext(ctx, s.dataRepository)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	if manager, ok := prof.(domain.PharmacyManager); ok {
		det.ManagerID = &manager.ID
	}

	pageInfo, err := stockRepo.GetMutationPageInfo(ctx, det)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	muts, err := stockRepo.ListMutations(ctx, det)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	return muts, pageInfo, nil
}

func (s *stockService) RequestStockTransferClosure(
	ctx context.Context,
	req domain.StockTransferRequest,
) domain.AtomicFunc[domain.StockMutation] {
	return func(dr domain.DataRepository) (domain.StockMutation, error) {
		stockRepo := dr.StockRepository()
		productRepo := dr.ProductRepository()
		pharmacyRepo := dr.PharmacyRepository()
		managerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		manager, err := managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		if req.SourcePharmacySlug == req.TargetPharmacySlug {
			return domain.StockMutation{}, apperror.NewTransferSameStock(nil)
		}

		product, err := productRepo.GetBySlug(ctx, req.ProductSlug)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		sourcePharma, err := pharmacyRepo.GetBySlug(ctx, req.SourcePharmacySlug)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		targetPharma, err := pharmacyRepo.GetBySlug(ctx, req.TargetPharmacySlug)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		if targetPharma.ManagerID != manager.ID {
			return domain.StockMutation{}, apperror.NewForbidden(nil)
		}

		source, err := stockRepo.GetByPharmacyAndProduct(ctx, sourcePharma.ID, product.ID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		target, err := stockRepo.GetByPharmacyAndProduct(ctx, targetPharma.ID, product.ID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		mut := domain.StockMutation{
			SourceID: source.ID,
			TargetID: target.ID,
			Method:   domain.StockMutationManual,
			Status:   domain.StockMutationStatusPending,
			Amount:   req.Amount,
		}

		mut, err = stockRepo.AddMutation(ctx, mut)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		return mut, nil
	}
}

func (s *stockService) RequestStockTransfer(
	ctx context.Context,
	req domain.StockTransferRequest,
) (domain.StockMutation, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.RequestStockTransferClosure(ctx, req),
	)
}

func (s *stockService) ApproveStockTransferClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[domain.StockMutation] {
	return func(dr domain.DataRepository) (domain.StockMutation, error) {
		stockRepo := dr.StockRepository()
		pharmacyRepo := dr.PharmacyRepository()
		managerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		manager, err := managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		mut, err := stockRepo.GetMutationByIDAndLock(ctx, id)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		sourceStock, err := stockRepo.GetByID(ctx, mut.SourceID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}
		sourcePharma, err := pharmacyRepo.GetByID(ctx, sourceStock.PharmacyID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}
		if sourcePharma.ManagerID != manager.ID {
			return domain.StockMutation{}, apperror.NewForbidden(err)
		}

		if mut.Status != domain.StockMutationStatusPending {
			return domain.StockMutation{}, apperror.NewNotPending(nil)
		}

		mut.Status = domain.StockMutationStatusApproved

		mut, err = stockRepo.UpdateMutation(ctx, mut)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		_, _, err = transferStock(dr, ctx, mut.SourceID, mut.TargetID, mut.Amount)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		return mut, nil
	}
}

func (s *stockService) ApproveStockTransfer(
	ctx context.Context,
	id int64,
) (domain.StockMutation, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.ApproveStockTransferClosure(ctx, id),
	)
}

func (s *stockService) CancelStockTransferClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[domain.StockMutation] {
	return func(dr domain.DataRepository) (domain.StockMutation, error) {
		stockRepo := dr.StockRepository()
		pharmacyRepo := dr.PharmacyRepository()
		managerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		manager, err := managerRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		mut, err := stockRepo.GetMutationByIDAndLock(ctx, id)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		sourceStock, err := stockRepo.GetByID(ctx, mut.SourceID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}
		sourcePharma, err := pharmacyRepo.GetByID(ctx, sourceStock.PharmacyID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		targetStock, err := stockRepo.GetByID(ctx, mut.TargetID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}
		targetPharma, err := pharmacyRepo.GetByID(ctx, targetStock.PharmacyID)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		if sourcePharma.ManagerID != manager.ID && targetPharma.ManagerID != manager.ID {
			return domain.StockMutation{}, apperror.NewForbidden(err)
		}

		if mut.Status != domain.StockMutationStatusPending {
			return domain.StockMutation{}, apperror.NewNotPending(nil)
		}

		mut.Status = domain.StockMutationStatusCancelled

		mut, err = stockRepo.UpdateMutation(ctx, mut)
		if err != nil {
			return domain.StockMutation{}, apperror.Wrap(err)
		}

		return mut, nil
	}
}

func (s *stockService) CancelStockTransfer(
	ctx context.Context,
	id int64,
) (domain.StockMutation, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.ApproveStockTransferClosure(ctx, id),
	)
}

func transferStock(
	dr domain.DataRepository,
	ctx context.Context,
	sourceID int64,
	targetID int64,
	amount int,
) (domain.Stock, domain.Stock, error) {
	stockRepo := dr.StockRepository()

	source, err := stockRepo.GetByIDAndLock(ctx, sourceID)
	if err != nil {
		return domain.Stock{}, domain.Stock{}, apperror.Wrap(err)
	}

	target, err := stockRepo.GetByIDAndLock(ctx, targetID)
	if err != nil {
		return domain.Stock{}, domain.Stock{}, apperror.Wrap(err)
	}

	if source.Stock < amount {
		return domain.Stock{}, domain.Stock{}, apperror.NewStockNotEnough(nil)
	}

	source.Stock -= amount
	target.Stock += amount

	source, err = stockRepo.Update(ctx, source)
	if err != nil {
		return domain.Stock{}, domain.Stock{}, apperror.Wrap(err)
	}

	target, err = stockRepo.Update(ctx, target)
	if err != nil {
		return domain.Stock{}, domain.Stock{}, apperror.Wrap(err)
	}

	return source, target, nil
}
