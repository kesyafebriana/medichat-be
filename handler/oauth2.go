package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/cryptoutil"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type OAuth2Handler struct {
	randomTokenProvider cryptoutil.RandomTokenProvider
	oAuth2Service       domain.OAuth2Service
}

type OAuth2HandlerOpts struct {
	RandomTokenProvider cryptoutil.RandomTokenProvider
	OAuth2Service       domain.OAuth2Service
}

func NewOAuth2Handler(opts OAuth2HandlerOpts) *OAuth2Handler {
	return &OAuth2Handler{
		randomTokenProvider: opts.RandomTokenProvider,
		oAuth2Service:       opts.OAuth2Service,
	}
}

func (h *OAuth2Handler) GetAuthURL(ctx *gin.Context) {
	session := sessions.Default(ctx)

	state, _ := h.randomTokenProvider.GenerateToken()
	session.Set(constants.SessionOAuth2State, state)
	err := session.Save()
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	url, err := h.oAuth2Service.GetAuthURL(ctx, state)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusSeeOther, dto.ResponseSeeOther(url))
}

func (h *OAuth2Handler) Callback(ctx *gin.Context) {
	var query dto.OAuth2CallbackQuery
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	_state := session.Get(constants.SessionOAuth2State)
	state, ok := _state.(string)
	if !ok {
		ctx.Error(apperror.NewTypeAssertionFailed(state, _state))
		ctx.Abort()
		return
	}

	opts := query.ToOpts()

	token, err := h.oAuth2Service.Callback(ctx, state, opts)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewAuthTokensResponse(token)))
}
