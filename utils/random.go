package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// init initializes the random number generator.
//
//
// @return Returns the seed for the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer between min and max
//
// @param min - Lower bound of the random number
//
// @return A random integer between min
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n.
//
// @param n - Length of the string to generate.
//
// @return A randomly generated string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner.
//
//
// @return Random string of owner or
func RandomOwner() string {
	return RandomString(6)
}

// RandomCurrency returns a random currency.
//
//
// @return A randomly selected currency. EUR USD CAD
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomMoney returns a random number between 0 and 1000.
//
//
// @return The random number between 0 and 1000
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
