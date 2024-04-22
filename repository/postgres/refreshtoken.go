package postgres

import (
	"context"
	"medichat-be/domain"
)

type refreshTokenRepository struct {
	querier Querier
}

func (r *refreshTokenRepository) Add(
	ctx context.Context,
	token domain.RefreshToken,
) (domain.RefreshToken, error) {
	q := `
		INSERT INTO refresh_tokens(account_id, token, client_ip, expired_at)
		VALUES
		($1, $2, $3, $4)
		RETURNING ` + refreshTokenColumns

	return queryOne(
		r.querier, ctx, q,
		refreshTokenScanDests,
		token.Account.ID, token.Token, token.ClientIP, token.ExpiredAt,
	)
}

func (r *refreshTokenRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.RefreshToken, error) {
	q := `
		SELECT ` + refreshTokenColumns + `
		FROM refresh_tokens
		WHERE id = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		refreshTokenScanDests,
		id,
	)
}

func (r *refreshTokenRepository) GetByTokenStr(
	ctx context.Context,
	tokenStr string,
) (domain.RefreshToken, error) {
	q := `
		SELECT ` + refreshTokenColumns + `
		FROM refresh_tokens
		WHERE token = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		refreshTokenScanDests,
		tokenStr,
	)
}

func (r *refreshTokenRepository) GetByTokenStrAndLock(
	ctx context.Context,
	tokenStr string,
) (domain.RefreshToken, error) {
	q := `
		SELECT ` + refreshTokenColumns + `
		FROM refresh_tokens
		WHERE token = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOne(
		r.querier, ctx, q,
		refreshTokenScanDests,
		tokenStr,
	)
}

func (r *refreshTokenRepository) SoftDeleteByID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE refresh_tokens
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return exec(
		r.querier, ctx, q,
		id,
	)
}

func (r *refreshTokenRepository) SoftDeleteByAccountID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE refresh_tokens
		SET deleted_at = now(),
			updated_at = now()
		WHERE account_id = $1
	`

	return exec(
		r.querier, ctx, q,
		id,
	)
}

func (r *refreshTokenRepository) SoftDeleteByClientIP(
	ctx context.Context,
	ip string,
) error {
	q := `
		UPDATE refresh_tokens
		SET deleted_at = now(),
			updated_at = now()
		WHERE client_ip = $1
	`

	return exec(
		r.querier, ctx, q,
		ip,
	)
}

var (
	refreshTokenColumns = " id, account_id, token, client_ip, expired_at "
)

func refreshTokenScanDests(t *domain.RefreshToken) []any {
	return []any{
		&t.ID, &t.Account.ID, &t.Token, &t.ClientIP, &t.ExpiredAt,
	}
}
