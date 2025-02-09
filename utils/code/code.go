// utils/code/code.go
package code

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateVerificationCode генерирует случайный 4-значный код
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000 // случайное число от 1000 до 9999
	return fmt.Sprintf("%d", code)
}
