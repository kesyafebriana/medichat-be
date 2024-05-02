package domain

import (
	"context"
	"time"
)

const (
	PharmacySortById        = "id"
	PharmacySortByName      = "name"
	PharmacySortByManagerId = "manager"
)

type PharmacyOperations struct {
	ID         int64
	PharmacyID int64

	Day       string
	StartTime time.Time
	EndTime   time.Time
}

type Pharmacy struct {
	ID                 int64
	ManagerID          int64
	Slug               string
	Name               string
	Address            string
	Coordinate         Coordinate
	PharmacistName     string
	PharmacistLicense  string
	PharmacistPhone    string
	PharmacyOperations []PharmacyOperations
}

type PharmacyCreateDetails struct {
	Name               string
	ManagerID          int64
	Slug               string
	Address            string
	Coordinate         Coordinate
	PharmacistName     string
	PharmacistLicense  string
	PharmacistPhone    string
	PharmacyOperations []PharmacyOperationCreateDetails
}

type PharmacyOperationCreateDetails struct {
	Slug       string
	PharmacyID int64
	Day        string
	StartTime  time.Time
	EndTime    time.Time
}

type PharmacyUpdateDetails struct {
	ID                int64
	ManagerID         int64
	Slug              string
	Name              string
	Address           string
	Coordinate        Coordinate
	PharmacistName    string
	PharmacistLicense string
	PharmacistPhone   string
}

type PharmacyOperationsUpdateDetails struct {
	ID         int64
	PharmacyId int64
	Slug       string
	Day        string
	StartTime  time.Time
	EndTime    time.Time
}

type PharmaciesQuery struct {
	ManagerID   *int64
	Day         *string
	StartTime   *string
	EndTime     *string
	ProductSlug *string
	Longitude   *float64
	Latitude    *float64
	Name        *string
	Page        int64
	Limit       int64
	SortBy      string
	SortType    string
}

type PharmacyRepository interface {
	GetPharmacies(ctx context.Context, query PharmaciesQuery) ([]Pharmacy, error)
	GetBySlug(ctx context.Context, slug string) (Pharmacy, error)

	Add(ctx context.Context, pharmacy PharmacyCreateDetails) (Pharmacy, error)
	Update(ctx context.Context, pharmacy PharmacyUpdateDetails) (Pharmacy, error)
	SoftDeleteBySlug(ctx context.Context, slug string) error

	GetPharmacyOperationsByPharmacyId(ctx context.Context, id int64) ([]PharmacyOperations, error)
	GetPharmacyOperationsByPharmacyIdAndLock(ctx context.Context, id int64) ([]PharmacyOperations, error)
	AddOperation(ctx context.Context, pharmacyOperation PharmacyOperationCreateDetails) (PharmacyOperations, error)
	UpdateOperation(ctx context.Context, pharmacyOperation PharmacyOperationsUpdateDetails) (PharmacyOperations, error)
	SoftDeleteOperationByID(ctx context.Context, id int64) error
}

type PharmacyService interface {
	CreatePharmacy(ctx context.Context, pharmacy PharmacyCreateDetails) (Pharmacy, error)
	GetPharmacies(ctx context.Context, query PharmaciesQuery) ([]Pharmacy, error)
	GetPharmacyBySlug(ctx context.Context, slug string) (Pharmacy, error)
	UpdatePharmacy(ctx context.Context, pharmacy PharmacyUpdateDetails) (Pharmacy, error)
	DeletePharmacyBySlug(ctx context.Context, slug string) error

	GetOperationsBySlug(ctx context.Context, slug string) ([]PharmacyOperations, error)
	AddOperation(ctx context.Context, pharmacyOperation PharmacyOperationCreateDetails) (PharmacyOperations, error)
	UpdateOperations(ctx context.Context, pharmacyOperation []PharmacyOperationsUpdateDetails) ([]PharmacyOperations, error)
}
