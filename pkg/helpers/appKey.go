package helpers

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomKey(klen int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	key := make([]byte, klen)
	for i := range key {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Handle the error appropriately (could be logged or returned)
			return ""
		}
		key[i] = charset[num.Int64()]
	}

	return string(key)
}
