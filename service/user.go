package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
)

type userService struct {
	dataRepository domain.DataRepository
}

type UserServiceOpts struct {
	DataRepository domain.DataRepository
}

func NewUserService(opts UserServiceOpts) *userService {
	return &userService{
		dataRepository: opts.DataRepository,
	}
}

func (s *userService) CreateClosure(
	ctx context.Context,
	dets domain.UserCreateDetails,
) domain.AtomicFunc[domain.User] {
	return func(dr domain.DataRepository) (domain.User, error) {
		accountRepo := dr.AccountRepository()
		userRepo := dr.UserRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		if accountID != dets.AccountID {
			return domain.User{}, apperror.NewForbidden(nil)
		}

		exists, err := userRepo.IsExistByAccountID(ctx, accountID)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}
		if exists {
			return domain.User{}, apperror.NewAlreadyExists("user")
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleUser {
			return domain.User{}, apperror.NewForbidden(nil)
		}

		user := domain.User{
			Account: domain.Account{
				ID: dets.AccountID,
			},
			DateOfBirth: dets.DateOfBirth,
		}

		user, err = userRepo.Add(ctx, user)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		locations := dets.Locations
		for i := 0; i < len(locations); i++ {
			locations[i].UserID = user.ID
		}

		err = userRepo.AddLocations(ctx, locations)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		account.Name = dets.Name

		account, err = accountRepo.Update(ctx, account)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		err = accountRepo.ProfileSetByID(ctx, accountID)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		account.ProfileSet = true
		user.Account = account

		return user, nil
	}
}

func (s *userService) CreateProfile(
	ctx context.Context,
	u domain.UserCreateDetails,
) (domain.User, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.CreateClosure(ctx, u),
	)
}

func (s *userService) UpdateClosure(
	ctx context.Context,
	u domain.UserUpdateDetails,
) domain.AtomicFunc[domain.User] {
	return func(dr domain.DataRepository) (domain.User, error) {
		accountRepo := dr.AccountRepository()
		userRepo := dr.UserRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		account, err := accountRepo.GetByIDAndLock(ctx, accountID)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleUser {
			return domain.User{}, apperror.NewForbidden(nil)
		}

		user, err := userRepo.GetByAccountIDAndLock(ctx, accountID)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		if u.ID != user.ID {
			return domain.User{}, apperror.NewForbidden(nil)
		}

		if u.Name != nil {
			account.Name = *u.Name
		}
		if u.DateOfBirth != nil {
			user.DateOfBirth = *u.DateOfBirth
		}

		account, err = accountRepo.Update(ctx, account)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		user, err = userRepo.Update(ctx, user)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		return user, nil
	}
}

func (s *userService) UpdateProfile(
	ctx context.Context,
	u domain.UserUpdateDetails,
) (domain.User, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.UpdateClosure(ctx, u),
	)
}

func (s *userService) GetProfile(
	ctx context.Context,
) (domain.User, error) {
	userRepo := s.dataRepository.UserRepository()

	accountID, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
		return domain.User{}, apperror.Wrap(err)
	}

	user, err := userRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return domain.User{}, apperror.Wrap(err)
	}

	locations, err := userRepo.GetLocationsByUserID(ctx, user.ID)
	if err != nil {
		return domain.User{}, apperror.Wrap(err)
	}

	user.Locations = locations

	return user, nil
}
