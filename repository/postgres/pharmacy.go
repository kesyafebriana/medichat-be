package postgres

import (
	"context"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
)

type pharmacyRepository struct {
	querier Querier
}

func (r *pharmacyRepository) GetBySlug(ctx context.Context, slug string) (domain.Pharmacy, error) {
	q := `
		SELECT ` + pharmacyColumns + `FROM pharmacies
		WHERE slug = $1
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy, slug,
	)
}

func (r *pharmacyRepository) Add(ctx context.Context, pharmacy domain.PharmacyCreateDetails) (domain.Pharmacy, error) {
	q := `
		INSERT INTO pharmacies(name, manager_id, address, coordinate, 
		pharmacist_name, pharmacist_license, pharmacist_phone, slug)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
	` + pharmacyColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy,
		pharmacy.Name, pharmacy.ManagerID, pharmacy.Address,
		postgis.NewPointFromCoordinate(pharmacy.Coordinate),
		pharmacy.PharmacistName, pharmacy.PharmacistLicense,
		pharmacy.PharmacistPhone, pharmacy.Slug,
	)
}

func (r *pharmacyRepository) Update(ctx context.Context, pharmacy domain.PharmacyUpdateDetails) (domain.Pharmacy, error) {
	q := `
		UPDATE pharmacies
		SET name = $1,
			address = $2,
			coordinate = $3,
			pharmacist_name = $4,
			pharmacist_license = $5,
			pharmacist_phone = $6
		WHERE slug = $7 RETURNING
	` + pharmacyColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy,
		pharmacy.Name, pharmacy.Address, postgis.NewPointFromCoordinate(pharmacy.Coordinate),
		pharmacy.PharmacistName, pharmacy.PharmacistLicense, pharmacy.PharmacistPhone,
		pharmacy.Slug,
	)
}

func (r *pharmacyRepository) SoftDeleteBySlug(ctx context.Context, slug string) error {
	q := `
		UPDATE pharmacies
		SET deleted_at = now(),
			updated_at = now()
		WHERE slug = $1
	`

	return exec(
		r.querier, ctx, q,
		slug,
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
		pharmacyOperation.EndTime, pharmacyOperation.PharmacyId, pharmacyOperation.ID,
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

func (r *pharmacyRepository) GetPharmacyOperationsByPharmacyIdAndLock(ctx context.Context, id int64) ([]domain.PharmacyOperations, error) {
	q := `
	SELECT ` + pharmacyOperationColumns + ` 
	FROM pharmacy_operations
	WHERE pharmacy_id = $1 AND deleted_at IS NULL
	FOR UPDATE
	`

	return queryFull(
		r.querier, ctx, q,
		ScanPharmacyOperation, id,
	)
}
