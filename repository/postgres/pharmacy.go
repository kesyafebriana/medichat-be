package postgres

import (
	"context"
	"fmt"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
	"strings"
	"time"
)

type pharmacyRepository struct {
	querier Querier
}

func (r *pharmacyRepository) GetPharmacies(ctx context.Context, query domain.PharmaciesQuery) ([]domain.Pharmacy, error) {
	sb := strings.Builder{}
	var args = make([]any, 0)
	var idx = 1
	offset := (query.Page - 1) * query.Limit

	sb.WriteString(`
		SELECT ` + pharmacyJoinedColumns + `
		FROM pharmacies p JOIN stocks as s ON s.pharmacy_id = p.id
		WHERE p.deleted_at IS NULL
	`)

	if query.ProductId != nil {
		fmt.Fprintf(&sb, ` AND s.product_id = $%d 
		`, idx)
		idx++
		args = append(args, *query.ProductId)
	}

	if query.Name != nil {
		fmt.Fprintf(&sb, ` AND p.name ILIKE $%d 
		`, idx)
		idx++
		args = append(args, *query.Name)
	}

	if query.ManagerID != nil {
		fmt.Fprintf(&sb, ` AND p.manager_id = $%d 
		`, idx)
		idx++
		args = append(args, *query.ManagerID)
	}

	if query.Longitude != nil && query.Latitude != nil {
		fmt.Fprintf(&sb, ` AND ST_DWithin(p.coordinate, ST_MakePoint($%d, $%d)::geography, 25000)
		`, idx, idx+1)
		idx += 2
		args = append(args, *query.Longitude, *query.Latitude)
	}

	if query.IsOpen != nil && *query.IsOpen {
		time := time.Now()

		fmt.Fprintf(&sb, `
		AND p.id IN (
		SELECT o.pharmacy_id
		FROM pharmacy_operations o
		WHERE o.deleted_at IS NULL
		AND o.day = '%s'
		AND o.start_time <= '0001-01-01 %s:00')
		`, time.Format("Monday"), time.Format("15:04"))
	}

	if (query.Day != nil || query.StartTime != nil || query.EndTime != nil) && query.IsOpen == nil {
		sb.WriteString(`
			AND p.id IN (
			SELECT o.pharmacy_id
			FROM pharmacy_operations o
			WHERE o.deleted_at IS NULL
		`)

		if query.Day != nil {
			fmt.Fprintf(&sb, `AND o.day = $%d
		`, idx)
			idx++
			args = append(args, *query.Day)
		}

		if query.StartTime != nil {
			fmt.Fprintf(&sb, `AND o.start_time <= '0001-01-01 %s:00'
		`, *query.StartTime)
		}

		if query.EndTime != nil {
			fmt.Fprintf(&sb, `AND o.end_time >= '0001-01-01 %s:00'
		`, *query.EndTime)
		}

		sb.WriteString(`)`)
	}


	if query.SortBy == domain.PharmacySortByName {
		fmt.Fprintf(&sb, " ORDER BY %s %s", query.SortBy, query.SortType)
	}

	if query.SortBy == domain.PharmacySortByDistance && query.Latitude!=nil && query.Longitude!=nil{
		fmt.Fprintf(&sb, " ORDER BY p.coordinate <-> ST_MakePoint(%f, %f)::geometry",*query.Longitude, *query.Latitude)
	}

	if query.Limit != 0 {
		fmt.Fprintf(&sb, " OFFSET %d LIMIT %d ", offset, query.Limit)
	}


	return queryFull(
		r.querier, ctx, sb.String(),
		scanPharmacy,
		args...,
	)
}

