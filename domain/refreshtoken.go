package domain

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID        int64
	Account   Account
	Token     string
	ClientIP  string
	ExpiredAt time.Time
}

type RefreshTokenRepository interface {
	Add(ctx context.Context, token RefreshToken) (RefreshToken, error)
	GetByID(ctx context.Context, id int64) (RefreshToken, error)
	GetByTokenStr(ctx context.Context, tokenStr string) (RefreshToken, error)
	GetByTokenStrAndLock(ctx context.Context, tokenStr string) (RefreshToken, error)
	SoftDeleteByID(ctx context.Context, id int64) error
	SoftDeleteByAccountID(ctx context.Context, id int64) error
	SoftDeleteByClientIP(ctx context.Context, ip string) error
}
