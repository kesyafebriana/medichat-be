package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockSrv domain.StockService
}

type StockHandlerOpts struct {
	StockSrv domain.StockService
}

func NewStockHandler(opts StockHandlerOpts) *StockHandler {
	return &StockHandler{
		stockSrv: opts.StockSrv,
	}
}

func (h *StockHandler) ListStocks(ctx *gin.Context) {
	var q dto.StockListQuery

	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}
}

func (h *StockHandler) AddStock(ctx *gin.Context) {
	var req dto.StockUpdateRequest

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}
}

func (h *StockHandler) UpdateStock(ctx *gin.Context) {
	var req dto.StockUpdateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	det := req.ToDetails()

	stock, err := h.stockSrv.Update(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewStockResponse(stock)),
	)
}

func (h *StockHandler) DeleteStock(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.stockSrv.DeleteByID(ctx, req.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusNoContent,
		nil,
	)
}

func (h *StockHandler) ListMutations(ctx *gin.Context) {
	var q dto.StockMutationListQuery

	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

}

func (h *StockHandler) RequestTransfer(ctx *gin.Context) {
	var req dto.StockTransferRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	rq := req.ToRequest()

	mut, err := h.stockSrv.RequestStockTransfer(ctx, rq)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewStockMutationResponse(mut)),
	)
}

func (h *StockHandler) ApproveTransfer(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	mut, err := h.stockSrv.ApproveStockTransfer(ctx, req.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewStockMutationResponse(mut)),
	)
}

func (h *StockHandler) CancelTransfer(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	mut, err := h.stockSrv.CancelStockTransfer(ctx, req.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewStockMutationResponse(mut)),
	)
}
