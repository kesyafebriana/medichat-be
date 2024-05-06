package domain

import "context"

type ShipmentMethod struct {
	ID   int64
	Name string
}

type ShipmentMethodRepository interface {
	GetShipmentMethodById(ctx context.Context, id int64) (ShipmentMethod, error)
}
