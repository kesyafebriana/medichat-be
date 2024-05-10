package dto

import (
	"medichat-be/domain"
	"medichat-be/util"
	"time"
)

type OrderResponse struct {
	ID int64 `json:"id"`

	User struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Pharmacy struct {
		ID        int64  `json:"id"`
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		ManagerID int64  `json:"manager_id"`
	} `json:"pharmacy"`
	Payment struct {
		ID            int64  `json:"id"`
		InvoiceNumber string `json:"invoice_number"`
	} `json:"payment"`
	ShipmentMethod struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"shipment_method"`

	Address    string        `json:"address"`
	Coordinate CoordinateDTO `json:"coordinate"`

	NItems      int `json:"n_items"`
	Subtotal    int `json:"subtotal"`
	ShipmentFee int `json:"shipment_fee"`
	Total       int `json:"total"`

	Status     string     `json:"status"`
	OrderedAt  time.Time  `json:"ordered_at"`
	FinishedAt *time.Time `json:"finished_at"`

	Items []OrderItemResponse `json:"items,omitempty"`
}

func NewOrderResponse(o domain.Order) OrderResponse {
	return OrderResponse{
		ID: o.ID,
		User: struct {
			ID   int64  "json:\"id\""
			Name string "json:\"name\""
		}(o.User),
		Pharmacy: struct {
			ID        int64  "json:\"id\""
			Slug      string "json:\"slug\""
			Name      string "json:\"name\""
			ManagerID int64  `json:"manager_id"`
		}(o.Pharmacy),
		Payment: struct {
			ID            int64  "json:\"id\""
			InvoiceNumber string "json:\"invoice_number\""
		}(o.Payment),
		ShipmentMethod: struct {
			ID   int64  "json:\"id\""
			Name string "json:\"name\""
		}(o.ShipmentMethod),
		Address:     o.Address,
		Coordinate:  CoordinateDTO(o.Coordinate),
		NItems:      o.NItems,
		Subtotal:    o.Subtotal,
		ShipmentFee: o.ShipmentFee,
		Total:       o.Total,
		Status:      o.Status,
		OrderedAt:   o.OrderedAt,
		FinishedAt:  o.FinishedAt,
		Items:       util.MapSlice(o.Items, NewOrderItemResponse),
	}
}

type OrderItemResponse struct {
	ID int64 `json:"id"`

	Product struct {
		ID             int64  `json:"id"`
		Slug           string `json:"slug"`
		Name           string `json:"name"`
		PhotoURL       string `json:"photo_url"`
		Classification string `json:"classification"`
	} `json:"product"`

	Price  int `json:"price"`
	Amount int `json:"amount"`
}

func NewOrderItemResponse(oi domain.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID: oi.ID,
		Product: struct {
			ID             int64  "json:\"id\""
			Slug           string "json:\"slug\""
			Name           string "json:\"name\""
			PhotoURL       string `json:"photo_url"`
			Classification string `json:"classification"`
		}{
			ID:             oi.Product.ID,
			Slug:           oi.Product.Slug,
			Name:           oi.Product.Name,
			PhotoURL:       oi.Product.PhotoURL,
			Classification: oi.Product.Classification,
		},
		Price:  oi.Price,
		Amount: oi.Amount,
	}
}

type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
}

func NewOrdersResponse(o domain.Orders) OrdersResponse {
	return OrdersResponse{
		Orders: util.MapSlice(o.Orders, NewOrderResponse),
		Total:  o.Total,
	}

}

type OrderListQuery struct {
	PharmacySlug *string `form:"pharmacy_slug"`
	Status       *string `form:"status"`

	Page  *int `form:"page" binding:"omitempty,min=1"`
	Limit *int `form:"limit" binding:"omitempty,min=1"`
}

func (q OrderListQuery) ToDetails() domain.OrderListDetails {
	ret := domain.OrderListDetails{
		PharmacySlug: q.PharmacySlug,
		Status:       q.Status,
		Page:         1,
		Limit:        10,
	}

	if q.Page != nil {
		ret.Page = *q.Page
	}
	if q.Limit != nil {
		ret.Limit = *q.Limit
	}

	return ret
}

type OrderItemCreateRequest struct {
	ProductSlug string `json:"product_slug"`
	Amount      int    `json:"amount"`
}

type OrderCreateRequest struct {
	PharmacySlug     string `json:"pharmacy_slug" binding:"required"`
	ShipmentMethodID int64  `json:"shipment_method_id" binding:"required"`

	Address    string        `json:"address" binding:"required"`
	Coordinate CoordinateDTO `json:"coordinate" binding:"required"`

	Items []OrderItemCreateRequest `json:"items" binding:"omitempty,dive,required"`
}

func (r OrderCreateRequest) ToDetails() domain.OrderCreateDetails {
	return domain.OrderCreateDetails{
		PharmacySlug:     r.PharmacySlug,
		ShipmentMethodID: r.ShipmentMethodID,
		Address:          r.Address,
		Coordinate:       r.Coordinate.ToCoordinate(),
		Items: util.MapSlice(r.Items, func(oi OrderItemCreateRequest) domain.OrderItemCreateDetails {
			return domain.OrderItemCreateDetails(oi)
		}),
	}
}

type OrdersCreateRequest struct {
	Orders []OrderCreateRequest `json:"orders" binding:"dive,required"`
}

func (r OrdersCreateRequest) ToDetails() []domain.OrderCreateDetails {
	return util.MapSlice(r.Orders, func(o OrderCreateRequest) domain.OrderCreateDetails {
		return o.ToDetails()
	})
}
