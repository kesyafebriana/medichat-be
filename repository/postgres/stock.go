package postgres

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/domain"
	"strings"

	"github.com/jackc/pgx/v5"
)

type stockRepository struct {
	querier Querier
}

func (r *stockRepository) GetByID(ctx context.Context, id int64) (domain.Stock, error) {
	q := `
		SELECT ` + stockColumns + `
		FROM stocks
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanStock,
		id,
	)
}

func (r *stockRepository) GetByPharmacyAndProduct(ctx context.Context, pharmacy_id int64, product_id int64) (domain.Stock, error) {
	q := `
		SELECT ` + stockColumns + `
		FROM stocks
		WHERE pharmacy_id = $1
			AND product_id = $2
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanStock,
		pharmacy_id, product_id,
	)
}

func (r *stockRepository) GetByIDAndLock(ctx context.Context, id int64) (domain.Stock, error) {
	q := `
		SELECT ` + stockColumns + `
		FROM stocks
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanStock,
		id,
	)
}

func (r *stockRepository) buildListQuery(sel string, det domain.StockListDetails) (*strings.Builder, pgx.NamedArgs) {
	var sb strings.Builder
	args := pgx.NamedArgs{}

	sb.WriteString(sel)
	sb.WriteString(`
		WHERE st.deleted_at IS NULL
	`)

	if det.PharmacySlug != nil {
		sb.WriteString(`
			AND ph.slug = @pharmacySlug
		`)
		args["pharmacySlug"] = *det.PharmacySlug
	}
	if det.ProductSlug != nil {
		sb.WriteString(`
			AND pd.slug = @productSlug
		`)
		args["productSlug"] = *det.ProductSlug
	}
	if det.ProductName != nil {
		sb.WriteString(`
			AND pd.name ILIKE '%' || @productName || '%'
		`)
		args["productName"] = *det.ProductName
	}
	if det.ManagerID != nil {
		sb.WriteString(`
			AND ph.manager_id = @managerID
		`)
		args["managerID"] = *det.ManagerID
	}

	return &sb, args
}

func (r *stockRepository) GetPageInfo(ctx context.Context, det domain.StockListDetails) (domain.PageInfo, error) {
	sb, args := r.buildListQuery(countStockJoined, det)

	return getPageInfo(
		r.querier, ctx, sb.String(),
		det.Page, det.Limit,
		args,
	)
}

func (r *stockRepository) List(ctx context.Context, det domain.StockListDetails) ([]domain.StockJoined, error) {
	sb, args := r.buildListQuery(selectStockJoined, det)
	offset := (det.Page - 1) * det.Limit

	sortCol := "pd.name"

	switch det.SortBy {
	case domain.StockSortByProductName:
		sortCol = "pd.name"
	case domain.StockSortByPharmacyName:
		sortCol = "ph.name"
	case domain.StockSortByPrice:
		sortCol = "st.price"
	case domain.StockSortByAmount:
		sortCol = "st.stock"
	}

	fmt.Fprintf(
		sb,
		` ORDER BY %s %s, st.id %s `,
		sortCol,
		getSortOrder(det.SortAsc),
		getSortOrder(det.SortAsc),
	)

	fmt.Fprintf(
		sb,
		` OFFSET %d LIMIT %d `,
		offset,
		det.Limit,
	)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanStockJoined,
		args,
	)
}

func (r *stockRepository) Add(ctx context.Context, s domain.Stock) (domain.Stock, error) {
	q := `
		INSERT INTO stocks (product_id, pharmacy_id, stock, price)
		VALUES
		($1, $2, $3, $4)
		RETURNING ` + stockColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanStock,
		s.ProductID, s.PharmacyID, s.Stock, s.Price,
	)
}

func (r *stockRepository) Update(ctx context.Context, s domain.Stock) (domain.Stock, error) {
	q := `
		UPDATE stocks
		SET stock = $2,
			price = $3,
			updated_at = now()
		WHERE id = $1
	`

	err := execOne(
		r.querier, ctx, q,
		s.ID, s.Stock, s.Price,
	)

	if err != nil {
		return domain.Stock{}, apperror.Wrap(err)
	}

	return s, nil
}

func (r *stockRepository) SoftDeleteByID(ctx context.Context, id int64) error {
	q := `
		UPDATE stocks
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}

func (r *stockRepository) GetMutationByID(ctx context.Context, id int64) (domain.StockMutation, error) {
	q := `
		SELECT ` + stockMutationColumns + `
		FROM stock_mutations
		WHERE id = $1
			AND deleted_at IS NULL
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanStockMutation,
		id,
	)
}

func (r *stockRepository) GetMutationByIDAndLock(ctx context.Context, id int64) (domain.StockMutation, error) {
	q := `
		SELECT ` + stockMutationColumns + `
		FROM stock_mutations
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanStockMutation,
		id,
	)
}

