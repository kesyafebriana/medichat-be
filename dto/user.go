package dto

import (
	"medichat-be/domain"
	"mime/multipart"
	"time"
)

type UserResponse struct {
	ID          int64                  `json:"id"`
	DateOfBirth string                 `json:"date_of_birth"`
	Locations   []UserLocationResponse `json:"locations,omitempty"`
}

func NewUserResponse(u domain.User) UserResponse {
	dob := ""
	if u.DateOfBirth != (time.Time{}) {
		dob = u.DateOfBirth.Format("2006-01-02")
	}

	var locations []UserLocationResponse
	if u.Locations != nil {
		locations = make([]UserLocationResponse, len(u.Locations))
		for i, ul := range u.Locations {
			locations[i] = NewUserLocationResponse(ul)
		}
	}

	return UserResponse{
		ID:          u.ID,
		DateOfBirth: dob,
		Locations:   locations,
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
		Photo *multipart.FileHeader `form:"photo"`
	},
	struct {
		AccountID   int64  `json:"account_id" binding:"required"`
		Name        string `json:"name" binding:"required,no_leading_trailing_space"`
		DateOfBirth string `json:"date_of_birth" binding:"required,no_leading_trailing_space"`
	},
]

func UserCreateRequestToDetails(r UserCreateRequest) (domain.UserCreateDetails, error) {
	dob, err := time.Parse("2006-01-02", r.Data.DateOfBirth)
	if err != nil {
		return domain.UserCreateDetails{}, err
	}

	ret := domain.UserCreateDetails{
		AccountID:   r.Data.AccountID,
		Name:        r.Data.Name,
		DateOfBirth: dob,
	}

	return ret, nil
}

type UserUpdateRequest = MultipartForm[
	struct {
		Photo *multipart.FileHeader `form:"photo"`
	},
	struct {
		ID          int64   `json:"id" binding:"required"`
		Name        *string `json:"name" binding:"omitempty,no_leading_trailing_space"`
		DateOfBirth *string `json:"date_of_birth" binding:"omitempty,no_leading_trailing_space"`
	},
]

func UserUpdateRequestToDetails(r UserUpdateRequest) (domain.UserUpdateDetails, error) {
	ret := domain.UserUpdateDetails{
		ID:          r.Data.ID,
		Name:        r.Data.Name,
		DateOfBirth: nil,
	}

	if r.Data.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *r.Data.DateOfBirth)
		if err != nil {
			return domain.UserUpdateDetails{}, err
		}
		ret.DateOfBirth = &dob
	}

	return ret, nil
}
