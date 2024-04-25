package domain

import (
	"context"
	"time"
)

type Doctor struct {
	ID             int64
	Account        Account
	Specialization Specialization

	STR            string
	WorkLocation   string
	Gender         string
	PhoneNumber    string
	IsActive       bool
	StartWorkDate  time.Time
	Price          int
	CertificateURL string
}

type DoctorCreateDetails struct {
	Name     string
	PhotoURL string

	SpecializationID int64
	STR              string
	WorkLocation     string
	Gender           string
	PhoneNumber      string
	IsActive         bool
	StartWorkDate    time.Time
	Price            int
	CertificateURL   string
}

type DoctorUpdateDetails struct {
	Name     *string
	PhotoURL *string

	WorkLocation *string
	Gender       *string
	PhoneNumber  *string
	Price        *int
}

func (d *Doctor) ApplyUpdate(det DoctorUpdateDetails) {
	if det.WorkLocation != nil {
		d.WorkLocation = *det.WorkLocation
	}
	if det.Gender != nil {
		d.Gender = *det.Gender
	}
	if det.PhoneNumber != nil {
		d.PhoneNumber = *det.PhoneNumber
	}
	if det.Price != nil {
		d.Price = *det.Price
	}
}

type DoctorRepository interface {
	GetByID(ctx context.Context, id int64) (Doctor, error)
	GetByIDAndLock(ctx context.Context, id int64) (Doctor, error)
	IsExistByID(ctx context.Context, id int64) (bool, error)

	GetByAccountID(ctx context.Context, id int64) (Doctor, error)
	GetByAccountIDAndLock(ctx context.Context, id int64) (Doctor, error)
	IsExistByAccountID(ctx context.Context, id int64) (bool, error)

	Add(ctx context.Context, d Doctor) (Doctor, error)
	Update(ctx context.Context, d Doctor) (Doctor, error)
}

type DoctorService interface {
	CreateProfile(ctx context.Context, det DoctorCreateDetails) (Doctor, error)
	UpdateProfile(ctx context.Context, det DoctorUpdateDetails) (Doctor, error)
	GetProfile(ctx context.Context) (Doctor, error)

	SetActiveStatus(ctx context.Context, active bool) error
}
