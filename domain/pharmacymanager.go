package domain

import "context"

type PharmacyManager struct {
	ID      int64
	Account Account
}

type PharmacyManagerRepository interface {
	GetByID(ctx context.Context, id int64) (PharmacyManager, error)
	GetByIDAndLock(ctx context.Context, id int64) (PharmacyManager, error)

	GetByAccountID(ctx context.Context, id int64) (PharmacyManager, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (PharmacyManager, error)

	Add(ctx context.Context, ph PharmacyManager) (PharmacyManager, error)
	Update(ctx context.Context, ph PharmacyManager) (PharmacyManager, error)
}