func (r *stockRepository) buildListMutationQuery(sel string, det domain.StockMutationListDetails) (*strings.Builder, pgx.NamedArgs) {
	var sb strings.Builder
	args := pgx.NamedArgs{}

	sb.WriteString(sel)
	sb.WriteString(`
		WHERE sm.deleted_at IS NULL
	`)

	if det.SourcePharmacySlug != nil {
		sb.WriteString(`
			AND ph1.slug = @sourcePharmacySlug
		`)
		args["sourcePharmacySlug"] = *det.SourcePharmacySlug
	}
	if det.TargetPharmacySlug != nil {
		sb.WriteString(`
			AND ph2.slug = @targetPharmacySlug
		`)
		args["targetPharmacySlug"] = *det.TargetPharmacySlug
	}
	if det.ProductSlug != nil {
		sb.WriteString(`
			AND pd.slug = @productSlug
		`)
		args["productSlug"] = *det.ProductSlug
	}
	if det.ProductName != nil {
		sb.WriteString(`
			AND pd.name ILIKE '%' || @productName || '%'
		`)
		args["productName"] = *det.ProductName
	}
	if det.Method != nil {
		sb.WriteString(`
			AND sm.method = @method
		`)
		args["method"] = *det.Method
	}
	if det.Status != nil {
		sb.WriteString(`
			AND sm.status = @status
		`)
		args["status"] = *det.Status
	}
	if det.ManagerID != nil {
		sb.WriteString(`
			AND (ph1.manager_id = @managerID OR ph2.manager_id = @managerID)
		`)
		args["managerID"] = *det.ManagerID
	}

	return &sb, args
}

func (r *stockRepository) GetMutationPageInfo(ctx context.Context, det domain.StockMutationListDetails) (domain.PageInfo, error) {
	sb, args := r.buildListMutationQuery(countStockMutationJoined, det)

	return getPageInfo(
		r.querier, ctx, sb.String(),
		det.Page, det.Limit,
		args,
	)
}

func (r *stockRepository) ListMutations(ctx context.Context, det domain.StockMutationListDetails) ([]domain.StockMutationJoined, error) {
	sb, args := r.buildListMutationQuery(selectStockMutationJoined, det)
	offset := (det.Page - 1) * det.Limit

	sortCol := "sm.created_at"

	switch det.SortBy {
	case domain.StockSortByProductName:
		sortCol = "pd.name"
	case domain.StockSortBySourcePharmacyName:
		sortCol = "ph1.name"
	case domain.StockSortByTargetPharmacyName:
		sortCol = "ph2.name"
	case domain.StockSortByAmount:
		sortCol = "sm.amount"
	}

	fmt.Fprintf(
		sb,
		` ORDER BY %s %s, sm.id %s `,
		sortCol,
		getSortOrder(det.SortAsc),
		getSortOrder(det.SortAsc),
	)

	fmt.Fprintf(
		sb,
		` OFFSET %d LIMIT %d `,
		offset,
		det.Limit,
	)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanStockMutationJoined,
		args,
	)
}

func (r *stockRepository) AddMutation(ctx context.Context, s domain.StockMutation) (domain.StockMutation, error) {
	q := `
		INSERT INTO stock_mutations (source_id, target_id, method, status, amount)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING ` + stockMutationColumns

	return queryOneFull(
		r.querier, ctx, q,
		scanStockMutation,
		s.SourceID, s.TargetID, s.Method, s.Status, s.Amount,
	)
}

func (r *stockRepository) UpdateMutation(ctx context.Context, s domain.StockMutation) (domain.StockMutation, error) {
	q := `
		UPDATE stock_mutations
		SET status = $2,
			updated_at = now()
		WHERE id = $1
	`

	err := execOne(
		r.querier, ctx, q,
		s.ID, s.Status,
	)

	if err != nil {
		return domain.StockMutation{}, apperror.Wrap(err)
	}

	return s, nil
}

func (r *stockRepository) SoftDeleteMutationByID(ctx context.Context, id int64) error {
	q := `
		UPDATE stock_mutations
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
	`

	return execOne(
		r.querier, ctx, q,
		id,
	)
}

func (r *stockRepository) GetNearestStockWithProduct(
	ctx context.Context,
	targetPharmacyID int64,
	productID int64,
	amount int,
) (domain.Stock, error) {
	q := `
		SELECT st.id, st.product_id, st.pharmacy_id, st.stock, st.price
		FROM stocks st
		WHERE st.deleted_at = NULL
			AND st.pharmacy_id != $1
			AND st.product_id = $2
			AND st.stock >= $3
	`

	return queryOneFull(
		r.querier, ctx, q,
		scanStock,
		targetPharmacyID, productID, amount,
	)
}
