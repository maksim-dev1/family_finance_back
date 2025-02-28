package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"family_finance_back/config"
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

	log.Printf("Подключение к SMTP-серверу %s", addr)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.cfg.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		log.Printf("Ошибка подключения через TLS: %v", err)
		return err
	}
	log.Printf("Установлено TLS-соединение с %s", addr)

	c, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		log.Printf("Ошибка создания SMTP клиента: %v", err)
		return err
	}
	defer c.Close()
	log.Printf("SMTP клиент успешно создан")

	if err = c.Auth(auth); err != nil {
		log.Printf("Ошибка аутентификации: %v", err)
		return err
	}
	log.Printf("Аутентификация пройдена")

	if err = c.Mail(s.cfg.SMTPUsername); err != nil {
		log.Printf("Ошибка отправки MAIL команды: %v", err)
		return err
	}
	log.Printf("MAIL команда отправлена")

	if err = c.Rcpt(to); err != nil {
		log.Printf("Ошибка отправки RCPT команды для %s: %v", to, err)
		return err
	}
	log.Printf("RCPT команда успешно принята для %s", to)

	w, err := c.Data()
	if err != nil {
		log.Printf("Ошибка вызова Data: %v", err)
		return err
	}
	log.Printf("Команда DATA принята, начинаем отправку сообщения")

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Printf("Ошибка записи данных: %v", err)
		return err
	}
	log.Printf("Сообщение успешно записано")

	if err = w.Close(); err != nil {
		log.Printf("Ошибка закрытия Data: %v", err)
		return err
	}
	log.Printf("Data успешно закрыто")

	if err = c.Quit(); err != nil {
		log.Printf("Ошибка завершения сессии: %v", err)
		return err
	}
	log.Printf("Сессия SMTP завершена успешно, сообщение отправлено")

	return nil
}
