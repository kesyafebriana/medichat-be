package dto

import (
	"medichat-be/domain"
	"medichat-be/util"
	"time"
)

type PharmacyResponse struct {
	ID                      int64                            `json:"id"`
	Name                    string                           `json:"name"`
	Slug                    string                           `json:"slug"`
	ManagerID               int64                            `json:"manager_id"`
	Address                 string                           `json:"address"`
	Coordinate              CoordinateDTO                    `json:"coordinate"`
	PharmacistName          string                           `json:"pharmacist_name"`
	PharmacistLicense       string                           `json:"pharmacist_license"`
	PharmacistPhone         string                           `json:"pharmacist_phone"`
	PharmacyOperations      []PharmacyOperationResponse      `json:"pharmacy_operations"`
	PharmacyShipmentMethods []PharmacyShipmentMethodResponse `json:"pharmacy_shipment_methods"`

	Distance *float64 `json:"distance,omitempty"`
}

type PharmacyResponseWithStock struct {
	ID                      int64                            `json:"id"`
	Name                    string                           `json:"name"`
	Slug                    string                           `json:"slug"`
	ManagerID               int64                            `json:"manager_id"`
	Address                 string                           `json:"address"`
	Coordinate              CoordinateDTO                    `json:"coordinate"`
	PharmacistName          string                           `json:"pharmacist_name"`
	PharmacistLicense       string                           `json:"pharmacist_license"`
	PharmacistPhone         string                           `json:"pharmacist_phone"`
	PharmacyOperations      []PharmacyOperationResponse      `json:"pharmacy_operations"`
	PharmacyShipmentMethods []PharmacyShipmentMethodResponse `json:"pharmacy_shipment_methods"`
	StockInfo               StockResponse                    `json:"stock"`

	Distance *float64 `json:"distance,omitempty"`
}

type PharmaciesResponse struct {
	Pharmacies []PharmacyResponse `json:"pharmacies"`
	PageInfo   PageInfoResponse   `json:"page_info"`
}

type PharmaciesStockResponse struct {
	Pharmacies []PharmacyResponseWithStock `json:"pharmacies"`
	PageInfo   PageInfoResponse            `json:"page_info"`
}

type PharmacySlugParams struct {
	Slug string `uri:"slug" binding:"required"`
}

func NewPharmacyResponse(pharmacy domain.Pharmacy) PharmacyResponse {
	return PharmacyResponse{
		ID:                      pharmacy.ID,
		ManagerID:               pharmacy.ManagerID,
		Slug:                    pharmacy.Slug,
		Name:                    pharmacy.Name,
		Address:                 pharmacy.Address,
		Coordinate:              CoordinateDTO(pharmacy.Coordinate),
		PharmacistName:          pharmacy.PharmacistName,
		PharmacistLicense:       pharmacy.PharmacistLicense,
		PharmacistPhone:         pharmacy.PharmacistLicense,
		PharmacyOperations:      util.MapSlice(pharmacy.PharmacyOperations, NewPharmacyOperationResponse),
		PharmacyShipmentMethods: util.MapSlice(pharmacy.PharmacyShipmentMethods, NewPharmacyShipmentMethodResponse),

		Distance: pharmacy.Distance,
	}
}

func NewPharmacyWithStockResponse(pharmacy domain.PharmacyStock) PharmacyResponseWithStock {
	return PharmacyResponseWithStock{
		ID:                      pharmacy.ID,
		ManagerID:               pharmacy.ManagerID,
		Slug:                    pharmacy.Slug,
		Name:                    pharmacy.Name,
		Address:                 pharmacy.Address,
		Coordinate:              CoordinateDTO(pharmacy.Coordinate),
		PharmacistName:          pharmacy.PharmacistName,
		PharmacistLicense:       pharmacy.PharmacistLicense,
		PharmacistPhone:         pharmacy.PharmacistLicense,
		PharmacyOperations:      util.MapSlice(pharmacy.PharmacyOperations, NewPharmacyOperationResponse),
		PharmacyShipmentMethods: util.MapSlice(pharmacy.PharmacyShipmentMethods, NewPharmacyShipmentMethodResponse),
		StockInfo:               NewStockResponse(pharmacy.Stock),

		Distance: pharmacy.Distance,
	}
}

func NewPharmaciesResponse(pharmacy []domain.Pharmacy, pageInfo domain.PageInfo) PharmaciesResponse {
	res := PharmaciesResponse{
		Pharmacies: []PharmacyResponse{},
		PageInfo:   PageInfoResponse{},
	}

	for _, v := range pharmacy {
		res.Pharmacies = append(res.Pharmacies, NewPharmacyResponse(v))
	}

	res.PageInfo = NewPageInfoResponse(pageInfo)

	return res
}

func NewPharmaciesStockResponse(pharmacy []domain.PharmacyStock, pageInfo domain.PageInfo) PharmaciesStockResponse {
	var res PharmaciesStockResponse

	for _, v := range pharmacy {
		res.Pharmacies = append(res.Pharmacies, NewPharmacyWithStockResponse(v))
	}

	res.PageInfo = NewPageInfoResponse(pageInfo)

	return res
}

