package domain

import (
	"context"
	"medichat-be/apperror"
	"time"
)

type AtomicFunc[T any] func(DataRepository) (T, error)
type AtomicFuncAny AtomicFunc[any]

type DataRepository interface {
	Atomic(ctx context.Context, fn AtomicFuncAny) (any, error)
	Sleep(ctx context.Context, duration time.Duration) error

	AccountRepository() AccountRepository
	ProductRepository() ProductRepository
	RefreshTokenRepository() RefreshTokenRepository
	ResetPasswordTokenRepository() ResetPasswordTokenRepository
	VerifyEmailTokenRepository() VerifyEmailTokenRepository
}

func RunAtomic[T any](
	dataRepo DataRepository,
	ctx context.Context,
	fn AtomicFunc[T],
) (T, error) {
	temp, err := dataRepo.Atomic(ctx, func(dr DataRepository) (any, error) {
		return fn(dr)
	})
	if err != nil {
		var t T
		return t, apperror.Wrap(err)
	}

	if temp == nil {
		var t T
		return t, nil
	}

	ret, ok := temp.(T)
	if !ok {
		var t T
		return t, apperror.NewTypeAssertionFailed(temp, ret)
	}

	return ret, nil
}
