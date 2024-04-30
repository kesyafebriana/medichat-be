package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type userService struct {
	dataRepository domain.DataRepository
	cloudProvider  util.CloudinaryProvider
}

type UserServiceOpts struct {
	DataRepository domain.DataRepository
	CloudProvider  util.CloudinaryProvider
}

func NewUserService(opts UserServiceOpts) *userService {
	return &userService{
		dataRepository: opts.DataRepository,
		cloudProvider:  opts.CloudProvider,
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
				ID: accountID,
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

		locations, err = userRepo.AddLocations(ctx, locations)
		if err != nil {
			return domain.User{}, apperror.Wrap(err)
		}

		if len(locations) > 0 {
			user.MainLocationID = locations[0].ID
			user, err = userRepo.Update(ctx, user)
			if err != nil {
				return domain.User{}, apperror.Wrap(err)
			}
		}

		account.Name = dets.Name
		if dets.Photo != nil {
			res, err := s.cloudProvider.UploadImage(ctx, dets.Photo, uploader.UploadParams{})
			if err != nil {
				return domain.User{}, apperror.Wrap(err)
			}
			account.PhotoURL = res.SecureURL
		}

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

		if u.Name != nil {
			account.Name = *u.Name
		}
		if u.Photo != nil {
			res, err := s.cloudProvider.UploadImage(ctx, u.Photo, uploader.UploadParams{})
			if err != nil {
				return domain.User{}, apperror.Wrap(err)
			}
			account.PhotoURL = res.SecureURL
		}

		if u.DateOfBirth != nil {
			user.DateOfBirth = *u.DateOfBirth
		}
		if u.MainLocationID != nil {
			user.MainLocationID = *u.MainLocationID
			location, err := userRepo.GetLocationByID(ctx, *u.MainLocationID)
			if err != nil {
				return domain.User{}, apperror.Wrap(err)
			}
			if location.UserID != user.ID {
				return domain.User{}, apperror.NewForbidden(nil)
			}
			if !location.IsActive {
				return domain.User{}, apperror.NewUserLocationIsNotActive(nil)
			}
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

func (s *userService) AddLocationClosure(
	ctx context.Context,
	ul domain.UserLocation,
) domain.AtomicFunc[domain.UserLocation] {
	return func(dr domain.DataRepository) (domain.UserLocation, error) {
		userRepo := dr.UserRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		ul.UserID = user.ID

		ul, err = userRepo.AddLocation(ctx, ul)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		return ul, nil
	}
}

func (s *userService) AddLocation(
	ctx context.Context,
	ul domain.UserLocation,
) (domain.UserLocation, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.AddLocationClosure(ctx, ul),
	)
}

func (s *userService) UpdateLocationClosure(
	ctx context.Context,
	det domain.UserLocationUpdateDetails,
) domain.AtomicFunc[domain.UserLocation] {
	return func(dr domain.DataRepository) (domain.UserLocation, error) {
		userRepo := dr.UserRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		ul, err := userRepo.GetLocationByIDAndLock(ctx, det.ID)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		if ul.UserID != user.ID {
			return domain.UserLocation{}, apperror.NewForbidden(nil)
		}

		if det.Alias != nil {
			ul.Alias = *det.Alias
		}
		if det.Address != nil {
			ul.Address = *det.Address
		}
		if det.Coordinate != nil {
			ul.Coordinate = *det.Coordinate
		}
		if det.IsActive != nil {
			ul.IsActive = *det.IsActive
			if !*det.IsActive {
				if user.MainLocationID == det.ID {
					return domain.UserLocation{}, apperror.NewUserLocationCannotDeleteMain(nil)
				}
			}
		}

		ul, err = userRepo.UpdateLocation(ctx, ul)
		if err != nil {
			return domain.UserLocation{}, apperror.Wrap(err)
		}

		if det.IsActive != nil && !*det.IsActive {
			exists, err := userRepo.IsAnyLocationActiveByUserID(ctx, user.ID)
			if err != nil {
				return domain.UserLocation{}, apperror.Wrap(err)
			}
			if !exists {
				return domain.UserLocation{}, apperror.NewUserLocationShouldHaveActive(nil)
			}
		}

		return ul, nil
	}
}

func (s *userService) UpdateLocation(
	ctx context.Context,
	det domain.UserLocationUpdateDetails,
) (domain.UserLocation, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.UpdateLocationClosure(ctx, det),
	)
}

func (s *userService) DeleteLocationByIDClosure(
	ctx context.Context,
	id int64,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		userRepo := dr.UserRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if user.MainLocationID == id {
			return nil, apperror.NewUserLocationCannotDeleteMain(nil)
		}

		ul, err := userRepo.GetLocationByIDAndLock(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		if ul.UserID != user.ID {
			return nil, apperror.NewForbidden(nil)
		}

		err = userRepo.SoftDeleteLocationByID(ctx, id)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		exists, err := userRepo.IsAnyLocationActiveByUserID(ctx, user.ID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if !exists {
			return nil, apperror.NewUserLocationShouldHaveActive(nil)
		}

		return nil, nil
	}
}

func (s *userService) DeleteLocationByID(
	ctx context.Context,
	id int64,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.DeleteLocationByIDClosure(ctx, id),
	)
	return err
}
