package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateVerificationCode генерирует криптографически стойкий 4-значный код.
func GenerateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return "", fmt.Errorf("ошибка генерации кода: %v", err)
	}

	code := n.Int64() + 1000
	return fmt.Sprintf("%d", code), nil
}
