package postgres

import (
	"database/sql"
	"medichat-be/domain"
)

func int64ScanDest(i *int64) []any {
	return []any{i}
}

func boolScanDest(b *bool) []any {
	return []any{b}
}

func stringScanDest(s *string) []any {
	return []any{s}
}

var (
	accountColumns                = " id, email, email_verified, name, photo_url, role, account_type, profile_set "
	accountWithCredentialsColumns = " id, email, email_verified, name, photo_url, role, account_type, profile_set, hashed_password "
)

func accountScanDests(u *domain.Account) []any {
	return []any{
		&u.ID, &u.Email, &u.EmailVerified, &u.Name, &u.PhotoURL, &u.Role, &u.AccountType, &u.ProfileSet,
	}
}

func scanAccountWithCredentials(r RowScanner, a *domain.AccountWithCredentials) error {
	var nullHashedPassword sql.NullString
	if err := r.Scan(
		&a.Account.ID, &a.Account.Email, &a.Account.EmailVerified,
		&a.Account.Name, &a.Account.PhotoURL,
		&a.Account.Role, &a.Account.AccountType, &a.Account.ProfileSet, &nullHashedPassword,
	); err != nil {
		return err
	}
	a.HashedPassword = toStringPtr(nullHashedPassword)
	return nil
}
