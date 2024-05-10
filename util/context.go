package util

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
)

func GetAccountIDFromContext(ctx context.Context) (int64, error) {
	val := ctx.Value(constants.ContextAccountID)
	id, ok := val.(int64)
	if !ok {
		return 0, apperror.NewTypeAssertionFailed(id, val)
	}

	return id, nil
}

func GetProfileFromContext(ctx context.Context, dr domain.DataRepository) (string, any, error) {
	accountRepo := dr.AccountRepository()
	userRepo := dr.UserRepository()
	doctorRepo := dr.DoctorRepository()
	managerRepo := dr.PharmacyManagerRepository()

	accountID, err := GetAccountIDFromContext(ctx)
	if err != nil {
		return "", nil, err
	}

	account, err := accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return "", nil, err
	}

	role := account.Role

	switch role {
	case domain.AccountRoleUser:
		user, err := userRepo.GetByAccountID(ctx, accountID)
		return role, user, err
	case domain.AccountRoleDoctor:
		doctor, err := doctorRepo.GetByAccountID(ctx, accountID)
		return role, doctor, err
	case domain.AccountRolePharmacyManager:
		manager, err := managerRepo.GetByAccountID(ctx, accountID)
		return role, manager, err
	default:
		return role, account, nil
	}
}
