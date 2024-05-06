package postgres

import (
	"context"
	"medichat-be/domain"
)

type pharmacyManagerRepository struct {
	querier Querier
}

func (r *pharmacyManagerRepository) GetByID(ctx context.Context, id int64) (domain.PharmacyManager, error) {
	q := `
		SELECT ` + pharmacyManagerJoinedColumns + `
		FROM pharmacy_managers p JOIN accounts a ON p.account_id = a.id
		WHERE p.id = $1
			AND p.deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacyManagerJoined,
		id,
	)
}

func (r *pharmacyManagerRepository) GetByIDAndLock(ctx context.Context, id int64) (domain.PharmacyManager, error) {
	q := `
		SELECT ` + pharmacyManagerJoinedColumns + `
		FROM pharmacy_managers p JOIN accounts a ON p.account_id = a.id
		WHERE p.id = $1
			AND p.deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacyManagerJoined,
		id,
	)
}

func (r *pharmacyManagerRepository) IsExistByID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM pharmacy_managers
			WHERE id = $1
				AND deleted_at IS NULL
		)
	`

	return queryOne(
		r.querier, ctx, q,
		boolScanDest,
		id,
	)
}

func (r *pharmacyManagerRepository) GetByAccountID(ctx context.Context, id int64) (domain.PharmacyManager, error) {
	q := `
		SELECT ` + pharmacyManagerJoinedColumns + `
		FROM pharmacy_managers p JOIN accounts a ON p.account_id = a.id
		WHERE p.account_id = $1
			AND p.deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacyManagerJoined,
		id,
	)
}

func (r *pharmacyManagerRepository) GetByAccountIDAndLock(ctx context.Context, id int64) (domain.PharmacyManager, error) {
	q := `
		SELECT ` + pharmacyManagerJoinedColumns + `
		FROM pharmacy_managers p JOIN accounts a ON p.account_id = a.id
		WHERE p.account_id = $1
			AND p.deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacyManagerJoined,
		id,
	)
}

func (r *pharmacyManagerRepository) IsExistByAccountID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM pharmacy_managers
			WHERE account_id = $1
				AND deleted_at IS NULL
		)
	`

	return queryOne(
		r.querier, ctx, q,
		boolScanDest,
		id,
	)
}

func (r *pharmacyManagerRepository) Add(ctx context.Context, ph domain.PharmacyManager) (domain.PharmacyManager, error) {
	q := `
		INSERT INTO pharmacy_managers(account_id)
		VALUES ($1)
		RETURNING
	` + pharmacyManagerColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacyManager,
		ph.Account.ID,
	)
}
