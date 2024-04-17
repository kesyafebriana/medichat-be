package postgres

import (
	"context"
	"medichat-be/domain"
)

type verifyEmailTokenRepository struct {
	querier Querier
}

func (r *verifyEmailTokenRepository) Add(
	ctx context.Context,
	token domain.VerifyEmailToken,
) (domain.VerifyEmailToken, error) {
	q := `
		INSERT INTO verify_email_tokens(account_id, token, expired_at)
		VALUES
		($1, $2, $3)
		RETURNING ` + verifyEmailTokenColumns

	return queryOne(
		r.querier, ctx, q,
		verifyEmailTokenScanDests,
		token.Account.ID, token.Token, token.ExpiredAt,
	)
}

func (r *verifyEmailTokenRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.VerifyEmailToken, error) {
	q := `
		SELECT ` + verifyEmailTokenColumns + `
		FROM verify_email_tokens
		WHERE id = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		verifyEmailTokenScanDests,
		id,
	)
}

func (r *verifyEmailTokenRepository) GetByTokenStr(
	ctx context.Context,
	tokenStr string,
) (domain.VerifyEmailToken, error) {
	q := `
		SELECT ` + verifyEmailTokenColumns + `
		FROM verify_email_tokens
		WHERE token = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		verifyEmailTokenScanDests,
		tokenStr,
	)
}

func (r *verifyEmailTokenRepository) GetByTokenStrAndLock(
	ctx context.Context,
	tokenStr string,
) (domain.VerifyEmailToken, error) {
	q := `
		SELECT ` + verifyEmailTokenColumns + `
		FROM verify_email_tokens
		WHERE token = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOne(
		r.querier, ctx, q,
		verifyEmailTokenScanDests,
		tokenStr,
	)
}

func (r *verifyEmailTokenRepository) SoftDeleteByID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE verify_email_tokens
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return exec(
		r.querier, ctx, q,
		id,
	)
}

func (r *verifyEmailTokenRepository) SoftDeleteByAccountID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE verify_email_tokens
		SET deleted_at = now(),
			updated_at = now()
		WHERE account_id = $1
	`

	return exec(
		r.querier, ctx, q,
		id,
	)
}

var (
	verifyEmailTokenColumns = " id, account_id, token, expired_at "
)

func verifyEmailTokenScanDests(t *domain.VerifyEmailToken) []any {
	return []any{
		&t.ID, &t.Account.ID, &t.Token, &t.ExpiredAt,
	}
}
