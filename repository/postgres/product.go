package postgres

import (
	"context"
	"fmt"
	"medichat-be/domain"
	"strings"
)

type productRepository struct {
	querier Querier
}

func (r *productRepository) GetByName(ctx context.Context, name string) (domain.Product, error) {
	q := `
		SELECT ` + productColumns + `
		FROM products
		WHERE name ilike $1
		AND deleted_at IS NULL
		`

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		name,
	)
}


func (r *productRepository) GetById(ctx context.Context, id int64) (domain.Product, error) {
	q := `
		SELECT ` + productColumns + `
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
		`

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		id,
	)
}

func (r *productRepository) Add(ctx context.Context, product domain.Product) (domain.Product, error) {
	q := `
		INSERT INTO products(name,category_id, product_detail_id, picture, slug, is_active)
		VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING ` + productColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		product.Name, product.ProductCategoryId, product.ProductDetailId, product.Picture, product.Slug, product.IsActive,
	)
}

func (r *productRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	q := `
		UPDATE products
		SET name = $1,
			category_id = $2,
			product_detail_id = $3,
			picture = $4
			slug = $5
			is_active = $6
		WHERE id = $5 RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		product.Name, product.ProductCategoryId, product.ProductDetailId, product.Picture, product.Slug, product.IsActive,
	)
}

func (r *productRepository) SoftDeleteBySlug(ctx context.Context, slug string) error {
	q := `
		UPDATE products
		SET deleted_at = now()
		WHERE slug = $1
	`

	return exec(
		r.querier, ctx, q,
		slug,
	)
}

func (r *productRepository) BulkSoftDeleteBySlug(ctx context.Context, slugs []string) error {
	sb := strings.Builder{}
	params := make([]interface{}, len(slugs))
	sb.WriteString(`
		UPDATE products
		SET deleted_at = now()
		WHERE slug IN (`)

	for i := 0; i < len(slugs); i++ {
		params[i] = slugs[i]
		sb.WriteString(fmt.Sprintf("$%d", i+1))
		if i != len(slugs)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	return exec(
		r.querier, ctx, sb.String(),
		params...,
	)
}


