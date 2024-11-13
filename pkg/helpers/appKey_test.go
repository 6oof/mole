package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAppKey(t *testing.T) {
	k := GenerateRandomKey(32)
	assert.Len(t, k, 32, "The generated app key should be 32 characters long.")

	k2 := GenerateRandomKey(8)
	assert.Len(t, k2, 8, "The generated app key should be 32 characters long.")

	k3 := GenerateRandomKey(8)
	assert.NotEqual(t, k2, k3, "Key is random")

}
