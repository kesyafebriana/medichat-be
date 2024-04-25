package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
)

type specializationService struct {
	dataRepository domain.DataRepository
}

type SpecializationServiceOpts struct {
	DataRepository domain.DataRepository
}

func NewSpecializationService(opts SpecializationServiceOpts) *specializationService {
	return &specializationService{
		dataRepository: opts.DataRepository,
	}
}

func (s *specializationService) GetAll(
	ctx context.Context,
) ([]domain.Specialization, error) {
	specializationRepo := s.dataRepository.SpecializationRepository()

	specializations, err := specializationRepo.GetAll(ctx)
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return specializations, nil
}
