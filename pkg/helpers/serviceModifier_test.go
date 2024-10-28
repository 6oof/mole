package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceModifierDecoratesCorrectly(t *testing.T) {
	mn := ServiceNameModifier("test-service", "testy")

	assert.Equal(t, mn, "mole-testy-test-service", "name should be modified correctly")
}
