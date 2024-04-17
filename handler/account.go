package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountSrv domain.AccountService
}

type AccountHandlerOpts struct {
	AccountSrv domain.AccountService
}

func NewAccountHandler(opts AccountHandlerOpts) *AccountHandler {
	return &AccountHandler{
		accountSrv: opts.AccountSrv,
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

	tokens, err := h.accountSrv.Login(ctx, creds)
	if err != nil {
		ctx.Error(apperror.Wrap(err))
		ctx.Abort()
		return
	}

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
