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

func float64ScanDest(f *float64) []any {
	return []any{f}
}

func stringScanDest(s *string) []any {
	return []any{s}
}

var (
	accountColumns                = " id, email, email_verified, name, photo_url, role, account_type, profile_set "
	accountWithCredentialsColumns = " id, email, email_verified, name, photo_url, role, account_type, profile_set, hashed_password "
)

func scanAccountPharmacy(r RowScanner, p *domain.Account) error {
	if err := r.Scan(&p.ID, &p.Email, &p.EmailVerified, &p.Name, &p.PhotoURL, &p.Role, &p.AccountType, &p.ProfileSet); err != nil {
		return err
	}
	return nil
}

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
	shipmentMethodColumns = " id, name "
)

func scanShipmentMethod(r RowScanner, s *domain.ShipmentMethod) error {
	if err := r.Scan(&s.ID, &s.Name); err != nil {
		return err
	}
	return nil
}

var (
	pharmacyManagerColumns       = " id, account_id "
	pharmacyManagerJoinedColumns = `
		p.id,
		p.account_id, a.email, a.email_verified, a.role, a.account_type,
		a.name, a.photo_url, a.profile_set
	`
)

func scanPharmacyManager(r RowScanner, p *domain.PharmacyManager) error {
	if err := r.Scan(&p.ID, &p.Account.ID); err != nil {
		return err
	}
	return nil
}

func scanPharmacyManagerJoined(r RowScanner, p *domain.PharmacyManager) error {
	if err := r.Scan(&p.ID, &p.Account.ID, &p.Account.Email, &p.Account.EmailVerified,
		&p.Account.Role, &p.Account.AccountType, &p.Account.Name, &p.Account.PhotoURL,
		&p.Account.ProfileSet); err != nil {
		return err
	}
	return nil
}

var (
	pharmacyColumns               = " id, manager_id, name, address, coordinate, pharmacist_name, pharmacist_license, pharmacist_phone, slug "
	pharmacyJoinedColumns         = " p.id, p.manager_id, p.name, p.address, p.coordinate, p.pharmacist_name, p.pharmacist_license, p.pharmacist_phone, p.slug "
	pharmacyOperationColumns      = " id, pharmacy_id, day, start_time, end_time "
	PharmacyShipmentMethodColumns = " id, pharmacy_id, shipment_method_id "
)

func scanPharmacy(r RowScanner, p *domain.Pharmacy) error {
	var pos postgis.Point
	if err := r.Scan(
		&p.ID, &p.ManagerID, &p.Name, &p.Address, &pos,
		&p.PharmacistName, &p.PharmacistLicense,
		&p.PharmacistPhone, &p.Slug,
	); err != nil {
		return err
	}
	p.Coordinate = pos.ToCoordinate()
	return nil
}

func scanPharmacyWithDistance(r RowScanner, p *domain.Pharmacy) error {
	var pos postgis.Point
	var dis float64
	if err := r.Scan(
		&p.ID, &p.ManagerID, &p.Name, &p.Address, &pos,
		&p.PharmacistName, &p.PharmacistLicense,
		&p.PharmacistPhone, &p.Slug, &dis,
	); err != nil {
		return err
	}
	p.Coordinate = pos.ToCoordinate()
	p.Distance = &dis
	return nil
}

func ScanPharmacyOperation(r RowScanner, p *domain.PharmacyOperations) error {
	if err := r.Scan(&p.ID, &p.PharmacyID, &p.Day, &p.StartTime, &p.EndTime); err != nil {
		return err
	}
	return nil
}

func ScanPharmacyShipmentMethod(r RowScanner, s *domain.PharmacyShipmentMethods) error {
	if err := r.Scan(&s.ID, &s.PharmacyID, &s.ShipmentMethodID); err != nil {
		return err
	}
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
	productColumns        = " id, name, slug, product_detail_id, category_id, picture, is_active  "
	productDetailsColumns = " id, generic_name, composition, content, manufacturer, description, product_classification, product_form, unit_in_pack, selling_unit, weight, height, length, width  "
)

func scanProduct(r RowScanner, c *domain.Product) error {
	var nullPhotoUrl sql.NullString
	if err := r.Scan(
		&c.ID, &c.Name, &c.Slug, &c.ProductDetailId, &c.ProductCategoryId, &nullPhotoUrl, &c.IsActive,
	); err != nil {
		return err
	}
	c.Picture = toStringPtr(nullPhotoUrl)
	return nil
}

func scanProductDetails(r RowScanner, d *domain.ProductDetails) error {
	if err := r.Scan(
		&d.ID, &d.GenericName, &d.Composition, &d.Content, &d.Manufacturer, &d.Description, &d.ProductClassification, &d.ProductForm, &d.UnitInPack, &d.SellingUnit, &d.Weight, &d.Height, &d.Length, &d.Width,
	); err != nil {
		return err
	}

	return nil
}

var (
	chatsColumns = " chat_room_id, type, message, file, user_id, user_name, created_at  "
	roomsColumns = " user_id, doctor_id, end_at  "
)

