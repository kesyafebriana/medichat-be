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

	adminAccessProvider           cryptoutil.JWTProvider
	userAccessProvider            cryptoutil.JWTProvider
	doctorAccessProvider          cryptoutil.JWTProvider
	pharmacyManagerAccessProvider cryptoutil.JWTProvider
	refreshProvider               cryptoutil.JWTProvider

	rptProvider cryptoutil.RandomTokenProvider
	rptLifespan time.Duration

	vetProvider cryptoutil.RandomTokenProvider
	vetLifespan time.Duration
}

type AccountServiceOpts struct {
	DataRepository domain.DataRepository
	PasswordHasher cryptoutil.PasswordHasher

	AdminAccessProvider           cryptoutil.JWTProvider
	UserAccessProvider            cryptoutil.JWTProvider
	DoctorAccessProvider          cryptoutil.JWTProvider
	PharmacyManagerAccessProvider cryptoutil.JWTProvider
	RefreshProvider               cryptoutil.JWTProvider

	RPTProvider cryptoutil.RandomTokenProvider
	RPTLifespan time.Duration

	VETProvider cryptoutil.RandomTokenProvider
	VETLifespan time.Duration
}

func NewAccountService(opts AccountServiceOpts) *accountService {
	return &accountService{
		dataRepository: opts.DataRepository,
		passwordHasher: opts.PasswordHasher,

		adminAccessProvider:           opts.AdminAccessProvider,
		userAccessProvider:            opts.UserAccessProvider,
		doctorAccessProvider:          opts.DoctorAccessProvider,
		pharmacyManagerAccessProvider: opts.PharmacyManagerAccessProvider,
		refreshProvider:               opts.RefreshProvider,

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

func (s *accountService) LoginClosure(
	ctx context.Context,
	creds domain.AccountLoginCredentials,
) domain.AtomicFunc[domain.AuthTokens] {
	return func(dr domain.DataRepository) (domain.AuthTokens, error) {
		accountRepo := s.dataRepository.AccountRepository()
		rtRepo := s.dataRepository.RefreshTokenRepository()

		ac, err := accountRepo.GetWithCredentialsByEmail(ctx, creds.Email)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		err = s.passwordHasher.CheckPassword(ac.HashedPassword, creds.Password)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		tokens, err := s.CreateTokensForAccount(ac.Account.ID, ac.Account.Role)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		rToken := domain.RefreshToken{
			Account:   ac.Account,
			Token:     tokens.RefreshToken,
			ClientIP:  creds.ClientIP,
			ExpiredAt: tokens.RefreshExpireAt,
		}

		_, err = rtRepo.Add(ctx, rToken)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		return tokens, nil
	}
}

func (s *accountService) Login(
	ctx context.Context,
	creds domain.AccountLoginCredentials,
) (domain.AuthTokens, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.LoginClosure(ctx, creds),
	)
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

		if !account.EmailVerified {
			return "", apperror.NewEmailNotVerified(nil)
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

		if !account.EmailVerified {
			return "", apperror.NewEmailNotVerified(nil)
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

func (s *accountService) RefreshTokensClosure(
	ctx context.Context,
	creds domain.AccountRefreshTokensCredentials,
) domain.AtomicFunc[domain.AuthTokens] {
	return func(dr domain.DataRepository) (domain.AuthTokens, error) {
		accountRepo := s.dataRepository.AccountRepository()
		rtRepo := s.dataRepository.RefreshTokenRepository()

		claims, err := s.refreshProvider.VerifyToken(creds.RefreshToken)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		rToken, err := rtRepo.GetByTokenStrAndLock(ctx, creds.RefreshToken)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		tokens, err := s.CreateTokensForAccount(account.ID, account.Role)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		err = rtRepo.SoftDeleteByID(ctx, rToken.ID)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		rToken = domain.RefreshToken{
			Account:   account,
			Token:     tokens.RefreshToken,
			ClientIP:  creds.ClientIP,
			ExpiredAt: tokens.RefreshExpireAt,
		}

		_, err = rtRepo.Add(ctx, rToken)
		if err != nil {
			return domain.AuthTokens{}, apperror.Wrap(err)
		}

		return tokens, nil
	}
}

func (s *accountService) RefreshTokens(
	ctx context.Context,
	creds domain.AccountRefreshTokensCredentials,
) (domain.AuthTokens, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.RefreshTokensClosure(ctx, creds),
	)
}

func (s *accountService) CreateTokensForAccount(
	accountID int64,
	role string,
) (domain.AuthTokens, error) {
	var accessProvider cryptoutil.JWTProvider
	var refreshProvider cryptoutil.JWTProvider

	switch role {
	case domain.AccountRoleAdmin:
		accessProvider = s.adminAccessProvider
	case domain.AccountRoleUser:
		accessProvider = s.userAccessProvider
	case domain.AccountRoleDoctor:
		accessProvider = s.doctorAccessProvider
	case domain.AccountRolePharmacyManager:
		accessProvider = s.pharmacyManagerAccessProvider
	default:
		return domain.AuthTokens{}, apperror.NewInternalFmt("unknown role")
	}
	refreshProvider = s.refreshProvider

	accessToken, err := accessProvider.CreateToken(accountID)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}
	accessClaims, err := accessProvider.VerifyToken(accessToken)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	refreshToken, err := refreshProvider.CreateToken(accountID)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}
	refreshClaims, err := refreshProvider.VerifyToken(refreshToken)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	tokens := domain.AuthTokens{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessClaims.ExpiresAt.Time,
		RefreshExpireAt: refreshClaims.ExpiresAt.Time,
	}

	return tokens, nil
}
