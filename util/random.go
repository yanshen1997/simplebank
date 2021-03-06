package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

var (
	currancy = []string{"USD", "ERU", "RMB"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return sb.String()
}

func GetRandomOwner() string {
	return RandomString(10)
}

func GetRandomBalance() int64 {
	return RandomInt(0, 1000)
}

func GetRandomCurrancy() string {
	return currancy[rand.Intn(len(currancy))]
}

func GetRandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(7))
}
