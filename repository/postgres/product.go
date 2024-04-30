package postgres

import (
	"context"
	"fmt"
	"medichat-be/domain"
	"strings"

	"github.com/jackc/pgx/v5"
)

type productRepository struct {
	querier Querier
}

func (r *productRepository) GetProducts(ctx context.Context, query domain.ProductsQuery) ([]domain.Product, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}
	offset := (query.Page - 1) * query.Limit

	sb.WriteString(`
		SELECT ` + productColumns + `
		FROM products p
		WHERE p.deleted_at IS NULL
	`)

	if query.Term != "" {
		sb.WriteString(` AND c.keyword ILIKE '%' || @name || '%' `)
		args["name"] = query.Term
	}

	if query.SortBy != domain.CategorySortByParent {
		fmt.Fprintf(&sb, " ORDER BY %s %s", query.SortBy, query.SortType)
	}

	if query.Limit != 0 {
		fmt.Fprintf(&sb, " OFFSET %d LIMIT %d ", offset, query.Limit)
	}

	return queryFull(
		r.querier, ctx, sb.String(),
		scanProduct,
		args,
	)
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (domain.Product, error) {
	q := `
		SELECT ` + productColumns + `
		FROM products
		WHERE slug = $1
		AND deleted_at IS NULL
		LIMIT 1
		`

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		slug,
	)
}

func (r *productRepository) GetByName(ctx context.Context, name string) (domain.Product, error) {
	q := `
		SELECT ` + productColumns + `
		FROM products
		WHERE name ilike $1
		AND deleted_at IS NULL
		LIMIT 1
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

func (r *productRepository) GetPageInfo(ctx context.Context, query domain.ProductsQuery) (domain.PageInfo, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}

	sb.WriteString(`
		SELECT COUNT(*) as total_data
		FROM products c
		WHERE c.deleted_at IS NULL
	`)

	if query.Term != "" {
		sb.WriteString(` AND c.name ILIKE '%' || @name || '%' `)
		args["name"] = query.Term
	}

	var totalData int64
	row := r.querier.QueryRowContext(ctx, sb.String(), args)
	err := row.Scan(&totalData)

	if err != nil {
		return domain.PageInfo{}, nil
	}

	return domain.PageInfo{
		CurrentPage: int(query.Page),
		ItemCount:   totalData,
	}, nil
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


