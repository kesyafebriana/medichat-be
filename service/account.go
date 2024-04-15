package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/cryptoutil"
	"medichat-be/domain"
)

type accountService struct {
	dataRepository domain.DataRepository
	passwordHasher cryptoutil.PasswordHasher

	adminAccessProvider  cryptoutil.JWTProvider
	adminRefreshProvider cryptoutil.JWTProvider

	userAccessProvider  cryptoutil.JWTProvider
	userRefreshProvider cryptoutil.JWTProvider

	doctorAccessProvider  cryptoutil.JWTProvider
	doctorRefreshProvider cryptoutil.JWTProvider

	pharmacyManagerAccessProvider  cryptoutil.JWTProvider
	pharmacyManagerRefreshProvider cryptoutil.JWTProvider
}

type AccountServiceOpts struct {
	DataRepository domain.DataRepository
	PasswordHasher cryptoutil.PasswordHasher

	AdminJWTProvider     cryptoutil.JWTProvider
	AdminRefreshProvider cryptoutil.JWTProvider

	UserJWTProvider     cryptoutil.JWTProvider
	UserRefreshProvider cryptoutil.JWTProvider

	DoctorJWTProvider     cryptoutil.JWTProvider
	DoctorRefreshProvider cryptoutil.JWTProvider

	PharmacyManagerJWTProvider     cryptoutil.JWTProvider
	PharmacyManagerRefreshProvider cryptoutil.JWTProvider
}

func NewAccountService(opts AccountServiceOpts) *accountService {
	return &accountService{
		dataRepository: opts.DataRepository,
		passwordHasher: opts.PasswordHasher,

		adminAccessProvider:  opts.AdminJWTProvider,
		adminRefreshProvider: opts.AdminRefreshProvider,

		userAccessProvider:  opts.UserJWTProvider,
		userRefreshProvider: opts.UserRefreshProvider,

		doctorAccessProvider:  opts.DoctorJWTProvider,
		doctorRefreshProvider: opts.DoctorRefreshProvider,

		pharmacyManagerAccessProvider:  opts.PharmacyManagerJWTProvider,
		pharmacyManagerRefreshProvider: opts.PharmacyManagerRefreshProvider,
	}
}

func (s *accountService) RegisterClosure(
	ctx context.Context,
	creds domain.AccountRegisterCredentials,
) domain.AtomicFunc[domain.Account] {
	return func(dr domain.DataRepository) (domain.Account, error) {
		accountRepo := dr.AccountRepository()

		exists, err := accountRepo.IsExistByEmail(ctx, creds.Account.Email)
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}
		if exists {
			return domain.Account{}, apperror.NewAlreadyExists(constants.EntityEmail)
		}

		hashedPassword, err := s.passwordHasher.HashPassword(creds.Password)
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}

		account, err := accountRepo.Add(ctx, domain.AccountWithCredentials{
			Account:        creds.Account,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}

		return account, nil
	}
}

func (s *accountService) Register(
	ctx context.Context,
	creds domain.AccountRegisterCredentials,
) (domain.Account, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.RegisterClosure(ctx, creds),
	)
}

func (s *accountService) Login(
	ctx context.Context,
	creds domain.AccountLoginCredentials,
) (domain.AuthTokens, error) {
	accountRepo := s.dataRepository.AccountRepository()

	ac, err := accountRepo.GetWithCredentialsByEmail(ctx, creds.Email)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	err = s.passwordHasher.CheckPassword(ac.HashedPassword, creds.Password)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	tokens, err := s.createTokensForAccount(ac.Account.ID, ac.Account.Role)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	return tokens, nil
}

func (s *accountService) createTokensForAccount(
	accountID int64,
	role string,
) (domain.AuthTokens, error) {
	var accessProvider cryptoutil.JWTProvider
	var refreshProvider cryptoutil.JWTProvider

	switch role {
	case domain.AccountRoleAdmin:
		accessProvider = s.adminAccessProvider
		refreshProvider = s.adminRefreshProvider
	case domain.AccountRoleUser:
		accessProvider = s.userAccessProvider
		refreshProvider = s.userRefreshProvider
	case domain.AccountRoleDoctor:
		accessProvider = s.doctorAccessProvider
		refreshProvider = s.doctorRefreshProvider
	case domain.AccountRolePharmacyManager:
		accessProvider = s.pharmacyManagerAccessProvider
		refreshProvider = s.pharmacyManagerRefreshProvider
	default:
		return domain.AuthTokens{}, apperror.NewInternalFmt("unknown role")
	}

	accessToken, err := accessProvider.CreateToken(accountID)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	refreshToken, err := refreshProvider.CreateToken(accountID)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	tokens := domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}
