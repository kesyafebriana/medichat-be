package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
)

type pharmacyService struct {
	dataRepository domain.DataRepository
}

type PharmacyServiceOpts struct {
	DataRepository domain.DataRepository
}

func NewPharmacyService(opts PharmacyServiceOpts) *pharmacyService {
	return &pharmacyService{
		dataRepository: opts.DataRepository,
	}
}

func (s *pharmacyService) GetPharmacies(ctx context.Context, query domain.PharmaciesQuery) ([]domain.Pharmacy, domain.PageInfo, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	shipmentRepo := s.dataRepository.ShipmentMethodRepository()

	p, err := pharmacyRepo.GetPharmacies(ctx, query)
	if err != nil {
		return []domain.Pharmacy{}, domain.PageInfo{}, apperror.Wrap(err)
	}

	for i, v := range p {
		o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, v.ID)
		if err != nil {
			return []domain.Pharmacy{}, domain.PageInfo{}, apperror.Wrap(err)
		}

		sh, err := pharmacyRepo.GetShipmentMethodsByPharmacyId(ctx, v.ID)
		if err != nil {
			return []domain.Pharmacy{}, domain.PageInfo{}, apperror.Wrap(err)
		}

		for i, x := range sh {
			shDetail, err := shipmentRepo.GetShipmentMethodById(ctx, x.ShipmentMethodID)
			if err != nil {
				return []domain.Pharmacy{}, domain.PageInfo{}, apperror.Wrap(err)
			}

			sh[i].Name = &shDetail.Name
		}

		p[i].PharmacyOperations = o
		p[i].PharmacyShipmentMethods = sh
	}

	pageInfo, err := pharmacyRepo.GetPageInfo(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo.ItemsPerPage = int(query.Limit)
	if query.Limit == 0 {
		pageInfo.ItemsPerPage = len(p)
	}

	if pageInfo.ItemsPerPage == 0 {
		pageInfo.PageCount = 0
	} else {
		pageInfo.PageCount = (int(pageInfo.ItemCount) + pageInfo.ItemsPerPage - 1) / pageInfo.ItemsPerPage
	}

	return p, pageInfo, nil
}

func (s *pharmacyService) GetPharmacyBySlug(ctx context.Context, slug string) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	shipmentRepo := s.dataRepository.ShipmentMethodRepository()

	p, err := pharmacyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, p.ID)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	sh, err := pharmacyRepo.GetShipmentMethodsByPharmacyId(ctx, p.ID)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	for i, v := range sh {
		shDetail, err := shipmentRepo.GetShipmentMethodById(ctx, v.ID)
		if err != nil {
			return domain.Pharmacy{}, apperror.Wrap(err)
		}

		sh[i].Name = &shDetail.Name
	}

	p.PharmacyOperations = o
	p.PharmacyShipmentMethods = sh

	return p, nil
}

func (s *pharmacyService) CreatePharmacy(ctx context.Context, pharmacy domain.PharmacyCreateDetails) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	shipmentRepo := s.dataRepository.ShipmentMethodRepository()
	var pharmacyOperations []domain.PharmacyOperations
	var PharmacyShipmentMethods []domain.PharmacyShipmentMethods

	p, err := pharmacyRepo.Add(ctx, pharmacy)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	for _, v := range pharmacy.PharmacyOperations {
		v.PharmacyID = p.ID

		o, err := pharmacyRepo.AddOperation(ctx, v)
		if err != nil {
			return domain.Pharmacy{}, apperror.Wrap(err)
		}

		pharmacyOperations = append(pharmacyOperations, o)
	}

	for _, v := range pharmacy.PharmacyShipmentMethods {
		v.PharmacyID = p.ID

		sh, err := pharmacyRepo.AddShipmentMethod(ctx, v)
		if err != nil {
			return domain.Pharmacy{}, apperror.Wrap(err)
		}

		shDetail, err := shipmentRepo.GetShipmentMethodById(ctx, sh.ShipmentMethodID)
		if err != nil {
			return domain.Pharmacy{}, apperror.Wrap(err)
		}

		sh.Name = &shDetail.Name
		PharmacyShipmentMethods = append(PharmacyShipmentMethods, sh)
	}

	p.PharmacyOperations = pharmacyOperations
	p.PharmacyShipmentMethods = PharmacyShipmentMethods

	return p, nil
}

func (s *pharmacyService) UpdatePharmacy(ctx context.Context, pharmacy domain.PharmacyUpdateDetails) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	shipmentRepo := s.dataRepository.ShipmentMethodRepository()

	p, err := pharmacyRepo.Update(ctx, pharmacy)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, p.ID)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	sh, err := pharmacyRepo.GetShipmentMethodsByPharmacyId(ctx, p.ID)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	for i, v := range sh {
		shDetail, err := shipmentRepo.GetShipmentMethodById(ctx, v.ID)
		if err != nil {
			return domain.Pharmacy{}, apperror.Wrap(err)
		}

		sh[i].Name = &shDetail.Name
	}

	p.PharmacyOperations = o
	p.PharmacyShipmentMethods = sh

	return p, nil
}

