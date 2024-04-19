package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GoogleHandler struct {
	googleSrv domain.GoogleService
	domain    string
}

type GoogleHandlerOpts struct {
	GoogleSrv domain.GoogleService
	Domain    string
}

func NewGoogleHandler(opts GoogleHandlerOpts) *GoogleHandler {
	return &GoogleHandler{
		googleSrv: opts.GoogleSrv,
		domain:    opts.Domain,
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
	opts.ClientIP = ctx.ClientIP()

	tokens, err := h.googleSrv.OAuth2Callback(ctx, state, opts)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.SetCookie(
		constants.CookieRefreshToken,
		tokens.RefreshToken,
		tokens.RefreshExpireAt.Second()-time.Now().Second(),
		"/",
		h.domain,
		false,
		true,
	)

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewAuthTokensResponse(tokens)))
}
