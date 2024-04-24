package handler

import (
	"log"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
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

	err := dto.ShouldBindMultipart(ctx, &req)
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

	log.Println(dets)

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

	err := dto.ShouldBindMultipart(ctx, &req)
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

	log.Println(dets)

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
