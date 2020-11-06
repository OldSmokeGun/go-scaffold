package utils

import (
	"math/rand"
)

func RandomString(length int) string {
	if length <= 0 {
		return ""
	}

	source := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_=")
	random := make([]byte, length)

	for i := 0; i < length; i++ {
		random[i] = source[rand.Intn(len(source))]
	}

	return string(random)
}