func (r *pharmacyRepository) GetPageInfo(ctx context.Context, query domain.PharmaciesQuery) (domain.PageInfo, error) {
	sb := strings.Builder{}
	var args = make([]any, 0)
	var idx = 1
	offset := (query.Page - 1) * query.Limit

	sb.WriteString(`
		SELECT COUNT(p.*) as total_data 
		FROM pharmacies p
		WHERE p.deleted_at IS NULL
	`)

	if query.Name != nil {
		fmt.Fprintf(&sb, ` AND p.name ILIKE $%d 
		`, idx)
		idx++
		args = append(args, *query.Name)
	}

	if query.ManagerID != nil {
		fmt.Fprintf(&sb, ` AND p.manager_id = $%d 
		`, idx)
		idx++
		args = append(args, *query.ManagerID)
	}

	if query.Longitude != nil && query.Latitude != nil {
		fmt.Fprintf(&sb, ` AND ST_DWithin(p.coordinate, ST_MakePoint($%d, $%d)::geography, 25000)
		`, idx, idx+1)
		idx += 2
		args = append(args, *query.Longitude, *query.Latitude)
	}

	if query.IsOpen != nil && *query.IsOpen {
		time := time.Now()

		fmt.Fprintf(&sb, `
		AND p.id IN (
		SELECT o.pharmacy_id
		FROM pharmacy_operations o
		WHERE o.deleted_at IS NULL
		AND o.day = '%s'
		AND o.start_time <= '0001-01-01 %s:00:00')
		`, time.Format("Monday"), time.Format("15"))
	}

	if (query.Day != nil || query.StartTime != nil || query.EndTime != nil) && query.IsOpen == nil {
		sb.WriteString(`
			AND p.id IN (
			SELECT o.pharmacy_id
			FROM pharmacy_operations o
			WHERE o.deleted_at IS NULL
		`)

		if query.Day != nil {
			fmt.Fprintf(&sb, `AND o.day = $%d
		`, idx)
			idx++
			args = append(args, *query.Day)
		}

		if query.StartTime != nil {
			fmt.Fprintf(&sb, `AND o.start_time <= '0001-01-01 %s:00'
		`, *query.StartTime)
		}

		if query.EndTime != nil {
			fmt.Fprintf(&sb, `AND o.end_time >= '0001-01-01 %s:00'
		`, *query.EndTime)
		}

		sb.WriteString(`)`)
	}

	if query.Limit != 0 {
		fmt.Fprintf(&sb, " OFFSET %d LIMIT %d ", offset, query.Limit)
	}

	var totalData int64
	row := r.querier.QueryRowContext(ctx, sb.String(), args...)
	err := row.Scan(&totalData)

	if err != nil {
		return domain.PageInfo{}, err
	}

	return domain.PageInfo{
		CurrentPage: int(query.Page),
		ItemCount:   totalData,
	}, nil
}

func (r *pharmacyRepository) GetBySlug(ctx context.Context, slug string) (domain.Pharmacy, error) {
	q := `
		SELECT ` + pharmacyColumns + `FROM pharmacies
		WHERE slug = $1 AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy, slug,
	)
}

func (r *pharmacyRepository) GetByID(ctx context.Context, id int64) (domain.Pharmacy, error) {
	q := `
		SELECT ` + pharmacyColumns + `FROM pharmacies
		WHERE id = $1 AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy, id,
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
			updated_at = now()
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
			pharmacy_id = $4,
			updated_at = now()
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

func (r *pharmacyRepository) GetShipmentMethodsByPharmacyId(ctx context.Context, id int64) ([]domain.PharmacyShipmentMethods, error) {
	q := `
		SELECT ` + PharmacyShipmentMethodColumns + `
		FROM pharmacy_shipment_methods
		WHERE deleted_at IS NULL AND pharmacy_id = $1
	`

	return queryFull(
		r.querier, ctx, q,
		ScanPharmacyShipmentMethod,
		id,
	)
}

func (r *pharmacyRepository) GetShipmentMethodsByPharmacyIdAndLock(ctx context.Context, id int64) ([]domain.PharmacyShipmentMethods, error) {
	q := `
		SELECT ` + PharmacyShipmentMethodColumns + `
		FROM pharmacy_shipment_methods
		WHERE deleted_at IS NULL AND pharmacy_id = $1
		FOR UPDATE
	`

	return queryFull(
		r.querier, ctx, q,
		ScanPharmacyShipmentMethod,
		id,
	)
}

func (r *pharmacyRepository) AddShipmentMethod(ctx context.Context, pharmacyCourier domain.PharmacyShipmentMethodsCreateDetails) (domain.PharmacyShipmentMethods, error) {
	q := `
		INSERT INTO pharmacy_shipment_methods(pharmacy_id, shipment_method_id)
		VALUES($1, $2)
		RETURNING
	` + PharmacyShipmentMethodColumns

	return queryOneFull(
		r.querier, ctx, q,
		ScanPharmacyShipmentMethod,
		pharmacyCourier.PharmacyID,
		pharmacyCourier.ShipmentMethodID,
	)
}

func (r *pharmacyRepository) SoftDeleteShipmentMethodByID(ctx context.Context, id int64) error {
	q := `
		UPDATE pharmacy_shipment_methods
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}
