package dto

import (
	"medichat-be/constants"
	"medichat-be/domain"
	"mime/multipart"
	"time"
)

type DoctorResponse struct {
	ID             int64                  `json:"id"`
	Specialization SpecializationResponse `json:"specialization"`

	STR            string `json:"str"`
	WorkLocation   string `json:"work_location"`
	Gender         string `json:"gender"`
	PhoneNumber    string `json:"phone_number"`
	IsActive       bool   `json:"is_active"`
	StartWorkDate  string `json:"start_working_date"`
	YearExperience int    `json:"year_experience"`
	Price          int    `json:"price"`
	CertificateURL string `json:"certificate_url"`
}

func NewDoctorResponse(d domain.Doctor) DoctorResponse {
	return DoctorResponse{
		ID:             d.ID,
		Specialization: NewSpecializationResponse(d.Specialization),
		STR:            d.STR,
		WorkLocation:   d.WorkLocation,
		Gender:         d.Gender,
		PhoneNumber:    d.PhoneNumber,
		IsActive:       d.IsActive,
		StartWorkDate:  d.StartWorkDate.Format("2006-01-02"),
		YearExperience: d.YearExperience,
		Price:          d.Price,
		CertificateURL: d.CertificateURL,
	}
}

type DoctorListQuery struct {
	SpecializationID  *int64  `form:"specialization_id"`
	Name              *string `form:"name"`
	Gender            *string `form:"gender"`
	MinPrice          *int    `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice          *int    `form:"max_price" binding:"omitempty,min=0"`
	MinYearExperience *int    `form:"min_year_experience" binding:"omitempty,min=0"`

	SortBy *string `form:"sort_by" binding:"omitempty,doctor_sort_by"`
	Sort   *string `form:"sort" binding:"omitempty,sort_order"`

	Cursor *string `form:"cursor"`
	Limit  *int    `form:"limit" binding:"omitempty,min=1"`
}

func (q *DoctorListQuery) ToDetails() domain.DoctorListDetails {
	ret := domain.DoctorListDetails{
		SpecializationID:  q.SpecializationID,
		Name:              q.Name,
		Gender:            q.Gender,
		MinPrice:          q.MinPrice,
		MaxPrice:          q.MaxPrice,
		MinYearExperience: q.MinYearExperience,

		SortBy:  constants.DoctorSortByName,
		SortAsc: true,

		Cursor: q.Cursor,
		Limit:  10,
	}

	if q.SortBy != nil {
		ret.SortBy = *q.SortBy
	}
	if q.Sort != nil && *q.Sort == constants.SortDesc {
		ret.SortAsc = false
	}
	if q.Limit != nil {
		ret.Limit = *q.Limit
	}

	return ret
}

type DoctorCreateRequest = MultipartForm[
	struct {
		Photo       *multipart.FileHeader `form:"photo"`
		Certificate *multipart.FileHeader `form:"certificate" binding:"required"`
	},
	struct {
		Name             string `json:"name" binding:"required,no_leading_trailing_space"`
		SpecializationID int64  `json:"specialization_id" binding:"required"`
		STR              string `json:"str" binding:"required,no_leading_trailing_space"`
		WorkLocation     string `json:"work_location" binding:"required,no_leading_trailing_space"`
		Gender           string `json:"gender" binding:"required,no_leading_trailing_space"`
		PhoneNumber      string `json:"phone_number" binding:"required,no_leading_trailing_space"`
		IsActive         bool   `json:"is_active"`
		StartWorkDate    string `json:"start_work_date" binding:"required,no_leading_trailing_space"`
		Price            int    `json:"price" binding:"min=0,max=10000000"`
	},
]

func DoctorCreateRequestToDetails(r DoctorCreateRequest) (domain.DoctorCreateDetails, error) {
	d := r.Data
	ret := domain.DoctorCreateDetails{
		Name:             d.Name,
		SpecializationID: d.SpecializationID,
		STR:              d.STR,
		WorkLocation:     d.WorkLocation,
		Gender:           d.Gender,
		PhoneNumber:      d.PhoneNumber,
		IsActive:         d.IsActive,
		Price:            d.Price,
	}

	swd, err := time.Parse("2006-01-02", d.StartWorkDate)
	if err != nil {
		return domain.DoctorCreateDetails{}, err
	}
	ret.StartWorkDate = swd

	return ret, nil
}

type DoctorUpdateRequest = MultipartForm[
	struct {
		Photo *multipart.FileHeader `form:"photo"`
	},
	struct {
		Name         *string `json:"name" binding:"omitempty,no_leading_trailing_space"`
		WorkLocation *string `json:"work_location" binding:"omitempty,no_leading_trailing_space"`
		Gender       *string `json:"gender" binding:"omitempty,no_leading_trailing_space"`
		PhoneNumber  *string `json:"phone_number" binding:"omitempty,no_leading_trailing_space"`
		Price        *int    `json:"price"`
	},
]

func DoctorUpdateRequestToDetails(r DoctorUpdateRequest) (domain.DoctorUpdateDetails, error) {
	d := r.Data
	ret := domain.DoctorUpdateDetails{
		Name:         d.Name,
		WorkLocation: d.WorkLocation,
		Gender:       d.Gender,
		PhoneNumber:  d.PhoneNumber,
		Price:        d.Price,
	}

	return ret, nil
}

type DoctorSetActiveRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}
