package domain

import "context"

type UserLocation struct {
	ID   int64
	User User

	Alias    string
	Name     string
	PhotoURL string
	IsActive bool
}

type UserLocationRepository interface {
	GetByID(ctx context.Context, id int64) (UserLocation, error)
	FindByUserID(ctx context.Context, id int64) ([]UserLocation, error)

	Add(ctx context.Context, ul UserLocation) (UserLocation, error)
	Update(ctx context.Context, ul UserLocation) (UserLocation, error)
}
