package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
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

func (s *pharmacyService) GetPharmacies(ctx context.Context, query domain.PharmaciesQuery) ([]domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	p, err := pharmacyRepo.GetPharmacies(ctx, query)
	if err != nil {
		return []domain.Pharmacy{}, apperror.Wrap(err)
	}

	return p, nil
}

func (s *pharmacyService) GetPharmacyBySlug(ctx context.Context, slug string) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	p, err := pharmacyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, p.ID)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	p.PharmacyOperations = o

	return p, nil
}

func (s *pharmacyService) CreatePharmacy(ctx context.Context, pharmacy domain.PharmacyCreateDetails) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()
	var pharmacyOperations []domain.PharmacyOperations
	pharmacy.Slug = util.GenerateSlug(pharmacy.Name)

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

func (s *pharmacyService) UpdatePharmacy(ctx context.Context, pharmacy domain.PharmacyUpdateDetails) (domain.Pharmacy, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	p, err := pharmacyRepo.Update(ctx, pharmacy)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	o, err := pharmacyRepo.GetPharmacyOperationsByPharmacyId(ctx, p.ID)
	if err != nil {
		return domain.Pharmacy{}, apperror.Wrap(err)
	}

	p.PharmacyOperations = o

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

func (s *pharmacyService) AddOperation(ctx context.Context, pharmacyOperation domain.PharmacyOperationCreateDetails) (domain.PharmacyOperations, error) {
	pharmacyRepo := s.dataRepository.PharmacyRepository()

	p, err := pharmacyRepo.GetBySlug(ctx, pharmacyOperation.Slug)
	if err != nil {
		return domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	pharmacyOperation.PharmacyID = p.ID

	o, err := pharmacyRepo.AddOperation(ctx, pharmacyOperation)
	if err != nil {
		return domain.PharmacyOperations{}, apperror.Wrap(err)
	}

	return o, nil
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
