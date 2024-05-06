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
	GetDistance(ctx context.Context, a, b Coordinate) (float64, error)

	AccountRepository() AccountRepository
	ProductRepository() ProductRepository
	ProductDetailsRepository() ProductDetailsRepository
	RefreshTokenRepository() RefreshTokenRepository
	ResetPasswordTokenRepository() ResetPasswordTokenRepository
	VerifyEmailTokenRepository() VerifyEmailTokenRepository
	CategoryRepository() CategoryRepository

	AdminRepository() AdminRepository
	UserRepository() UserRepository
	DoctorRepository() DoctorRepository
	PharmacyManagerRepository() PharmacyManagerRepository

	SpecializationRepository() SpecializationRepository
	PharmacyRepository() PharmacyRepository
	ShipmentMethodRepository() ShipmentMethodRepository

	StockRepository() StockRepository

	PaymentRepository() PaymentRepository
	OrderRepository() OrderRepository
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
