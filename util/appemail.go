package util

import (
	"bytes"
	"html/template"

	"gopkg.in/gomail.v2"
)

type AppEmail interface {
	NewVerifyAccountEmail(fullname, email string, confirmURL string) *gomail.Message
	NewPasswordResetEmail(resetUrl string) *gomail.Message
}

type appEmail struct {
	verifyAccountTemplate *template.Template
	passwordResetTemplate *template.Template
}

func NewAppEmail() (*appEmail, error) {
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
	}, nil
}

func (a *appEmail) NewVerifyAccountEmail(fullname, email string, confirmURL string) *gomail.Message {
	var body bytes.Buffer
	a.verifyAccountTemplate.Execute(&body, struct {
		Fullname   string
		Email      string
		ConfirmURL string
	}{
		Fullname:   fullname,
		Email:      email,
		ConfirmURL: confirmURL,
	})
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "Welcome on Medichat")
	mailer.SetBody("text/html", body.String())
	return mailer
}

func (a *appEmail) NewPasswordResetEmail(resetUrl string) *gomail.Message {
	var body bytes.Buffer
	a.passwordResetTemplate.Execute(&body, struct {
		ResetURL string
	}{
		ResetURL: resetUrl,
	})
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "Medichat Account Password Reset")
	mailer.SetBody("text/html", body.String())
	return mailer
}
