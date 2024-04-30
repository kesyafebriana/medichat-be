package dto

import (
	"medichat-be/domain"
	"medichat-be/util"
	"time"
)

type PharmacyResponse struct {
	ID                 int64                       `json:"id"`
	Name               string                      `json:"name"`
	ManagerID          int64                       `json:"manager_id"`
	Address            string                      `json:"address"`
	Coordinate         CoordinateDTO               `json:"coordinate"`
	PharmacistName     string                      `json:"pharmacist_name"`
	PharmacistLicense  string                      `json:"pharmacist_license"`
	PharmacistPhone    string                      `json:"pharmacist_phone"`
	PharmacyOperations []PharmacyOperationResponse `json:"pharmacy_operations"`
}

func NewPharmacyResponse(pharmacy domain.Pharmacy) PharmacyResponse {
	return PharmacyResponse{
		ID:                 pharmacy.ID,
		ManagerID:          pharmacy.ManagerID,
		Name:               pharmacy.Name,
		Address:            pharmacy.Address,
		Coordinate:         CoordinateDTO(pharmacy.Coordinate),
		PharmacistName:     pharmacy.PharmacistName,
		PharmacistLicense:  pharmacy.PharmacistLicense,
		PharmacistPhone:    pharmacy.PharmacistLicense,
		PharmacyOperations: util.MapSlice(pharmacy.PharmacyOperations, NewPharmacyOperationResponse),
	}
}

type PharmacyOperationResponse struct {
	ID        int64  `json:"id"`
	Day       string `json:"day"`
	StartTime string `json:"start_time`
	EndTime   string `json:"end_time"`
}

func NewPharmacyOperationResponse(pharmacyOperation domain.PharmacyOperations) PharmacyOperationResponse {
	return PharmacyOperationResponse{
		ID:        pharmacyOperation.ID,
		Day:       pharmacyOperation.Day,
		StartTime: pharmacyOperation.StartTime.Format("07:00"),
		EndTime:   pharmacyOperation.EndTime.Format("07:00"),
	}
}

type PharmacyOperationCreateRequest struct {
	Day       string `json:"day" binding:"required,no_leading_trailing_space"`
	StartTime string `json:"start_time" binding:"required,no_leading_trailing_space"`
	EndTime   string `json:"end_time" binding:"required,no_leading_trailing_space"`
}

func (p PharmacyOperationCreateRequest) ToEntity() domain.PharmacyOperations {
	starTime, _ := time.Parse("03:04", p.StartTime)
	endTime, _ := time.Parse("03:04", p.EndTime)

	return domain.PharmacyOperations{
		Day:       p.Day,
		StartTime: starTime,
		EndTime:   endTime,
	}
}

type PharmacyCreateRequest struct {
	Name               string                           `json:"name" binding:"required,no_leading_trailing_space"`
	ManagerID          int64                            `json:"manager_id" binding:"required"`
	Address            string                           `json:"address" binding:"required,no_leading_trailing_space"`
	Coordinate         CoordinateDTO                    `json:"coordinate" binding:"required"`
	PharmacistName     string                           `json:"pharmacist_name" binding:"required,no_leading_trailing_space"`
	PharmacistLicense  string                           `json:"pharmacist_license" binding:"required,no_leading_trailing_space"`
	PharmacistPhone    string                           `json:"pharmacist_phone" binding:"required,no_leading_trailing_space"`
	PharmacyOperations []PharmacyOperationCreateRequest `json:"pharmacy_operations" binding:"required,min=1,dive,required"`
}

func PharmacyCreateToDetails(p PharmacyCreateRequest) domain.PharmacyCreateDetails {
	return domain.PharmacyCreateDetails{
		Name:              p.Name,
		ManagerID:         p.ManagerID,
		Address:           p.Address,
		Coordinate:        p.Coordinate.ToCoordinate(),
		PharmacistName:    p.PharmacistName,
		PharmacistPhone:   p.PharmacistPhone,
		PharmacistLicense: p.PharmacistLicense,
		PharmacyOperations: util.MapSlice(p.PharmacyOperations, func(p PharmacyOperationCreateRequest) domain.PharmacyOperations {
			return p.ToEntity()
		}),
	}
}

type PharmacyUpdateRequest struct {
	Name              *string        `json:"name" binding:"omitempty,no_leading_trailing_space"`
	Address           *string        `json:"address" binding:"omitempty,no_leading_trailing_space"`
	Coordinate        *CoordinateDTO `json:"coordinate" binding:"omitempty"`
	PharmacistName    *string        `json:"pharmacist_name" binding:"omitempty,no_leading_trailing_space"`
	PharmacistLicense *string        `json:"pharmacist_license" binding:"omitempty,no_leading_trailing_space"`
	PharmacistPhone   *string        `json:"pharmacist_phone" binding:"omitempty,no_leading_trailing_space"`
}

func PharmacyUpdateRequestToDetails(p PharmacyUpdateRequest) domain.PharmacyUpdateDetails {
	return domain.PharmacyUpdateDetails{
		Name:              p.Name,
		Address:           p.Address,
		Coordinate:        (*domain.Coordinate)(p.Coordinate),
		PharmacistName:    p.PharmacistName,
		PharmacistLicense: p.PharmacistLicense,
		PharmacistPhone:   p.PharmacistPhone,
	}
}

type PharmacyOperationUpdateRequest struct {
	Day       *string `json:"day" binding:"omitempty,no_leading_trailing_space"`
	StartTime *string `json:"start_time" binding:"omitempty,no_leading_trailing_space"`
	EndTime   *string `json:"end_time" binding:"omitempty,no_leading_trailing_space"`
}