func scanChats(r RowScanner, c *domain.Chat) error {
	if err := r.Scan(
		&c.RoomId, &c.Type, &c.Message, &c.File, &c.UserId, &c.UserName, &c.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

func scanRooms(r RowScanner, c *domain.Room) error {
	if err := r.Scan(
		&c.ID, &c.UserId, &c.DoctorId, &c.EndAt,
	); err != nil {
		return err
	}
	return nil
}

var (
	stockColumns = `
		id, product_id, pharmacy_id, stock, price
	`
	stockMutationColumns = `
		id, source_id, target_id, method, status, amount
	`

	selectStockJoined = `
		SELECT
			st.id,
			pd.id, pd.slug, pd.name,
			ph.id, ph.slug, ph.name,
			st.stock, st.price
		FROM stocks st
			JOIN pharmacies ph ON st.pharmacy_id = ph.id
			JOIN products pd ON st.product_id = pd.id
	`

	countStockJoined = `
		SELECT COUNT(st.id)
		FROM stocks st
			JOIN pharmacies ph ON st.pharmacy_id = ph.id
			JOIN products pd ON st.product_id = pd.id
	`

	selectStockMutationJoined = `
		SELECT 
			sm.id, 
			st1.id, ph1.id, ph1.slug, ph1.name,
			st2.id, ph2.id, ph2.slug, ph2.name,
			pd.id, pd.slug, pd.name,
			sm.method, sm.status, sm.amount, sm.created_at
		FROM stock_mutations sm
			JOIN stocks st1 ON sm.source_id = st1.id
			JOIN stocks st2 ON sm.target_id = st2.id
			JOIN pharmacies ph1 ON st1.pharmacy_id = ph1.id
			JOIN pharmacies ph2 ON st2.pharmacy_id = ph2.id
			JOIN products pd ON st1.product_id = pd.id
	`

	countStockMutationJoined = `
		SELECT COUNT(sm.id)
		FROM stock_mutations sm
			JOIN stocks st1 ON sm.source_id = st1.id
			JOIN stocks st2 ON sm.target_id = st2.id
			JOIN pharmacies ph1 ON st1.pharmacy_id = ph1.id
			JOIN pharmacies ph2 ON st2.pharmacy_id = ph2.id
			JOIN products pd ON st1.product_id = pd.id
	`
)

func scanStock(r RowScanner, s *domain.Stock) error {
	return r.Scan(
		&s.ID, &s.ProductID, &s.PharmacyID, &s.Stock, &s.Price,
	)
}

func scanStockMutation(r RowScanner, sm *domain.StockMutation) error {
	return r.Scan(
		&sm.ID, &sm.SourceID, &sm.TargetID, &sm.Method, &sm.Status, &sm.Amount,
	)
}

func scanStockJoined(r RowScanner, s *domain.StockJoined) error {
	return r.Scan(
		&s.ID,
		&s.Product.ID, &s.Product.Slug, &s.Product.Name,
		&s.Pharmacy.ID, &s.Pharmacy.Slug, &s.Pharmacy.Name,
		&s.Stock, &s.Price,
	)
}

func scanStockMutationJoined(r RowScanner, sm *domain.StockMutationJoined) error {
	return r.Scan(
		&sm.ID,
		&sm.Source.ID, &sm.Source.PharmacyID, &sm.Source.PharmacySlug, &sm.Source.PharmacyName,
		&sm.Target.ID, &sm.Target.PharmacyID, &sm.Target.PharmacySlug, &sm.Target.PharmacyName,
		&sm.Product.ID, &sm.Product.Slug, &sm.Product.Name,
		&sm.Method, &sm.Status, &sm.Amount, &sm.Timestamp,
	)
}

var (
	paymentColumns = `
		id, invoice_number, user_id, file_url, is_confirmed, amount
	`

	selectPaymentJoined = `
		SELECT
			p.id, p.invoice_number, u.id, a.name, 
			p.file_url, p.is_confirmed, p.amount
		FROM
			payments p
			JOIN users u ON p.user_id = u.id
			JOIN accounts a ON u.account_id = a.id
	`

	countPaymentJoined = `
		SELECT COUNT(p.id)
		FROM
			payments p
			JOIN users u ON p.user_id = u.id
	`
)

func scanPayment(r RowScanner, p *domain.Payment) error {
	nullURL := sql.NullString{}
	if err := r.Scan(
		&p.ID, &p.InvoiceNumber, &p.User.ID, &nullURL, &p.IsConfirmed, &p.Amount,
	); err != nil {
		return err
	}
	p.FileURL = toStringPtr(nullURL)
	return nil
}

func scanPaymentJoined(r RowScanner, p *domain.Payment) error {
	nullURL := sql.NullString{}
	if err := r.Scan(
		&p.ID, &p.InvoiceNumber, &p.User.ID, &p.User.Name, &nullURL, &p.IsConfirmed, &p.Amount,
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
			u.id, a.name,
			ph.id, ph.slug, ph.name, ph.manager_id,
			py.id, py.invoice_number,
			sm.id, sm.name,
			o.address, o.coordinate, 
			o.n_items, o.subtotal, o.shipment_fee, o.total, 
			o.status, o.ordered_at, o.finished_at
		FROM orders o
			JOIN users u ON o.user_id = u.id
			JOIN accounts a ON u.account_id = a.id
			JOIN pharmacies ph ON o.pharmacy_id = ph.id
			JOIN payments py ON o.payment_id = py.id
			JOIN shipment_methods sm ON o.shipment_method_id = sm.id
	`

	countOrderJoined = `
		SELECT COUNT(o.id)
		FROM orders o
			JOIN users u ON o.user_id = u.id
			JOIN pharmacies ph ON o.pharmacy_id = ph.id
			JOIN payments py ON o.payment_id = py.id
			JOIN shipment_methods sm ON o.shipment_method_id = sm.id
	`

	orderItemColumns = `
		id, order_id, product_id, price, amount
	`

	selectOrderItemJoined = `
		SELECT
			oi.id, oi.order_id,
			pd.id, pd.slug, pd.name, pd.picture,
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
		&ph.ID, &ph.Slug, &ph.Name, &ph.ManagerID,
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
		&pd.ID, &pd.Slug, &pd.Name, &pd.PhotoURL,
		&oi.Price, &oi.Amount,
	)
}
