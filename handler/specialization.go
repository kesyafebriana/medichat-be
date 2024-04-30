package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpecializationHandler struct {
	specializationSrv domain.SpecializationService
}

type SpecializationHandlerOpts struct {
	SpecializationSrv domain.SpecializationService
}

func NewSpecializationHandler(opts SpecializationHandlerOpts) *SpecializationHandler {
	return &SpecializationHandler{
		specializationSrv: opts.SpecializationSrv,
	}
}

func (h *SpecializationHandler) GetAll(ctx *gin.Context) {
	specializations, err := h.specializationSrv.GetAll(ctx)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(util.MapSlice(
			specializations,
			dto.NewSpecializationResponse,
		)),
	)
}
