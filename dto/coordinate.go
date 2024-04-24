package dto

import "medichat-be/domain"

type CoordinateDTO struct {
	Longitude float64 `json:"lon" binding:"required,gte=-180.0,lte=180.0"`
	Latitude  float64 `json:"lat" binding:"required,gte=-90.0,lte=90.0"`
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
