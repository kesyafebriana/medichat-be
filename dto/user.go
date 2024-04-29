package dto

import (
	"medichat-be/domain"
	"medichat-be/util"
	"mime/multipart"
	"time"
)

type UserResponse struct {
	ID             int64                  `json:"id"`
	DateOfBirth    string                 `json:"date_of_birth"`
	MainLocationID int64                  `json:"main_location_id"`
	Locations      []UserLocationResponse `json:"locations,omitempty"`
}

func NewUserResponse(u domain.User) UserResponse {
	return UserResponse{
		ID:             u.ID,
		DateOfBirth:    u.DateOfBirth.Format("2006-01-02"),
		MainLocationID: u.MainLocationID,
		Locations:      util.MapSlice(u.Locations, NewUserLocationResponse),
	}
}

type UserLocationResponse struct {
	ID int64 `json:"id"`

	Alias      string        `json:"alias"`
	Address    string        `json:"address"`
	Coordinate CoordinateDTO `json:"coordinate"`
	IsActive   bool          `json:"is_active"`
}

func NewUserLocationResponse(ul domain.UserLocation) UserLocationResponse {
	return UserLocationResponse{
		ID:         ul.ID,
		Alias:      ul.Alias,
		Address:    ul.Address,
		Coordinate: NewCoordinateDTO(ul.Coordinate),
		IsActive:   ul.IsActive,
	}
}

type UserCreateRequest = MultipartForm[
	struct {
		Photo *multipart.FileHeader `form:"photo" binding:"omitempty,content_type=image/png"`
	},
	struct {
		Name        string                      `json:"name" binding:"required,no_leading_trailing_space"`
		DateOfBirth string                      `json:"date_of_birth" binding:"required,no_leading_trailing_space"`
		Locations   []UserLocationCreateRequest `json:"locations" binding:"required,min=1,dive,required"`
	},
]

func UserCreateRequestToDetails(r UserCreateRequest) (domain.UserCreateDetails, error) {
	dob, err := time.Parse("2006-01-02", r.Data.DateOfBirth)
	if err != nil {
		return domain.UserCreateDetails{}, err
	}

	ret := domain.UserCreateDetails{
		Name:        r.Data.Name,
		DateOfBirth: dob,
		Locations: util.MapSlice(
			r.Data.Locations,
			func(ul UserLocationCreateRequest) domain.UserLocation {
				return ul.ToEntity()
			},
		),
	}

	if r.Form.Photo != nil {
		f, err := r.Form.Photo.Open()
		if err != nil {
			return domain.UserCreateDetails{}, err
		}
		ret.Photo = f
	}

	return ret, nil
}

type UserUpdateRequest = MultipartForm[
	struct {
		Photo *multipart.FileHeader `form:"photo" binding:"omitempty,content_type=image/png"`
	},
	struct {
		Name           *string `json:"name" binding:"omitempty,no_leading_trailing_space"`
		DateOfBirth    *string `json:"date_of_birth" binding:"omitempty,no_leading_trailing_space"`
		MainLocationID *int64  `json:"main_location_id"`
	},
]

func UserUpdateRequestToDetails(r UserUpdateRequest) (domain.UserUpdateDetails, error) {
	ret := domain.UserUpdateDetails{
		Name:           r.Data.Name,
		DateOfBirth:    nil,
		MainLocationID: r.Data.MainLocationID,
	}

	if r.Data.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *r.Data.DateOfBirth)
		if err != nil {
			return domain.UserUpdateDetails{}, err
		}
		ret.DateOfBirth = &dob
	}

	if r.Form.Photo != nil {
		f, err := r.Form.Photo.Open()
		if err != nil {
			return domain.UserUpdateDetails{}, err
		}
		ret.Photo = f
	}

	return ret, nil
}

type UserLocationCreateRequest struct {
	Alias      string        `json:"alias" binding:"required"`
	Address    string        `json:"address" binding:"required"`
	Coordinate CoordinateDTO `json:"coordinate" binding:"required"`
	IsActive   bool          `json:"is_active" binding:"required"`
}

func (r UserLocationCreateRequest) ToEntity() domain.UserLocation {
	return domain.UserLocation{
		Alias:      r.Alias,
		Address:    r.Address,
		Coordinate: r.Coordinate.ToCoordinate(),
		IsActive:   r.IsActive,
	}
}

type UserLocationUpdateRequest struct {
	ID int64 `json:"id" binding:"required"`

	Alias      *string        `json:"alias"`
	Address    *string        `json:"address"`
	Coordinate *CoordinateDTO `json:"coordinate"`
	IsActive   *bool          `json:"is_active"`
}

func (r UserLocationUpdateRequest) ToDetails() domain.UserLocationUpdateDetails {
	return domain.UserLocationUpdateDetails{
		ID:         r.ID,
		Alias:      r.Alias,
		Address:    r.Address,
		Coordinate: (*domain.Coordinate)(r.Coordinate),
		IsActive:   r.IsActive,
	}
}
