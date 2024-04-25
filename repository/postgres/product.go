package postgres

import (
	"context"
	"medichat-be/domain"
)

type productRepository struct {
	querier Querier
}

func (r *productRepository) GetByName(
	ctx context.Context,
	name string,
) ([]domain.Product, error) {
	q := `
		SELECT *
		FROM products
		WHERE name ilike $1
			AND deleted_at IS NULL
	`

	return queryFull(
		r.querier, ctx, q, func(rs RowScanner, t *domain.Product) error {
			scanDests := []any{
		&t.ID, &t.Name, &t.Picture, &t.ProductCategoryId, &t.ProductDetailId, t.IsActive,
	}
			return rs.Scan(scanDests...)
		},name,
	)
}
