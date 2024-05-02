package postgres

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/domain"
	"strings"

	"github.com/jackc/pgx/v5"
)

type paymentRepository struct {
	querier Querier
}

func (r *paymentRepository) buildListQuery(sel string, dets domain.PaymentListDetails) (strings.Builder, pgx.NamedArgs) {
	var sb strings.Builder
	args := pgx.NamedArgs{}

	sb.WriteString(sel)
	sb.WriteString(`
		WHERE deleted_at IS NULL
	`)

	if dets.IsConfirmed != nil {
		sb.WriteString(`
			AND is_confirmed = @isConfirmed
		`)
		args["isConfirmed"] = *dets.IsConfirmed
	}

	return sb, args
}

func (r *paymentRepository) GetPageInfo(ctx context.Context, dets domain.PaymentListDetails) (domain.PageInfo, error) {
	sb, args := r.buildListQuery(`
		SELECT COUNT(*)
		FROM payments
	`, dets)

	count, err := queryOne(
		r.querier, ctx, sb.String(),
		int64ScanDest,
		args,
	)
	if err != nil {
		return domain.PageInfo{}, apperror.Wrap(err)
	}

	return domain.PageInfo{
		CurrentPage:  dets.Page,
		ItemsPerPage: dets.Limit,
		ItemCount:    count,
		PageCount:    int((count - 1 + int64(dets.Limit)) / int64(dets.Limit)),
	}, nil

}

func (r *paymentRepository) List(ctx context.Context, dets domain.PaymentListDetails) ([]domain.Payment, error) {
	sb, args := r.buildListQuery(`
		SELECT `+paymentColumns+`
		FROM payments
	`, dets)
	offset := (dets.Page - 1) * dets.Limit

	fmt.Fprintf(
		&sb,
		` OFFSET %d LIMIT %d `,
		offset,
		dets.Limit,
	)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanPayment,
		args,
	)
}

func (r *paymentRepository) GetByID(ctx context.Context, id int64) (domain.Payment, error) {
	q := `
		SELECT ` + paymentColumns + `
		FROM payments
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPayment,
		id,
	)
}

func (r *paymentRepository) GetByInvoiceNumber(ctx context.Context, num string) (domain.Payment, error) {
	q := `
		SELECT ` + paymentColumns + `
		FROM payments
		WHERE invoice_number = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanPayment,
		num,
	)
}

func (r *paymentRepository) Add(ctx context.Context, p domain.Payment) (domain.Payment, error) {
	q := `
		INSERT INTO payments(invoice_number, file_url, is_confirmed, amount)
		VALUES
		($1, $2, $3, $4)
		RETURNING ` + paymentColumns

	nullURL := fromStringPtr(p.FileURL)
	return queryOneFull(
		r.querier, ctx, q,
		scanPayment,
		p.InvoiceNumber, nullURL, p.IsConfirmed, p.Amount,
	)
}

func (r *paymentRepository) Update(ctx context.Context, p domain.Payment) (domain.Payment, error) {
	q := `
		UPDATE payments
		SET file_url = $2,
			is_confirmed = $3,
			updated_at = now()
		WHERE id = $1
	`

	nullURL := fromStringPtr(p.FileURL)
	err := execOne(
		r.querier, ctx, q,
		p.Amount, nullURL,
	)
	if err != nil {
		return domain.Payment{}, apperror.Wrap(err)
	}

	return p, nil
}
