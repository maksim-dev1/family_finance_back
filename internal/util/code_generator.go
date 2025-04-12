package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%05d", rand.Intn(100000))
}
