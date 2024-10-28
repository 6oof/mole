package helpers

import (
	"math/rand"
	"time"
)

func GenerateAppKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const klen = 32

	seed := rand.New(rand.NewSource(time.Now().Unix()))

	key := make([]byte, klen)
	for i := range key {
		key[i] = charset[seed.Intn(len(charset))]
	}

	return string(key)
}
