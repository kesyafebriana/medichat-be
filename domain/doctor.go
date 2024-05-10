package domain

import (
	"context"
	"mime/multipart"
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
	YearExperience int
	Price          int
	CertificateURL string
}

type DoctorCreateDetails struct {
	Name  string
	Photo multipart.File

	SpecializationID int64
	STR              string
	WorkLocation     string
	Gender           string
	PhoneNumber      string
	IsActive         bool
	StartWorkDate    time.Time
	Price            int
	Certificate      multipart.File
}

type DoctorUpdateDetails struct {
	Name  *string
	Photo multipart.File

	WorkLocation *string
	Gender       *string
	PhoneNumber  *string
	Price        *int
}

type DoctorListDetails struct {
	SpecializationID  *int64
	Name              *string
	Gender            *string
	MinPrice          *int
	MaxPrice          *int
	MinYearExperience *int

	SortBy  string
	SortAsc bool

	Cursor   any
	CursorID *int64
	Limit    int
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
	List(ctx context.Context, det DoctorListDetails) ([]Doctor, error)

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
	List(ctx context.Context, det DoctorListDetails) ([]Doctor, error)

	GetByID(ctx context.Context, id int64) (Doctor, error)

	CreateProfile(ctx context.Context, det DoctorCreateDetails) (Doctor, error)
	UpdateProfile(ctx context.Context, det DoctorUpdateDetails) (Doctor, error)
	GetProfile(ctx context.Context) (Doctor, error)

	SetActiveStatus(ctx context.Context, active bool) error
}
