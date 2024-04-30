package domain

import "context"

const (
	StockMutationAutomatic       = "automatic"
	StockMutationManual          = "manual"
	StockMutationStatusApproved  = "approved"
	StockMutationStatusPending   = "pending"
	StockMutationStatusCancelled = "cancelled"
)

type Stock struct {
	ID int64

	ProductID  int64
	PharmacyID int64

	Stock int
	Price int
}

type StockUpdateDetail struct {
	ID int64

	Stock *int
	Price *int
}

type StockMutation struct {
	ID int64

	SourceID int64
	TargetID int64

	Method string
	Status string

	Amount int
}

type StockTransferRequest struct {
	SourceID int64
	TargetID int64
	Amount   int
}

type StockListDetails struct {
	ProductID  *int64
	PharmacyID *int64

	SortBy  string
	SortAsc bool

	Page  int
	Limit int
}

type StockMutationListDetails struct {
	SourceID *int64
	TargetID *int64

	SourcePharmacyID *int64
	TargetPharmacyID *int64

	Method *string
	Status *string

	SortBy  string
	SortAsc bool

	Page  int
	Limit int
}

type StockRepository interface {
	GetByID(ctx context.Context, id int64) (Stock, error)
	GetByIDAndLock(ctx context.Context, id int64) (Stock, error)
	List(ctx context.Context, det StockListDetails) ([]Stock, error)

	Add(ctx context.Context, s Stock) (Stock, error)
	Update(ctx context.Context, s Stock) (Stock, error)
	SoftDeleteByID(ctx context.Context, id int64) error

	GetMutationByID(ctx context.Context, id int64) (StockMutation, error)
	GetMutationByIDAndLock(ctx context.Context, id int64) (StockMutation, error)
	ListMutations(ctx context.Context, det StockMutationListDetails) ([]StockMutation, error)

	AddMutation(ctx context.Context, s StockMutation) (StockMutation, error)
	UpdateMutation(ctx context.Context, s StockMutation) (StockMutation, error)
	SoftDeleteMutationByID(ctx context.Context, id int64) error
}

type StockService interface {
	GetByID(ctx context.Context, id int64) (Stock, error)
	List(ctx context.Context, det StockListDetails) ([]Stock, error)

	Add(ctx context.Context, s Stock) (Stock, error)
	Update(ctx context.Context, det StockUpdateDetail) (Stock, error)
	DeleteByID(ctx context.Context, id int64) error

	GetMutationByID(ctx context.Context, id int64) (StockMutation, error)
	ListMutations(ctx context.Context, det StockMutationListDetails) ([]StockMutation, error)

	RequestStockTransfer(ctx context.Context, r StockTransferRequest) (StockMutation, error)
	ApproveStockTransfer(ctx context.Context, id int64) (StockMutation, error)
	CancelStockTransfer(ctx context.Context, id int64) (StockMutation, error)
}
