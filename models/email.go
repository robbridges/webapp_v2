package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@webgallery.com"
)

type EmailService struct {
	DefaultSender string

	dialer *mail.Dialer
}

type SMTPConfig struct {
	HOST     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.HOST, config.Port, config.Username, config.Password),
	}

	return &es
}
