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
	ID        int64
	ManagerID int64

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
	Address            string
	Coordinate         Coordinate
	PharmacistName     string
	PharmacistLicense  string
	PharmacistPhone    string
	PharmacyOperations []PharmacyOperations
}

type PharmacyOperationCreateDetails struct {
	PharmacyID int64
	Day        string
	StartTime  time.Time
	EndTime    time.Time
}

type PharmacyUpdateDetails struct {
	ID        int64
	ManagerID int64

	Name              *string
	Address           *string
	Coordinate        *Coordinate
	PharmacistName    *string
	PharmacistLicense *string
	PharmacistPhone   *string
}

type PharmacyOperationsUpdateDetails struct {
	ID int64

	Day       *string
	StartTime *time.Time
	EndTime   *time.Time
}

type PharmaciesQuery struct {
	ManagerID *int64
	City      *string
	Page      int64
	Limit     int64
	SortBy    string
	SortType  string
}

type PharmacyRepository interface {
	// GetPharmacies(ctx context.Context, query PharmaciesQuery) ([]Pharmacy, error)
	// GetByID(ctx context.Context, id int64) (Pharmacy, error)
	// GetByName(ctx context.Context, name string) (Pharmacy, error)

	Add(ctx context.Context, pharmacy PharmacyCreateDetails) (Pharmacy, error)
	// Update(ctx context.Context, pharmacy Pharmacy) (Pharmacy, error)
	// SoftDeleteById(ctx context.Context, id int64) error
	// BulkSoftDelete(ctx context.Context, ids []int64) error

	// GetOperationsByPharmacyID(ctx context.Context, id int64) ([]PharmacyOperations, error)
	// GetOperationsByDay(ctx context.Context, day string) (PharmacyOperations, error)

	AddOperation(ctx context.Context, pharmacyOperation PharmacyOperations) (PharmacyOperations, error)
	// AddOperations(ctx context.Context, pharmacyOperations []PharmacyOperations) ([]PharmacyOperations, error)
	// UpdateOperation(ctx context.Context, pharmacyOperation PharmacyOperations) (PharmacyOperations, error)
	// SoftDeleteOperationByID(ctx context.Context, id int64) error
}

type PharmacyService interface {
	CreatePharmacy(ctx context.Context, pharmacy PharmacyCreateDetails) (Pharmacy, error)
	// GetPharmacies(ctx context.Context, query PharmaciesQuery) ([]Pharmacy, error)
	// UpdatePharmacy(ctx context.Context, pharmacy PharmacyUpdateDetails) (Pharmacy, error)
	// DeletePharmacy(ctx context.Context, id int64) error

	// AddOperation(ctx context.Context, pharmacyOperation PharmacyOperationCreateDetails) (PharmacyOperations, error)
	// UpdateOperation(ctx context.Context, pharmacyOperation PharmacyOperationsUpdateDetails) (PharmacyOperations, error)
	// DeleteOperationByID(ctx context.Context, id int64) error
}
