package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentSrv domain.PaymentService
}

type PaymentHandlerOpts struct {
	PaymentSrv domain.PaymentService
}

func NewPaymentHandler(opts PaymentHandlerOpts) *PaymentHandler {
	return &PaymentHandler{
		paymentSrv: opts.PaymentSrv,
	}
}

func (h *PaymentHandler) ListPayments(ctx *gin.Context) {
	var q dto.PaymentListQuery

	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	det := q.ToDetails()

	payments, page, err := h.paymentSrv.List(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(map[string]any{
			"page_info": dto.NewPageInfoResponse(page),
			"payments":  util.MapSlice(payments, dto.NewPaymentResponse),
		}),
	)
}

func (h *PaymentHandler) GetPaymentByInvoiceNumber(ctx *gin.Context) {
	var uri dto.PaymentInvoiceNumberURI

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	payment, err := h.paymentSrv.GetByInvoiceNumber(ctx, uri.InvoiceNumber)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewPaymentResponse(payment)),
	)
}

func (h *PaymentHandler) UploadPayment(ctx *gin.Context) {
	var uri dto.PaymentInvoiceNumberURI
	var req dto.PaymentUploadRequest

	err := util.LimitContentLength(ctx, constants.MaxFileSize)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = dto.ShouldBindMultipart(ctx, &req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	f, err := req.Form.File.Open()
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.paymentSrv.UploadPayment(ctx, uri.InvoiceNumber, f)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(nil),
	)
}

func (h *PaymentHandler) ConfirmPayment(ctx *gin.Context) {
	var uri dto.PaymentInvoiceNumberURI

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.paymentSrv.ConfirmPayment(ctx, uri.InvoiceNumber)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(nil),
	)
}
