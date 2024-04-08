package postgres

import (
	"context"
	"database/sql"
	"medichat-be/apperror"
	"medichat-be/domain"
	"time"
)

type dataRepository struct {
	conn    *sql.DB
	querier Querier
}

func NewDataRepository(db *sql.DB) *dataRepository {
	return &dataRepository{
		conn:    db,
		querier: db,
	}
}

func (r *dataRepository) Atomic(
	ctx context.Context,
	fn domain.AtomicFunc[any],
) (any, error) {
	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, apperror.Wrap(err)
	}
	defer tx.Rollback()

	txRepo := &dataRepository{
		conn:    r.conn,
		querier: tx,
	}

	ret, err := fn(txRepo)
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return ret, nil
}

func (r *dataRepository) Sleep(ctx context.Context, duration time.Duration) error {
	q := `
		SELECT pg_sleep($1)
	`

	return exec(
		r.querier, ctx, q,
		duration.Seconds(),
	)
}
