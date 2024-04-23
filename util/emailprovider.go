package util

import (
	"gopkg.in/gomail.v2"
)

type EmailProvider interface {
	SendEmail(to string, message *gomail.Message) error
}

type emailProviderImpl struct {
	sender string
	dialer *gomail.Dialer
}

type EmailProviderOpts struct {
	Username    string
	Password    string
	EmailSender string
}

func NewGmailProvider(opts EmailProviderOpts) *emailProviderImpl {
	dialer := gomail.NewDialer("smtp.gmail.com", 587, opts.Username, opts.Password)
	return &emailProviderImpl{
		sender: opts.EmailSender,
		dialer: dialer,
	}
}

func (p *emailProviderImpl) SendEmail(to string, message *gomail.Message) error {
	message.SetHeader("From", p.sender)
	message.SetHeader("To", to)
	err := p.dialer.DialAndSend(message)
	if err != nil {
		return err
	}
	return nil
}
