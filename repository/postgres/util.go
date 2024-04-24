package postgres

import (
	"database/sql"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
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

var (
	userColumns = `
		id, account_id, date_of_birth
	`

	userJoinedColumns = `
		u.id,
		u.account_id, a.email, a.email_verified, a.role, a.account_type,
		a.name, a.photo_url, u.date_of_birth
	`

	userLocationColumns = `
		id, user_id, alias, address, coordinate, is_active
	`
)

func scanUser(r RowScanner, u *domain.User) error {
	a := &u.Account
	return r.Scan(
		&u.ID, &a.ID, &u.DateOfBirth,
	)
}

func scanUserJoined(r RowScanner, u *domain.User) error {
	a := &u.Account
	return r.Scan(
		&u.ID,
		&a.ID, &a.Email, &a.EmailVerified, &a.Role, &a.AccountType,
		&a.Name, &a.PhotoURL, &u.DateOfBirth,
	)
}

func scanUserLocation(r RowScanner, ul *domain.UserLocation) error {
	var p postgis.Point
	if err := r.Scan(
		&ul.ID, &ul.UserID, &ul.Alias, &ul.Address, &p, &ul.IsActive,
	); err != nil {
		return err
	}
	ul.Coordinate = p.ToCoordinate()
	return nil
}
