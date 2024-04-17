package domain

import (
	"context"
	"time"
)

type VerifyEmailToken struct {
	ID        int64
	Account   Account
	Token     string
	ExpiredAt time.Time
}

type VerifyEmailTokenRepository interface {
	Add(ctx context.Context, token VerifyEmailToken) (VerifyEmailToken, error)
	GetByID(ctx context.Context, id int64) (VerifyEmailToken, error)
	GetByTokenStr(ctx context.Context, tokenStr string) (VerifyEmailToken, error)
	GetByTokenStrAndLock(ctx context.Context, tokenStr string) (VerifyEmailToken, error)
	SoftDeleteByID(ctx context.Context, id int64) error
	SoftDeleteByAccountID(ctx context.Context, id int64) error
}
