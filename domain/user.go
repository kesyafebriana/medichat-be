package domain

import "context"

type User struct {
	ID int64

	DateOfBirth int64
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByIDAndLock(ctx context.Context, id int64) (User, error)

	GetByAccountID(ctx context.Context, id int64) (User, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (User, error)

	Add(ctx context.Context, u User) (User, error)
	Update(ctx context.Context, u User) (User, error)
}
