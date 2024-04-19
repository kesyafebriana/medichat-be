package domain

import "context"

type Admin struct {
	ID      int64
	Account Account
}

type AdminRepository interface {
	GetByID(ctx context.Context, id int64) (Admin, error)
	GetByIDAndLock(ctx context.Context, id int64) (Admin, error)

	GetByAccountID(ctx context.Context, id int64) (Admin, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (Admin, error)

	Add(ctx context.Context, a Admin) (Admin, error)
	Update(ctx context.Context, a Admin) (Admin, error)
}
