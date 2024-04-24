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

func (r *categoryRepository) GetCategories(ctx context.Context, query domain.CategoriesQuery) ([]domain.CategoryWithParentName, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}
	offset := (query.Page - 1) * query.Limit

	sb.WriteString(`
		SELECT ` + categoryWithParentNameColumns + `
		FROM categories c LEFT JOIN categories c2 
			ON c.parent_id = c2.id
		WHERE c.deleted_at IS NULL
	`)

	sb.WriteString(` AND c.name ILIKE '%' || @name || '%' `)
	args["name"] = query.Term

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
		INSERT INTO categories(parent_id, name)
		VALUES
		($1, $2)
		RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		category.ParentID, category.Name,
	)
}

func (r *categoryRepository) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	q := `
		UPDATE categories
		SET name = $1, 
			parent_id = $2
		WHERE id = $3
		RETURNING ` + categoryColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanCategory,
		category.Name, category.ParentID, category.ID,
	)
}

func (r *categoryRepository) SoftDeleteById(ctx context.Context, id int64) error {
	q := `
		UPDATE categories
		SET deleted_at = now()
		WHERE id = $1 
	`

	return exec(
		r.querier, ctx, q,
		id,
	)
}

func (r *categoryRepository) BulkSoftDelete(ctx context.Context, ids []int64) error {
	sb := strings.Builder{}
	params := make([]interface{}, len(ids))
	sb.WriteString(`
		UPDATE categories
		SET deleted_at = now()
		WHERE id IN (`)

	for i := 0; i < len(ids); i++ {
		params[i] = ids[i]
		sb.WriteString(fmt.Sprintf("$%d", i+1))
		if i != len(ids)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	return exec(
		r.querier, ctx, sb.String(),
		params...,
	)
}
