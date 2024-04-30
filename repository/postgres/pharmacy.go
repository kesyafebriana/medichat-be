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

func (r *pharmacyRepository) AddOperation(ctx context.Context, pharmacyOperation domain.PharmacyOperationCreateDetails) (domain.PharmacyOperations, error) {
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

func (r *pharmacyRepository) UpdateOperation(ctx context.Context, pharmacyOperation domain.PharmacyOperationsUpdateDetails) (domain.PharmacyOperations, error) {
	q := `
		UPDATE pharmacy_operations
		SET	day = $1,
			start_time = $2,
			end_time = $3,
			pharmacy_id = $4
		WHERE id = $5 RETURNING
	` + pharmacyOperationColumns

	return queryOneFull(
		r.querier, ctx, q,
		ScanPharmacyOperation,
		pharmacyOperation.Day, pharmacyOperation.StartTime,
		pharmacyOperation.EndTime, pharmacyOperation.PharmacyID, pharmacyOperation.ID,
	)
}

func (r *pharmacyRepository) SoftDeleteOperationByID(ctx context.Context, id int64) error {
	q := `
		UPDATE pharmacy_operations
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}

func (r *pharmacyRepository) GetPharmacyOperationsByPharmacyId(ctx context.Context, id int64) ([]domain.PharmacyOperations, error) {
	q := `
	SELECT ` + pharmacyOperationColumns + ` 
	FROM pharmacy_operations
	WHERE pharmacy_id = $1 AND deleted_at IS NULL
	`

	return queryFull(
		r.querier, ctx, q,
		ScanPharmacyOperation, id,
	)
}