func (s *pharmacyService) DeletePharmacyBySlug(ctx context.Context, slug string) error {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	_, err := pharmacyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return apperror.Wrap(err)
	}

	err = pharmacyRepo.SoftDeleteBySlug(ctx, slug)
	if err != nil {
		return apperror.Wrap(err)
	}

	return nil
}

func (s *pharmacyService) GetOperationsBySlug(ctx context.Context, slug string) ([]domain.PharmacyOperations, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	p, err := pharmacyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return []domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, p.ID)
	if err != nil {
		return []domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	return o, nil
}

func (s *pharmacyService) UpdateOperations(ctx context.Context, pharmacyOperations []domain.PharmacyOperationsUpdateDetails) ([]domain.PharmacyOperations, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	var res []domain.PharmacyOperations

	p, err := pharmacyRepo.GetBySlug(ctx, pharmacyOperations[0].Slug)
	if err != nil {
		return []domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyIdAndLock(ctx, p.ID)
	if err != nil {
		return []domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	oDays := make(map[string]bool)
	oIds := make(map[string]int64)
	for _, x := range o {
		oDays[x.Day] = true
		oIds[x.Day] = x.ID
	}

	for _, v := range pharmacyOperations {
		v.PharmacyId = p.ID

		if oDays[v.Day] {
			v.ID = oIds[v.Day]
			oUpdate, err := pharmacyRepo.UpdateOperation(ctx, v)
			if err != nil {
				return []domain.PharmacyOperations{}, apperror.Wrap(err)
			}

			res = append(res, oUpdate)
			continue
		}

		oCreate, err := pharmacyRepo.AddOperation(ctx, domain.PharmacyOperationCreateDetails{
			PharmacyID: v.PharmacyId,
			Day:        v.Day,
			StartTime:  v.StartTime,
			EndTime:    v.EndTime,
		})
		if err != nil {
			return []domain.PharmacyOperations{}, apperror.Wrap(err)
		}

		res = append(res, oCreate)
	}

	for _, x := range o {
		found := false
		for _, v := range pharmacyOperations {
			if x.Day == v.Day {
				found = true
				break
			}
		}

		if !found {
			err := pharmacyRepo.SoftDeleteOperationByID(ctx, x.ID)
			if err != nil {
				return []domain.PharmacyOperations{}, apperror.Wrap(err)
			}
		}
	}

	return res, nil
}

func (s *pharmacyService) GetShipmentMethodBySlug(ctx context.Context, slug string) ([]domain.PharmacyShipmentMethods, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	shipmentRepo := s.dataRepository.ShipmentMethodRepository()

	p, err := pharmacyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
	}

	sh, err := pharmacyRepo.GetShipmentMethodsByPharmacyId(ctx, p.ID)
	if err != nil {
		return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
	}

	for i, v := range sh {
		shDetail, err := shipmentRepo.GetShipmentMethodById(ctx, v.ShipmentMethodID)
		if err != nil {
			return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
		}

		sh[i].Name = &shDetail.Name
	}

	return sh, nil
}

func (s *pharmacyService) UpdateShipmentMethod(ctx context.Context, shipmentMethods []domain.PharmacyShipmentMethodsUpdateDetails) ([]domain.PharmacyShipmentMethods, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	shipmentRepo := s.dataRepository.ShipmentMethodRepository()
	var res []domain.PharmacyShipmentMethods

	p, err := pharmacyRepo.GetBySlug(ctx, shipmentMethods[0].Slug)
	if err != nil {
		return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
	}

	sh, err := pharmacyRepo.GetShipmentMethodsByPharmacyIdAndLock(ctx, p.ID)
	if err != nil {
		return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
	}

	for _, v := range shipmentMethods {
		v.PharmacyID = p.ID
		found := false

		for _, x := range sh {
			if x.ShipmentMethodID == v.ShipmentMethodID {
				found = true
				break
			}
		}

		if !found {
			newSh, err := pharmacyRepo.AddShipmentMethod(ctx, domain.PharmacyShipmentMethodsCreateDetails{
				PharmacyID:       v.PharmacyID,
				ShipmentMethodID: v.ShipmentMethodID,
			})
			if err != nil {
				return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
			}
			res = append(res, newSh)
		}
	}

	for _, x := range sh {
		found := false
		for _, v := range shipmentMethods {
			if x.ShipmentMethodID == v.ShipmentMethodID {
				found = true
				break
			}
		}

		if found {
			res = append(res, domain.PharmacyShipmentMethods{
				ID:               x.ID,
				PharmacyID:       x.PharmacyID,
				ShipmentMethodID: x.ShipmentMethodID,
			})
		}

		if !found {
			err := pharmacyRepo.SoftDeleteShipmentMethodByID(ctx, x.ID)
			if err != nil {
				return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
			}
		}
	}

	for i, v := range res {
		shDetail, err := shipmentRepo.GetShipmentMethodById(ctx, v.ShipmentMethodID)
		if err != nil {
			return []domain.PharmacyShipmentMethods{}, apperror.Wrap(err)
		}

		res[i].Name = &shDetail.Name
	}

	return res, nil
}
