package smtp

import (
    "crypto/tls"
    "fmt"
    "net/smtp"
)

// Mailer отвечает за отправку email через SMTP.
type Mailer struct {
    host     string
    port     int
    username string
    password string
}

// NewMailer создаёт новый экземпляр Mailer.
func NewMailer(host string, port int, username, password string) *Mailer {
    return &Mailer{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }
}

// SendEmail отправляет email указанному получателю.
func (m *Mailer) SendEmail(to, subject, body string) error {
    from := m.username
    msg := "From: " + from + "\r\n" +
        "To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "MIME-version: 1.0;\r\n" +
        "Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
        body

    addr := fmt.Sprintf("%s:%d", m.host, m.port)
    // Настраиваем TLS (в продакшене рекомендуется проверка сертификата)
    tlsconfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         m.host,
    }

    // Устанавливаем TLS соединение
    conn, err := tls.Dial("tcp", addr, tlsconfig)
    if err != nil {
        return err
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, m.host)
    if err != nil {
        return err
    }
    defer client.Quit()

    auth := smtp.PlainAuth("", m.username, m.password, m.host)
    if err = client.Auth(auth); err != nil {
        return err
    }

    if err = client.Mail(from); err != nil {
        return err
    }
    if err = client.Rcpt(to); err != nil {
        return err
    }

    w, err := client.Data()
    if err != nil {
        return err
    }

    _, err = w.Write([]byte(msg))
    if err != nil {
        return err
    }

    err = w.Close()
    if err != nil {
        return err
    }

    return nil
}
