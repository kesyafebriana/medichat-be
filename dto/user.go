package dto

import (
	"medichat-be/domain"
	"time"
)

type UserResponse struct {
	ID          int64  `json:"id"`
	DateOfBirth string `json:"date_of_birth"`
}

func NewUserResponse(u domain.User) UserResponse {
	dob := ""
	if u.DateOfBirth != (time.Time{}) {
		dob = u.DateOfBirth.Format("2006-01-02")
	}
	return UserResponse{
		ID:          u.ID,
		DateOfBirth: dob,
	}
}

type UserCreateRequest struct {
	AccountID   int64  `json:"account_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}

func (r *UserCreateRequest) ToDetails() (domain.UserCreateDetails, error) {
	dob, err := time.Parse("2006-01-02", r.DateOfBirth)
	if err != nil {
		return domain.UserCreateDetails{}, err
	}

	ret := domain.UserCreateDetails{
		AccountID:   r.AccountID,
		Name:        r.Name,
		DateOfBirth: dob,
	}

	return ret, nil
}

type UserUpdateRequest struct {
	ID          int64   `json:"id" binding:"required"`
	Name        *string `json:"name"`
	DateOfBirth *string `json:"date_of_birth"`
}

func (r *UserUpdateRequest) ToDetails() (domain.UserUpdateDetails, error) {
	ret := domain.UserUpdateDetails{
		ID:          r.ID,
		Name:        r.Name,
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
