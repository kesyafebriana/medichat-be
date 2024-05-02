package domain

import (
	"context"
	"time"
)

const (
	OrderStatusWaitingPayment      = "Waiting for payment"
	OrderStatusWaitingConfirmation = "Waiting for payment confirmation"
	OrderStatusProcessed           = "Processed"
	OrderStatusSent                = "Sent"
	OrderStatusConfirmed           = "Order confirmed"
	OrderStatusCancelled           = "Cancelled"
)

type Order struct {
	ID int64

	User struct {
		ID   int64
		Name string
	}
	Pharmacy struct {
		ID   int64
		Slug string
		Name string
	}
	Payment struct {
		ID            int64
		InvoiceNumber string
	}
	ShipmentMethod struct {
		ID   int64
		Name string
	}

	Address    string
	Coordinate Coordinate

	NItems      int
	Subtotal    int
	ShipmentFee int
	Total       int

	Status     string
	OrderedAt  time.Time
	FinishedAt time.Time

	Items []OrderItem
}

type OrderItem struct {
	ID int64

	OrderID   int64
	ProductID int64
	Price     int

	ProductName string

	Amount int
}

type OrderListDetails struct {
	UserID            *int64
	PharmacyID        *int64
	PharmacyManagerID *int64
	Status            *string
}

type OrderCreateDetails struct {
}

type OrderRepository interface {
	GetPageInfo(ctx context.Context, dets OrderListDetails) (PageInfo, error)
	List(ctx context.Context, dets OrderListDetails) ([]Order, error)
	GetByID(ctx context.Context, id int64) (Order, error)
	GetByIDAndLock(ctx context.Context, id int64) (Order, error)

	Add(ctx context.Context, order Order) (Order, error)
	UpdateStatusByID(ctx context.Context, id int64, status string) error

	ListItemsByOrderID(ctx context.Context, id int64) ([]OrderItem, error)
	AddItem(ctx context.Context, item OrderItem) (OrderItem, error)
}

type OrderService interface {
	List(ctx context.Context, dets OrderListDetails) ([]Order, PageInfo, error)
	GetByID(ctx context.Context, id int64) (Order, error)

	AddOrder(ctx context.Context, dets OrderCreateDetails) (Order, error)
	SendOrder(ctx context.Context, id int64) error
	ConfirmSentOrder(ctx context.Context, id int64) error
	CancelOrder(ctx context.Context, id int64) error
}
