package domain

import "time"

const (
	PharmacySortById        = "id"
	PharmacySortByName      = "name"
	PharmacySortByManagerId = "manager"
)

type PharmacyDetail struct {
	ID int64

	PharmacistName    string
	PharmacistLicense string
	PharmacistPhone   string
}

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
	PharmacyDetail     PharmacyDetail
	PharmacyOperations []PharmacyOperations
}

type PharmacyCreateDetails struct {
	Name               string
	Address            string
	Coordinate         Coordinate
	PharmacyDetail     PharmacyDetail
	PharmacyOperations []PharmacyOperations
}

type PharmacyUpdateDetails struct {
	ID int64

	Name       *string
	Address    *string
	Coordinate *Coordinate
}

type PharmacyDetailUpdateDetails struct {
	ID int64

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
