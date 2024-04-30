package postgres

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
	"strings"
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
		INSERT INTO users(account_id, date_of_birth, main_location_id)
		VALUES
		($1, $2, $3)
		RETURNING ` + userColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanUser,
		u.Account.ID, u.DateOfBirth, u.MainLocationID,
	)
}

func (r *userRepository) Update(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {
	q := `
		UPDATE users
		SET date_of_birth = $2,
			main_location_id = $3,
			updated_at = now()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	err := execOne(
		r.querier, ctx, q,
		u.ID, u.DateOfBirth, u.MainLocationID,
	)
	if err != nil {
		return domain.User{}, apperror.Wrap(err)
	}

	return u, nil
}

func (r *userRepository) GetLocationsByUserID(
	ctx context.Context,
	id int64,
) ([]domain.UserLocation, error) {
	q := `
		SELECT ` + userLocationColumns + `
		FROM user_locations
		WHERE user_id = $1
			AND deleted_at IS NULL
	`

	return queryFull(
		r.querier, ctx, q,
		scanUserLocation,
		id,
	)
}

func (r *userRepository) GetLocationByID(
	ctx context.Context,
	id int64,
) (domain.UserLocation, error) {
	q := `
		SELECT ` + userLocationColumns + `
		FROM user_locations
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanUserLocation,
		id,
	)
}

func (r *userRepository) GetLocationByIDAndLock(
	ctx context.Context,
	id int64,
) (domain.UserLocation, error) {
	q := `
		SELECT ` + userLocationColumns + `
		FROM user_locations
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanUserLocation,
		id,
	)
}

func (r *userRepository) IsAnyLocationActiveByUserID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM user_locations
			WHERE user_id = $1 
				AND is_active = true
				AND deleted_at IS NULL
		)
	`

	return queryOne(
		r.querier, ctx, q,
		boolScanDest,
		id,
	)
}

func (r *userRepository) AddLocation(
	ctx context.Context,
	ul domain.UserLocation,
) (domain.UserLocation, error) {
	q := `
		INSERT INTO user_locations(user_id, alias, address, coordinate, is_active)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING ` + userLocationColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanUserLocation,
		ul.UserID, ul.Alias, ul.Address,
		postgis.NewPointFromCoordinate(ul.Coordinate),
		ul.IsActive,
	)
}

func (r *userRepository) AddLocations(
	ctx context.Context,
	uls []domain.UserLocation,
) ([]domain.UserLocation, error) {
	var sb strings.Builder
	l := len(uls)
	args := make([]any, l*5)

	sb.WriteString(`
		INSERT INTO user_locations(user_id, alias, address, coordinate, is_active)
		VALUES
	`)

	for i, ul := range uls {
		if i > 0 {
			sb.WriteString(`,
			`)
		}
		j := i * 5

		fmt.Fprintf(
			&sb, " ($%d, $%d, $%d, $%d, $%d) ",
			j+1, j+2, j+3, j+4, j+5,
		)

		args[j] = ul.UserID
		args[j+1] = ul.Alias
		args[j+2] = ul.Address
		args[j+3] = postgis.NewPointFromCoordinate(ul.Coordinate)
		args[j+4] = ul.IsActive
	}

	sb.WriteString(` RETURNING ` + userLocationColumns)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanUserLocation,
		args...,
	)
}

func (r *userRepository) UpdateLocation(
	ctx context.Context,
	ul domain.UserLocation,
) (domain.UserLocation, error) {
	q := `
		UPDATE user_locations
		SET alias = $2,
			address = $3,
			coordinate = $4,
			is_active = $5
		WHERE ID = $1
			AND deleted_at IS NULL
	`

	err := execOne(
		r.querier, ctx, q,
		ul.ID, ul.Alias, ul.Address,
		postgis.NewPointFromCoordinate(ul.Coordinate),
		ul.IsActive,
	)
	if err != nil {
		return domain.UserLocation{}, apperror.Wrap(err)
	}

	return ul, nil
}

func (r *userRepository) SoftDeleteLocationByID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE user_locations
		SET deleted_at = now(),
			updated_at = now()
		WHERE ID = $1
			AND deleted_at IS NULL
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}
