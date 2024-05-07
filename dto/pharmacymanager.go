package dto

import (
	"medichat-be/domain"
	"mime/multipart"
)

type PharmacyManagerResponse struct {
	ID int64 `json:"id"`
}

type PharmacyManagerAccountResponse struct {
	AccountID  int64  `json:"account_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	ProfileSet bool   `json:"profile_set"`
}

func NewPharmacyManagerResponse(p domain.PharmacyManager) PharmacyManagerResponse {
	return PharmacyManagerResponse{
		ID: p.ID,
	}
}
func NewPharmacyManagersAccountResponse(p []domain.Account) []PharmacyManagerAccountResponse {
	var res []PharmacyManagerAccountResponse

	for _, v := range p {
		res = append(res, PharmacyManagerAccountResponse{
			AccountID:  v.ID,
			Email:      v.Email,
			Name:       v.Name,
			ProfileSet: v.ProfileSet,
		})
	}

	return res
}

type PharmacyManagerCreateRequest = MultipartForm[
	struct {
		Photo *multipart.FileHeader `form:"photo" binding:"omitempty,content_type=image/png"`
	},
	struct {
		Name string `json:"name" binding:"required,no_leading_trailing_space"`
	},
]

func PharmacyManagerCreateRequestToDetails(r PharmacyManagerCreateRequest) (domain.PharmacyManagerCreateDetails, error) {
	ret := domain.PharmacyManagerCreateDetails{
		Name: r.Data.Name,
	}

	if r.Form.Photo != nil {
		f, err := r.Form.Photo.Open()
		if err != nil {
			return domain.PharmacyManagerCreateDetails{}, err
		}
		ret.Photo = f
	}

	return ret, nil
}
