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

	pharmacy, err := h.pharmacySrv.CreatePharmacy(ctx, domain.Pharmacy{
		Name: "Kesya",
		ManagerID: 1,
		// PharmacistName:    req.Data.PharmacistName,
		// PharmacistPhone:   req.Data.PharmacistPhone,
		// PharmacistLicense: req.Data.PharmacistLicense,
	})

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewPharmacyResponse(pharmacy)))
}
