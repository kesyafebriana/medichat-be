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
	accountColumns                = " id, email, email_verified, role, account_type "
	accountWithCredentialsColumns = " id, email, email_verified, role, account_type, hashed_password "
)

func accountScanDests(u *domain.Account) []any {
	return []any{
		&u.ID, &u.Email, &u.EmailVerified, &u.Role, &u.AccountType,
	}
}

func scanAccountWithCredentials(r RowScanner, a *domain.AccountWithCredentials) error {
	var nullHashedPassword sql.NullString
	if err := r.Scan(
		&a.Account.ID, &a.Account.Email, &a.Account.EmailVerified,
		&a.Account.Role, &a.Account.AccountType, &nullHashedPassword,
	); err != nil {
		return err
	}
	a.HashedPassword = toStringPtr(nullHashedPassword)
	return nil
}

var (
	categoryColumns               = " id, parent_id, name "
	categoryWithParentNameColumns = " c.id, c.parent_id, c.name, c2.name as parent_name "
)

func scanCategory(r RowScanner, c *domain.Category) error {
	var nullParentId sql.NullInt64
	if err := r.Scan(
		&c.ID, &nullParentId, &c.Name,
	); err != nil {
		return err
	}
	c.ParentID = toInt64Ptr(nullParentId)
	return nil
}

func scanCategoryWithParentName(r RowScanner, c *domain.CategoryWithParentName) error {
	var nullParentId sql.NullInt64
	var nullParentName sql.NullString
	if err := r.Scan(
		&c.Category.ID, &nullParentId, &c.Category.Name, &nullParentName,
	); err != nil {
		return err
	}
	c.Category.ParentID = toInt64Ptr(nullParentId)
	c.ParentName = toStringPtr(nullParentName)
	return nil
}
