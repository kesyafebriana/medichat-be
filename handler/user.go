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

type UserHandler struct {
	userSrv domain.UserService
}

type UserHandlerOpts struct {
	UserSrv domain.UserService
}

func NewUserHandler(opts UserHandlerOpts) *UserHandler {
	return &UserHandler{
		userSrv: opts.UserSrv,
	}
}

func (h *UserHandler) CreateProfile(ctx *gin.Context) {
	var req dto.UserCreateRequest

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

	dets, err := dto.UserCreateRequestToDetails(req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	_, err = h.userSrv.CreateProfile(ctx, dets)
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

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req dto.UserUpdateRequest

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

	dets, err := dto.UserUpdateRequestToDetails(req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	_, err = h.userSrv.UpdateProfile(ctx, dets)
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

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	profile, err := h.userSrv.GetProfile(ctx)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewProfileResponse(profile)),
	)
}

func (h *UserHandler) AddLocation(ctx *gin.Context) {
	var req dto.UserLocationCreateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	ul := req.ToEntity()

	ul, err = h.userSrv.AddLocation(ctx, ul)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewUserLocationResponse(ul)),
	)
}

func (h *UserHandler) UpdateLocation(ctx *gin.Context) {
	var req dto.UserLocationUpdateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	det := req.ToDetails()

	ul, err := h.userSrv.UpdateLocation(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(dto.NewUserLocationResponse(ul)),
	)
}

func (h *UserHandler) DeleteLocation(ctx *gin.Context) {
	var req dto.IDPathRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.userSrv.DeleteLocationByID(ctx, req.ID)
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
