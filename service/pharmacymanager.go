package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/util"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type pharmacyManagerService struct {
	dataRepository domain.DataRepository
	cloudProvider  util.CloudinaryProvider
}

type PharmacyManagerServiceOpts struct {
	DataRepository domain.DataRepository
	CloudProvider  util.CloudinaryProvider
}

func NewPharmacyManagerService(opts PharmacyManagerServiceOpts) *pharmacyManagerService {
	return &pharmacyManagerService{
		dataRepository: opts.DataRepository,
		cloudProvider:  opts.CloudProvider,
	}
}

func (s *pharmacyManagerService) GetAll(ctx context.Context, query domain.PharmacyManagerQuery) ([]domain.Account, error) {
	accountRepo := s.dataRepository.AccountRepository()

	p, err := accountRepo.GetAllPharmacyManager(ctx, query)
	if err != nil {
		return []domain.Account{}, apperror.Wrap(err)
	}

	return p, nil
}

func (s *pharmacyManagerService) CreateClosure(
	ctx context.Context,
	creds domain.AccountRegisterCredentials,
) domain.AtomicFunc[domain.Account] {
	return func(dr domain.DataRepository) (domain.Account, error) {
		accountRepo := dr.AccountRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleAdmin {
			return domain.Account{}, apperror.NewForbidden(nil)
		}

		exists, err := accountRepo.IsExistByEmail(ctx, creds.Account.Email)
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}
		if exists {
			return domain.Account{}, apperror.NewAlreadyExists(constants.EntityEmail)
		}

		account, err = accountRepo.Add(ctx, domain.AccountWithCredentials{
			Account: creds.Account,
		})
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}

		return account, nil
	}
}

func (s *pharmacyManagerService) CreatePharmacyManager(
	ctx context.Context,
	creds domain.AccountRegisterCredentials,
) (domain.Account, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.CreateClosure(ctx, creds),
	)
}

func (s *pharmacyManagerService) CreateProfileClosure(
	ctx context.Context,
	dets domain.PharmacyManagerCreateDetails,
) domain.AtomicFunc[domain.PharmacyManager] {
	return func(dr domain.DataRepository) (domain.PharmacyManager, error) {
		accountRepo := dr.AccountRepository()
		pharmacyManagerRepo := dr.PharmacyManagerRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.PharmacyManager{}, apperror.Wrap(err)
		}

		exists, err := pharmacyManagerRepo.IsExistByAccountID(ctx, accountID)
		if err != nil {
			return domain.PharmacyManager{}, apperror.Wrap(err)
		}
		if exists {
			return domain.PharmacyManager{}, apperror.NewAlreadyExists("pharmacy manager")
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return domain.PharmacyManager{}, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRolePharmacyManager {
			return domain.PharmacyManager{}, apperror.NewForbidden(nil)
		}

		pharmacyManager := domain.PharmacyManager{
			Account: domain.Account{
				ID: accountID,
			},
		}

		account.Name = dets.Name
		if dets.Photo != nil {
			res, err := s.cloudProvider.UploadImage(ctx, dets.Photo, uploader.UploadParams{})
			if err != nil {
				return domain.PharmacyManager{}, apperror.Wrap(err)
			}
			account.PhotoURL = res.SecureURL
		}

		account, err = accountRepo.Update(ctx, account)
		if err != nil {
			return domain.PharmacyManager{}, apperror.Wrap(err)
		}

		err = accountRepo.ProfileSetByID(ctx, accountID)
		if err != nil {
			return domain.PharmacyManager{}, apperror.Wrap(err)
		}

		account.ProfileSet = true
		pharmacyManager.Account = account

		return pharmacyManager, nil
	}
}

func (s *pharmacyManagerService) CreateProfilePharmacyManager(
	ctx context.Context,
	p domain.PharmacyManagerCreateDetails,
) (domain.PharmacyManager, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.CreateProfileClosure(ctx, p),
	)
}

func (s *pharmacyManagerService) DeletePharmacyManager(
	ctx context.Context,
	id int64,
) error {
	accountRepo := s.dataRepository.AccountRepository()
	pharmacyManagerRepo := s.dataRepository.PharmacyManagerRepository()

	exist, err := pharmacyManagerRepo.IsExistByAccountID(ctx, id)
	if err != nil {
		return apperror.Wrap(err)
	}

	if exist {
		err = pharmacyManagerRepo.DeleteByAccountId(ctx, id)
		if err != nil {
			return apperror.Wrap(err)
		}
	}

	err = accountRepo.SoftDeleteById(ctx, id)
	if err != nil {
		return apperror.Wrap(err)
	}

	return nil
}
