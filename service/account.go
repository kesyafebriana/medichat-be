package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/cryptoutil"
	"medichat-be/domain"
	"time"
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

	rptProvider cryptoutil.RandomTokenProvider
	rptLifespan time.Duration

	vetProvider cryptoutil.RandomTokenProvider
	vetLifespan time.Duration
}

type AccountServiceOpts struct {
	DataRepository domain.DataRepository
	PasswordHasher cryptoutil.PasswordHasher

	AdminAccessProvider  cryptoutil.JWTProvider
	AdminRefreshProvider cryptoutil.JWTProvider

	UserAccessProvider  cryptoutil.JWTProvider
	UserRefreshProvider cryptoutil.JWTProvider

	DoctorAccessProvider  cryptoutil.JWTProvider
	DoctorRefreshProvider cryptoutil.JWTProvider

	PharmacyManagerAccessProvider  cryptoutil.JWTProvider
	PharmacyManagerRefreshProvider cryptoutil.JWTProvider

	RPTProvider cryptoutil.RandomTokenProvider
	RPTLifespan time.Duration

	VETProvider cryptoutil.RandomTokenProvider
	VETLifespan time.Duration
}

func NewAccountService(opts AccountServiceOpts) *accountService {
	return &accountService{
		dataRepository: opts.DataRepository,
		passwordHasher: opts.PasswordHasher,

		adminAccessProvider:  opts.AdminAccessProvider,
		adminRefreshProvider: opts.AdminRefreshProvider,

		userAccessProvider:  opts.UserAccessProvider,
		userRefreshProvider: opts.UserRefreshProvider,

		doctorAccessProvider:  opts.DoctorAccessProvider,
		doctorRefreshProvider: opts.DoctorRefreshProvider,

		pharmacyManagerAccessProvider:  opts.PharmacyManagerAccessProvider,
		pharmacyManagerRefreshProvider: opts.PharmacyManagerRefreshProvider,

		rptProvider: opts.RPTProvider,
		rptLifespan: opts.RPTLifespan,

		vetProvider: opts.VETProvider,
		vetLifespan: opts.VETLifespan,
	}
}

func (s *accountService) RegisterClosure(
	ctx context.Context,
	creds domain.AccountRegisterCredentials,
) domain.AtomicFunc[domain.Account] {
	return func(dr domain.DataRepository) (domain.Account, error) {
		accountRepo := dr.AccountRepository()

		if creds.Account.Role == domain.AccountRoleAdmin ||
			creds.Account.Role == domain.AccountRolePharmacyManager {
			return domain.Account{}, apperror.NewAppError(
				apperror.CodeBadRequest,
				"cannot register privileged account",
				nil,
			)
		}

		exists, err := accountRepo.IsExistByEmail(ctx, creds.Account.Email)
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
		}
		if exists {
			return domain.Account{}, apperror.NewAlreadyExists(constants.EntityEmail)
		}

		account, err := accountRepo.Add(ctx, domain.AccountWithCredentials{
			Account: creds.Account,
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

func (s *accountService) GetResetPasswordTokenClosure(
	ctx context.Context,
	email string,
) domain.AtomicFunc[string] {
	return func(dr domain.DataRepository) (string, error) {
		accountRepo := s.dataRepository.AccountRepository()
		rptRepo := s.dataRepository.ResetPasswordTokenRepository()

		account, err := accountRepo.GetByEmail(ctx, email)
		if err != nil {
			return "", apperror.Wrap(err)
		}

		err = rptRepo.SoftDeleteByAccountID(ctx, account.ID)
		if err != nil {
			return "", apperror.Wrap(err)
		}

		tokenStr, err := s.rptProvider.GenerateToken()
		if err != nil {
			return "", apperror.Wrap(err)
		}

		token := domain.ResetPasswordToken{
			Account:   account,
			Token:     tokenStr,
			ExpiredAt: time.Now().Add(s.rptLifespan),
		}

		_, err = rptRepo.Add(ctx, token)
		if err != nil {
			return "", apperror.Wrap(err)
		}

		return tokenStr, nil
	}
}

func (s *accountService) GetResetPasswordToken(
	ctx context.Context,
	email string,
) (string, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.GetResetPasswordTokenClosure(ctx, email),
	)
}

func (s *accountService) ResetPasswordClosure(
	ctx context.Context,
	creds domain.AccountResetPasswordCredentials,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		accountRepo := dr.AccountRepository()
		rptRepo := dr.ResetPasswordTokenRepository()

		token, err := rptRepo.GetByTokenStrAndLock(ctx, creds.ResetPasswordToken)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByIDAndLock(ctx, token.Account.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		if account.Email != creds.Email {
			return nil, apperror.NewForbidden(nil)
		}

		hashedPassword, err := s.passwordHasher.HashPassword(creds.NewPassword)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = accountRepo.UpdatePasswordByID(ctx, account.ID, hashedPassword)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = rptRepo.SoftDeleteByID(ctx, token.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
	}
}

func (s *accountService) ResetPassword(
	ctx context.Context,
	creds domain.AccountResetPasswordCredentials,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.ResetPasswordClosure(ctx, creds),
	)

	return err
}

func (s *accountService) GetVerifyEmailTokenClosure(
	ctx context.Context,
	email string,
) domain.AtomicFunc[string] {
	return func(dr domain.DataRepository) (string, error) {
		accountRepo := s.dataRepository.AccountRepository()
		vetRepo := s.dataRepository.VerifyEmailTokenRepository()

		account, err := accountRepo.GetByEmail(ctx, email)
		if err != nil {
			return "", apperror.Wrap(err)
		}

		if account.EmailVerified {
			return "", apperror.NewEmailAlreadyVerified(nil)
		}

		err = vetRepo.SoftDeleteByAccountID(ctx, account.ID)
		if err != nil {
			return "", apperror.Wrap(err)
		}

		tokenStr, err := s.vetProvider.GenerateToken()
		if err != nil {
			return "", apperror.Wrap(err)
		}

		token := domain.VerifyEmailToken{
			Account:   account,
			Token:     tokenStr,
			ExpiredAt: time.Now().Add(s.rptLifespan),
		}

		_, err = vetRepo.Add(ctx, token)
		if err != nil {
			return "", apperror.Wrap(err)
		}

		return tokenStr, nil
	}
}

func (s *accountService) GetVerifyEmailToken(
	ctx context.Context,
	email string,
) (string, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.GetVerifyEmailTokenClosure(ctx, email),
	)
}

func (s *accountService) VerifyEmailClosure(
	ctx context.Context,
	creds domain.AccountVerifyEmailCredentials,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		accountRepo := dr.AccountRepository()
		rptRepo := dr.VerifyEmailTokenRepository()

		token, err := rptRepo.GetByTokenStrAndLock(ctx, creds.VerifyEmailToken)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByIDAndLock(ctx, token.Account.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		if account.Email != creds.Email {
			return nil, apperror.NewForbidden(nil)
		}

		if account.EmailVerified {
			return "", apperror.NewEmailAlreadyVerified(nil)
		}

		hashedPassword, err := s.passwordHasher.HashPassword(creds.Password)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = accountRepo.VerifyEmailByID(ctx, account.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = accountRepo.UpdatePasswordByID(ctx, account.ID, hashedPassword)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = rptRepo.SoftDeleteByID(ctx, token.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
	}
}

func (s *accountService) VerifyEmail(
	ctx context.Context,
	creds domain.AccountVerifyEmailCredentials,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.VerifyEmailClosure(ctx, creds),
	)

	return err
}
