package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
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

	det := q.ToDetails()

	stocks, page, err := h.stockSrv.List(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(map[string]any{
			"page_info": dto.NewPageInfoResponse(page),
			"stocks":    util.MapSlice(stocks, dto.NewStockJoinedResponse),
		}),
	)
}

func (h *StockHandler) AddStock(ctx *gin.Context) {
	var req dto.StockCreateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	det := req.ToDetails()

	stock, err := h.stockSrv.Add(ctx, det)
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

func (h *StockHandler) GetStockByID(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	stock, err := h.stockSrv.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewStockResponse(stock)),
	)
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

	det := q.ToDetails()

	muts, page, err := h.stockSrv.ListMutations(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(map[string]any{
			"page_info":       dto.NewPageInfoResponse(page),
			"stock_mutations": util.MapSlice(muts, dto.NewStockMutationJoinedResponse),
		}),
	)
}

func (h *StockHandler) GetMutationByID(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	mut, err := h.stockSrv.GetMutationByID(ctx, req.ID)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewStockMutationResponse(mut)),
	)
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
