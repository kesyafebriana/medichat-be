package domain

import (
	"context"
	"mime/multipart"
)

type PharmacyManager struct {
	ID      int64
	Account Account
}

type PharmacyManagerCreateDetails struct {
	Name  string
	Photo multipart.File
}

type PharmacyManagerRepository interface {
	GetByID(ctx context.Context, id int64) (PharmacyManager, error)
	GetByIDAndLock(ctx context.Context, id int64) (PharmacyManager, error)
	IsExistByID(ctx context.Context, id int64) (bool, error)

	GetByAccountID(ctx context.Context, id int64) (PharmacyManager, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (PharmacyManager, error)
	IsExistByAccountID(ctx context.Context, id int64) (bool, error)

	Add(ctx context.Context, ph PharmacyManager) (PharmacyManager, error)
	//soft delete
}

type PharmacyManagerService interface {
	CreatePharmacyManager(ctx context.Context, creds AccountRegisterCredentials) (Account, error)
	CreateProfilePharmacyManager(ctx context.Context, p PharmacyManagerCreateDetails) (PharmacyManager, error)
	// GetPharmacyManagers(ctx context.Context) (PharmacyManager, error)
}
