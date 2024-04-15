package domain

import (
	"context"
	"medichat-be/apperror"
	"time"
)

type AtomicFunc[T any] func(DataRepository) (T, error)

type DataRepository interface {
	Atomic(ctx context.Context, fn AtomicFunc[any]) (any, error)
	Sleep(ctx context.Context, duration time.Duration) error

	AccountRepository() AccountRepository
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

	ret, ok := temp.(T)
	if !ok {
		var t T
		return t, apperror.NewTypeAssertionFailed(temp, ret)
	}

	return ret, nil
}
