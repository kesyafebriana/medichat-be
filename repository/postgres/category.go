package postgres

import (
	"context"
	"fmt"
	"medichat-be/domain"
	"strings"

	"github.com/jackc/pgx/v5"
)

type categoryRepository struct {
	querier Querier
}

func (r *categoryRepository) GetCategoriesWithParentName(ctx context.Context, query domain.CategoriesQuery) ([]domain.CategoryWithParentName, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}
	offset := (query.Page - 1) * query.Limit

	sb.WriteString(`
		SELECT ` + categoryWithParentNameColumns + `
		FROM categories c LEFT JOIN categories c2
			ON c.parent_id = c2.id
		WHERE c.deleted_at IS NULL
	`)

	if query.Term != "" {
		sb.WriteString(` AND c.name ILIKE '%' || @name || '%' `)
		args["name"] = query.Term
	}

	if query.ParentSlug != "" {
		sb.WriteString(` AND c2.slug = @parentSlug `)
		args["parentSlug"] = query.ParentSlug
	}

	if query.Level != 0 {
		var key string
		if query.Level == 2 {
			key = "NOT"
		}
		fmt.Fprintf(&sb, `AND c.parent_id IS %s NULL `, key)
	}

	if query.ParentId != nil {
		sb.WriteString(` AND c.parent_id = @parentId `)
		args["parentId"] = *query.ParentId
	}

	if query.SortBy != domain.CategorySortByParent {
		fmt.Fprintf(&sb, " ORDER BY %s %s", query.SortBy, query.SortType)
	}

	if query.Limit != 0 {
		fmt.Fprintf(&sb, " OFFSET %d LIMIT %d ", offset, query.Limit)
	}

	return queryFull(
		r.querier, ctx, sb.String(),
		scanCategoryWithParentName,
		args,
	)
}

func (r *categoryRepository) GetCategories(ctx context.Context, query domain.CategoriesQuery) ([]domain.Category, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}
	offset := (query.Page - 1) * query.Limit

	sb.WriteString(`
		SELECT ` + categoryColumns + `
		FROM categories 
		WHERE deleted_at IS NULL
	`)

	sb.WriteString(` AND name ILIKE '%' || @name || '%' `)
	args["name"] = query.Term

	if query.Level != 0 {
		var key string
		if query.Level == 2 {
			key = "NOT"
		}
		fmt.Fprintf(&sb, `AND parent_id IS %s NULL `, key)
	}

	if query.ParentId != nil {
		sb.WriteString(` AND parent_id = @parentId `)
		args["parentId"] = *query.ParentId
	}

	if query.SortBy != domain.CategorySortByParent {
		fmt.Fprintf(&sb, " ORDER BY %s %s", query.SortBy, query.SortType)
	}

	if query.Limit != 0 {
		fmt.Fprintf(&sb, " OFFSET %d LIMIT %d ", offset, query.Limit)
	}

	return queryFull(
		r.querier, ctx, sb.String(),
		scanCategory,
		args,
	)
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (domain.Category, error) {
	q := `
		SELECT ` + categoryColumns + `
		FROM categories
		WHERE deleted_at IS NULL AND slug = $1
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		slug,
	)
}

func (r *categoryRepository) GetBySlugWithParentName(ctx context.Context, slug string) (domain.CategoryWithParentName, error) {
	q := `
		SELECT ` + categoryWithParentNameColumns + `
		FROM categories c LEFT JOIN categories c2 
			ON c.parent_id = c2.id
		WHERE c.deleted_at IS NULL AND c.slug = $1
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanCategoryWithParentName,
		slug,
	)
}

func (r *categoryRepository) GetPageInfo(ctx context.Context, query domain.CategoriesQuery) (domain.PageInfo, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}

	sb.WriteString(`
		SELECT COUNT(*) as total_data 
		FROM categories c LEFT JOIN categories c2 
			ON c.parent_id = c2.id
		WHERE c.deleted_at IS NULL
	`)

	if query.Term != "" {
		sb.WriteString(` AND c.name ILIKE '%' || @name || '%' `)
		args["name"] = query.Term
	}

	if query.ParentSlug != "" {
		sb.WriteString(` AND c2.slug = @parentSlug `)
		args["parentSlug"] = query.ParentSlug
	}

	if query.Level != 0 {
		var key string
		if query.Level == 2 {
			key = "NOT"
		}
		fmt.Fprintf(&sb, `AND c.parent_id IS %s NULL `, key)
	}

	if query.ParentId != nil {
		sb.WriteString(` AND c.parent_id = @parentId `)
		args["parentId"] = *query.ParentId
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

func (r *categoryRepository) GetByName(ctx context.Context, name string) (domain.Category, error) {
	q := `
		SELECT ` + categoryColumns + `
		FROM categories
		WHERE name = $1 AND deleted_at IS NULL
		`

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		name,
	)
}

func (r *categoryRepository) GetById(ctx context.Context, id int64) (domain.Category, error) {
	q := `
		SELECT ` + categoryColumns + `
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL
		`

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		id,
	)
}

func (r *categoryRepository) Add(ctx context.Context, category domain.Category) (domain.Category, error) {
	q := `
		INSERT INTO categories(parent_id, name, slug, photo_url)
		VALUES
		($1, $2, $3, $4)
		RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		category.ParentID, category.Name, category.Slug, category.PhotoUrl,
	)
}

func (r *categoryRepository) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	q := `
		UPDATE categories
		SET name = $1, 
			parent_id = $2,
			slug = $3,
			photo_url = $4
		WHERE id = $5 RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		category.Name, category.ParentID, category.Slug, category.PhotoUrl, category.ID,
	)
}

func (r *categoryRepository) SoftDeleteBySlug(ctx context.Context, slug string) error {
	q := `
		UPDATE categories
		SET deleted_at = now()
		WHERE slug = $1 
	`

	return exec(
		r.querier, ctx, q,
		slug,
	)
}

func (r *categoryRepository) BulkSoftDeleteBySlug(ctx context.Context, slugs []string) error {
	sb := strings.Builder{}
	params := make([]interface{}, len(slugs))
	sb.WriteString(`
		UPDATE categories
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
