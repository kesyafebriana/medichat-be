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

func (s *pharmacyService) CreatePharmacy(ctx context.Context, pharmacy domain.PharmacyCreateDetails) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	var pharmacyOperations []domain.PharmacyOperations

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

	p.PharmacyOperations = pharmacyOperations

	return p, nil
}

func (s *pharmacyService) AddOperation(ctx context.Context, pharmacyOperation domain.PharmacyOperationCreateDetails) (domain.PharmacyOperations, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	o, err := pharmacyRepo.AddOperation(ctx, pharmacyOperation)
	if err != nil {
		return domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	return o, nil
}

func (s *pharmacyService) GetOperationsById(ctx context.Context, id int64) ([]domain.PharmacyOperations, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, id)
	if err != nil {
		return []domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	return o, nil
}

func (s *pharmacyService) UpdateOperations(ctx context.Context, pharmacyOperations []domain.PharmacyOperationsUpdateDetails) ([]domain.PharmacyOperations, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	var res []domain.PharmacyOperations

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyIdAndLock(ctx, pharmacyOperations[0].PharmacyID)
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
			PharmacyID: v.PharmacyID,
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
