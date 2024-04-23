package domain

import (
	"context"
	"time"
)

type User struct {
	ID      int64
	Account Account

	DateOfBirth time.Time
}

type UserCreateDetails struct {
	AccountID   int64
	Name        string
	PhotoURL    string
	DateOfBirth time.Time
}

type UserUpdateDetails struct {
	ID          int64
	Name        *string
	PhotoURL    *string
	DateOfBirth *time.Time
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByIDAndLock(ctx context.Context, id int64) (User, error)
	IsExistByID(ctx context.Context, id int64) (bool, error)

	GetByAccountID(ctx context.Context, id int64) (User, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (User, error)
	IsExistByAccountID(ctx context.Context, id int64) (bool, error)

	Add(ctx context.Context, u User) (User, error)
	Update(ctx context.Context, u User) (User, error)
}

type UserService interface {
	CreateProfile(ctx context.Context, u UserCreateDetails) (User, error)
	UpdateProfile(ctx context.Context, u UserUpdateDetails) (User, error)
}
