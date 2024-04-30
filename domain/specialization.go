package domain

import "context"

type Specialization struct {
	ID   int64
	Name string
}

type SpecializationRepository interface {
	GetAll(ctx context.Context) ([]Specialization, error)
	GetByID(ctx context.Context, id int64) (Specialization, error)
}

type SpecializationService interface {
	GetAll(ctx context.Context) ([]Specialization, error)
}
