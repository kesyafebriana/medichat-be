package dto

import "medichat-be/domain"

type AccountLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=24"`
}

func (r *AccountLoginRequest) ToCredentials() domain.AccountLoginCredentials {
	return domain.AccountLoginCredentials{
		Email:    r.Email,
		Password: r.Password,
	}
}

type AccountRegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,account_role"`
}

func (r *AccountRegisterRequest) ToCredentials() domain.AccountRegisterCredentials {
	return domain.AccountRegisterCredentials{
		Account: domain.Account{
			Email:         r.Email,
			EmailVerified: false,
			Role:          r.Role,
			AccountType:   domain.AccountTypeRegular,
		},
	}
}

type AccountForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AccountGetVerifyEmailTokenRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AccountCheckResetPasswordQuery struct {
	Email              string `form:"email" binding:"required,email"`
	ResetPasswordToken string `form:"reset_password_token" binding:"required"`
}

type AccountCheckVerifyEmailQuery struct {
	Email            string `form:"email" binding:"required,email"`
	VerifyEmailToken string `form:"verify_email_token" binding:"required"`
}

type AccountResetPasswordRequest struct {
	Email              string `json:"email" binding:"required,email"`
	NewPassword        string `json:"new_password" binding:"required,password"`
	ResetPasswordToken string `json:"reset_password_token" binding:"required"`
}

func (r *AccountResetPasswordRequest) ToCredentials() domain.AccountResetPasswordCredentials {
	return domain.AccountResetPasswordCredentials{
		Email:              r.Email,
		NewPassword:        r.NewPassword,
		ResetPasswordToken: r.ResetPasswordToken,
	}
}

type AccountVerifyEmailRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,password"`
	VerifyEmailToken string `json:"verify_email_token" binding:"required"`
}

func (r *AccountVerifyEmailRequest) ToCredentials() domain.AccountVerifyEmailCredentials {
	return domain.AccountVerifyEmailCredentials{
		Email:            r.Email,
		Password:         r.Password,
		VerifyEmailToken: r.VerifyEmailToken,
	}
}

type AccountResponse struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
	AccountType   string `json:"account_type"`
}

func NewAccountResponse(u domain.Account) AccountResponse {
	return AccountResponse{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Role:          u.Role,
		AccountType:   u.AccountType,
	}
}
