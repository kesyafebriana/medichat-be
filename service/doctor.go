package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
)

type doctorService struct {
	dataRepository domain.DataRepository
}

type DoctorServiceOpts struct {
	DataRepository domain.DataRepository
}

func NewDoctorService(opts DoctorServiceOpts) *doctorService {
	return &doctorService{
		dataRepository: opts.DataRepository,
	}
}

func (s *doctorService) List(
	ctx context.Context,
	det domain.DoctorListDetails,
) ([]domain.Doctor, error) {
	doctorRepo := s.dataRepository.DoctorRepository()

	doctors, err := doctorRepo.List(ctx, det)
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return doctors, nil
}

func (s *doctorService) CreateClosure(
	ctx context.Context,
	dets domain.DoctorCreateDetails,
) domain.AtomicFunc[domain.Doctor] {
	return func(dr domain.DataRepository) (domain.Doctor, error) {
		accountRepo := dr.AccountRepository()
		doctorRepo := dr.DoctorRepository()
		specializationRepo := dr.SpecializationRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		exists, err := doctorRepo.IsExistByAccountID(ctx, accountID)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}
		if exists {
			return domain.Doctor{}, apperror.NewAlreadyExists("doctor")
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleDoctor {
			return domain.Doctor{}, apperror.NewForbidden(nil)
		}

		_, err = specializationRepo.GetByID(ctx, dets.SpecializationID)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		doctor := domain.Doctor{
			ID: accountID,
			Account: domain.Account{
				ID: accountID,
			},
			Specialization: domain.Specialization{
				ID: dets.SpecializationID,
			},
			STR:            dets.STR,
			WorkLocation:   dets.WorkLocation,
			Gender:         dets.Gender,
			PhoneNumber:    dets.PhoneNumber,
			IsActive:       dets.IsActive,
			StartWorkDate:  dets.StartWorkDate,
			Price:          dets.Price,
			CertificateURL: dets.CertificateURL,
		}

		doctor, err = doctorRepo.Add(ctx, doctor)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		account.Name = dets.Name

		account, err = accountRepo.Update(ctx, account)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		err = accountRepo.ProfileSetByID(ctx, accountID)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		account.ProfileSet = true
		doctor.Account = account

		return doctor, nil
	}
}

func (s *doctorService) CreateProfile(
	ctx context.Context,
	dets domain.DoctorCreateDetails,
) (domain.Doctor, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.CreateClosure(ctx, dets),
	)
}

func (s *doctorService) UpdateClosure(
	ctx context.Context,
	dets domain.DoctorUpdateDetails,
) domain.AtomicFunc[domain.Doctor] {
	return func(dr domain.DataRepository) (domain.Doctor, error) {
		accountRepo := dr.AccountRepository()
		doctorRepo := dr.DoctorRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleDoctor {
			return domain.Doctor{}, apperror.NewForbidden(nil)
		}

		doctor, err := doctorRepo.GetByAccountIDAndLock(ctx, accountID)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		if dets.Name != nil {
			account.Name = *dets.Name
		}
		doctor.ApplyUpdate(dets)

		account, err = accountRepo.Update(ctx, account)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		doctor, err = doctorRepo.Update(ctx, doctor)
		if err != nil {
			return domain.Doctor{}, apperror.Wrap(err)
		}

		return doctor, nil
	}
}

func (s *doctorService) UpdateProfile(
	ctx context.Context,
	dets domain.DoctorUpdateDetails,
) (domain.Doctor, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.UpdateClosure(ctx, dets),
	)
}

func (s *doctorService) GetProfile(ctx context.Context) (domain.Doctor, error) {
	doctorRepo := s.dataRepository.DoctorRepository()

	accountID, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
		return domain.Doctor{}, apperror.Wrap(err)
	}

	doctor, err := doctorRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return domain.Doctor{}, apperror.Wrap(err)
	}

	return doctor, nil
}

func (s *doctorService) SetActiveStatusClosure(
	ctx context.Context,
	active bool,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		accountRepo := dr.AccountRepository()
		doctorRepo := dr.DoctorRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleDoctor {
			return nil, apperror.NewForbidden(nil)
		}

		doctor, err := doctorRepo.GetByAccountIDAndLock(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		doctor.IsActive = active

		_, err = doctorRepo.Update(ctx, doctor)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		// TODO: notify firebase

		return nil, nil
	}
}

func (s *doctorService) SetActiveStatus(
	ctx context.Context,
	active bool,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.SetActiveStatusClosure(ctx, active),
	)
	return err
}
