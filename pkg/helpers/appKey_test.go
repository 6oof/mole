package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAppKey(t *testing.T) {
	k := GenerateAppKey()

	assert.Len(t, k, 32, "The generated app key should be 32 characters long.")
}
