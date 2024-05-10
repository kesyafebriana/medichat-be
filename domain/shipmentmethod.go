package domain

import "context"

const (
	ShipmentOfficialInstantID  = 1
	ShipmentOfficialInstantFee = 2500

	ShipmentOfficialSameDayID  = 2
	ShipmentOfficialSameDayFee = 1000
)

type ShipmentMethod struct {
	ID   int64
	Name string
}

type ShipmentMethodRepository interface {
	GetShipmentMethodById(ctx context.Context, id int64) (ShipmentMethod, error)
}
