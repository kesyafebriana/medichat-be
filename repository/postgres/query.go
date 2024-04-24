package postgres

import (
	"context"
	"database/sql"
	"medichat-be/apperror"
)

func query[T any](
	querier Querier,
	ctx context.Context,
	query string,
	scanDestsFunc func(*T) []any,
	args ...any,
) ([]T, error) {
	return queryFull(
		querier,
		ctx,
		query,
		func(rs RowScanner, t *T) error {
			scanDests := scanDestsFunc(t)
			return rs.Scan(scanDests...)
		},
		args...,
	)
}

func queryFull[T any](
	querier Querier,
	ctx context.Context,
	query string,
	scanFunc func(RowScanner, *T) error,
	args ...any,
) ([]T, error) {
	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Wrap(err)
	}
	defer rows.Close()

	var ret []T

	for rows.Next() {
		var t T

		err = scanFunc(rows, &t)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		ret = append(ret, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return ret, nil
}

func queryOne[T any](
	querier Querier,
	ctx context.Context,
	query string,
	scanDestsFunc func(*T) []any,
	args ...any,
) (T, error) {
	return queryOneFull(
		querier,
		ctx,
		query,
		func(rs RowScanner, t *T) error {
			scanDests := scanDestsFunc(t)
			return rs.Scan(scanDests...)
		},
		args...,
	)
}

func queryOneFull[T any](
	querier Querier,
	ctx context.Context,
	query string,
	scanFunc func(RowScanner, *T) error,
	args ...any,
) (T, error) {
	var ret T
	var empty T

	err := scanFunc(querier.QueryRowContext(ctx, query, args...), &ret)
	if err == sql.ErrNoRows {
		return empty, apperror.NewNotFound()
	}
	if err != nil {
		return empty, apperror.Wrap(err)
	}

	return ret, nil
}

func exec(
	querier Querier,
	ctx context.Context,
	query string,
	args ...any,
) error {
	_, err := querier.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.Wrap(err)
	}

	return nil
}

func execOne(
	querier Querier,
	ctx context.Context,
	query string,
	args ...any,
) error {
	res, err := querier.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.Wrap(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap(err)
	}
	if rows != 1 {
		return apperror.NewInternalFmt("query: rows affected is not 1 (got %d)", rows)
	}

	return nil
}
