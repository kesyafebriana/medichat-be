package util

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type AppEmail interface {
	NewVerifyAccountEmail(fullname, email string, verifyEmailToken string) *gomail.Message
	NewPasswordResetEmail(email, resetPasswordToken string) *gomail.Message
}

type appEmail struct {
	verifyAccountTemplate *template.Template
	passwordResetTemplate *template.Template
	feVerificationURL     string
	feResetPasswordURL    string
}

type AppEmailOpts struct {
	FEVerivicationURL  string
	FEResetPasswordURL string
}

func NewAppEmail(opts AppEmailOpts) (*appEmail, error) {
	verifyAccountTemplate, err := template.ParseFiles("templates/verify-account-email.html")
	if err != nil {
		return nil, err
	}

	passwordResetTemplate, err := template.ParseFiles("templates/password-reset-email.html")
	if err != nil {
		return nil, err
	}

	return &appEmail{
		verifyAccountTemplate: verifyAccountTemplate,
		passwordResetTemplate: passwordResetTemplate,
		feVerificationURL:     opts.FEVerivicationURL,
		feResetPasswordURL:    opts.FEResetPasswordURL,
	}, nil
}

func (a *appEmail) NewVerifyAccountEmail(fullname, email, verifyEmailToken string) *gomail.Message {
	var body bytes.Buffer
	a.verifyAccountTemplate.Execute(&body, struct {
		Fullname   string
		Email      string
		ConfirmURL string
	}{
		Fullname:   fullname,
		Email:      email,
		ConfirmURL: fmt.Sprintf("%s?email=%s&verify_email_token=%s", a.feVerificationURL, email, verifyEmailToken),
	})
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "Welcome on Medichat")
	mailer.SetBody("text/html", body.String())
	return mailer
}

func (a *appEmail) NewPasswordResetEmail(email, resetPasswordToken string) *gomail.Message {
	var body bytes.Buffer
	a.passwordResetTemplate.Execute(&body, struct {
		ResetURL string
	}{
		ResetURL: fmt.Sprintf("%s?email=%s&reset_password_token=%s", a.feResetPasswordURL, email, resetPasswordToken),
	})
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "Medichat Account Password Reset")
	mailer.SetBody("text/html", body.String())
	return mailer
}
