package postgres

import (
	"context"
	"medichat-be/domain"
)

type specializationRepository struct {
	querier Querier
}

func (r *specializationRepository) GetAll(
	ctx context.Context,
) ([]domain.Specialization, error) {
	q := `
		SELECT ` + specializationColumns + `
		FROM specializations
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`

	return queryFull(
		r.querier, ctx, q,
		scanSpecialization,
	)
}

func (r *specializationRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.Specialization, error) {
	q := `
		SELECT ` + specializationColumns + `
		FROM specializations
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanSpecialization,
		id,
	)
}
