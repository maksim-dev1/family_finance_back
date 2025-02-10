package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// SendEmail отправляет письмо через SMTP.
func SendEmail(smtpHost string, smtpPort int, username, password, to, subject, body string) error {
	log.Printf("[INFO] Отправка email на %s через %s:%d", to, smtpHost, smtpPort)
	auth := smtp.PlainAuth("", username, password, smtpHost)
	from := username

	// Формирование заголовков сообщения
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=utf-8",
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Подключение с использованием TLS
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		log.Printf("[ERROR] Ошибка при подключении к SMTP серверу: %v", err)
		return fmt.Errorf("не удалось подключиться к SMTP серверу: %v", err)
	}
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("[ERROR] Ошибка при создании SMTP клиента: %v", err)
		return fmt.Errorf("не удалось создать SMTP клиента: %v", err)
	}

	// Аутентификация
	if err := client.Auth(auth); err != nil {
		log.Printf("[ERROR] Ошибка аутентификации: %v", err)
		return fmt.Errorf("не удалось пройти аутентификацию: %v", err)
	}

	// Установка отправителя и получателя
	if err := client.Mail(from); err != nil {
		log.Printf("[ERROR] Ошибка установки отправителя: %v", err)
		return fmt.Errorf("не удалось установить отправителя: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Printf("[ERROR] Ошибка установки получателя: %v", err)
		return fmt.Errorf("не удалось установить получателя: %v", err)
	}

	// Отправка письма
	w, err := client.Data()
	if err != nil {
		log.Printf("[ERROR] Ошибка при получении канала для данных: %v", err)
		return fmt.Errorf("не удалось получить канал для данных: %v", err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("[ERROR] Ошибка при записи сообщения в канал: %v", err)
		return fmt.Errorf("не удалось записать сообщение: %v", err)
	}
	err = w.Close()
	if err != nil {
		log.Printf("[ERROR] Ошибка при закрытии канала: %v", err)
		return fmt.Errorf("не удалось закрыть канал: %v", err)
	}

	// Завершаем соединение
	client.Quit()

	log.Println("[INFO] Email успешно отправлен")
	return nil
}
