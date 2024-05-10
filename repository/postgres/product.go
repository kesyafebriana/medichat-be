package postgres

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/domain"
	"strings"

	"github.com/jackc/pgx/v5"
)

type productRepository struct {
	querier Querier
}

func (r *productRepository) GetProductsFromArea(ctx context.Context, query domain.ProductsQuery) ([]domain.Product, error) {
	sb := strings.Builder{}
	args := make([]any, 0)
	var idx = 1
	offset := (query.Page - 1) * query.Limit

	_, err := fmt.Fprintf(&sb, ` SELECT DISTINCT p.id, p.name, p.slug, p.product_detail_id, p.category_id, p.picture, p.is_active from products p inner join
		(select pharma.id as pharmacy_id,pharma.name, stock.product_id as product_id from
			(SELECT id, name, address, coordinate FROM pharmacies
			WHERE deleted_at IS null AND ST_DWithin(coordinate, ST_MakePoint($%d, $%d)::geography, 25000)) as pharma
			inner join (select s.id,s.product_id,s.pharmacy_id from stocks s where s.deleted_at is null and stock >=0) as stock
		on pharma.id = stock.pharmacy_id) as pp
		on p.id = pp.product_id
	`, idx, idx+1)
	idx += 2
	args = append(args, *query.Longitude, *query.Latitude)
	if err != nil {
		return []domain.Product{}, apperror.Wrap(err)
	}

	if query.Term != "" {

		fmt.Fprintf(&sb, `AND c.keyword ILIKE $%d`, idx)
		args = append(args, query.Term)
		idx += 1

	}

	if query.SortBy != domain.CategorySortByParent {
		query.SortBy = "p." + query.SortBy
		fmt.Fprintf(&sb, " ORDER BY %v %v", query.SortBy, query.SortType)

	}

	if query.Limit != 0 {
		fmt.Fprintf(&sb, " OFFSET $%d LIMIT $%d ", idx, idx+1)
		args = append(args, offset, query.Limit)
		idx += 2
	}

	return queryFull(
		r.querier, ctx, sb.String(),
		scanProduct,
		args...,
	)
}

func (r *productRepository) GetPageInfoFromArea(ctx context.Context, query domain.ProductsQuery) (domain.PageInfo, error) {
	sb := strings.Builder{}
	args := make([]any, 0)
	var idx = 1

	_, err := fmt.Fprintf(&sb, ` SELECT COUNT(DISTINCT p.id) as total_data from products p inner join
		(select pharma.id as pharmacy_id,pharma.name, stock.product_id as product_id from
			(SELECT id, name, address, coordinate FROM pharmacies
			WHERE deleted_at IS null AND ST_DWithin(coordinate, ST_MakePoint($%d, $%d)::geography, 25000)) as pharma
			inner join (select s.id,s.product_id,s.pharmacy_id from stocks s where s.deleted_at is null and stock >=0) as stock
		on pharma.id = stock.pharmacy_id) as pp
		on p.id = pp.product_id
	`, idx, idx+1)
	idx += 2
	args = append(args, *query.Longitude, *query.Latitude)
	if err != nil {
		return domain.PageInfo{}, apperror.Wrap(err)
	}

	if query.Term != "" {
		fmt.Fprintf(&sb, `AND c.name ILIKE $%d`, idx)
		args = append(args, query.Term)
		idx += 1
	}

	var totalData int64
	row := r.querier.QueryRowContext(ctx, sb.String(), args...)
	err = row.Scan(&totalData)

	if err != nil {
		return domain.PageInfo{}, nil
	}

	return domain.PageInfo{
		CurrentPage: int(query.Page),
		ItemCount:   totalData,
	}, nil
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
		sb.WriteString(` AND (p.keyword ILIKE '%' || @name || '%' `)
		sb.WriteString(` OR p.name ILIKE '%' || @name || '%') `)
		args["name"] = query.Term
	}

	if query.CategoryID != nil {
		sb.WriteString(` AND p.category_id = @categoryID `)
		args["categoryID"] = *query.CategoryID
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
		sb.WriteString(` AND (c.keyword ILIKE '%' || @name || '%' `)
		sb.WriteString(` OR c.name ILIKE '%' || @name || '%') `)
		args["name"] = query.Term
	}

	if query.CategoryID != nil {
		sb.WriteString(` AND c.category_id = @categoryID `)
		args["categoryID"] = *query.CategoryID
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
	picture := ""
	if product.Picture != nil {
		picture = *product.Picture
	}
	q := `
		INSERT INTO products(name,category_id, product_detail_id, picture, slug, is_active, keyword)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING ` + productColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		product.Name, product.ProductCategoryId, product.ProductDetailId, picture, product.Slug, product.IsActive, product.KeyWord,
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
			keyword = $7
		WHERE id = $8 RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanProduct,
		product.Name, product.ProductCategoryId, product.ProductDetailId, product.Picture, product.Slug, product.IsActive, product.KeyWord, product.ID,
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
