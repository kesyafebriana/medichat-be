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
		o, err := pharmacyRepo.AddOperation(ctx, domain.PharmacyOperations{
			PharmacyID: p.ID,
			Day:        v.Day,
			StartTime:  v.StartTime,
			EndTime:    v.EndTime,
		})
		if err != nil {
			return domain.Pharmacy{}, apperror.Wrap(err)
		}

		pharmacyOperations = append(pharmacyOperations, o)
	}

	p.PharmacyOperations = pharmacyOperations

	return p, nil
}
