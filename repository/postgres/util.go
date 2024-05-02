package postgres

import (
	"database/sql"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
)

func getSortOrder(asc bool) string {
	if asc {
		return "ASC"
	}
	return "DESC"
}

func getSortCursorCmp(asc bool) string {
	if asc {
		return ">"
	}
	return "<"
}

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
	userColumns = `
		id, account_id, date_of_birth, main_location_id
	`

	userJoinedColumns = `
		u.id,
		u.account_id, a.email, a.email_verified, a.role, a.account_type,
		a.name, a.photo_url, a.profile_set, u.date_of_birth, u.main_location_id
	`

	userLocationColumns = `
		id, user_id, alias, address, coordinate, is_active
	`
)

func scanUser(r RowScanner, u *domain.User) error {
	a := &u.Account
	return r.Scan(
		&u.ID, &a.ID, &u.DateOfBirth, &u.MainLocationID,
	)
}

func scanUserJoined(r RowScanner, u *domain.User) error {
	a := &u.Account
	return r.Scan(
		&u.ID,
		&a.ID, &a.Email, &a.EmailVerified, &a.Role, &a.AccountType,
		&a.Name, &a.PhotoURL, &a.ProfileSet, &u.DateOfBirth, &u.MainLocationID,
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

var (
	doctorColumns = `
		id, account_id, specialization_id, str, work_location, gender,
		phone_number, is_active, start_work_date, price, certificate_url,
		now()::date - start_work_date as year_experience
	`

	doctorJoinedColumns = `
		d.id, 
		d.account_id, a.email, a.email_verified, a.role, a.account_type, 
		a.name, a.photo_url, a.profile_set,
		d.specialization_id, s.name, 
		d.str, d.work_location, d.gender, d.phone_number, d.is_active, 
		d.start_work_date, d.price, d.certificate_url,
		(now()::date - d.start_work_date) / 365 as year_experience
	`
)

func scanDoctor(r RowScanner, d *domain.Doctor) error {
	a := &d.Account
	s := &d.Specialization
	return r.Scan(
		&d.ID, &a.ID, &s.ID, &d.STR, &d.WorkLocation, &d.Gender,
		&d.PhoneNumber, &d.IsActive, &d.StartWorkDate, &d.Price,
		&d.CertificateURL, &d.YearExperience,
	)
}

func scanDoctorJoined(r RowScanner, d *domain.Doctor) error {
	a := &d.Account
	s := &d.Specialization
	return r.Scan(
		&d.ID,
		&a.ID, &a.Email, &a.EmailVerified, &a.Role, &a.AccountType,
		&a.Name, &a.PhotoURL, &a.ProfileSet,
		&s.ID, &s.Name,
		&d.STR, &d.WorkLocation, &d.Gender,
		&d.PhoneNumber, &d.IsActive, &d.StartWorkDate, &d.Price,
		&d.CertificateURL, &d.YearExperience,
	)
}

var (
	specializationColumns = `
		id, name
	`
)

func scanSpecialization(r RowScanner, s *domain.Specialization) error {
	return r.Scan(
		&s.ID, &s.Name,
	)
}

var (
	paymentColumns = `
		id, invoice_number, file_url, is_confirmed, amount
	`
)

func scanPayment(r RowScanner, p *domain.Payment) error {
	nullURL := sql.NullString{}
	if err := r.Scan(
		&p.ID, &p.InvoiceNumber, &nullURL, &p.IsConfirmed, &p.Amount,
	); err != nil {
		return err
	}
	p.FileURL = toStringPtr(nullURL)
	return nil
}

var (
	orderColumns = `
		id, user_id, pharmacy_id, payment_id, shipment_method_id,
		address, coordinate, 
		n_items, subtotal, shipment_fee, total,
		status, ordered_at, finished_at
	`

	selectOrderJoined = `
		SELECT
			o.id, 
			u.id, u.name
			ph.id, ph.slug, ph.name,
			py.id, py.invoice_number,
			sm.id, sm.name,
			o.address, o.coordinate, 
			o.n_items, o.subtotal, o.shipment_fee, o.total, 
			o.status, o.ordered_at, o.finished_at
		FROM orders o
			JOIN users u ON o.user_id = u.id
			JOIN pharmacies ph ON o.pharmacy_id = ph.id
			JOIN payments py ON o.payment_id = py.id
			JOIN shipment_methods sm ON o.shimpent_method_id = sm.id
	`

	countOrderJoined = `
		SELECT COUNT(o.id)
		FROM orders o
			JOIN users u ON o.user_id = u.id
			JOIN pharmacies ph ON o.pharmacy_id = ph.id
			JOIN payments py ON o.payment_id = py.id
			JOIN shipment_methods sm ON o.shimpent_method_id = sm.id
	`

	orderItemColumns = `
		id, order_id, product_id, price, amount
	`

	selectOrderItemJoined = `
		SELECT
			oi.id, oi.order_id,
			pd.id, pd.slug, pd.name,
			oi.price, oi.amount
		FROM order_items oi
			JOIN products pd ON oi.product_id = pd.id
	`
)

func scanOrder(r RowScanner, o *domain.Order) error {
	u := &o.User
	ph := &o.Pharmacy
	py := &o.Payment
	sm := &o.ShipmentMethod
	nullFinished := sql.NullTime{}
	point := postgis.Point{}
	if err := r.Scan(
		&o.ID, &u.ID, &ph.ID, &py.ID, &sm.ID,
		&o.Address, &point,
		&o.NItems, &o.Subtotal, &o.ShipmentFee, &o.Total,
		&o.Status, &o.OrderedAt, &nullFinished,
	); err != nil {
		return apperror.Wrap(err)
	}
	o.Coordinate = point.ToCoordinate()
	o.FinishedAt = toTimePtr(nullFinished)
	return nil
}

func scanOrderJoined(r RowScanner, o *domain.Order) error {
	u := &o.User
	ph := &o.Pharmacy
	py := &o.Payment
	sm := &o.ShipmentMethod
	nullFinished := sql.NullTime{}
	point := postgis.Point{}
	if err := r.Scan(
		&o.ID,
		&u.ID, &u.Name,
		&ph.ID, &ph.Slug, &ph.Name,
		&py.ID, &py.InvoiceNumber,
		&sm.ID, &sm.Name,
		&o.Address, &point,
		&o.NItems, &o.Subtotal, &o.ShipmentFee, &o.Total,
		&o.Status, &o.OrderedAt, &nullFinished,
	); err != nil {
		return apperror.Wrap(err)
	}
	o.Coordinate = point.ToCoordinate()
	o.FinishedAt = toTimePtr(nullFinished)
	return nil
}

func scanOrderItem(r RowScanner, oi *domain.OrderItem) error {
	pd := &oi.Product
	return r.Scan(
		&oi.ID, &oi.OrderID,
		&pd.ID,
		&oi.Price, &oi.Amount,
	)
}

func scanOrderItemJoined(r RowScanner, oi *domain.OrderItem) error {
	pd := &oi.Product
	return r.Scan(
		&oi.ID, &oi.OrderID,
		&pd.ID, &pd.Slug, &pd.Name,
		&oi.Price, &oi.Amount,
	)
}
