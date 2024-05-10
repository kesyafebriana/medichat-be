package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
	"net/http"

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
	var q dto.OrderListQuery

	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	det := q.ToDetails()

	orders, page, err := h.orderSrv.List(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(map[string]any{
			"page_info": dto.NewPageInfoResponse(page),
			"orders":    util.MapSlice(orders, dto.NewOrderResponse),
		}),
	)
}

func (h *OrderHandler) GetOrderByID(ctx *gin.Context) {
	var uri dto.IDPathRequest

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	order, err := h.orderSrv.GetByID(ctx, uri.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewOrderResponse(order)),
	)
}

func (h *OrderHandler) GetCartInfo(ctx *gin.Context) {
	var req dto.OrdersCreateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	dets := req.ToDetails()

	orders, err := h.orderSrv.GetCartInfo(ctx, dets)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewOrdersResponse(orders)),
	)
}

func (h *OrderHandler) AddOrders(ctx *gin.Context) {
	var req dto.OrdersCreateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	dets := req.ToDetails()

	orders, err := h.orderSrv.AddOrders(ctx, dets)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewOrdersResponse(orders)),
	)
}

func (h *OrderHandler) SendOrder(ctx *gin.Context) {
	var uri dto.IDPathRequest

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.orderSrv.SendOrder(ctx, uri.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseOk(nil),
	)
}

func (h *OrderHandler) FinishOrder(ctx *gin.Context) {
	var uri dto.IDPathRequest

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.orderSrv.FinishOrder(ctx, uri.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseOk(nil),
	)
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	var uri dto.IDPathRequest

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.orderSrv.CancelOrder(ctx, uri.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseOk(nil),
	)
}
