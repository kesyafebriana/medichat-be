package service

import (
	"context"
	"log"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/util"
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

func (s *pharmacyManagerService) CreateClosure(
	ctx context.Context,
	creds domain.AccountRegisterCredentials,
) domain.AtomicFunc[domain.Account] {
	return func(dr domain.DataRepository) (domain.Account, error) {
		accountRepo := dr.AccountRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		log.Print(accountID)
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
