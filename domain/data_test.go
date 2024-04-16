package domain_test

import (
	"context"
	"errors"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/mocks/domainmocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrSomething = errors.New("something")
)

func TestRunAtomic(t *testing.T) {
	t.Run("should return int when T is int", func(t *testing.T) {
		// given
		dataRepo := new(domainmocks.DataRepository)
		ctx := context.Background()
		want := 10
		fn := func(dr domain.DataRepository) (int, error) {
			return want, nil
		}

		dataRepo.On("Atomic", ctx, mock.Anything).Return(
			want, nil,
		)

		// when
		got, _ := domain.RunAtomic(dataRepo, ctx, fn)

		// then
		assert.Equal(t, want, got)
	})

	t.Run("should return int when T is string", func(t *testing.T) {
		// given
		dataRepo := new(domainmocks.DataRepository)
		ctx := context.Background()
		want := "hello"
		fn := func(dr domain.DataRepository) (string, error) {
			return want, nil
		}

		dataRepo.On("Atomic", ctx, mock.Anything).Return(
			want, nil,
		)

		// when
		got, _ := domain.RunAtomic(dataRepo, ctx, fn)

		// then
		assert.Equal(t, want, got)
	})

	t.Run("should return int when atomic func returns nil", func(t *testing.T) {
		// given
		dataRepo := new(domainmocks.DataRepository)
		ctx := context.Background()
		var want any = nil
		fn := func(dr domain.DataRepository) (any, error) {
			return want, nil
		}

		dataRepo.On("Atomic", ctx, mock.Anything).Return(
			want, nil,
		)

		// when
		got, _ := domain.RunAtomic(dataRepo, ctx, fn)

		// then
		assert.Equal(t, want, got)
	})

	t.Run("should return error when T is int and atomic func returns error", func(t *testing.T) {
		// given
		dataRepo := new(domainmocks.DataRepository)
		ctx := context.Background()
		fn := func(dr domain.DataRepository) (int, error) {
			return 0, ErrSomething
		}

		dataRepo.On("Atomic", ctx, mock.Anything).Return(
			0, ErrSomething,
		)

		// when
		_, err := domain.RunAtomic(dataRepo, ctx, fn)

		// then
		assert.EqualError(t, err, apperror.Wrap(ErrSomething).Error())
	})
}
