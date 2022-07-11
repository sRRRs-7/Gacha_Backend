package utils

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

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	sb.Grow(10)
	k := len(alphabet)

	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomItemName() string {
	return fmt.Sprintf("srrrs-%s", RandomString(5))
}

func RandomItemUrl() string {
	return fmt.Sprintf("http://srrrs/%s", RandomString(10))
}

func RandomCategory() string {
	return RandomString(6)
}