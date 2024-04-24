package dto

import "medichat-be/domain"

type CoordinateDTO struct {
	Longitude float64 `json:"lon" binding:"required,min=-180,max=180"`
	Latitude  float64 `json:"lat" binding:"required,min=-90,max=90"`
}

func NewCoordinateDTO(c domain.Coordinate) CoordinateDTO {
	return CoordinateDTO{
		Longitude: c.Longitude,
		Latitude:  c.Latitude,
	}
}

func (c CoordinateDTO) ToCoordinate() domain.Coordinate {
	return domain.Coordinate{
		Longitude: c.Longitude,
		Latitude:  c.Latitude,
	}
}
