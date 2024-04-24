package postgres

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
)

type accountRepository struct {
	querier Querier
}

func (r *accountRepository) GetByEmail(
	ctx context.Context,
	email string,
) (domain.Account, error) {
	q := `
		SELECT ` + accountColumns + `
		FROM accounts
		WHERE email = $1
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		accountScanDests,
		email,
	)
}

func (r *accountRepository) GetByEmailAndLock(
	ctx context.Context,
	email string,
) (domain.Account, error) {
	q := `
		SELECT ` + accountColumns + `
		FROM accounts
		WHERE email = $1
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOne(
		r.querier, ctx, q,
		accountScanDests,
		email,
	)
}

func (r *accountRepository) GetWithCredentialsByEmail(
	ctx context.Context,
	email string,
) (domain.AccountWithCredentials, error) {
	q := `
		SELECT ` + accountWithCredentialsColumns + `
		FROM accounts
		WHERE email = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanAccountWithCredentials,
		email,
	)
}

func (r *accountRepository) IsExistByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	q := `
		SELECT EXISTS(
			SELECT id
			FROM accounts
			WHERE email = $1
				AND deleted_at IS NULL
		)
	`

	return queryOne(
		r.querier, ctx, q,
		boolScanDest,
		email,
	)
}

func (r *accountRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.Account, error) {
	q := `
		SELECT ` + accountColumns + `
		FROM accounts
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOne(
		r.querier, ctx, q,
		accountScanDests,
		id,
	)
}

func (r *accountRepository) GetByIDAndLock(
	ctx context.Context,
	id int64,
) (domain.Account, error) {
	q := `
		SELECT ` + accountColumns + `
		FROM accounts
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOne(
		r.querier, ctx, q,
		accountScanDests,
		id,
	)
}

func (r *accountRepository) GetWithCredentialsByID(
	ctx context.Context,
	id int64,
) (domain.AccountWithCredentials, error) {
	q := `
		SELECT ` + accountWithCredentialsColumns + `
		FROM accounts
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanAccountWithCredentials,
		id,
	)
}

func (r *accountRepository) IsExistByID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS(
			SELECT id
			FROM accounts
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

func (r *accountRepository) Add(
	ctx context.Context,
	creds domain.AccountWithCredentials,
) (domain.Account, error) {
	q := `
		INSERT INTO accounts(email, email_verified, role, account_type, hashed_password, name, photo_url, profile_set)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING ` + accountColumns

	return queryOne(
		r.querier, ctx, q,
		accountScanDests,
		creds.Account.Email, creds.Account.EmailVerified,
		creds.Account.Role, creds.Account.AccountType,
		fromStringPtr(creds.HashedPassword),
		creds.Account.Name, creds.Account.PhotoURL, creds.Account.ProfileSet,
	)
}

func (r *accountRepository) Update(
	ctx context.Context,
	a domain.Account,
) (domain.Account, error) {
	q := `
		UPDATE accounts
		SET name = $2,
			photo_url = $3
		WHERE id = $1
	`

	err := execOne(
		r.querier, ctx, q,
		a.ID, a.Name, a.PhotoURL,
	)
	if err != nil {
		return domain.Account{}, apperror.Wrap(err)
	}

	return a, nil
}

func (r *accountRepository) UpdatePasswordByID(
	ctx context.Context,
	id int64,
	newHashedPassword string,
) error {
	q := `
		UPDATE accounts
		SET hashed_password = $1,
			updated_at = now()
		WHERE id = $2
	`

	return execOne(
		r.querier, ctx, q,
		newHashedPassword, id,
	)
}

func (r *accountRepository) VerifyEmailByID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE accounts
		SET email_verified = true,
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}

func (r *accountRepository) ProfileSetByID(
	ctx context.Context,
	id int64,
) error {
	q := `
		UPDATE accounts
		SET profile_set = true,
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}
