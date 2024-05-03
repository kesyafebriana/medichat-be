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

func (h *PharmacyHandler) GetPharmacies(ctx *gin.Context) {
	var q dto.PharmacyListQuery

	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	query, err := q.ToDetails()
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	p, pInfo, err := h.pharmacySrv.GetPharmacies(ctx, query)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewPharmaciesResponse(p, pInfo)),
	)
}

func (h *PharmacyHandler) GetPharmacyBySlug(ctx *gin.Context) {
	var uri dto.PharmacySlugParams

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	pharmacy, err := h.pharmacySrv.GetPharmacyBySlug(ctx, uri.Slug)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewPharmacyResponse(pharmacy)),
	)
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

func (h *PharmacyHandler) UpdatePharmacy(ctx *gin.Context) {
	var uri dto.PharmacySlugParams

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	var req dto.PharmacyUpdateRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	pharmacy, err := h.pharmacySrv.UpdatePharmacy(ctx, dto.PharmacyUpdateRequestToDetails(req, uri.Slug))

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewPharmacyResponse(pharmacy)))
}

func (h *PharmacyHandler) DeletePharmacy(ctx *gin.Context) {
	var uri dto.PharmacySlugParams

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.pharmacySrv.DeletePharmacyBySlug(ctx, uri.Slug)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusNoContent,
		nil,
	)
}

func (h *PharmacyHandler) GetPharmacyOperations(ctx *gin.Context) {
	var req dto.PharmacySlugParams

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	pharmacyOperations, err := h.pharmacySrv.GetOperationsBySlug(ctx, req.Slug)

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
	var uri dto.PharmacySlugParams

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
		reqEntity = append(reqEntity, dto.PharmacyOperationRequestToDetails(v, uri.Slug))
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
