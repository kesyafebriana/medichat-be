package postgres

import (
	"context"
	"medichat-be/domain"
)

type shipmentMethodRepository struct {
	querier Querier
}

func (r *shipmentMethodRepository) GetShipmentMethodById(ctx context.Context, id int64) (domain.ShipmentMethod, error) {
	q := `
		SELECT` + shipmentMethodColumns + `
		FROM shipment_methods
		WHERE deleted_at IS NULL AND id = $1
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanShipmentMethod,
		id,
	)
}
