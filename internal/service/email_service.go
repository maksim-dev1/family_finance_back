package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"family_finance_back/config"

	"github.com/jordan-wright/email"
)

type EmailService interface {
	SendCode(to, code string) error
}

type emailService struct {
	cfg config.Config
}

func NewEmailService(cfg config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendCode(to, code string) error {
	e := email.NewEmail()
	e.From = s.cfg.SMTPUsername
	e.To = []string{to}
	e.Subject = "Ваш код авторизации"
	e.Text = []byte(fmt.Sprintf("Ваш код: %s\nКод действителен 90 секунд", code))

	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	// Устанавливаем TLS-конфигурацию
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.cfg.SMTPHost,
	}

	// Устанавливаем соединение с SMTP-сервером
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return err
	}
	defer c.Quit()

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = e.SendWithTLS(addr, auth, tlsConfig); err != nil {
		return err
	}

	return nil
}
