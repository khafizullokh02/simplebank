package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate random numbers between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// random string of n charracters
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generate random owner name
func RandomOwner() string {
	return RandomString(6)
}

// random money generator
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generate a random currency
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
