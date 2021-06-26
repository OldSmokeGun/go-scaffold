package pkg

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	if length <= 0 {
		return ""
	}

	source := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_=")
	random := make([]byte, length)

	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		random[i] = source[newRand.Intn(len(source))]
	}

	return string(random)
}
