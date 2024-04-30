package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PharmacyHandler struct {
	pharmacySrv domain.PharmacyService
	domain      string
}

type PharmacyHandlerOpts struct {
	PharmacySrv domain.PharmacyService
	Domain      string
}

func NewPharmacyHandler(opts PharmacyHandlerOpts) *PharmacyHandler {
	return &PharmacyHandler{
		pharmacySrv: opts.PharmacySrv,
		domain:      opts.Domain,
	}
}

func (h *PharmacyHandler) CreatePharmacy(ctx *gin.Context) {
	var req dto.PharmacyCreateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	pharmacy, err := h.pharmacySrv.CreatePharmacy(ctx, dto.PharmacyCreateToDetails(req))

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewPharmacyResponse(pharmacy)))
}

func (h *PharmacyHandler) GetPharmacyOperations(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	pharmacyOperations, err := h.pharmacySrv.GetOperationsById(ctx, req.ID)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewPharmacyOperationsResponse(pharmacyOperations)),
	)
}

func (h *PharmacyHandler) UpdatePharmacyOperations(ctx *gin.Context) {
	var reqEntity []domain.PharmacyOperationsUpdateDetails
	var uri dto.IDPathRequest

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	var req []dto.PharmacyOperationUpdateRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	for _, v := range req {
		reqEntity = append(reqEntity, dto.PharmacyOperationRequestToDetails(v, uri.ID))
	}

	newO, err := h.pharmacySrv.UpdateOperations(ctx, reqEntity)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewPharmacyOperationsResponse(newO)),
	)
}
