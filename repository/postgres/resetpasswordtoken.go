package postgres

import (
	"context"
	"medichat-be/domain"
)

type resetPasswordTokenRepository struct {
	querier Querier
}

func (r *resetPasswordTokenRepository) Add(
	ctx context.Context,
	token domain.ResetPasswordToken,
) (domain.ResetPasswordToken, error) {
	q := `
		INSERT INTO reset_password_tokens(account_id, token, expired_at)
		VALUES
		($1, $2, $3)
		RETURNING ` + resetPasswordTokenColumns

	return queryOne(
		r.querier, ctx, q,
		resetPasswordTokenScanDests,
		token.Account.ID, token.Token, token.ExpiredAt,
	)
}

func (r *resetPasswordTokenRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.ResetPasswordToken, error) {
	q := `
		SELECT ` + resetPasswordTokenColumns + `
		FROM reset_password_tokens
		WHERE id = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		resetPasswordTokenScanDests,
		id,
	)
}

func (r *resetPasswordTokenRepository) GetByTokenStr(
	ctx context.Context,
	tokenStr string,
) (domain.ResetPasswordToken, error) {
	q := `
		SELECT ` + resetPasswordTokenColumns + `
		FROM reset_password_tokens
		WHERE token = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		resetPasswordTokenScanDests,
		tokenStr,
	)
}

func (r *resetPasswordTokenRepository) GetByTokenStrAndLock(
	ctx context.Context,
	tokenStr string,
) (domain.ResetPasswordToken, error) {
	q := `
		SELECT ` + resetPasswordTokenColumns + `
		FROM reset_password_tokens
		WHERE token = $1
			AND expired_at > now() 
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOne(
		r.querier, ctx, q,
		resetPasswordTokenScanDests,
		tokenStr,
	)
}

func (r *resetPasswordTokenRepository) SoftDeleteByID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE reset_password_tokens
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return exec(
		r.querier, ctx, q,
		id,
	)
}

var (
	resetPasswordTokenColumns = " id, account_id, token, expired_at "
)

func resetPasswordTokenScanDests(t *domain.ResetPasswordToken) []any {
	return []any{
		&t.ID, &t.Account.ID, &t.Token, &t.ExpiredAt,
	}
}
