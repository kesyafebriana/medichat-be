package testdata

import (
	"context"
	"medichat-be/domain"
	"medichat-be/mocks/domainmocks"

	"github.com/stretchr/testify/mock"
)

type DataRepositoryMockOpts struct {
	AccountRepository            domain.AccountRepository
	ResetPasswordTokenRepository domain.ResetPasswordTokenRepository
	VerifyEmailTokenRepository   domain.VerifyEmailTokenRepository
}

func NewDataRepositoryMock(opts DataRepositoryMockOpts) *domainmocks.DataRepository {
	dataRepo := new(domainmocks.DataRepository)

	dataRepo.On("AccountRepository").
		Return(opts.AccountRepository)
	dataRepo.On("ResetPasswordTokenRepository").
		Return(opts.ResetPasswordTokenRepository)
	dataRepo.On("VerifyEmailTokenRepository").
		Return(opts.VerifyEmailTokenRepository)

	return dataRepo
}

type Result[T any] struct {
	Val T
	Err error
}

type WantValue[T any] struct {
	Val T
	Err int
}

func OnDataRepositoryAtomic[T any](
	dr *domainmocks.DataRepository,
	ctx context.Context,
	fn domain.AtomicFunc[T],
) {
	ret, err := fn(dr)
	dr.On("Atomic", ctx, mock.Anything).Return(ret, err)
}
