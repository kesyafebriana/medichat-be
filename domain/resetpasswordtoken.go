package domain

import (
	"context"
	"time"
)

type ResetPasswordToken struct {
	ID        int64
	Account   Account
	Token     string
	ExpiredAt time.Time
}

type ResetPasswordTokenRepository interface {
	Add(ctx context.Context, token ResetPasswordToken) (ResetPasswordToken, error)
	GetByID(ctx context.Context, id int64) (ResetPasswordToken, error)
	GetByTokenStr(ctx context.Context, tokenStr string) (ResetPasswordToken, error)
	GetByTokenStrAndLock(ctx context.Context, tokenStr string) (ResetPasswordToken, error)
	SoftDeleteByID(ctx context.Context, id int64) error
}
