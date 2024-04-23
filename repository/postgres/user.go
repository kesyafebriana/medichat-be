package postgres

import (
	"context"
	"medichat-be/domain"
)

type userRepository struct {
	querier Querier
}

func (r *userRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.User, error) {
	q := `
		SELECT ` + userJoinedColumns + `
		FROM users u JOIN accounts a ON u.account_id = a.id
		WHERE u.id = $1
			AND u.deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanUserJoined,
		id,
	)
}

func (r *userRepository) GetByIDAndLock(
	ctx context.Context,
	id int64,
) (domain.User, error) {
	q := `
		SELECT ` + userJoinedColumns + `
		FROM users u JOIN accounts a ON u.account_id = a.id
		WHERE u.id = $1
			AND u.deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanUserJoined,
		id,
	)
}

func (r *userRepository) IsExistByID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM users
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

func (r *userRepository) GetByAccountID(
	ctx context.Context,
	id int64,
) (domain.User, error) {
	q := `
		SELECT ` + userJoinedColumns + `
		FROM users u JOIN accounts a ON u.account_id = a.id
		WHERE u.account_id = $1
			AND u.deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanUserJoined,
		id,
	)
}

func (r *userRepository) GetByAccountIDAndLock(
	ctx context.Context,
	id int64,
) (domain.User, error) {
	q := `
		SELECT ` + userJoinedColumns + `
		FROM users u JOIN accounts a ON u.account_id = a.id
		WHERE u.account_id = $1
			AND u.deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanUserJoined,
		id,
	)
}

func (r *userRepository) IsExistByAccountID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM users
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

func (r *userRepository) Add(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {
	q := `
		INSERT INTO users(account_id, date_of_birth)
		VALUES
		($1, $2)
		RETURNING ` + userColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanUser,
		u.Account.ID, u.DateOfBirth,
	)
}

func (r *userRepository) Update(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {
	q := `
		UPDATE users
		SET date_of_birth = $2,
			updated_at = now()
		WHERE id = $1
		RETURNING ` + userColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanUser,
		u.ID, u.DateOfBirth,
	)
}
