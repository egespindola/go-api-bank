package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomString generates a random with n length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)])
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "BRL"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomMoney generates a random amount
func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}
