package domain

import (
	"context"
	"mime/multipart"
	"time"
)

type User struct {
	ID      int64
	Account Account

	DateOfBirth    time.Time
	MainLocationID int64
	Locations      []UserLocation
}

type UserLocation struct {
	ID     int64
	UserID int64

	Alias      string
	Address    string
	Coordinate Coordinate
	IsActive   bool
}

type UserCreateDetails struct {
	Name        string
	Photo       multipart.File
	DateOfBirth time.Time
	Locations   []UserLocation
}

type UserUpdateDetails struct {
	Name           *string
	Photo          multipart.File
	DateOfBirth    *time.Time
	MainLocationID *int64
}

type UserLocationUpdateDetails struct {
	ID int64

	Alias      *string
	Address    *string
	Coordinate *Coordinate
	IsActive   *bool
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByIDAndLock(ctx context.Context, id int64) (User, error)
	IsExistByID(ctx context.Context, id int64) (bool, error)

	GetByAccountID(ctx context.Context, id int64) (User, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (User, error)
	IsExistByAccountID(ctx context.Context, id int64) (bool, error)

	Add(ctx context.Context, u User) (User, error)
	Update(ctx context.Context, u User) (User, error)

	GetLocationsByUserID(ctx context.Context, id int64) ([]UserLocation, error)
	GetLocationByID(ctx context.Context, id int64) (UserLocation, error)
	GetLocationByIDAndLock(ctx context.Context, id int64) (UserLocation, error)
	IsAnyLocationActiveByUserID(Ctx context.Context, id int64) (bool, error)

	AddLocation(ctx context.Context, ul UserLocation) (UserLocation, error)
	AddLocations(ctx context.Context, uls []UserLocation) ([]UserLocation, error)
	UpdateLocation(ctx context.Context, ul UserLocation) (UserLocation, error)
	SoftDeleteLocationByID(ctx context.Context, id int64) error
}

type UserService interface {
	CreateProfile(ctx context.Context, u UserCreateDetails) (User, error)
	UpdateProfile(ctx context.Context, u UserUpdateDetails) (User, error)
	GetProfile(ctx context.Context) (User, error)

	AddLocation(ctx context.Context, ul UserLocation) (UserLocation, error)
	UpdateLocation(ctx context.Context, ul UserLocationUpdateDetails) (UserLocation, error)
	DeleteLocationByID(ctx context.Context, id int64) error
}
