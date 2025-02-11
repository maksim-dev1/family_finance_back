package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
)

// GenerateVerificationCode генерирует криптографически стойкий 4-значный код.
func GenerateVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		// В случае ошибки можно вернуть фиксированный код или обработать ошибку
		return "1000"
	}
	code := n.Int64() + 1000
	return fmt.Sprintf("%d", code)
}

// IsValidEmail проверяет корректность email-адреса.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
