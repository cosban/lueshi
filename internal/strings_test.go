package internal

import (
	"testing"

	"github.com/cosban/assert"
)

func TestHash(t *testing.T) {
	result := Hash("bees", "honey")

	expected := "03b5c64d26016b15f14dac64a3deafccaa9ef72655ae09753c554e6b940ac5a2"

	assert.Equals(t, expected, result)
	assert.Equals(t, 64, len(result))
}
