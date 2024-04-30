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

func (s *pharmacyService) CreatePharmacy(ctx context.Context, pharmacy domain.Pharmacy) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	p, err := pharmacyRepo.Add(ctx, pharmacy)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	return p, nil
}
