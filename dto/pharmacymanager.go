package dto

import (
	"medichat-be/constants"
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
	PhotoUrl   string `json:"photo_url"`
	ProfileSet bool   `json:"profile_set"`
}

type PharmacyManagerIdParams struct {
	Id int64 `uri:"id" binding:"required"`
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
			PhotoUrl:   v.PhotoURL,
		})
	}

	return res
}

type PharmacyManagersResponse struct {
	PharmacyManagers []PharmacyManagerAccountResponse `json:"pharmacy_managers"`
	PageInfo         PageInfoResponse                 `json:"page_info"`
}

func NewPharmacyManagersWithPage(p []domain.Account, pI domain.PageInfo) PharmacyManagersResponse {
	return PharmacyManagersResponse{
		PharmacyManagers: NewPharmacyManagersAccountResponse(p),
		PageInfo:         NewPageInfoResponse(pI),
	}
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

type GetPharmacyManagerQuery struct {
	Page       int64   `form:"page" binding:"numeric,omitempty,min=1"`
	Limit      int64   `form:"limit" binding:"numeric,omitempty,min=1"`
	Level      int64   `form:"level" binding:"numeric,omitempty,oneof=1 2"`
	SortBy     string  `form:"sort_by" binding:"omitempty,oneof=created_at"`
	SortType   string  `form:"sort_type" binding:"omitempty,oneof=ASC DESC"`
	Term       string  `form:"term"`
	ProfileSet *string `form:"profile_set"`
}

func (q *GetPharmacyManagerQuery) ToPharmacyManagerQuery() domain.PharmacyManagerQuery {
	var page int64 = q.Page
	var sortBy string = q.SortBy
	var sortType string = q.SortType

	if q.Page == 0 || q.Limit == 0 {
		page = 1
	}
	if q.SortBy == "" {
		sortBy = domain.PharmacyManagerSortByCreatedAt
	}
	if q.SortType == "" {
		sortType = constants.SortDesc
	}

	if q.SortBy == domain.PharmacyManagerSortByCreatedAt {
		if sortType == "ASC" {
			sortType = constants.SortAsc
		} else {
			sortType = constants.SortDesc
		}
	}

	return domain.PharmacyManagerQuery{
		Page:       page,
		Limit:      q.Limit,
		Level:      q.Level,
		Term:       q.Term,
		SortBy:     sortBy,
		SortType:   sortType,
		ProfileSet: q.ProfileSet,
	}
}
