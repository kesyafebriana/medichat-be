package domain

import "context"

const (
	AccountTypeRegular = "regular"
	AccountTypeGoogle  = "google"

	AccountRoleAdmin           = "admin"
	AccountRoleUser            = "user"
	AccountRoleDoctor          = "doctor"
	AccountRolePharmacyManager = "pharmacy_manager"
)

type Account struct {
	ID            int64
	Email         string
	EmailVerified bool
	Role          string
	AccountType   string
}

type AccountWithCredentials struct {
	Account        Account
	HashedPassword string
}

type AccountLoginCredentials struct {
	Email    string
	Password string
}

type AccountRegisterCredentials struct {
	Account  Account
	Password string
}

type AccountRepository interface {
	GetByEmail(ctx context.Context, email string) (Account, error)
	GetWithCredentialsByEmail(ctx context.Context, email string) (AccountWithCredentials, error)
	IsExistByEmail(ctx context.Context, email string) (bool, error)

	GetByID(ctx context.Context, id int64) (Account, error)
	GetWithCredentialsByID(ctx context.Context, id int64) (AccountWithCredentials, error)
	IsExistByID(ctx context.Context, id int64) (bool, error)

	Add(ctx context.Context, creds AccountWithCredentials) (Account, error)
	UpdateByPassword(ctx context.Context, id int64, newHashedPassword string) error
}

type AccountService interface {
	Register(ctx context.Context, creds AccountRegisterCredentials) (Account, error)
	Login(ctx context.Context, creds AccountLoginCredentials) (AuthTokens, error)
}
