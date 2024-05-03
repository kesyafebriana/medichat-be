package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PharmacyManagerHandler struct {
	pharmacyManagerSrv domain.PharmacyManagerService
}

type PharmacyManagerHandlerOpts struct {
	PharmacyManagerSrv domain.PharmacyManagerService
}

func NewPharmacyManagerHandler(opts PharmacyManagerHandlerOpts) *PharmacyManagerHandler {
	return &PharmacyManagerHandler{
		pharmacyManagerSrv: opts.PharmacyManagerSrv,
	}
}

func (h *PharmacyManagerHandler) CreateProfile(ctx *gin.Context) {
	var req dto.AccountRegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	creds := req.ToCredentials()

	_, err = h.pharmacyManagerSrv.CreatePharmacyManager(ctx, creds)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(nil),
	)
}
