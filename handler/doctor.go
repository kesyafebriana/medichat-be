package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorSrv domain.DoctorService
}

type DoctorHandlerOpts struct {
	DoctorSrv domain.DoctorService
}

func NewDoctorHandler(opts DoctorHandlerOpts) *DoctorHandler {
	return &DoctorHandler{
		doctorSrv: opts.DoctorSrv,
	}
}

func (h *DoctorHandler) ListDoctors(ctx *gin.Context) {
	var q dto.DoctorListQuery

	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	det, err := q.ToDetails()
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	doctors, err := h.doctorSrv.List(ctx, det)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(util.MapSlice(doctors, func(d domain.Doctor) any {
			return dto.NewProfileResponse(d)
		})),
	)
}

func (h *DoctorHandler) CreateProfile(ctx *gin.Context) {
	var req dto.DoctorCreateRequest

	err := dto.ShouldBindMultipart(ctx, &req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	dets, err := dto.DoctorCreateRequestToDetails(req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	_, err = h.doctorSrv.CreateProfile(ctx, dets)
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

func (h *DoctorHandler) UpdateProfile(ctx *gin.Context) {
	var req dto.DoctorUpdateRequest

	err := dto.ShouldBindMultipart(ctx, &req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	dets, err := dto.DoctorUpdateRequestToDetails(req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	_, err = h.doctorSrv.UpdateProfile(ctx, dets)
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

func (h *DoctorHandler) GetProfile(ctx *gin.Context) {
	profile, err := h.doctorSrv.GetProfile(ctx)
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

func (h *DoctorHandler) SetActiveStatus(ctx *gin.Context) {
	var req dto.DoctorSetActiveRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.doctorSrv.SetActiveStatus(ctx, *req.IsActive)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(nil),
	)
}
