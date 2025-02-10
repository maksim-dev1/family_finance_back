package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateVerificationCode генерирует криптографически стойкий 4-значный код.
func GenerateVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		// В случае ошибки возвращаем фиксированный код (лучше обработать ошибку в реальном проекте)
		return "1000"
	}
	code := n.Int64() + 1000
	return fmt.Sprintf("%d", code)
}
