package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GoogleHandler struct {
	googleSrv domain.GoogleService
}

type GoogleHandlerOpts struct {
	GoogleSrv domain.GoogleService
}

func NewGoogleHandler(opts GoogleHandlerOpts) *GoogleHandler {
	return &GoogleHandler{
		googleSrv: opts.GoogleSrv,
	}
}

func (h *GoogleHandler) OAuth2Callback(ctx *gin.Context) {
	var query dto.OAuth2CallbackQuery
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	tempState := session.Get(constants.SessionOAuth2State)
	state, ok := tempState.(string)
	if !ok {
		ctx.Error(apperror.NewTypeAssertionFailed(state, tempState))
		ctx.Abort()
		return
	}

	opts := query.ToOpts()

	token, err := h.googleSrv.OAuth2Callback(ctx, state, opts)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewAuthTokensResponse(token)))
}
