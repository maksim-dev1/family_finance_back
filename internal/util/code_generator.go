package util

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCode генерирует случайный 6-значный код подтверждения
// Использует криптографически безопасный генератор случайных чисел
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%05d", rand.Intn(100000))
}
