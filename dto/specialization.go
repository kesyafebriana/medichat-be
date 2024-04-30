package dto

import "medichat-be/domain"

type SpecializationResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewSpecializationResponse(s domain.Specialization) SpecializationResponse {
	return SpecializationResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}
