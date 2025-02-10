package utils

import (
	"fmt"
	"net/smtp"
)

// SendEmail отправляет письмо через SMTP.
func SendEmail(smtpHost string, smtpPort int, username, password, to, subject, body string) error {
	auth := smtp.PlainAuth("", username, password, smtpHost)
	from := username

	// Формирование заголовков сообщения
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `text/plain; charset="utf-8"`,
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message)); err != nil {
		return fmt.Errorf("не удалось отправить email: %v", err)
	}
	return nil
}
