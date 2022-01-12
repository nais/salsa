package digest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDigest(t *testing.T) {
	hello := "hello"
	result := Hash(hello)
	assert.Falsef(t, Verify(result, hello), "Should return false for verifying OK")

	result = Hash("yello")
	assert.True(t, Verify(result, hello), "Should return true for verifying NOT OK")
}
