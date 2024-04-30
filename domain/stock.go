package domain

import (
	"context"
	"time"
)

const (
	StockMutationAutomatic = "automatic"
	StockMutationManual    = "manual"

	StockMutationStatusApproved  = "approved"
	StockMutationStatusPending   = "pending"
	StockMutationStatusCancelled = "cancelled"

	StockSortByProductName  = "product_name"
	StockSortByPharmacyName = "pharmacy_name"
	StockSortByPrice        = "price"
	StockSortByAmount       = "amount"

	StockSortBySourcePharmacyName = "source_pharmacy_name"
	StockSortByTargetPharmacyName = "target_pharmacy_name"
)

type Stock struct {
	ID int64

	ProductID  int64
	PharmacyID int64

	Stock int
	Price int
}

type StockJoined struct {
	ID int64

	Product struct {
		ID   int64
		Slug string
		Name string
	}
	Pharmacy struct {
		ID   int64
		Name string
	}

	Stock int
	Price int
}

type StockCreateDetail struct {
	ProductSlug string
	PharmacyID  int64

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

	Timestamp time.Time
}

type StockMutationJoined struct {
	ID int64

	Source struct {
		ID           int64
		PharmacyID   int64
		PharmacyName string
	}
	Target struct {
		ID           int64
		PharmacyID   int64
		PharmacyName string
	}
	Product struct {
		ID   int64
		Slug string
		Name string
	}

	Method string
	Status string

	Amount int

	Timestamp time.Time
}

type StockTransferRequest struct {
	SourcePharmacyID int64
	TargetPharmacyID int64
	ProductSlug      string
	Amount           int
}

type StockListDetails struct {
	ProductSlug *string
	ProductName *string
	PharmacyID  *int64

	SortBy  string
	SortAsc bool

	Page  int
	Limit int
}

type StockMutationListDetails struct {
	ProductSlug *string
	ProductName *string

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
	GetByPharmacyAndProduct(ctx context.Context, pharmacy_id int64, product_id int64) (Stock, error)
	GetByIDAndLock(ctx context.Context, id int64) (Stock, error)
	GetPageInfo(ctx context.Context, det StockListDetails) (PageInfo, error)
	List(ctx context.Context, det StockListDetails) ([]StockJoined, error)

	Add(ctx context.Context, s Stock) (Stock, error)
	Update(ctx context.Context, s Stock) (Stock, error)
	SoftDeleteByID(ctx context.Context, id int64) error

	GetMutationByID(ctx context.Context, id int64) (StockMutation, error)
	GetMutationByIDAndLock(ctx context.Context, id int64) (StockMutation, error)
	GetMutationPageInfo(ctx context.Context, det StockMutationListDetails) (PageInfo, error)
	ListMutations(ctx context.Context, det StockMutationListDetails) ([]StockMutationJoined, error)

	AddMutation(ctx context.Context, s StockMutation) (StockMutation, error)
	UpdateMutation(ctx context.Context, s StockMutation) (StockMutation, error)
	SoftDeleteMutationByID(ctx context.Context, id int64) error
}

type StockService interface {
	GetByID(ctx context.Context, id int64) (Stock, error)
	List(ctx context.Context, det StockListDetails) ([]StockJoined, error)

	Add(ctx context.Context, s Stock) (Stock, error)
	Update(ctx context.Context, det StockUpdateDetail) (Stock, error)
	DeleteByID(ctx context.Context, id int64) error

	GetMutationByID(ctx context.Context, id int64) (StockMutation, error)
	ListMutations(ctx context.Context, det StockMutationListDetails) ([]StockMutationJoined, error)

	RequestStockTransfer(ctx context.Context, r StockTransferRequest) (StockMutation, error)
	ApproveStockTransfer(ctx context.Context, id int64) (StockMutation, error)
	CancelStockTransfer(ctx context.Context, id int64) (StockMutation, error)
}
