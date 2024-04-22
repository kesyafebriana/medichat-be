package util

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/constants"
)

func GetAccountIDFromContext(ctx context.Context) (int64, error) {
	val := ctx.Value(constants.ContextAccountID)
	id, ok := val.(int64)
	if !ok {
		return 0, apperror.NewTypeAssertionFailed(id, val)
	}

	return id, nil
}
