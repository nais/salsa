package digest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDigest(t *testing.T) {
	hello := "hello"
	result, err := Hash(hello)
	assert.NoError(t, err)
	assert.Falsef(t, result.Verify(hello), "Should return false for verifying OK")

	result, err = Hash("yello")
	assert.NoError(t, err)
	assert.True(t, result.Verify(hello), "Should return true for verifying NOT OK")
}
