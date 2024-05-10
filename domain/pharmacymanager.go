package domain

import (
	"context"
	"mime/multipart"
)

const (
	PharmacyManagerSortByCreatedAt     = "created_at"
)

type PharmacyManager struct {
	ID      int64
	Account Account
}

type PharmacyManagerCreateDetails struct {
	Name  string
	Photo multipart.File
}

type PharmacyManagerQuery struct {
	Page       int64
	Limit      int64
	Level      int64
	Term       string
	SortBy     string
	SortType   string
	ProfileSet *string
}

type PharmacyManagerRepository interface {
	GetByID(ctx context.Context, id int64) (PharmacyManager, error)
	GetByIDAndLock(ctx context.Context, id int64) (PharmacyManager, error)
	IsExistByID(ctx context.Context, id int64) (bool, error)

	GetByAccountID(ctx context.Context, id int64) (PharmacyManager, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (PharmacyManager, error)
	IsExistByAccountID(ctx context.Context, id int64) (bool, error)

	Add(ctx context.Context, ph PharmacyManager) (PharmacyManager, error)
	DeleteByAccountId(ctx context.Context, id int64) error
}

type PharmacyManagerService interface {
	CreatePharmacyManager(ctx context.Context, creds AccountRegisterCredentials) (Account, error)
	CreateProfilePharmacyManager(ctx context.Context, p PharmacyManagerCreateDetails) (PharmacyManager, error)
	GetAll(ctx context.Context, query PharmacyManagerQuery) ([]Account, PageInfo, error)
	DeletePharmacyManager(ctx context.Context, id int64) error
}
