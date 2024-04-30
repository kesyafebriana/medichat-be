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
	categoryColumns               = " id, parent_id, name, slug, photo_url "
	categoryWithParentNameColumns = " c.id, c.parent_id, c.name, c2.name as parent_name, c.slug, c.photo_url "
)

func scanCategory(r RowScanner, c *domain.Category) error {
	var nullParentId sql.NullInt64
	var nullPhotoUrl sql.NullString
	if err := r.Scan(
		&c.ID, &nullParentId, &c.Name, &c.Slug, &nullPhotoUrl,
	); err != nil {
		return err
	}
	c.ParentID = toInt64Ptr(nullParentId)
	c.PhotoUrl = toStringPtr(nullPhotoUrl)
	return nil
}

func scanCategoryWithParentName(r RowScanner, c *domain.CategoryWithParentName) error {
	var nullParentId sql.NullInt64
	var nullParentName sql.NullString
	var nullPhotoUrl sql.NullString
	if err := r.Scan(
		&c.Category.ID, &nullParentId, &c.Category.Name, &nullParentName, &c.Category.Slug, &nullPhotoUrl,
	); err != nil {
		return err
	}
	c.Category.ParentID = toInt64Ptr(nullParentId)
	c.ParentName = toStringPtr(nullParentName)
	c.Category.PhotoUrl = toStringPtr(nullPhotoUrl)
	return nil
}

var (
	productColumns               = " id, name, product_detail_id, category_id, picture, is_active  "
	productDetailsColumns        = " id, generic_name, content, manufacturer, description, product_classification, product_form, unit_in_pack, selling_unit, weight, height, length, width  "
)

func scanProduct(r RowScanner, c *domain.Product) error {
	var nullPhotoUrl sql.NullString
	if err := r.Scan(
		&c.ID, &c.Name, &c.ProductDetailId, &c.ProductCategoryId, &nullPhotoUrl,&c.IsActive,
	); err != nil {
		return err
	}
	c.Picture = toStringPtr(nullPhotoUrl)
	return nil
}

func scanProductDetails(r RowScanner, d *domain.ProductDetails) error {
	if err := r.Scan(
		&d.ID, &d.GenericName, &d.Content, &d.Manufacturer, &d.Description,&d.ProductClassification,&d.ProductForm,&d.UnitInPack,&d.SellingUnit,&d.Weight,&d.Height,&d.Length,&d.Width,
	); err != nil {
		return err
	}

	return nil
}
