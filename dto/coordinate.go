package dto

import "medichat-be/domain"

type CoordinateDTO struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
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