type PharmacyOperationResponse struct {
	ID        int64  `json:"id"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func NewPharmacyOperationResponse(pharmacyOperation domain.PharmacyOperations) PharmacyOperationResponse {
	return PharmacyOperationResponse{
		ID:        pharmacyOperation.ID,
		Day:       pharmacyOperation.Day,
		StartTime: pharmacyOperation.StartTime.Format("15:04"),
		EndTime:   pharmacyOperation.EndTime.Format("15:04"),
	}
}

func NewPharmacyOperationsResponse(pharmacyOperations []domain.PharmacyOperations) []PharmacyOperationResponse {
	var res []PharmacyOperationResponse

	for _, v := range pharmacyOperations {
		res = append(res, PharmacyOperationResponse{
			ID:        v.ID,
			Day:       v.Day,
			StartTime: v.StartTime.Format("15:04"),
			EndTime:   v.EndTime.Format("15:04"),
		})
	}

	return res
}

type PharmacyShipmentMethodResponse struct {
	ID               int64  `json:"id"`
	PharmacyID       int64  `json:"pharmacy_id"`
	ShipmentMethodID int64  `json:"shipment_method_id"`
	ShipmentMethod   string `json:"shipment_method"`
}

func NewPharmacyShipmentMethodResponse(shipmentMethod domain.PharmacyShipmentMethods) PharmacyShipmentMethodResponse {
	return PharmacyShipmentMethodResponse{
		ID:               shipmentMethod.ID,
		PharmacyID:       shipmentMethod.PharmacyID,
		ShipmentMethodID: shipmentMethod.ShipmentMethodID,
		ShipmentMethod:   *shipmentMethod.Name,
	}
}

func NewPharmacyShipmentMethodsResponse(shipmentMethod []domain.PharmacyShipmentMethods) []PharmacyShipmentMethodResponse {
	var res []PharmacyShipmentMethodResponse

	for _, v := range shipmentMethod {
		res = append(res, PharmacyShipmentMethodResponse{
			ID:               v.ID,
			PharmacyID:       v.PharmacyID,
			ShipmentMethodID: v.ShipmentMethodID,
			ShipmentMethod:   *v.Name,
		})
	}

	return res
}

type PharmacyOperationCreateRequest struct {
	Day       string `json:"day" binding:"required,no_leading_trailing_space"`
	StartTime string `json:"start_time" binding:"required,no_leading_trailing_space"`
	EndTime   string `json:"end_time" binding:"required,no_leading_trailing_space"`
}

type PharmacyShipmentMethodCreateRequest struct {
	ShipmentID int64 `json:"shipment_method_id" binding:"required,no_leading_trailing_space"`
}

func (p PharmacyOperationCreateRequest) ToEntity() domain.PharmacyOperationCreateDetails {
	starTime, _ := time.Parse("15:04", p.StartTime)
	endTime, _ := time.Parse("15:04", p.EndTime)

	return domain.PharmacyOperationCreateDetails{
		Day:       p.Day,
		StartTime: starTime,
		EndTime:   endTime,
	}
}

func (sh PharmacyShipmentMethodCreateRequest) ToEntity() domain.PharmacyShipmentMethodsCreateDetails {
	return domain.PharmacyShipmentMethodsCreateDetails{
		ShipmentMethodID: sh.ShipmentID,
	}
}

type PharmacyCreateRequest struct {
	Name                    string                                `json:"name" binding:"required,no_leading_trailing_space"`
	Address                 string                                `json:"address" binding:"required,no_leading_trailing_space"`
	Coordinate              CoordinateDTO                         `json:"coordinate" binding:"required"`
	PharmacistName          string                                `json:"pharmacist_name" binding:"required,no_leading_trailing_space"`
	PharmacistLicense       string                                `json:"pharmacist_license" binding:"required,no_leading_trailing_space"`
	PharmacistPhone         string                                `json:"pharmacist_phone" binding:"required,no_leading_trailing_space"`
	PharmacyOperations      []PharmacyOperationCreateRequest      `json:"pharmacy_operations" binding:"required,min=1,dive,required"`
	PharmacyShipmentMethods []PharmacyShipmentMethodCreateRequest `json:"pharmacy_shipment_methods" binding:"required,min=1,dive,required"`
}

func PharmacyCreateToDetails(p PharmacyCreateRequest) domain.PharmacyCreateDetails {
	return domain.PharmacyCreateDetails{
		Name:              p.Name,
		Address:           p.Address,
		Coordinate:        p.Coordinate.ToCoordinate(),
		PharmacistName:    p.PharmacistName,
		PharmacistPhone:   p.PharmacistPhone,
		PharmacistLicense: p.PharmacistLicense,
		PharmacyOperations: util.MapSlice(p.PharmacyOperations, func(p PharmacyOperationCreateRequest) domain.PharmacyOperationCreateDetails {
			return p.ToEntity()
		}),
		PharmacyShipmentMethods: util.MapSlice(p.PharmacyShipmentMethods, func(sh PharmacyShipmentMethodCreateRequest) domain.PharmacyShipmentMethodsCreateDetails {
			return sh.ToEntity()
		}),
	}
}

type PharmacyUpdateRequest struct {
	Name              string        `json:"name" binding:"omitempty,no_leading_trailing_space"`
	Address           string        `json:"address" binding:"omitempty,no_leading_trailing_space"`
	Coordinate        CoordinateDTO `json:"coordinate" binding:"omitempty"`
	PharmacistName    string        `json:"pharmacist_name" binding:"omitempty,no_leading_trailing_space"`
	PharmacistLicense string        `json:"pharmacist_license" binding:"omitempty,no_leading_trailing_space"`
	PharmacistPhone   string        `json:"pharmacist_phone" binding:"omitempty,no_leading_trailing_space"`
}

func PharmacyUpdateRequestToDetails(p PharmacyUpdateRequest, slug string) domain.PharmacyUpdateDetails {
	return domain.PharmacyUpdateDetails{
		Slug:              slug,
		Name:              p.Name,
		Address:           p.Address,
		Coordinate:        (domain.Coordinate)(p.Coordinate),
		PharmacistName:    p.PharmacistName,
		PharmacistLicense: p.PharmacistLicense,
		PharmacistPhone:   p.PharmacistPhone,
	}
}

type PharmacyOperationUpdateRequest struct {
	Day       string `json:"day" binding:"omitempty,no_leading_trailing_space"`
	StartTime string `json:"start_time" binding:"omitempty,no_leading_trailing_space"`
	EndTime   string `json:"end_time" binding:"omitempty,no_leading_trailing_space"`
}

func PharmacyOperationRequestToDetails(p PharmacyOperationUpdateRequest, slug string) domain.PharmacyOperationsUpdateDetails {
	starTime, _ := time.Parse("15:04", p.StartTime)
	endTime, _ := time.Parse("15:04", p.EndTime)

	return domain.PharmacyOperationsUpdateDetails{
		Slug:      slug,
		Day:       p.Day,
		StartTime: starTime,
		EndTime:   endTime,
	}
}

type PharmacyShipmentMethodUpdateRequest struct {
	Slug             string `json:"slug" binding:"omitempty"`
	PharmacyID       int64  `json:"pharmacy_id" binding:"omitempty"`
	ShipmentMethodID int64  `json:"shipment_method_id" binding:"omitempty,no_leading_trailing_space"`
}

func PharmacyShipmentMethodRequestToDetails(p PharmacyShipmentMethodUpdateRequest, slug string) domain.PharmacyShipmentMethodsUpdateDetails {
	return domain.PharmacyShipmentMethodsUpdateDetails{
		PharmacyID:       p.PharmacyID,
		ShipmentMethodID: p.ShipmentMethodID,
		Slug:             slug,
	}
}

type PharmacyListQuery struct {
	ManagerID   *int64   `form:"manager_id"`
	Name        *string  `form:"name"`
	Term        *string  `form:"term"`
	Day         *string  `form:"day"`
	StartTime   *string  `form:"start_time"`
	EndTime     *string  `form:"end_time"`
	Longitude   *float64 `form:"long"`
	Latitude    *float64 `form:"lat"`
	SortBy      *string  `form:"sort_by"`
	Sort        *string  `form:"sort"`
	Limit       *int     `form:"limit"`
	Page        *int     `form:"page"`
	IsOpen      *bool    `form:"is_open"`
	ProductSlug *string  `form:"product_slug"`
}

func (p PharmacyListQuery) ToDetails() (domain.PharmaciesQuery, error) {
	query := domain.PharmaciesQuery{
		Name:        p.Name,
		Term:        p.Term,
		Day:         p.Day,
		ManagerID:   p.ManagerID,
		Longitude:   p.Longitude,
		Latitude:    p.Latitude,
		Limit:       10,
		Page:        1,
		SortBy:      "name",
		SortType:    "asc",
		IsOpen:      p.IsOpen,
		ProductSlug: p.ProductSlug,
	}

	if p.Sort != nil {
		query.SortType = *p.Sort
	}

	if p.SortBy != nil {
		query.SortBy = *p.SortBy
	}

	if p.Limit != nil {
		query.Limit = *p.Limit
	}

	if p.Page != nil {
		query.Page = *p.Page
	}

	if p.StartTime != nil {
		_, err := time.Parse("15:04", *p.StartTime)
		if err != nil {
			return domain.PharmaciesQuery{}, err
		}

		query.StartTime = p.StartTime
	}

	if p.EndTime != nil {
		_, err := time.Parse("15:04", *p.EndTime)
		if err != nil {
			return domain.PharmaciesQuery{}, err
		}

		query.EndTime = p.EndTime
	}

	return query, nil
}
