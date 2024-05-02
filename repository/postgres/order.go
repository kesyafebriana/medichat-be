package postgres

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
	"strings"

	"github.com/jackc/pgx/v5"
)

type orderRepository struct {
	querier Querier
}

func (r *orderRepository) buildListQuery(sel string, dets domain.OrderListDetails) (*strings.Builder, pgx.NamedArgs) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}

	sb.WriteString(sel)

	if dets.PharmacyManagerID != nil {
		sb.WriteString(`
			JOIN pharmacy_managers pm ON ph.manager_id = pm.id
		`)
	}

	sb.WriteString(`
		WHERE deleted_at IS NULL
	`)

	if dets.UserID != nil {
		sb.WriteString(`
			AND u.id = @userID
		`)
		args["userID"] = *dets.UserID
	}
	if dets.PharmacyID != nil {
		sb.WriteString(`
			AND ph.id = @pharmacyID
		`)
		args["pharmacyID"] = *dets.PharmacyID
	}
	if dets.PharmacyManagerID != nil {
		sb.WriteString(`
			AND pm.id = @pmID
		`)
		args["pmID"] = *dets.PharmacyManagerID
	}
	if dets.Status != nil {
		sb.WriteString(`
			AND o.status = @status
		`)
		args["status"] = *dets.Status
	}

	return &sb, args
}

func (r *orderRepository) GetPageInfo(ctx context.Context, dets domain.OrderListDetails) (domain.PageInfo, error) {
	sb, args := r.buildListQuery(countOrderJoined, dets)

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

func (r *orderRepository) List(ctx context.Context, dets domain.OrderListDetails) ([]domain.Order, error) {
	sb, args := r.buildListQuery(countOrderJoined, dets)
	offset := (dets.Page - 1) * dets.Limit

	sb.WriteString(`
		ORDER BY ordered_at DESC
	`)

	fmt.Fprintf(
		sb,
		` OFFSET %d LIMIT %d `,
		offset,
		dets.Limit,
	)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanOrderJoined,
		args,
	)
}

func (r *orderRepository) GetByID(ctx context.Context, id int64) (domain.Order, error) {
	q := selectOrderJoined + `
		WHERE deleted_at IS NULL
			AND id = $1
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanOrderJoined,
		id,
	)
}

func (r *orderRepository) GetByIDAndLock(ctx context.Context, id int64) (domain.Order, error) {
	q := selectOrderJoined + `
		WHERE deleted_at IS NULL
			AND id = $1
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanOrderJoined,
		id,
	)
}

func (r *orderRepository) Add(ctx context.Context, o domain.Order) (domain.Order, error) {
	q := `
		INSERT orders(
			user_id, pharmacy_id, payment_id, shipment_method_id,
			address, coordinate,
			n_items, subtotal, shipment_fee, total,
			status, ordered_at, finished_at
		)
		VALUES
		(
			$1, $2, $3, $4,
			$5, $6,
			$7, $8, $9, $10,
			$11, $12, $13
		)
		RETURNING ` + orderColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanOrder,
		o.User.ID, o.Pharmacy.ID, o.Payment.ID, o.ShipmentMethod.ID,
		o.Address, postgis.NewPointFromCoordinate(o.Coordinate),
		o.NItems, o.Subtotal, o.ShipmentFee, o.Total,
		o.Status, o.OrderedAt, fromTimePtr(o.FinishedAt),
	)
}

func (r *orderRepository) UpdateStatusByID(ctx context.Context, id int64, status string) error {
	q := `
		UPDATE orders
		SET status = $2,
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id, status,
	)
}

func (r *orderRepository) UpdateStatusByPaymentID(ctx context.Context, id int64, status string) error {
	q := `
		UPDATE orders
		SET status = $2,
			updated_at = now()
		WHERE payment_id = $1
	`

	return exec(
		r.querier, ctx, q,
		id, status,
	)
}

func (r *orderRepository) ListItemsByOrderID(ctx context.Context, id int64) ([]domain.OrderItem, error) {
	q := selectOrderItemJoined + `
		WHERE AND oi.id = $1
			AND oi.deleted_at IS NULL
	`

	return queryFull(
		r.querier, ctx, q,
		scanOrderItemJoined,
		id,
	)
}

func (r *orderRepository) AddItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	q := `
		INSERT INTO order_items(order_id, product_id, price, amount)
		VALUES
		($1, $2, $3, $4)
		RETURNING ` + orderItemColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanOrderItem,
		item.OrderID, item.Product.ID, item.Price, item.Amount,
	)
}
