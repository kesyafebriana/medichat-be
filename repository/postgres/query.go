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
	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Wrap(err)
	}
	defer rows.Close()

	var ret []T

	for rows.Next() {
		var t T
		scanDests := scanDestsFunc(&t)

		err = rows.Scan(scanDests...)
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
	var ret T
	var empty T
	scanDests := scanDestsFunc(&ret)

	err := querier.QueryRowContext(ctx, query, args...).Scan(scanDests...)
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
	res, err := querier.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.Wrap(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap(err)
	}
	if n == 0 {
		return apperror.NewNotFound()
	}

	return nil
}
