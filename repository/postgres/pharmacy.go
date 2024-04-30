package postgres

import (
	"context"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
)

type pharmacyRepository struct {
	querier Querier
}

func (r *pharmacyRepository) Add(ctx context.Context, pharmacy domain.PharmacyCreateDetails) (domain.Pharmacy, error) {
	q := `
		INSERT INTO pharmacies(name, manager_id, address, coordinate, 
		pharmacist_name, pharmacist_license, pharmacist_phone)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING
	` + pharmacyColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy,
		pharmacy.Name, pharmacy.ManagerID, pharmacy.Address,
		postgis.NewPointFromCoordinate(pharmacy.Coordinate),
		pharmacy.PharmacistName, pharmacy.PharmacistLicense, pharmacy.PharmacistPhone,
	)
}

func (r *pharmacyRepository) AddOperation(ctx context.Context, pharmacyOperation domain.PharmacyOperations) (domain.PharmacyOperations, error) {
	q := `
		INSERT INTO pharmacy_operations(pharmacy_id, day, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING
	` + pharmacyOperationColumns

	return queryOneFull(
		r.querier, ctx, q,
		ScanPharmacyOperation,
		pharmacyOperation.PharmacyID, pharmacyOperation.Day,
		pharmacyOperation.StartTime, pharmacyOperation.EndTime,
	)
}
