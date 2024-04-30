package postgres

import (
	"context"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
)

type pharmacyRepository struct {
	querier Querier
}

func (r *pharmacyRepository) Add(ctx context.Context, pharmacy domain.Pharmacy) (domain.Pharmacy, error) {
	q := `
		INSERT INTO pharmacies(name, manager_id, address, logo_url, coordinate, 
		pharmacist_name, pharmacist_license, pharmacist_phone)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
	` + pharmacyColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanPharmacy,
		pharmacy.Name, pharmacy.ManagerID, pharmacy.Address, pharmacy.LogoURL,
		postgis.NewPointFromCoordinate(pharmacy.Coordinate),
		pharmacy.PharmacistName, pharmacy.PharmacistLicense, pharmacy.PharmacistPhone,
	)
}
