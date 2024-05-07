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

func (h *PharmacyManagerHandler) GetAll(ctx *gin.Context) {
	p, err := h.pharmacyManagerSrv.GetAll(ctx)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.NewPharmacyManagersAccountResponse(p))
}

func (h *PharmacyManagerHandler) CreateAccount(ctx *gin.Context) {
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
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(nil),
	)
}

func (h *PharmacyManagerHandler) CreateProfile(ctx *gin.Context) {
	var req dto.PharmacyManagerCreateRequest

	err := util.LimitContentLength(ctx, constants.MaxFileSize)
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

	dets, err := dto.PharmacyManagerCreateRequestToDetails(req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	_, err = h.pharmacyManagerSrv.CreateProfilePharmacyManager(ctx, dets)
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
