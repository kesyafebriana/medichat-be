package dto

import (
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
		Price:          d.Price,
		CertificateURL: d.CertificateURL,
	}
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
