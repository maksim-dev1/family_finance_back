package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"myapp/config"
)

type EmailService interface {
	SendCode(to, code string) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendCode(to, code string) error {
	subject := "Ваш код подтверждения"
	body := fmt.Sprintf("Ваш код подтверждения: %s", code)
	msg := "From: " + s.cfg.SMTPUsername + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.cfg.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Auth(auth); err != nil {
		return err
	}
	if err = c.Mail(s.cfg.SMTPUsername); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
