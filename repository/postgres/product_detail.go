package postgres

import (
	"context"
	"medichat-be/domain"
)

type productDetailRepository struct {
	querier Querier
}

func (r *productDetailRepository) GetById(ctx context.Context, id int64) (domain.ProductDetails, error) {
	q := `
		SELECT ` + productDetailsColumns + `
		FROM product_details
		WHERE id = $1 AND deleted_at IS NULL
		`

	return queryOneFull(
		r.querier, ctx, q,
		scanProductDetails,
		id,
	)
}

func (r *productDetailRepository) Add(ctx context.Context, detail domain.ProductDetails) (domain.ProductDetails, error) {
	q := `
		INSERT INTO product_details(generic_name, content, manufacturer, description, product_classification, product_form, unit_in_pack, selling_unit, weight, height, length, width, composition)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING ` + productDetailsColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanProductDetails,
		detail.GenericName, detail.Content, detail.Manufacturer, detail.Description, detail.ProductClassification, detail.ProductForm, detail.UnitInPack, detail.SellingUnit, detail.Weight, detail.Height, detail.Length, detail.Width, detail.Composition,
	)
}

func (r *productDetailRepository) Update(ctx context.Context, detail domain.ProductDetails) (domain.ProductDetails, error) {
	q := `
		UPDATE product_details
		SET generic_name = $1,
			content = $2,
			manufacturer = $3,
			description = $4,
			product_classification = $5,
			product_form = $6,
			unit_in_pack= $7,
			selling_unit = $8,
			weight = $9,
			height = $10,
			length = $11,
			width = $12,
			composition = $13,
		WHERE id = $14 RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanProductDetails,
		detail.GenericName, detail.Content, detail.Manufacturer, detail.Description, detail.ProductClassification, detail.ProductForm, detail.UnitInPack, detail.SellingUnit, detail.Weight, detail.Height, detail.Length, detail.Composition, detail.Width,detail.ID,
	)
}


