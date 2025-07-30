package internal

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomCode() string {
	rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // 0 ~ 999999
}
