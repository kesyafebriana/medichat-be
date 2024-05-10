package postgres

import (
	"context"
	"database/sql"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/repository/postgis"
	"time"
)

type dataRepository struct {
	conn    *sql.DB
	querier Querier
}

func NewDataRepository(db *sql.DB) *dataRepository {
	return &dataRepository{
		conn:    db,
		querier: db,
	}
}

func (r *dataRepository) Atomic(
	ctx context.Context,
	fn domain.AtomicFuncAny,
) (any, error) {
	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, apperror.Wrap(err)
	}
	defer tx.Rollback()

	txRepo := &dataRepository{
		conn:    r.conn,
		querier: tx,
	}

	ret, err := fn(txRepo)
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return ret, nil
}

func (r *dataRepository) Sleep(ctx context.Context, duration time.Duration) error {
	q := `
		SELECT pg_sleep($1)
	`

	return exec(
		r.querier, ctx, q,
		duration.Seconds(),
	)
}

func (r *dataRepository) GetDistance(ctx context.Context, a, b domain.Coordinate) (float64, error) {
	q := `
		SELECT ST_Distance($1::geography, $2::geography)
	`

	return queryOne(
		r.querier, ctx, q,
		float64ScanDest,
		postgis.NewPointFromCoordinate(a),
		postgis.NewPointFromCoordinate(b),
	)
}

func (r *dataRepository) ProductRepository() domain.ProductRepository {
	return &productRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) ProductDetailsRepository() domain.ProductDetailsRepository {
	return &productDetailRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) ChatRepository() domain.ChatRepository {
	return &chatRepository{
		querier: r.querier,
	}
}


func (r *dataRepository) AccountRepository() domain.AccountRepository {
	return &accountRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) RefreshTokenRepository() domain.RefreshTokenRepository {
	return &refreshTokenRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) ResetPasswordTokenRepository() domain.ResetPasswordTokenRepository {
	return &resetPasswordTokenRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) VerifyEmailTokenRepository() domain.VerifyEmailTokenRepository {
	return &verifyEmailTokenRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) CategoryRepository() domain.CategoryRepository {
	return &categoryRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) AdminRepository() domain.AdminRepository {
	return nil
}

func (r *dataRepository) UserRepository() domain.UserRepository {
	return &userRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) DoctorRepository() domain.DoctorRepository {
	return &doctorRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) PharmacyManagerRepository() domain.PharmacyManagerRepository {
	return &pharmacyManagerRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) SpecializationRepository() domain.SpecializationRepository {
	return &specializationRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) PharmacyRepository() domain.PharmacyRepository {
	return &pharmacyRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) ShipmentMethodRepository() domain.ShipmentMethodRepository {
	return &shipmentMethodRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) StockRepository() domain.StockRepository {
	return &stockRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) PaymentRepository() domain.PaymentRepository {
	return &paymentRepository{
		querier: r.querier,
	}
}

func (r *dataRepository) OrderRepository() domain.OrderRepository {
	return &orderRepository{
		querier: r.querier,
	}
}
