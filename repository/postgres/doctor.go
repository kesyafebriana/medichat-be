package postgres

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"strings"
)

type doctorRepository struct {
	querier Querier
}

func (r *doctorRepository) List(
	ctx context.Context,
	det domain.DoctorListDetails,
) ([]domain.Doctor, error) {
	var sb strings.Builder
	var args = make([]any, 0)
	var idx = 1

	sb.WriteString(`
		SELECT ` + doctorJoinedColumns + `
		FROM doctors d JOIN accounts a ON d.account_id = a.id
			JOIN specializations s ON d.specialization_id = s.id
		WHERE d.deleted_at IS NULL
	`)

	if det.SpecializationID != nil {
		fmt.Fprintf(&sb, ` AND d.specialization_id $%d 
		`, idx)
		idx++
		args = append(args, *det.SpecializationID)
	}
	if det.Name != nil {
		fmt.Fprintf(&sb, ` AND a.name ILIKE $%d 
		`, idx)
		idx++
		args = append(args, *det.Name)
	}
	if det.Gender != nil {
		fmt.Fprintf(&sb, ` AND d.gender = $%d 
		`, idx)
		idx++
		args = append(args, *det.Gender)
	}
	if det.MinPrice != nil {
		fmt.Fprintf(&sb, ` AND d.price >= $%d 
		`, idx)
		idx++
		args = append(args, *det.MinPrice)
	}
	if det.MaxPrice != nil {
		fmt.Fprintf(&sb, ` AND d.price <= $%d 
		`, idx)
		idx++
		args = append(args, *det.MaxPrice)
	}
	if det.MinYearExperience != nil {
		fmt.Fprintf(&sb, ` AND now()::date - d.start_work_date >= $%d `, idx)
		idx++
		args = append(args, *det.MinYearExperience*365)
	}

	sortCol := "a.name"
	sortAsc := det.SortAsc
	switch det.SortBy {
	case constants.DoctorSortByName:
		sortCol = "a.name"
	case constants.DoctorSortByPrice:
		sortCol = "d.price"
	case constants.DoctorSortByStartWorkDate:
		sortCol = "d.start_work_date"
	}

	if det.CursorID != nil && det.Cursor != nil {
		fmt.Fprintf(
			&sb,
			` AND (%s, d.id) %s ($%d, $%d) `,
			sortCol, getSortCursorCmp(sortAsc), idx, idx+1,
		)
		idx += 2
		args = append(args, det.Cursor, *det.CursorID)
	}

	fmt.Fprintf(
		&sb,
		` ORDER BY is_active desc, %s %s, d.id %s`,
		sortCol,
		getSortOrder(sortAsc),
		getSortOrder(det.SortAsc),
	)

	fmt.Fprintf(&sb, ` LIMIT $%d `, idx)
	idx++
	args = append(args, det.Limit)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanDoctorJoined,
		args...,
	)
}

func (r *doctorRepository) GetByID(
	ctx context.Context,
	id int64,
) (domain.Doctor, error) {
	q := `
		SELECT ` + doctorJoinedColumns + `
		FROM doctors d JOIN accounts a ON d.account_id = a.id
			JOIN specializations s ON d.specialization_id = s.id
		WHERE d.id = $1
			AND d.deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanDoctorJoined,
		id,
	)
}

func (r *doctorRepository) GetByIDAndLock(
	ctx context.Context,
	id int64,
) (domain.Doctor, error) {
	q := `
		SELECT ` + doctorJoinedColumns + `
		FROM doctors d JOIN accounts a ON d.account_id = a.id
			JOIN specializations s ON d.specialization_id = s.id
		WHERE d.id = $1
			AND d.deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanDoctorJoined,
		id,
	)
}

func (r *doctorRepository) IsExistByID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM doctors
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

func (r *doctorRepository) GetByAccountID(
	ctx context.Context,
	id int64,
) (domain.Doctor, error) {
	q := `
		SELECT ` + doctorJoinedColumns + `
		FROM doctors d JOIN accounts a ON d.account_id = a.id
			JOIN specializations s ON d.specialization_id = s.id
		WHERE d.account_id = $1
			AND d.deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanDoctorJoined,
		id,
	)
}

func (r *doctorRepository) GetByAccountIDAndLock(
	ctx context.Context,
	id int64,
) (domain.Doctor, error) {
	q := `
		SELECT ` + doctorJoinedColumns + `
		FROM doctors d JOIN accounts a ON d.account_id = a.id
			JOIN specializations s ON d.specialization_id = s.id
		WHERE d.account_id = $1
			AND d.deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanDoctorJoined,
		id,
	)
}

func (r *doctorRepository) IsExistByAccountID(
	ctx context.Context,
	id int64,
) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT id
			FROM doctors
			WHERE account_id = $1
				AND deleted_at IS NULL
		)
	`

	return queryOne(
		r.querier, ctx, q,
		boolScanDest,
		id,
	)
}

func (r *doctorRepository) Add(
	ctx context.Context,
	d domain.Doctor,
) (domain.Doctor, error) {
	q := `
		INSERT INTO doctors(
			account_id, specialization_id, str, work_location, gender,
			phone_number, is_active, start_work_date, price, certificate_url
		)		
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING ` + doctorColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanDoctor,
		d.Account.ID, d.Specialization.ID, d.STR, d.WorkLocation, d.Gender,
		d.PhoneNumber, d.IsActive, d.StartWorkDate, d.Price, d.CertificateURL,
	)
}

func (r *doctorRepository) Update(
	ctx context.Context,
	d domain.Doctor,
) (domain.Doctor, error) {
	q := `
		UPDATE doctors
		SET work_location = $2,
			gender = $3,
			phone_number = $4,
			price = $5,
			is_active = $6,
			updated_at = now()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	err := execOne(
		r.querier, ctx, q,
		d.ID, d.WorkLocation, d.Gender, d.PhoneNumber, d.Price, d.IsActive,
	)
	if err != nil {
		return domain.Doctor{}, apperror.Wrap(err)
	}

	return d, nil
}
