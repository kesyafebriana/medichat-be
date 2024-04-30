package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountSrv domain.AccountService
	domain     string
}

type AccountHandlerOpts struct {
	AccountSrv domain.AccountService
	Domain     string
}

func NewAccountHandler(opts AccountHandlerOpts) *AccountHandler {
	return &AccountHandler{
		accountSrv: opts.AccountSrv,
		domain:     opts.Domain,
	}
}

func (h *AccountHandler) Register(ctx *gin.Context) {
	var req dto.AccountRegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	creds := req.ToCredentials()

	_, err = h.accountSrv.Register(ctx, creds)
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

func (h *AccountHandler) Login(ctx *gin.Context) {
	var req dto.AccountLoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	creds := req.ToCredentials()
	creds.ClientIP = ctx.ClientIP()

	tokens, err := h.accountSrv.Login(ctx, creds)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
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

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewAuthTokensResponse(tokens)),
	)
}

func (h *AccountHandler) ForgetPassword(ctx *gin.Context) {
	var req dto.AccountForgetPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	token, err := h.accountSrv.GetResetPasswordToken(ctx, req.Email)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(token),
	)
}

func (h *AccountHandler) ResetPassword(ctx *gin.Context) {
	var req dto.AccountResetPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	creds := req.ToCredentials()

	err = h.accountSrv.ResetPassword(ctx, creds)
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

func (h *AccountHandler) CheckResetPasswordToken(ctx *gin.Context) {
	var query dto.AccountCheckResetPasswordQuery

	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.accountSrv.CheckResetPasswordToken(ctx, query.Email, query.ResetPasswordToken)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(nil),
	)
}

func (h *AccountHandler) GetVerifyEmailToken(ctx *gin.Context) {
	var req dto.AccountGetVerifyEmailTokenRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	token, err := h.accountSrv.GetVerifyEmailToken(ctx, req.Email)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.ResponseCreated(token),
	)
}

func (h *AccountHandler) VerifyEmail(ctx *gin.Context) {
	var req dto.AccountVerifyEmailRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	creds := req.ToCredentials()

	err = h.accountSrv.VerifyEmail(ctx, creds)
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

func (h *AccountHandler) CheckVerifyEmailToken(ctx *gin.Context) {
	var query dto.AccountCheckVerifyEmailQuery

	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.accountSrv.CheckVerifyEmailToken(ctx, query.Email, query.VerifyEmailToken)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(nil),
	)
}

func (h *AccountHandler) RefreshTokens(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(constants.CookieRefreshToken)
	if err != nil {
		ctx.Error(apperror.NewUnauthorized(err))
		ctx.Abort()
		return
	}

	creds := domain.AccountRefreshTokensCredentials{
		RefreshToken: refreshToken,
		ClientIP:     ctx.ClientIP(),
	}

	tokens, err := h.accountSrv.RefreshTokens(ctx, creds)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
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

	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk(dto.NewAuthTokensResponse(tokens)),
	)
}

func (h *AccountHandler) GetProfile(ctx *gin.Context) {
	profile, err := h.accountSrv.GetProfile(ctx)
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
