package handler

import (
	"medichat-be/domain"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderSrv domain.OrderService
}

type OrderHandlerOpts struct {
	OrderSrv domain.OrderService
}

func NewOrderHandler(opts OrderHandlerOpts) *OrderHandler {
	return &OrderHandler{
		orderSrv: opts.OrderSrv,
	}
}

func (h *OrderHandler) ListOrders(ctx *gin.Context) {

}

func (h *OrderHandler) GetOrderByID(ctx *gin.Context) {

}

func (h *OrderHandler) GetCartInfo(ctx *gin.Context) {

}

func (h *OrderHandler) AddOrder(ctx *gin.Context) {

}

func (h *OrderHandler) SendOrder(ctx *gin.Context) {

}

func (h *OrderHandler) FinishOrder(ctx *gin.Context) {

}

func (h *OrderHandler) CancelOrders(ctx *gin.Context) {

}
