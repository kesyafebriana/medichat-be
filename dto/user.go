package dto

import (
	"medichat-be/domain"
	"time"
)

type UserResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhotoURL    string `json:"photo_url"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

func NewUserResponse(u domain.User) UserResponse {
	dob := ""
	if u.DateOfBirth != (time.Time{}) {
		dob = u.DateOfBirth.Format("2006-01-02")
	}
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		PhotoURL:    u.PhotoURL,
		DateOfBirth: dob,
	}
}

type UserUpdateRequest struct {
	ID          int64   `json:"id" binding:"required"`
	Name        *string `json:"name"`
	PhotoURL    *string `json:"photo_url"`
	DateOfBirth *string `json:"date_of_birth"`
}

func (r *UserUpdateRequest) ToDetails() (domain.UserUpdateDetails, error) {
	ret := domain.UserUpdateDetails{
		ID:          r.ID,
		Name:        r.Name,
		PhotoURL:    r.PhotoURL,
		DateOfBirth: nil,
	}

	if r.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *r.DateOfBirth)
		if err != nil {
			return domain.UserUpdateDetails{}, err
		}
		ret.DateOfBirth = &dob
	}

	return ret, nil
}
