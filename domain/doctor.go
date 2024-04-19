package domain

import (
	"context"
	"time"
)

type Doctor struct {
	ID             int64
	Specialization Specialization

	STR           string
	WorkLocation  string
	Gender        string
	PhoneNumber   string
	IsActive      bool
	StartWorkDate time.Time
	Price         int64
}

type DoctorRepository interface {
	GetByID(ctx context.Context, id int64) (Doctor, error)
	GetByIDAndLock(ctx context.Context, id int64) (Doctor, error)

	GetByAccountID(ctx context.Context, id int64) (Doctor, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (Doctor, error)

	Add(ctx context.Context, d Doctor) (Doctor, error)
	Update(ctx context.Context, d Doctor) (Doctor, error)
}
